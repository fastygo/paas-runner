package cli

import (
	"flag"
	"fmt"

	"github.com/paas/paas-runner/internal/config"
	"github.com/paas/paas-runner/internal/dsl"
)

func validateCommand(args []string) error {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}

	rest := fs.Args()
	if len(rest) == 0 {
		return fmt.Errorf("missing extension name")
	}

	extensionName := normalizeExtensionName(rest[0])

	projectConfig, err := config.LoadProjectConfig()
	if err != nil {
		return err
	}

	content, _, err := findExtension(extensionName, projectConfig.ExtensionsDirOrDefault(), config.UserExtensionsDir())
	if err != nil {
		return err
	}

	extension, err := dsl.ParseExtension(content)
	if err != nil {
		return err
	}

	return dsl.ValidateExtension(extension)
}
