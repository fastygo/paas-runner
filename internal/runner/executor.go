package runner

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/paas/paas-runner/internal/dsl"
	"github.com/paas/paas-runner/internal/output"
)

type Executor struct {
	Local   *LocalRunner
	Remote  *RemoteRunner
	Printer *output.Printer
}

type ExecOptions struct {
	DryRun bool
}

func (e *Executor) Execute(ctx context.Context, extension dsl.Extension, baseEnv map[string]string, options ExecOptions) error {
	if e.Printer == nil {
		e.Printer = output.NewPrinter(nil)
	}

	total := len(extension.Steps)
	if total == 0 {
		return fmt.Errorf("extension %q has no steps", extension.ID)
	}

	env := clone(baseEnv)
	captures := make(map[string]string)

	started := time.Now()

	for i, step := range extension.Steps {
		for key, value := range captures {
			env[key] = value
		}

		description := step.Description
		if description == "" {
			if step.ID != "" {
				description = step.ID
			} else {
				description = fmt.Sprintf("step-%d", i+1)
			}
		}

		e.Printer.PrintStepHeader(i+1, total, description)

		shouldRun, err := dsl.EvaluateWhen(step.When, env)
		if err != nil {
			return fmt.Errorf("step %d: %w", i+1, err)
		}
		if !shouldRun {
			e.Printer.PrintSkipped("condition is false")
			continue
		}

		stepEnv, err := dslSubstituteMap(step.Env, env)
		if err != nil {
			return fmt.Errorf("step %d: %w", i+1, err)
		}

		merged := mergeMaps(env, stepEnv)
		command, err := dsl.SubstituteVariables(step.Run, merged)
		if err != nil {
			return fmt.Errorf("step %d: %w", i+1, err)
		}

		e.Printer.PrintCommand(command)

		if options.DryRun {
			e.Printer.PrintDryRun()
			if step.Capture != "" {
				key := "STEP_" + strings.ToUpper(step.Capture)
				env[key] = ""
			}
			continue
		}

		start := time.Now()
		result, err := e.runStep(ctx, step, merged, command)
		if err != nil {
			return err
		}

		elapsed := time.Since(start)

		if result.exitCode == 0 {
			e.Printer.PrintResult(output.StatusSuccess, result.exitCode, elapsed)
			if step.Capture != "" {
				value := captureLast(result.captured)
				key := "STEP_" + strings.ToUpper(step.Capture)
				env[key] = value
				captures[key] = value
			}
			continue
		}

		if step.IgnoreError {
			e.Printer.PrintResult(output.StatusWarning, result.exitCode, elapsed)
			if step.Capture != "" {
				value := captureLast(result.captured)
				key := "STEP_" + strings.ToUpper(step.Capture)
				env[key] = value
				captures[key] = value
			}
			continue
		}

		e.Printer.PrintResult(output.StatusFailed, result.exitCode, elapsed)
		if result.err != nil {
			e.Printer.PrintError(result.err)
		}

		return fmt.Errorf("step %d failed", i+1)
	}

	e.Printer.PrintSummary(time.Since(started))

	return nil
}

type stepResult struct {
	exitCode int
	err      error
	captured []string
}

func (e *Executor) runStep(ctx context.Context, step dsl.Step, env map[string]string, command string) (stepResult, error) {
	selected := Runner(e.Local)
	if !step.Local {
		selected = Runner(e.Remote)
	}

	if selected == nil {
		return stepResult{}, fmt.Errorf("runner is not configured")
	}

	if step.Local {
		e.Local.SetWorkdir(step.Workdir)
	} else {
		e.Remote.SetWorkdir(step.Workdir)
	}

	lines := make([]string, 0)
	envSlice := toEnvSlice(env)

	result, err := selected.Run(ctx, command, envSlice, func(stream output.Stream, line string) {
		e.Printer.PrintStream(stream, line)
		if step.Capture != "" && stream == output.Stdout {
			lines = append(lines, line)
		}
	})

	return stepResult{exitCode: result.ExitCode, err: err, captured: lines}, nil
}

func dslSubstituteMap(values map[string]string, env map[string]string) (map[string]string, error) {
	out := make(map[string]string)
	for key, value := range values {
		resolved, err := dsl.SubstituteVariables(value, env)
		if err != nil {
			return nil, fmt.Errorf("env %q: %w", key, err)
		}

		out[key] = resolved
	}

	return out, nil
}

func captureLast(lines []string) string {
	for i := len(lines) - 1; i >= 0; i-- {
		if value := strings.TrimSpace(lines[i]); value != "" {
			return value
		}
	}

	return ""
}

func toEnvSlice(env map[string]string) []string {
	out := make([]string, 0, len(env))
	for key, value := range env {
		out = append(out, key+"="+value)
	}

	return out
}

func clone(input map[string]string) map[string]string {
	output := make(map[string]string, len(input))
	for key, value := range input {
		output[key] = value
	}

	return output
}

func mergeMaps(base map[string]string, overlay map[string]string) map[string]string {
	out := clone(base)
	for key, value := range overlay {
		out[key] = value
	}

	return out
}
