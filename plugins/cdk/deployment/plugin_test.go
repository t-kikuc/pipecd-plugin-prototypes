package deployment

import (
	"context"
	"testing"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/stretchr/testify/assert"
)

func Test_FetchDefinedStages(t *testing.T) {
	p := &Plugin{}
	want := []string{
		"CDK_DEPLOY",
		"CDK_DIFF",
		"CDK_ROLLBACK",
	}
	got := p.FetchDefinedStages()

	assert.Equal(t, want, got)
}

func Test_DetermineStrategy(t *testing.T) {
	p := &Plugin{}
	got, err := p.DetermineStrategy(context.Background(), nil, nil)

	assert.NoError(t, err)
	assert.Nil(t, got)
}

func TestPlugin_BuildPipelineSyncStages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *sdk.BuildPipelineSyncStagesInput
		want  *sdk.BuildPipelineSyncStagesResponse
	}{
		{
			name: "single stage without rollback",
			input: &sdk.BuildPipelineSyncStagesInput{
				Request: sdk.BuildPipelineSyncStagesRequest{
					Stages: []sdk.StageConfig{
						{
							Name:  "CDK_DEPLOY",
							Index: 1,
						},
					},
					Rollback: false,
				},
			},
			want: &sdk.BuildPipelineSyncStagesResponse{
				Stages: []sdk.PipelineStage{
					{
						Name:               "CDK_DEPLOY",
						Index:              1,
						Rollback:           false,
						Metadata:           map[string]string{},
						AvailableOperation: sdk.ManualOperationNone,
					},
				},
			},
		},
		{
			name: "multiple stages without rollback",
			input: &sdk.BuildPipelineSyncStagesInput{
				Request: sdk.BuildPipelineSyncStagesRequest{
					Stages: []sdk.StageConfig{
						{
							Name:  "CDK_DEPLOY",
							Index: 1,
						},
						{
							Name:  "CDK_DIFF",
							Index: 3,
						},
					},
					Rollback: false,
				},
			},
			want: &sdk.BuildPipelineSyncStagesResponse{
				Stages: []sdk.PipelineStage{
					{
						Name:               "CDK_DEPLOY",
						Index:              1,
						Rollback:           false,
						Metadata:           map[string]string{},
						AvailableOperation: sdk.ManualOperationNone,
					},
					{
						Name:               "CDK_DIFF",
						Index:              3,
						Rollback:           false,
						Metadata:           map[string]string{},
						AvailableOperation: sdk.ManualOperationNone,
					},
				},
			},
		},
		{
			name: "multiple stages with rollback",
			input: &sdk.BuildPipelineSyncStagesInput{
				Request: sdk.BuildPipelineSyncStagesRequest{
					Stages: []sdk.StageConfig{
						{
							Name:  "CDK_DEPLOY",
							Index: 2,
						},
						{
							Name:  "CDK_DIFF",
							Index: 3,
						},
					},
					Rollback: true,
				},
			},
			want: &sdk.BuildPipelineSyncStagesResponse{
				Stages: []sdk.PipelineStage{
					{
						Name:               "CDK_DEPLOY",
						Index:              2,
						Rollback:           false,
						Metadata:           map[string]string{},
						AvailableOperation: sdk.ManualOperationNone,
					},
					{
						Name:               "CDK_DIFF",
						Index:              3,
						Rollback:           false,
						Metadata:           map[string]string{},
						AvailableOperation: sdk.ManualOperationNone,
					},
					{
						Name:               "CDK_ROLLBACK",
						Index:              2,
						Rollback:           true,
						Metadata:           map[string]string{},
						AvailableOperation: sdk.ManualOperationNone,
					},
				},
			},
		},
	}

	p := &Plugin{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := p.BuildPipelineSyncStages(t.Context(), nil, tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPlugin_BuildQuickSyncStages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *sdk.BuildQuickSyncStagesInput
		want  *sdk.BuildQuickSyncStagesResponse
	}{
		{
			name: "no rollback",
			input: &sdk.BuildQuickSyncStagesInput{
				Request: sdk.BuildQuickSyncStagesRequest{
					Rollback: false,
				},
			},
			want: &sdk.BuildQuickSyncStagesResponse{
				Stages: []sdk.QuickSyncStage{
					{
						Name:               "CDK_DEPLOY",
						Description:        "Sync by 'cdk deploy'",
						Rollback:           false,
						Metadata:           map[string]string{},
						AvailableOperation: sdk.ManualOperationNone,
					},
				},
			},
		},
		{
			name: "with rollback",
			input: &sdk.BuildQuickSyncStagesInput{
				Request: sdk.BuildQuickSyncStagesRequest{
					Rollback: true,
				},
			},
			want: &sdk.BuildQuickSyncStagesResponse{
				Stages: []sdk.QuickSyncStage{
					{
						Name:               "CDK_DEPLOY",
						Description:        "Sync by 'cdk deploy'",
						Rollback:           false,
						Metadata:           map[string]string{},
						AvailableOperation: sdk.ManualOperationNone,
					},
					{
						Name:               "CDK_ROLLBACK",
						Description:        "Rollback by 'cdk deploy' for the previous CDK files",
						Rollback:           true,
						Metadata:           map[string]string{},
						AvailableOperation: sdk.ManualOperationNone,
					},
				},
			},
		},
	}

	p := &Plugin{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := p.BuildQuickSyncStages(t.Context(), nil, tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
