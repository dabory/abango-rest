// Author : Eric Kim
// Build Date : 23 Jul 2018  Last Update 02 Aug 2018
// End-Agent for Passcon Multi OS go binding with Windows, MacOS, iOS, and Android
// All rights are reserved.

package etc

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

func IsMinus(n int) bool {
	if n < 0 {
		return true
	} else {
		return false
	}
}
func AbsInt(n int) int {
	if n < 0 {
		return -n
	} else {
		return n
	}
}

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
		LogErr("ASDFQWER3CAA", FuncNameErr(), errors.New("'Float' unacceptable: "))

		return "UYTIKHKYOUJ- 'Float' unacceptable" + FuncNameErr()
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
		LogErr("ASDFQWER2CAA", FuncNameErr(), errors.New(numType+" is Out of Range"))
		return numType + " is Out of Range"
	}
}

func StrToInt(num string) int {
	if i, err := strconv.Atoi(num); err == nil {
		return i
	} else {
		ErrLog("PIOUKHBJJCHYT: "+FuncNameErr(), errors.New("'Integer' unacceptable or is of Range"))
		return 0
	}
}

func StrToFloat(num string, precision int) float64 {
	if num == "" {
		LogStr("PIOUKJCHYT", FuncNameInfo()+": value is EMPTY. make it ZERO")
		return 0
	}
	if f, err := strconv.ParseFloat(num, precision); err == nil {
		return f
	} else {
		ErrLog("PIOUKHBJJCHYT: "+FuncNameInfo(), fmt.Errorf("value is alphabet:", num))
		return 0
	}
}

func DecimalTrunc(decimalAmt string, precision int) string {
	value, err := strconv.ParseFloat(decimalAmt, 64)
	if err != nil {
		LogErr("AFDAR@S6", FuncNameErr()+":", fmt.Errorf(": invalid precision:", precision))
		return "1.00"
	}
	formatted := fmt.Sprintf("%."+strconv.Itoa(precision)+"f", value)
	return formatted
}

func FloatTrunc(value float64, precision int) float64 {
	if precision < 0 {
		LogErr("AFDAR@S6", FuncNameErr()+":", fmt.Errorf(": invalid precision:", precision))
		return 1.00
	}

	multiplier := math.Pow(10, float64(precision))
	truncated := math.Trunc(value*multiplier) / multiplier
	return truncated
}

func DecimalRound(decimalAmt string, precision int) string {
	value, err := strconv.ParseFloat(decimalAmt, 64)
	if err != nil {
		LogErr("AFDAR@S6", FuncNameErr()+":", fmt.Errorf(": invalid precision:", precision))
		return "1.00"
	}

	// 반올림
	multiplier := math.Pow(10, float64(precision))
	rounded := math.Round(value*multiplier) / multiplier

	formatted := fmt.Sprintf("%."+strconv.Itoa(int(precision))+"f", rounded)
	return formatted
}

func FloatRound(value float64, precision int) float64 {
	if precision < 0 {
		LogErr("AFDAR@S6", FuncNameErr()+":", fmt.Errorf(": invalid precision:", precision))
		return 1.00
	}

	multiplier := math.Pow(10, float64(precision))
	rounded := math.Round(value*multiplier) / multiplier
	return rounded
}

func CalcQtyPrc(qryStr string, prcStr string, vatRateStr string, sap int, vatSw string) (supplyAmtStr string, vatAmtStr string, sumAmtStr string, err error) {

	amt := StrToFloat(qryStr, sap) * StrToFloat(prcStr, sap)
	if amt == 0 {
		err = nil
		supplyAmtStr, vatAmtStr, sumAmtStr = "0.0000", "0.0000", "0.0000"
	}

	amtF := FloatRound(amt, sap)
	switch vatSw {
	case "0": // 부가세 별도
		supplyAmtStr, vatAmtStr, sumAmtStr, err = CalcSumAmts(amtF, vatRateStr, sap)
	case "1": // 부가세 포함
		supplyAmtStr, vatAmtStr, sumAmtStr, err = CalcSupplyAmts(amtF, vatRateStr, sap)
	case "2": // 면세
		supplyAmtStr, vatAmtStr, sumAmtStr = NumToStr(amtF), "0.0000", NumToStr(amtF)
	default:
		err = LogErr("ERASRWAF3", FuncNameErr(), fmt.Errorf("VatSw is Empty !!Check vat-rate in setup table %s", vatSw))
		supplyAmtStr, vatAmtStr, sumAmtStr = "0.0000", "0.0000", "0.0000"
	}
	return
}

