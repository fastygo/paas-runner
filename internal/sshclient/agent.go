package sshclient

import (
	"fmt"
	"io"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func LoadAgentSigners() ([]ssh.Signer, io.Closer, error) {
	socket := os.Getenv("SSH_AUTH_SOCK")
	if socket == "" {
		return nil, nil, fmt.Errorf("SSH_AUTH_SOCK is not set")
	}

	conn, err := net.Dial("unix", socket)
	if err != nil {
		return nil, nil, fmt.Errorf("connect to SSH agent: %w", err)
	}

	client := agent.NewClient(conn)
	signers, err := client.Signers()
	if err != nil {
		_ = conn.Close()
		return nil, nil, fmt.Errorf("load SSH agent signers: %w", err)
	}

	return signers, conn, nil
}
