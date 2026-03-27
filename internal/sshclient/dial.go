package sshclient

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/paas/paas-runner/internal/config"
	"golang.org/x/crypto/ssh"
)

func Dial(server config.ServerConfig, prompt PassphrasePrompt) (*ssh.Client, error) {
	server = server.ApplyDefaults()

	if server.Host == "" {
		return nil, fmt.Errorf("server host is required")
	}

	callback, err := HostKeyCallback(server.HostKeyCheck)
	if err != nil {
		return nil, err
	}

	auth, cleanup, err := resolveAuthMethods(server, prompt)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup.Close()
	}

	port := server.Port
	if port == 0 {
		port = 22
	}

	address := net.JoinHostPort(server.Host, fmt.Sprintf("%d", port))
	clientCfg := &ssh.ClientConfig{
		User:            server.User,
		Auth:            auth,
		HostKeyCallback: callback,
		Timeout:         30 * time.Second,
	}

	return ssh.Dial("tcp", address, clientCfg)
}

func resolveAuthMethods(server config.ServerConfig, prompt PassphrasePrompt) ([]ssh.AuthMethod, io.Closer, error) {
	methods := make([]ssh.AuthMethod, 0)

	if server.Key != "" {
		signer, err := LoadSignerFromPath(server.Key, prompt)
		if err != nil {
			return nil, nil, err
		}
		methods = append(methods, ssh.PublicKeys(signer))
		return methods, nil, nil
	}

	if signers, closer, err := LoadAgentSigners(); err == nil && len(signers) > 0 {
		methods = append(methods, ssh.PublicKeys(signers...))
		return methods, closer, nil
	}

	if signers, err := LoadDefaultSigners(prompt); err == nil && len(signers) > 0 {
		methods = append(methods, ssh.PublicKeys(signers...))
	}

	if len(methods) == 0 {
		return nil, nil, fmt.Errorf("no SSH auth methods discovered")
	}

	return methods, nil, nil
}
