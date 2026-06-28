# API Reference

This reference lists the public types and functions exported by the **Nano Scrape** Go package.

## Packages

- `github.com/halas77/nano-scrape/engine`

## Types

- `engine.Client`
  - Fields: `Header http.Header`
  - Methods:
    - `NewClient(baseHeader ...http.Header) *Client`
    - `ProxyRotator(proxies ...string)` – enable round-robin proxy rotation.
    - `Execute(method, targetURL string, body ...io.Reader) (io.Reader, error)` – perform HTTP request.
    - `SendJSON(method, targetURL string, payload any) (io.Reader, error)` – transmit structured data as JSON.
    - `SendForm(method, targetURL string, payload map[string]string) (io.Reader, error)` – submit x-www-form-urlencoded data.
    - `CookiesFor(rawURL string) []*http.Cookie`
- `engine.Tag`
  - Methods for DOM traversal and selection:
    - `Select(selector string, f func(*Tag))`
    - `SelectAll(selector string) *Tags`
    - `SelectFirst(selector string) *Tag`
    - `Find(name string, attrs []*Attribute, cb TagCallback)`
    - `FindAll(name string, attribute ...[]*Attribute) *Tags`
    - `FindFirst(name string, attr ...[]*Attribute) *Tag`
    - `Text() string`
    - `Print(depth ...uint16) string`
- `engine.Tags`
  - Collection helpers:
    - `First() *Tag`
    - `Map(mapping map[string]string) []map[string]string`
    - `ToJSON() ([]byte, error)`
    - `ToCSV() (string, error)`
    - `ToMD() (string, error)`
    - `WriteJSON(filename string) error`
    - `WriteCSV(filename string) error`
    - `WriteMD(filename string) error`
- `engine.Attribute`
  - Simple key/value pair used for attribute filtering.

## Export Functions (Top‑Level)

- `engine.ExportJSON(data []map[string]string) ([]byte, error)`
- `engine.ExportCSV(data []map[string]string) (string, error)`
- `engine.ExportMD(data []map[string]string) (string, error)`
- `engine.WriteMappedJSON(filename string, data []map[string]string) error`
- `engine.WriteMappedCSV(filename string, data []map[string]string) error`
- `engine.WriteMappedMD(filename string, data []map[string]string) error`

All functions return errors where appropriate; handle them in production code.

For complete type definitions and comments, refer to the source files in the `engine` directory.
