package config

import (
	"errors"
)

type AppRunDeploymentInput struct {
	// Path to the AppRun configuration file.
	ConfigFile string `json:"configFile"`
}

func (i *AppRunDeploymentInput) validate() error {
	if i.ConfigFile == "" {
		return errors.New("configFile is required")
	}
	return nil
}

// AppRunDeployStageOptions contains all configurable values for a APPRUN_DEPLOY stage.
type AppRunDeployStageOptions struct {
}
