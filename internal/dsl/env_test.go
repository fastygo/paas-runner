package dsl

import "testing"

func TestBuildBaseEnvUsesNormalizedInputKeys(t *testing.T) {
	defaultTag := ""
	inputs := []Input{
		{
			Name:     "registry-host",
			Type:     "text",
			Required: true,
		},
		{
			Name:    "tag",
			Type:    "text",
			Default: &defaultTag,
		},
	}

	projectDefaults := map[string]string{
		"PAAS_VERSION": "1.0.0",
	}
	serverEnv := map[string]string{
		"SERVER_HOST": "deploy.example.com",
	}
	processEnv := map[string]string{
		"INPUT_TAG": "from-process",
	}
	cliInputs := BuildInputEnvFromCLI(map[string]string{
		"registry-host": "registry.example.com",
		"tag":           "from-cli",
	})
	stepCaptures := map[string]string{
		"STEP_LAST_IMAGE": "registry.example.com/app:sha-123456789abc",
	}

	env, err := BuildBaseEnv(projectDefaults, serverEnv, processEnv, inputs, cliInputs, stepCaptures)
	if err != nil {
		t.Fatalf("BuildBaseEnv returned error: %v", err)
	}

	if got := env["INPUT_REGISTRY_HOST"]; got != "registry.example.com" {
		t.Fatalf("INPUT_REGISTRY_HOST = %q, want %q", got, "registry.example.com")
	}

	if got := env["INPUT_TAG"]; got != "from-cli" {
		t.Fatalf("INPUT_TAG = %q, want %q", got, "from-cli")
	}

	if got := env["STEP_LAST_IMAGE"]; got != "registry.example.com/app:sha-123456789abc" {
		t.Fatalf("STEP_LAST_IMAGE = %q, want captured value", got)
	}
}

func TestBuildBaseEnvMissingRequiredInput(t *testing.T) {
	inputs := []Input{
		{
			Name:     "app-id",
			Type:     "text",
			Required: true,
		},
	}

	_, err := BuildBaseEnv(nil, nil, nil, inputs, nil, nil)
	if err == nil {
		t.Fatal("expected missing required input error")
	}

	missing, ok := err.(*MissingRequiredInputError)
	if !ok {
		t.Fatalf("expected MissingRequiredInputError, got %T", err)
	}

	if missing.Input != "app-id" {
		t.Fatalf("missing.Input = %q, want %q", missing.Input, "app-id")
	}
}
