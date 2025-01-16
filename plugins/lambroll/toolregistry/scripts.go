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

package toolregistry

const installScript = `
cd {{ .TmpDir }}
curl -L https://github.com/fujiwara/lambroll/releases/download/{{ .Version }}/lambroll_{{ .Version }}_{{ .Os }}_{{ .Arch }}.tar.gz -o lambroll_{{ .Version }}_{{ .Os }}_{{ .Arch }}.tar.gz
tar -zxvf lambroll_{{ .Version }}_{{ .Os }}_{{ .Arch }}.tar.gz
mv lambroll {{ .OutPath }}
`