# Tor Client in Go

This is a simple Tor client written in Go that connects to a Tor SOCKS proxy
server and sends an HTTP HEAD request through the proxy to a specified destination.

## Features

- **SOCKS4 Protocol:** Connects to a Tor SOCKS4 proxy server to route traffic.
- **HTTP Request:** Sends an HTTP HEAD request to the destination server
  through the proxy.
- **IPv4 and IPv6 Support:** Supports both IPv4 and IPv6 addresses for
  destination connections.
- **Verbose Output:** Optional verbose mode for detailed logging.

## Requirements

- Go (version 1.11 or higher recommended)
- Access to a running Tor SOCKS proxy (default configuration uses `127.0.0.1:9050`)

## Usage

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/andrewthecoder/go-tormask.git
   ```

2. Navigate to the directory:

   ```bash
   cd go-tormask
   ```

3. Build the binary:

   ```bash
   go build
   ```

### Command Line Usage

```bash
./go-tormask [-i IP ADDRESS | -u URL] -p PORT [-v]
```

- `-i IP ADDRESS`: IP address of the destination server.
- `-u URL`: URL of the destination server.
- `-p PORT`: Port of the destination server.
- `-v`: Optional verbose mode for detailed logging.

## Examples

Connect to IPV4 destination server:

```bash
./go-tormask -i 89.207.132.170 -p 80 -v
```

Connect to IPV6 destination server:

```bash
./go-tormask -i 5be8:dde9:7f0b:d5a7:bd01:b3be:9c69:573b -p 80 -v
```

Connect to URL destination server:

```bash
./go-tormask -u https://www.google.com -p 80 -v
```
