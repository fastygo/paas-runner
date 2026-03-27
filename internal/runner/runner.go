package runner

import (
	"context"

	"github.com/paas/paas-runner/internal/output"
)

type Result struct {
	ExitCode int
}

type Runner interface {
	Run(ctx context.Context, command string, env []string, onLine func(output.Stream, string)) (Result, error)
	Close() error
}
