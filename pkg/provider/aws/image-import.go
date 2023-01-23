package aws

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/crc/crc-cloud/pkg/bundle"
	bundleExtract "github.com/crc/crc-cloud/pkg/bundle/extract"
	providerAPI "github.com/crc/crc-cloud/pkg/manager/provider/api"
	"github.com/crc/crc-cloud/pkg/util"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ebs"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	defaultVolumeSize int = 100
)

type importRequest struct {
	projectName           string
	bundleDownloadURL     string
	shasumfileDownloadURL string
}

func fillImportRequest(projectName, bundleDownloadURL, shasumfileDownloadURL string) (*importRequest, error) {
	return &importRequest{
		projectName:           projectName,
		bundleDownloadURL:     bundleDownloadURL,
		shasumfileDownloadURL: shasumfileDownloadURL,
	}, nil
}

func (r importRequest) runFunc(ctx *pulumi.Context) error {
	amiName, err := bundle.GetDescription(r.bundleDownloadURL)
	if err != nil {
		return err
	}
	// A random generated name is used for temporary assets
	// associated with the import process
	id := randomID()
	bucket, bucketResources, err := createTempBucket(ctx, id)
	if err != nil {
		return err
	}
	vmieRole, roleDependecy, err := createVMIEmportExportRole(ctx, id)
	if err != nil {
		return err
	}
	bundleExtractAssets, bootkey, err := bundleExtract.Extract(ctx, r.bundleDownloadURL, r.shasumfileDownloadURL)
	if err != nil {
		return err
	}

	// This code takes too long to upload the image
	// diskImageUpload, err := s3.NewBucketObject(ctx,
	// 	"crcCloudImporterDiskUpload",
	// 	&s3.BucketObjectArgs{
	// 		Key:    pulumi.String(bundleExtract.ExtractedDiskRawFileName),
	// 		Bucket: bucket.ID(),
	// 		Source: pulumi.NewFileAsset(bundleExtract.ExtractedDiskRawFileName),
	// 	},
	// 	pulumi.DependsOn([]pulumi.Resource{bundleExtractAssets}))

	// This option will use the underneath aws cli and s3 mulitpart upload
	diskImageUpload, err := uploadDisk(ctx, id, []pulumi.Resource{bundleExtractAssets})
	if err != nil {
		return err
	}
	// Register the AMI requires disk is uploaded to S3 bucket and policy for the bucket
	// is attached to the role
	ami, err := registerAMI(ctx, *amiName, bucket,
		vmieRole, []pulumi.Resource{diskImageUpload, bucketResources, roleDependecy})
	if err != nil {
		return err
	}
	ctx.Export(providerAPI.OutputBootKey, *bootkey)
	ctx.Export(providerAPI.OutputImageID, ami.ID())
	return nil
}

// random name for temporary assets requried for importing the image
func randomID() string {
	var letters = []rune(
		"abcdefghijklmnopqrstuvwxyz1234567890")
	b := make([]rune, 7)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return fmt.Sprintf("crc-cloud-%s", string(b))
}

// This function creates a temporary bucket to upload the disk image to be imported
// It returns the bucket resource, the generated bucket name and error if any
func createTempBucket(ctx *pulumi.Context, bucketName string) (*s3.BucketV2, pulumi.Resource, error) {
	bucket, err := s3.NewBucketV2(ctx,
		"crcCloudImporterTempBucket",
		&s3.BucketV2Args{
			Bucket: pulumi.String(bucketName),
			Tags: pulumi.StringMap{
				"Product": pulumi.String("crc-cloud"),
			},
		})
	if err != nil {
		return nil, nil, err
	}
	bucketACL, err := s3.NewBucketAclV2(ctx,
		"crcCloudImporterTempBucketACL",
		&s3.BucketAclV2Args{
			Bucket: bucket.Bucket,
			Acl:    pulumi.String("private"),
		})
	return bucket, bucketACL, err
}

// pulumi s3.NewBucketObject is really slow,
// we take advangate of aws cli to run this as a cmd
// as cli defaults upload with multipart 10
func uploadDisk(ctx *pulumi.Context, bucketName string,
	dependsOn []pulumi.Resource) (pulumi.Resource, error) {

	uploadCommand := fmt.Sprintf("aws s3 cp disk.raw s3://%s/disk.raw --only-show-errors", bucketName)
	return util.LocalExecWithDependencies(
		ctx,
		"crcCloudImporterDiskUploadByCli",
		pulumi.String(uploadCommand),
		nil,
		dependsOn)
}