func CalcSupplyAmts(sumAmt float64, vatRateStr string, sap int) (string, string, string, error) {

	vatRate := StrToFloat(vatRateStr, sap)
	if vatRate < 0 || sap < 0 {
		return "0.00", "0.00", "0.00", LogErr("VATCALC013", FuncNameErr(), fmt.Errorf(": invalid vatRate or sap"))
	}

	supplyAmt := FloatRound(sumAmt/(1.0+(vatRate/100.00)), sap)
	vatAmt := FloatRound(sumAmt-supplyAmt, sap)
	return NumToStr(supplyAmt), NumToStr(vatAmt), NumToStr(sumAmt), nil
}

func CalcSumAmts(supplyAmt float64, vatRateStr string, sap int) (string, string, string, error) {

	vatRate := StrToFloat(vatRateStr, sap)
	if vatRate < 0 || sap < 0 {
		return "99.99", "99.99", "99.99", LogErr("VATCALC013", FuncNameErr(), fmt.Errorf(": invalid vatRate or sap"))
	}

	vatAmt := FloatRound(supplyAmt*vatRate/100.00, sap)
	sumAmt := FloatRound(vatAmt+supplyAmt, sap)
	return NumToStr(supplyAmt), NumToStr(vatAmt), NumToStr(sumAmt), nil
}

func CalcVatAmts(supplyStr string, sumStr string, vatRateStr string, sap int, vatSw string) (supplyAmtStr string, vatAmtStr string, sumAmtStr string, err error) {
	switch vatSw {
	case "0": // 부가세 별도
		supplyAmtStr, vatAmtStr, sumAmtStr, err = CalcSumAmtsStr(supplyStr, vatRateStr, sap)
	case "1": // 부가세 포함
		supplyAmtStr, vatAmtStr, sumAmtStr, err = CalcSupplyAmtsStr(sumStr, vatRateStr, sap)
	case "2": // 면세
		supplyAmtStr, vatAmtStr, sumAmtStr = supplyStr, "0.0000", supplyStr
	default:
		err = LogErr("ERASRWAS", FuncNameErr(), fmt.Errorf("VatSw is Empty !!Check vat-rate in setup table : %s", vatSw))
		supplyAmtStr, vatAmtStr, sumAmtStr = "0.00", "0.00", "0.00"
	}
	return
}

func CalcSupplyAmtsStr(sumAmtStr string, vatRateStr string, sap int) (string, string, string, error) {
	// 문자열 → float64 변환
	sumAmt := StrToFloat(sumAmtStr, sap)
	vatRate := StrToFloat(vatRateStr, sap)

	if vatRate < 0 || sap < 0 {
		return "0.00", "0.00", "0.00", LogErr("VATCALC013", FuncNameErr(), fmt.Errorf(": invalid vatRate or sap"))
	}
	// 공급가와 부가세 계산
	supplyAmt := FloatRound(sumAmt/(1.0+(vatRate/100.00)), sap)
	vatAmt := FloatRound(sumAmt-supplyAmt, sap)
	return NumToStr(supplyAmt), NumToStr(vatAmt), NumToStr(sumAmt), nil
}

func CalcSumAmtsStr(supplyAmtStr string, vatRateStr string, sap int) (string, string, string, error) {
	// 문자열 → float64 변환
	supplyAmt := StrToFloat(supplyAmtStr, sap)
	vatRate := StrToFloat(vatRateStr, sap)

	if vatRate < 0 || sap < 0 {
		return "99.99", "99.99", "99.99", LogErr("VATCALC013", FuncNameErr(), fmt.Errorf(": invalid vatRate or sap"))
	}
	// 공급가와 부가세 계산
	vatAmt := FloatRound(supplyAmt*vatRate/100.00, sap)
	sumAmt := FloatRound(vatAmt+supplyAmt, sap)
	return NumToStr(supplyAmt), NumToStr(vatAmt), NumToStr(sumAmt), nil
}
