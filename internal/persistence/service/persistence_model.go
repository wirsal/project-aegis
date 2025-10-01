package service

import "time"

type TransactionLogModel struct {
	TrxID            string    `db:"trx_id"`
	CardOrg          string    `db:"card_org"`
	CardType         string    `db:"card_type"`
	CardNumber       string    `db:"card_number"`
	CardExpDate      string    `db:"card_expdate"`
	TrxDate          string    `db:"trx_date"`
	TrxTime          string    `db:"trx_time"`
	TrxDatetime      time.Time `db:"trx_datetime"`
	MerchOrg         string    `db:"merch_org"`
	MerchID          string    `db:"merch_id"`
	TrxCardType      string    `db:"trx_cardtype"`
	TrxCode          int       `db:"trx_code"`
	TrxRespCode      string    `db:"trx_respcode"`
	TrxDeclineReason int       `db:"trx_declinereason"`
	TrxReffNumber    string    `db:"trx_reffnumber"`
	TrxAmt           float32   `db:"trx_amt"`
	TrxBillAmt       float64   `db:"trx_billamt"`
	TrxOrgAmt        float32   `db:"trx_orgamt"`
	TrxConvRate      float64   `db:"trx_convrate"`
	TrxCurrency      int       `db:"trx_currency"`
	TrxChbCurr       int       `db:"trx_chbcurr"`
	TrxMerchant      string    `db:"trx_merchant"`
	TrxMerchName     string    `db:"trx_merchname"`
	TrxAcqID         string    `db:"trx_acqid"`
	TrxFwdID         string    `db:"trx_fwdid"`
	TrxMCC           int       `db:"trx_mcc"`
	TrxCountryCode   int       `db:"trx_countrycode"`
	TrxAuthCode      string    `db:"trx_authcode"`
	TrxTerminal      string    `db:"trx_terminal"`
	TrxPinCap        int       `db:"trx_pincap"`
	TrxPosMode       int       `db:"trx_posmode"`
	TrxPosData       string    `db:"trx_posdata"`
	TrxInstallment   string    `db:"trx_installment"`
	TrxStip          string    `db:"trx_stip"`
	TrxCvvResult     string    `db:"trx_cvv_result"`
	TrxCvv2Result    string    `db:"trx_cvv2_result"`
	TrxCavvResult    string    `db:"trx_cavv_result"`
	TrxArqcResult    string    `db:"trx_arqc_result"`
	TrxChipLength    int       `db:"trx_chip_length"`
	TrxChipData      string    `db:"trx_chip_data"`
}
