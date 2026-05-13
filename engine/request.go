package engine

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
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
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(r.Method, r.Url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: status code %d\n", resp.StatusCode)
		return nil, errors.New("Server responded with status code: " + string(rune(resp.StatusCode)))
	}

	bodyBytes, err := io.ReadAll(resp.Body)

	defer resp.Body.Close()

	return bodyBytes, nil
}

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
