package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/t-kikuc/pipecd-plugin-prototypes/cdk/config"
)

// TODO: Add tests for other funcs.

func TestStacksArgs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    config.DeploymentInput
		expected string
	}{
		{
			name:     "empty stacks list, return --all",
			input:    config.DeploymentInput{},
			expected: "--all",
		},
		{
			name: "single stack, return stack name",
			input: config.DeploymentInput{
				Stacks: []string{"stack1"},
			},
			expected: "stack1",
		},
		{
			name: "multiple stacks, return stack names separated by spaces",
			input: config.DeploymentInput{
				Stacks: []string{"stack1", "stack2", "stack3"},
			},
			expected: "stack1 stack2 stack3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := stacksArgs(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestContextsArgs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    config.DeploymentInput
		expected string
	}{
		{
			name:     "empty contexts list, return empty string",
			input:    config.DeploymentInput{},
			expected: "",
		},
		{
			name: "single context, return --context flag",
			input: config.DeploymentInput{
				Contexts: []string{"key1=value1"},
			},
			expected: "--context key1=value1",
		},
		{
			name: "multiple contexts, return --context flag with spaces",
			input: config.DeploymentInput{
				Contexts: []string{"key1=value1", "key2=value2", "key3=value3"},
			},
			expected: "--context key1=value1 --context key2=value2 --context key3=value3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := contextsArgs(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}
