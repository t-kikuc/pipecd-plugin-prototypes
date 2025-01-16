package config

import (
	"encoding/json"
	"fmt"

	pipedconfig "github.com/pipe-cd/pipecd/pkg/configv1"
)

// LambrollDeployTargetConfig represents PipedDeployTarget.Config for lambroll plugin.
type LambrollDeployTargetConfig struct {
	// Version is the version of lambroll to use. e.g. "v1.1.3"
	Version string `json:"version"`
	// TODO: Add fields if needed.
}

func ParseDeployTargetConfig(deployTarget pipedconfig.PipedDeployTarget) (LambrollDeployTargetConfig, error) {
	var cfg LambrollDeployTargetConfig

	if err := json.Unmarshal(deployTarget.Config, &cfg); err != nil {
		return LambrollDeployTargetConfig{}, fmt.Errorf("failed to unmarshal deploy target configuration: %w", err)
	}

	return cfg, nil
}
