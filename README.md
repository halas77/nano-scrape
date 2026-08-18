# ⚡Nano Scrape

##  Overview

**Nano Scrape** makes extracting web data effortless. Built to harness Go’s native concurrency, it delivers blazing-fast performance without compromising simplicity. Whether you are parsing raw HTML strings or scraping live web pages at scale, Nano Scrape handles everything—from automatic network stream cleanup to proxy rotation and structured data exports.

## Features

*  **Blazing Fast:** Powered by Go’s lightweight goroutines for high-throughput, concurrent crawling.
*  **Flexible Parsing:** Effortlessly parse local HTML strings, bytes, standard `io.Reader` streams, or direct web URLs.
*  **Zero Resource Leaks:** Fully managed HTTP connections so you never have to worry about unclosed network streams.
*  **Intuitive Traversal:** Query the DOM using clean, CSS-like and attribute-based selectors.
*  **Export Ready:** One-line exports directly to **JSON**, **CSV**, or **Markdown**.

## Installation

Add **Nano Scrape** to your Go module using `go get`:

```bash
go get github.com/halas77/nano-scrape

```
