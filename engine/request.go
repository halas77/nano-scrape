package engine

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync/atomic"
	"time"
)

type Request struct {
	Url         string
	Method      string
	Header      http.Header
	ContentType string
}

func InitRequest(url string, method string, header http.Header) Request {
	req := Request{Url: url, Method: method, Header: header}
	return req
}

func (r Request) Execute(body ...map[string][]string) ([]byte, error) {

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Timeout: 10 * time.Second,
		Jar:     jar,
	}

	req, err := http.NewRequest(r.Method, r.Url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if isBotBlock(resp) {
		fmt.Println("-------- There is a bot ---------------")
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: status code %d\n", resp.StatusCode)
		return nil, errors.New("Server responded with status code: " + string(rune(resp.StatusCode)))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return bodyBytes, nil
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

/*
func Browser() {
	// Create context and allocate a browser instance
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set a timeout
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var htmlContent string

	// Navigate and wait for the page to evaluate JavaScript/solve challenges
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://target-website.com`),
		// Wait for a specific element that proves you passed the challenge
		chromedp.WaitVisible(`#main-content-loaded`, chromedp.ByID),
		// Grab the rendered HTML
		chromedp.OuterHTML(`html`, &htmlContent),
	)
	if err != nil {
		log.Fatalf("Failed to bypass via Chromedp: %v", err)
	}

	fmt.Println("Successfully bypassed and fetched page!")
}

*/

// ProxyRotator holds a slice of proxy URLs and rotates through them sequentially
type ProxyRotator struct {
	proxies []string
	index   atomic.Uint64
}

func NewProxyRotator(proxies []string) *ProxyRotator {
	return &ProxyRotator{proxies: proxies}
}

// GetProxyFunc returns a function compatible with http.Transport.Proxy
func (pr *ProxyRotator) GetProxyFunc() func(*http.Request) (*url.URL, error) {
	return func(req *http.Request) (*url.URL, error) {
		if len(pr.proxies) == 0 {
			return nil, nil // Fallback to direct connection if pool is empty
		}
		// Thread-safe increment and modulo to get the next proxy
		idx := pr.index.Add(1) % uint64(len(pr.proxies))
		return url.Parse(pr.proxies[idx])
	}
}

func main() {
	proxyList := []string{
		"http://username:password@proxy1.example.com:8080",
		"http://username:password@proxy2.example.com:8080",
		"http://username:password@proxy3.example.com:8080",
	}

	rotator := NewProxyRotator(proxyList)

	// Best Practice: Always set explicit timeouts on your client when scraping
	transport := &http.Transport{
		Proxy:           rotator.GetProxyFunc(),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Optional: if dealing with messy SSL certificates
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   15 * time.Millisecond, // Crucial: don't let dead proxies hang your scraper
	}

	// Example request
	resp, err := client.Get("https://httpbin.org/ip")
	if err != nil {
		fmt.Printf("Error fetching URL: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Request successful!")
}
