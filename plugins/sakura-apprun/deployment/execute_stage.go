package deployment

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister"
	apprun "github.com/sacloud/apprun-api-go"
	v1 "github.com/sacloud/apprun-api-go/apis/v1"

	apprunconfig "github.com/t-kikuc/pipecd-plugin-prototypes/sakura-apprun/config"
)

type deployExecutor struct {
	appDir string
	input  apprunconfig.AppRunDeploymentInput
	slp    logpersister.StageLogPersister
}

func (s *DeploymentServiceServer) executeStage(ctx context.Context, slp logpersister.StageLogPersister, input *deployment.ExecutePluginInput) (model.StageStatus, error) {
	cfg, err := config.DecodeYAML[*apprunconfig.AppRunApplicationSpec](input.GetTargetDeploymentSource().GetApplicationConfig())
	if err != nil {
		slp.Errorf("Failed while decoding application config (%v)", err)
		return model.StageStatus_STAGE_FAILURE, err
	}

	e := &deployExecutor{
		input:  cfg.Spec.Input,
		slp:    slp,
		appDir: string(input.GetTargetDeploymentSource().GetApplicationDirectory()),
	}

	slp.Infof("[DEBUG sakura-apprun] ### pipedv1 executeStage() > %s xxxxx ###", input.GetStage().GetName())

	switch input.GetStage().GetName() {
	case stageDeploy.String():
		return e.ensureSync(ctx), nil
	case stageRollback.String():
		e.appDir = string(input.GetRunningDeploymentSource().GetApplicationDirectory())
		return e.ensureRollback(ctx, input.GetDeployment().GetRunningCommitHash()), nil
	default:
		return model.StageStatus_STAGE_FAILURE, status.Error(codes.InvalidArgument, "unsupported stage")
	}
}

func (e *deployExecutor) ensureSync(ctx context.Context) model.StageStatus {
	cli := &apprun.Client{} // Uses env vars by deafult (SAKURACLOUD_ACCESS_TOKEN, SAKURACLOUD_ACCESS_TOKEN_SECRET)

	e.slp.Infof("[DEBUG] load manifest from %s/%s ###", e.appDir, e.input.ConfigFile)
	manifest, err := loadManifest(fmt.Sprintf("%s/%s", e.appDir, e.input.ConfigFile))
	if err != nil {
		e.slp.Errorf("Failed to load manifest (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	e.slp.Infof("[DEBUG] manifest.patchBody: %+v", manifest.patchBody)

	op := apprun.NewApplicationOp(cli)

	exists, id, err := existsApplication(ctx, op, *manifest.patchBody.Name)
	if err != nil {
		e.slp.Errorf("Failed to check if application exists (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	if exists {
		e.slp.Infof("Start updating the existing application. id: %s", id)
		if _, err := op.Update(ctx, id, manifest.patchBody); err != nil {
			e.slp.Errorf("Failed to update application (%v)", err)
			return model.StageStatus_STAGE_FAILURE
		}
		e.slp.Success("Successfully updated the existing application")
		return model.StageStatus_STAGE_SUCCESS
	} else {
		e.slp.Infof("Start creating a new application. name: %s", *manifest.patchBody.Name)
		resp, err := op.Create(ctx, manifest.toCreateBody())
		if err != nil {
			e.slp.Errorf("Failed to create application (%v)", err)
			return model.StageStatus_STAGE_FAILURE
		}
		e.slp.Successf("Successfully created a new application. id: %s", resp.Id)
		return model.StageStatus_STAGE_SUCCESS
	}
}

func (e *deployExecutor) ensureRollback(ctx context.Context, runningCommitHash string) model.StageStatus {
	panic("rollbakc is not implemented yet")
}

// TODO: Optimize performance
func existsApplication(ctx context.Context, op apprun.ApplicationAPI, name string) (exist bool, id string, err error) {
	listReq := &v1.ListApplicationsParams{
		PageSize:  p(100),
		SortField: p("name"),
		SortOrder: p(v1.ListApplicationsParamsSortOrder(v1.ListApplicationsParamsSortOrderAsc)),
	}
	page := 0
	for {
		page++
		listReq.PageNum = p(page)
		res, err := op.List(ctx, listReq)
		if err != nil {
			return false, "", err
		}
		if len(*res.Data) == 0 {
			return false, "", nil
		}

		for _, app := range *res.Data {
			if *app.Name == name {
				return true, *app.Id, nil
			}
		}
	}
}

func p[T any](v T) *T {
	return &v
}
