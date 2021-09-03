// Author : Eric Kim
// Build Date : 23 Jul 2018  Last Update 02 Aug 2018
// End-Agent for Passcon Multi OS go binding with Windows, MacOS, iOS, and Android
// All rights are reserved.

package etc

import (
	"fmt"
	"reflect"
	"strconv"
)

func DecimalToInt(decimal string) int {
	if tmp, err := strconv.ParseFloat(decimal, 4); err == nil {
		return int(tmp)
	} else {
		ErrLog("LJLJLJLJIUU-decimal unacceptable", err)
		return 999999999
	}
}

func DecimalToFloat64(decimal string) float64 {
	if num, err := strconv.ParseFloat(decimal, 4); err == nil {
		return num
	} else {
		ErrLog("UYTIUYKHKIU-decimal unacceptable", err)
		return 999999999999.9999
	}
}

func NumToDecimal(num interface{}, precision string) string {

	switch reflect.TypeOf(num).String() {
	case "int":
		tmp := float64(num.(int))
		return fmt.Sprintf("%"+precision+"f", tmp)
	case "int32":
		tmp := float64(num.(int))
		return fmt.Sprintf("%"+precision+"f", tmp)
	case "int64":
		tmp := float64(num.(int))
		return fmt.Sprintf("%"+precision+"f", tmp)
	case "float32":
		return fmt.Sprintf("%"+precision+"f", num)
	case "float64":
		return fmt.Sprintf("%"+precision+"f", num)
	default:
		return "Out of Range"
	}
}

func NumToStr(num interface{}) string {
	numType := reflect.TypeOf(num).String()
	switch reflect.TypeOf(num).String() {
	case "int":
		return fmt.Sprintf("%d", num)
	case "int32":
		return fmt.Sprintf("%d", num)
	case "int64":
		return fmt.Sprintf("%d", num)
	case "float32":
		return fmt.Sprintf("%f", num)
	case "float64":
		return fmt.Sprintf("%f", num)
	default:
		return numType + " is Out of Range"
	}
}

func StrToInt(num string) int {
	if i, err := strconv.Atoi(num); err == nil {
		return i
	} else {
		ErrLog("PIOUKHBJJCHYT-'Integer' unacceptable", err)
		return 0
	}
}

func StrToFloat(num string, decimalP int) float64 {
	if f, err := strconv.ParseFloat(num, decimalP); err == nil {
		return f
	} else {
		ErrLog("PIOUKHBJJCHYT- 'Float' unacceptable", err)
		return 0
	}
}
