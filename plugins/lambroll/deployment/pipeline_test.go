package deployment

import (
	"testing"

	"github.com/stretchr/testify/assert"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

func TestBuildQuickSyncStages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		autoRollback bool
		expected     []sdk.PipelineStage
	}{
		{
			name:         "without auto rollback",
			autoRollback: false,
			expected: []sdk.PipelineStage{
				{
					Name:     "LAMBROLL_DEPLOY",
					Index:    0,
					Rollback: false,
				},
			},
		},
		{
			name:         "with auto rollback",
			autoRollback: true,
			expected: []sdk.PipelineStage{
				{
					Name:     "LAMBROLL_DEPLOY",
					Index:    0,
					Rollback: false,
				},
				{
					Name:     "LAMBROLL_ROLLBACK",
					Rollback: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := buildQuickSyncStages(tt.autoRollback)
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
		expected     []sdk.PipelineStage
	}{
		{
			name: "without auto rollback",
			stages: []sdk.StageConfig{
				{
					Name:  "Stage 1",
					Index: 0,
				},
				{
					Name:  "Stage 2",
					Index: 1,
				},
			},
			autoRollback: false,
			expected: []sdk.PipelineStage{
				{
					Name:     "Stage 1",
					Index:    0,
					Rollback: false,
				},
				{
					Name:     "Stage 2",
					Index:    1,
					Rollback: false,
				},
			},
		},
		{
			name: "with auto rollback",
			stages: []sdk.StageConfig{
				{
					Name:  "Stage 1",
					Index: 0,
				},
				{
					Name:  "Stage 2",
					Index: 1,
				},
			},
			autoRollback: true,
			expected: []sdk.PipelineStage{
				{
					Name:     "Stage 1",
					Index:    0,
					Rollback: false,
				},
				{
					Name:     "Stage 2",
					Index:    1,
					Rollback: false,
				},
				{
					Name:     "LAMBROLL_ROLLBACK",
					Index:    0,
					Rollback: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := buildPipelineStages(tt.stages, tt.autoRollback)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
