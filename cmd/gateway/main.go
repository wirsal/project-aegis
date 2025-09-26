package main

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Import paket handler dan service dari direktori internal
	"github.com/wirsal/project-aegis/internal/gateway/handler"
	"github.com/wirsal/project-aegis/internal/gateway/service"

	// Ganti dengan path .proto Anda yang sudah di-generate
	pb "github.com/wirsal/project-aegis/api/protos"
)

const (
	tcpPort               = ":3333"
	ruleEngineGRPCAddress = "localhost:50051" // Alamat gRPC Rule Engine Service
)

func main() {
	tcpHandler()
	kafkaHandler()

}

func tcpHandler() {
	// LANGKAH 1: Buat koneksi gRPC ke Rule Engine Service
	conn, err := grpc.Dial(ruleEngineGRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("FATAL: Gagal terhubung ke gRPC server: %v", err)
	}
	defer conn.Close()

	// Buat gRPC client dari koneksi yang sudah ada
	ruleEngineClient := pb.NewRuleEngineClient(conn)
	log.Println("Berhasil terhubung ke Rule Engine Service.")

	// LANGKAH 2: Inisialisasi service layer (si Koki)
	// Suntikkan (inject) gRPC client ke dalam service
	gatewayService := service.NewGatewayService(ruleEngineClient)

	// LANGKAH 3: Inisialisasi handler layer (si Pelayan)
	// Suntikkan (inject) service ke dalam handler
	tcpHandler := handler.NewTCPHandler(gatewayService)

	// LANGKAH 4: Jalankan server
	// Handler kini siap untuk menerima koneksi TCP dan meneruskannya ke service
	tcpHandler.StartServer(tcpPort)
}

func kafkaHandler() {
	// Implementasi Kafka handler di sini
}
