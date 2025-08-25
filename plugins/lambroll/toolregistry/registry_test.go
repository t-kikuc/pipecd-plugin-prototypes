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

func TestRegistry_Lambroll(t *testing.T) {
	t.Parallel()

	c := toolregistrytest.NewTestToolRegistry(t)

	r := NewRegistry(c)

	p, err := r.Lambroll(context.Background(), "v1.1.3")
	require.NoError(t, err)
	require.NotEmpty(t, p)

	out, err := exec.CommandContext(context.Background(), p, "version").CombinedOutput()
	require.NoError(t, err)

	expected := "lambroll v1.1.3"

	assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(string(out)))
}
