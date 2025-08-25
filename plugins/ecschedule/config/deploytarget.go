package config

// EcscheduleDeployTargetConfig represents PipedDeployTarget.Config for ecspresso plugin.
type EcscheduleDeployTargetConfig struct {
	// Version is the version of ecspresso to use. e.g. "2.4.5"
	// Do not specify the prefix "v".
	Version string `json:"version"`
	// TODO: Add fields if needed.
}
