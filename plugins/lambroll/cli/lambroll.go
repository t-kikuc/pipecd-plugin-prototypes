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
	execPath string
	dir      string

	functionFilePath string
	srcPath          string

	options options
}

func NewLambroll(execPath, dir, functionFilePath, srcPath string, opts ...Option) *Lambroll {
	opt := options{}
	for _, o := range opts {
		o(&opt)
	}

	return &Lambroll{
		execPath:         execPath,
		dir:              dir,
		functionFilePath: functionFilePath,
		srcPath:          srcPath,
		options:          opt,
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
		"--function", e.functionFilePath,
		"--src", e.srcPath,
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
		"--function", e.functionFilePath,
		"--src", e.srcPath,
	}

	cmd := exec.CommandContext(ctx, e.execPath, args...)
	cmd.Dir = e.dir
	cmd.Stdout = w
	cmd.Stderr = w

	fmt.Fprintf(w, "execute 'lambroll %s'\n", strings.Join(args, " "))
	return cmd.Run()
}
