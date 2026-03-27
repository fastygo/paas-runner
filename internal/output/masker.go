package output

import "strings"

type SecretMasker struct {
	secrets []string
}

func NewSecretMasker() *SecretMasker {
	return &SecretMasker{secrets: make([]string, 0)}
}

func (m *SecretMasker) AddSecret(value string) {
	if value == "" || len(value) < 3 {
		return
	}

	m.secrets = append(m.secrets, value)
}

func (m *SecretMasker) AddFromEnv(env map[string]string) {
	for key, value := range env {
		keyLower := strings.ToLower(key)
		if strings.Contains(keyLower, "pass") || strings.Contains(keyLower, "token") || strings.Contains(keyLower, "secret") || strings.Contains(keyLower, "key") {
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
