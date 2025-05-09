package tunnel

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/dnc-data-mcp/config"
	"golang.org/x/crypto/ssh"
)

type SSHTunnel struct {
	Local  *net.TCPListener
	Config *config.Config
	client *ssh.Client
}

func NewSSHTunnel(cfg *config.Config) (*SSHTunnel, error) {
	// Expand the private key path if it contains ~
	keyPath := cfg.Default.SSHPrivateKey
	if strings.HasPrefix(keyPath, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("error getting home directory: %v", err)
		}
		keyPath = filepath.Join(homeDir, keyPath[1:])
	}

	log.Printf("Using SSH key: %s\n", keyPath)

	// Read the SSH key
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	// Configure SSH client
	sshConfig := &ssh.ClientConfig{
		User: cfg.Default.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to SSH server
	addr := fmt.Sprintf("%s:%d", cfg.Default.SSHHost, cfg.Default.SSHPort)
	log.Printf("Connecting to SSH server: %s\n", addr)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to SSH server: %v", err)
	}
	log.Printf("Connected to SSH server: %s\n", addr)

	// Start local listener on port 5433
	local, err := net.Listen("tcp", "localhost:5433")
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("unable to start local listener: %v", err)
	}
	log.Printf("Started local listener on: %s\n", local.Addr().String())

	tunnel := &SSHTunnel{
		Local:  local.(*net.TCPListener),
		Config: cfg,
		client: client,
	}

	// Start forwarding
	go tunnel.forward()

	return tunnel, nil
}

func (t *SSHTunnel) forward() {
	// Use the database server from config
	dbAddr := fmt.Sprintf("%s:%d",
		t.Config.Database.ROTraffic.Server, // Use the server from config
		t.Config.Database.ROTraffic.Port)   // Use the port from config
	log.Printf("Starting tunnel forwarding from %s to %s\n",
		t.Local.Addr().String(),
		dbAddr)

	for {
		local, err := t.Local.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				log.Printf("Local listener closed, stopping tunnel\n")
				return
			}
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go func() {
			defer local.Close()
			log.Printf("Attempting to connect to remote address: %s\n", dbAddr)

			remote, err := t.client.Dial("tcp", dbAddr)
			if err != nil {
				log.Printf("Remote dial error: %s\n", err)
				return
			}
			defer remote.Close()

			log.Printf("Connected to remote address: %s\n", dbAddr)

			// Create channels to handle connection closure
			done := make(chan struct{})
			go func() {
				defer close(done)
				n, err := copyConn(local, remote)
				if err != nil {
					log.Printf("Error copying local to remote: %v (copied %d bytes)\n", err, n)
				}
			}()

			go func() {
				defer close(done)
				n, err := copyConn(remote, local)
				if err != nil {
					log.Printf("Error copying remote to local: %v (copied %d bytes)\n", err, n)
				}
			}()

			// Wait for either direction to complete
			<-done
		}()
	}
}

// GetLocalEndpoint returns the local endpoint for the tunnel
func (t *SSHTunnel) GetLocalEndpoint() string {
	return t.Local.Addr().String()
}

func copyConn(writer, reader net.Conn) (int64, error) {
	return io.Copy(writer, reader)
}

func (t *SSHTunnel) Close() error {
	if t.Local != nil {
		t.Local.Close()
	}
	if t.client != nil {
		t.client.Close()
	}
	return nil
}
