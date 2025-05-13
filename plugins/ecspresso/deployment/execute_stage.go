package deployment

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
	"github.com/t-kikuc/pipecd-plugin-prototypes/ecspresso/cli"
	ecspressoconfig "github.com/t-kikuc/pipecd-plugin-prototypes/ecspresso/config"
	"github.com/t-kikuc/pipecd-plugin-prototypes/ecspresso/toolregistry"
)

type deployExecutor struct {
	appDir        string
	ecspressoPath string
	input         ecspressoconfig.EcspressoDeploymentInput
	slp           sdk.StageLogPersister
}

func (e *deployExecutor) initEcspressoCommand(ctx context.Context) (cmd *cli.Ecspresso, ok bool) {
	cmd = cli.NewEcspresso(
		e.ecspressoPath,
		e.appDir,
		e.input.Config,
	)

	if ok := showUsingVersion(ctx, cmd, e.slp); !ok {
		return nil, false
	}

	return cmd, true
}

func (p *Plugin) executeStage(ctx context.Context, input *sdk.ExecuteStageInput[ecspressoconfig.EcspressoApplicationSpec], dts []*sdk.DeployTarget[ecspressoconfig.EcspressoDeployTargetConfig]) (sdk.StageStatus, error) {
	toolRegistry := toolregistry.NewRegistry(input.Client.ToolRegistry())
	req := input.Request

	e := &deployExecutor{
		// input:  req.
		slp:    input.Client.LogPersister(),
		appDir: string(req.TargetDeploymentSource.ApplicationDirectory),
	}
	var err error
	e.ecspressoPath, err = toolRegistry.Ecspresso(ctx, dts[0].Config.Version)
	if err != nil {
		return sdk.StageStatusFailure, err
	}

	switch req.StageName {
	case stageEcspressoDeploy:
		return e.ensureSync(ctx), nil
	case stageEcspressoDiff:
		return e.ensureDiff(ctx), nil
	case stageEcspressoRollback:
		return e.ensureRollback(ctx, req.RunningDeploymentSource.CommitHash), nil
	default:
		return sdk.StageStatusFailure, status.Error(codes.InvalidArgument, "unsupported stage")
	}
}

func (e *deployExecutor) ensureSync(ctx context.Context) sdk.StageStatus {
	cmd, ok := e.initEcspressoCommand(ctx)
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
	cmd, ok := e.initEcspressoCommand(ctx)
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
	// Nothing to do if this is the first deployment.
	if runningCommitHash == "" {
		e.slp.Errorf("Unable to determine the last deployed commit to rollback. It seems this is the first deployment.")
		return sdk.StageStatusFailure
	}

	e.slp.Infof("Start rolling back to the state defined at commit %s", runningCommitHash)

	cmd, ok := e.initEcspressoCommand(ctx)
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

func showUsingVersion(ctx context.Context, cmd *cli.Ecspresso, slp sdk.StageLogPersister) (ok bool) {
	version, err := cmd.Version(ctx)
	if err != nil {
		slp.Errorf("Failed to check ecspresso version (%v)", err)
		return false
	}
	slp.Infof("Using ecspresso version %q to execute the ecspresso commands", version)
	return true
}
