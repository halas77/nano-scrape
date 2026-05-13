package engine

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func Fetch() {
	// 1. Create a client with a timeout (Best practice)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 2. Create the request
	req, err := http.NewRequest("GET", "http://127.0.0.1:5500/examples/basic/index.html", nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	// 3. Execute the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	// Always close the body when finished
	defer resp.Body.Close()

	// 4. Accessing Headers
	fmt.Println("--- Individual Headers ---")

	// Get a specific header (Case-insensitive)
	contentType := resp.Header.Get("Content-Type")
	server := resp.Header.Get("Server")

	if strings.Contains(contentType, "text/html") {
		// Process as HTML string
		fmt.Println(resp.Body)
		html.Parse(resp.Body)
	}

	fmt.Printf("Content-Type: %s\n", contentType)
	fmt.Printf("Server: %s\n", server)

	fmt.Println("\n--- All Response Headers ---")

	// Loop through all headers
	for key, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}
}
