// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package compute

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/internal"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// To get more information about Snapshot, see:
//
// * [API documentation](https://cloud.google.com/compute/docs/reference/rest/v1/snapshots)
// * How-to Guides
//   - [Official Documentation](https://cloud.google.com/compute/docs/disks/create-snapshots)
//
// ## Example Usage
//
// ```go
// package main
//
// import (
//
//	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/compute"
//	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
//
// )
//
//	func main() {
//		pulumi.Run(func(ctx *pulumi.Context) error {
//			_, err := compute.LookupSnapshot(ctx, &compute.LookupSnapshotArgs{
//				Name: pulumi.StringRef("my-snapshot"),
//			}, nil)
//			if err != nil {
//				return err
//			}
//			_, err = compute.LookupSnapshot(ctx, &compute.LookupSnapshotArgs{
//				Filter:     pulumi.StringRef("name != my-snapshot"),
//				MostRecent: pulumi.BoolRef(true),
//			}, nil)
//			if err != nil {
//				return err
//			}
//			return nil
//		})
//	}
//
// ```
func LookupSnapshot(ctx *pulumi.Context, args *LookupSnapshotArgs, opts ...pulumi.InvokeOption) (*LookupSnapshotResult, error) {
	opts = internal.PkgInvokeDefaultOpts(opts)
	var rv LookupSnapshotResult
	err := ctx.Invoke("gcp:compute/getSnapshot:getSnapshot", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getSnapshot.
type LookupSnapshotArgs struct {
	// A filter to retrieve the compute snapshot.
	// See [gcloud topic filters](https://cloud.google.com/sdk/gcloud/reference/topic/filters) for reference.
	// If multiple compute snapshot match, either adjust the filter or specify `mostRecent`. One of `name` or `filter` must be provided.
	Filter *string `pulumi:"filter"`
	// If `filter` is provided, ensures the most recent snapshot is returned when multiple compute snapshot match.
	//
	// ***
	MostRecent *bool `pulumi:"mostRecent"`
	// The name of the compute snapshot. One of `name` or `filter` must be provided.
	Name *string `pulumi:"name"`
	// The ID of the project in which the resource belongs.
	// If it is not provided, the provider project is used.
	Project *string `pulumi:"project"`
}

// A collection of values returned by getSnapshot.
type LookupSnapshotResult struct {
	ChainName         string  `pulumi:"chainName"`
	CreationTimestamp string  `pulumi:"creationTimestamp"`
	Description       string  `pulumi:"description"`
	DiskSizeGb        int     `pulumi:"diskSizeGb"`
	Filter            *string `pulumi:"filter"`
	// The provider-assigned unique ID for this managed resource.
	Id                       string                               `pulumi:"id"`
	LabelFingerprint         string                               `pulumi:"labelFingerprint"`
	Labels                   map[string]string                    `pulumi:"labels"`
	Licenses                 []string                             `pulumi:"licenses"`
	MostRecent               *bool                                `pulumi:"mostRecent"`
	Name                     *string                              `pulumi:"name"`
	Project                  *string                              `pulumi:"project"`
	SelfLink                 string                               `pulumi:"selfLink"`
	SnapshotEncryptionKeys   []GetSnapshotSnapshotEncryptionKey   `pulumi:"snapshotEncryptionKeys"`
	SnapshotId               int                                  `pulumi:"snapshotId"`
	SourceDisk               string                               `pulumi:"sourceDisk"`
	SourceDiskEncryptionKeys []GetSnapshotSourceDiskEncryptionKey `pulumi:"sourceDiskEncryptionKeys"`
	StorageBytes             int                                  `pulumi:"storageBytes"`
	StorageLocations         []string                             `pulumi:"storageLocations"`
	Zone                     string                               `pulumi:"zone"`
}

func LookupSnapshotOutput(ctx *pulumi.Context, args LookupSnapshotOutputArgs, opts ...pulumi.InvokeOption) LookupSnapshotResultOutput {
	return pulumi.ToOutputWithContext(context.Background(), args).
		ApplyT(func(v interface{}) (LookupSnapshotResult, error) {
			args := v.(LookupSnapshotArgs)
			r, err := LookupSnapshot(ctx, &args, opts...)
			var s LookupSnapshotResult
			if r != nil {
				s = *r
			}
			return s, err
		}).(LookupSnapshotResultOutput)
}

// A collection of arguments for invoking getSnapshot.
type LookupSnapshotOutputArgs struct {
	// A filter to retrieve the compute snapshot.
	// See [gcloud topic filters](https://cloud.google.com/sdk/gcloud/reference/topic/filters) for reference.
	// If multiple compute snapshot match, either adjust the filter or specify `mostRecent`. One of `name` or `filter` must be provided.
	Filter pulumi.StringPtrInput `pulumi:"filter"`
	// If `filter` is provided, ensures the most recent snapshot is returned when multiple compute snapshot match.
	//
	// ***
	MostRecent pulumi.BoolPtrInput `pulumi:"mostRecent"`
	// The name of the compute snapshot. One of `name` or `filter` must be provided.
	Name pulumi.StringPtrInput `pulumi:"name"`
	// The ID of the project in which the resource belongs.
	// If it is not provided, the provider project is used.
	Project pulumi.StringPtrInput `pulumi:"project"`
}

func (LookupSnapshotOutputArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*LookupSnapshotArgs)(nil)).Elem()
}

