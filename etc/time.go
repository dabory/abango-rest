// Author : Eric Kim
// Build Date : 23 Jul 2018  Last Update 02 Aug 2018
// End-Agent for Passcon Multi OS go binding with Windows, MacOS, iOS, and Android
// All rights are reserved.

package etc

import (
	"errors"
	"fmt"
	time "time"
)

func PrevYyyyMm(yyyymm string) (string, error) {
	if t, err := time.Parse("200601", yyyymm); err == nil {
		return t.AddDate(0, -1, 0).Format("200601"), nil
	} else {
		return "", ErrLog(FuncRun("adsa16fasf : "+yyyymm, FuncNameErr()), err)
	}
}

func NextYyyyMm(yyyymm string) (string, error) { // 아직 사용않지만 나중을 위해 추가해둠.
	if t, err := time.Parse("200601", yyyymm); err == nil {
		return t.AddDate(0, 1, 0).Format("200601"), nil
	} else {
		return "", ErrLog(FuncRun("adsa16fas9f : "+yyyymm, FuncNameErr()), err)
	}
}

func getNow() time.Time {
	loc, _ := time.LoadLocation("UTC")
	return time.Now().In(loc)
}

func GetNowUnix(sec ...int) int64 {
	var ret int64
	if sec == nil {
		ret = time.Now().UTC().Unix()
	} else {
		ret = time.Now().Add(time.Duration(sec[0]) * time.Second).UTC().Unix()
	}
	return ret
}

func GetNowDate(i int) string {
	format := "060102"
	if i == 8 {
		format = "20060102"
	} else if i == 6 {
		format = "060102"
	}
	return time.Now().Format(format)
}

func YyyyMmToUnixInterval(yyyymm string) (string, string, error) {

	if len(yyyymm) != 6 {
		return "", "", ErrLog(FuncRun("xcawrq ", FuncNameErr()), errors.New("yyyymmdd format size should be 8"))
	}

	firstDayStr := ""
	lastDayStr := ""
	date, err := time.Parse("200601", yyyymm)
	if err != nil {
		return "", "", err
	}
	firstDay := time.Date(date.Year(), date.Month(), 1, -0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, 0).Add(-time.Second)
	firstDayStr = NumToStr(firstDay.Unix())
	lastDayStr = NumToStr(lastDay.Unix())
	fmt.Println("firstDay", firstDay)
	fmt.Println("lastDay", lastDay)

	return firstDayStr, lastDayStr, nil
}
func YyyyMmDdToUnixInterval(yyyymmdd string) (string, string, error) {

	if len(yyyymmdd) != 8 {
		return "", "", ErrLog(FuncRun("xcawrq ", FuncNameErr()), errors.New("yyyymmdd format size should be 8"))
	}

	firstDayStr := ""
	lastDayStr := ""

	monthStr := yyyymmdd[6:8]
	if monthStr == "00" {
		date, err := time.Parse("200601", yyyymmdd[:6])
		if err != nil {
			return "", "", err
		}
		firstDay := time.Date(date.Year(), date.Month(), 1, -0, 0, 0, 0, time.UTC)
		lastDay := firstDay.AddDate(0, 1, 0).Add(-time.Second)
		firstDayStr = NumToStr(firstDay.Unix())
		lastDayStr = NumToStr(lastDay.Unix())
		// fmt.Println("firstDay", firstDay)
		// fmt.Println("lastDay", lastDay)
	} else {
		date, err := time.Parse("20060102", yyyymmdd)
		if err != nil {
			return "", "", err
		}
		firstDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
		lastDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, time.UTC)
		firstDayStr = NumToStr(firstDay.Unix())
		lastDayStr = NumToStr(lastDay.Unix())
		// fmt.Println("firstDay", firstDay)
		// fmt.Println("lastDay", lastDay)
	}
	return firstDayStr, lastDayStr, nil
}

func YyyyMmToStringInterval(yyyymm string) (string, string, error) {
	if len(yyyymm) != 6 {
		return "", "", LogErr("EWRWQFAEEZ", FuncNameErr(), errors.New("yyyymm should be 6 char but curr one is length "+NumToStr(len(yyyymm))))
	}
	date, err := time.Parse("200601", yyyymm)
	if err != nil {
		return "", "", err
	}
	firstDay := time.Date(date.Year(), date.Month(), 1, -0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, 0).Add(-time.Second)
	return firstDay.Format("20060102"), lastDay.Format("20060102"), nil
}

func UnixToYyyyMm(unix int64) string {
	if !IsUnixTime(unix) {
		LogErr("JHLSAKAS", "", errors.New(NumToStr(unix)+" is Not a Unix Time"))
		return "ERRORS"
	}
	return time.Unix(unix, 0).Format("200601")
}

func IsUnixTime(timestamp int64) bool {
	minUnixTime := int64(0)
	maxUnixTime := time.Now().Unix() + 10*365*24*60*60 // 10 years in seconds
	return timestamp >= minUnixTime && timestamp <= maxUnixTime
}
