package parser

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	pb "github.com/wirsal/project-aegis/api/protos"
	"github.com/wirsal/project-aegis/pkg/codec"
)

// ParseAndMapTransaction takes a raw data string and maps it to a Transaction protobuf struct.
func ParseAndMapTransaction(data string) (*pb.Transaction, error) {
	log.Println("Parsing and mapping transaction data from ASCII string...")

	parseFloat := func(s string) float32 {
		val, _ := strconv.ParseFloat(s, 32)
		return float32(val)
	}

	trx := &pb.Transaction{

		TrxDate:          julianToDate(codec.Hex2string_comp3(data[0:4])[0:7]),
		TrxTime:          strToTime(codec.Hex2string_comp3(data[4:8])[:7]),
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

	switch trx.TrxCardType {
	case "1", "3": //VISA
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
	case "2", "4": //Mastercard
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
	case "6", "7":

	}

	trx.TrxKey = generateDeterministicTrxKey(trx.CardNumber, trx.TrxDate, trx.TrxTime, trx.TrxAuthCode, trx.TrxAmount)
	log.Printf("Parsed Transaction Data: %+v", trx)
	return trx, nil
}

// All helper functions are now neatly contained in the parser package.
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

func generateDeterministicTrxKey(cardNumber, trxDate, trxTime, trxAuth string, amount float32) string {
	// 1. Gabungkan semua field penting menjadi satu string yang stabil dan unik.
	//    Gunakan pemisah agar tidak ambigu (misal, "123"+"45" vs "12"+"345").
	inputString := fmt.Sprintf("%s-%s-%s-%s-%.2f", cardNumber, trxDate, trxTime, trxAuth, amount)

	// 2. Buat hash dari string tersebut menggunakan SHA-256
	hasher := sha256.New()
	hasher.Write([]byte(inputString))
	hashBytes := hasher.Sum(nil)

	// 3. Konversi hash menjadi string heksadesimal.
	//    SHA-256 menghasilkan 64 karakter, jadi kita potong agar sesuai batasan 30 karakter.
	return hex.EncodeToString(hashBytes)[:30]
}

func julianToDate(julianStr string) string {
	year, err := strconv.Atoi(julianStr[0:4])
	if err != nil {
		log.Printf("ERROR: Gagal parse tahun dari '%s': %v", julianStr, err)
		return ""
	}

	dayOfYear, err := strconv.Atoi(julianStr[4:])
	if err != nil {
		log.Printf("ERROR: Gagal parse hari dari '%s': %v", julianStr, err)
		return ""
	}

	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	gregorianDate := startDate.AddDate(0, 0, dayOfYear-1)

	return gregorianDate.Format("2006-01-02")
}

func strToTime(inputTime string) string {
	timeToParse := inputTime[1:]
	inputLayout := "150405"

	parsedTime, err := time.Parse(inputLayout, timeToParse)
	if err != nil {
		log.Fatalf("Gagal parse waktu: %v", err)
	}
	outputLayout := "15:04:05"

	return parsedTime.Format(outputLayout)
}
