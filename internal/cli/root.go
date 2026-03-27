package cli

import (
	"fmt"
)

func Run(args []string) error {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		printUsage()
		return nil
	}

	command := args[0]
	switch command {
	case "run":
		return runCommand(args[1:])
	case "validate":
		return validateCommand(args[1:])
	case "list":
		return listCommand()
	case "init":
		return initCommand(args[1:])
	case "servers":
		return serversCommand()
	case "help":
		printUsage()
		return nil
	default:
		return fmt.Errorf("unknown command %q", command)
	}
}

func printUsage() {
	fmt.Println("paas CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  paas run <extension> [--server name] [--input key=value] [--dry-run]")
	fmt.Println("  paas validate <extension>")
	fmt.Println("  paas list")
	fmt.Println("  paas init")
	fmt.Println("  paas servers")
}
