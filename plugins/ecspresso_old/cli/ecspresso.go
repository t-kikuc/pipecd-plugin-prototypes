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

type Ecspresso struct {
	execPath   string
	dir        string
	configPath string

	options options
}

func NewEcspresso(execPath, dir, configPath string, opts ...Option) *Ecspresso {
	opt := options{}
	for _, o := range opts {
		o(&opt)
	}

	return &Ecspresso{
		execPath:   execPath,
		dir:        dir,
		configPath: configPath,
		options:    opt,
	}
}

func (e *Ecspresso) Version(ctx context.Context) (string, error) {
	args := []string{"version"}
	cmd := exec.CommandContext(ctx, e.execPath, args...)
	cmd.Dir = e.dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}

	return strings.TrimSpace(string(out)), nil
}

func (e *Ecspresso) Deploy(ctx context.Context, w io.Writer) error {
	args := []string{
		"deploy",
		"--config", e.configPath,
	}

	cmd := exec.CommandContext(ctx, e.execPath, args...)
	cmd.Dir = e.dir
	cmd.Stdout = w
	cmd.Stderr = w

	fmt.Fprintf(w, "execute: 'ecspresso %s'", strings.Join(args, " "))
	return cmd.Run()
}

func (e *Ecspresso) Diff(ctx context.Context, w io.Writer) error {
	args := []string{
		"diff",
	}

	cmd := exec.CommandContext(ctx, e.execPath, args...)
	cmd.Dir = e.dir
	cmd.Stdout = w
	cmd.Stderr = w

	fmt.Fprintf(w, "execute: 'ecspresso %s'\n", strings.Join(args, " "))
	return cmd.Run()
}
