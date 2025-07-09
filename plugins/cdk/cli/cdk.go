package cli

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	sdktoolregistry "github.com/pipe-cd/piped-plugin-sdk-go/toolregistry"
	"github.com/t-kikuc/pipecd-plugin-prototypes/cdk/config"
	"github.com/t-kikuc/pipecd-plugin-prototypes/cdk/toolregistry"
)

type options struct {
}

type Option func(*options)

type CDK struct {
	execPath string
	dir      string

	dtCfg config.DeployTargetConfig

	options options
}

func NewCDK(
	ctx context.Context,
	tr *sdktoolregistry.ToolRegistry,
	appDir string,
	dtCfg config.DeployTargetConfig,
	opts ...Option,
) (*CDK, error) {
	cdktr := toolregistry.NewRegistry(tr)
	cdkPath, err := cdktr.CDK(ctx, dtCfg.NodeVersion, dtCfg.CDKVersion)
	if err != nil {
		return nil, err
	}

	opt := options{}
	for _, o := range opts {
		o(&opt)
	}

	return &CDK{
		execPath: cdkPath,
		dir:      appDir,
		dtCfg:    dtCfg,
		options:  opt,
	}, nil
}

func (c *CDK) Version(ctx context.Context) (string, error) {
	args := []string{"--version"}
	cmd := exec.CommandContext(ctx, c.execPath, args...)
	cmd.Dir = c.dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}

	return strings.TrimSpace(string(out)), nil
}

func (c *CDK) Deploy(ctx context.Context, w io.Writer, input config.DeploymentInput) error {
	args := []string{
		"deploy",
		stacksArgs(input),
		contextsArgs(input),
		"--require-approval", "never", // Skip approval for security-sensitive changes
		// "--no-rollback",
		"--profile", c.dtCfg.Profile,
	}

	return c.execute(ctx, w, args)
}

func (c *CDK) Diff(ctx context.Context, w io.Writer, input config.DeploymentInput) error {
	args := []string{
		"diff",
		stacksArgs(input),
		contextsArgs(input),
		"--profile", c.dtCfg.Profile,
	}
	return c.execute(ctx, w, args)
}

func (c *CDK) execute(ctx context.Context, w io.Writer, args []string) error {
	if err := c.npmInstall(ctx, w); err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, c.execPath, args...)
	cmd.Dir = c.dir
	cmd.Stdout = w
	cmd.Stderr = w

	fmt.Fprintf(w, "execute 'cdk %s'\n", strings.Join(args, " "))
	return cmd.Run()
}

// FIXME: Use npm installed by toolregistry
func (c *CDK) npmInstall(ctx context.Context, w io.Writer) error {
	cmd := exec.CommandContext(ctx, "npm", "install")
	cmd.Dir = c.dir
	cmd.Stdout = w
	cmd.Stderr = w

	fmt.Fprintf(w, "execute 'npm install'\n")
	return cmd.Run()
}

func stacksArgs(input config.DeploymentInput) string {
	if len(input.Stacks) == 0 {
		return "--all"
	}
	return strings.Join(input.Stacks, " ")
}

func contextsArgs(input config.DeploymentInput) string {
	if len(input.Contexts) == 0 {
		return ""
	}
	return "--context " + strings.Join(input.Contexts, " --context ")
}
