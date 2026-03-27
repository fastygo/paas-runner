package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/paas/paas-runner/internal/config"
	"github.com/paas/paas-runner/internal/extensions"
)

const defaultConfigTemplate = `server: production

defaults:
  INPUT_APP_NAME: starter
  INPUT_REGISTRY_HOST: registry.example.com
  INPUT_IMAGE_REPOSITORY: starter/app
  INPUT_DASHBOARD_URL: https://dashboard.example.com

extensions_dir: .paas/extensions
`

func initCommand(args []string) error {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	extract := fs.String("extract", "", "extract built-in extension name to project extensions directory")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := ensureDir(".paas"); err != nil {
		return err
	}

	configPath := config.ProjectConfigPath()
	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("%s already exists", configPath)
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := os.WriteFile(configPath, []byte(defaultConfigTemplate), 0o644); err != nil {
		return err
	}

	if *extract == "" {
		return nil
	}

	name := normalizeExtensionName(*extract)
	content, err := extensions.Read(name + ".yml")
	if err != nil {
		return err
	}

	targetDir := filepath.Join(".paas", "extensions")
	if err := ensureDir(targetDir); err != nil {
		return err
	}

	target := filepath.Join(targetDir, name+".yml")
	if err := os.WriteFile(target, content, 0o644); err != nil {
		return fmt.Errorf("write extension %q: %w", name, err)
	}

	return nil
}
