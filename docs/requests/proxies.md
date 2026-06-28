# Proxies

## ProxyRotator

`ProxyRotator` is a thread-safe utility used to balance outgoing web requests evenly across a pool of proxy servers. It ensures that concurrent requests do not conflict when choosing the next proxy.

## Functions

### `NewProxyRotator`

```go
func NewProxyRotator(proxies []string) *ProxyRotator
```

Creates and initializes a new proxy rotation manager with a defined list of proxy addresses.

- **What it takes:** A slice of proxy URL strings (e.g., `[]string{"http://proxy1.com:8080", "http://proxy2.com:8080"}`).
- **What it returns:** An initialized pointer to a `ProxyRotator`.

#### `💡 Example`

```go
proxyList := []string{
    "http://proxy-us.example.com:3128",
    "http://proxy-eu.example.com:3128",
}
rotator := engine.NewProxyRotator(proxyList)
```

### `GetProxyFunc`

```go
func (pr *ProxyRotator) GetProxyFunc() func(*http.Request) (*url.URL, error)
```

Generates a configuration function compatible with standard Go network transports (`http.Transport.Proxy`). Every time Go starts a new network request, this function automatically cycles to the next proxy in your list using a Round-Robin algorithm.

- **What it takes:** Nothing.
- **What it returns:** A proxy selection function that automatically cycles through your proxy pool safely across multiple goroutines.

### `ProxyRotator`

```go
func (c *Client) ProxyRotator(proxies ...string)
```

A convenience method attached directly to your custom `Client`. It creates a rotation pool and automatically attaches it to the client's internal network configuration, complete with optimal timeout settings.

- **What it takes:** One or more proxy URL strings passed as separate arguments.
- **What it returns:** Nothing. It configures your active client internally.

#### `💡 Example`

```go
client := engine.NewClient()

// Directly apply rotating proxies to your client instance
client.ProxyRotator(
    "http://123.45.67.89:8080",
    "http://98.76.54.321:8080",
)

// All future requests made by this client will cycle through those proxies
resp, err := client.Execute("GET", "https://example.com")
```

This standalone function demonstrates the entire pipeline: initializing a client, attaching rotating proxies, downloading a webpage safely using `io.Reader` streams, and finding specific elements within the parsed HTML.

```go
package main

import (
	"fmt"
	"log"
	"strings"
)

func ScrapeTargetWithProxies() {
	proxies := []string{
		"http://proxy-us.example.com:3128",
		"http://proxy-eu.example.com:3128",
	}
	target := "https://example.com"

	// 1. Initialize your custom HTTP client wrapper
	client := engine.NewClient()

	// 2. Attach your proxy list using the ProxyRotator method.
	// This automatically configures the internal Go http.Transport to cycle proxies safely.
	client.ProxyRotator(proxies...)

	fmt.Println("🚀 Request pipeline initialized with proxy rotation...")

	// 3. Execute a network request to pull down the webpage data stream.
	// This returns an io.Reader (specifically a *bytes.Buffer), safely closing the network socket internally.
	responseStream, err := client.Execute("GET", target)
	if err != nil {
		log.Fatalf("❌ Network request failed: %v", err)
	}

	// 4. Pass the streaming data directly into the HTML Parser.
	// InitDocument automatically reads from the io.Reader stream.
	rootTag, err := engine.InitDocument(responseStream)
	if err != nil {
		log.Fatalf("❌ Failed to parse HTML content: %v", err)
	}

	fmt.Printf("✅ Successfully fetched and parsed: %s (Root Tag: %s)\n", target, rootTag.Name)

	// 5. Use the Tag query functions to extract data from the document tree.
	// Find the very first <h1> element on the page
	firstHeading := rootTag.FindFirst("h1")
	if firstHeading != nil {
		fmt.Printf("🎯 First Heading found on page: %s\n", firstHeading.Name)
	}

	// Stream and print all links found on the page using a callback
	fmt.Println("🔗 Listing all links found on the page:")
	rootTag.Find("a", nil, func(linkTag *engine.Tag) {
		// You can access link attributes via linkTag.Attrs
		fmt.Printf("   Found link element: %s\n", linkTag.Name)
	})
}
```
