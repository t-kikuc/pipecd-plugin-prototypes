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

type CdkCli struct {
	execPath string
	dir      string

	functionFilePath string // TODO
	srcPath          string // TODO

	options options
}

func NewCDK(execPath, dir, functionFilePath, srcPath string, opts ...Option) *CdkCli {
	opt := options{}
	for _, o := range opts {
		o(&opt)
	}

	return &CdkCli{
		execPath:         execPath,
		dir:              dir,
		functionFilePath: functionFilePath,
		srcPath:          srcPath,
		options:          opt,
	}
}

func (e *CdkCli) Version(ctx context.Context) (string, error) {
	args := []string{"version"}
	cmd := exec.CommandContext(ctx, e.execPath, args...)
	cmd.Dir = e.dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}

	return strings.TrimSpace(string(out)), nil
}

func (e *CdkCli) Deploy(ctx context.Context, w io.Writer) error {
	args := []string{
		"deploy",
		// "--function", e.functionFilePath,
		// "--src", e.srcPath,
	}

	cmd := exec.CommandContext(ctx, e.execPath, args...)
	cmd.Dir = e.dir
	cmd.Stdout = w
	cmd.Stderr = w

	fmt.Fprintf(w, "execute 'cdk %s'\n", strings.Join(args, " "))
	return cmd.Run()
}

func (e *CdkCli) Diff(ctx context.Context, w io.Writer) error {
	args := []string{
		"diff",
		// "--function", e.functionFilePath,
		// "--src", e.srcPath,
	}

	cmd := exec.CommandContext(ctx, e.execPath, args...)
	cmd.Dir = e.dir
	cmd.Stdout = w
	cmd.Stderr = w

	fmt.Fprintf(w, "execute 'cdk %s'\n", strings.Join(args, " "))
	return cmd.Run()
}
