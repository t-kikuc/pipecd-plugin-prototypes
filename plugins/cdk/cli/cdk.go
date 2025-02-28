package cli

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type options struct {
}

type Option func(*options)

type CDK struct {
	stacks   string
	contexts string

	region  string
	profile string

	execPath string
	dir      string

	options options
}

func NewCDK(stacks, contexts []string, region, profile, execPath, dir string, opts ...Option) *CDK {
	opt := options{}
	for _, o := range opts {
		o(&opt)
	}

	return &CDK{
		stacks:   strings.Join(stacks, " "),
		contexts: strings.Join(contexts, " --context "),
		region:   region,
		profile:  profile,
		execPath: execPath,
		dir:      dir,
		options:  opt,
	}
}

func (e *CDK) Version(ctx context.Context) (string, error) {
	args := []string{"--version"}
	cmd := exec.CommandContext(ctx, e.execPath, args...)
	cmd.Dir = e.dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}

	return strings.TrimSpace(string(out)), nil
}

func (e *CDK) Deploy(ctx context.Context, w io.Writer) error {
	if err := e.npmInstall(ctx, w); err != nil {
		return err
	}

	args := []string{
		"deploy",
		"--stacks", e.stacks,
		"--contexts", e.contexts,
		"--require-approval", "never", // Skip approval for security-sensitive changes
		// "--no-rollback",
		"--region", e.region,
		"--profile", e.profile,
	}

	cmd := exec.CommandContext(ctx, e.execPath, args...)
	cmd.Dir = e.dir
	cmd.Stdout = w
	cmd.Stderr = w

	fmt.Fprintf(w, "execute 'cdk %s'", strings.Join(args, " "))
	return cmd.Run()
}

func (e *CDK) Diff(ctx context.Context, w io.Writer) error {
	if err := e.npmInstall(ctx, w); err != nil {
		return err
	}

	args := []string{
		"diff",
		"--stacks", e.stacks,
		"--contexts", e.contexts,
		"--region", e.region,
		"--profile", e.profile,
	}

	cmd := exec.CommandContext(ctx, e.execPath, args...)
	cmd.Dir = e.dir
	cmd.Stdout = w
	cmd.Stderr = w

	fmt.Fprintf(w, "execute 'cdk %s'\n", strings.Join(args, " "))
	return cmd.Run()
}

// FIXME: Use npm installed by toolregistry
func (e *CDK) npmInstall(ctx context.Context, w io.Writer) error {
	cmd := exec.CommandContext(ctx, "npm", "install")
	cmd.Dir = e.dir
	cmd.Stdout = w
	cmd.Stderr = w

	fmt.Fprintf(w, "execute 'npm install'\n")
	return cmd.Run()
}
