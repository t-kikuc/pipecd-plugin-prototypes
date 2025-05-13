package deployment

import (
	"context"

	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
	ecspressoconfig "github.com/t-kikuc/pipecd-plugin-prototypes/ecspresso/config"
)

// Plugin implements the sdk.DeploymentPlugin interface.
type Plugin struct{}

var _ sdk.DeploymentPlugin[sdk.ConfigNone, ecspressoconfig.EcspressoDeployTargetConfig, ecspressoconfig.EcspressoApplicationSpec] = (*Plugin)(nil)

type toolRegistry interface {
	Ecspresso(ctx context.Context, version string) (path string, err error)
}

// DetermineVersions determines the versions of the resources that will be deployed.
// This implements sdk.DeploymentPlugin.
func (p *Plugin) DetermineVersions(ctx context.Context, _ *sdk.ConfigNone, input *sdk.DetermineVersionsInput[ecspressoconfig.EcspressoApplicationSpec]) (*sdk.DetermineVersionsResponse, error) {
	return determineVersions(input)
}

// DetermineStrategy determines the strategy to deploy the resources.
// This implements sdk.DeploymentPlugin.
func (p *Plugin) DetermineStrategy(ctx context.Context, _ *sdk.ConfigNone, _ *sdk.DetermineStrategyInput[ecspressoconfig.EcspressoApplicationSpec]) (*sdk.DetermineStrategyResponse, error) {
	return determineStrategy()
}

// BuildQuickSyncStages builds the stages that will be executed during the quick sync process.
// This implements sdk.DeploymentPlugin.
func (p *Plugin) BuildQuickSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildQuickSyncStagesInput) (*sdk.BuildQuickSyncStagesResponse, error) {
	return buildQuickSyncStages(input)
}

// BuildPipelineSyncStages builds the stages that will be executed by the plugin.
// This implements sdk.StagePlugin.
func (p *Plugin) BuildPipelineSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	return buildPipelineSyncStages(input)
}

// FetchDefinedStages returns the list of stages that the plugin can execute.
// This implements sdk.StagePlugin.
func (p *Plugin) FetchDefinedStages() []string {
	return fetchDefinedStages()
}

// ExecuteStage executes the given stage.
// This implements sdk.StagePlugin.
func (p *Plugin) ExecuteStage(ctx context.Context, _ *sdk.ConfigNone, dts []*sdk.DeployTarget[ecspressoconfig.EcspressoDeployTargetConfig], input *sdk.ExecuteStageInput[ecspressoconfig.EcspressoApplicationSpec]) (*sdk.ExecuteStageResponse, error) {
	status, err := p.executeStage(ctx, input, dts)
	if err != nil {
		return nil, err
	}
	return &sdk.ExecuteStageResponse{
		Status: status,
	}, nil
}
