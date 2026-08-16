# Exporting

The `nano-scrape` library provides built-in tools to export your scraped HTML elements into structured data formats, specifically **JSON**, **CSV**, and **Markdown**. You can either export raw parsed tags directly or map them into custom key-value pairs before serialization.

## TagExport

`TagExport` is a cleanly structured representation of a `Tag` designed for direct serialization. It filters out internal DOM references and focuses entirely on the key fields of the element.

```go
type TagExport struct {
	Name       string            `json:"name"`
	Text       string            `json:"text"`
	Class      string            `json:"class,omitempty"`
	ID         string            `json:"id,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}
```

## Raw Export Functions

These methods are attached to `Tag` or `Tags` and allow you to serialize or write HTML elements in their raw parsed form.

### `Export` (Single Tag)

```go
func (t Tag) Export() TagExport
```

Converts a single parsed `Tag` into a serializable `TagExport` structure.

- **What it takes:** Nothing (called on a `Tag` instance).
- **What it returns:** A `TagExport` struct with populated name, inner text, class, id, and attributes.

### `Export` (Tag Collection)

```go
func (ts *Tags) Export() []TagExport
```

Converts a collection of tags (`Tags`) into a slice of serializable `TagExport` structures.

- **What it takes:** Nothing (called on a `*Tags` collection).
- **What it returns:** A slice of `TagExport` structs.

### `ToJSON`

```go
func (ts *Tags) ToJSON() ([]byte, error)
```

Converts the tags collection directly into a pretty-printed JSON byte slice.

- **What it takes:** Nothing.
- **What it returns:** A pretty-printed JSON byte slice (`[]byte`) and an `error` if serialization fails.

### `ToCSV`

```go
func (ts *Tags) ToCSV() (string, error)
```

Converts the tags collection directly into a CSV string.

- **What it takes:** Nothing.
- **What it returns:** A CSV string and an `error`.
- **Note:** The CSV columns are `name`, `text`, `class`, `id`, and `attributes`. Semicolons are used to separate key-value attribute pairs (e.g. `href=http://example.com; target=_blank`).

### `ToMD`

```go
func (ts *Tags) ToMD() (string, error)
```

Converts the tags collection directly into a Markdown table representation.

- **What it takes:** Nothing.
- **What it returns:** A Markdown table string and an `error`.
- **Note:** Columns match the CSV format. Any internal pipe characters (`|`) are escaped, and newline characters are removed to prevent breaking the Markdown table layout.

### `WriteJSON` / `WriteCSV` / `WriteMD`

```go
func (ts *Tags) WriteJSON(filename string) error
func (ts *Tags) WriteCSV(filename string) error
func (ts *Tags) WriteMD(filename string) error
```

Convenience methods to serialize and save the tags collection directly to a file on your system.

- **What it takes:** `filename`: The destination path of the file.
- **What it returns:** An `error` if writing to the file or serialization fails.

## Mapping & Mapped Export Functions

Often, you do not want the raw tag dump. Instead, you want to parse out specific fields (like the URL, the price text, or custom identifiers) and group them into records.

### `Map`

```go
func (ts *Tags) Map(mapping map[string]string) []map[string]string
```

Extracts structured data from a collection of elements using a key-to-selector mapping scheme.

- **What it takes:** `mapping`: A map of user-defined field names to query selectors.
- **What it returns:** A slice of key-value maps representing the extracted data rows.

#### Selector Syntax in Mappings

- **Text Extraction:** Providing a standard CSS selector (e.g., `h2.title`) extracts the inner text of the first matching child element.
- **Attribute Extraction:** Appending `@attribute_name` (e.g., `a.link @href`) extracts that specific attribute value from the matched element.
- **Self-Attribute Extraction:** Providing just `@attribute_name` (e.g., `@data-sku`) extracts the attribute directly from the root element of each iteration.

### `ExportJSON` (Mapped)

```go
func ExportJSON(data []map[string]string) ([]byte, error)
```

Converts a slice of mapped key-value pairs into a pretty-printed JSON byte slice.

- **What it takes:** `data`: A slice of string maps.
- **What it returns:** A pretty-printed JSON byte slice and an `error`.

### `ExportCSV` (Mapped)

```go
func ExportCSV(data []map[string]string) (string, error)
```

Converts a slice of mapped key-value pairs into a CSV string.

- **What it takes:** `data`: A slice of string maps.
- **What it returns:** A CSV string containing the header row followed by values, and an `error`.
- **Note:** Headers are gathered from all keys in the data slice and sorted alphabetically for deterministic output.

### `ExportMD` (Mapped)

```go
func ExportMD(data []map[string]string) (string, error)
```

Converts a slice of mapped key-value pairs into a Markdown table.

