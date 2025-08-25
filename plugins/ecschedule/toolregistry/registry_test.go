package toolregistry

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/pipe-cd/piped-plugin-sdk-go/toolregistry/toolregistrytest"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestRegistry_Ecschedule(t *testing.T) {
	t.Parallel()

	c := toolregistrytest.NewTestToolRegistry(t)

	r := NewRegistry(c)

	p, err := r.Ecschedule(context.Background(), "v0.12.0")
	require.NoError(t, err)
	require.NotEmpty(t, p)

	out, err := exec.CommandContext(context.Background(), p, "-version").CombinedOutput()
	require.NoError(t, err)

	expected := "ecschedule v0.12.0 (rev:586cc61)"

	assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(string(out)))
}
