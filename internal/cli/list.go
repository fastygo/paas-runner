package cli

import (
	"fmt"

	"github.com/paas/paas-runner/internal/config"
)

func listCommand() error {
	projectConfig, err := config.LoadProjectConfig()
	if err != nil {
		return err
	}

	extensions, err := listExtensions(projectConfig.ExtensionsDirOrDefault(), config.UserExtensionsDir())
	if err != nil {
		return err
	}

	if len(extensions) == 0 {
		fmt.Println("No extensions found")
		return nil
	}

	fmt.Println("Available extensions:")
	for _, item := range extensions {
		fmt.Printf("  %s (%s)\n", item.Name, item.Source)
	}

	return nil
}
