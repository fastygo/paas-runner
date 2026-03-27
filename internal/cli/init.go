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
  APP_NAME: starter
  APP_ID: starter-6e62b32b
  IMAGE_REPO: phpfasty/starter
  REGISTRY_HOST: buildy-apps.registry.twcstorage.ru

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

	if err := os.WriteFile(config.ProjectConfigPath(), []byte(defaultConfigTemplate), 0o644); err != nil {
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
