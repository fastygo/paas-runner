package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ensureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

func normalizeExtensionName(raw string) string {
	raw = strings.TrimSpace(raw)
	if strings.HasSuffix(raw, ".yml") {
		return strings.TrimSuffix(raw, filepath.Ext(raw))
	}

	return raw
}

func parseInputFlag(raw string) (string, string, error) {
	parts := strings.SplitN(raw, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid input value %q (expected key=value)", raw)
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	if key == "" {
		return "", "", fmt.Errorf("empty input key")
	}

	return key, value, nil
}
