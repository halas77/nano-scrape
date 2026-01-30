package fetcher

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DefaultFetcher is the default implementation of the Fetcher interface
// using the standard net/http package.
type DefaultFetcher struct {
	client *http.Client
}

// NewDefaultFetcher creates a new DefaultFetcher with a production-grade http.Client.
func NewDefaultFetcher() *DefaultFetcher {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return &DefaultFetcher{
		client: &http.Client{
			Transport: transport,
			// We do not set a global timeout here to allow the context to control cancellation.
		},
	}
}

// Do executes an HTTP request and returns the response.
func (f *DefaultFetcher) Do(ctx context.Context, req *Request) (*Response, error) {
	method := req.Method
	if method == "" {
		method = http.MethodGet
	}

	var bodyReader io.Reader
	if len(req.Body) > 0 {
		bodyReader = bytes.NewReader(req.Body)
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, req.URL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("fetcher: failed to create request: %w", err)
	}

	// Set headers if provided
	if req.Headers != nil {
		for k, v := range req.Headers {
			httpReq.Header.Set(k, v)
		}
	}

	// Politeness: Set default User-Agent if not provided
	if httpReq.Header.Get("User-Agent") == "" {
		httpReq.Header.Set("User-Agent", "Goscrape/1.0")
	}

	resp, err := f.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("fetcher: failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fetcher: failed to read response body: %w", err)
	}

	finalURL := req.URL
	if resp.Request != nil && resp.Request.URL != nil {
		finalURL = resp.Request.URL.String()
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       responseBody,
		Headers:    resp.Header,
		URL:        finalURL,
	}, nil
}
