package cli

import (
	"fmt"

	"github.com/paas/paas-runner/internal/config"
)

func serversCommand() error {
	userConfig, err := config.LoadUserConfig()
	if err != nil {
		return err
	}

	if len(userConfig.Servers) == 0 {
		fmt.Println("No servers configured")
		return nil
	}

	fmt.Println("Configured servers:")
	for name, server := range userConfig.Servers {
		fmt.Printf("  %s -> %s@%s:%d\n", name, server.User, server.Host, server.Port)
	}

	return nil
}
