package config

import (
	"encoding/json"
	"fmt"

	pipedconfig "github.com/pipe-cd/pipecd/pkg/configv1"
)

// EcspressoDeployTargetConfig represents PipedDeployTarget.Config for ecspresso plugin.
type EcspressoDeployTargetConfig struct {
	// Version is the version of ecspresso to use. e.g. "2.4.5"
	// Do not specify the prefix "v".
	Version string `json:"version"`
	// TODO: Add fields if needed.
}

func ParseDeployTargetConfig(deployTarget pipedconfig.PipedDeployTarget) (EcspressoDeployTargetConfig, error) {
	var cfg EcspressoDeployTargetConfig

	if err := json.Unmarshal(deployTarget.Config, &cfg); err != nil {
		return EcspressoDeployTargetConfig{}, fmt.Errorf("failed to unmarshal deploy target configuration: %w", err)
	}

	return cfg, nil
}
