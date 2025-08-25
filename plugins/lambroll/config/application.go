package config

import (
	"errors"
)

type LambrollDeploymentInput struct {
	// FunctionFile is the path to the lambroll function file. This will be used as `lambroll deploy --function <functionFile>`
	FunctionFile string `json:"functionFile"`

	// Source is the path to the lambroll source directory. This will be used as `lambroll deploy --src <source>`
	Source string `json:"source"`
}

func (i *LambrollDeploymentInput) validate() error {
	if i.FunctionFile == "" {
		return errors.New("functionFile is required")
	}
	if i.Source == "" {
		return errors.New("source is required")
	}
	return nil
}

// LambrollDeployStageOptions contains all configurable values for a LAMBROLL_SYNC stage.
type LambrollDeployStageOptions struct {
}

// LambrollDiffStageOptions contains all configurable values for a LAMBROLL_PLAN stage.
type LambrollDiffStageOptions struct {
	// Exit the pipeline if the result is "No Changes" with success status.
	ExitOnNoChanges bool `json:"exitOnNoChanges"`
}
