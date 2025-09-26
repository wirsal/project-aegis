package service

import (
	"context"
	"log"
	"strconv"

	"golang.org/x/text/encoding/charmap" // Diperlukan untuk EBCDIC -> ASCII
	"google.golang.org/grpc"

	pb "github.com/wirsal/project-aegis/api/protos"
	"github.com/wirsal/project-aegis/pkg/codec"
)

// Definisi interface RuleEngineClient tetap sama
type RuleEngineClient interface {
	AnalyzeTransaction(ctx context.Context, in *pb.Transaction, opts ...grpc.CallOption) (*pb.RiskResult, error)
}

// Struct GatewayService tetap sama
type GatewayService struct {
	ruleEngineClient RuleEngineClient
}

// Fungsi NewGatewayService tetap sama
func NewGatewayService(client RuleEngineClient) *GatewayService {
	return &GatewayService{
		ruleEngineClient: client,
	}
}

// ProcessAndForwardMessage sekarang berisi semua logika.
func (s *GatewayService) ProcessAndForwardMessage(ctx context.Context, rawMessage []byte) error {
	log.Println("Memulai translasi dan parsing pesan kustom...")

	// 1. Lakukan translasi dari EBCDIC ke ASCII
	decoder := charmap.CodePage1047.NewDecoder()
	asciiBytes, err := decoder.Bytes(rawMessage)
	if err != nil {
		log.Printf("ERROR: Gagal melakukan translasi EBCDIC ke ASCII: %v", err)
		return err
	}

	log.Println("Pesan berhasil ditranslasi ke ASCII.")

	// 2. Panggil fungsi ReadTransaction (sekarang menjadi parser) dengan data ASCII
	protoMsg, err := s.ParseAndMapTransaction(string(asciiBytes))
	if err != nil {
		log.Printf("ERROR: Gagal mem-parsing atau memetakan transaksi: %v", err)
		return err
	}

	log.Printf("✅ Pesan berhasil di-parsing. Mengirim transaksi (Reff: %s)...", protoMsg.TrxReffNumber)
	riskResult, err := s.ruleEngineClient.AnalyzeTransaction(ctx, protoMsg)
	if err != nil {
		log.Printf("ERROR: Panggilan gRPC ke Rule Engine gagal: %v", err)
		return err
	}

	log.Printf("✅ Respon dari Rule Engine diterima. Skor Risiko: %d, Level: %s", riskResult.RiskScore, riskResult.RiskLevel)
	return nil
}

// ReadTransaction diubah namanya menjadi ParseAndMapTransaction untuk kejelasan
func (s *GatewayService) ParseAndMapTransaction(data string) (*pb.Transaction, error) {
	println("data, ", data)
	log.Println("Mem-parsing dan memetakan data transaksi dari string ASCII...")

	// Helper untuk parseFloat yang aman
	parseFloat := func(s string) float32 {
		val, _ := strconv.ParseFloat(s, 32)
		return float32(val)
	}

	// Catatan: Fungsi codec.Hex2* Anda sekarang akan bekerja pada string ASCII.
	// Pastikan fungsi tersebut sesuai. Jika fungsi codec Anda mengharapkan
	// string HEX, maka data ASCII perlu diubah ke HEX terlebih dahulu.
	// Namun, kita asumsikan di sini data sudah siap pakai.

	// Mapping langsung dari data parsing ke struct Protobuf
	trx := &pb.Transaction{
		TrxDate:          codec.Hex2string_comp3(data[0:4])[0:7],
		TrxTime:          codec.Hex2string_comp3(data[4:8])[:7],
		CardOrg:          codec.Hex2string_comp3(data[8:10])[:3],
		CardType:         codec.Hex2string_comp3(data[10:12])[:3],
		CardNumber:       codec.Hex2string_comp3(data[12:21])[1:17],
		MerchOrg:         codec.Hex2string_comp3(data[44:46])[:3],
		MerchNumber:      codec.Hex2string_comp3(data[46:51])[:9],
		TrxCode:          codec.Hex2string_comp3(data[55:57])[1:3],
		TrxReffNumber:    codec.Hex2string_comp3(data[57:64])[:13],
		TrxAmount:        parseFloat(codec.Hex2string_comp3(data[70:77])[:13]),
		CardExpired:      codec.Hex2string(data[92:96]),
		MerchCategory:    codec.Hex2string_comp3(data[96:99])[1:5],
		TrxCountry:       codec.Hex2string_comp3(data[99:101])[:3],
		TrxAuthCode:      codec.Hex2string(data[117:123]),
		TrxCardType:      codec.Hex2string_comp3(data[155:156])[1:],
		TrxCurrency:      codec.Hex2string_comp3(data[149:151])[:3],
		TrxInstallment:   isInstallment(data[275:276]),
		TrxDeclineReason: declineReason(data[171:173]),
		TrxTerminalId:    codec.Hex2string(data[126:134]),
		TrxMerchantId:    codec.Hex2string(data[134:149]),
		TrxAcqId:         codec.Hex2string_comp3(data[105:111])[:11],
		TrxFwdId:         codec.Hex2string_comp3(data[111:116])[:10],
		TrxChbCurr:       codec.Hex2string_comp3(data[151:153])[:3],
	}

	// Logika switch untuk field-field kondisional
	switch trx.TrxCardType {
	case "1", "3":
		trx.TrxRespCode = codec.Hex2string(data[619:621])
		trx.TrxMerchantName = codec.Hex2string(data[621:661])
		trx.TrxPinCap = codec.Hex2string(data[103:104])
		trx.TrxPosMode = codec.Hex2string(data[101:103])
		trx.TrxPosData = codec.Hex2string_comp3(data[677:682])
		trx.TrxStip = codec.Hex2string(data[661:662])
		trx.TrxCavResult = codec.Hex2string(data[694:695])
		trx.TrxArqcResult = codec.Hex2string(data[693:694])
	case "2", "4":

	case "6", "7":
	}

	log.Printf("Data Transaksi yang Diparsing: %+v", trx)
	return trx, nil
}

func isInstallment(instlInd string) string {
	instl := codec.Hex2string(instlInd)
	switch instl {
	case "I", "Y":
		return "Y"
	default:
		return "N"
	}
}

func declineReason(decCode string) string {
	val := codec.Hex2string_comp3(decCode)
	sign := val[len(val)-1:]
	code, _ := strconv.Atoi(val[:len(val)-1])

	switch sign {
	case "C", "c":
		return strconv.Itoa(code * -1)
	default:
		return strconv.Itoa(code)
	}
}
