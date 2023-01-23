package setup

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/crc/crc-cloud/pkg/bundle"
	"github.com/crc/crc-cloud/pkg/util"
	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

//go:embed clustersetup.sh
var script []byte

type SetupData struct {
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
	dependencies := []pulumi.Resource{}
	// Load pull secret content
	privateKey, err := ioutil.ReadFile(bootingPrivateKeyFilePath)
	if err != nil {
		return nil, err
	}

	c := remote.ConnectionArgs{
		Host:       publicIP,
		PrivateKey: pulumi.String(privateKey),
		User:       pulumi.String(bundle.ImageUsername),
		Port:       pulumi.Float64(bundle.ImageSSHPort),
	}

	_ = newPublicKey.ApplyT(
		func(pubKey string) (string, error) {
			pubKeyFilename, err := util.WriteTempFile(pubKey)
			if err != nil {
				return "", err
			}
			pubKeyRemoteCopyResource, err :=
				remote.NewCopyFile(ctx, "uploadNewPublicKey",
					&remote.CopyFileArgs{
						Connection: c,
						LocalPath:  pulumi.String(pubKeyFilename),
						RemotePath: pulumi.String("id_rsa.pub"),
					})
			if err != nil {
				return "", err
			}
			overrideKeyCommand := "cat /home/core/id_rsa.pub >> /home/core/.ssh/authorized_keys"
			_, err =
				remote.NewCommand(ctx, "addPublicKeyAsAuthorized",
					&remote.CommandArgs{
						Connection: c,
						Create:     pulumi.String(overrideKeyCommand),
					},
					pulumi.DependsOn([]pulumi.Resource{pubKeyRemoteCopyResource}))
			if err != nil {
				return "", err
			}
			return "", nil
		})
	return dependencies, nil
}

func Setup(ctx *pulumi.Context,
	publicIP, privateKey *pulumi.StringOutput,
	data SetupData) ([]pulumi.Resource, error) {
	dependencies := []pulumi.Resource{}

	// Load pull secret content
	pullsecret, err := ioutil.ReadFile(data.OCPPullSecretFilePath)
	if err != nil {
		return nil, err
	}
	pullSecretEncoded := base64.StdEncoding.EncodeToString([]byte(pullsecret))

	c := remote.ConnectionArgs{
		Host:       publicIP,
		PrivateKey: privateKey,
		User:       pulumi.String(bundle.ImageUsername),
		Port:       pulumi.Float64(bundle.ImageSSHPort),
	}

	_ = pulumi.All(
		data.Password,
		data.PublicIP,
		data.PrivateIP).ApplyT(
		func(args []interface{}) (string, error) {

			execScriptENVS := map[string]string{
				"IIP":            args[2].(string),
				"EIP":            args[1].(string),
				"PULL_SECRET":    pullSecretEncoded,
				"PASS_DEVELOPER": args[0].(string),
				"PASS_KUBEADMIN": args[0].(string),
				"PASS_REDHAT":    args[0].(string)}

			clusterSetupfileName, err := util.WriteTempFile(string(script))
			if err != nil {
				return "", err
			}
			clusterSetupRemoteCopyResource, err :=
				remote.NewCopyFile(ctx, "uploadClusterSetupScript",
					&remote.CopyFileArgs{
						Connection: c,
						LocalPath:  pulumi.String(clusterSetupfileName),
						RemotePath: pulumi.String("/var/home/core/cluster_setup.sh"),
					})
			if err != nil {
				return "", err
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
				return "", err
			}
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
			execClusterSetupCommand = fmt.Sprintf("%s %s",
				execClusterSetupCommand,
				"/var/home/core/cluster_setup.sh")

			execClusterSetup, err :=
				remote.NewCommand(ctx, "runClusterSetupScript",
					&remote.CommandArgs{
						Connection: c,
						Create:     pulumi.String(execClusterSetupCommand),
						// https://github.com/pulumi/pulumi-command/issues/48
						// Environment: pulumi.ToStringMap(execScriptENVS),
					},
					pulumi.DependsOn([]pulumi.Resource{scriptXRightsCommandResource}))
			if err != nil {
				execClusterSetup.Stderr.ApplyT(func(err string) error {
					return ctx.Log.Error(err, nil)
				})
				return "", err
			}
			return "", nil
			// getKCCmd := ("cat /opt/kubeconfig")
			// getKC, err :=
			// 	remote.NewCommand(ctx, "getKCCmd",
			// 		&remote.CommandArgs{
			// 			Connection: c,
			// 			Create:     pulumi.String(getKCCmd)},
			// 		pulumi.DependsOn([]pulumi.Resource{execClusterSetup}))
			// if err != nil {
			// 	return "", err
			// }

			// return getKC.Stdout.ApplyT(func(kubeconfig string) (string, error) {
			// 	return kubeconfig, nil
			// }), nil
		}).(pulumi.StringOutput)
	// _= pulumi.All(data.PublicIP,kubeconfig).ApplyT(
	// 	func(args []interface{}) (string, error) {
	// 	}
	// )

	// )
	// // pulumi.StringPtrInput
	// // pulumi.StringOutput
	// kubeconfig.
	return dependencies, nil
}
