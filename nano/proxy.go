package nano

import (
	"net/http"
	"net/url"
	"sync/atomic"
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
