package fetcher

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestDefaultFetcher_Do(t *testing.T) {
	tests := []struct {
		name          string
		handler       http.HandlerFunc
		reqURL        string // overrides server URL if set
		ctxTimeout    time.Duration
		expectedCode  int
		expectedBody  []byte
		expectedError string
	}{
		{
			name: "Success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Hello, World!"))
			},
			expectedCode: http.StatusOK,
			expectedBody: []byte("Hello, World!"),
		},
		{
			name: "Server Error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal Error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: []byte("Internal Error"),
		},
		{
			name: "Timeout",
			handler: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(100 * time.Millisecond)
				w.WriteHeader(http.StatusOK)
			},
			ctxTimeout:    50 * time.Millisecond,
			expectedError: "context deadline exceeded",
		},
		{
			name:          "Invalid URL",
			reqURL:        "://invalid-url\n", // Control characters or invalid scheme
			expectedError: "parse",            // Expecting a URL parsing error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test server if a handler is provided
			var ts *httptest.Server
			var targetURL string

			if tt.handler != nil {
				ts = httptest.NewServer(tt.handler)
				defer ts.Close()
				targetURL = ts.URL
			}

			// Override URL if specified (for invalid URL tests)
			if tt.reqURL != "" {
				targetURL = tt.reqURL
			}

			// Create Fetcher
			f := NewDefaultFetcher()

			// Setup Context
			ctx := context.Background()
			if tt.ctxTimeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, tt.ctxTimeout)
				defer cancel()
			}

			// Execute
			req := &Request{
				URL: targetURL,
			}
			resp, err := f.Do(ctx, req)

			// Assertions
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.expectedError)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				if resp != nil {
					assert.Equal(t, tt.expectedCode, resp.StatusCode)
					assert.Equal(t, tt.expectedBody, resp.Body)
				}
			}
		})
	}
}
