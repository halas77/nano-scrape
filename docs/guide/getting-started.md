# Getting Started

Welcome to **Nano Scrape** – a high‑performance, lightweight Go library for web scraping.

## Installation

```bash
go get github.com/halas77/nano-scrape
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/halas77/nano-scrape/nano"
)

func main() {
    // Load an HTML document (local string or remote URL)
    doc, err := nano.InitDocument(`<html><body><div class="price">$199.99</div></body></html>`)
    if err != nil { panic(err) }

    // Use a CSS‑like selector to extract the price
    priceTag := doc.SelectFirst(".price")
    fmt.Println("Price:", priceTag.Text())
}
```

The example demonstrates how to:

1. Initialise a document from a string.
2. Select a node with a CSS selector.
3. Extract the text content.

## Core Concepts

- **Request** – Handles HTTP GET/POST with optional proxy rotation.
- **Tag** – Represents an HTML element; provides traversal, selection, and text extraction.
- **Tags** – A collection of `Tag` objects with helpers for bulk operations.
- **Export** – Convert scraped data to JSON, CSV, or Markdown.

Explore the next sections for deeper details on each feature.
