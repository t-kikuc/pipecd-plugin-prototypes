package deployment

import (
	"testing"

	v1 "github.com/sacloud/apprun-api-go/apis/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadManifest(t *testing.T) {
	expected := &v1.PatchApplicationBody{
		Name:           p("app3"),
		TimeoutSeconds: p(300),
		Port:           p(8080),
		MinScale:       p(0),
		MaxScale:       p(1),
		Components: &[]v1.PatchApplicationBodyComponent{
			{
				Name: "test-apprun:v2",
				DeploySource: v1.PatchApplicationBodyComponentDeploySource{
					ContainerRegistry: &v1.PatchApplicationBodyComponentDeploySourceContainerRegistry{
						Image:    "test-apprun.sakuracr.jp/test-apprun:v2",
						Password: nil,
						Server:   p("test-apprun.sakuracr.jp"),
						Username: p("test-user"),
					},
				},
				Env:       &[]v1.PatchApplicationBodyComponentEnv{},
				MaxCpu:    v1.PatchApplicationBodyComponentMaxCpu("1"),
				MaxMemory: v1.PatchApplicationBodyComponentMaxMemory("512Mi"),
				Probe:     nil,
			},
		},
	}

	manifest, err := loadManifest("testdata/app.json")
	require.NoError(t, err)
	assert.Equal(t, expected, manifest.patchBody)
}
