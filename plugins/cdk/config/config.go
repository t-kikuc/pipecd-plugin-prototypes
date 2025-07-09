package config

// Config represents the plugin-scoped configuration.
// type Config struct{}

// CDKDeployTargetConfig represents PipedDeployTarget.Config for cdk plugin.
type CDKDeployTargetConfig struct {
	// Version is the version of cdk to use. e.g. "2.1001.0"
	Version string `json:"version"`
	// NodeVersion is the version of node to use. e.g. "v22.14.0"
	NodeVersion string `json:"nodeVersion"`
	// Region is the AWS region to deploy to. e.g. "us-east-1"
	Region string `json:"region"`
	// Profile is the AWS profile to use. e.g. "my-profile"
	Profile string `json:"profile"`

	// TODO: Add fields if needed.
}

// CDKApplicationSpec represents an application configuration for cdk application.
type CDKApplicationSpec struct {
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
	// Stacks is the list of stacks to deploy.
	// If you want to deploy all stacks, set "--all".
	Stacks []string `json:"stacks"`
	// Contexts is the list of context to pass to the cdk deploy command.
	// Each context is a key-value pair like "bucketName=my-bucket"
	Contexts []string `json:"contexts"`
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
