package deployment

import (
	"github.com/pipe-cd/pipecd/pkg/model"
	config "github.com/t-kikuc/pipecd-plugin-prototypes/sakura-apprun/config"
)

func determineStrategy(spec config.AppRunApplicationSpec) (strategy model.SyncStrategy, summary string, err error) {
	if spec.Pipeline == nil || len(spec.Pipeline.Stages) == 0 {
		return model.SyncStrategy_QUICK_SYNC,
			"Quick sync by creating or updating an AppRun application because no pipeline was configured",
			nil
	} else {
		return model.SyncStrategy_PIPELINE,
			"Sync with the specified pipeline",
			nil
	}
}
