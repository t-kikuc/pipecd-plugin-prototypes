package config

import (
	"fmt"
)

// Config represents the plugin-scoped configuration.
// type Config struct{}

// DeployTargetConfig represents PipedDeployTarget.Config for cdk plugin.
type DeployTargetConfig struct {
	// CDKVersion is the version of cdk to use. e.g. "2.1001.0"
	// This field is required.
	CDKVersion string `json:"cdkVersion"`
	// NodeVersion is the version of node to use. e.g. "v22.14.0"
	// This field is required.
	NodeVersion string `json:"nodeVersion"`
	// Profile is the AWS profile to use. e.g. "my-profile"
	// This field is required.
	Profile string `json:"profile"`

	// TODO: Add fields if needed.
}

func (c *DeployTargetConfig) Validate() error {
	if c.CDKVersion == "" {
		return fmt.Errorf("cdkVersion is required")
	}
	if c.NodeVersion == "" {
		return fmt.Errorf("nodeVersion is required")
	}
	if c.Profile == "" {
		return fmt.Errorf("profile is required")
	}
	return nil
}

// ApplicationSpec represents an application configuration for cdk application.
type ApplicationSpec struct {
	// Input for cdk deployment. e.g. cdk version
	Input DeploymentInput `json:"input"`
	// Configuration for quick sync.
	QuickSync DeployStageOptions `json:"quickSync"`
}

func (s *ApplicationSpec) Validate() error {
	if err := s.Input.validate(); err != nil {
		return err
	}
	return nil
}

type DeploymentInput struct {
	// Stacks is the list of stacks to deploy.
	// If you want to deploy all stacks, set "--all".
	Stacks []string `json:"stacks"`
	// Contexts is the list of context to pass to the cdk deploy command.
	// Each context is a key-value pair like "bucketName=my-bucket"
	Contexts []string `json:"contexts"`
}

func (i *DeploymentInput) validate() error {
	return nil
}

// DeployStageOptions contains all configurable values for a CDK_SYNC stage.
type DeployStageOptions struct {
}

// DiffStageOptions contains all configurable values for a CDK_PLAN stage.
type DiffStageOptions struct {
	// Exit the pipeline if the result is "No Changes" with success status.
	// ExitOnNoChanges bool `json:"exitOnNoChanges"`
}
