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

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/pipe-cd/pipecd/pkg/plugin/toolregistry/toolregistrytest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistry_Lambroll(t *testing.T) {
	t.Parallel()

	c, err := toolregistrytest.NewToolRegistry(t)
	require.NoError(t, err)

	r := NewRegistry(c)

	t.Cleanup(func() { c.Close() })

	p, err := r.Lambroll(context.Background(), "v1.1.3")
	require.NoError(t, err)
	require.NotEmpty(t, p)

	out, err := exec.CommandContext(context.Background(), p, "version").CombinedOutput()
	require.NoError(t, err)

	expected := "lambroll v1.1.3"

	assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(string(out)))
}
