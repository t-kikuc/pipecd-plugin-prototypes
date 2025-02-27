// Package toolregistry installs and manages the needed tools
// such as kubectl, helm... for executing tasks in pipeline.
package toolregistry

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"
)

type client interface {
	InstallTool(ctx context.Context, name, version, script string) (path string, err error)
}

func NewRegistry(client client) *Registry {
	return &Registry{
		client: client,
	}
}

// Registry provides functions to get path to the needed tools.
type Registry struct {
	client client
}

func (r *Registry) CDK(ctx context.Context, nodeVersion, version string) (path string, err error) {
	script := `
		cd {{ .TmpDir }}
		curl -L https://nodejs.org/dist/{{ .NodeVersion }}/node-{{ .NodeVersion }}-{{ .Os }}-{{ .Arch }}.tar.gz -o node-{{ .NodeVersion }}-{{ .Os }}-{{ .Arch }}.tar.gz
		tar -zxf node-{{ .NodeVersion }}-{{ .Os }}-{{ .Arch }}.tar.gz

		node-{{ .NodeVersion }}-{{ .Os }}-{{ .Arch }}/bin/npm install aws-cdk@{{ .Version }}
		mkdir -p $(dirname {{ .OutPath }})
		cp -r {{ .TmpDir }}/node_modules/aws-cdk $(dirname {{ .OutPath }})/
		ln -sf $(dirname {{ .OutPath }})/aws-cdk/bin/cdk {{ .OutPath }}
	`

	tmpl, err := template.New("node-install").Delims("[[", "]]").Parse(strings.ReplaceAll(script, "{{ .NodeVersion }}", "[[ .NodeVersion ]]"))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"NodeVersion": nodeVersion,
	})
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return r.client.InstallTool(ctx, "cdk", version, buf.String())
}
