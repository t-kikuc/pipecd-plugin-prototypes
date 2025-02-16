package config

import (
	"encoding/json"
	"fmt"

	pipedconfig "github.com/pipe-cd/pipecd/pkg/configv1"
)

// AppRunDeployTargetConfig represents PipedDeployTarget.Config for AppRun plugin.
type AppRunDeployTargetConfig struct {
	// TODO: Add fields if needed.
}

func ParseDeployTargetConfig(deployTarget pipedconfig.PipedDeployTarget) (AppRunDeployTargetConfig, error) {
	var cfg AppRunDeployTargetConfig

	if err := json.Unmarshal(deployTarget.Config, &cfg); err != nil {
		return AppRunDeployTargetConfig{}, fmt.Errorf("failed to unmarshal deploy target configuration: %w", err)
	}

	return cfg, nil
}
