package deployment

import (
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
	ecspressoconfig "github.com/t-kikuc/pipecd-plugin-prototypes/ecspresso/config"
)

const (
	// stageEcspressoDeploy executes "deploy"
	stageEcspressoDeploy = "ECSPRESSO_DEPLOY"
	// stageEcspressoDiff executes "diff"
	stageEcspressoDiff = "ECSPRESSO_DIFF"

	// stageEcspressoRollback rollbacks the deployment.
	stageEcspressoRollback = "ECSPRESSO_ROLLBACK"
)

func determineVersions(input *sdk.DetermineVersionsInput[ecspressoconfig.EcspressoApplicationSpec]) (*sdk.DetermineVersionsResponse, error) {
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

func determineStrategy() (*sdk.DetermineStrategyResponse, error) {
	return &sdk.DetermineStrategyResponse{
		Strategy: sdk.SyncStrategyPipelineSync,
		Summary:  "PipelineSync with the specified pipeline",
	}, nil
}

func buildQuickSyncStages(input *sdk.BuildQuickSyncStagesInput) (*sdk.BuildQuickSyncStagesResponse, error) {
	stages := make([]sdk.QuickSyncStage, 0, 2)

	stages = append(stages, sdk.QuickSyncStage{
		Name:               stageEcspressoDeploy,
		Description:        "Sync by executing 'ecspresso deploy'",
		Rollback:           false,
		Metadata:           make(map[string]string, 0),
		AvailableOperation: sdk.ManualOperationNone,
	},
	)

	if input.Request.Rollback {
		stages = append(stages, sdk.QuickSyncStage{
			Name:               stageEcspressoRollback,
			Description:        "Rollback the deployment",
			Rollback:           true,
			Metadata:           make(map[string]string, 0),
			AvailableOperation: sdk.ManualOperationNone,
		})
	}

	return &sdk.BuildQuickSyncStagesResponse{
		Stages: stages,
	}, nil
}

func buildPipelineSyncStages(input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
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

	if input.Request.Rollback {
		stages = append(stages, sdk.PipelineStage{
			Index:    len(input.Request.Stages),
			Name:     stageEcspressoRollback,
			Rollback: true,
		})
	}

	return &sdk.BuildPipelineSyncStagesResponse{
		Stages: stages,
	}, nil
}

func fetchDefinedStages() []string {
	return []string{
		stageEcspressoDeploy,
		stageEcspressoDiff,
		stageEcspressoRollback,
	}
}
