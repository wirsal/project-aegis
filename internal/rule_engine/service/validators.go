package service

import (
	"strconv"
	"strings"

	pb "github.com/wirsal/project-aegis/api/protos"
)

// validateRule tetap sama
func (s *Service) validateRule(rule Rule, trx *pb.Transaction) bool {
	trxTime, _ := strconv.Atoi(trx.TrxTime)
	return vInList(rule.Org, trx.CardOrg, "000") &&
		vInList(rule.Type, trx.CardType, "000") &&
		vInList(rule.MerchCategory, trx.MerchCategory, "0000") &&
		// vInList(rule.TransCode, trx.tra)
		vInclusionExclusion(rule.CountryCode, trx.TrxCountry, "A000") &&
		vInclusionExclusion(rule.CurrencyCode, trx.TrxCurrency, "A000") &&
		vInRange(rule.Amount, int64(trx.TrxAmount)) &&
		vInList(rule.PosCondCode, trx.TrxPosMode, "AA") &&
		vInList(rule.RespCode, trx.TrxRespCode, "AA") &&
		vInRange(rule.TimeStamp, int64(trxTime)) &&
		vInList(rule.InstallmentInd, trx.TrxInstallment, "-")
}

// vInList checks if an input exists in a semicolon-separated list.
func vInList(ruleValue, inputValue, wildcard string) bool {
	cleanRule := strings.Trim(ruleValue, " ;")

	if cleanRule == wildcard {
		return true
	}

	// 2. If the rule becomes empty after cleaning, consider it a match (no rule)
	if cleanRule == "" {
		return true
	}

	list := strings.SplitSeq(cleanRule, ";")
	for item := range list {
		// Trim spaces per item to handle cases like "001; 002"
		if strings.TrimSpace(item) == inputValue {
			return true
		}
	}
	return false
}

// vInRange checks if an input is within a min-max range.
// Used for vCrLimit, vTrxAmt, vTimeStamp.
func vInRange(ruleValue string, inputValue int64) bool {
	parts := strings.Split(ruleValue, "-")
	if len(parts) != 2 {
		return false // Incorrect rule format
	}
	min, _ := strconv.ParseInt(parts[0], 10, 64)
	max, _ := strconv.ParseInt(parts[1], 10, 64)
	return inputValue >= min && inputValue <= max
}

// vInclusionExclusion implements include/exclude logic.
// It checks a rule value like "I360;840" or "E360".
func vInclusionExclusion(ruleValue, inputValue, wildcard string) bool {
	// 1. Handle wildcard case first
	if ruleValue == wildcard {
		return true
	}

	// 2. Separate the indicator ('I' or 'E') from the list of values
	indicator := ruleValue[0:1]
	listStr := ruleValue[1:]
	list := strings.Split(listStr, ";")

	// 3. Search if the input value is in the list
	found := false
	for _, item := range list {
		if strings.TrimSpace(item) == inputValue {
			found = true
			break
		}
	}

	// 4. Apply the final logic
	// If 'I' (Include), it must be found.
	// If 'E' (Exclude), it must NOT be found.
	switch indicator {
	case "I":
		return found
	case "E":
		return !found
	}

	// Default to false if the indicator is not 'I' or 'E'
	return false
}
