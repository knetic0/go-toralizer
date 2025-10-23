# Toraliz - SOCKS4 Tor Proxy Client

A lightweight SOCKS4 client written in Go that connects to websites through the Tor network.

## Features

- 🧅 Direct SOCKS4 protocol implementation
- 🔒 Anonymous HTTP requests through Tor
- 📦 No external dependencies (pure Go)
- ⚡ Simple and educational codebase

## Prerequisites

- Go 1.16+
- Tor service running locally (default: `127.0.0.1:9050`)

## Installation
```bash
go build -o toraliz
```

## Usage
```bash
./toraliz <host> <port>
```

**Example:**
```bash
./toraliz check.torproject.org 80
```

## How it works

1. Connects to local Tor SOCKS4 proxy (port 9050)
2. Sends SOCKS4 CONNECT request
3. Establishes tunnel to destination
4. Sends HTTP GET request
5. Retrieves and displays response

## Topics for GitHub

`tor` `socks4` `proxy` `golang` `privacy` `anonymity` `networking` `socks-proxy`

## License

MIT
