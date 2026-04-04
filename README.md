# Mirage

Universal command‑line HTTP client with browser fingerprint spoofing (TLS/JA\* + HTTP headers).

Mirage lets you send HTTP(S) requests that look like real browser traffic, by impersonating Chrome/Firefox on different OSes. It controls both high‑level HTTP headers (User‑Agent, client hints) and low‑level TLS parameters that form JA\* fingerprints.

---

## Features

- Single entrypoint:
  - `mirage http URL` – perform one HTTP(S) request.
- Browser impersonation:
  - `--fp chrome-android`, `--fp firefox-linux`, etc.
  - Controls both HTTP headers and TLS/JA\* fingerprint used for detection.
- Flexible request building:
  - Any HTTP method via `--method` (`GET`, `POST`, `PUT`, ...).
  - Multiple headers: `--header "Key:Value"` (repeatable).
  - Multiple cookies: `--cookie "name=value"` (repeatable).
  - Request body from `stdin` via `--body` flag.
- Output modes:
  - `meta` – status and timings (for scripting).
  - `resp` – only response body.
  - `full` – request + status + headers + body (curl -v style).
- Output formats:
  - `plain` – human‑readable text (default).
  - `json` – structured JSON with request/response fields (for tools and logs).

---

## Installation

Assuming you have Go installed (1.21+):

You can build mirage using folowing commands:

```bash
git clone https://github.com/shpak86/mirage
cd mirage
go build -o mirage .
```

To install mirage use go install:

```bash
go install github.com/shpak86/mirage/cmd/mirage
```

Precompiled binaries for different platforms ad OSes are available [here](https://github.com/shpak86/mirage/releases).

---

## Usage

Basic syntax:

```bash
mirage http URL [flags]
```

### Main flags

- `-f, --fp string`  
  Fingerprint profile in format `PLATFORM-OS`.  
  Platforms: `chrome`, `chrome58`, `chrome62`, `chrome70`, `chrome72`, `chrome83`, `chrome87`, `chrome96`, `chrome100`, `chrome102`, `chrome106`, `chrome120`, `chrome120`, `chrome145`, `firefox`, `firefox55`, `firefox56`, `firefox63`, `firefox65`, `firefox99`, `firefox102`, `firefox105`, `firefox120`, `firefox148`.
  OS: `linux`, `windows`, `mac`, `android`, `macos`.  
  Controls both HTTP headers (UA, client hints) and low-level TLS/JA\* fingerprint parameters used for browser impersonation.  
  Default: `chrome-android`.

- `-m, --method string`  
  HTTP method (`GET`, `POST`, `PUT`, `DELETE`, `HEAD`, `PATCH`, `OPTIONS`, `CONNECT`).  
  Default: `GET`.

- `-H, --header KEY:VALUE`  
  Set HTTP header. Can be repeated multiple times.  
  Example: `-H "X-Custom-Header:value1" -H "X-Custom-Header:value2"`.

- `-C, --cookie "name=value"`  
  Add cookie. Can be repeated multiple times.

- `-b, --body`  
  Read request body from `stdin` (e.g. `echo 'json' | mirage http ... -b`).

- `-o, --output string`  
  Output mode:
  - `meta` – status + timings (currently status).
  - `resp` – response body (default).
  - `full` – request + response (status, headers, body).

- `-F, --format string`  
  Output format:
  - `plain` – text.
  - `json` – structured JSON.
  Default: `plain`.

- `--proxy string`  
  HTTP proxy URL (e.g. `socks5://127.0.0.1:8080`).
  If set, all requests will be routed through the specified proxy.
  Default: no proxy.

- `--version string`  
  HTTP protocol version (`1.1`, `2`, `3`).
  If not specified, the version is determined automatically based on the fingerprint and server capabilities.
  Default: auto.

---

## Examples

### Simple GET

```bash
mirage http https://example.com -m GET -f chrome-android
```

### GET with custom headers and cookies

```bash
mirage http https://httpbin.org/headers \
  -m GET \
  -f firefox-linux \
  -H "X-Debug:1" \
  -H "Accept:application/json" \
  -C "session=313373" \
  -o full \
  -F plain
```

### POST with body from stdin

```bash
echo '{"user":"alice","pass":"secret"}' | \
  mirage http https://httpbin.org/post \
    -m POST \
    -f chrome-windows \
    -b \
    -F json \
    -o full
```

In `json` format with `-o full`, Mirage prints a JSON object with `request` and `response` sections (method, URL, headers, body, status, headers, body).

---

## How fingerprinting works (high level)

Internally, Mirage uses the `surf` client to impersonate real browsers:

- Parses `--fp` into `browser` and `os` parts (e.g. `chrome-android`).
- Applies OS preset: `Android()`, `Windows()`, `Linux()`, `MacOS()`.
- Applies browser preset: `Chrome()` or `Firefox()`, including TLS/JA\* and HTTP/2 fingerprint tuning.
- Builds the request and forwards your headers/cookies/body.

The result is traffic that looks much closer to real browser sessions than a default `net/http` client.

---

## Roadmap

Planned ideas:

- Additional protocols: WebSocket (`mirage ws`), raw TCP (`mirage tcp`).
- More fingerprint presets and better diagnostics (JA3/JA4 preview).
- Built‑in benchmarking mode to compare how different fingerprints behave against a target.
