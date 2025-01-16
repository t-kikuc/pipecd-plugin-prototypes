package deployment

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
)

func TestBuildQuickSyncStages(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name         string
		autoRollback bool
		expected     []*model.PipelineStage
	}{
		{
			name:         "without auto rollback",
			autoRollback: false,
			expected: []*model.PipelineStage{
				{
					Id:        "EcspressoDeploy",
					Name:      "ECSPRESSO_DEPLOY",
					Desc:      "Sync by executing 'ecspresso deploy'",
					Index:     0,
					Rollback:  false,
					Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
					Metadata:  nil,
					CreatedAt: now.Unix(),
					UpdatedAt: now.Unix(),

					Visible: true, // TODO: This is for debug with v0 UI
				},
			},
		},
		{
			name:         "with auto rollback",
			autoRollback: true,
			expected: []*model.PipelineStage{
				{
					Id:        "EcspressoDeploy",
					Name:      "ECSPRESSO_DEPLOY",
					Desc:      "Sync by executing 'ecspresso deploy'",
					Index:     0,
					Rollback:  false,
					Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
					Metadata:  nil,
					CreatedAt: now.Unix(),
					UpdatedAt: now.Unix(),

					Visible: true, // TODO: This is for debug with v0 UI
				},
				{
					Id:        "EcspressoRollback",
					Name:      "ECSPRESSO_ROLLBACK",
					Desc:      "Rollback the deployment",
					Rollback:  true,
					Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
					CreatedAt: now.Unix(),
					UpdatedAt: now.Unix(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := buildQuickSyncStages(tt.autoRollback, now)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestBuildPipelineStages(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name         string
		stages       []*deployment.BuildPipelineSyncStagesRequest_StageConfig
		autoRollback bool
		expected     []*model.PipelineStage
	}{
		{
			name: "without auto rollback",
			stages: []*deployment.BuildPipelineSyncStagesRequest_StageConfig{
				{
					Id:    "stage-1",
					Name:  "Stage 1",
					Desc:  "Description 1",
					Index: 0,
				},
				{
					Id:    "stage-2",
					Name:  "Stage 2",
					Desc:  "Description 2",
					Index: 1,
				},
			},
			autoRollback: false,
			expected: []*model.PipelineStage{
				{
					Id:        "stage-1",
					Name:      "Stage 1",
					Desc:      "Description 1",
					Index:     0,
					Rollback:  false,
					Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
					CreatedAt: now.Unix(),
					UpdatedAt: now.Unix(),

					Visible: true, // TODO: This is for debug with v0 UI
				},
				{
					Id:        "stage-2",
					Name:      "Stage 2",
					Desc:      "Description 2",
					Index:     1,
					Rollback:  false,
					Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
					CreatedAt: now.Unix(),
					UpdatedAt: now.Unix(),

					Visible: true, // TODO: This is for debug with v0 UI
				},
			},
		},
		{
			name: "with auto rollback",
			stages: []*deployment.BuildPipelineSyncStagesRequest_StageConfig{
				{
					Id:    "stage-1",
					Name:  "Stage 1",
					Desc:  "Description 1",
					Index: 0,
				},
				{
					Id:    "stage-2",
					Name:  "Stage 2",
					Desc:  "Description 2",
					Index: 1,
				},
			},
			autoRollback: true,
			expected: []*model.PipelineStage{
				{
					Id:        "stage-1",
					Name:      "Stage 1",
					Desc:      "Description 1",
					Index:     0,
					Rollback:  false,
					Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
					CreatedAt: now.Unix(),
					UpdatedAt: now.Unix(),

					Visible: true, // TODO: This is for debug with v0 UI
				},
				{
					Id:        "stage-2",
					Name:      "Stage 2",
					Desc:      "Description 2",
					Index:     1,
					Rollback:  false,
					Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
					CreatedAt: now.Unix(),
					UpdatedAt: now.Unix(),

					Visible: true, // TODO: This is for debug with v0 UI
				},
				{
					Id:        "EcspressoRollback",
					Name:      "ECSPRESSO_ROLLBACK",
					Desc:      "Rollback the deployment",
					Index:     0,
					Rollback:  true,
					Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
					CreatedAt: now.Unix(),
					UpdatedAt: now.Unix(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := buildPipelineStages(tt.stages, tt.autoRollback, now)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
