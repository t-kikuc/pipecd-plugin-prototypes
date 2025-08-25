package deployment

import (
	"context"
	"fmt"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/t-kikuc/pipecd-plugin-prototypes/ecschedule/cli"
	"github.com/t-kikuc/pipecd-plugin-prototypes/ecschedule/config"
	"github.com/t-kikuc/pipecd-plugin-prototypes/ecschedule/toolregistry"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type deployExecutor struct {
	appDir         string
	ecschedulePath string
	input          config.EcscheduleDeploymentInput
	lp             sdk.StageLogPersister
}

func executeStage(ctx context.Context, dtCfgs []*sdk.DeployTarget[config.EcscheduleDeployTargetConfig], input *sdk.ExecuteStageInput[config.EcscheduleDeploymentInput]) (sdk.StageStatus, error) {
	if len(dtCfgs) != 1 {
		return sdk.StageStatusFailure, status.Error(codes.InvalidArgument, "the number of deploy target must be one for this plugin.")
	}
	e := &deployExecutor{
		input:  *input.Request.TargetDeploymentSource.ApplicationConfig.Spec,
		lp:     input.Client.LogPersister(),
		appDir: input.Request.TargetDeploymentSource.ApplicationDirectory,
	}
	toolRegistry := toolregistry.NewRegistry(input.Client.ToolRegistry())
	p, err := toolRegistry.Ecschedule(ctx, dtCfgs[0].Config.Version)
	if err != nil {
		return sdk.StageStatusFailure, err
	}
	e.ecschedulePath = p

	switch input.Request.StageName {
	case stageApply:
		// return e.ensureSync(ctx), nil
		return e.ensureApply(ctx), nil
	case stageDiff:
		return e.ensureDiff(ctx), nil
	case stageRollback:
		e.appDir = input.Request.RunningDeploymentSource.ApplicationDirectory
		return e.ensureRollback(ctx, input.Request.RunningDeploymentSource.CommitHash), nil
	default:
		return sdk.StageStatusFailure, status.Error(codes.InvalidArgument, "unsupported stage")
	}
}

func (e *deployExecutor) initEcscheduleCommand(ctx context.Context) (cmd *cli.Ecschedule, ok bool) {
	cmd = cli.NewEcschedule(
		e.ecschedulePath,
		e.appDir,
		e.input.Config,
	)

	if ok := showUsingVersion(ctx, cmd, e.lp); !ok {
		return nil, false
	}

	return cmd, true
}

func (e *deployExecutor) ensureApply(ctx context.Context) sdk.StageStatus {
	cmd, ok := e.initEcscheduleCommand(ctx)
	if !ok {
		return sdk.StageStatusFailure
	}

	if err := cmd.Apply(ctx, e.lp); err != nil {
		e.lp.Error(fmt.Sprintf("Failed to apply changes (%v)", err))
		return sdk.StageStatusFailure
	}

	e.lp.Success("Successfully applied changes")
	return sdk.StageStatusSuccess
}

func (e *deployExecutor) ensureDiff(ctx context.Context) sdk.StageStatus {
	cmd, ok := e.initEcscheduleCommand(ctx)
	if !ok {
		return sdk.StageStatusFailure
	}

	if err := cmd.Diff(ctx, e.lp); err != nil {
		e.lp.Error(fmt.Sprintf("Failed to apply changes (%v)", err))
		return sdk.StageStatusFailure
	}

	e.lp.Success("Successfully executed 'diff'")
	return sdk.StageStatusSuccess
}

func (e *deployExecutor) ensureRollback(ctx context.Context, runningCommitHash string) sdk.StageStatus {
	// There is nothing to do if this is the first deployment.
	if runningCommitHash == "" {
		e.lp.Error("Unable to determine the last deployed commit to rollback. It seems this is the first deployment.")
		return sdk.StageStatusFailure
	}

	e.lp.Info(fmt.Sprintf("Start rolling back to the state defined at commit %s", runningCommitHash))

	cmd, ok := e.initEcscheduleCommand(ctx)
	if !ok {
		return sdk.StageStatusFailure
	}

	if err := cmd.Apply(ctx, e.lp); err != nil {
		e.lp.Error(fmt.Sprintf("Failed to apply changes (%v)", err))
		return sdk.StageStatusFailure
	}

	e.lp.Success("Successfully rolled back the changes")
	return sdk.StageStatusSuccess
}

func showUsingVersion(ctx context.Context, cmd *cli.Ecschedule, lp sdk.StageLogPersister) (ok bool) {
	version, err := cmd.Version(ctx)
	if err != nil {
		lp.Error(fmt.Sprintf("Failed to check ecschedule version (%v)", err))
		return false
	}
	lp.Info(fmt.Sprintf("Using ecschedule version %q to execute ecschedule commands", version))
	return true
}
