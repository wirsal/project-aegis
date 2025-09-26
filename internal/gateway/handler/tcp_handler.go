package handler

import (
	"context"
	"io" // Diperlukan untuk io.ReadAll
	"log"
	"net"
)

// GatewayServiceDef mendefinisikan interface yang harus dipenuhi oleh service layer.
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

// HandleConnection sekarang jauh lebih sederhana.
func (h *TCPHandler) HandleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("Menangani koneksi baru dari: %s", conn.RemoteAddr())

	// Tidak ada lagi 'for' loop, karena kita berasumsi 1 koneksi = 1 pesan.

	// 1. Baca SEMUA data dari koneksi sampai koneksi ditutup oleh klien.
	rawBody, err := io.ReadAll(conn)
	if err != nil {
		log.Printf("ERROR: Gagal membaca data dari koneksi: %v", err)
		return // Keluar dari fungsi karena koneksi bermasalah.
	}

	// Cek apakah ada data yang diterima sebelum melanjutkan
	if len(rawBody) == 0 {
		log.Printf("Koneksi ditutup tanpa menerima data.")
		return
	}

	log.Printf("Menerima total %d byte data.", len(rawBody))

	// 2. Teruskan seluruh body mentah ke service layer untuk diproses.
	ctx := context.Background()
	if err := h.service.ProcessAndForwardMessage(ctx, rawBody); err != nil {
		log.Printf("ERROR: Gagal memproses pesan: %v", err)
	}

	log.Printf("Selesai memproses pesan dari %s", conn.RemoteAddr())
}

// StartServer memulai TCP listener dan menerima koneksi.
func (h *TCPHandler) StartServer(port string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("FATAL: Gagal memulai TCP server di port %s: %v", port, err)
	}
	defer listener.Close()
	log.Printf("🚀 Gateway Service berjalan di port TCP %s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("ERROR: Gagal menerima koneksi: %v", err)
			continue
		}
		// Setiap koneksi baru akan ditangani oleh HandleConnection.
		go h.HandleConnection(conn)
	}
}
