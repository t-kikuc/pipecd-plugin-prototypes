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

// EcspressoApplicationSpec represents an application configuration for Ecspresso application.
type EcspressoApplicationSpec struct {
	config.GenericApplicationSpec
	// Input for ecspresso deployment. e.g. ecspresso version, workspace
	Input EcspressoDeploymentInput `json:"input"`
	// Configuration for quick sync.
	QuickSync EcspressoDeployStageOptions `json:"quickSync"`
}

func (s *EcspressoApplicationSpec) Validate() error {
	// TODO: Validate EcspressoApplicationSpec fields.
	return nil
}

type EcspressoDeploymentInput struct {
	// Config is the path to the ecspresso config file. This will be used as `ecspresso deploy --config <Config>`
	Config string `json:"config"`
}

// EcspressoDeployStageOptions contains all configurable values for a ECSPRESSO_SYNC stage.
type EcspressoDeployStageOptions struct {
}

// EcspressoDiffStageOptions contains all configurable values for a ECSPRESSO_PLAN stage.
type EcspressoDiffStageOptions struct {
	// Exit the pipeline if the result is "No Changes" with success status.
	ExitOnNoChanges bool `json:"exitOnNoChanges"`
}