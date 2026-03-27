package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	projectConfigFile = ".paas/config.yml"
	userConfigFile    = "servers.yml"
	userConfigDir     = "paas"
)

func ProjectConfigPath() string {
	return projectConfigFile
}

func UserConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(home, ".config", userConfigDir, userConfigFile)
}

func UserExtensionsDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(home, ".config", userConfigDir, "extensions")
}

func DefaultConfig() ProjectConfig {
	return ProjectConfig{
		Server:        "",
		Defaults:      map[string]string{"PAAS_VERSION": "1.0.0"},
		ExtensionsDir: ".paas/extensions",
	}
}

func LoadProjectConfig() (ProjectConfig, error) {
	cfg := DefaultConfig()

	raw, err := os.ReadFile(ProjectConfigPath())
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}

		return cfg, fmt.Errorf("read project config: %w", err)
	}

	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return cfg, fmt.Errorf("parse project config: %w", err)
	}

	cfg.Normalize()

	return cfg, nil
}

func LoadUserConfig() (UserConfig, error) {
	cfg := UserConfig{Servers: map[string]ServerConfig{}}

	path := UserConfigPath()
	if path == "" {
		return cfg, nil
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("read user config: %w", err)
	}

	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return cfg, fmt.Errorf("parse user config: %w", err)
	}

	normalized := make(map[string]ServerConfig)
	for name, server := range cfg.Servers {
		if name == "" {
			continue
		}

		normalized[name] = server.ApplyDefaults()
	}

	cfg.Servers = normalized

	return cfg, nil
}
