// Copyright 2024 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
