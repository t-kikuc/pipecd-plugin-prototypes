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
	// cp -L {{ .TmpDir }}/node_modules/.bin/cdk {{ .OutPath }}

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

	fmt.Printf("[DEBUG] script: %s\n----------\n", buf.String())

	return r.client.InstallTool(ctx, "cdk", version, buf.String())
}

func (r *Registry) CDK__(ctx context.Context, nodeVersion, version string) (path string, err error) {

	nodePath, err := r.client.InstallTool(ctx, "node", nodeVersion, nodeInstallScript)
	if err != nil {
		fmt.Printf("[DEBUG-ERROR] nodePath: %s, err: %s\n", nodePath, err)
		return "", err
	}

	fmt.Printf("------------\n[DEBUG] nodePath: %s\n", nodePath)

	// cdkInstallScript := fmt.Sprintf(`
	// cd {{ .TmpDir }}
	// export PATH="{{ .TmpDir }}/bin/node-v{{ .Version }}:$PATH"
	// export npm_config_prefix="{{ .TmpDir }}"
	// npm install aws-cdk@{{ .Version }}
	// cp -L {{ .TmpDir }}/node_modules/.bin/cdk {{ .OutPath }}
	// `)
	// // cp -L %s/aws-cdk/bin/cdk {{ .OutPath }}
	cdkInstallScript := fmt.Sprintf(`
		cd {{ .TmpDir }}
		mkdir -p node_modules
		echo "Node version: $(node --version)"
		echo "NPM version: $(npm --version)"
		%s install --no-package-lock aws-cdk@{{ .Version }}
		cp -L {{ .TmpDir }}/node_modules/aws-cdk/bin/cdk {{ .OutPath }}
	`, nodePath)
	fmt.Printf("[DEBUG] cdkInstallScript: %s\n", cdkInstallScript)

	return r.client.InstallTool(ctx, "cdk", version, cdkInstallScript)
}
