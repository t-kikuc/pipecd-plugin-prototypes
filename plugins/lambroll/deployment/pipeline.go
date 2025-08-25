package deployment

import (
	"slices"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

const (
	// stageDeploy executes "deploy"
	stageDeploy string = "LAMBROLL_DEPLOY"
	// stageDiff executes "diff"
	stageDiff string = "LAMBROLL_DIFF"

	// stageRollback rollbacks the deployment.
	stageRollback string = "LAMBROLL_ROLLBACK"
)

func buildQuickSyncStages(autoRollback bool) []sdk.QuickSyncStage {
	out := make([]sdk.QuickSyncStage, 0, 2)

	out = append(out, sdk.QuickSyncStage{
		Name:     stageDeploy,
		Rollback: false,
	})

	// Append ROLLBACK stage if auto rollback is enabled.
	if autoRollback {
		out = append(out, sdk.QuickSyncStage{
			Name:     stageRollback,
			Rollback: true,
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
		// Use the minimum index of all stages in order to ... // TODO: Add comment
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
