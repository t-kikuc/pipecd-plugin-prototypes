package deployment

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister"
	"github.com/t-kikuc/pipecd-plugin-prototypes/cdk/cli"
	ecspconfig "github.com/t-kikuc/pipecd-plugin-prototypes/cdk/config"
)

type deployExecutor struct {
	appDir  string
	cdkPath string
	input   ecspconfig.CDKDeploymentInput
	slp     logpersister.StageLogPersister
}

func (e *deployExecutor) initCDKCommand(ctx context.Context) (cmd *cli.CDK, ok bool) {
	cmd = cli.NewCDK(
		e.cdkPath,
		e.appDir,
		e.input.FunctionFile,
		e.input.Source,
	)

	if ok := showUsingVersion(ctx, cmd, e.slp); !ok {
		return nil, false
	}

	return cmd, true
}

func (s *DeploymentServiceServer) executeStage(ctx context.Context, slp logpersister.StageLogPersister, input *deployment.ExecutePluginInput) (model.StageStatus, error) {
	cfg, err := config.DecodeYAML[*ecspconfig.CDKApplicationSpec](input.GetTargetDeploymentSource().GetApplicationConfig())
	if err != nil {
		slp.Errorf("Failed while decoding application config (%v)", err)
		return model.StageStatus_STAGE_FAILURE, err
	}

	e := &deployExecutor{
		input:  cfg.Spec.Input,
		slp:    slp,
		appDir: string(input.GetTargetDeploymentSource().GetApplicationDirectory()),
	}
	e.cdkPath, err = s.toolRegistry.CDK(ctx, s.deployTargetConfig.NodeVersion, s.deployTargetConfig.Version)
	if err != nil {
		return model.StageStatus_STAGE_FAILURE, err
	}

	slp.Infof("[DEBUG cdk] ### pipedv1 executeStage() > %s ###", input.GetStage().GetName())

	switch input.GetStage().GetName() {
	case stageDeploy.String():
		return e.ensureSync(ctx), nil
	case stageDiff.String():
		return e.ensureDiff(ctx), nil
	case stageRollback.String():
		e.appDir = string(input.GetRunningDeploymentSource().GetApplicationDirectory())
		return e.ensureRollback(ctx, input.GetDeployment().GetRunningCommitHash()), nil
	default:
		return model.StageStatus_STAGE_FAILURE, status.Error(codes.InvalidArgument, "unsupported stage")
	}
}

func (e *deployExecutor) ensureSync(ctx context.Context) model.StageStatus {
	cmd, ok := e.initCDKCommand(ctx)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	if err := cmd.Deploy(ctx, e.slp); err != nil {
		e.slp.Errorf("Failed to apply changes (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	e.slp.Success("Successfully applied changes")
	return model.StageStatus_STAGE_SUCCESS
}

func (e *deployExecutor) ensureDiff(ctx context.Context) model.StageStatus {
	cmd, ok := e.initCDKCommand(ctx)
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

	cmd, ok := e.initCDKCommand(ctx)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	if err := cmd.Deploy(ctx, e.slp); err != nil {
		e.slp.Errorf("Failed to apply changes (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	e.slp.Success("Successfully rolled back the changes")
	return model.StageStatus_STAGE_SUCCESS
}

func showUsingVersion(ctx context.Context, cmd *cli.CDK, slp logpersister.StageLogPersister) (ok bool) {
	version, err := cmd.Version(ctx)
	if err != nil {
		slp.Errorf("Failed to check cdk version (%v)", err)
		return false
	}
	slp.Infof("Using cdk version %q to execute the cdk commands", version)
	return true
}
