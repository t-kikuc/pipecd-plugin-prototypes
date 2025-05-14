package config

// EcspressoDeployTargetConfig represents PipedDeployTarget.Config for ecspresso plugin.
type EcspressoDeployTargetConfig struct {
	// Version is the version of ecspresso to use. e.g. "2.4.5"
	// Do not specify the prefix "v".
	Version string `json:"version"`
	// TODO: Add fields if needed.
}

// EcspressoApplicationSpec represents an app configuration for Ecspresso application (i.e. spec.plugins.xxx of app.pipecd.yaml).
type EcspressoApplicationSpec struct {
	// Input for ecspresso deployment. e.g. ecspresso version, workspace
	Input EcspressoDeploymentInput `json:"input"`
	// Configuration for quick sync.
	QuickSync EcspressoDeployStageOptions `json:"quickSync"`
}

func (s *EcspressoApplicationSpec) Validate() error {
	// TODO: Validate if needed
	return nil
}

// EcspressoDeploymentInput is the input for ecspresso stages.
type EcspressoDeploymentInput struct {
	// Config is the path to the ecspresso config file. This will be used as `ecspresso deploy --config <config>`
	Config string `json:"config"`
}

// EcspressoDeployStageOptions contains all configurable fields for an ECSPRESSO_DEPLOY stage.
type EcspressoDeployStageOptions struct {
}

// EcspressoDiffStageOptions contains all configurable fields for an ECSPRESSO_DIFF stage.
type EcspressoDiffStageOptions struct {
	// Exit the pipeline if the result is "No Changes" with success status.
	ExitOnNoChanges bool `json:"exitOnNoChanges"`
}
