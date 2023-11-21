package setup

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/crc/crc-cloud/pkg/bundle"
	"github.com/crc/crc-cloud/pkg/util"
	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

//go:embed clustersetup.sh
var script []byte

type Data struct {
	PrivateIP             *pulumi.StringOutput
	PublicIP              *pulumi.StringOutput
	OCPPullSecretFilePath string
	Password              *pulumi.StringOutput
	// DeveloperPassword     pulumi.StringOutput
	// KubeadminPassword     pulumi.StringOutput
	// RedHatPassword        pulumi.StringOutput
}

// Switch the fixed initial key with self created one
func SwapKeys(ctx *pulumi.Context, publicIP *pulumi.StringOutput,
	bootingPrivateKeyFilePath string, newPublicKey *pulumi.StringOutput) ([]pulumi.Resource, error) {
	// Load pull secret content
	privateKey, err := os.ReadFile(bootingPrivateKeyFilePath)
	if err != nil {
		return nil, err
	}

	c := remote.ConnectionArgs{
		Host:           publicIP,
		PrivateKey:     pulumi.String(privateKey),
		User:           pulumi.String(bundle.ImageUsername),
		Port:           pulumi.Float64(bundle.ImageSSHPort),
		DialErrorLimit: pulumi.Int(-1),
		PerDialTimeout: pulumi.Int(0),
	}

	pubKeyFilename := newPublicKey.ApplyT(util.WriteTempFile).(pulumi.StringOutput)
	pubKeyRemoteCopyResource, err :=
		remote.NewCopyFile(ctx, "uploadNewPublicKey",
			&remote.CopyFileArgs{
				Connection: c,
				LocalPath:  pubKeyFilename,
				RemotePath: pulumi.String("id_rsa.pub"),
			},
			pulumi.IgnoreChanges([]string{"localPath"}))
	if err != nil {
		return nil, err
	}
	overrideKeyCommand := "cat /home/core/id_rsa.pub >> /home/core/.ssh/authorized_keys"
	overrideKeyResource, err :=
		remote.NewCommand(ctx, "addPublicKeyAsAuthorized",
			&remote.CommandArgs{
				Connection: c,
				Create:     pulumi.String(overrideKeyCommand),
			},
			pulumi.DependsOn([]pulumi.Resource{pubKeyRemoteCopyResource}))
	if err != nil {
		return nil, err
	}
	return []pulumi.Resource{overrideKeyResource}, nil
}

func Setup(ctx *pulumi.Context,
	publicIP, privateKey *pulumi.StringOutput,
	data Data) (pulumi.StringPtrInput, []pulumi.Resource, error) {
	// Load pull secret content
	pullsecret, err := os.ReadFile(data.OCPPullSecretFilePath)
	if err != nil {
		return nil, nil, err
	}
	pullSecretEncoded := base64.StdEncoding.EncodeToString([]byte(pullsecret))

	c := remote.ConnectionArgs{
		Host:           publicIP,
		PrivateKey:     privateKey,
		User:           pulumi.String(bundle.ImageUsername),
		Port:           pulumi.Float64(bundle.ImageSSHPort),
		DialErrorLimit: pulumi.Int(-1),
		PerDialTimeout: pulumi.Int(0)}

	clusterSetupfileName, err := util.WriteTempFile(string(script))
	if err != nil {
		return nil, nil, err
	}
	clusterSetupRemoteCopyResource, err :=
		remote.NewCopyFile(ctx, "uploadClusterSetupScript",
			&remote.CopyFileArgs{
				Connection: c,
				LocalPath:  pulumi.String(clusterSetupfileName),
				RemotePath: pulumi.String("/var/home/core/cluster_setup.sh"),
			},
			pulumi.IgnoreChanges([]string{"localPath"}))
	if err != nil {
		return nil, nil, err
	}
	scriptXRightsCommand := "chmod +x /var/home/core/cluster_setup.sh"
	scriptXRightsCommandResource, err :=
		remote.NewCommand(ctx, "setXRightsForClusterSetupScript",
			&remote.CommandArgs{
				Connection: c,
				Create:     pulumi.String(scriptXRightsCommand),
			},
			pulumi.DependsOn([]pulumi.Resource{clusterSetupRemoteCopyResource}))
	if err != nil {
		return nil, nil, err
	}
	execClusterSetupCommand := pulumi.All(
		data.Password,
		data.PublicIP,
		data.PrivateIP).ApplyT(
		func(args []interface{}) string {

			execScriptENVS := map[string]string{
				"IIP":            args[2].(string),
				"EIP":            args[1].(string),
				"PULL_SECRET":    pullSecretEncoded,
				"PASS_DEVELOPER": args[0].(string),
				"PASS_KUBEADMIN": args[0].(string),
				"PASS_REDHAT":    args[0].(string)}

			// https://github.com/pulumi/pulumi-command/issues/48
			// using Environment from remote Command would require customice
			// ssh server config + restart it. So we workaround it just adding the
			// envs from the map to the cmd but keep it as map in case it could be solved
			// by any other way
			execClusterSetupCommand := "sudo "
			for k, v := range execScriptENVS {
				execClusterSetupCommand =
					fmt.Sprintf("%s %s=\"%s\"", execClusterSetupCommand, k, v)
			}
			return fmt.Sprintf("%s %s",
				execClusterSetupCommand,
				"/var/home/core/cluster_setup.sh")
		}).(pulumi.StringOutput)
	execClusterSetup, err :=
		remote.NewCommand(ctx, "runClusterSetupScript",
			&remote.CommandArgs{
				Connection: c,
				Create:     execClusterSetupCommand,
				// https://github.com/pulumi/pulumi-command/issues/48
				// Environment: pulumi.ToStringMap(execScriptENVS),
			},
			pulumi.IgnoreChanges([]string{"create"}),
			pulumi.DependsOn([]pulumi.Resource{scriptXRightsCommandResource}))
	if err != nil {
		execClusterSetup.Stderr.ApplyT(func(err string) error {
			return ctx.Log.Error(err, nil)
		})
		return nil, nil, err
	}
	// Once the cluster setup is comleted we
	// get the kubeconfig file from the host running the cluster
	// then we replace the internal access with the public IP
	// the resulting kubeconfig file can be used to access the cluster
	getKCCmd := ("cat /opt/kubeconfig")
	getKC, err :=
		remote.NewCommand(ctx, "getKCCmd",
			&remote.CommandArgs{
				Connection: c,
				Create:     pulumi.String(getKCCmd)},
			pulumi.DependsOn([]pulumi.Resource{execClusterSetup}))
	if err != nil {
		return nil, nil, err
	}
	kubeconfig := pulumi.All(getKC.Stdout, data.PublicIP).ApplyT(
		func(args []interface{}) string {
			return strings.ReplaceAll(args[0].(string),
				"https://api.crc.testing:6443",
				fmt.Sprintf("https://api.%s.nip.io:6443", args[1].(string)))
		}).(pulumi.StringOutput)
	return kubeconfig, []pulumi.Resource{execClusterSetup}, nil
}