// A collection of values returned by getSnapshot.
type LookupSnapshotResultOutput struct{ *pulumi.OutputState }

func (LookupSnapshotResultOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*LookupSnapshotResult)(nil)).Elem()
}

func (o LookupSnapshotResultOutput) ToLookupSnapshotResultOutput() LookupSnapshotResultOutput {
	return o
}

func (o LookupSnapshotResultOutput) ToLookupSnapshotResultOutputWithContext(ctx context.Context) LookupSnapshotResultOutput {
	return o
}

func (o LookupSnapshotResultOutput) ChainName() pulumi.StringOutput {
	return o.ApplyT(func(v LookupSnapshotResult) string { return v.ChainName }).(pulumi.StringOutput)
}

func (o LookupSnapshotResultOutput) CreationTimestamp() pulumi.StringOutput {
	return o.ApplyT(func(v LookupSnapshotResult) string { return v.CreationTimestamp }).(pulumi.StringOutput)
}

func (o LookupSnapshotResultOutput) Description() pulumi.StringOutput {
	return o.ApplyT(func(v LookupSnapshotResult) string { return v.Description }).(pulumi.StringOutput)
}

func (o LookupSnapshotResultOutput) DiskSizeGb() pulumi.IntOutput {
	return o.ApplyT(func(v LookupSnapshotResult) int { return v.DiskSizeGb }).(pulumi.IntOutput)
}

func (o LookupSnapshotResultOutput) Filter() pulumi.StringPtrOutput {
	return o.ApplyT(func(v LookupSnapshotResult) *string { return v.Filter }).(pulumi.StringPtrOutput)
}

// The provider-assigned unique ID for this managed resource.
func (o LookupSnapshotResultOutput) Id() pulumi.StringOutput {
	return o.ApplyT(func(v LookupSnapshotResult) string { return v.Id }).(pulumi.StringOutput)
}

func (o LookupSnapshotResultOutput) LabelFingerprint() pulumi.StringOutput {
	return o.ApplyT(func(v LookupSnapshotResult) string { return v.LabelFingerprint }).(pulumi.StringOutput)
}

func (o LookupSnapshotResultOutput) Labels() pulumi.StringMapOutput {
	return o.ApplyT(func(v LookupSnapshotResult) map[string]string { return v.Labels }).(pulumi.StringMapOutput)
}

func (o LookupSnapshotResultOutput) Licenses() pulumi.StringArrayOutput {
	return o.ApplyT(func(v LookupSnapshotResult) []string { return v.Licenses }).(pulumi.StringArrayOutput)
}

func (o LookupSnapshotResultOutput) MostRecent() pulumi.BoolPtrOutput {
	return o.ApplyT(func(v LookupSnapshotResult) *bool { return v.MostRecent }).(pulumi.BoolPtrOutput)
}

func (o LookupSnapshotResultOutput) Name() pulumi.StringPtrOutput {
	return o.ApplyT(func(v LookupSnapshotResult) *string { return v.Name }).(pulumi.StringPtrOutput)
}

func (o LookupSnapshotResultOutput) Project() pulumi.StringPtrOutput {
	return o.ApplyT(func(v LookupSnapshotResult) *string { return v.Project }).(pulumi.StringPtrOutput)
}

func (o LookupSnapshotResultOutput) SelfLink() pulumi.StringOutput {
	return o.ApplyT(func(v LookupSnapshotResult) string { return v.SelfLink }).(pulumi.StringOutput)
}

func (o LookupSnapshotResultOutput) SnapshotEncryptionKeys() GetSnapshotSnapshotEncryptionKeyArrayOutput {
	return o.ApplyT(func(v LookupSnapshotResult) []GetSnapshotSnapshotEncryptionKey { return v.SnapshotEncryptionKeys }).(GetSnapshotSnapshotEncryptionKeyArrayOutput)
}

func (o LookupSnapshotResultOutput) SnapshotId() pulumi.IntOutput {
	return o.ApplyT(func(v LookupSnapshotResult) int { return v.SnapshotId }).(pulumi.IntOutput)
}

func (o LookupSnapshotResultOutput) SourceDisk() pulumi.StringOutput {
	return o.ApplyT(func(v LookupSnapshotResult) string { return v.SourceDisk }).(pulumi.StringOutput)
}

func (o LookupSnapshotResultOutput) SourceDiskEncryptionKeys() GetSnapshotSourceDiskEncryptionKeyArrayOutput {
	return o.ApplyT(func(v LookupSnapshotResult) []GetSnapshotSourceDiskEncryptionKey { return v.SourceDiskEncryptionKeys }).(GetSnapshotSourceDiskEncryptionKeyArrayOutput)
}

func (o LookupSnapshotResultOutput) StorageBytes() pulumi.IntOutput {
	return o.ApplyT(func(v LookupSnapshotResult) int { return v.StorageBytes }).(pulumi.IntOutput)
}

func (o LookupSnapshotResultOutput) StorageLocations() pulumi.StringArrayOutput {
	return o.ApplyT(func(v LookupSnapshotResult) []string { return v.StorageLocations }).(pulumi.StringArrayOutput)
}

func (o LookupSnapshotResultOutput) Zone() pulumi.StringOutput {
	return o.ApplyT(func(v LookupSnapshotResult) string { return v.Zone }).(pulumi.StringOutput)
}

func init() {
	pulumi.RegisterOutputType(LookupSnapshotResultOutput{})
}