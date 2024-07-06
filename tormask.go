package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

// Constants for proxy connection
const (
	ProxyIP   = "127.0.0.1"
	ProxyPort = 9050
	User      = "tormask"
)

// SOCKS4 request and response structures
type ProxyRequest struct {
	Vn      byte
	Cd      byte
	DstPort uint16
	DstIP   [4]byte // Use [4]byte for IPv4, [16]byte for IPv6
	UserID  [8]byte
}

type ProxyResponse struct {
	Vn byte
	Cd byte
	_  [6]byte
}

// ResolveHost resolves a hostname to an IP address
func ResolveHost(host string) (string, error) {
	addrs, err := net.LookupHost(host)
	if err != nil {
		return "", err
	}

	// Prefer IPv4 over IPv6 for SOCKS4 protocol
	for _, addr := range addrs {
		if strings.Contains(addr, ":") { // IPv6 addresses contain colons
			continue
		}
		return addr, nil
	}

	return "", fmt.Errorf("no IPv4 address found for %s", host)
}

// CreateRequest creates a SOCKS4 request
func CreateRequest(dstip string, dstport int) (*ProxyRequest, error) {
	ip := net.ParseIP(dstip)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address")
	}

	var dstIP [4]byte
	if ip.To4() != nil {
		copy(dstIP[:], ip.To4())
	} else if ip.To16() != nil {
		copy(dstIP[:], ip.To16())
	} else {
		return nil, fmt.Errorf("unsupported IP address format")
	}

	req := &ProxyRequest{
		Vn:      4,
		Cd:      1,
		DstPort: uint16(dstport),
		DstIP:   dstIP,
	}
	copy(req.UserID[:], User)

	return req, nil
}

func main() {
	// Command-line argument parsing
	url := flag.String("u", "", "URL to connect to")
	ip := flag.String("i", "", "IP address to connect to")
	port := flag.Int("p", 0, "Port to connect to")
	verbose := flag.Bool("v", false, "Verbose output")
	flag.Parse()

	// Validate arguments
	if (*url == "" && *ip == "") || *port == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-u url | -i ip] -p port [-v]\n", os.Args[0])
		os.Exit(1)
	}

	host := *url
	if *ip != "" {
		host = *ip
	}

	// Resolve host if necessary
	if net.ParseIP(host) == nil {
		resolvedIP, err := ResolveHost(host)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error resolving host: %v\n", err)
			os.Exit(1)
		}
		host = resolvedIP
	}

	// Create SOCKS4 request
	req, err := CreateRequest(host, *port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating request: %v\n", err)
		os.Exit(1)
	}

	// Connect to the proxy server
	proxyAddr := fmt.Sprintf("%s:%d", ProxyIP, ProxyPort)
	conn, err := net.DialTimeout("tcp", proxyAddr, 10*time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to proxy: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	if *verbose {
		fmt.Println("Connected to proxy...")
	}

	// Send SOCKS4 request
	if err := binary.Write(conn, binary.BigEndian, req); err != nil {
		fmt.Fprintf(os.Stderr, "Error sending request: %v\n", err)
		os.Exit(1)
	}

	// Read SOCKS4 response
	var resp ProxyResponse
	if err := binary.Read(conn, binary.BigEndian, &resp); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading response: %v\n", err)
		os.Exit(1)
	}

	// Check if connection through proxy was successful
	if resp.Cd != 90 {
		fmt.Fprintf(os.Stderr, "Failed to connect: %d\n", resp.Cd)
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("Connected through proxy to %s:%d\n", host, *port)
	}

	// Send HTTP HEAD request
	httpReq := fmt.Sprintf("HEAD / HTTP/1.1\r\nHost: %s\r\n\r\n", host)
	if _, err := conn.Write([]byte(httpReq)); err != nil {
		fmt.Fprintf(os.Stderr, "Error sending HTTP request: %v\n", err)
		os.Exit(1)
	}

	// Read and print HTTP response
	respBuf := make([]byte, 4096)
	n, err := conn.Read(respBuf)
	if err != nil && err != io.EOF {
		fmt.Fprintf(os.Stderr, "Error reading HTTP response: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(respBuf[:n]))
}
