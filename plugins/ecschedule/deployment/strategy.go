package deployment

import (
	"github.com/pipe-cd/pipecd/pkg/model"
	ecspconfig "github.com/t-kikuc/pipecd-plugin-prototypes/ecschedule/config"
)

func determineStrategy(spec ecspconfig.EcscheduleApplicationSpec) (strategy model.SyncStrategy, summary string, err error) {
	if spec.Pipeline == nil || len(spec.Pipeline.Stages) == 0 {
		return model.SyncStrategy_QUICK_SYNC,
			"Quick sync by executing 'ecschedule apply' because no pipeline was configured",
			nil
	} else {
		return model.SyncStrategy_PIPELINE,
			"Sync with the specified pipeline",
			nil
	}
}
