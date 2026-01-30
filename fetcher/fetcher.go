package fetcher

import "context"

// Request represents a web request to be fetched.
type Request struct {
	URL     string
	Method  string
	Headers map[string]string
	Body    []byte
}

// Response represents the result of a fetch operation.
type Response struct {
	StatusCode int
	Body       []byte
	Headers    map[string][]string
	URL        string
}

// Fetcher defines the interface for fetching web resources.
type Fetcher interface {
	Do(ctx context.Context, req *Request) (*Response, error)	
}