// from an image as a raw on a s3 bucket this function will import it as a snapshot
// and the register the snapshot as an AMI
func registerAMI(ctx *pulumi.Context, amiName string,
	bucket *s3.BucketV2, vmieRole *iam.Role,
	dependsOn []pulumi.Resource) (*ec2.Ami, error) {
	snapshot, err := ebs.NewSnapshotImport(ctx,
		"crcCloudImporterSnapshotImport",
		&ebs.SnapshotImportArgs{
			DiskContainer: &ebs.SnapshotImportDiskContainerArgs{
				Format: pulumi.String("RAW"),
				UserBucket: &ebs.SnapshotImportDiskContainerUserBucketArgs{
					S3Bucket: bucket.Bucket,
					S3Key:    pulumi.String(bundleExtract.ExtractedDiskRawFileName),
				},
			},
			RoleName: vmieRole.Name,
			Tags: pulumi.StringMap{
				"Name": pulumi.String(amiName),
			},
		},
		pulumi.DependsOn(dependsOn),
		// This allows to mask the import operation with a create and destroy
		// keeping only the AMI the other resources are ephermeral only tied to
		// the execution
		pulumi.RetainOnDelete(true))
	if err != nil {
		return nil, err
	}
	return ec2.NewAmi(ctx,
		"crcCloudImporterAMIRegister",
		&ec2.AmiArgs{
			EbsBlockDevices: ec2.AmiEbsBlockDeviceArray{
				&ec2.AmiEbsBlockDeviceArgs{
					DeviceName: pulumi.String("/dev/xvda"),
					SnapshotId: snapshot.ID(),
					VolumeSize: pulumi.Int(defaultVolumeSize),
				},
			},
			Name:               pulumi.String(amiName),
			Description:        pulumi.String(amiName),
			RootDeviceName:     pulumi.String("/dev/xvda"),
			VirtualizationType: pulumi.String("hvm"),
			// Required by c6a instances
			EnaSupport: pulumi.Bool(true),
			Tags: pulumi.StringMap{
				"Name": pulumi.String(amiName),
			},
		},
		// This allows to mask the import operation with a create and destroy
		// keeping only the AMI the other resources are ephermeral only tied to
		// the execution
		pulumi.RetainOnDelete(true))
}

// https://docs.aws.amazon.com/vm-import/latest/userguide/required-permissions.html
func createVMIEmportExportRole(ctx *pulumi.Context,
	id string) (*iam.Role, pulumi.Resource, error) {
	tpJSON, err := trustPolicyContent()
	if err != nil {
		return nil, nil, err
	}
	role, err := iam.NewRole(ctx,
		"crcCloudImporterRole",
		&iam.RoleArgs{
			Name:             pulumi.String(id),
			AssumeRolePolicy: pulumi.String(*tpJSON),
		})
	if err != nil {
		return nil, nil, err
	}
	rolePolicy, err := rolePolicyContent(id)
	if err != nil {
		return nil, nil, err
	}
	rolePolicyAttachment, err := iam.NewRolePolicy(ctx,
		"rolePolicy",
		&iam.RolePolicyArgs{
			Role:   role.ID(),
			Policy: pulumi.String(*rolePolicy),
		})
	return role, rolePolicyAttachment, err
}

func trustPolicyContent() (*string, error) {
	tmpJSON0, err := json.Marshal(map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Sid":    "",
				"Effect": "Allow",
				"Principal": map[string]interface{}{
					"Service": "vmie.amazonaws.com",
				},
				"Action": "sts:AssumeRole",
				"Condition": map[string]interface{}{
					"StringEquals": map[string]interface{}{
						"sts:ExternalId": "vmimport",
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	json := string(tmpJSON0)
	return &json, nil
}

// TODO review s3 actions
func rolePolicyContent(bucketName string) (*string, error) {
	bucketNameARN := fmt.Sprintf("arn:aws:s3:::%s", bucketName)
	bucketNameRecursiveARN := fmt.Sprintf("arn:aws:s3:::%s/*", bucketName)
	tmpJSON0, err := json.Marshal(map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Effect": "Allow",
				"Action": []string{
					"s3:GetBucketLocation",
					"s3:GetObject",
					"s3:ListBucket",
				},
				"Resource": []string{
					bucketNameARN,
					bucketNameRecursiveARN,
				},
			},
			{
				"Effect": "Allow",
				"Action": []string{
					"ec2:ModifySnapshotAttribute",
					"ec2:CopySnapshot",
					"ec2:RegisterImage",
					"ec2:Describe*",
				},
				"Resource": "*",
			},
		},
	})
	if err != nil {
		return nil, err
	}
	json := string(tmpJSON0)
	return &json, nil
}
