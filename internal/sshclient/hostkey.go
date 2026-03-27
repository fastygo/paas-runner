package sshclient

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func HostKeyCallback(mode string) (ssh.HostKeyCallback, error) {
	mode = strings.ToLower(mode)
	switch mode {
	case "", "strict":
		path, err := knownHostsPath()
		if err != nil {
			return nil, err
		}
		return knownhosts.New(path)
	case "tofu":
		return tofuCallback()
	case "insecure":
		fmt.Fprintln(os.Stderr, "WARNING: SSH host key verification is disabled")
		return ssh.InsecureIgnoreHostKey(), nil
	default:
		return nil, fmt.Errorf("invalid host_key_check %q", mode)
	}
}

func knownHostsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".ssh", "known_hosts"), nil
}

func tofuCallback() (ssh.HostKeyCallback, error) {
	path, err := knownHostsPath()
	if err != nil {
		return nil, err
	}

	strict, err := knownhosts.New(path)
	if err != nil {
		return nil, err
	}

	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		if err := strict(hostname, remote, key); err == nil {
			return nil
		}

		keyErr, ok := err.(*knownhosts.KeyError)
		if !ok {
			return err
		}

		// Unknown host key: accept and persist.
		if len(keyErr.Want) == 0 {
			return appendKnownHost(path, hostname, remote, key)
		}

		return err
	}, nil
}

func appendKnownHost(path string, hostname string, remote net.Addr, key ssh.PublicKey) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}

	host := hostname
	if host == "" && remote != nil {
		host = remote.String()
	}

	if idx := strings.Index(host, ":"); idx > 0 {
		host = host[:idx]
	}

	line := fmt.Sprintf("%s %s", host, strings.TrimSpace(string(ssh.MarshalAuthorizedKey(key))))

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(line + "\n")
	return err
}
