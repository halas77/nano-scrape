# HTTP Request & Client

`engine.Request` provides a reusable HTTP client with cookie‑jar support, custom headers and optional proxy rotation.

## Initialise a client

```go
req := engine.InitRequest()
// Optionally add default headers
req.Header.Set("User-Agent", "NanoScrape/1.0")
```

## Proxy Rotation

```go
proxies := []string{"http://1.2.3.4:8080", "http://5.6.7.8:3128"}
req.ProxyRotator(proxies...)
```

The client will automatically pick the next proxy for each request.

## GET request

```go
body, err := req.Execute("https://example.com", "GET")
if err != nil { panic(err) }
fmt.Println(string(body))
```

## POST JSON

```go
payload := map[string]string{"email": "you@example.com", "password": "secret"}
resp, err := req.MakeJSONPostRequest("https://example.com/login", "POST", payload)
```

## POST Form

```go
form := map[string]string{"q": "search term"}
resp, err := req.MakeFormPostRequest("https://example.com/search", "POST", form)
```

## Cookies

Retrieve cookies for a specific URL:

```go
cookies := req.CookiesFor("https://example.com")
for _, c := range cookies {
    fmt.Printf("%s=%s; ", c.Name, c.Value)
}
```

## Bot Block Detection

`engine` automatically checks for common bot‑block responses (e.g., Cloudflare 403/429). If a block is detected, the request returns an error.

Read the full API reference for more advanced usage.
