# Attribute Traversal

## Functions

### `Find`

```go
func (t *Tag) Find(name string, attrs []*Attribute, cb TagCallback)
```

Searches through the HTML document tree to find elements that match the specified tag name and attributes. Every time a match is found, it executes a custom callback function.

- **What it takes:** \* `name`: The target HTML tag name to look for (e.g., `"div"`, `"span"`).
  - `attrs`: A slice of specific attributes the tag must have to match.
  - `cb`: A callback function (`TagCallback`) that runs instantly when a matching tag is found.
- **What it returns:** Nothing. It yields results immediately to the callback function during the search.

### `💡 Example`

```go
// Process matching elements on the fly using a callback function
page.Find("div", nil, func(matchedTag *engine.Tag) {
    // This block runs automatically every time a <div> is encountered
    fmt.Println("Streamed match found for tag:", matchedTag.Name)
})
```

### `FindAll`

```go
func (t *Tag) FindAll(name string, attribute ...[]*Attribute) *Tags
```

Searches the HTML tree and collects **all** matching elements into a single list collection.

- **What it takes:**
  - `name`: The target HTML tag name to look for.
  - `attribute`: _(Optional)_ A slice of specific attributes to filter the results.
- **What it returns:** A pointer to a `Tags` collection containing every matching element found.

#### `💡 Example`

```go
// Gather all links on the webpage into a slice collection
allLinks := page.FindAll("a")
fmt.Printf("Found %d total links on the page.\n", len(*allLinks))
```

### `FindFirst`

```go
func (t *Tag) FindFirst(name string, attr ...[]*Attribute) *Tag
```

Searches the HTML tree and returns only the **very first** element that matches the criteria. It automatically limits the search depth to maximize performance once a match is found.

- **What it takes:**
  - `name`: The target HTML tag name to look for.
  - `attr`: _(Optional)_ A slice of specific attributes to filter the result.
- **What it returns:** A pointer to the first matching `Tag` element found, or `nil` if no match exists.

#### `💡 Example`

```go
// Quickly grab the first instance of a heading tag
firstHeading := page.FindFirst("h1")
if firstHeading != nil {
    fmt.Println("First heading found:", firstHeading.Name)
}
```

## Attribute

`Attribute` represents a key-value pair used to target specific HTML element attributes (such as `class`, `id`, `href`, or `data-*` properties) when searching a document.

```go
type Attribute struct {
    Key   string // The attribute name (e.g., "class", "id", "href")
    Value string // The expected value of the attribute (e.g., "details", "main-content")
}
```

#### `💡 Usage Example`

You can pass a slice of `*Attribute` pointers into filtering functions like `FindAll` or `FindFirst` to narrow down your search results to specific elements.

```go
package main

import (
	"fmt"
	"log"
)

func main() {
	// Assume 'scrape' is an HTML document already parsed via InitDocument or LoadDocument
	scrape, _ := engine.LoadDocument("https://example.com")

	// 1. Define the specific attributes you want to filter by
	params := []*engine.Attribute{
		{
			Key:   "class",
			Value: "details",
		},
	}

	// 2. Search for <div> tags that specifically match your attribute parameters
	mainElements := scrape.FindAll("div", params)

	// 3. Print out the matching HTML blocks
	fmt.Println(mainElements.Print())
}
```

### `Find` / `FindAll` with Text Filtering

When you pass an `Attribute` with the special key `"string"`, the search engine start inspecting the direct plain text inside the tag itself.

- **Exact Match Target:** It evaluates the immediate text children belonging to that tag.
- **Strict Depth Control:** It intentionally ignores text contained inside nested children elements (like a `<span>` or `<strong>` inside a `<div>`). It only looks at the string fragments that belong to the parent tag directly.

#### `💡 Usage Example`

Given this HTML snippet, the text `"World"` and `"Hello"` belong directly to the `<div>` tag, while `"Span's Inner text"`belongs to the `<span>` tag.

```go
package main

import (
	"fmt"
	"log"
)

func main() {
	// Sample HTML structure
	input := `<div class="details"> Hello <span> Span's Inner text </span> World </div>`

	scrape, _ := engine.InitDocument(input)

	// Configure the parameter to look for the raw string "World" inside a tag
	params := []*engine.Attribute{
		{
			Key:   "string",
			Value: "World",
		},
	}

	// This WILL find the <div> because "World" is part of its immediate text content.
	// It will NOT match the <span> because the search skips nested child strings.
	matches := scrape.FindAll("div", params)

	fmt.Println(matches.Print())
}
```

#### 🔍 Search Filtering Breakdown

- **`<div>` inspection:** Contains text segments `" Hello "` and `" World "`. The value matches your parameter. **(Match Found)**
- **`<span>` inspection:** Contains text segment `" Span's Inner text "`. If you were to search for `<span>` with a `"string"` value of `"Span's Inner text"`, it would match. Searching for `<div>` with a `"string"` value of `"Span's Inner text"` returns nothing because it belongs to the child.

## Combining HTML Attribute and Text Filtering

You can combine standard HTML attributes (like `class` or `id`) with the special `"string"` key inside the same parameter slice. When you pass multiple attributes, the search engine treats them as an **AND** condition. An HTML tag must match **both** the structural attributes and contain the specified immediate text to be returned.

#### `💡 Usage Example`

This example demonstrates how to find an element that simultaneously has a specific class name and contains a specific direct text fragment.

```go
package main

import (
	"fmt"
)

func main() {
	// Sample HTML structure with two similar divs
	input := `
		<div class="details"> Hello <span> Span's Inner text </span> World </div>
		<div class="sidebar"> Hello <span> Span's Inner text </span> World </div>
	`

	scrape, _ := engine.InitDocument(input)

	// Configure parameters to match: class="details" AND containing the text "World"
	params := []*engine.Attribute{
		{
			Key:   "class",
			Value: "details",
		},
		{
			Key:   "string",
			Value: "World",
		},
	}

	// This will ONLY find the first <div> because it satisfies both criteria.
	// The second <div> is skipped because its class is "sidebar".
	matches := scrape.FindAll("div", params)
}
```

#### `🔍 Filter Multi-Match Rule`

- **Tag Name:** Must be a `<div>`.
- **Criteria 1 (`class`):** Must exactly equal `"details"`.
- **Criteria 2 (`string`):** Must directly contain the text block `"World"` (ignoring nested elements).