- **What it takes:** `data`: A slice of string maps.
- **What it returns:** A Markdown table string and an `error`.
- **Note:** Columns are sorted alphabetically. Any pipe (`|`) or newline (`\n`) characters in the values are cleaned to keep the table formatting intact.

### `WriteMappedJSON` / `WriteMappedCSV` / `WriteMappedMD`

```go
func WriteMappedJSON(filename string, data []map[string]string) error
func WriteMappedCSV(filename string, data []map[string]string) error
func WriteMappedMD(filename string, data []map[string]string) error
```

Convenience functions to serialize and save mapped key-value data directly to system files.

- **What it takes:**
  - `filename`: The destination path of the file.
  - `data`: A slice of string maps.
- **What it returns:** An `error` if serialization or file writing fails.

## Step-by-Step Example

This example demonstrates how to parse a raw HTML snippet containing product cards, select elements, map the target information, and save the data in different formats.

### 1. The HTML Source

Imagine we are scraping the following list of products:

```html
<div class="product-grid">
  <div class="product-card" data-sku="SKU-8891">
    <h3 class="name">Wireless Charger</h3>
    <span class="price">$19.99</span>
    <a class="details-btn" href="/products/wireless-charger">View Details</a>
  </div>
  <div class="product-card" data-sku="SKU-4432">
    <h3 class="name">Bluetooth Headset</h3>
    <span class="price">$49.99</span>
    <a class="details-btn" href="/products/bluetooth-headset">View Details</a>
  </div>
</div>
```

### 2. The Go Code

The following program loads the HTML, queries all `.product-card` elements, maps their child contents using custom selectors, and writes the output to both JSON and Markdown table files:

```go
package main

import (
	"fmt"
	"log"

	"github.com/halas77/nano-scrape/nano"
)

func main() {
	htmlContent := `
	<div class="product-grid">
	  <div class="product-card" data-sku="SKU-8891">
	    <h3 class="name">Wireless Charger</h3>
	    <span class="price">$19.99</span>
	    <a class="details-btn" href="/products/wireless-charger">View Details</a>
	  </div>
	  <div class="product-card" data-sku="SKU-4432">
	    <h3 class="name">Bluetooth Headset</h3>
	    <span class="price">$49.99</span>
	    <a class="details-btn" href="/products/bluetooth-headset">View Details</a>
	  </div>
	</div>
	`

	// 1. Initialize the document
	doc, err := nano.InitDocument(htmlContent)
	if err != nil {
		log.Fatalf("Failed to initialize document: %v", err)
	}

	// 2. Select all product card elements
	cards := doc.SelectAll(".product-card")
	fmt.Printf("Selected %d product cards\n", len(*cards))

	// 3. Define mapping to extract details
	// "@data-sku" matches attribute data-sku on the card itself
	// "h3.name" matches text inside <h3 class="name">
	// "span.price" matches text inside <span class="price">
	// "a.details-btn @href" matches href attribute of <a class="details-btn">
	mapping := map[string]string{
		"sku":   "@data-sku",
		"name":  "h3.name",
		"price": "span.price",
		"url":   "a.details-btn @href",
	}

	// 4. Map the collection to structured slices of key-value maps
	mappedData := cards.Map(mapping)

	// 5. Convert mapped data to a Markdown table string and print it
	mdTable, err := nano.ExportMD(mappedData)
	if err != nil {
		log.Fatalf("Failed to export to Markdown: %v", err)
	}
	fmt.Println("\nGenerated Markdown Table:")
	fmt.Println(mdTable)

	// 6. Write mapped data directly to files
	err = nano.WriteMappedJSON("products.json", mappedData)
	if err != nil {
		log.Fatalf("Failed to write JSON file: %v", err)
	}
	fmt.Println("Successfully wrote products.json")

	err = nano.WriteMappedMD("products.md", mappedData)
	if err != nil {
		log.Fatalf("Failed to write Markdown file: %v", err)
	}
	fmt.Println("Successfully wrote products.md")
}
```

### 3. Output Formats

#### JSON Output (`products.json`)

```json
[
  {
    "name": "Wireless Charger",
    "price": "$19.99",
    "sku": "SKU-8891",
    "url": "/products/wireless-charger"
  },
  {
    "name": "Bluetooth Headset",
    "price": "$49.99",
    "sku": "SKU-4432",
    "url": "/products/bluetooth-headset"
  }
]
```

#### Markdown Output (`products.md`)

| name              | price  | sku      | url                         |
| ----------------- | ------ | -------- | --------------------------- |
| Wireless Charger  | $19.99 | SKU-8891 | /products/wireless-charger  |
| Bluetooth Headset | $49.99 | SKU-4432 | /products/bluetooth-headset |
