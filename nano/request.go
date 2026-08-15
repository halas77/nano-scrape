package nano

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	BaseHeader http.Header
	client     *http.Client
}

func NewClient(baseHeader ...http.Header) *Client {
	jar, _ := cookiejar.New(nil)

	var hdr http.Header
	if len(baseHeader) > 0 && baseHeader[0] != nil {
		hdr = baseHeader[0].Clone()
	} else {
		hdr = make(http.Header)
	}

	return &Client{
		BaseHeader: hdr,
		client: &http.Client{
			Timeout: 10 * time.Second,
			Jar:     jar,
		},
	}
}

func (c *Client) ProxyRotator(proxies ...string) error {
	if c == nil {
		return fmt.Errorf("client instance is nil")
	}

	rotator := NewProxyRotator(proxies)
	c.client.Transport = &http.Transport{
		Proxy:               rotator.GetProxyFunc(),
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	return nil
}

func (c *Client) Execute(method, targetURL string, body io.Reader, extraHeaders ...map[string]string) (io.Reader, error) {
	req, err := http.NewRequest(method, targetURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header = c.BaseHeader.Clone()

	if len(extraHeaders) > 0 && extraHeaders[0] != nil {
		for k, v := range extraHeaders[0] {
			req.Header.Set(k, v)
		}
	}

	setDefaultHeader(req.Header, "User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0 Safari/537.36")
	setDefaultHeader(req.Header, "Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	setDefaultHeader(req.Header, "Accept-Language", "en-US,en;q=0.5")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request execution failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("server returned bad status: %d", resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		return nil, err
	}

	return buf, nil
}

func (c *Client) Get(targetURL string) (io.Reader, error) {
	return c.Execute(http.MethodGet, targetURL, nil)
}

func (c *Client) SendJSON(method, targetURL string, payload any) (io.Reader, error) {
	jsonValue, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json: %w", err)
	}

	return c.Execute(method, targetURL, bytes.NewReader(jsonValue), map[string]string{
		"Content-Type": "application/json",
	})
}

func (c *Client) SendForm(method, targetURL string, payload map[string]string) (io.Reader, error) {
	formData := url.Values{}
	for key, val := range payload {
		formData.Set(key, val)
	}

	return c.Execute(method, targetURL, strings.NewReader(formData.Encode()), map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	})
}

func setDefaultHeader(h http.Header, key, value string) {
	if h.Get(key) == "" {
		h.Set(key, value)
	}
}

func (c *Client) CookiesFor(rawURL string) []*http.Cookie {
	if c == nil || c.client == nil || c.client.Jar == nil || rawURL == "" {
		return nil
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil
	}
	return c.client.Jar.Cookies(u)
}
