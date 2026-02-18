// Author : Eric Kim
// Build Date : 23 Jul 2018  Last Update 02 Aug 2018
// End-Agent for Passcon Multi OS go binding with Windows, MacOS, iOS, and Android
// All rights are reserved.

package etc

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"html"
	"os"
	"reflect"
	"regexp"
	"strings"
	time "time"

	"github.com/microcosm-cc/bluemonday"
	// "xorm.io/xorm"
)

var EmailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func TwoArrayFromStrings(tableNames, dateFields string) ([]string, []string, error) {
	tableList := strings.Split(tableNames, ",")
	dateList := strings.Split(dateFields, ",")

	if len(tableList) != len(dateList) {
		// 영문 에러 메시지로 변경
		return nil, nil, fmt.Errorf(
			"configuration mismatch: count of tables (%d) does not match count of date fields (%d)",
			len(tableList), len(dateList),
		)
	}

	return tableList, dateList, nil
}

// IsValidEmail은 문자열이 유효한 이메일 형식인지 확인합니다.
func IsValidEmail(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	return EmailRegex.MatchString(email)
}

func ExtractWithoutHttp(raw string) string {
	raw = strings.TrimPrefix(raw, "http://")
	raw = strings.TrimPrefix(raw, "https://")
	return raw
}

func ExtractDomainSimple(uri string) string {
	// Remove scheme if present
	uri = strings.TrimPrefix(uri, "http://")
	uri = strings.TrimPrefix(uri, "https://")
	uri = strings.TrimPrefix(uri, "//")

	// Remove path and query parameters
	if idx := strings.Index(uri, "/"); idx != -1 {
		uri = uri[:idx]
	}
	if idx := strings.Index(uri, "?"); idx != -1 {
		uri = uri[:idx]
	}

	return uri
}

func SafeTrimUTF8(s string, start, size int) string {
	if start < 0 || size < 0 {
		return ""
	}

	runes := []rune(s)
	length := len(runes)

	if start >= length {
		return ""
	}

	end := start + size
	if end > length {
		end = length
	}

	return string(runes[start:end])
}

func SafeTrim(s string, start, size int) string {
	if start < 0 || size < 0 {
		return ""
	}
	length := len(s)
	if start >= length {
		return ""
	}
	end := start + size
	if end > length {
		end = length
	}
	return s[start:end]
}

func StripHtml(cont string, max int) string {
	// 1. HTML 엔티티 디코딩
	unescaped := html.UnescapeString(cont)

	// 2. HTML 태그 제거
	p := bluemonday.StripTagsPolicy()
	sanitized := p.Sanitize(unescaped)

	// 3. 최대 길이 제한 (유니코드 고려)
	runes := []rune(sanitized)
	if len(runes) > max {
		return string(runes[:max])
	}
	return sanitized
}

func Sanitize(cont string) string {
	p := bluemonday.UGCPolicy()
	return p.Sanitize(cont)
}

func AddStrIfNotExist(s *string, target string) {
	if !strings.Contains(*s, target) {
		*s += target
	}
}

// 24.8.22 이렇게 수정하고 점진적으로 파라를 수정할 것. 파라수정이 다되면 이 평션은 제거랗 것.
func QueryNameToTableAndPrefix(tgtTblName string, tgtTblFullName string) (string, string) {
	if tgtTblFullName == "" {
		prefix := strings.Replace(tgtTblName, "-", "_", -1)
		tbl := ""
		if prefix == "widget_taxo" {
			tbl = "pro_"
		} else if prefix == "sso_app" || prefix == "api23_app" {
			tbl = "main_"
		} else {
			tbl = "dbr_"
		}
		// fmt.Println("qrytableName:", tbl+prefix)
		// fmt.Println("tgtPrefix:", prefix)
		return tbl + prefix, prefix
	} else {
		return tgtTblFullName, tgtTblName
	}

}

func HasPickActPage(uri string, table string) bool {
	if table == "member" { // member-act로 호출하면 MemberSecuredAct로 처리해줌.
		if uri == "/"+table+"-pick" || uri == "/"+table+"-act" || uri == "/"+table+"-page" {
			// if uri == "/"+table+"-pick" || uri == "/"+table+"-act" || uri == "/"+table+"-page" || uri == "/"+table+"-secured-pick" || uri == "/"+table+"-secured-page" || uri == "/"+table+"-secured-act" {
			return true
		} else {
			return false
		}
	} else {
		if uri == "/"+table+"-pick" || uri == "/"+table+"-act" || uri == "/"+table+"-page" {
			return true
		} else {
			return false
		}
	}
}

func HasPickActPageDelPage(uri string, table string) bool {
	if uri == "/"+table+"-pick" || uri == "/"+table+"-act" || uri == "/"+table+"-page" || uri == "/"+table+"-del-page" {
		return true
	} else {
		return false
	}
}

