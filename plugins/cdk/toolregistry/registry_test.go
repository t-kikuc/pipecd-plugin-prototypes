package toolregistry

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/pipe-cd/piped-plugin-sdk-go/toolregistry/toolregistrytest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistry_CDK(t *testing.T) {
	t.Parallel()

	c := toolregistrytest.NewTestToolRegistry(t)
	require.NotNil(t, c)

	r := NewRegistry(c)

	p, err := r.CDK(context.Background(), "v22.13.0", "2.1000.0")
	require.NoError(t, err)
	require.NotEmpty(t, p)

	out, err := exec.CommandContext(context.Background(), p, "--version").CombinedOutput()
	require.NoError(t, err)

	expected := "2.1000.0 (build 1c60bc6)"
	assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(string(out)))
}
