package config

import (
	"encoding/json"
	"fmt"

	pipedconfig "github.com/pipe-cd/pipecd/pkg/configv1"
)

// CDKDeployTargetConfig represents PipedDeployTarget.Config for cdk plugin.
type CDKDeployTargetConfig struct {
	// Version is the version of cdk to use. e.g. "v1.1.3"
	Version string `json:"version"`
	// TODO: Add fields if needed.
}

func ParseDeployTargetConfig(deployTarget pipedconfig.PipedDeployTarget) (CDKDeployTargetConfig, error) {
	var cfg CDKDeployTargetConfig

	if err := json.Unmarshal(deployTarget.Config, &cfg); err != nil {
		return CDKDeployTargetConfig{}, fmt.Errorf("failed to unmarshal deploy target configuration: %w", err)
	}

	return cfg, nil
}
