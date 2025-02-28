package config

import (
	config "github.com/pipe-cd/pipecd/pkg/configv1"
)

// CDKApplicationSpec represents an application configuration for cdk application.
type CDKApplicationSpec struct {
	config.GenericApplicationSpec
	// Input for cdk deployment. e.g. cdk version
	Input CDKDeploymentInput `json:"input"`
	// Configuration for quick sync.
	QuickSync CDKDeployStageOptions `json:"quickSync"`
}

func (s *CDKApplicationSpec) Validate() error {
	if err := s.Input.validate(); err != nil {
		return err
	}
	return nil
}

type CDKDeploymentInput struct {
}

func (i *CDKDeploymentInput) validate() error {
	return nil
}

// CDKDeployStageOptions contains all configurable values for a CDK_SYNC stage.
type CDKDeployStageOptions struct {
}

// CDKDiffStageOptions contains all configurable values for a CDK_PLAN stage.
type CDKDiffStageOptions struct {
	// Exit the pipeline if the result is "No Changes" with success status.
	// ExitOnNoChanges bool `json:"exitOnNoChanges"`
}
