package service

import (
	"context"
	"log"
	"math"
	"strconv"

	"google.golang.org/grpc"

	"github.com/segmentio/ksuid"
	pb "github.com/wirsal/project-aegis/api/protos"
	"github.com/wirsal/project-aegis/pkg/codec"
)

type RuleEngineClient interface {
	AnalyzeTransaction(ctx context.Context, in *pb.Transaction, opts ...grpc.CallOption) (*pb.RiskResult, error)
}

type GatewayService struct {
	ruleEngineClient RuleEngineClient
}

func NewGatewayService(client RuleEngineClient) *GatewayService {
	return &GatewayService{
		ruleEngineClient: client,
	}
}

func (s *GatewayService) ProcessAndForwardMessage(ctx context.Context, rawMessage []byte) error {
	log.Println("Starting translation and parsing of custom message...")

	protoMsg, err := s.ParseAndMapTransaction(string(rawMessage))
	if err != nil {
		log.Printf("ERROR: Failed to parse or map the transaction: %v", err)
		return err
	}

	log.Printf("✅ Message parsed successfully. Sending transaction (Reff: %s)...", protoMsg.TrxReffNumber)
	riskResult, err := s.ruleEngineClient.AnalyzeTransaction(ctx, protoMsg)
	if err != nil {
		log.Printf("ERROR: gRPC call to Rule Engine failed: %v", err)
		return err
	}

	log.Printf("✅ Response from Rule Engine received. Risk Score: %d, Level: %s", riskResult.RiskScore, riskResult.RiskLevel)
	return nil
}

func (s *GatewayService) ParseAndMapTransaction(data string) (*pb.Transaction, error) {
	log.Println("Parsing and mapping transaction data from ASCII string...")

	// Helper for safe float parsing
	parseFloat := func(s string) float32 {
		val, _ := strconv.ParseFloat(s, 32)
		return float32(val)
	}

	// Direct mapping from parsed data to the Protobuf struct
	trx := &pb.Transaction{
		TrxId:            generateTrxId(),
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
	// Switch logic for conditional fields
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

	log.Printf("Parsed Transaction Data: %+v", trx)
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
	// 1. Parse the main amount with 32-bit precision. The result is still a float64 for calculations.
	orgAmount, err := strconv.ParseFloat(amountStr, 32)
	if err != nil {
		return 0.0 // The float32 type will be inferred automatically
	}

	// 2. Guard Clause remains the same
	isDefaultRate := ((trxCardType == "1" || trxCardType == "3") && tempConvRate == "9999999") ||
		merchOrg != "000" ||
		tempConvRate == "0000000"

	if isDefaultRate {
		// Convert to float32 upon returning the value
		return float32(orgAmount)
	}

	// 3. Parse the rate with 32-bit precision.
	convRate, err := strconv.ParseFloat(tempConvRate, 32)
	if err != nil || convRate <= 0 {
		// Convert to float32 upon returning the value
		return float32(orgAmount)
	}

	// 4. Calculations remain in float64 for best precision during the operation
	divisor := math.Pow(10, float64(decx))
	finalRate := convRate / divisor

	if finalRate <= 0 {
		return float32(orgAmount)
	}

	// 5. Calculate the final result as a float64, then convert to float32 on return
	result := orgAmount / finalRate
	return float32(result)
}

// Generate unique transaction ID using KSUID
func generateTrxId() string {
	id := ksuid.New()
	return id.String()
}
