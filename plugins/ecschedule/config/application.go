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

package config

import (
	config "github.com/pipe-cd/pipecd/pkg/configv1"
)

// EcscheduleApplicationSpec represents an application configuration for Ecschedule application.
type EcscheduleApplicationSpec struct {
	config.GenericApplicationSpec
	// Input for ecschedule deployment. e.g. ecschedule version
	Input EcscheduleDeploymentInput `json:"input"`
	// Configuration for quick sync.
	QuickSync EcscheduleApplyStageOptions `json:"quickSync"`
}

func (s *EcscheduleApplicationSpec) Validate() error {
	// TODO: Validate EcscheduleApplicationSpec fields.
	return nil
}

type EcscheduleDeploymentInput struct {
	// Config is the path to the ecschedule config file. This will be used as `ecschedule deploy --config <Config>`
	Config string `json:"config"`
}

// EcscheduleApplyStageOptions contains all configurable values for a ECSPRESSO_SYNC stage.
type EcscheduleApplyStageOptions struct {
}

// EcscheduleDiffStageOptions contains all configurable values for a ECSPRESSO_PLAN stage.
type EcscheduleDiffStageOptions struct {
	// Exit the pipeline if the result is "No Changes" with success status.
	ExitOnNoChanges bool `json:"exitOnNoChanges"`
}
