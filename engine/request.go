package engine

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
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
