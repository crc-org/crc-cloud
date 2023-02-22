package extract

import (
	_ "embed"
	"fmt"

	"github.com/crc/crc-cloud/pkg/util"
	"github.com/crc/crc-cloud/pkg/util/command"
	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	ExtractedDiskRawFileName = "disk.raw"
	bootKeyfilename          = "id_ecdsa"
)

//go:embed extract.sh
var script []byte

func Extract(ctx *pulumi.Context,
	bundleDownloadURL, shasumfileDownloadURL string) (pulumi.Resource, *pulumi.StringOutput, error) {

	// Write to temp file to be executed locally
	scriptfileName, err := util.WriteTempFile(string(script))
	if err != nil {
		return nil, nil, err
	}
	scriptXRightsCommand, err :=
		local.NewCommand(ctx, "extractScriptXRights",
			&local.CommandArgs{
				Create: pulumi.String(fmt.Sprintf("chmod +x %s", scriptfileName)),
			},
			command.DefaultTimeouts())
	if err != nil {
		return nil, nil, err
	}
	execScriptENVS := map[string]string{
		"BUNDLE_DOWNLOAD_URL":     bundleDownloadURL,
		"SHASUMFILE_DOWNLOAD_URL": shasumfileDownloadURL}
	execScriptCommand, err :=
		local.NewCommand(ctx, "execExtractScript",
			&local.CommandArgs{
				Create:      pulumi.String(fmt.Sprintf(". %s", scriptfileName)),
				Environment: pulumi.ToStringMap(execScriptENVS),
			},
			command.DefaultTimeouts(),
			pulumi.DependsOn([]pulumi.Resource{scriptXRightsCommand}))
	if err != nil {
		return nil, nil, err
	}
	bootKeyContentCommand, err :=
		local.NewCommand(ctx, "execBootKeyContent",
			&local.CommandArgs{
				Create: pulumi.String(fmt.Sprintf("cat %s", bootKeyfilename)),
			},
			command.DefaultTimeouts(),
			pulumi.DependsOn([]pulumi.Resource{execScriptCommand}))
	if err != nil {
		return nil, nil, err
	}
	return bootKeyContentCommand, &bootKeyContentCommand.Stdout, nil
}
