package deployment

import (
	"slices"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

const (
	stageDeploy string = "APPRUN_DEPLOY"

	stageRollback string = "APPRUN_ROLLBACK"
)

func buildQuickSyncStages(autoRollback bool) []sdk.QuickSyncStage {
	out := make([]sdk.QuickSyncStage, 0, 2)

	out = append(out, sdk.QuickSyncStage{
		Name:        stageDeploy,
		Description: "Create or Update an AppRun application",
		Rollback:    false,
		Metadata:    nil,
	})

	// Append ROLLBACK stage if auto rollback is enabled.
	if autoRollback {
		out = append(out, sdk.QuickSyncStage{
			Name:        stageRollback,
			Description: "Rollback the deployment",
			Rollback:    true,
			Metadata:    nil,
		})
	}

	return out
}

func buildPipelineStages(stages []sdk.StageConfig, autoRollback bool) []sdk.PipelineStage {
	out := make([]sdk.PipelineStage, 0, len(stages)+1)

	for _, s := range stages {
		stage := sdk.PipelineStage{
			Name:     s.Name,
			Index:    s.Index,
			Rollback: false,
		}
		out = append(out, stage)
	}

	if autoRollback {
		minIndex := slices.MinFunc(stages, func(a, b sdk.StageConfig) int {
			return int(a.Index - b.Index)
		}).Index

		out = append(out, sdk.PipelineStage{
			Name:     stageRollback,
			Index:    minIndex,
			Rollback: true,
		})
	}

	return out
}
