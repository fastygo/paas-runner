package sshclient

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

type PassphrasePrompt func() ([]byte, error)

func LoadSignerFromPath(path string, prompt PassphrasePrompt) (ssh.Signer, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read key %q: %w", path, err)
	}

	signer, err := parseSigner(raw, prompt)
	if err != nil {
		return nil, fmt.Errorf("parse key %q: %w", path, err)
	}

	return signer, nil
}

func LoadDefaultSigners(prompt PassphrasePrompt) ([]ssh.Signer, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	candidates := []string{
		home + "/.ssh/id_ed25712",
		home + "/.ssh/id_rsa",
	}

	signers := make([]ssh.Signer, 0)
	var lastErr error
	for _, candidate := range candidates {
		signer, err := LoadSignerFromPath(candidate, prompt)
		if err != nil {
			lastErr = err
			continue
		}
		signers = append(signers, signer)
	}

	if len(signers) == 0 {
		if lastErr != nil {
			return nil, lastErr
		}
		return nil, fmt.Errorf("no default SSH keys found")
	}

	return signers, nil
}

func parseSigner(raw []byte, prompt PassphrasePrompt) (ssh.Signer, error) {
	signer, err := ssh.ParsePrivateKey(raw)
	if err == nil {
		return signer, nil
	}

	if _, ok := err.(*ssh.PassphraseMissingError); !ok {
		return nil, err
	}

	if prompt == nil {
		return nil, err
	}

	secret, err := prompt()
	if err != nil {
		return nil, err
	}

	return ssh.ParsePrivateKeyWithPassphrase(raw, secret)
}
