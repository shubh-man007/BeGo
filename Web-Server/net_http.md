### 🌐 net/http Package – Deep Dive Summary

---

## Core Concepts

The `net/http` package in Go provides fundamental primitives for building HTTP servers and clients. It abstracts away low-level TCP handling, allowing developers to focus on request/response logic.

---

## HTTP Client-Side (Making Requests)

### `http.Request`

Represents an HTTP request.

**Fields:**

- `Method` → GET, POST, etc.
- `URL` → Full URL (parsed into a `*url.URL` struct)
- `Header` → Key-value map of headers
- `Body` → An `io.ReadCloser` for request body (e.g., for POST requests)
- `Context()` → Holds cancelation signals/deadlines

### Creating Requests

```go
req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
```

Adds context to the request → good for timeouts, cancellations.

### `http.Client`

Sends the request and receives the response.

```go
client := http.Client{Timeout: 10 * time.Second}
resp, err := client.Do(req)
```

### `http.Response`

Represents an HTTP response.

**Fields:**

- `StatusCode` → e.g., 200, 404
- `Header` → Key-value map of headers
- `Body` → `io.ReadCloser` (must close after reading)
- `ContentLength`, `Request`, etc.

### Example: Simple GET Request

```go
resp, err := http.Get("https://example.com")
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
fmt.Println(string(body))
```

---

## HTTP Server-Side

### `http.Server`

Represents the server instance.

```go
s := &http.Server{
    Addr: ":8080",
    Handler: mux, // router
    ReadTimeout: 10 * time.Second,
    WriteTimeout: 10 * time.Second,
}
log.Fatal(s.ListenAndServe())
```

### `http.ServeMux` (Multiplexer)

Maps URL paths to handlers (routes).

```go
mux := http.NewServeMux()
mux.HandleFunc("/hello", handlerFunc)
```

---

## Concurrency & Goroutines

- Every incoming request is handled **in a new goroutine**.
- The handler function is safe to be long-running or blocking.
- This enables parallel processing of client requests.

---

## `context.Context`

- Propagated in HTTP requests via `http.NewRequestWithContext()`.
- Useful for timeouts, cancellations, deadlines.
- Server can check `.Done()` channel to know if client disconnected.

```go
ctx := req.Context()
select {
  case <-ctx.Done():
    log.Println("Client cancelled the request")
}
```

---

## Middleware Pattern

```go
func loggingMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    log.Println(r.Method, r.URL.Path)
    next.ServeHTTP(w, r)
  })
}
```

Use by wrapping handlers:

```go
mux.Handle("/hello", loggingMiddleware(http.HandlerFunc(handler)))
```

---

## Handler Interface and Custom Routing

### `http.Handler` Interface

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

Any type implementing `ServeHTTP()` is a valid HTTP handler.

### `http.HandlerFunc`

```go
type HandlerFunc func(ResponseWriter, *Request)
```

`HandlerFunc` allows regular functions to satisfy `Handler`.

### How Routing Works

```go
func myHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello!")
}

mux := http.NewServeMux()
mux.Handle("/", http.HandlerFunc(myHandler))
```

Or simply:

```go
mux.HandleFunc("/", myHandler)
```

`HandleFunc` wraps `myHandler` with `http.HandlerFunc`, which implements `ServeHTTP()`.

### Behind the scenes:

```go
mux.Handle("/", http.HandlerFunc(func(w, r) {
    w.Write([]byte("Hello"))
}))
```

When a request hits `/`, it calls `ServeHTTP(w, r)` on your handler.

---

## Summary of Concurrency

- Server spawns **one goroutine per incoming request**.
- ResponseWriter is safe per request.
- Clients can also run multiple requests concurrently via goroutines.

---

## Additional Utilities

- `http.ServeFile(w, r, "index.html")` → serves static file
- `http.Redirect(w, r, "/new", http.StatusFound)` → HTTP redirect
- `http.StripPrefix`, `http.FileServer` → serve assets

---