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

type Request struct {
	Header http.Header
	client *http.Client
}

func InitRequest(header ...http.Header) *Request {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Timeout: 10 * time.Second,
		Jar:     jar,
	}

	var hdr http.Header
	if len(header) > 0 {
		hdr = header[0]
	}
	req := Request{Header: hdr, client: client}
	return &req
}

func (r *Request) ProxyRotator(proxies ...string) {
	rotator := NewProxyRotator(proxies)
	r.client.Transport = &http.Transport{
		Proxy:               rotator.GetProxyFunc(),
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}
}

func (r *Request) Execute(url string, method string, body ...io.Reader) ([]byte, error) {
	client := r.client

	var reqBody io.Reader
	if len(body) > 0 && body[0] != nil {
		reqBody = body[0]
	}
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Set(key, value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// if isBotBlock(resp) {
	// 	fmt.Println("-------- There is a bot ---------------")
	// }

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: status code %d\n", resp.StatusCode)
		return nil, fmt.Errorf("server responded with status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}

func (r *Request) MakeJSONPostRequest(url string, method string, payload map[string]string) ([]byte, error) {
	jsonValue, _ := json.Marshal(payload)
	if r.Header == nil {
		r.Header = make(http.Header)
	}
	r.Header.Set("Content-Type", "application/json")

	reqBody := bytes.NewReader(jsonValue)
	resp, err := r.Execute(url, method, reqBody)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *Request) MakeFormPostRequest(URL string, method string, payload map[string]string) ([]byte, error) {

	formData := url.Values{}
	for key, val := range payload {
		formData.Set(key, val)
	}

	if r.Header == nil {
		r.Header = make(http.Header)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	reqBody := strings.NewReader(formData.Encode())
	resp, err := r.Execute(URL, method, reqBody)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func isBotBlock(resp *http.Response) bool {
	// Check status codes
	if resp.StatusCode == 403 || resp.StatusCode == 429 {
		return true
	}

	// Check headers (e.g., Cloudflare)
	if resp.Header.Get("Server") == "cloudflare" {
		// Cloudflare uses 503 or 403 for challenges
		if resp.StatusCode == 503 || resp.StatusCode == 403 {
			return true
		}
	}

	// Check body content (Read a snippet of the body)
	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096)) // Read first 4KB safely

	if bytes.Contains(bodyBytes, []byte("captcha")) ||
		bytes.Contains(bodyBytes, []byte("javascript challenge")) ||
		bytes.Contains(bodyBytes, []byte("checking your browser")) {
		return true
	}

	return false
}

// Cookies returns cookies currently stored in the jar for Request.Url.
// func (r *Request) Cookies() []*http.Cookie {
// 	if r == nil || r.client == nil || r.client.Jar == nil {
// 		return nil
// 	}

// 	u, err := url.Parse(r.Url)
// 	if err != nil {
// 		return nil
// 	}

// 	return r.client.Jar.Cookies(u)
// }

// CookiesFor returns cookies currently stored in the jar for a specific URL.
func (r *Request) CookiesFor(rawURL string) []*http.Cookie {
	if r == nil || r.client == nil || r.client.Jar == nil || rawURL == "" {
		return nil
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil
	}

	return r.client.Jar.Cookies(u)
}
