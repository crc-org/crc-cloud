package extract

import (
	_ "embed"
	"fmt"

	"github.com/crc/crc-cloud/pkg/util"
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
		util.LocalExec(ctx,
			"extractScriptXRights",
			pulumi.String(fmt.Sprintf("chmod +x %s", scriptfileName)), nil)
	if err != nil {
		return nil, nil, err
	}
	execScriptENVS := map[string]string{
		"BUNDLE_DOWNLOAD_URL":     bundleDownloadURL,
		"SHASUMFILE_DOWNLOAD_URL": shasumfileDownloadURL}
	execScriptCommand, err :=
		util.LocalExecWithDependencies(ctx,
			"execExtractScript",
			pulumi.String(fmt.Sprintf(". %s", scriptfileName)),
			pulumi.ToStringMap(execScriptENVS),
			[]pulumi.Resource{scriptXRightsCommand})
	if err != nil {
		return nil, nil, err
	}
	bootKeyContentCommand, err :=
		util.LocalExecWithDependencies(ctx,
			"execBootKeyContent",
			pulumi.String(fmt.Sprintf("cat %s", bootKeyfilename)),
			nil,
			[]pulumi.Resource{execScriptCommand})
	if err != nil {
		return nil, nil, err
	}
	return bootKeyContentCommand, &bootKeyContentCommand.Stdout, nil
}
