package service

import "testing"

func TestVInList(t *testing.T) {
	// Definisikan semua kasus uji dalam sebuah "tabel"
	tests := []struct {
		name       string
		ruleValue  string
		inputValue string
		wildcard   string
		want       bool // Hasil yang diharapkan
	}{
		{"Simple Match", "001;002;003", "002", "000", true},
		{"No Match", "001;002;003", "004", "000", false},
		{"Wildcard Match", "000", "123", "000", true},
		{"Wildcard With Semicolon", "000;", "123", "000", true},
		{"Empty Rule String", "", "456", "000", true},
		{"Rule is Just Semicolon", ";", "456", "000", true},
		{"Trailing Semicolon Match", "001;002;", "002", "000", true},
		{"Item with Spaces", "001 ; 002", "002", "000", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := vInList(tt.ruleValue, tt.inputValue, tt.wildcard); got != tt.want {
				t.Errorf("vInList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVInRange(t *testing.T) {
	tests := []struct {
		name       string
		ruleValue  string
		inputValue int64
		want       bool
	}{
		{"Value in Range", "100-500", 250, true},
		{"Value Below Range", "100-500", 50, false},
		{"Value Above Range", "100-500", 600, false},
		{"Value equals Min", "100-500", 100, true},
		{"Value equals Max", "100-500", 500, true},
		{"Invalid Rule Format", "100", 150, false},
		{"Negative Range", "-10-10", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := vInRange(tt.ruleValue, tt.inputValue); got != tt.want {
				t.Errorf("vInRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVInclusionExclusion(t *testing.T) {
	tests := []struct {
		name       string
		ruleValue  string
		inputValue string
		wildcard   string
		want       bool
	}{
		{"Include - Match", "I360;840", "360", "A000", true},
		{"Include - No Match", "I360;840", "123", "A000", false},
		{"Exclude - Match (should be false)", "E360;840", "360", "A000", false},
		{"Exclude - No Match (should be true)", "E360;840", "123", "A000", true},
		{"Wildcard Match", "A000", "any_value", "A000", true},
		{"Single Item Include", "I360", "360", "A000", true},
		{"Single Item Exclude", "E360", "360", "A000", false},
		{"Invalid Indicator", "X360", "360", "A000", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := vInclusionExclusion(tt.ruleValue, tt.inputValue, tt.wildcard); got != tt.want {
				t.Errorf("vInclusionExclusion() = %v, want %v", got, tt.want)
			}
		})
	}
}
