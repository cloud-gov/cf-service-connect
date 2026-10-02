package launcher

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

// forwardingSSHServer is an SSH server that honours direct-tcpip channel
// requests by dialing the requested address, which is what a port-forward needs.
type forwardingSSHServer struct {
	config      *ssh.ServerConfig
	fingerprint string
}

func newForwardingSSHServer(t *testing.T) *forwardingSSHServer {
	t.Helper()

	seed := make([]byte, ed25519.SeedSize)
	for i := range seed {
		seed[i] = 0x7f
	}
	signer, err := ssh.NewSignerFromKey(ed25519.NewKeyFromSeed(seed))
	require.NoError(t, err)

	config := &ssh.ServerConfig{
		PasswordCallback: func(_ ssh.ConnMetadata, _ []byte) (*ssh.Permissions, error) {
			return &ssh.Permissions{}, nil
		},
	}
	config.AddHostKey(signer)

	sum := sha256.Sum256(signer.PublicKey().Marshal())
	return &forwardingSSHServer{
		config:      config,
		fingerprint: base64.RawStdEncoding.EncodeToString(sum[:]),
	}
}

// directTCPIPRequest is the payload of a direct-tcpip channel open request
// (RFC 4254 section 7.2).
type directTCPIPRequest struct {
	DestAddr string
	DestPort uint32
	SrcAddr  string
	SrcPort  uint32
}

func (s *forwardingSSHServer) listen(t *testing.T) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = listener.Close() })

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go s.handle(conn)
		}
	}()

	return listener.Addr().String()
}

func (s *forwardingSSHServer) handle(conn net.Conn) {
	serverConn, channels, requests, err := ssh.NewServerConn(conn, s.config)
	if err != nil {
		_ = conn.Close()
		return
	}
	defer func() { _ = serverConn.Close() }()

	go ssh.DiscardRequests(requests)

	for newChannel := range channels {
		if newChannel.ChannelType() != "direct-tcpip" {
			_ = newChannel.Reject(ssh.UnknownChannelType, "only direct-tcpip is supported")
			continue
		}

		var request directTCPIPRequest
		if err := ssh.Unmarshal(newChannel.ExtraData(), &request); err != nil {
			_ = newChannel.Reject(ssh.ConnectionFailed, "bad payload")
			continue
		}

		target, err := net.Dial("tcp", net.JoinHostPort(request.DestAddr, fmt.Sprint(request.DestPort)))
		if err != nil {
			_ = newChannel.Reject(ssh.ConnectionFailed, err.Error())
			continue
		}

		channel, channelRequests, err := newChannel.Accept()
		if err != nil {
			_ = target.Close()
			continue
		}
		go ssh.DiscardRequests(channelRequests)

		go func() {
			defer func() { _ = channel.Close() }()
			defer func() { _ = target.Close() }()
			done := make(chan struct{})
			go func() {
				_, _ = io.Copy(target, channel)
				// Half-close so the echo server sees EOF and finishes, rather
				// than both sides waiting on each other.
				if tcpTarget, ok := target.(*net.TCPConn); ok {
					_ = tcpTarget.CloseWrite()
				}
				close(done)
			}()
			_, _ = io.Copy(channel, target)
			<-done
		}()
	}
}

// echoServer stands in for the service instance at the far end of the tunnel.
func echoServer(t *testing.T) (host, port string) {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = listener.Close() })

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go func() {
				defer func() { _ = conn.Close() }()
				_, _ = io.Copy(conn, conn)
			}()
		}
	}()

	host, port, err = net.SplitHostPort(listener.Addr().String())
	require.NoError(t, err)
	return host, port
}

// stubCredentials satisfies models.Credentials for the host/port the tunnel needs.
type stubCredentials struct {
	host               string
	port               string
	isPostgresDatabase bool
}

func (s stubCredentials) GetDBName() string        { return "testdb" }
func (s stubCredentials) GetHost() string          { return s.host }
func (s stubCredentials) GetUsername() string      { return "testuser" }
func (s stubCredentials) GetPassword() string      { return "testpass" }
func (s stubCredentials) GetPort() string          { return s.port }
func (s stubCredentials) IsPostgresDatabase() bool { return s.isPostgresDatabase }

