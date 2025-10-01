-- Table structure for table `transaction_log`
CREATE TABLE transaction_log (
    id SERIAL PRIMARY KEY,
    trx_key VARCHAR(30) NOT NULL UNIQUE,
    card_org VARCHAR(3) NOT NULL,
    card_type VARCHAR(3) NOT NULL,
    card_number VARCHAR(16) NOT NULL,
    card_expdate VARCHAR(4) NOT NULL,
    trx_date DATE NOT NULL,
    trx_time TIME NOT NULL,
    trx_datetime TIMESTAMP NOT NULL,
    merch_org VARCHAR(3) NOT NULL,
    merch_id VARCHAR(9) NOT NULL,
    trx_cardtype VARCHAR(1) NOT NULL,
    trx_code INT NOT NULL,
    trx_respcode VARCHAR(11) NOT NULL,
    trx_declinereason INT NOT NULL,
    trx_reffnumber VARCHAR(15) NOT NULL,
    trx_amt NUMERIC(12,2) NOT NULL,
    trx_billamt NUMERIC(12,2) NOT NULL,
    trx_orgamt NUMERIC(12,2) NOT NULL,
    trx_convrate NUMERIC(12,6) NOT NULL,
    trx_currency INT NOT NULL,
    trx_chbcurr INT NOT NULL,
    trx_merchant VARCHAR(15) NOT NULL,
    trx_merchname VARCHAR(50) NOT NULL,
    trx_acqid VARCHAR(15) NOT NULL,
    trx_fwdid VARCHAR(15) NOT NULL,
    trx_mcc INT NOT NULL,
    trx_countrycode INT NOT NULL,
    trx_authcode VARCHAR(6) NOT NULL,
    trx_terminal VARCHAR(10) NOT NULL,
    trx_pincap INT NOT NULL,
    trx_posmode INT NOT NULL,
    trx_posdata VARCHAR(30) NOT NULL,
    trx_installment VARCHAR(1) NOT NULL,
    trx_stip VARCHAR(1) NOT NULL,
    trx_cvv_result VARCHAR(1) NOT NULL,
    trx_cvv2_result VARCHAR(1) NOT NULL,
    trx_cavv_result VARCHAR(1) NOT NULL,
    trx_arqc_result VARCHAR(1) NOT NULL,
    trx_chip_length INT NOT NULL,
    trx_chip_data VARCHAR(1000) NOT NULL,
    date_add TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_trx_key ON transaction_log (trx_key);
CREATE INDEX idx_cardacct_authcode_date ON transaction_log (card_number, trx_authcode, trx_date);



--- Table structure for table `risk_results`
CREATE TABLE risk_results (
    rr_id BIGSERIAL PRIMARY KEY,
    rr_key VARCHAR(51) NOT NULL UNIQUE,
    rr_card VARCHAR(16) NOT NULL DEFAULT '',
    rr_desc VARCHAR(50) DEFAULT '',
    rr_desc_add1 VARCHAR(400),
    rr_desc_add2 VARCHAR(400),
    rr_desc_add3 VARCHAR(400),
    rr_curr_code VARCHAR(3) DEFAULT '',
    rr_amount VARCHAR(16),
    rr_amount_add1 VARCHAR(20),
    rr_amount_add2 VARCHAR(20),
    rr_datetime TIMESTAMP,
    rr_date_add1 DATE,
    rr_date_add2 DATE,
    rr_rule_code VARCHAR(10) NOT NULL,
    rr_type VARCHAR(10) NOT NULL,
    rr_date_proc TIMESTAMP NOT NULL,
    rr_date_vald TIMESTAMP NOT NULL,
    rr_date_write TIMESTAMP NOT NULL
);
-- Comment untuk kolom tinyint
COMMENT ON COLUMN risk_results.rr_notif_act IS '0 : Not Proses; 1 : email; 2 : HP; 3 : semuanya';

-- -- Indexes
-- CREATE INDEX idx_rr_datetime ON risk_results (rr_datetime);
-- CREATE INDEX idx_rr_notif_code ON risk_results (rr_notif_code);
-- CREATE INDEX idx_cek1 ON risk_results (rr_card, rr_datetime, rr_notif_code);
-- CREATE INDEX idx_cek2 ON risk_results (rr_datetime, rr_notif_code);
-- CREATE INDEX idx_fk_rr_otorisasi_input1 ON risk_results (rr_card, rr_datetime);
-- CREATE INDEX idx_type_card ON risk_results (rr_datetime, rr_key);
