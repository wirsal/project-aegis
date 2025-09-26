package codec

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/indece-official/go-ebcdic"
	"github.com/wirsal/project-aegis/pkg/logger"
)

// Convert COMP3 EBCDIC to HEX
// then transalate HEX to ASCII
func Hex2string_comp3(str string) string {
	input := []byte(str)
	return hex.EncodeToString(input)
}

// Convert Alphanumeric EBCDIC  to HEX
// then transalate HEX to ASCII
func Hex2string(str string) string {
	input := []byte(str)

	// Cek apakah semua byte adalah 0x00 (zero value)
	isZero := true
	for _, b := range input {
		if b != 0x00 {
			isZero = false
			break
		}
	}

	if isZero {
		return strings.Repeat(" ", len(input)) // return spasi sebanyak panjang input
	}

	// Lanjutkan konversi normal
	decoded, err := ebcdic.Decode(input, ebcdic.EBCDIC1141)
	if err != nil {
		return ""
	}
	return cleanString(decoded)
}
func cleanString(input string) string {
	replacer := strings.NewReplacer(
		"§", "@",
		"\r\n", " ",
		"\r", " ",
		"\n", " ",
		"|", " ",
		"\x00", " ",
	)
	return replacer.Replace(input)
}

func ParseComp3SignedMode(str, mode string) string {
	val := Hex2string_comp3(str)
	if len(val) == 0 {
		return ""
	}
	if mode == "" {
		return val
	} else {
		sign := val[len(val)-1:]
		number := val[:len(val)-1]

		switch mode {
		case "d", "D":
			if strings.EqualFold(sign, "d") {
				return "-" + number
			}
			return " " + number
		case "c", "C":
			if strings.EqualFold(sign, "c") {
				return "+" + number
			}
			return "-" + number
		default:
			return number
		}
	}
}

func SafeDecode(funcName string, fn func() string) (result string) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(fmt.Sprintf("[%s] decode error", funcName), fmt.Errorf("%v", r))
			result = "" // fallback kalau error
		}
	}()
	return fn()
}
