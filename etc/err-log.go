// Author : Eric Kim
// Build Date : 23 Jul 2018  Last Update 02 Aug 2018
// End-Agent for Passcon Multi OS go binding with Windows, MacOS, iOS, and Android
// All rights are reserved.

package etc

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
)

// var InfoLog *log.Logger

func InitLog(path string, showstdout string) error {

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return MyErr("CXZFDREADSF-MyLog file could not be opened: ", err, true)
	}
	// defer file.Close()
	// 파일과 화면에 같이 출력하기 위해 MultiWriter 생성

	if showstdout == "Yes" {
		multiWriter := io.MultiWriter(file, os.Stdout)
		log.SetOutput(multiWriter)
	} else {
		log.SetOutput(file)
	}
	// InfoLog = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	return nil
}

func PageCntErr(index string, tablename string) string {
	return ErrStr(index, " Count Query Error "+tablename+" ")
}

func PageRead(index string, tablename string) string {
	return ErrStr(index, " Page Read from "+tablename+" ")
}

func PageQryErr(index string, tablename string) string {
	return ErrStr(index, " Page Query Error "+tablename+" ")
}

func RecRead(index string, msgString string) string {
	return LogStr(index, " Read "+msgString+" ")
}

func RecNotFound(index string, msgString string) string {
	return LogStr(index, " Not Found "+msgString+" ")
}

func RecReadErr(index string, msgString string) string {
	return ErrStr(index, " B error in Reading "+msgString+" ")
}

func RecAdded(index string, msgString string) string {
	return LogStr(index, " Add "+msgString+" ")
}

func RecNotAdded(index string, msgString string) string {
	return LogStr(index, " Not Added "+msgString+" ")
}

func RecAddErr(index string, msgString string) string {
	return ErrStr(index, " DB error in Adding "+msgString+" ")
}

func RecEdited(index string, msgString string) string {
	return LogStr(index, " Edited "+msgString+" ")
}

func RecNotEdited(index string, msgString string) string {
	return LogStr(index, " Not Edited-'Same Contents Update' is NOT necessary with "+msgString+" ")
}

func RecEditErr(index string, msgString string) string {
	return ErrStr(index, " DB error in Editing "+msgString+" ")
}

func RecDeleted(index string, msgString string) string {
	return LogStr(index, " Delete "+msgString+" ")
}

func RecNotDeleted(index string, msgString string) string {
	return LogStr(index, " Not Deleted "+msgString+" ")
}

func RecDelErr(index string, msgString string) string {
	return LogStr(index, " DB error in Deleting "+msgString+" ")
}

func FuncRun(index string, funcname string) string {
	return LogStr(index, " Func : "+funcname+" ")
}

func FuncRunErr(index string, funcname string) string {
	return ErrStr(index, "-> "+funcname+" ")
}

func JsonFormatErr(index string, structname string) string {
	return ErrStr(index, " Json Format mismatches in "+structname+" ")
}

func ErrStr(index string, s string) string { // nㅣl 아님 경우만 처리(!!중요)
	msg := s
	str := index + " @ " + msg
	log.Println("[Err]: " + str)

	return msg
}

func LogStr(index string, s string) string { // nㅣl 아님 경우만 처리(!!중요)
	msg := s
	str := index + " @ " + msg
	log.Println("[Log]: " + str)

	return msg
}

func LogErr(index string, s string, err error) error { // nㅣl 아님 경우만 처리(!!중요)
	var errStr string
	if err != nil {
		errStr = err.Error()
	} else {
		log.Println("========= Fatal: error is nil LogErr ==========")
	}
	msg := s + " * " + errStr
	str := index + " @ " + msg
	log.Println("[Err]: " + str)
	return errors.New(msg)
}

func LogWarning(index string, s string) { // nㅣl 아님 경우만 처리(!!중요)
	str := index + " @ " + s
	log.Println("[Waring]: " + str)
	return
}

func FuncNameErr() string { // nil 아닌 경우만 처리(!!중요)
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	// fmt.Printf("FuncPath: %s:%d %s\n", frame.File, frame.Line, frame.Function)
	// fmt.Println("Lastindex:", s[strings.LastIndex(s, "/")+1:])
	WhereAmI(GetCurrentDepth())
	s := frame.Function
	return s[strings.LastIndex(s, "/")+1:] + "\n"
}

func FuncNameInfo() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	s := frame.Function
	return s[strings.LastIndex(s, "/")+1:] + " : "
}

func CallerFuncName() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	s := frame.Function
	return s[strings.LastIndex(s, "/")+1:]
}

func CcallerFuncName() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(4, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	s := frame.Function
	return s[strings.LastIndex(s, "/")+1:]
}

