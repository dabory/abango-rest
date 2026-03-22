// Author : Eric Kim
// Build Date : 23 Jul 2018  Last Update 02 Aug 2018
// End-Agent for Passcon Multi OS go binding with Windows, MacOS, iOS, and Android
// All rights are reserved.

// Fundation Functions: 가장 기저가 되는 함수들 (뭔가 추가되면 기저가 아닐 수도)
package fnd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/translate"
	"github.com/dabory/abango-rest"
	e "github.com/dabory/abango-rest/etc"
	"golang.org/x/text/language"
)

var (
	LANG_DIR string
)

var GtransClient *translate.Client
var GtransCtx = context.Background()

func LangStr(y *abango.Controller, index string, engMsg string, args ...interface{}) string {
	langMsg := SelfLangHandler(y, engMsg)
	str := "@" + index + ": " + fmt.Sprintf(langMsg, args...)
	return str
}

// SelfLangHandler: Redis(L1) -> File(L2) -> GoogleTranslate(L3) -> Save(File+Redis)
// func SelfLangHandler(y *abango.Controller, engMsg string, isTheme bool) string { //니중에 Theme별로도 가능함.
func SelfLangHandler(y *abango.Controller, engMsg string) string {

	if y.Gtb.CountryCode == "" {
		y.Gtb.CountryCode = "ko-KR"
	}

	fmt.Println("y.Gtb.CountryCode:", y.Gtb.CountryCode)
	fmt.Println("y.Gtb.CountryCode:", y.Gtb.CountryCode)
	fmt.Println("y.Gtb.CountryCode:", y.Gtb.CountryCode)
	fmt.Println("y.Gtb.CountryCode:", y.Gtb.CountryCode)
	// if countryCode == "en-US" {
	// 	return engMsg // 모든 메지지는 미국식 영어로 되어 있나는 전제
	// }

	countryCode := strings.Replace(y.Gtb.CountryCode, "_", "-", -1)

	md5hash := e.Md5Hashed(engMsg, 32)
	// md5hash := e.Md5Hashed(y.Gtb.HostUser+"::"+engMsg, 32) -- 나준에 theme(HostUser별오도 가능함)
	// 3) redis key: {lang}::{md5}
	redisKey := countryCode + "::" + md5hash

	// 3) Redis 먼저 조회

	fmt.Println("KKK", y.Gtb.CountryCode)
	fmt.Println("KKK", y.Gtb.CountryCode)
	fmt.Println("KKK", y.Gtb.CountryCode)
	fmt.Println("KKK", y.Gtb.CountryCode)

	if v, err := abango.QdbView(redisKey); err == nil {
		return v
	}

	// if abango.QDB != nil {
	// 	if v, err := abango.QDB.Get(context.Background(), redisKey).Result(); err == nil && v != "" {
	// 		return v
	// 	}
	// }

	// 4) 파일 경로: {LANG_DIR}/{lang}/{md5앞2}/{md5}
	prefix2 := md5hash[:2]
	dirPath := filepath.Join(LANG_DIR, countryCode, prefix2)
	filePath := filepath.Join(dirPath, md5hash)
	// fmt.Println("filePath:", filePath)

	// 4) 파일 있으면 읽어서 반환 (+ Redis 저장)
	if b, err := os.ReadFile(filePath); err == nil {
		txt := strings.TrimSpace(string(b))
		if txt != "" {
			if err := abango.QdbUpdate(redisKey, txt); err != nil {
				return e.LogStr("OIUOIGUYE", engMsg+" was NOT Updated")
			}
			return txt
		}
	}

	fmt.Println("여기", y.Gtb.CountryCode)

	target, err := language.Parse(countryCode)
	if err != nil {
		return engMsg
	}

	opts := &translate.Options{Format: translate.Text}
	resp, err := GtransClient.Translate(GtransCtx, []string{engMsg}, target, opts)
	if err != nil || len(resp) == 0 || strings.TrimSpace(resp[0].Text) == "" { // 이거 이헣게 복잡할 필요가 있나 ?
		return e.LogStr("OIUOIJHGUYTG", engMsg+" was NOT Translated")
	}
	translated := strings.TrimSpace(resp[0].Text)

	if err := e.StrToFile(filePath, translated); err != nil {
		e.LogErr("OIUOIJHGUYD", e.FuncNameErr()+": File was NOT saved!", err)
		return engMsg
	}

	if err := abango.QdbUpdate(redisKey, translated); err != nil {
		return e.LogStr("OIUOIGUYE", engMsg+" was NOT Updated")
	}

	fmt.Println("저기", y.Gtb.CountryCode)
	return translated
}
