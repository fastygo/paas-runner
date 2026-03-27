package runner

import (
	"context"
	"os"
	"os/exec"

	"github.com/paas/paas-runner/internal/output"
)

type LocalRunner struct {
	shell   string
	useWsl  bool
	workdir string
}

func NewLocalRunner(shell string, useWsl bool) *LocalRunner {
	return &LocalRunner{shell: shell, useWsl: useWsl}
}

func (r *LocalRunner) SetWorkdir(workdir string) {
	r.workdir = workdir
}

func (r *LocalRunner) Run(ctx context.Context, command string, env []string, onLine func(output.Stream, string)) (Result, error) {
	cmd := func() *exec.Cmd {
		if r.useWsl {
			return exec.CommandContext(ctx, "wsl", "bash", "-lc", command)
		}
		return exec.CommandContext(ctx, r.shell, "-c", command)
	}()

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, env...)

	if r.workdir != "" {
		cmd.Dir = r.workdir
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return Result{ExitCode: -1}, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return Result{ExitCode: -1}, err
	}

	if err := cmd.Start(); err != nil {
		return Result{ExitCode: -1}, err
	}

	done := make(chan struct{})
	go func() {
		scanReader(stdout, output.Stdout, onLine)
		done <- struct{}{}
	}()
	go func() {
		scanReader(stderr, output.Stderr, onLine)
		done <- struct{}{}
	}()

	<-done
	<-done

	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return Result{ExitCode: exitErr.ExitCode()}, err
		}

		return Result{ExitCode: -1}, err
	}

	return Result{ExitCode: 0}, nil
}

func (r *LocalRunner) Close() error {
	return nil
}