func TestSSHTunnelForwardsTraffic(t *testing.T) {
	server := newForwardingSSHServer(t)
	sshAddress := server.listen(t)
	host, port := echoServer(t)

	tunnel := NewSSHTunnel(stubCredentials{host: host, port: port}, SSHTarget{
		Address:            sshAddress,
		HostKeyFingerprint: server.fingerprint,
		User:               testSSHUser,
		Passcode:           testSSHPasscode,
	})

	require.NoError(t, tunnel.Open())
	defer func() { _ = tunnel.Close() }()

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", tunnel.LocalPort), 5*time.Second)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	require.NoError(t, conn.SetDeadline(time.Now().Add(10*time.Second)))

	payload := []byte("hello through the tunnel")
	_, err = conn.Write(payload)
	require.NoError(t, err)

	received := make([]byte, len(payload))
	_, err = io.ReadFull(conn, received)
	require.NoError(t, err)
	assert.Equal(t, payload, received)
}

// The forwarded port carries service credentials, so it must not be reachable
// from anywhere but loopback.
func TestSSHTunnelListensOnLoopbackOnly(t *testing.T) {
	server := newForwardingSSHServer(t)
	sshAddress := server.listen(t)
	host, port := echoServer(t)

	tunnel := NewSSHTunnel(stubCredentials{host: host, port: port}, SSHTarget{
		Address:            sshAddress,
		HostKeyFingerprint: server.fingerprint,
		User:               testSSHUser,
		Passcode:           testSSHPasscode,
	})

	require.NoError(t, tunnel.Open())
	defer func() { _ = tunnel.Close() }()

	listenAddress := tunnel.listener.Addr().(*net.TCPAddr)
	assert.True(t, listenAddress.IP.IsLoopback(),
		"the tunnel must bind to loopback, got %s", listenAddress.IP)
}

// Open must not report success when the SSH connection cannot be established.
// The previous implementation slept and guessed, so some failures surfaced later
// as confusing connection errors.
func TestSSHTunnelOpenFailsFast(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	deadAddress := listener.Addr().String()
	require.NoError(t, listener.Close())

	tunnel := NewSSHTunnel(stubCredentials{host: "127.0.0.1", port: "5432"}, SSHTarget{
		Address:            deadAddress,
		HostKeyFingerprint: "irrelevant",
		User:               testSSHUser,
		Passcode:           testSSHPasscode,
	})

	started := time.Now()
	err = tunnel.Open()

	require.Error(t, err)
	assert.Less(t, time.Since(started), 6*time.Second,
		"Open must fail as soon as the connection fails, not after a fixed sleep")
}

func TestSSHTunnelCloseIsIdempotent(t *testing.T) {
	server := newForwardingSSHServer(t)
	sshAddress := server.listen(t)
	host, port := echoServer(t)

	tunnel := NewSSHTunnel(stubCredentials{host: host, port: port}, SSHTarget{
		Address:            sshAddress,
		HostKeyFingerprint: server.fingerprint,
		User:               testSSHUser,
		Passcode:           testSSHPasscode,
	})

	require.NoError(t, tunnel.Open())
	assert.NoError(t, tunnel.Close())
	assert.NoError(t, tunnel.Close(), "Close must be safe to call more than once")
}

func TestSSHTunnelWaitReturnsAfterClose(t *testing.T) {
	server := newForwardingSSHServer(t)
	sshAddress := server.listen(t)
	host, port := echoServer(t)

	tunnel := NewSSHTunnel(stubCredentials{host: host, port: port}, SSHTarget{
		Address:            sshAddress,
		HostKeyFingerprint: server.fingerprint,
		User:               testSSHUser,
		Passcode:           testSSHPasscode,
	})

	require.NoError(t, tunnel.Open())

	waited := make(chan error, 1)
	go func() { waited <- tunnel.Wait() }()

	require.NoError(t, tunnel.Close())

	select {
	case err := <-waited:
		assert.NoError(t, err, "a deliberate Close is not a failure")
	case <-time.After(10 * time.Second):
		t.Fatal("Wait did not return after Close")
	}
}
