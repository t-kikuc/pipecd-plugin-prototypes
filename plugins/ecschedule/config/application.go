package config

type EcscheduleDeploymentInput struct {
	// Config is the path to the ecschedule config file. This will be used as `ecschedule deploy --config <Config>`
	Config string `json:"config"`
}

// EcscheduleApplyStageOptions contains all configurable values for a ECSPRESSO_SYNC stage.
type EcscheduleApplyStageOptions struct {
}

// EcscheduleDiffStageOptions contains all configurable values for a ECSPRESSO_PLAN stage.
type EcscheduleDiffStageOptions struct {
	// Exit the pipeline if the result is "No Changes" with success status.
	ExitOnNoChanges bool `json:"exitOnNoChanges"`
}
