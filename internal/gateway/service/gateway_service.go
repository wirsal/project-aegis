package service

import (
	"context"
	"log"
	"math"
	"strconv"

	// Diperlukan untuk EBCDIC -> ASCII
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

	protoMsg, err := s.ParseAndMapTransaction(string(rawMessage))
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
	log.Println("Mem-parsing dan memetakan data transaksi dari string ASCII...")

	// Helper untuk parseFloat yang aman
	parseFloat := func(s string) float32 {
		val, _ := strconv.ParseFloat(s, 32)
		return float32(val)
	}

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
	decx, _ := strconv.Atoi(codec.Hex2string(data[84:85]))
	tempConvRate := codec.Hex2string_comp3(data[85:92])
	trx.TrxOrgAmount = trxOrgAmountSimple(trx.TrxBillAmount, trx.TrxCardType, trx.MerchOrg, tempConvRate, decx)
	// Logika switch untuk field-field kondisional
	switch trx.TrxCardType {
	case "1", "3":
		//VISA

		trx.TrxPosMode = codec.Hex2string(data[101:103])
		trx.TrxPinCap = codec.Hex2string(data[103:104])
		trx.TrxRespCode = codec.Hex2string(data[619:621])
		trx.TrxMerchantName = codec.Hex2string(data[621:661])
		trx.TrxStip = codec.Hex2string(data[661:662])
		trx.TrxPosData = codec.Hex2string_comp3(data[677:682])
		trx.TrxCvvResult = codec.Hex2string(data[690:691])
		trx.TrxCvv2Result = codec.Hex2string(data[692:693])
		trx.TrxArqcResult = codec.Hex2string(data[693:694])
		trx.TrxCavResult = codec.Hex2string(data[694:695])
		//
	case "2", "4":
		//Mastercard
		trx.TrxPosMode = codec.Hex2string(data[102:104])
		trx.TrxPinCap = codec.Hex2string(data[104:105])
		trx.TrxRespCode = codec.Hex2string(data[657:659])
		trx.TrxMerchantName = codec.Hex2string(data[766:806])
		trx.TrxStip = codec.Hex2string(data[615:616])
		trx.TrxPosData = codec.Hex2string_comp3(data[631:657])
		trx.TrxCvvResult = codec.Hex2string(data[659:660])
		trx.TrxCvv2Result = codec.Hex2string(data[661:662])
		trx.TrxArqcResult = codec.Hex2string(data[664:665])
		trx.TrxCavResult = codec.Hex2string(data[665:666])
		//chipdata
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
func trxOrgAmountSimple(amountStr, trxCardType, merchOrg, tempConvRate string, decx int) float32 {
	// 1. Parse amount dengan presisi 32-bit. Hasilnya tetap float64 untuk kalkulasi.
	orgAmount, err := strconv.ParseFloat(amountStr, 32)
	if err != nil {
		return 0.0 // Tipe float32 akan diinferensikan secara otomatis
	}

	// 2. Guard Clause tetap sama
	isDefaultRate := ((trxCardType == "1" || trxCardType == "3") && tempConvRate == "9999999") ||
		merchOrg != "000" ||
		tempConvRate == "0000000"

	if isDefaultRate {
		// Konversi ke float32 saat mengembalikan nilai
		return float32(orgAmount)
	}

	// 3. Parse rate dengan presisi 32-bit.
	convRate, err := strconv.ParseFloat(tempConvRate, 32)
	if err != nil || convRate <= 0 {
		// Konversi ke float32 saat mengembalikan nilai
		return float32(orgAmount)
	}

	// 4. Kalkulasi tetap menggunakan float64 untuk presisi terbaik selama perhitungan
	divisor := math.Pow(10, float64(decx))
	finalRate := convRate / divisor

	if finalRate <= 0 {
		return float32(orgAmount)
	}

	// 5. Hitung hasil akhir sebagai float64, lalu konversi ke float32 saat return
	result := orgAmount / finalRate
	return float32(result)
}
