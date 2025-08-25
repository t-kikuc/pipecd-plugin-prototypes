package deployment

import (
	"slices"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

const (
	// stageApply executes "apply"
	stageApply string = "ECSCHEDULE_APPLY"
	// stageDiff executes "diff"
	stageDiff string = "ECSCHEDULE_DIFF"

	// stageRollback rollbacks the deployment.
	stageRollback string = "ECSCHEDULE_ROLLBACK"
)

func buildQuickSyncStages(autoRollback bool) []sdk.QuickSyncStage {
	out := make([]sdk.QuickSyncStage, 0, 2)

	out = append(out, sdk.QuickSyncStage{
		Name:               stageApply,
		Description:        "Sync by executing 'ecschedule apply'",
		Rollback:           false,
		Metadata:           nil,
		AvailableOperation: sdk.ManualOperationNone,
	})

	// Append ROLLBACK stage if auto rollback is enabled.
	if autoRollback {
		out = append(out, sdk.QuickSyncStage{
			Name:               stageRollback,
			Description:        "Rollback the deployment",
			Rollback:           true,
			Metadata:           nil,
			AvailableOperation: sdk.ManualOperationNone,
		})
	}

	return out
}

func buildPipelineStages(stages []sdk.StageConfig, autoRollback bool) []sdk.PipelineStage {
	out := make([]sdk.PipelineStage, 0, len(stages)+1)

	for _, s := range stages {
		stage := sdk.PipelineStage{
			Index:              s.Index,
			Name:               s.Name,
			Rollback:           false,
			Metadata:           nil,
			AvailableOperation: sdk.ManualOperationNone,
		}
		out = append(out, stage)
	}

	if autoRollback {
		// Use the minimum index of all stages in order to ... // TODO: Add comment
		minIndex := slices.MinFunc(stages, func(a, b sdk.StageConfig) int {
			return a.Index - b.Index
		}).Index

		out = append(out, sdk.PipelineStage{
			Name:               stageRollback,
			Index:              int(minIndex),
			Rollback:           true,
			Metadata:           nil,
			AvailableOperation: sdk.ManualOperationNone,
		})
	}

	return out
}
