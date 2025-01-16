package deployment

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister"
	"github.com/t-kikuc/pipecd-plugin-prototypes/ecschedule/cli"
	ecspconfig "github.com/t-kikuc/pipecd-plugin-prototypes/ecschedule/config"
)

type deployExecutor struct {
	appDir         string
	ecschedulePath string
	input          ecspconfig.EcscheduleDeploymentInput
	slp            logpersister.StageLogPersister
}

func (e *deployExecutor) initEcscheduleCommand(ctx context.Context) (cmd *cli.Ecschedule, ok bool) {
	cmd = cli.NewEcschedule(
		e.ecschedulePath,
		e.appDir,
		e.input.Config,
	)

	if ok := showUsingVersion(ctx, cmd, e.slp); !ok {
		return nil, false
	}

	return cmd, true
}

func (s *DeploymentServiceServer) executeStage(ctx context.Context, slp logpersister.StageLogPersister, input *deployment.ExecutePluginInput) (model.StageStatus, error) {
	cfg, err := config.DecodeYAML[*ecspconfig.EcscheduleApplicationSpec](input.GetTargetDeploymentSource().GetApplicationConfig())
	if err != nil {
		slp.Errorf("Failed while decoding application config (%v)", err)
		return model.StageStatus_STAGE_FAILURE, err
	}

	e := &deployExecutor{
		input:  cfg.Spec.Input,
		slp:    slp,
		appDir: string(input.GetTargetDeploymentSource().GetApplicationDirectory()),
	}
	e.ecschedulePath, err = s.toolRegistry.Ecschedule(ctx, s.deployTargetConfig.Version)
	if err != nil {
		return model.StageStatus_STAGE_FAILURE, err
	}

	slp.Infof("[DEBUG] ### pipedv1 executeStage() ###")

	switch input.GetStage().GetName() {
	case stageApply.String():
		// return e.ensureSync(ctx), nil
		return e.ensureApply(ctx), nil
	case stageDiff.String():
		return e.ensureDiff(ctx), nil
	case stageRollback.String():
		e.appDir = string(input.GetRunningDeploymentSource().GetApplicationDirectory())
		return e.ensureRollback(ctx, input.GetDeployment().GetRunningCommitHash()), nil
	default:
		return model.StageStatus_STAGE_FAILURE, status.Error(codes.InvalidArgument, "unsupported stage")
	}
}

func (e *deployExecutor) ensureApply(ctx context.Context) model.StageStatus {
	cmd, ok := e.initEcscheduleCommand(ctx)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	if err := cmd.Apply(ctx, e.slp); err != nil {
		e.slp.Errorf("Failed to apply changes (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	e.slp.Success("Successfully applied changes")
	return model.StageStatus_STAGE_SUCCESS
}

func (e *deployExecutor) ensureDiff(ctx context.Context) model.StageStatus {
	cmd, ok := e.initEcscheduleCommand(ctx)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	if err := cmd.Diff(ctx, e.slp); err != nil {
		e.slp.Errorf("Failed to apply changes (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	e.slp.Success("Successfully executed 'diff'")
	return model.StageStatus_STAGE_SUCCESS
}

func (e *deployExecutor) ensureRollback(ctx context.Context, runningCommitHash string) model.StageStatus {
	// There is nothing to do if this is the first deployment.
	if runningCommitHash == "" {
		e.slp.Errorf("Unable to determine the last deployed commit to rollback. It seems this is the first deployment.")
		return model.StageStatus_STAGE_FAILURE
	}

	e.slp.Infof("Start rolling back to the state defined at commit %s", runningCommitHash)

	cmd, ok := e.initEcscheduleCommand(ctx)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	if err := cmd.Apply(ctx, e.slp); err != nil {
		e.slp.Errorf("Failed to apply changes (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	e.slp.Success("Successfully rolled back the changes")
	return model.StageStatus_STAGE_SUCCESS
}

func showUsingVersion(ctx context.Context, cmd *cli.Ecschedule, slp logpersister.StageLogPersister) (ok bool) {
	version, err := cmd.Version(ctx)
	if err != nil {
		slp.Errorf("Failed to check ecschedule version (%v)", err)
		return false
	}
	slp.Infof("Using ecschedule version %q to execute ecschedule commands", version)
	return true
}
