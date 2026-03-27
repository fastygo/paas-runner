package runner

import (
	"context"
	"strings"

	"github.com/paas/paas-runner/internal/output"
	"golang.org/x/crypto/ssh"
)

type RemoteRunner struct {
	client  *ssh.Client
	workdir string
}

func NewRemoteRunner(client *ssh.Client) *RemoteRunner {
	return &RemoteRunner{client: client}
}

func (r *RemoteRunner) SetWorkdir(workdir string) {
	r.workdir = workdir
}

func (r *RemoteRunner) Run(ctx context.Context, command string, env []string, onLine func(output.Stream, string)) (Result, error) {
	session, err := r.client.NewSession()
	if err != nil {
		return Result{ExitCode: -1}, err
	}
	defer session.Close()

	for _, entry := range env {
		if parts := strings.SplitN(entry, "=", 2); len(parts) == 2 {
			_ = session.Setenv(parts[0], parts[1])
		}
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		return Result{ExitCode: -1}, err
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		return Result{ExitCode: -1}, err
	}

	wrapped := buildRemoteCommand(command, r.workdir, env)
	if err := session.Start(wrapped); err != nil {
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

	err = session.Wait()
	if err != nil {
		if exitErr, ok := err.(*ssh.ExitError); ok {
			return Result{ExitCode: exitErr.ExitStatus()}, err
		}
		return Result{ExitCode: -1}, err
	}

	if ctx.Err() != nil {
		return Result{ExitCode: -1}, ctx.Err()
	}

	return Result{ExitCode: 0}, nil
}

func (r *RemoteRunner) Close() error {
	if r.client == nil {
		return nil
	}

	return r.client.Close()
}

func buildRemoteCommand(command string, workdir string, env []string) string {
	var b strings.Builder

	for _, entry := range env {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 {
			continue
		}

		b.WriteString("export ")
		b.WriteString(parts[0])
		b.WriteString("=")
		b.WriteString(singleQuote(parts[1]))
		b.WriteString("; ")
	}

	if workdir != "" {
		b.WriteString("cd ")
		b.WriteString(singleQuote(workdir))
		b.WriteString(" && ")
	}

	b.WriteString("bash -lc ")
	b.WriteString(singleQuote(command))

	return b.String()
}

func singleQuote(value string) string {
	escaped := strings.ReplaceAll(value, "'", "'\"'\"'")
	return "'" + escaped + "'"
}
