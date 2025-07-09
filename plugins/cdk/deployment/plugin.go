package deployment

import (
	"context"
	"fmt"
	"slices"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	config "github.com/t-kikuc/pipecd-plugin-prototypes/cdk/config"
)

const (
	// stageDeploy executes "cdk deploy"
	stageDeploy string = "CDK_DEPLOY"
	// stageDiff executes "cdk diff"
	stageDiff string = "CDK_DIFF"

	// stageRollback rollbacks by executing "cdk deploy" for the previous success commit.
	stageRollback string = "CDK_ROLLBACK"
)

// Plugin implements sdk.DeploymentPlugin for Terraform.
type Plugin struct{}

var _ sdk.DeploymentPlugin[sdk.ConfigNone, config.DeployTargetConfig, config.ApplicationSpec] = (*Plugin)(nil)

// FetchDefinedStages implements sdk.DeploymentPlugin.
func (p *Plugin) FetchDefinedStages() []string {
	return []string{
		stageDeploy,
		stageDiff,
		stageRollback,
	}
}

// BuildPipelineSyncStages implements sdk.DeploymentPlugin.
func (p *Plugin) BuildPipelineSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	reqStages := input.Request.Stages
	out := make([]sdk.PipelineStage, 0, len(reqStages)+1)

	for _, s := range reqStages {
		out = append(out, sdk.PipelineStage{
			Name:               s.Name,
			Index:              s.Index,
			Rollback:           false,
			Metadata:           make(map[string]string),
			AvailableOperation: sdk.ManualOperationNone,
		})
	}
	if input.Request.Rollback {
		minIndex := slices.MinFunc(reqStages, func(a, b sdk.StageConfig) int { return a.Index - b.Index }).Index
		out = append(out, sdk.PipelineStage{
			Name:               stageRollback,
			Index:              minIndex,
			Rollback:           true,
			Metadata:           make(map[string]string),
			AvailableOperation: sdk.ManualOperationNone,
		})
	}
	return &sdk.BuildPipelineSyncStagesResponse{
		Stages: out,
	}, nil
}

// BuildQuickSyncStages implements sdk.DeploymentPlugin.
func (p *Plugin) BuildQuickSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildQuickSyncStagesInput) (*sdk.BuildQuickSyncStagesResponse, error) {
	stages := make([]sdk.QuickSyncStage, 0, 2)
	stages = append(stages, sdk.QuickSyncStage{
		Name:               stageDeploy,
		Description:        "Sync by 'cdk deploy'",
		Rollback:           false,
		Metadata:           map[string]string{},
		AvailableOperation: sdk.ManualOperationNone,
	})

	if input.Request.Rollback {
		stages = append(stages, sdk.QuickSyncStage{
			Name:               stageRollback,
			Description:        "Rollback by 'cdk deploy' for the previous CDK files",
			Rollback:           true,
			Metadata:           map[string]string{},
			AvailableOperation: sdk.ManualOperationNone,
		})
	}
	return &sdk.BuildQuickSyncStagesResponse{
		Stages: stages,
	}, nil
}

// DetermineStrategy implements sdk.DeploymentPlugin.
// It returns (nil, nil) because this plugin does not have specific logic for DetermineStrategy.
func (p *Plugin) DetermineStrategy(ctx context.Context, _ *sdk.ConfigNone, input *sdk.DetermineStrategyInput[config.ApplicationSpec]) (*sdk.DetermineStrategyResponse, error) {
	return nil, nil
}

// DetermineVersions implements sdk.DeploymentPlugin.
func (p *Plugin) DetermineVersions(ctx context.Context, _ *sdk.ConfigNone, input *sdk.DetermineVersionsInput[config.ApplicationSpec]) (*sdk.DetermineVersionsResponse, error) {
	// TODO: implement
	return &sdk.DetermineVersionsResponse{}, nil
}

// ExecuteStage executes a stage.
func (p *Plugin) ExecuteStage(ctx context.Context, _ *sdk.ConfigNone, dts []*sdk.DeployTarget[config.DeployTargetConfig], input *sdk.ExecuteStageInput[config.ApplicationSpec]) (*sdk.ExecuteStageResponse, error) {
	switch input.Request.StageName {
	case stageDeploy:
		return &sdk.ExecuteStageResponse{
			Status: executeDeploy(ctx, dts[0], input),
		}, nil
	case stageDiff:
		panic("unimplemented")
		// return executeDiff(ctx, dts[0], input)
	case stageRollback:
		panic("unimplemented")
		// return executeRollback(ctx, dts[0], input)
	default:
		return nil, fmt.Errorf("unknown stage: %s", input.Request.StageName)
	}
}
