# CSS Traversal

`nano-scrape` allows you to traverse and query HTML documents using standard CSS selectors. Under the hood, this is powered by the robust `cascadia` CSS selector library.

## Supported Selectors

You can use standard CSS selectors to find elements. Here are the most common patterns:

| Selector Type  | Syntax   | Example          | Description                                |
| -------------- | -------- | ---------------- | ------------------------------------------ |
| **Class**      | `.class` | `.price`         | Matches elements with `class="price"`      |
| **ID**         | `#id`    | `#price`         | Matches elements with `id="price"`         |
| **Descendant** | `A B`    | `article .title` | Matches `.title` nested inside `<article>` |

## Core Selection Functions

There are three primary methods on the `*Tag` struct for CSS selection.

### 1. `SelectFirst`

Retrieves the **first** element that matches the CSS selector.

- **Signature**: `func (t *Tag) SelectFirst(selector string) *Tag`
- **Performance Note**: This method is highly optimized. It sets an internal limit to stop traversal and exits early as soon as the first match is found.

### 2. `SelectAll`

Retrieves **all** elements matching the CSS selector.

- **Signature**: `func (t *Tag) SelectAll(selector string) *Tags`
- **Result**: Returns a `Tags` collection wrapper, which provides further methods like `.First()`, `.Export()`, `.Map()`, `.ToJSON()`, etc.

### 3. `Select`

Iterates over all matching elements and runs a custom callback function for each match.

- **Signature**: `func (t *Tag) Select(selector string, f func(*Tag))`

## Step-by-Step Example

Let's look at how to use these functions with a practical example.

### 1. The HTML Source

Assume we have the following HTML structure:

```html
<div id="store">
  <h1 class="store-title">Tech Emporium</h1>
  <div class="product-list">
    <div class="product-item" data-id="101">
      <span class="name">Wireless Mouse</span>
      <span class="price">$29.99</span>
    </div>
    <div class="product-item" data-id="102">
      <span class="name">Mechanical Keyboard</span>
      <span class="price">$89.99</span>
    </div>
  </div>
</div>
```

### 2. The Go Code

Here is how we load the document and extract details step by step:

```go
package main

import (
	"fmt"
	"github.com/halas77/nano-scrape/nano"
)

func main() {
	htmlContent := `
	<div id="store">
	  <h1 class="store-title">Tech Emporium</h1>
	  <div class="product-list">
	    <div class="product-item" data-id="101">
	      <span class="name">Wireless Mouse</span>
	      <span class="price">$29.99</span>
	    </div>
	    <div class="product-item" data-id="102">
	      <span class="name">Mechanical Keyboard</span>
	      <span class="price">$89.99</span>
	    </div>
	  </div>
	</div>
	`

	// Step 1: Initialize the document
	doc, err := nano.InitDocument(htmlContent)
	if err != nil {
		panic(err)
	}

	// Step 2: Use SelectFirst to get a single, specific element (Store Title)
	titleTag := doc.SelectFirst("h1.store-title")
	if titleTag != nil {
		fmt.Println("Store Title:", titleTag.Text())
		// Output: Store Title: Tech Emporium
	}

	// Step 3: Use SelectAll to get a collection of elements (Product Items)
	products := doc.SelectAll("div.product-item")
	fmt.Printf("Found %d products.\n", len(*products))
	// Output: Found 2 products.

	// Step 4: Loop through the collected tags
	for _, item := range *products {
		name := item.SelectFirst(".name").Text()
		price := item.SelectFirst(".price").Text()
		fmt.Printf("- %s: %s\n", name, price)
	}
	// Output:
	// - Wireless Mouse: $29.99
	// - Mechanical Keyboard: $89.99

	// Step 5: Use Select to iterate directly using a callback
	fmt.Println("Iterating with Select callback:")
	doc.Select(".product-item", func(item *nano.Tag) {
		name := item.SelectFirst(".name").Text()
		fmt.Println("Product Item:", name)
	})
}
```

---

## Important Tips

> [!TIP]
> **Use Early Exits for Performance**
> When you only need the first match (e.g. searching for a title, a main container, or single attribute value), prefer `SelectFirst`. Because it avoids traversing the rest of the HTML tree after finding the match, it is much faster and uses fewer CPU allocations.

> [!IMPORTANT]
> **Always Check for nil**
> If a selector does not find any match, `SelectFirst` returns `nil`. Always check that the returned `*Tag` is not `nil` before calling methods like `.Text()` to avoid runtime panics.
