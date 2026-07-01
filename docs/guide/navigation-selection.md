# Navigation & Selection

Learn how to traverse the DOM and select elements using CSS‑like selectors.

## Selecting Elements
```go
// Load the document
doc, _ := engine.InitDocument(htmlString)

// Select all product cards
cards := doc.SelectAll(".product-card")
```

## Traversal Helpers
- `Select` – iterate over matching elements.
- `SelectFirst` – return the first match.
- `Find`, `FindAll` – advanced attribute‑based queries.

## Example: Extract Titles
```go
cards.Select("h2.title", func(t *engine.Tag) {
    fmt.Println("Title:", t.Text())
})
```
