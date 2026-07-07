# Requests

## Client

`Client` handles web requests, session state (cookies).

```go
type Client struct {
    Header http.Header // Global headers appended automatically to every outgoing request
}
```

## Core Functions

There are primary methods on the `*Client` struct for initiating and making HTTP request.

### 1. `NewClient`

```go
func NewClient(baseHeader ...http.Header) *Client
```

Creates and configures a new HTTP client. It includes a built-in session cookie manager and defaults to a 10-second timeout.

- **What it takes:** _(Optional)_ An existing `http.Header` map to pre-populate global request headers.
- **What it returns:** An initialized pointer to a `Client`.

#### 💡 Example

```go
// Create a client with baseline headers
customHeader := make(http.Header)
customHeader.Set("Authorization", "Bearer token123")

client := engine.NewClient(customHeader)
```

---

### 2. `Get`

```go
func (c *Client) Get(targetURL string) (io.Reader, error)

```

A convenience wrapper around `Execute` specifically designed for issuing HTTP `GET` requests.

- **What it takes:**
- `targetURL`: The full destination URL string.

- **What it returns:** An `io.Reader` holding the response body, or an error if the request fails or returns a non-200 status code.

#### `💡 Example`

```go
responseStream, err := client.Get("https://api.example.com/data")
if err != nil {
    log.Fatal(err)
}

```

---

<!-- ### `Execute`

```go
func (c *Client) Execute(method, targetURL string, body ...io.Reader) (io.Reader, error)
```

Executes a generic HTTP request. It manages underlying network memory cleanup automatically.

- **What it takes:** \* `method`: The HTTP verb (e.g., `"GET"`, `"POST"`).
- `targetURL`: The full destination URL string.
- `body`: _(Optional)_ An optional input stream data payload.
- **What it returns:** An `io.Reader` holding the response body, or an error if the request fails or returns a non-200 status code.

#### `💡 Example`

```go
responseStream, err := client.Execute("GET", "<https://api.example.com/data>")
if err != nil {
    log.Fatal(err)
}
```

--- -->

### 3. `SendJSON`

```go
func (c *Client) SendJSON(method, targetURL string, payload any) (io.Reader, error)
```

Converts a Go struct or map into a JSON payload and transmits it with the appropriate header types (`application/json`).

- **What it takes:** \* `method`: HTTP verbs supporting data payloads (e.g., `"POST"`, `"PUT"`).
- `targetURL`: The destination URL.
- `payload`: Any Go data structure (struct, map, slice) to convert to JSON.
- **What it returns:** An `io.Reader` containing the server's response stream, or an error.

#### `💡 Example`

```go
userData := map[string]string{"username": "admin"}
resp, err := client.SendJSON("POST", "<https://example.com/api/login>", userData)
```

---

### 4. `SendForm`

```go
func (c *Client) SendForm(method, targetURL string, payload map[string]string) (io.Reader, error)
```

Submits standard key-value map data acting like a native web browser form submission (`x-www-form-urlencoded`).

- **What it takes:** \* `method`: HTTP verbs supporting data payloads (e.g., `"POST"`).
- `targetURL`: The destination URL.
- `payload`: A map of string keys and values representing form inputs.
- **What it returns:** An `io.Reader` containing the server's response stream, or an error.

#### 💡 Example

```go
formData := map[string]string{"email": "user@test.com", "submit": "true"}
resp, err := client.SendForm("POST", "<https://example.com/submit>", formData)
```

---

### 5. `CookiesFor`

```go
func (c *Client) CookiesFor(rawURL string) []*http.Cookie
```

Inspects and extracts the cookies currently saved in the client session for a given domain/URL.

- **What it takes:** A target URL string.
- **What it returns:** A slice of `http.Cookie` structs currently tracked for that domain.

#### `💡 Example`

```go
cookies := client.CookiesFor("<https://example.com>")
for _, cookie := range cookies {
    fmt.Printf("Active cookie session: %s=%s\n", cookie.Name, cookie.Value)
}
```
