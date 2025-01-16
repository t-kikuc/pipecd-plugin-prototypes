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

type Ecschedule struct {
	execPath   string
	dir        string
	configPath string

	options options
}

func NewEcschedule(execPath, dir, configPath string, opts ...Option) *Ecschedule {
	opt := options{}
	for _, o := range opts {
		o(&opt)
	}

	return &Ecschedule{
		execPath:   execPath,
		dir:        dir,
		configPath: configPath,
		options:    opt,
	}
}

func (e *Ecschedule) Version(ctx context.Context) (string, error) {
	args := []string{"-version"}
	cmd := exec.CommandContext(ctx, e.execPath, args...)
	cmd.Dir = e.dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}

	return strings.TrimSpace(string(out)), nil
}

func (e *Ecschedule) Apply(ctx context.Context, w io.Writer) error {
	args := []string{
		"apply",
		"-conf", e.configPath,
		"-all",
	}

	cmd := exec.CommandContext(ctx, e.execPath, args...)
	cmd.Dir = e.dir
	cmd.Stdout = w
	cmd.Stderr = w

	fmt.Fprintf(w, "execute 'ecschedule %s'", strings.Join(args, " "))
	return cmd.Run()
}

func (e *Ecschedule) Diff(ctx context.Context, w io.Writer) error {
	args := []string{
		"diff",
		"-conf", e.configPath,
		"-all",
	}

	cmd := exec.CommandContext(ctx, e.execPath, args...)
	cmd.Dir = e.dir
	cmd.Stdout = w
	cmd.Stderr = w

	fmt.Fprintf(w, "execute 'ecschedule %s'\n", strings.Join(args, " "))
	return cmd.Run()
}
