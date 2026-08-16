# Init Document

## Tag

`Tag` represents a single HTML element parsed from a document or webpage. You can use its public fields to inspect the element's properties.

```go
type Tag struct {
    Name  string
    Attrs []html.Attribute
    Class string
    Id    string
}
```

## Core Functions

There are two primary methods on the `*Tag` struct for initiating scraping.

### 1. `InitDocument`

```go
func InitDocument(input any) (*Tag, error)
```

Converts local HTML data into a queryable `Tag` structure.

- **What it takes:** A `string` of raw HTML, a `[]byte` slice, or any standard Go reading stream (`io.Reader`).
- **What it returns:** The root HTML `Tag` of the document, or an error if the input cannot be processed.

### 2. `LoadDocument`

```go
func LoadDocument(url string) (*Tag, error)
```

Downloads an HTML document from the internet and prepares it for scraping.

- **What it takes:** A standard website URL string (e.g., `"https://example.com"`).
- **What it returns:** The root HTML `Tag` of the downloaded webpage, or an error if the website is unreachable.
- **Note:** This function completely manages the web connection for you. You do not need to worry about closing network streams or handling resource leaks.

#### `Usage Example`

```go
package main

import (
	"fmt"
	"log"
)

func main() {
	// Example 1: Parsing HTML from a local string variable
	htmlData := `<div id="main" class="content">Hello World</div>`

	page1, err := nano.InitDocument(htmlData)
	if err != nil {
		log.Fatalf("Error parsing string: %v", err)
	}
	fmt.Println("Parsed local tag:", page1.Name) // Outputs: "div"

	// Example 2: Parsing HTML directly from a live website URL
	page2, err := nano.LoadDocument("https://example.com")
	if err != nil {
		log.Fatalf("Error loading URL: %v", err)
	}
	fmt.Println("Parsed web tag:", page2.Name) // Outputs: "html"
}
```
