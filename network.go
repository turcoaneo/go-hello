package main

import (
	"fmt"
	"log/slog"
	"net"
)

// Listener wraps a TCP listener and its connection channel
type Listener struct {
	Port   int
	ConnCh chan net.Conn
}

// StartListener creates a TCP listener on the given port
func StartListener(port int) (*Listener, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	connCh := make(chan net.Conn)
	listener := &Listener{Port: port, ConnCh: connCh}

	// Accept loop
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				continue
			}
			connCh <- conn
		}
	}()

	return listener, nil
}

// WriteHTTPResponse sends a minimal HTTP response with a dynamic body
func WriteHTTPResponse(conn net.Conn, body string, logger *slog.Logger) {
	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		fmt.Sprintf("Content-Length: %d\r\n", len(body)) +
		"\r\n" +
		body

	_, err := conn.Write([]byte(response))
	if err != nil {
		logger.Error("Failed to send feedback", "error", err)
	}
	errConn := conn.Close()
	if errConn != nil {
		logger.Error("Failed to close connection", "error", errConn)
	}
}
