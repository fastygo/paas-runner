package output

import "strings"

type SecretMasker struct {
	secrets []string
}

func NewSecretMasker() *SecretMasker {
	return &SecretMasker{secrets: make([]string, 0)}
}

func (m *SecretMasker) AddSecret(value string) {
	if value == "" {
		return
	}

	for _, existing := range m.secrets {
		if existing == value {
			return
		}
	}

	m.secrets = append(m.secrets, value)
}

func (m *SecretMasker) AddFromEnv(env map[string]string) {
	for key, value := range env {
		if isSensitiveEnvKey(key) {
			m.AddSecret(value)
		}
	}
}

func (m *SecretMasker) AddInputSecrets(inputs map[string]bool, values map[string]string) {
	for key, isSecret := range inputs {
		if !isSecret {
			continue
		}

		if value, ok := values[key]; ok {
			m.AddSecret(value)
		}
	}
}

func (m *SecretMasker) Mask(value string) string {
	result := value

	for _, secret := range m.secrets {
		result = strings.ReplaceAll(result, secret, "***")
	}

	return result
}

func isSensitiveEnvKey(key string) bool {
	keyLower := strings.ToLower(key)

	sensitiveParts := []string{
		"pass",
		"password",
		"token",
		"secret",
	}

	for _, part := range sensitiveParts {
		if strings.Contains(keyLower, part) {
			return true
		}
	}

	return false
}
