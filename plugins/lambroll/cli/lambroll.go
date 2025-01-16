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

type Lambroll struct {
	execPath   string
	dir        string
	configPath string

	options options
}

func NewLambroll(execPath, dir, configPath string, opts ...Option) *Lambroll {
	opt := options{}
	for _, o := range opts {
		o(&opt)
	}

	return &Lambroll{
		execPath:   execPath,
		dir:        dir,
		configPath: configPath,
		options:    opt,
	}
}

func (e *Lambroll) Version(ctx context.Context) (string, error) {
	args := []string{"version"}
	cmd := exec.CommandContext(ctx, e.execPath, args...)
	cmd.Dir = e.dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}

	return strings.TrimSpace(string(out)), nil
}

func (e *Lambroll) Deploy(ctx context.Context, w io.Writer) error {
	args := []string{
		"deploy",
		"--config", e.configPath,
	}

	cmd := exec.CommandContext(ctx, e.execPath, args...)
	cmd.Dir = e.dir
	cmd.Stdout = w
	cmd.Stderr = w

	fmt.Fprintf(w, "execute 'lambroll %s'", strings.Join(args, " "))
	return cmd.Run()
}

func (e *Lambroll) Diff(ctx context.Context, w io.Writer) error {
	args := []string{
		"diff",
	}

	cmd := exec.CommandContext(ctx, e.execPath, args...)
	cmd.Dir = e.dir
	cmd.Stdout = w
	cmd.Stderr = w

	fmt.Fprintf(w, "execute 'lambroll %s'\n", strings.Join(args, " "))
	return cmd.Run()
}
