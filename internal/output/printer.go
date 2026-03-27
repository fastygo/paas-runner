package output

import (
	"fmt"
	"io"
	"time"
)

type Stream int

const (
	Stdout Stream = iota
	Stderr
)

type StepStatus int

const (
	StatusSuccess StepStatus = iota
	StatusWarning
	StatusFailed
)

type Printer struct {
	out     io.Writer
	success int
	warning int
	failed  int
	masker  *SecretMasker
}

func NewPrinter(out io.Writer) *Printer {
	return &Printer{out: out, masker: NewSecretMasker()}
}

func (p *Printer) SetMasker(masker *SecretMasker) {
	if masker != nil {
		p.masker = masker
	}
}

func (p *Printer) PrintHeader(serverName, extensionName, extensionDescription string) {
	fmt.Fprintf(p.out, "[paas] Server: %s\n", serverName)

	if extensionDescription != "" {
		fmt.Fprintf(p.out, "[paas] Extension: %s (%s)\n\n", extensionName, extensionDescription)
		return
	}

	fmt.Fprintf(p.out, "[paas] Extension: %s\n\n", extensionName)
}

func (p *Printer) PrintStepHeader(stepNo, total int, description string) {
	if description == "" {
		description = "Step"
	}

	fmt.Fprintf(p.out, "  [%d/%d] %s\n", stepNo, total, description)
}

func (p *Printer) PrintCommand(command string) {
	fmt.Fprintf(p.out, "        $ %s\n", p.mask(command))
}

func (p *Printer) PrintSkipped(reason string) {
	if reason == "" {
		reason = "condition is false"
	}

	fmt.Fprintf(p.out, "        - Skipped (%s)\n\n", reason)
}

func (p *Printer) PrintStream(stream Stream, line string) {
	_ = stream
	if line == "" {
		return
	}

	fmt.Fprintf(p.out, "        %s\n", p.mask(line))
}

func (p *Printer) PrintDryRun() {
	fmt.Fprintf(p.out, "        (skipped)\n\n")
}

func (p *Printer) PrintResult(status StepStatus, exitCode int, elapsed time.Duration) {
	switch status {
	case StatusSuccess:
		p.success++
		fmt.Fprintf(p.out, "        ✓ Done (%s)\n\n", elapsed.Truncate(time.Millisecond))
	case StatusWarning:
		p.warning++
		fmt.Fprintf(p.out, "        ⚠ Ignored failure (exit code %d, %s)\n\n", exitCode, elapsed.Truncate(time.Millisecond))
	case StatusFailed:
		p.failed++
		fmt.Fprintf(p.out, "        ✗ Failed (exit code %d, %s)\n\n", exitCode, elapsed.Truncate(time.Millisecond))
	}
}

func (p *Printer) PrintError(err error) {
	if err != nil {
		fmt.Fprintf(p.out, "        Error: %v\n", p.mask(err.Error()))
	}
}

func (p *Printer) PrintSummary(total time.Duration) {
	if p.failed == 0 && p.warning == 0 {
		fmt.Fprintf(p.out, "All %d steps succeeded (%s)\n", p.success, total.Truncate(time.Millisecond))
		return
	}

	fmt.Fprintf(p.out, "%d succeeded, %d warning, %d failed (%s)\n", p.success, p.warning, p.failed, total.Truncate(time.Millisecond))
}

func (p *Printer) mask(value string) string {
	if p.masker == nil {
		return value
	}

	return p.masker.Mask(value)
}
