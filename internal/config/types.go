package config

import "fmt"

const (
	defaultExtensionsDir = ".paas/extensions"
)

type ProjectConfig struct {
	Server        string            `yaml:"server"`
	Defaults      map[string]string `yaml:"defaults"`
	ExtensionsDir string            `yaml:"extensions_dir"`
}

type ServerConfig struct {
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	User          string `yaml:"user"`
	Key           string `yaml:"key"`
	DashboardUser string `yaml:"dashboard_user"`
	DashboardPass string `yaml:"dashboard_pass"`
	HostKeyCheck  string `yaml:"host_key_check"`
}

type UserConfig struct {
	Servers map[string]ServerConfig `yaml:"servers"`
}

func (cfg *ProjectConfig) Normalize() {
	if cfg.Defaults == nil {
		cfg.Defaults = map[string]string{}
	}
	if cfg.ExtensionsDir == "" {
		cfg.ExtensionsDir = defaultExtensionsDir
	}
}

func (s ServerConfig) ApplyDefaults() ServerConfig {
	if s.User == "" {
		s.User = "root"
	}

	if s.Port == 0 {
		s.Port = 22
	}

	if s.HostKeyCheck == "" {
		s.HostKeyCheck = "strict"
	}

	return s
}

func (s ServerConfig) ToEnv() map[string]string {
	env := make(map[string]string)

	if s.Host != "" {
		env["SERVER_HOST"] = s.Host
	}
	if s.Port != 0 {
		env["SERVER_PORT"] = fmt.Sprintf("%d", s.Port)
	}
	if s.User != "" {
		env["SERVER_USER"] = s.User
	}
	if s.Key != "" {
		env["SERVER_KEY"] = s.Key
	}
	if s.DashboardUser != "" {
		env["DASHBOARD_USER"] = s.DashboardUser
	}
	if s.DashboardPass != "" {
		env["DASHBOARD_PASS"] = s.DashboardPass
	}
	if s.HostKeyCheck != "" {
		env["HOST_KEY_CHECK"] = s.HostKeyCheck
	}

	return env
}

func (c ProjectConfig) ExtensionsDirOrDefault() string {
	if c.ExtensionsDir == "" {
		return defaultExtensionsDir
	}

	return c.ExtensionsDir
}

func (c ProjectConfig) ServerOrDefault(fallback string) string {
	if c.Server == "" {
		return fallback
	}

	return c.Server
}
