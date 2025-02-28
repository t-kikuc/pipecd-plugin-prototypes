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
	execPath string
	dir      string

	options options
}

func NewCDK(execPath, dir string, opts ...Option) *CDK {
	opt := options{}
	for _, o := range opts {
		o(&opt)
	}

	return &CDK{
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
	args := []string{
		"deploy",
	}

	cmd := exec.CommandContext(ctx, e.execPath, args...)
	cmd.Dir = e.dir
	cmd.Stdout = w
	cmd.Stderr = w

	fmt.Fprintf(w, "execute 'cdk %s'", strings.Join(args, " "))
	return cmd.Run()
}

func (e *CDK) Diff(ctx context.Context, w io.Writer) error {
	args := []string{
		"diff",
	}

	cmd := exec.CommandContext(ctx, e.execPath, args...)
	cmd.Dir = e.dir
	cmd.Stdout = w
	cmd.Stderr = w

	fmt.Fprintf(w, "execute 'cdk %s'\n", strings.Join(args, " "))
	return cmd.Run()
}
