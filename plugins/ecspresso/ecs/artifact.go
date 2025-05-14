package ecs

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

// ExtractArtifactVersions extracts artifact versions from the given task definition.
// For example, when taskdef.ContainerDefinitions=[{Image: "gcr.io/pipecd/helloworld:v1.0.0"}, {Image: "gcr.io/pipecd/some-sidecar:v1.0.1"}],
// then the result is:
//
//	[{
//	  Kind: model.ArtifactVersion_CONTAINER_IMAGE,
//	  Version: "v1.0.0",
//	  Name: "helloworld",
//	  Url: "gcr.io/pipecd/helloworld:v1.0.0"
//	}, {
//	  Kind: model.ArtifactVersion_CONTAINER_IMAGE,
//	  Version: "v1.0.1",
//	  Name: "some-sidecar",
//	  Url: "gcr.io/pipecd/some-sidecar:v1.0.1"
//	}]
func ExtractArtifactVersions(taskdef *ecs.RegisterTaskDefinitionInput) ([]*sdk.ArtifactVersion, error) {
	if len(taskdef.ContainerDefinitions) == 0 {
		return nil, fmt.Errorf("container definition must not be empty")
	}

	// Remove duplicate images.
	imageMap := map[string]struct{}{}
	for _, cd := range taskdef.ContainerDefinitions {
		imageMap[*cd.Image] = struct{}{}
	}

	versions := make([]*sdk.ArtifactVersion, 0, len(imageMap))
	for img := range imageMap {
		name, tag := parseContainerImage(img)
		if name == "" {
			return nil, fmt.Errorf("image name must not be empty")
		}

		versions = append(versions, &sdk.ArtifactVersion{
			Kind:    sdk.ArtifactKindContainerImage,
			Version: tag,
			Name:    name,
			URL:     img,
		})
	}

	return versions, nil
}

// parseContainerImage extracts name and tag from the given image.
// For example, when image="gcr.io/pipecd/helloworld:v1.0.0", then (name="helloworld", tag="v1.0.0")
func parseContainerImage(image string) (name, tag string) {
	parts := strings.Split(image, ":")
	if len(parts) == 2 {
		tag = parts[1]
	}
	paths := strings.Split(parts[0], "/")
	name = paths[len(paths)-1]
	return
}
