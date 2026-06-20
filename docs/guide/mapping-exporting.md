# Mapping & Exporting

Convert scraped data into structured formats (JSON, CSV, Markdown) using the powerful mapping API.

## Mapping a Collection
```go
// Assume `products` is a *engine.Tags collection
mapping := map[string]string{
    "sku":   "@data-sku",
    "title": "h3.name",
    "price": "span.price",
    "url":   "a.link @href",
}

mapped := products.Map(mapping)
```

The map keys become the column names, and the selectors extract values. The `@` prefix fetches an attribute value.

## Export Functions
- `ExportJSON(data []map[string]string) ([]byte, error)` – returns pretty‑printed JSON.
- `ExportCSV(data []map[string]string) (string, error)` – returns CSV string with headers.
- `ExportMD(data []map[string]string) (string, error)` – returns a Markdown table.

## Writing Directly to Files
```go
engine.WriteMappedJSON("products.json", mapped)
engine.WriteMappedCSV("products.csv", mapped)
engine.WriteMappedMD("products.md", mapped)
```

## Exporting Raw Tag Collections
If you want to dump the raw `Tags` collection without custom mapping:
```go
jsonBytes, _ := products.ToJSON()
csvStr, _ := products.ToCSV()
mdStr, _ := products.ToMD()
```

These helpers automatically include all tag fields (name, text, class, id, attributes).

Explore the **API Reference** for the full list of export utilities.
