package deployment

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

func TestBuildQuickSyncStages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		autoRollback bool
		expected     []*sdk.PipelineStage
	}{
		{
			name:         "without auto rollback",
			autoRollback: false,
			expected: []*sdk.PipelineStage{
				{
					Index:    0,
					Name:     "ECSPRESSO_DEPLOY",
					Rollback: false,
				},
			},
		},
		{
			name:         "with auto rollback",
			autoRollback: true,
			expected: []*sdk.PipelineStage{
				{
					Index:    0,
					Name:     "ECSPRESSO_DEPLOY",
					Rollback: false,
				},
				{
					Index:    1,
					Name:     "ECSPRESSO_ROLLBACK",
					Rollback: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := buildQuickSyncStages(&sdk.BuildQuickSyncStagesInput{
				Request: sdk.BuildQuickSyncStagesRequest{
					Rollback: tt.autoRollback,
				},
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestBuildPipelineStages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		stages       []sdk.StageConfig
		autoRollback bool
		expected     []*sdk.PipelineStage
	}{
		{
			name: "without auto rollback",
			stages: []sdk.StageConfig{
				{
					Index: 0,
					Name:  "Stage 1",
				},
				{
					Index: 1,
					Name:  "Stage 2",
				},
			},
			autoRollback: false,
			expected: []*sdk.PipelineStage{
				{
					Index:    0,
					Name:     "Stage 1",
					Rollback: false,
				},
				{
					Index:    1,
					Name:     "Stage 2",
					Rollback: false,
				},
			},
		},
		{
			name: "with auto rollback",
			stages: []sdk.StageConfig{
				{
					Index: 0,
					Name:  "Stage 1",
				},
				{
					Index: 1,
					Name:  "Stage 2",
				},
			},
			autoRollback: true,
			expected: []*sdk.PipelineStage{
				{
					Index:    0,
					Name:     "Stage 1",
					Rollback: false,
				},
				{
					Index:    1,
					Name:     "Stage 2",
					Rollback: false,
				},
				{
					Index:    2,
					Name:     "ECSPRESSO_ROLLBACK",
					Rollback: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := buildPipelineSyncStages(&sdk.BuildPipelineSyncStagesInput{
				Request: sdk.BuildPipelineSyncStagesRequest{
					Stages:   tt.stages,
					Rollback: tt.autoRollback,
				},
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
