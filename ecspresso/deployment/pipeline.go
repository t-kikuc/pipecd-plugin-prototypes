// Copyright 2024 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	// stageEcspressoDeploy executes "deploy"
	stageEcspressoDeploy stage = "ECSPRESSO_DEPLOY"
	// stageEcspressoDiff executes "diff"
	stageEcspressoDiff stage = "ECSPRESSO_DIFF"

	// stageEcspressoRollback rollbacks the deployment.
	stageEcspressoRollback stage = "ECSPRESSO_ROLLBACK"
)

var allStages = []string{
	string(stageEcspressoDeploy),
	string(stageEcspressoDiff),
	string(stageEcspressoRollback),
}

var (
	predefinedStageEcspressoDeploy = model.PipelineStage{
		Id:       "EcspressoDeploy",
		Name:     string(stageEcspressoDeploy),
		Desc:     "Sync by executing 'ecspresso deploy'",
		Rollback: false,
	}
	predefinedStageEcspressoRollback = model.PipelineStage{
		Id:       "EcspressoRollback",
		Name:     string(stageEcspressoRollback),
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
		Id:        predefinedStageEcspressoDeploy.GetId(),
		Name:      predefinedStageEcspressoDeploy.GetName(),
		Desc:      predefinedStageEcspressoDeploy.GetDesc(),
		Rollback:  predefinedStageEcspressoDeploy.GetRollback(),
		Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
		Metadata:  nil,
		CreatedAt: now.Unix(),
		UpdatedAt: now.Unix(),
	})

	// Append ROLLBACK stage if auto rollback is enabled.
	if autoRollback {
		out = append(out, &model.PipelineStage{
			Id:        predefinedStageEcspressoRollback.GetId(),
			Name:      predefinedStageEcspressoRollback.GetName(),
			Desc:      predefinedStageEcspressoRollback.GetDesc(),
			Rollback:  predefinedStageEcspressoRollback.GetRollback(),
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
		}
		out = append(out, stage)
	}

	if autoRollback {
		// Use the minimum index of all stages in order to ... // TODO: Add comment
		minIndex := slices.MinFunc(stages, func(a, b *deployment.BuildPipelineSyncStagesRequest_StageConfig) int {
			return int(a.GetIndex() - b.GetIndex())
		}).GetIndex()

		out = append(out, &model.PipelineStage{
			Id:        predefinedStageEcspressoRollback.GetId(),
			Name:      predefinedStageEcspressoRollback.GetName(),
			Desc:      predefinedStageEcspressoRollback.GetDesc(),
			Index:     minIndex,
			Rollback:  predefinedStageEcspressoRollback.GetRollback(),
			Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
			CreatedAt: now.Unix(),
			UpdatedAt: now.Unix(),
		})
	}

	return out
}
