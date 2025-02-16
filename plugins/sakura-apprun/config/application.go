package config

import (
	"errors"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
)

// AppRunApplicationSpec represents an application configuration for app-run application.
type AppRunApplicationSpec struct {
	config.GenericApplicationSpec
	// Input for AppRun deployment.
	Input AppRunDeploymentInput `json:"input"`
	// Configuration for quick sync.
	QuickSync AppRunDeployStageOptions `json:"quickSync"`
}

func (s *AppRunApplicationSpec) Validate() error {
	if err := s.Input.validate(); err != nil {
		return err
	}
	return nil
}

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
