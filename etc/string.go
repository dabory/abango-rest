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
	"os"
	"reflect"
	"strings"
	time "time"

	"github.com/microcosm-cc/bluemonday"
	// "xorm.io/xorm"
)

// func LastQry(qry xorm.Session) string {
// 	ret, _ := qry.LastSQL()
// 	fmt.Println("\n" + ret + "\n")
// 	return ret
// }

func StripHtml(cont string, max int) string {
	p := bluemonday.StripTagsPolicy()
	s := p.Sanitize(cont)
	if len([]rune(s)) > max { // len(s) 로 크기를 한글의 경우 확인하면 안된다.
		return string([]rune(s)[:max])
	} else {
		return s
	}
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
		LogNil("== config Key: " + swName + " is ON ==")
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
