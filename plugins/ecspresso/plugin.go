package main

import (
	"context"

	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
	ecspressoconfig "github.com/t-kikuc/pipecd-plugin-prototypes/ecspresso/config"
	"github.com/t-kikuc/pipecd-plugin-prototypes/ecspresso/deployment"
)

const (
// stageEcspressoSync string = "ECSPRESSO_SYNC"
)

// plugin implements the sdk.DeploymentPlugin interface.
type plugin struct{}

// DetermineVersions determines the versions of the resources that will be deployed.
// This implements sdk.DeploymentPlugin.
func (p *plugin) DetermineVersions(ctx context.Context, _ *sdk.ConfigNone, input *sdk.DetermineVersionsInput[ecspressoconfig.EcspressoApplicationSpec]) (*sdk.DetermineVersionsResponse, error) {
	return deployment.DetermineVersions(input)
}

// DetermineStrategy determines the strategy to deploy the resources.
// This implements sdk.DeploymentPlugin.
func (p *plugin) DetermineStrategy(ctx context.Context, _ *sdk.ConfigNone, _ *sdk.DetermineStrategyInput[ecspressoconfig.EcspressoApplicationSpec]) (*sdk.DetermineStrategyResponse, error) {
	return deployment.DetermineStrategy()
}

// BuildQuickSyncStages builds the stages that will be executed during the quick sync process.
// This implements sdk.DeploymentPlugin.
func (p *plugin) BuildQuickSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildQuickSyncStagesInput) (*sdk.BuildQuickSyncStagesResponse, error) {
	return deployment.BuildQuickSyncStages(input)
}

// BuildPipelineSyncStages builds the stages that will be executed by the plugin.
// This implements sdk.StagePlugin.
func (p *plugin) BuildPipelineSyncStages(ctx context.Context, _ sdk.ConfigNone, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	return deployment.BuildPipelineSyncStages(input)
}

// FetchDefinedStages returns the list of stages that the plugin can execute.
// This implements sdk.StagePlugin.
func (p *plugin) FetchDefinedStages() []string {
	return deployment.FetchDefinedStages()
}

// ExecuteStage executes the given stage.
// This implements sdk.StagePlugin.
func (p *plugin) ExecuteStage(ctx context.Context, _ sdk.ConfigNone, _ sdk.DeployTargetsNone, input *sdk.ExecuteStageInput[struct{}]) (*sdk.ExecuteStageResponse, error) {
	// TODO implement
	return &sdk.ExecuteStageResponse{
		Status: sdk.StageStatusSuccess,
	}, nil
}
