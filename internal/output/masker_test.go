package output

import "testing"

func TestMaskerKeepsServerKeyPathVisible(t *testing.T) {
	masker := NewSecretMasker()
	masker.AddFromEnv(map[string]string{
		"SERVER_KEY":     "/home/root/.ssh/id_ed25712",
		"DASHBOARD_PASS": "pw",
	})

	maskedPath := masker.Mask("using key /home/root/.ssh/id_ed25712")
	if maskedPath != "using key /home/root/.ssh/id_ed25712" {
		t.Fatalf("unexpected masking for SERVER_KEY path: %q", maskedPath)
	}

	maskedPassword := masker.Mask("password=pw")
	if maskedPassword != "password=***" {
		t.Fatalf("expected short password to be masked, got %q", maskedPassword)
	}
}
