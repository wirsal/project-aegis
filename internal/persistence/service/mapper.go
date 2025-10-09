package service

import (
	"fmt"
	"log"
	"strconv"
	"time"

	pb "github.com/wirsal/project-aegis/api/protos"
)

const insertRawQuery = `
    INSERT INTO transaction_log (
        trx_key, card_org, card_type, card_number, card_expdate, trx_date, trx_time, trx_datetime,
        merch_org, merch_id, trx_cardtype, trx_code, trx_respcode, trx_declinereason,
        trx_reffnumber, trx_amt, trx_billamt, trx_orgamt, trx_convrate, trx_currency,
        trx_chbcurr, trx_merchant, trx_merchname, trx_acqid, trx_fwdid, trx_mcc, trx_countrycode,
        trx_authcode, trx_terminal, trx_pincap, trx_posmode, trx_posdata, trx_installment,
        trx_stip, trx_cvv_result, trx_cvv2_result, trx_cavv_result, trx_arqc_result,
        trx_chip_length, trx_chip_data
    ) VALUES (
        :trx_key, :card_org, :card_type, :card_number, :card_expdate, :trx_date, :trx_time, :trx_datetime,
        :merch_org, :merch_id, :trx_cardtype, :trx_code, :trx_respcode, :trx_declinereason,
        :trx_reffnumber, :trx_amt, :trx_billamt, :trx_orgamt, :trx_convrate, :trx_currency,
        :trx_chbcurr, :trx_merchant, :trx_merchname, :trx_acqid, :trx_fwdid, :trx_mcc, :trx_countrycode,
        :trx_authcode, :trx_terminal, :trx_pincap, :trx_posmode, :trx_posdata, :trx_installment,
        :trx_stip, :trx_cvv_result, :trx_cvv2_result, :trx_cavv_result, :trx_arqc_result,
        :trx_chip_length, :trx_chip_data
    )`
const insertRiskResultQuery = `
    INSERT INTO risk_results (
        rr_key, rr_card, rr_desc, rr_desc_add1, rr_desc_add2, rr_desc_add3,
        rr_curr_code, rr_amount, rr_amount_add1, rr_amount_add2, rr_datetime,
        rr_date_add1, rr_date_add2, rr_rule_code, rr_rule_type, 
        rr_date_proc, rr_date_vald, rr_date_write
    ) VALUES (
        :rr_key, :rr_card, :rr_desc, :rr_desc_add1, :rr_desc_add2, :rr_desc_add3,
        :rr_curr_code, :rr_amount, :rr_amount_add1, :rr_amount_add2, :rr_datetime,
        :rr_date_add1, :rr_date_add2, :rr_rule_code, :rr_rule_type, 
        :rr_date_proc, :rr_date_vald, :rr_date_write
    )`

func mapToTransactionLogModel(trx *pb.Transaction) (*TransactionLogModel, error) {
	trxDateTime, err := time.Parse("2006-01-02 15:04:05", trx.TrxDate+" "+trx.TrxTime)
	if err != nil {
		log.Printf("Warning: could not parse trx_datetime for TrxKey %s, using current time. Error: %v", trx.TrxKey, err)
		trxDateTime = time.Now()
	}

	model := &TransactionLogModel{
		TrxKey:           trx.TrxKey,
		CardOrg:          trx.CardOrg,
		CardType:         trx.CardType,
		CardNumber:       trx.CardNumber,
		CardExpDate:      trx.CardExpired,
		TrxDate:          trx.TrxDate,
		TrxTime:          trx.TrxTime,
		TrxDatetime:      trxDateTime,
		MerchOrg:         trx.MerchOrg,
		MerchID:          trx.MerchNumber,
		TrxCardType:      trx.TrxCardType,
		TrxCode:          safeAtoi(trx.TrxCode),
		TrxRespCode:      trx.TrxRespCode,
		TrxDeclineReason: safeAtoi(trx.TrxDeclineReason),
		TrxReffNumber:    trx.TrxReffNumber,
		TrxAmt:           trx.TrxAmount,
		TrxBillAmt:       safeParseFloat(trx.TrxBillAmount),
		TrxOrgAmt:        trx.TrxOrgAmount,
		TrxConvRate:      0.0, // Tetap default 0.0
		TrxCurrency:      safeAtoi(trx.TrxCurrency),
		TrxChbCurr:       safeAtoi(trx.TrxChbCurr),
		TrxMerchant:      trx.TrxMerchantId,
		TrxMerchName:     trx.TrxMerchantName,
		TrxAcqID:         trx.TrxAcqId,
		TrxFwdID:         trx.TrxFwdId,
		TrxMCC:           safeAtoi(trx.MerchCategory),
		TrxCountryCode:   safeAtoi(trx.TrxCountry),
		TrxAuthCode:      trx.TrxAuthCode,
		TrxTerminal:      trx.TrxTerminalId,
		TrxPinCap:        safeAtoi(trx.TrxPinCap),
		TrxPosMode:       safeAtoi(trx.TrxPosMode),
		TrxPosData:       trx.TrxPosData,
		TrxInstallment:   trx.TrxInstallment,
		TrxStip:          trx.TrxStip,
		TrxCvvResult:     trx.TrxCvvResult,
		TrxCvv2Result:    trx.TrxCvv2Result,
		TrxCavvResult:    trx.TrxCavResult,
		TrxArqcResult:    trx.TrxArqcResult,
		TrxChipLength:    len(trx.TrxChipData),
		TrxChipData:      trx.TrxChipData,
	}
	return model, nil
}

func mapToRiskResultModel(req *pb.StoreTransactionRequest) (*RiskResultModel, error) {
	trxData := req.GetTransactionData()
	riskData := req.GetRiskData()

	trxDateTime, err := time.Parse("2006-01-02 15:04:05", trxData.TrxDate+" "+trxData.TrxTime)
	if err != nil {
		log.Printf("Warning: could not parse rr_datetime for TrxKey %s, using current time. Error: %v", trxData.TrxKey, err)
		trxDateTime = time.Now()
	}

	model := &RiskResultModel{
		Key:       trxData.TrxKey,
		Card:      trxData.CardNumber,
		Desc:      trxData.TrxMerchantName,
		CurrCode:  trxData.TrxCurrency,
		Amount:    fmt.Sprintf("%.2f", trxData.TrxAmount),
		DateTime:  trxDateTime,
		RuleCode:  riskData.RuleCode,
		RuleType:  riskData.RuleType,
		DateProc:  time.Now(),
		DateVald:  time.Now(),
		DateWrite: time.Now(),
	}
	return model, nil
}

func safeAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func safeParseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return f
}
