package deployment

import (
	"context"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	config "github.com/t-kikuc/pipecd-plugin-prototypes/sakura-apprun/config"
)

type Plugin struct{}

var _ sdk.DeploymentPlugin[sdk.ConfigNone, config.AppRunDeployTargetConfig, config.AppRunDeploymentInput] = &Plugin{}

// DetermineStrategy implements sdk.DeploymentPlugin.
func (p *Plugin) DetermineStrategy(ctx context.Context, _ *sdk.ConfigNone, input *sdk.DetermineStrategyInput[config.AppRunDeploymentInput]) (*sdk.DetermineStrategyResponse, error) {
	return nil, nil
}

// DetermineVersions implements sdk.DeploymentPlugin.
func (p *Plugin) DetermineVersions(ctx context.Context, _ *sdk.ConfigNone, input *sdk.DetermineVersionsInput[config.AppRunDeploymentInput]) (*sdk.DetermineVersionsResponse, error) {
	return &sdk.DetermineVersionsResponse{
		// TODO: implement
		Versions: nil,
	}, nil
}

// BuildPipelineSyncStages implements sdk.DeploymentPlugin.
func (p *Plugin) BuildPipelineSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	stages := buildPipelineStages(input.Request.Stages, input.Request.Rollback)
	return &sdk.BuildPipelineSyncStagesResponse{
		Stages: stages,
	}, nil
}

// BuildQuickSyncStages implements sdk.DeploymentServiceServer.
func (p *Plugin) BuildQuickSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildQuickSyncStagesInput) (*sdk.BuildQuickSyncStagesResponse, error) {
	stages := buildQuickSyncStages(input.Request.Rollback)
	return &sdk.BuildQuickSyncStagesResponse{
		Stages: stages,
	}, nil
}

// FetchDefinedStages implements sdk.DeploymentServiceServer.
func (p *Plugin) FetchDefinedStages() []string {
	return []string{
		stageDeploy,
		stageRollback,
	}
}

// ExecuteStage performs stage-defined tasks.
// It returns stage status after execution without error.
// An error will be returned only if the given stage is not supported.
func (p *Plugin) ExecuteStage(ctx context.Context, _ *sdk.ConfigNone, dtCfgs []*sdk.DeployTarget[config.AppRunDeployTargetConfig], input *sdk.ExecuteStageInput[config.AppRunDeploymentInput]) (response *sdk.ExecuteStageResponse, _ error) {
	status, err := executeStage(ctx, dtCfgs, input)
	if err != nil {
		return nil, err
	}
	return &sdk.ExecuteStageResponse{
		Status: status,
	}, nil
}
