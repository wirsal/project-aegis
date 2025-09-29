package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/wirsal/project-aegis/api/protos"
	"github.com/wirsal/project-aegis/internal/rule_engine/service"
)

const (
	port = ":50051"
)

func main() {
	// 1. Buka listener di port gRPC
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("FATAL: Gagal membuka listener di port %s: %v", port, err)
	}

	// 2. Buat instance server gRPC baru
	grpcServer := grpc.NewServer()

	// 3. Buat instance dari implementasi service kita
	// PERUBAHAN: Masukkan path ke file "rules.json" sebagai argumen.
	ruleEngineServer, err := service.NewRuleEngineServer("./configs/rules.json")
	if err != nil {
		log.Fatalf("FATAL: Gagal inisialisasi Rule Engine Service: %v", err)
	}

	// 4. Daftarkan service kita ke server gRPC
	pb.RegisterRuleEngineServer(grpcServer, ruleEngineServer)

	log.Printf("🚀 Rule Engine Service berjalan dan mendengarkan di port %s", port)

	// 5. Mulai melayani permintaan
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("FATAL: Gagal menjalankan server gRPC: %v", err)
	}
}
