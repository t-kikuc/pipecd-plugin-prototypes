package deployment

import (
	"os"

	v1 "github.com/sacloud/apprun-api-go/apis/v1"
	"sigs.k8s.io/yaml"
)

type apprunManifest struct {
	patchBody *v1.PatchApplicationBody
}

func loadManifest(path string) (*apprunManifest, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var patchBody v1.PatchApplicationBody
	if err := yaml.Unmarshal(content, &patchBody); err != nil {
		return nil, err
	}

	return &apprunManifest{patchBody: &patchBody}, nil
}

func (m *apprunManifest) toCreateBody() *v1.PostApplicationBody {
	pb := *m.patchBody
	return &v1.PostApplicationBody{
		Name:           *pb.Name,
		TimeoutSeconds: *pb.TimeoutSeconds,
		Port:           *pb.Port,
		MinScale:       *pb.MinScale,
		MaxScale:       *pb.MaxScale,
	}
}
