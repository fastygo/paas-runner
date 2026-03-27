package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func findBash() (string, bool, error) {
	if path, err := exec.LookPath("bash"); err == nil {
		return path, false, nil
	}

	if programFiles := os.Getenv("ProgramFiles"); programFiles != "" {
		gitBash := filepath.Join(programFiles, "Git", "bin", "bash.exe")
		if _, err := os.Stat(gitBash); err == nil {
			return gitBash, false, nil
		}
	}

	if path, err := exec.LookPath("wsl"); err == nil {
		return path, true, nil
	}

	return "", false, fmt.Errorf("bash not found; install Git for Windows or WSL")
}
