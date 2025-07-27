### Go URL Shortener

  Tiny URL shortener I built while learning Go. It creates short codes for long URLs and redirects users to the original link.

  ## Features

  * `POST /shorten` to create a short code for a valid http/https URL
  * `GET /{code}` returns a 302 redirect to the original URL
  * In‑memory store protected by RWMutex
  * Crypto safe base62 code generator
  * Basic global rate limiting with `golang.org/x/time/rate`
  * Graceful shutdown on SIGINT or SIGTERM
  * chi middleware for request logging and panic recovery

  ## Run locally

  ```bash
  go run ./cmd/api
  ```

  Server listens on `:8000`.

  ## API

  ### POST /shorten

  **Request body**

  ```json
  { "text": "https://google.com" }
  ```

  **PowerShell example**

  ```powershell
  Invoke-RestMethod -Method Post -Uri http://localhost:8000/shorten `
    -ContentType 'application/json' `
    -Body '{"text":"https://google.com"}'
  ```

  **curl.exe on Windows**

  ```powershell
  curl.exe -X POST http://localhost:8000/shorten ^
    -H "Content-Type: application/json" ^
    -d "{\"text\":\"https://google.com\"}"
  ```

  **Response**

  ```json
  { "code": "abc123" }
  ```

  ### GET /{code}

  Visit in a browser or curl:

  ```
  http://localhost:8000/abc123
  ```

  Response headers should include:

  ```
  HTTP/1.1 302 Found
  Location: https://google.com
  ```

  ## Project layout

  ```
  cmd/api/main.go          # server bootstrap and graceful shutdown
  internal/handlers        # HTTP handlers (Shorten, Redirect)
  internal/store           # MemStore (map + mutex)
  internal/util            # JSON helpers
  ```

  ## Development notes

  * Race detector: `go test -race ./...` (once tests are added)
  * Lint/vet: `go vet ./...`
  * Replace MemStore with a SQL store later without touching handlers

  ## Future ideas

  * Click stats endpoint and link expiry
  * Per-IP rate limiting
* SQLite or Postgres persistence
* Dockerfile + CI pipeline
* Tiny HTML front end

## License

MIT. See `LICENSE`.



