package handler

import (
	"context"
	"encoding/binary"
	"io" // Diperlukan untuk io.ReadAll
	"log"
	"net"
)

// GatewayServiceDef defines the interface that must be implemented by the service layer.
type GatewayServiceDef interface {
	ProcessAndForwardMessage(ctx context.Context, rawMessage []byte) error
}

type TCPHandler struct {
	service GatewayServiceDef
}

func NewTCPHandler(svc GatewayServiceDef) *TCPHandler {
	return &TCPHandler{
		service: svc,
	}
}

// HandleConnection to handle one TCP connection.
func (h *TCPHandler) HandleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("Handling new connection from: %s", conn.RemoteAddr())

	for {
		// 1. Read 2-byte length header
		header := make([]byte, 2)
		if _, err := io.ReadFull(conn, header); err != nil {
			if err == io.EOF {
				log.Printf("Connection closed by client: %s", conn.RemoteAddr())
			} else {
				log.Printf("ERROR: Failed to read header from %s:", err)
			}
			break
		}
		length := int(binary.BigEndian.Uint16(header))

		if length <= 0 {
			log.Printf("Received header, expected message length: %d bytes", length)
		} else {
			// 2. Read the entire message body based on the length
			rawBody := make([]byte, length)
			if _, err := io.ReadFull(conn, rawBody); err != nil {
				log.Printf("ERROR: Failed to read message body: %v", err)
				break
			}

			ctx := context.Background()
			if err := h.service.ProcessAndForwardMessage(ctx, rawBody); err != nil {
				log.Printf("ERROR: Failed to process message: %v", err)
			}
		}

	}
}

// StartServer starts the TCP listener and accepts connections.
func (h *TCPHandler) StartServer(port string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("FATAL: Failed to start TCP server on port %s: %v", port, err)
	}
	defer listener.Close()
	log.Printf("🚀 Gateway Service is running on TCP port %s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("ERROR: Failed to accept connection: %v", err)
			continue
		}
		// Each new connection will be handled by HandleConnection.
		go h.HandleConnection(conn)
	}
}
