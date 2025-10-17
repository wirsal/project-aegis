package service

type Rule struct {
	Priority       int    `json:"priority"`
	Status         int    `json:"status"`
	RuleCode       string `json:"rule_code"`
	RuleDesc       string `json:"rule_desc"`
	RuleType       string `json:"rule_type"`
	Org            string `json:"org"`
	Type           string `json:"type"`
	BlockCode      string `json:"block_code"`
	CrLimit        string `json:"cr_limit"`
	MerchCategory  string `json:"merch_category"`
	TransCode      string `json:"trans_code"`
	CountryCode    string `json:"country_code"`
	CurrencyCode   string `json:"currency_code"`
	Amount         string `json:"amount"`
	PosCondCode    string `json:"pos_cond_code"`
	RespCode       string `json:"resp_code"`
	TimeStamp      string `json:"time_stamp"`
	InstallmentInd string `json:"installment_ind"`
	FirstUsageFlag string `json:"first_usage_flag"`
	CardList       string `json:"card_list"`
	Channel        string `json:"channel"`
	Templates_id   string `json:"templates_id"`
}
