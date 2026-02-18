// Author : Eric Kim
// Build Date : 23 Jul 2018  Last Update 02 Aug 2018
// End-Agent for Passcon Multi OS go binding with Windows, MacOS, iOS, and Android
// All rights are reserved.

package abg

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

var GtransClient *translate.Client
var GtransCtx = context.Background()

func LangStr(y *abango.Controller, index string, engMsg string, args ...interface{}) string {
	langMsg := SelfLangHandler(y, engMsg)
	str := "@" + index + ": " + fmt.Sprintf(langMsg, args...)
	return str
}

// SelfLangHandler: Redis -> File -> GoogleTranslate -> Save(File+Redis)
func SelfLangHandler(y *abango.Controller, engMsg string) string {
	y.Gtb.LangCode = "ko"
	langCode := y.Gtb.LangCode

	if langCode == "en" {
		return engMsg
	}

	md5hash := e.Md5Hashed(engMsg, 32)

	// 3) redis key: {lang}::{md5}
	redisKey := langCode + "::" + md5hash

	// 3) Redis 먼저 조회
	if abango.QDB != nil {
		if v, err := abango.QDB.Get(context.Background(), redisKey).Result(); err == nil && v != "" {
			return v
		}
	}

	// 4) 파일 경로: {LANG_DIR}/{lang}/{md5앞2}/{md5}
	prefix2 := md5hash[:2]
	dirPath := filepath.Join(LANG_DIR, langCode, prefix2)
	filePath := filepath.Join(dirPath, md5hash)
	// fmt.Println("filePath:", filePath)

	// 4) 파일 있으면 읽어서 반환 (+ Redis 저장)
	if b, err := os.ReadFile(filePath); err == nil {
		txt := strings.TrimSpace(string(b))
		if txt != "" {
			if abango.QDB != nil {
				_ = abango.QDB.Set(context.Background(), redisKey, txt, 0).Err()
			}
			return txt
		}
	}

	// 5) 파일 없으면 Google Translate fallback
	// (최소 비용 구조: 여기서만 외부 호출)

	target, err := language.Parse(langCode)
	if err != nil {
		return engMsg
	}

	opts := &translate.Options{Format: translate.Text}
	resp, err := GtransClient.Translate(GtransCtx, []string{engMsg}, target, opts)
	if err != nil || len(resp) == 0 || strings.TrimSpace(resp[0].Text) == "" { // 이거 이헣게 복잡할 필요가 있나 ?
		return e.LogStr("2398ujwdl", engMsg+" was NOT translated")
	}
	translated := strings.TrimSpace(resp[0].Text)

	if err := e.StrToFile(filePath, translated); err != nil {
		e.LogErr("lskj3sㅂ3s", e.FuncNameErr()+": File was NOT saved!", err)
		return engMsg
	}
	// // 6) 파일 저장 (디렉토리 없으면 생성)
	// _ = os.MkdirAll(dirPath, 0755)
	// _ = os.WriteFile(filePath, []byte(translated), 0644)

	// 7) Redis 저장
	if abango.QDB != nil {
		_ = abango.QDB.Set(GtransCtx, redisKey, translated, 0).Err()
	}

	// 8) 반환
	return translated
}
