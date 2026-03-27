package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/paas/paas-runner/internal/config"
	"github.com/paas/paas-runner/internal/dsl"
	"github.com/paas/paas-runner/internal/output"
	"github.com/paas/paas-runner/internal/runner"
	"github.com/paas/paas-runner/internal/sshclient"
	"golang.org/x/term"
)

type stringList []string

func (s *stringList) String() string {
	return strings.Join(*s, ",")
}

func (s *stringList) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func runCommand(args []string) error {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	selectedServer := fs.String("server", "", "server name")
	dryRun := fs.Bool("dry-run", false, "dry run")
	inputValues := stringList{}
	fs.Var(&inputValues, "input", "override input as key=value")

	if err := fs.Parse(args); err != nil {
		return err
	}

	rest := fs.Args()
	if len(rest) == 0 {
		return fmt.Errorf("missing extension name")
	}

	extensionName := normalizeExtensionName(rest[0])

	parsedInputs, err := parseInputs(inputValues)
	if err != nil {
		return err
	}

	projectConfig, err := config.LoadProjectConfig()
	if err != nil {
		return err
	}

	userConfig, err := config.LoadUserConfig()
	if err != nil {
		return err
	}

	extensionBytes, _, err := findExtension(extensionName, projectConfig.ExtensionsDirOrDefault(), config.UserExtensionsDir())
	if err != nil {
		return err
	}

	extension, err := dsl.ParseExtension(extensionBytes)
	if err != nil {
		return err
	}

	if err := dsl.ValidateExtension(extension); err != nil {
		return err
	}

	hasRemote := false
	hasLocal := false
	for _, step := range extension.Steps {
		if step.Local {
			hasLocal = true
			continue
		}
		hasRemote = true
	}

	var selected config.ServerConfig
	if hasRemote {
		if *dryRun {
			if *selectedServer != "" {
				selected = resolveOptionalServer(*selectedServer, userConfig)
			}
		} else {
			selected, err = resolveServer(*selectedServer, projectConfig, userConfig)
			if err != nil {
				return err
			}
		}
	}

	processEnv := dsl.ProcessEnvironment()
	projectDefaults := config.DefaultConfig().Defaults
	for key, value := range projectConfig.Defaults {
		projectDefaults[key] = value
	}

	baseEnv, err := dsl.BuildBaseEnv(projectDefaults, selected.ToEnv(), processEnv, extension.Inputs, parsedInputs, nil)
	if err != nil {
		if missing, ok := err.(*dsl.MissingRequiredInputError); ok {
			return fmt.Errorf("missing required input %q for extension %q\nPass it via: paas run %s --input %s=<value>", missing.Input, extension.ID, extensionName, missing.Input)
		}
		return err
	}

	printer := output.NewPrinter(os.Stdout)
	masker := output.NewSecretMasker()

	for key, value := range selected.ToEnv() {
		masker.AddFromEnv(map[string]string{key: value})
	}
	for _, input := range extension.Inputs {
		if input.Type == "password" {
			key := dsl.NormalizeInputEnvKey(input.Name)
			masker.AddSecret(baseEnv[key])
		}
	}
	printer.SetMasker(masker)

	serverLabel := selected.Host
	if serverLabel == "" {
		serverLabel = "local"
	}
	printer.PrintHeader(serverLabel, extension.ID, extension.Description)

	var localRunner *runner.LocalRunner
	if hasLocal && !*dryRun {
		shellPath, useWsl, err := findBash()
		if err != nil {
			return fmt.Errorf("bash is required for local steps: %w", err)
		}
		localRunner = runner.NewLocalRunner(shellPath, useWsl)
	}

	ex := &runner.Executor{
		Local:   localRunner,
		Printer: printer,
	}

	if hasRemote {
		if selected.Host != "" {
			prompt := func() ([]byte, error) {
				if !term.IsTerminal(int(os.Stdin.Fd())) {
					return nil, fmt.Errorf("stdin is not a terminal; cannot read passphrase")
				}

				fmt.Fprintf(os.Stderr, "Passphrase for %q: ", selected.Key)
				secret, err := term.ReadPassword(int(os.Stdin.Fd()))
				fmt.Fprintln(os.Stderr)
				return secret, err
			}

			client, err := sshclient.Dial(selected, prompt)
			if err != nil {
				return err
			}

			ex.Remote = runner.NewRemoteRunner(client)
			defer ex.Remote.Close()
		}
	}

	if err := ex.Execute(context.Background(), extension, baseEnv, runner.ExecOptions{DryRun: *dryRun}); err != nil {
		return err
	}

	return nil
}

func parseInputs(values []string) (map[string]string, error) {
	rawInputs := make(map[string]string)

	for _, raw := range values {
		key, value, err := parseInputFlag(raw)
		if err != nil {
			return nil, err
		}

		rawInputs[key] = value
	}

	return dsl.BuildInputEnvFromCLI(rawInputs), nil
}

func resolveServer(name string, project config.ProjectConfig, user config.UserConfig) (config.ServerConfig, error) {
	serverName := name
	if serverName == "" {
		serverName = project.Server
	}

	if serverName == "" {
		return config.ServerConfig{}, fmt.Errorf("server required for remote execution")
	}

	selected, ok := user.Servers[serverName]
	if !ok {
		return config.ServerConfig{}, fmt.Errorf("server %q not found in user config", serverName)
	}

	return selected.ApplyDefaults(), nil
}

func resolveOptionalServer(name string, user config.UserConfig) config.ServerConfig {
	if name == "" {
		return config.ServerConfig{}
	}

	selected, ok := user.Servers[name]
	if !ok {
		return config.ServerConfig{}
	}

	return selected.ApplyDefaults()
}
