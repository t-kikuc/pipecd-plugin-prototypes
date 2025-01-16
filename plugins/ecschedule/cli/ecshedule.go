// Copyright 2024 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
