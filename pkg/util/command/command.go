package command

import (
	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	// https://www.pulumi.com/docs/intro/concepts/resources/options/customtimeouts/
	commandTimeout string = "20m"
)

func CopyFile(ctx *pulumi.Context, resourceName string,
	host pulumi.StringInput, username, privateKey pulumi.StringPtrInput, port pulumi.Float64PtrInput,
	localPath, remotePath string) (*remote.CopyFile, error) {
	return CopyFileWithDependencies(ctx, resourceName,
		host, username, privateKey, port,
		localPath, remotePath, []pulumi.Resource{})
}

func CopyFileWithDependencies(ctx *pulumi.Context, resourceName string,
	host pulumi.StringInput, username, privateKey pulumi.StringPtrInput, port pulumi.Float64PtrInput,
	localPath, remotePath string,
	dependecies []pulumi.Resource) (*remote.CopyFile, error) {
	return remote.NewCopyFile(ctx, resourceName, &remote.CopyFileArgs{
		Connection: remote.ConnectionArgs{
			Host:       host,
			PrivateKey: privateKey,
			User:       username,
			Port:       port,
		},
		LocalPath:  pulumi.String(localPath),
		RemotePath: pulumi.String(remotePath),
	},
		DefaultTimeouts(),
		pulumi.DependsOn(dependecies))
}

func RemoteExec(ctx *pulumi.Context, resourceName string,
	host pulumi.StringInput, username, privateKey pulumi.StringPtrInput, port pulumi.Float64PtrInput,
	remoteCommand pulumi.StringPtrInput, environment pulumi.StringMapInput) (*remote.Command, error) {
	return RemoteExecWithDependencies(ctx, resourceName,
		host, username, privateKey, port,
		remoteCommand, environment, []pulumi.Resource{})
}

func RemoteExecWithDependencies(ctx *pulumi.Context, resourceName string,
	host pulumi.StringInput, username, privateKey pulumi.StringPtrInput, port pulumi.Float64PtrInput,
	remoteCommand pulumi.StringPtrInput, environment pulumi.StringMapInput,
	dependecies []pulumi.Resource) (*remote.Command, error) {
	return remote.NewCommand(ctx, resourceName,
		&remote.CommandArgs{
			Connection: remote.ConnectionArgs{
				Host:       host,
				PrivateKey: privateKey,
				User:       username,
				Port:       port,
			},
			Create:      remoteCommand,
			Update:      remoteCommand,
			Environment: environment,
		},
		DefaultTimeouts(),
		pulumi.DependsOn(dependecies))
}

func DefaultTimeouts() pulumi.ResourceOption {
	return pulumi.Timeouts(
		&pulumi.CustomTimeouts{
			Create: commandTimeout,
			Update: commandTimeout,
			Delete: commandTimeout})
}
