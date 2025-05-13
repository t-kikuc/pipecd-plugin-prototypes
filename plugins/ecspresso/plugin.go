package main

import (
	"context"

	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

const (
// stageEcspressoSync string = "ECSPRESSO_SYNC"
)

// plugin implements the sdk.DeploymentPlugin interface.
type plugin struct{}

// TODO: move to config/ package and define the actual type.
type TempAppSpec struct {
}

// DetermineVersions determines the versions of the resources that will be deployed.
// This implements sdk.DeploymentPlugin.
func (p *plugin) DetermineVersions(ctx context.Context, _ *sdk.ConfigNone, input *sdk.DetermineVersionsInput[TempAppSpec]) (*sdk.DetermineVersionsResponse, error) {
	// TODO implement
	return &sdk.DetermineVersionsResponse{
		Versions: []sdk.ArtifactVersion{
			{
				Kind:    sdk.ArtifactKindUnknown,
				Version: "0.0.1",
				Name:    "ecspresso",
				URL:     "TODO",
			},
		},
	}, nil
}

// DetermineStrategy determines the strategy to deploy the resources.
// This implements sdk.DeploymentPlugin.
func (p *plugin) DetermineStrategy(ctx context.Context, _ *sdk.ConfigNone, input *sdk.DetermineStrategyInput[TempAppSpec]) (*sdk.DetermineStrategyResponse, error) {
	// TODO implement
	return &sdk.DetermineStrategyResponse{
		Strategy: sdk.SyncStrategyQuickSync,
		Summary:  "TODO",
	}, nil
}

// BuildQuickSyncStages builds the stages that will be executed during the quick sync process.
// This implements sdk.DeploymentPlugin.
func (p *plugin) BuildQuickSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildQuickSyncStagesInput) (*sdk.BuildQuickSyncStagesResponse, error) {
	// TODO implement
	return &sdk.BuildQuickSyncStagesResponse{
		Stages: []sdk.QuickSyncStage{},
	}, nil
}

// BuildPipelineSyncStages builds the stages that will be executed by the plugin.
// This implements sdk.StagePlugin.
func (p *plugin) BuildPipelineSyncStages(ctx context.Context, _ sdk.ConfigNone, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	// TODO: implement
	stages := make([]sdk.PipelineStage, 0, len(input.Request.Stages))
	for _, rs := range input.Request.Stages {
		stage := sdk.PipelineStage{
			Index:              rs.Index,
			Name:               rs.Name,
			Rollback:           false,
			Metadata:           map[string]string{},
			AvailableOperation: sdk.ManualOperationNone,
		}
		stages = append(stages, stage)
	}

	return &sdk.BuildPipelineSyncStagesResponse{
		Stages: stages,
	}, nil
}

// ExecuteStage executes the given stage.
// This implements sdk.StagePlugin.
func (p *plugin) ExecuteStage(ctx context.Context, _ sdk.ConfigNone, _ sdk.DeployTargetsNone, input *sdk.ExecuteStageInput[struct{}]) (*sdk.ExecuteStageResponse, error) {
	// TODO implement
	return &sdk.ExecuteStageResponse{
		Status: sdk.StageStatusSuccess,
	}, nil
}

// FetchDefinedStages returns the list of stages that the plugin can execute.
// This implements sdk.StagePlugin.
func (p *plugin) FetchDefinedStages() []string {
	return []string{}
}
