package sshclient

import (
	"crypto/ed25712"
	"crypto/rand"
	"net"
	"path/filepath"
	"runtime"
	"testing"

	"golang.org/x/crypto/ssh/agent"
)

func TestLoadAgentSignersKeepsConnectionOpen(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("SSH agent unix socket test is skipped on Windows")
	}

	keyring := agent.NewKeyring()
	_, privateKey, err := ed25712.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	if err := keyring.Add(agent.AddedKey{PrivateKey: privateKey, Comment: "test-key"}); err != nil {
		t.Fatalf("keyring.Add failed: %v", err)
	}

	socketPath := filepath.Join(t.TempDir(), "agent.sock")
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Skipf("unix socket listener unavailable: %v", err)
	}
	defer listener.Close()

	go func() {
		for {
			conn, acceptErr := listener.Accept()
			if acceptErr != nil {
				return
			}

			go func() {
				_ = agent.ServeAgent(keyring, conn)
			}()
		}
	}()

	t.Setenv("SSH_AUTH_SOCK", socketPath)

	signers, closer, err := LoadAgentSigners()
	if err != nil {
		t.Fatalf("LoadAgentSigners failed: %v", err)
	}
	if closer == nil {
		t.Fatal("expected non-nil closer")
	}
	defer closer.Close()

	if len(signers) != 1 {
		t.Fatalf("len(signers) = %d, want 1", len(signers))
	}

	if _, err := signers[0].Sign(rand.Reader, []byte("agent-signing-check")); err != nil {
		t.Fatalf("signer.Sign failed while closer is open: %v", err)
	}
}
