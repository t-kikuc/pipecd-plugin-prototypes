package deployment

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	apprun "github.com/sacloud/apprun-api-go"
	v1 "github.com/sacloud/apprun-api-go/apis/v1"

	"github.com/t-kikuc/pipecd-plugin-prototypes/sakura-apprun/config"
)

type deployExecutor struct {
	appDir string
	input  config.AppRunDeploymentInput
	slp    sdk.StageLogPersister
}

func executeStage(ctx context.Context, dtCfgs []*sdk.DeployTarget[config.AppRunDeployTargetConfig], input *sdk.ExecuteStageInput[config.AppRunDeploymentInput]) (sdk.StageStatus, error) {
	e := &deployExecutor{
		input:  *input.Request.TargetDeploymentSource.ApplicationConfig.Spec,
		slp:    input.Client.LogPersister(),
		appDir: string(input.Request.TargetDeploymentSource.ApplicationDirectory),
	}

	switch input.Request.StageName {
	case stageDeploy:
		return e.ensureSync(ctx), nil
	case stageRollback:
		e.appDir = input.Request.RunningDeploymentSource.ApplicationDirectory
		return e.ensureRollback(ctx, input.Request.RunningDeploymentSource.CommitHash), nil
	default:
		return sdk.StageStatusFailure, status.Error(codes.InvalidArgument, "unsupported stage")
	}
}

func (e *deployExecutor) ensureSync(ctx context.Context) sdk.StageStatus {
	cli := &apprun.Client{} // Uses env vars by deafult (SAKURACLOUD_ACCESS_TOKEN, SAKURACLOUD_ACCESS_TOKEN_SECRET)

	manifest, err := loadManifest(fmt.Sprintf("%s/%s", e.appDir, e.input.ConfigFile))
	if err != nil {
		e.slp.Errorf("Failed to load manifest (%v)", err)
		return sdk.StageStatusFailure
	}

	op := apprun.NewApplicationOp(cli)

	exists, id, err := existsApplication(ctx, op, *manifest.patchBody.Name)
	if err != nil {
		e.slp.Errorf("Failed to check if application exists (%v)", err)
		return sdk.StageStatusFailure
	}

	if exists {
		e.slp.Infof("Start updating the existing application. id: %s", id)
		if _, err := op.Update(ctx, id, manifest.patchBody); err != nil {
			e.slp.Errorf("Failed to update application (%v)", err)
			return sdk.StageStatusFailure
		}
		e.slp.Success("Successfully updated the existing application")
		return sdk.StageStatusSuccess
	} else {
		e.slp.Infof("Start creating a new application. name: %s", *manifest.patchBody.Name)
		resp, err := op.Create(ctx, manifest.toCreateBody())
		if err != nil {
			e.slp.Errorf("Failed to create application (%v)", err)
			return sdk.StageStatusFailure
		}
		e.slp.Successf("Successfully created a new application. id: %s", resp.Id)
		return sdk.StageStatusSuccess
	}
}

func (e *deployExecutor) ensureRollback(ctx context.Context, runningCommitHash string) sdk.StageStatus {
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
