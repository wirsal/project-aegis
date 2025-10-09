package service

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	pb "github.com/wirsal/project-aegis/api/protos"
)

func TestStoreRawTransaction(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected...", err)
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	s := &Service{db: sqlxDB}

	sampleTrx := &pb.Transaction{
		TrxKey:        "test-key-123",
		CardOrg:       "001",
		CardType:      "002",
		CardNumber:    "4000123456789010",
		CardExpired:   "1228",
		TrxDate:       "2025-10-09",
		TrxTime:       "10:30:00",
		MerchOrg:      "003",
		MerchNumber:   "merchant001",
		TrxCardType:   "C",
		TrxCode:       "01",
		TrxRespCode:   "00",
		TrxAmount:     15000.50,
		TrxBillAmount: "15000.50",
		MerchCategory: "5411",
		TrxCountry:    "360",
		TrxChipData:   "somechipdata",
	}

	model, _ := mapToTransactionLogModel(sampleTrx)

	reboundQuery := sqlx.Rebind(sqlx.QUESTION, insertRawQuery)

	mock.ExpectExec(regexp.QuoteMeta(reboundQuery)).
		WithArgs(
			model.TrxKey, model.CardOrg, model.CardType, model.CardNumber, model.CardExpDate, model.TrxDate,
			model.TrxTime, sqlmock.AnyArg(),
			model.MerchOrg, model.MerchID, model.TrxCardType,
			model.TrxCode, model.TrxRespCode, model.TrxDeclineReason, model.TrxReffNumber, model.TrxAmt,
			model.TrxBillAmt, model.TrxOrgAmt, model.TrxConvRate, model.TrxCurrency, model.TrxChbCurr,
			model.TrxMerchant, model.TrxMerchName, model.TrxAcqID, model.TrxFwdID, model.TrxMCC,
			model.TrxCountryCode, model.TrxAuthCode, model.TrxTerminal, model.TrxPinCap, model.TrxPosMode,
			model.TrxPosData, model.TrxInstallment, model.TrxStip, model.TrxCvvResult, model.TrxCvv2Result,
			model.TrxCavvResult, model.TrxArqcResult, model.TrxChipLength, model.TrxChipData,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
		// -------------------------

	err = s.StoreRawTransaction(context.Background(), sampleTrx)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
