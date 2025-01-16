// Copyright 2024 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deployment

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister"
	"github.com/t-kikuc/pipecd-plugin-prototypes/ecspresso/cli"
	ecspconfig "github.com/t-kikuc/pipecd-plugin-prototypes/ecspresso/config"
)

type deployExecutor struct {
	appDir        string
	ecspressoPath string
	input         ecspconfig.EcspressoDeploymentInput
	slp           logpersister.StageLogPersister
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

func (s *DeploymentServiceServer) executeStage(ctx context.Context, slp logpersister.StageLogPersister, input *deployment.ExecutePluginInput) (model.StageStatus, error) {
	cfg, err := config.DecodeYAML[*ecspconfig.EcspressoApplicationSpec](input.GetTargetDeploymentSource().GetApplicationConfig())
	if err != nil {
		slp.Errorf("Failed while decoding application config (%v)", err)
		return model.StageStatus_STAGE_FAILURE, err
	}

	e := &deployExecutor{
		input:  cfg.Spec.Input,
		slp:    slp,
		appDir: string(input.GetTargetDeploymentSource().GetApplicationDirectory()),
	}
	e.ecspressoPath, err = s.toolRegistry.Ecspresso(ctx, s.deployTargetConfig.Version)
	if err != nil {
		return model.StageStatus_STAGE_FAILURE, err
	}

	slp.Infof("[DEBUG] ### pipedv1 executeStage() ###")

	switch input.GetStage().GetName() {
	case stageEcspressoDeploy.String():
		return e.ensureSync(ctx), nil
	case stageEcspressoDiff.String():
		return e.ensureDiff(ctx), nil
	case stageEcspressoRollback.String():
		e.appDir = string(input.GetRunningDeploymentSource().GetApplicationDirectory())
		return e.ensureRollback(ctx, input.GetDeployment().GetRunningCommitHash()), nil
	default:
		return model.StageStatus_STAGE_FAILURE, status.Error(codes.InvalidArgument, "unsupported stage")
	}
}

func (e *deployExecutor) ensureSync(ctx context.Context) model.StageStatus {
	cmd, ok := e.initEcspressoCommand(ctx)
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
	cmd, ok := e.initEcspressoCommand(ctx)
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

	cmd, ok := e.initEcspressoCommand(ctx)
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

func showUsingVersion(ctx context.Context, cmd *cli.Ecspresso, slp logpersister.StageLogPersister) (ok bool) {
	version, err := cmd.Version(ctx)
	if err != nil {
		slp.Errorf("Failed to check ecspresso version (%v)", err)
		return false
	}
	slp.Infof("Using ecspresso version %q to execute the ecspresso commands", version)
	return true
}
