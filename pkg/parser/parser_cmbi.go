package parser

import (
	"log"
	"strconv"

	"github.com/wirsal/project-aegis/api/protos"
	"github.com/wirsal/project-aegis/pkg/codec"
)

func ParseTransaction(data string) *protos.Transaction {
	trx := &protos.Transaction{}
	trx.TrxDate = codec.Hex2string_comp3(data[0:4])[0:7]
	trx.TrxTime = codec.Hex2string_comp3(data[4:8])[:7]
	trx.CardOrg = codec.Hex2string_comp3(data[8:10])[:3]
	trx.CardType = codec.Hex2string_comp3(data[10:12])[:3]
	trx.CardNumber = codec.Hex2string_comp3(data[12:21])[1:17]
	trx.MerchOrg = codec.Hex2string_comp3(data[44:46])[:3]
	trx.MerchNumber = codec.Hex2string_comp3(data[46:51])[:9]
	trx.TrxCode = codec.Hex2string_comp3(data[55:57])[1:3]
	trx.TrxReffNumber = codec.Hex2string_comp3(data[57:64])[:13]

	return trx
}

func ReadTransaction(data string) {
	log.Println("Membaca transaksi dari pesan...")

	log.Printf("trxDate : %s", codec.Hex2string_comp3(data[0:4])[0:7])
	log.Printf("timeTrx : %s", codec.Hex2string_comp3(data[4:8])[:7])
	log.Printf("cardOrg : %s", codec.Hex2string_comp3(data[8:10])[:3])
	log.Printf("cardType : %s", codec.Hex2string_comp3(data[10:12])[:3])
	log.Printf("cardNumber : %s", codec.Hex2string_comp3(data[12:21])[1:17])
	log.Printf("merchOrg : %s", codec.Hex2string_comp3(data[44:46])[:3])
	merchOrg := codec.Hex2string_comp3(data[44:46])[:3]
	log.Printf("merchNumber : %s", codec.Hex2string_comp3(data[46:51])[:9])
	log.Printf("trxCode : %s", codec.Hex2string_comp3(data[55:57])[1:3])
	log.Printf("trxReffNum : %s", codec.Hex2string_comp3(data[57:64])[:13])
	log.Printf("trxAmt : %s", codec.Hex2string_comp3(data[70:77])[:13])
	log.Printf("trxBill : %s", "")
	log.Printf("carExpDate : %s", codec.Hex2string(data[92:96]))
	log.Printf("mcc : %s", codec.Hex2string_comp3(data[96:99])[1:5])
	log.Printf("trxCountry : %s", codec.Hex2string_comp3(data[99:101])[:3])
	log.Printf("trxAuthcode : %s", codec.Hex2string(data[117:123]))
	log.Printf("trx_cardtype : %s", codec.Hex2string_comp3(data[155:156])[1:])
	trx_cardtype := codec.Hex2string_comp3(data[155:156])[1:]

	switch trx_cardtype {
	case "1", "3":
		log.Printf("trx_respcode : %s", codec.Hex2string(data[619:621]))
		log.Printf("trx_merchname : %s", codec.Hex2string(data[621:661]))
		log.Printf("trx_pincap : %s", codec.Hex2string(data[103:104]))
		log.Printf("trx_posmode : %s", codec.Hex2string(data[101:103]))
		log.Printf("trx_posdata : %s", codec.Hex2string_comp3(data[677:682]))
		log.Printf("trx_stip : %s", codec.Hex2string(data[661:662]))
		// log.Printf("trx_cvv_result : %s", codec.Hex2string_comp3(data[690:691]))
		// log.Printf("trx_cvv2_result : %s", codec.Hex2string_comp3(data[692:693]))
		log.Printf("trx_cavv_result : %s", codec.Hex2string(data[694:695]))
		log.Printf("trx_arqc_result : %s", codec.Hex2string(data[693:694]))
	case "2", "4":
	case "6", "7":
	}
	log.Printf("trx_currency : %s", codec.Hex2string_comp3(data[149:151])[:3])
	log.Printf("trx_installment : %s", isInstallment(codec.Hex2string(data[275:276])))
	log.Printf("trx_declinereason : %s", declineReason(data[171:173]))
	log.Printf("trx_terminal : %s", codec.Hex2string(data[126:134]))
	log.Printf("trx_merchant : %s", codec.Hex2string(data[134:149]))
	log.Printf("tempConvRate : %s", codec.Hex2string(data[85:92]))
	tempConvRate := codec.Hex2string(data[85:92])
	decx, _ := strconv.Atoi(codec.Hex2string_comp3(data[84:85]))
	println("decx", decx)
	trxConvRate := 0
	if ((trx_cardtype == "1" || trx_cardtype == "3") && tempConvRate == "9999999") || merchOrg != "000" {

	} else {
		if tempConvRate != "0000000" {
			trxConvRate, _ := strconv.ParseFloat(tempConvRate, 64)
			trxConvRate = trxConvRate / 10
		} else {
			trxConvRate = 1
		}
	}
	trxOrgAmnt := "123"
	println("trxOrgAmnt", trxOrgAmnt)
	if trxConvRate > 0 {
		trxOrgAmnt = "2345"
	}

	log.Printf("trx_acqid : %s", codec.Hex2string_comp3(data[105:111])[:11])
	log.Printf("trx_fwdid : %s", codec.Hex2string_comp3(data[111:116])[:10])
	log.Printf("trx_chbcurr : %s", codec.Hex2string_comp3(data[151:153])[:3])

}

func isInstallment(instlInd string) string {
	switch instlInd {
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
