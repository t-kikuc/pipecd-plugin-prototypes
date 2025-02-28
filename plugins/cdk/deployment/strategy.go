package deployment

import (
	"github.com/pipe-cd/pipecd/pkg/model"
	cdkconfig "github.com/t-kikuc/pipecd-plugin-prototypes/cdk/config"
)

func determineStrategy(spec cdkconfig.CDKApplicationSpec) (strategy model.SyncStrategy, summary string, err error) {
	if spec.Pipeline == nil || len(spec.Pipeline.Stages) == 0 {
		return model.SyncStrategy_QUICK_SYNC,
			"Quick sync by executing 'cdk deploy' because no pipeline was configured",
			nil
	} else {
		return model.SyncStrategy_PIPELINE,
			"Sync with the specified pipeline",
			nil
	}
}
