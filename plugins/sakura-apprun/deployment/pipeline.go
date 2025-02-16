package deployment

import (
	"fmt"
	"slices"
	"time"

	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
)

type stage string

const (
	stageDeploy stage = "APPRUN_DEPLOY"

	stageRollback stage = "APPRUN_ROLLBACK"
)

var allStages = []string{
	string(stageDeploy),
	string(stageRollback),
}

var (
	predefinedStageDeploy = model.PipelineStage{
		Id:       "AppRunDeploy",
		Name:     string(stageDeploy),
		Desc:     "Create or Update an AppRun application",
		Rollback: false,
	}
	predefinedStageRollback = model.PipelineStage{
		Id:       "AppRunRollback",
		Name:     string(stageRollback),
		Desc:     "Rollback the deployment",
		Rollback: true,
	}
)

func (s stage) String() string {
	return string(s)
}

func buildQuickSyncStages(autoRollback bool, now time.Time) []*model.PipelineStage {
	out := make([]*model.PipelineStage, 0, 2)

	out = append(out, &model.PipelineStage{
		Id:        predefinedStageDeploy.GetId(),
		Name:      predefinedStageDeploy.GetName(),
		Desc:      predefinedStageDeploy.GetDesc(),
		Rollback:  predefinedStageDeploy.GetRollback(),
		Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
		Metadata:  nil,
		CreatedAt: now.Unix(),
		UpdatedAt: now.Unix(),

		Visible: true, // TODO: for debug.
	})

	// Append ROLLBACK stage if auto rollback is enabled.
	if autoRollback {
		out = append(out, &model.PipelineStage{
			Id:        predefinedStageRollback.GetId(),
			Name:      predefinedStageRollback.GetName(),
			Desc:      predefinedStageRollback.GetDesc(),
			Rollback:  predefinedStageRollback.GetRollback(),
			Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
			CreatedAt: now.Unix(),
			UpdatedAt: now.Unix(),
		})
	}

	return out
}

func buildPipelineStages(stages []*deployment.BuildPipelineSyncStagesRequest_StageConfig, autoRollback bool, now time.Time) []*model.PipelineStage {
	out := make([]*model.PipelineStage, 0, len(stages)+1)

	for _, s := range stages {
		id := s.GetId()
		if id == "" {
			id = fmt.Sprintf("stage-%d", s.GetIndex())
		}
		stage := &model.PipelineStage{
			Id:        id,
			Name:      s.GetName(),
			Desc:      s.GetDesc(),
			Index:     s.GetIndex(),
			Rollback:  false,
			Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
			CreatedAt: now.Unix(),
			UpdatedAt: now.Unix(),

			Visible: true, // TODO: for debug.
		}
		out = append(out, stage)
	}

	if autoRollback {
		minIndex := slices.MinFunc(stages, func(a, b *deployment.BuildPipelineSyncStagesRequest_StageConfig) int {
			return int(a.GetIndex() - b.GetIndex())
		}).GetIndex()

		out = append(out, &model.PipelineStage{
			Id:        predefinedStageRollback.GetId(),
			Name:      predefinedStageRollback.GetName(),
			Desc:      predefinedStageRollback.GetDesc(),
			Index:     minIndex,
			Rollback:  predefinedStageRollback.GetRollback(),
			Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
			CreatedAt: now.Unix(),
			UpdatedAt: now.Unix(),
		})
	}

	return out
}
