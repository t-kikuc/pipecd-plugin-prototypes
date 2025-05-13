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