func YesToTrue(yes string, swName string) bool {
	if yes == "Yes" {
		fmt.Println("== config Key: " + swName + " is ON ==")
		return true
	} else {
		return false
	}
}

func TimeFormatGet(format string) string {
	rtn := ""
	if format == "" {
		rtn = "060102"
	} else if format == "YYMMDD" {
		rtn = "060102"
	} else if format == "HHMMSS" {
		rtn = "15:04:05"
	} else if format == "YYYYMMDD" {
		rtn = "20060102"
	} else if format == "YY-MM-DD" {
		rtn = "06-01-02"
	} else if format == "YY.MM.DD" {
		rtn = "06.01.02"
	} else if format == "YYMM" {
		rtn = "0601"
	} else if format == "YY" {
		rtn = "06"
	}
	return rtn
}

func SlipNoPlusRandom() string {
	t := time.Now()
	date := t.Format(TimeFormatGet("YYMMDD"))
	second := strings.Replace(t.Format(TimeFormatGet("HHMMSS")), ":", "", -1)
	return date + "-" + second + "-" + RandString(6)
}

func ToPdpSerial(serial int) string {
	padSize := "8"
	if serial > 100000000 {
		padSize = "9"
	}
	return fmt.Sprintf("%0"+padSize+"d", serial)
}

func RandBytes(i int) []byte {

	return []byte(RandString(i))
}

func RandString(i int) string {
	b := make([]byte, i)
	rand.Read(b)
	return (base64.URLEncoding.EncodeToString(b))[0:i]
}

func GetAskName() string {

	i := len(os.Args)
	if i < 2 {
		MyErr("ZMXCDKDALKSJD", errors.New("command arguments are less then 2"), true)
	} else {

		return os.Args[i-1] // 맨뒤의 Arg 를 json 화일명으로 간주
	}
	return ""
}

// func (t *EnvConf) StrToStruct(str string) {
// 	if err := json.Unmarshal([]byte(str), &t); err != nil {
// 		MyErr("sjdfljsf", err, true)
// 	}
// 	return
// }

// func structToMap(in interface{}, tag string) (map[string]interface{}, string) {
// 	out := make(map[string]interface{})

// 	v := reflect.ValueOf(in)
// 	if v.Kind() == reflect.Ptr {
// 		v = v.Elem()
// 	}

// 	// we only accept structs
// 	if v.Kind() != reflect.Struct {
// 		fmt.Errorf("ToMap only accepts structs; got %T", v)
// 		return nil, MyErr("only accepts structs", nil, false)
// 	}

// 	typ := v.Type()
// 	for i := 0; i < v.NumField(); i++ {
// 		// gets us a StructField
// 		fi := typ.Field(i)
// 		if tagv := fi.Tag.Get(tag); tagv != "" {
// 			out[tagv] = v.Field(i).Interface()
// 		}
// 	}
// 	return out, ""
// }

func ParentDir(params ...string) string {

	var workdir string
	if len(params) == 0 {
		workdir, _ = os.Getwd()
	} else {
		workdir = params[0]
	}

	sp := strings.Split(workdir, "/")
	parentdir := ""
	for i := 1; i < len(sp)-1; i++ {
		parentdir += "/" + sp[i]
	}
	return parentdir
}

func TableName(m interface{}) string {
	s := reflect.TypeOf(m).Name()
	return SnakeString(s)
}

func PageName(m interface{}) string {
	s := reflect.TypeOf(m).Elem().Name()
	return SnakeString(s)
}

// snake string, XxYy to xx_yy , XxYY to xx_yy
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// camel string, xx_yy to XxYy
func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	flag, num := true, len(s)-1
	for i := 0; i <= num; i++ {
		d := s[i]
		if d == '_' {
			flag = true
			continue
		} else if flag {
			if d >= 'a' && d <= 'z' {
				d = d - 32
			}
			flag = false
		}
		data = append(data, d)
	}
	return string(data[:])
}

func StructToMap(in interface{}, tag string) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		fmt.Errorf("ToMap only accepts structs; got %T", v)
		return nil, MyErr("HJKMNHGYUH-only accepts structs:", nil, false)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if tagv := fi.Tag.Get(tag); tagv != "" {
			out[tagv] = v.Field(i).Interface()
		}
	}
	return out, nil
}
func FileNameFromPath(input string) string {
	result := strings.ReplaceAll(input, "\\", "/")
	parts := strings.Split(result, "/")
	return parts[len(parts)-1]
}

// return the source filename after the last slash
func ChopPath(original string) string {
	i := strings.LastIndex(original, "/")
	if i == -1 {
		return original
	} else {
		return original[i+1:]
	}
}

func TransTblName(word string) string {
	switch strings.ToLower(word) {
	case "dbr_sales":
		return "매출등록"
	case "dbr_purch":
		return "매입등록"
	}
	return "NoName"
}
