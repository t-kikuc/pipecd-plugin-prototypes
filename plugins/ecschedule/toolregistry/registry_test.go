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

func TestRegistry_Ecschedule(t *testing.T) {
	t.Parallel()

	c, err := toolregistrytest.NewToolRegistry(t)
	require.NoError(t, err)

	r := NewRegistry(c)

	t.Cleanup(func() { c.Close() })

	p, err := r.Ecschedule(context.Background(), "v0.12.0")
	require.NoError(t, err)
	require.NotEmpty(t, p)

	out, err := exec.CommandContext(context.Background(), p, "-version").CombinedOutput()
	require.NoError(t, err)

	expected := "ecschedule v0.12.0 (rev:586cc61)"

	assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(string(out)))
}
