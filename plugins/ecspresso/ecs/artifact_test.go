package ecs

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
	"github.com/stretchr/testify/assert"
)

func TestExtractArtifactVersions(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		input       ecs.RegisterTaskDefinitionInput
		expected    []*sdk.ArtifactVersion
		expectedErr bool
	}{
		{
			name: "one container",
			input: ecs.RegisterTaskDefinitionInput{
				ContainerDefinitions: []types.ContainerDefinition{
					{
						Image: aws.String("gcr.io/pipecd/helloworld:v1.0.0"),
					},
				},
			},
			expected: []*sdk.ArtifactVersion{
				{
					Kind:    sdk.ArtifactKindContainerImage,
					Version: "v1.0.0",
					Name:    "helloworld",
					URL:     "gcr.io/pipecd/helloworld:v1.0.0",
				},
			},
			expectedErr: false,
		},
		{
			name:        "missing containerDefinitions",
			input:       ecs.RegisterTaskDefinitionInput{},
			expected:    nil,
			expectedErr: true,
		},
		{
			name:        "missing image name",
			input:       ecs.RegisterTaskDefinitionInput{},
			expected:    nil,
			expectedErr: true,
		},
		{
			name: "multiple containers",
			input: ecs.RegisterTaskDefinitionInput{
				ContainerDefinitions: []types.ContainerDefinition{
					{
						Image: aws.String("gcr.io/pipecd/helloworld:v1.0.0"),
					},
					{
						Image: aws.String("gcr.io/pipecd/my-service:v1.2.3"),
					},
				},
			},
			expected: []*sdk.ArtifactVersion{
				{
					Kind:    sdk.ArtifactKindContainerImage,
					Version: "v1.0.0",
					Name:    "helloworld",
					URL:     "gcr.io/pipecd/helloworld:v1.0.0",
				},
				{
					Kind:    sdk.ArtifactKindContainerImage,
					Version: "v1.2.3",
					Name:    "my-service",
					URL:     "gcr.io/pipecd/my-service:v1.2.3",
				},
			},
			expectedErr: false,
		},
		{
			name: "multiple containers with the same image returns only one version",
			input: ecs.RegisterTaskDefinitionInput{
				ContainerDefinitions: []types.ContainerDefinition{
					{
						Image: aws.String("gcr.io/pipecd/helloworld:v1.0.0"),
					},
					{
						Image: aws.String("gcr.io/pipecd/helloworld:v1.0.0"),
					},
				},
			},
			expected: []*sdk.ArtifactVersion{
				{
					Kind:    sdk.ArtifactKindContainerImage,
					Version: "v1.0.0",
					Name:    "helloworld",
					URL:     "gcr.io/pipecd/helloworld:v1.0.0",
				},
			},
			expectedErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			versions, err := ExtractArtifactVersions(&tc.input)
			assert.Equal(t, tc.expectedErr, err != nil)
			assert.ElementsMatch(t, tc.expected, versions)
		})
	}
}
