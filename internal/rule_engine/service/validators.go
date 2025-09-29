package service

import (
	"strconv"
	"strings"
)

// vInList checks if an input exists in a semicolon-separated list.
func vInList(ruleValue, inputValue, wildcard string) bool {
	// 1. Trim excess spaces and semicolons from the start/end
	cleanRule := strings.Trim(ruleValue, " ;")

	if cleanRule == wildcard {
		return true
	}

	// 2. If the rule becomes empty after cleaning, consider it a match (no rule)
	// This handles cases like ruleValue = ";" or ruleValue = ""
	if cleanRule == "" {
		return true
	}

	list := strings.Split(cleanRule, ";")
	for _, item := range list {
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
