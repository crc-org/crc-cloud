package manager

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	providerAPI "github.com/crc/crc-cloud/pkg/manager/provider/api"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
)

func validateParams(args map[string]string, required []string) error {
	var requiredMissing []string
	for _, r := range required {
		_, ok := args[r]
		if !ok {
			requiredMissing = append(requiredMissing, r)
		}
	}

	if len(requiredMissing) > 0 {
		return fmt.Errorf("required fields missing: %v", requiredMissing)
	}
	return nil
}

func upStack(targetStack providerAPI.Stack) (auto.UpResult, error) {
	ctx := context.Background()
	objectStack := getStack(ctx, targetStack)
	stdoutStreamer := optup.ProgressStreams(os.Stdout)
	return objectStack.Up(ctx, stdoutStreamer)
}

func destroyStack(targetStack providerAPI.Stack) (err error) {
	ctx := context.Background()
	objectStack := getStack(ctx, targetStack)
	stdoutStreamer := optdestroy.ProgressStreams(os.Stdout)
	if _, err = objectStack.Destroy(ctx, stdoutStreamer); err != nil {
		return
	}
	err = objectStack.Workspace().RemoveStack(ctx, targetStack.StackName)
	return
}

// this function gets our stack ready for update/destroy by prepping the workspace, init/selecting the stack
// and doing a refresh to make sure state and cloud resources are in sync
func getStack(ctx context.Context, target providerAPI.Stack) auto.Stack {
	// create or select a stack with an inline Pulumi program
	s, err := auto.UpsertStackInlineSource(ctx, target.StackName,
		target.ProjectName, target.DeployFunc, getOpts(target)...)
	if err != nil {
		logging.Errorf("%v", err)
		os.Exit(1)
	}
	if err = postStack(ctx, target, &s); err != nil {
		logging.Errorf("%v", err)
		os.Exit(1)
	}
	return s
}

func getOpts(target providerAPI.Stack) []auto.LocalWorkspaceOption {
	return []auto.LocalWorkspaceOption{
		auto.Project(workspace.Project{
			Name:    tokens.PackageName(target.ProjectName),
			Runtime: workspace.NewProjectRuntimeInfo("go", nil),
			Backend: &workspace.ProjectBackend{
				URL: target.BackedURL,
			},
		}),
		auto.WorkDir(filepath.Join(".")),
		// auto.SecretsProvider("awskms://alias/pulumi-secret-encryption"),
	}
}

func postStack(ctx context.Context, target providerAPI.Stack, stack *auto.Stack) (err error) {
	w := stack.Workspace()
	// for inline source programs, we must manage plugins ourselves
	if err = w.InstallPlugin(ctx, target.Plugin.Name, target.Plugin.Version); err != nil {
		return
	}
	_, err = stack.Refresh(ctx)
	return
}

func writeOutputs(stackResult auto.UpResult,
	destinationFolder string, results map[string]string) (err error) {
	for k, v := range results {
		if err = writeOutput(stackResult, k, destinationFolder, v); err != nil {
			return err
		}
	}
	return
}

func writeOutput(stackResult auto.UpResult, outputkey,
	destinationFolder, destinationFilename string) error {
	value, ok := stackResult.Outputs[outputkey].Value.(string)
	if ok {
		err := os.WriteFile(path.Join(destinationFolder, destinationFilename), []byte(value), 0600)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("output value %s not found", outputkey)
	}
	return nil
}
