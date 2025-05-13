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

func TestRegistry_Ecspresso(t *testing.T) {
	t.Parallel()

	c := toolregistrytest.NewTestToolRegistry(t)

	r := NewRegistry(c)

	p, err := r.Ecspresso(context.Background(), "2.4.5")
	require.NoError(t, err)
	require.NotEmpty(t, p)

	out, err := exec.CommandContext(context.Background(), p, "version").CombinedOutput()
	require.NoError(t, err)

	expected := "ecspresso v2.4.5"

	assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(string(out)))
}
