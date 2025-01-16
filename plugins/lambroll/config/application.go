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

// LambrollApplicationSpec represents an application configuration for lambroll application.
type LambrollApplicationSpec struct {
	config.GenericApplicationSpec
	// Input for lambroll deployment. e.g. lambroll version
	Input LambrollDeploymentInput `json:"input"`
	// Configuration for quick sync.
	QuickSync LambrollDeployStageOptions `json:"quickSync"`
}

func (s *LambrollApplicationSpec) Validate() error {
	// TODO: Validate LambrollApplicationSpec fields.
	return nil
}

type LambrollDeploymentInput struct {
	// Config is the path to the lambroll config file. This will be used as `lambroll deploy --config <Config>`
	Config string `json:"config"`
}

// LambrollDeployStageOptions contains all configurable values for a LAMBROLL_SYNC stage.
type LambrollDeployStageOptions struct {
}

// LambrollDiffStageOptions contains all configurable values for a LAMBROLL_PLAN stage.
type LambrollDiffStageOptions struct {
	// Exit the pipeline if the result is "No Changes" with success status.
	ExitOnNoChanges bool `json:"exitOnNoChanges"`
}
