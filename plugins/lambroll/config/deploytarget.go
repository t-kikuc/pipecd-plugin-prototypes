package config

// LambrollDeployTargetConfig represents PipedDeployTarget.Config for lambroll plugin.
type LambrollDeployTargetConfig struct {
	// Version is the version of lambroll to use. e.g. "v1.1.3"
	Version string `json:"version"`
	// TODO: Add fields if needed.
}