func LogCritical(index string, s string, err error) { //에러 ㄱ계를 추적
	var errStr string
	if err != nil {
		errStr = err.Error()
	} else {
		log.Println("========= Fatal: error is nil LogCritical==========")
	}
	str := index + " @ " + s + " * " + errStr
	log.Println("[Fatal]: " + str)

	whereami(2)
	whereami(3)
	whereami(4)
	fmt.Println(strings.Repeat("=", 80))

}

func LogFatal(index string, s string, err error) { //Critical 동일하지만 프로세스 중단
	var errStr string
	if err != nil {
		errStr = err.Error()
	} else {
		log.Println("========= Fatal: error is nil LogFatal ==========")
	}
	str := index + " @ " + s + " * " + errStr
	log.Println("[Fatal]: " + str)

	whereami(2)
	whereami(3)
	whereami(4)
	fmt.Println(strings.Repeat("=", 80))

	os.Exit(100)
}

// ==== 아래건은 모두 옛날 것이라 차츰 Deprecate 할 것 =====
func LogNil(s string) error {
	// log.Logger
	log.Println("[OK]: " + s)
	return nil
}

func AokLog(s string) {
	log.Println("[Abango-OK]: " + s)
}

func ErrLog(s string, err error) error { // // nㅣl처리 아주 중요함 ( 이건 이제 더 사용하지 말것)
	var errStr string
	if err != nil {
		errStr = err.Error()
	} else {
		log.Println("========= Fatal: error is nil ErrLog==========")
	}

	str := "[Err]: " + s + " (Err): " + errStr
	log.Println(str)
	return errors.New(str)
}

func ChkLog(point string, x ...interface{}) {
	log.Println("[CHECK:" + point + "] " + fmt.Sprintf("%v", x))
}

func MyErr(s string, e error, eout bool) error {
	fmt.Println("[MyErr] Position -> ", s, strings.Repeat("=", 40))

	emsg := ""
	if e != nil {
		emsg = "Error: " + e.Error()
	} else {
		emsg = "ERROR is Nil: Wrong Error Check: Check err != OR err == is correct !  "
	}
	fmt.Println(emsg, "\n")
	whereami(2)
	whereami(3)
	whereami(4)
	fmt.Println(strings.Repeat("=", 80))

	if e != nil && eout == true { // quit running if it is FATAL ERROR
		log.Println("[FATAL-ERROR] : EXIT 100")
		os.Exit(100)
	}
	return errors.New(emsg)
}

func agErr(s string, e error, amsg *string) string {
	fmt.Println("== agErr ", strings.Repeat("=", 90))
	// fpcs := make([]uintptr, 1)
	// n := runtime.Callers(2, fpcs)
	// if n == 0 {
	// 	fmt.Println("MSG: NO CALLER")
	// }
	// // caller := runtime.FuncForPC(fpcs[0] - 1)
	// caller := runtime.FuncForPC(fpcs[0])
	// // fmt.Println(caller.FileLine(fpcs[0] - 1))
	// fmt.Println(caller.FileLine(fpcs[0]))
	// fmt.Println(caller.Name())
	emsg := ""
	if e != nil {
		emsg = "Error: " + e.Error() + " in " + s
	} else {
		emsg = "Error: error is nil" + " in " + s // e 가 nil 인 상태에서 Error() 인용시 runtime error
	}
	fmt.Println(emsg, "\n")
	whereami(2)
	whereami(3)
	fmt.Println(strings.Repeat("=", 100))
	return emsg
}

func whereami(i int) {
	function, file, line, _ := runtime.Caller(i)
	fmt.Printf("  %d.File: %s - %d  %s\n   func: %s \n", i, ChopPath(file), line, file, runtime.FuncForPC(function).Name())
}

func WhereAmI(depthList ...int) {
	var depth int
	if depthList == nil {
		depth = 1
	} else {
		depth = depthList[0]
	}

	fmt.Printf("\n==Func Location Start==\n")
	for i := 0; i < depth+1; i++ {
		if i != 0 && i != 1 && i < 7 {
			function, file, line, _ := runtime.Caller(i)
			fmt.Printf("==Level %d :: ", i)
			fmt.Printf("File: %s, Trace at %d Line\nFunction: %s \n", file, line, runtime.FuncForPC(function).Name())
		}
	}
	fmt.Printf("==End==\n\n")
	return
}

func QryPathSql(path string, sql string) string { // // nㅣl처리 아주 중요함 ( 이건 이제 더 사용하지 말것)
	return "====" + path + "====\n" + sql + "\n" + "===================\n"
}

func GetCurrentDepth() int {
	// Use runtime.Callers to measure current call stack depth
	pc := make([]uintptr, 10) // Adjust size as needed
	depth := runtime.Callers(1, pc)
	return depth - 2 // Exclude the current and runtime.Callers frames
}
