package deployment

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/t-kikuc/pipecd-plugin-prototypes/lambroll/cli"
	config "github.com/t-kikuc/pipecd-plugin-prototypes/lambroll/config"
	"github.com/t-kikuc/pipecd-plugin-prototypes/lambroll/toolregistry"
)

type deployExecutor struct {
	appDir       string
	lambrollPath string
	input        config.LambrollDeploymentInput
	slp          sdk.StageLogPersister
}

func (e *deployExecutor) initLambrollCommand(ctx context.Context) (cmd *cli.Lambroll, ok bool) {
	cmd = cli.NewLambroll(
		e.lambrollPath,
		e.appDir,
		e.input.FunctionFile,
		e.input.Source,
	)

	if ok := showUsingVersion(ctx, cmd, e.slp); !ok {
		return nil, false
	}

	return cmd, true
}

func executeStage(ctx context.Context, dtCfgs []*sdk.DeployTarget[config.LambrollDeployTargetConfig], input *sdk.ExecuteStageInput[config.LambrollDeploymentInput]) (sdk.StageStatus, error) {
	if len(dtCfgs) != 1 {
		return sdk.StageStatusFailure, status.Error(codes.InvalidArgument, "the number of deploy target must be one for this plugin.")
	}

	e := &deployExecutor{
		input:  *input.Request.TargetDeploymentSource.ApplicationConfig.Spec,
		slp:    input.Client.LogPersister(),
		appDir: input.Request.TargetDeploymentSource.ApplicationDirectory,
	}
	toolRegistry := toolregistry.NewRegistry(input.Client.ToolRegistry())
	p, err := toolRegistry.Lambroll(ctx, dtCfgs[0].Config.Version)
	if err != nil {
		return sdk.StageStatusFailure, err
	}
	e.lambrollPath = p

	switch input.Request.StageName {
	case stageDeploy:
		return e.ensureSync(ctx), nil
	case stageDiff:
		return e.ensureDiff(ctx), nil
	case stageRollback:
		e.appDir = input.Request.RunningDeploymentSource.ApplicationDirectory
		return e.ensureRollback(ctx, input.Request.RunningDeploymentSource.CommitHash), nil
	default:
		return sdk.StageStatusFailure, status.Error(codes.InvalidArgument, "unsupported stage")
	}
}

func (e *deployExecutor) ensureSync(ctx context.Context) sdk.StageStatus {
	cmd, ok := e.initLambrollCommand(ctx)
	if !ok {
		return sdk.StageStatusFailure
	}

	if err := cmd.Deploy(ctx, e.slp); err != nil {
		e.slp.Errorf("Failed to apply changes (%v)", err)
		return sdk.StageStatusFailure
	}

	e.slp.Success("Successfully applied changes")
	return sdk.StageStatusSuccess
}

func (e *deployExecutor) ensureDiff(ctx context.Context) sdk.StageStatus {
	cmd, ok := e.initLambrollCommand(ctx)
	if !ok {
		return sdk.StageStatusFailure
	}

	if err := cmd.Diff(ctx, e.slp); err != nil {
		e.slp.Errorf("Failed to apply changes (%v)", err)
		return sdk.StageStatusFailure
	}

	e.slp.Success("Successfully executed 'diff'")
	return sdk.StageStatusSuccess
}

func (e *deployExecutor) ensureRollback(ctx context.Context, runningCommitHash string) sdk.StageStatus {
	// There is nothing to do if this is the first deployment.
	if runningCommitHash == "" {
		e.slp.Errorf("Unable to determine the last deployed commit to rollback. It seems this is the first deployment.")
		return sdk.StageStatusFailure
	}

	e.slp.Infof("Start rolling back to the state defined at commit %s", runningCommitHash)

	cmd, ok := e.initLambrollCommand(ctx)
	if !ok {
		return sdk.StageStatusFailure
	}

	if err := cmd.Deploy(ctx, e.slp); err != nil {
		e.slp.Errorf("Failed to apply changes (%v)", err)
		return sdk.StageStatusFailure
	}

	e.slp.Success("Successfully rolled back the changes")
	return sdk.StageStatusSuccess
}

func showUsingVersion(ctx context.Context, cmd *cli.Lambroll, slp sdk.StageLogPersister) (ok bool) {
	version, err := cmd.Version(ctx)
	if err != nil {
		slp.Errorf("Failed to check lambroll version (%v)", err)
		return false
	}
	slp.Infof("Using lambroll version %q to execute the lambroll commands", version)
	return true
}
