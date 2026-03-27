package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/paas/paas-runner/internal/extensions"
)

type ExtensionSource struct {
	Name   string
	Source string
}

func findExtension(name string, projectDir string, userDir string) ([]byte, string, error) {
	normalized := name
	if filepath.Ext(normalized) != ".yml" {
		normalized += ".yml"
	}

	projectPath := filepath.Join(projectDir, normalized)
	if data, err := os.ReadFile(projectPath); err == nil {
		return data, "project", nil
	}

	userPath := filepath.Join(userDir, normalized)
	if data, err := os.ReadFile(userPath); err == nil {
		return data, "user", nil
	}

	if data, err := extensions.Read(normalized); err == nil {
		return data, "embedded", nil
	}

	return nil, "", fmt.Errorf("extension %q not found", name)
}

func listExtensions(projectDir string, userDir string) ([]ExtensionSource, error) {
	found := map[string]ExtensionSource{}

	projectEntries, err := os.ReadDir(projectDir)
	if err == nil {
		collectEntries(projectEntries, "project", found)
	}

	userEntries, err := os.ReadDir(userDir)
	if err == nil {
		collectEntries(userEntries, "user", found)
	}

	builtins, err := extensions.List()
	if err == nil {
		for _, name := range builtins {
			if _, ok := found[name]; !ok {
				found[name] = ExtensionSource{Name: name, Source: "embedded"}
			}
		}
	}

	out := make([]ExtensionSource, 0, len(found))
	for _, item := range found {
		out = append(out, item)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})

	return out, nil
}

func collectEntries(entries []os.DirEntry, source string, target map[string]ExtensionSource) {
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if filepath.Ext(name) != ".yml" {
			continue
		}
		name = name[:len(name)-4]

		if _, ok := target[name]; ok {
			continue
		}

		target[name] = ExtensionSource{Name: name, Source: source}
	}
}
