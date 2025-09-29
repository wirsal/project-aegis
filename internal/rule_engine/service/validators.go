package service

import (
	"strconv"
	"strings"
)

// vInList memeriksa apakah input ada di dalam daftar yang dipisahkan semicolon.
// Ini bisa digunakan untuk vCardOrg, vCardType, vMcc, dll.
func vInList(ruleValue, inputValue, wildcard string) bool {
	if ruleValue == wildcard || ruleValue == wildcard+";" {
		return true
	}
	list := strings.Split(ruleValue, ";")
	for _, item := range list {
		if item == inputValue {
			return true
		}
	}
	return false
}

// vInRange memeriksa apakah input berada dalam rentang min-max.
// Digunakan untuk vCrLimit, vTrxAmt, vTimeStamp.
func vInRange(ruleValue string, inputValue int64) bool {
	parts := strings.Split(ruleValue, "-")
	if len(parts) != 2 {
		return false // Format aturan salah
	}
	min, _ := strconv.ParseInt(parts[0], 10, 64)
	max, _ := strconv.ParseInt(parts[1], 10, 64)

	return inputValue >= min && inputValue <= max
}

// vCountry mengimplementasikan logika include/exclude untuk negara.
func vCountry(ruleValue, inputValue string) bool {
	if ruleValue == "A000" {
		return true
	}
	indicator := ruleValue[0:1]
	listStr := ruleValue[1:]
	list := strings.Split(listStr, ";")

	found := false
	for _, item := range list {
		if item == inputValue {
			found = true
			break
		}
	}

	// Jika 'I' (Include), harus ditemukan. Jika 'E' (Exclude), tidak boleh ditemukan.
	return (indicator == "I" && found) || (indicator == "E" && !found)
}
