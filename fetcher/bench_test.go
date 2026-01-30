package fetcher

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkDefaultFetcher_Do(b *testing.B) {
	// Setup a lightweight test server
	// We use a simple handler to minimize server-side processing time
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("benchmark payload"))
	}))
	defer ts.Close()

	// Initialize the fetcher
	f := NewDefaultFetcher()
	
	// Create a reusable request object (struct copy is cheap, strictly we should pass a new pointer or ensure safety, 
	// but Do takes *Request and treats it essentially as read-only for input fields except strictly Body if it was a stream, 
	// but here we have no body or simple body).
	// effectively Do(ctx, req) uses req properties.
	
	reqBase := &Request{
		URL: ts.URL,
	}

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		// Create a context for each operation or reuse background? 
		// Context creation has some overhead, but is part of the cost of using the library.
		// We'll reuse Background for stability or creating new value ctx if we want to be realistic.
		// The requirement asks to benchmark the overhead of the Do method.
		
		ctx := context.Background()
		
		for pb.Next() {
			_, err := f.Do(ctx, reqBase)
			if err != nil {
				b.Errorf("request failed: %v", err)
			}
		}
	})
}
