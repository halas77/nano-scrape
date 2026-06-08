package engine

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
	Header http.Header
	client *http.Client
}

// NewClient initializes a new HTTP engine with default configurations.
func NewClient(baseHeader ...http.Header) *Client {
	jar, _ := cookiejar.New(nil)

	var hdr http.Header
	if len(baseHeader) > 0 && baseHeader[0] != nil {
		hdr = baseHeader[0].Clone()
	} else {
		hdr = make(http.Header)
	}

	return &Client{
		Header: hdr,
		client: &http.Client{
			Timeout: 10 * time.Second,
			Jar:     jar,
		},
	}
}

// ProxyRotator assigns a custom transport to handle proxy rotations.
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

func (c *Client) Execute(method, targetURL string, body ...io.Reader) ([]byte, error) {
	if c == nil {
		return nil, fmt.Errorf("client instance is nil")
	}

	var reqBody io.Reader
	if len(body) > 0 {
		reqBody = body[0]
	}

	req, err := http.NewRequest(method, targetURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header = c.Header.Clone()

	setDefaultHeader(req.Header, "User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	setDefaultHeader(req.Header, "Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	setDefaultHeader(req.Header, "Accept-Language", "en-US,en;q=0.5")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request execution failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("server returned bad status: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return bodyBytes, nil
}

// SendJSON automates marshaling structures into JSON payloads.
// Accepts HTTP methods like http.MethodPost, http.MethodPut, or http.MethodPatch.
func (c *Client) SendJSON(method, targetURL string, payload any) ([]byte, error) {
	if c == nil {
		return nil, fmt.Errorf("client instance is nil")
	}

	jsonValue, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json: %w", err)
	}

	reqBody := bytes.NewReader(jsonValue)
	return c.executeWithModifier(targetURL, method, reqBody, func(h http.Header) {
		h.Set("Content-Type", "application/json")
	})
}

// SendForm encodes payload maps into standard x-www-form-urlencoded payloads.
// Accepts HTTP methods like http.MethodPost, http.MethodPut, or http.MethodPatch.
func (c *Client) SendForm(method, targetURL string, payload map[string]string) ([]byte, error) {
	if c == nil {
		return nil, fmt.Errorf("client instance is nil")
	}

	formData := url.Values{}
	for key, val := range payload {
		formData.Set(key, val)
	}

	reqBody := strings.NewReader(formData.Encode())
	return c.executeWithModifier(targetURL, method, reqBody, func(h http.Header) {
		h.Set("Content-Type", "application/x-www-form-urlencoded")
	})
}

// CookiesFor retrieves stored cookies for a specific target URL.
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

// IsBotBlock inspects response structures without draining the stream permanently.
func IsBotBlock(resp *http.Response) bool {
	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
		return true
	}

	if resp.Header.Get("Server") == "cloudflare" &&
		(resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusForbidden) {
		return true
	}

	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	resp.Body = io.NopCloser(io.MultiReader(bytes.NewReader(bodyBytes), resp.Body))

	lowered := bytes.ToLower(bodyBytes)
	return bytes.Contains(lowered, []byte("captcha")) ||
		bytes.Contains(lowered, []byte("javascript challenge")) ||
		bytes.Contains(lowered, []byte("checking your browser"))
}

func setDefaultHeader(h http.Header, key, value string) {
	if h.Get(key) == "" {
		h.Set(key, value)
	}
}

func (c *Client) executeWithModifier(targetURL, method string, body io.Reader, modifyHeaders func(http.Header)) ([]byte, error) {
	oldHeaders := c.Header.Clone()
	modifyHeaders(c.Header)
	defer func() { c.Header = oldHeaders }()

	return c.Execute(method, targetURL, body)
}
