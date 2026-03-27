package runner

import (
	"bufio"
	"io"
	"strings"

	"github.com/paas/paas-runner/internal/output"
)

func scanReader(reader io.Reader, stream output.Stream, onLine func(output.Stream, string)) {
	if onLine == nil {
		onLine = func(output.Stream, string) {}
	}

	br := bufio.NewReader(reader)
	for {
		line, err := br.ReadString('\n')
		if len(line) > 0 {
			onLine(stream, strings.TrimRight(line, "\r\n"))
		}

		if err != nil {
			return
		}
	}
}
