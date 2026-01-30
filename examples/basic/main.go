package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/halas77/goscrape/fetcher"
)

// ANSI Color Codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
)

func main() {
	printHeader()

	// 1. Initialize the fetcher
	f := fetcher.NewDefaultFetcher()

	// 2. define the target URL
	url := "https://example.com"
	
	printSection("REQUEST DETAILS")
	fmt.Printf("%sTarget URL:%s %s\n", Bold, Reset, url)
	fmt.Printf("%sMethod:%s     GET\n", Bold, Reset)
	fmt.Printf("%sTimeout:%s    10s\n", Bold, Reset)
	fmt.Println()

	// 3. Create a request object
	req := &fetcher.Request{
		URL:    url,
		Method: "GET",
	}

	// 4. Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start := time.Now()
	// 5. Execute the request
	resp, err := f.Do(ctx, req)
	duration := time.Since(start)

	if err != nil {
		printError(err)
		return
	}

	// 6. Print the results
	printSection("RESPONSE DETAILS")
	
	// Colorize Status Code
	statusColor := Green
	if resp.StatusCode >= 400 {
		statusColor = Red
	} else if resp.StatusCode >= 300 {
		statusColor = Yellow
	}

	fmt.Printf("%sStatus Code:%s  %s%d%s\n", Bold, Reset, statusColor, resp.StatusCode, Reset)
	fmt.Printf("%sLatency:%s      %v\n", Bold, Reset, duration)
	fmt.Printf("%sBody Size:%s    %d bytes\n", Bold, Reset, len(resp.Body))
	fmt.Println()

	// Print a preview of the body with a nice box
	printSection("BODY PREVIEW")
	previewLen := min(len(resp.Body), 1000)
	content := string(resp.Body[:previewLen])
	fmt.Printf("%s%s%s\n", Cyan, content, Reset)
}

func printHeader() {
	fmt.Println()
	fmt.Printf("%s╔════════════════════════════════════════╗%s\n", Cyan, Reset)
	fmt.Printf("%s║           GOSCRAPE FETCHER             ║%s\n", Cyan, Reset)
	fmt.Printf("%s╚════════════════════════════════════════╝%s\n", Cyan, Reset)
	fmt.Println()
}

func printSection(title string) {
	fmt.Printf("%s%s%s\n", Purple, title, Reset)
	fmt.Printf("%s%s%s\n", Purple, strings.Repeat("-", len(title)), Reset)
}

func printError(err error) {
	fmt.Println()
	fmt.Printf("%s❌ ERROR: %v%s\n", Red, err, Reset)
	fmt.Println()
	log.Fatalf("Execution failed")
}



