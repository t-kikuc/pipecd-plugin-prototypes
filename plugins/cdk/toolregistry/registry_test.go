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

func TestRegistry_CDK(t *testing.T) {
	t.Parallel()

	c, err := toolregistrytest.NewToolRegistry(t)
	require.NoError(t, err)

	r := NewRegistry(c)

	t.Cleanup(func() { c.Close() })

	p, err := r.CDK(context.Background(), "v1.1.3")
	require.NoError(t, err)
	require.NotEmpty(t, p)

	out, err := exec.CommandContext(context.Background(), p, "version").CombinedOutput()
	require.NoError(t, err)

	expected := "cdk v1.1.3"

	assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(string(out)))
}
