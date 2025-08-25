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
		expected     []sdk.QuickSyncStage
	}{
		{
			name:         "without auto rollback",
			autoRollback: false,
			expected: []sdk.QuickSyncStage{
				{
					Name:        "APPRUN_DEPLOY",
					Description: "Create or Update an AppRun application",
					Rollback:    false,
					Metadata:    nil,
				},
			},
		},
		{
			name:         "with auto rollback",
			autoRollback: true,
			expected: []sdk.QuickSyncStage{
				{
					Name:        "APPRUN_DEPLOY",
					Description: "Create or Update an AppRun application",
					Rollback:    false,
					Metadata:    nil,
				},
				{
					Name:        "APPRUN_ROLLBACK",
					Description: "Rollback the deployment",
					Rollback:    true,
					Metadata:    nil,
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
					Name: "Stage 1",
				},
				{
					Name: "Stage 2",
				},
			},
			autoRollback: false,
			expected: []sdk.PipelineStage{
				{
					Name: "Stage 1",
				},
				{
					Name: "Stage 2",
				},
			},
		},
		{
			name: "with auto rollback",
			stages: []sdk.StageConfig{
				{
					Name: "Stage 1",
				},
				{
					Name: "Stage 2",
				},
			},
			autoRollback: true,
			expected: []sdk.PipelineStage{
				{
					Name: "Stage 1",
				},
				{
					Name: "Stage 2",
				},
				{
					Name:     "APPRUN_ROLLBACK",
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
