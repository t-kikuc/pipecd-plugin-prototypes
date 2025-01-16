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
