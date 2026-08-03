package engine

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync/atomic"
	"testing"

	"github.com/things-go/go-socks5"
)

type ProxyRotator struct {
	proxies []string
	index   atomic.Uint64
}

func NewProxyRotator(proxies []string) *ProxyRotator {
	return &ProxyRotator{proxies: proxies}
}

func (pr *ProxyRotator) GetProxyFunc() func(*http.Request) (*url.URL, error) {
	return func(req *http.Request) (*url.URL, error) {
		if len(pr.proxies) == 0 {
			return nil, nil
		}

		idx := pr.index.Add(1) % uint64(len(pr.proxies))
		return url.Parse(pr.proxies[idx])
	}
}

// Helper to launch a local SOCKS5 server returning its socks5:// address
func startMockSocks5Server(t *testing.T, hitCount *int32) (string, func()) {
	t.Helper()

	server := socks5.NewServer()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen for socks5: %v", err)
	}

	go func() {
		// Increment counter when a request connects through this proxy
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			atomic.AddInt32(hitCount, 1)
			go server.ServeConn(conn)
		}
	}()

	addr := fmt.Sprintf("socks5://%s", listener.Addr().String())
	cleanup := func() {
		listener.Close()
	}

	return addr, cleanup
}
