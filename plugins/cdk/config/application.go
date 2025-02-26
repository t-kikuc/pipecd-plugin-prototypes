package config

import (
	"errors"

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
	// FunctionFile is the path to the cdk function file. This will be used as `cdk deploy --function <functionFile>`
	FunctionFile string `json:"functionFile"`

	// Source is the path to the cdk source directory. This will be used as `cdk deploy --src <source>`
	Source string `json:"source"`
}

func (i *CDKDeploymentInput) validate() error {
	if i.FunctionFile == "" {
		return errors.New("functionFile is required")
	}
	if i.Source == "" {
		return errors.New("source is required")
	}
	return nil
}

// CDKDeployStageOptions contains all configurable values for a CDK_SYNC stage.
type CDKDeployStageOptions struct {
}

// CDKDiffStageOptions contains all configurable values for a CDK_PLAN stage.
type CDKDiffStageOptions struct {
	// Exit the pipeline if the result is "No Changes" with success status.
	ExitOnNoChanges bool `json:"exitOnNoChanges"`
}
