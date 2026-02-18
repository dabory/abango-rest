package abango

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"

	e "github.com/dabory/abango-rest/etc"
)

type RunConf struct {
	RunMode     string
	DevPrefix   string
	ProdPrefix  string
	ConfPostFix string
}

func init() {
}

func RunServicePoint(RestHandler func(ask *AbangoAsk)) {

	var wg sync.WaitGroup
	if XConfig["XDBOn"] == "Yes" { // 필요없어진 건 같으니까 2024년에 지울것.
		MyLinkXDB()
	}
	if XConfig["CrystalDBOn"] == "Yes" { // 필요없어진 건 같으니까 2024년에 지울것.
		MyLinkCrystalDB()
	}
	if XConfig["IsQryFromQDB"] == "Yes" {
		QDBOn = true
	}

	if XConfig["KafkaOn"] == "Yes" {
		// wg.Add(1)
		// go func() {
		// KafkaSvcStandBy(KafkaHandler)
		// 	wg.Done()
		// }()
	}
	if XConfig["gRpcOn"] == "Yes" {
		// wg.Add(1)
		// go func() {
		// 	// GrpcSvcStandBy(GrpcHandler)
		// 	wg.Done()
		// }()
	}
	if XConfig["RestOn"] == "Yes" {
		wg.Add(1)
		go func() {
			RestSvcStandBy(RestHandler)
			wg.Done()
		}()
	}

	wg.Wait()
}

func RunRouterPostNormal(askuri string, body string) (string, string, string) {

	if err := GetXConfig(); err == nil {
		e.InitLog(XConfig["LogFilePath"], XConfig["ShowLogStdout"])
		log.Print("============ RunEndRequest Begins ==============")

		apiMethod := "POST"
		askBytes := []byte(body)
		restUri := XConfig["RestConnect"] + askuri
		if askuri == "upload-file" {
			if retBytes, retstaBytes, err := e.UploadFileResponse(apiMethod, restUri, askBytes); err == nil {
				var out bytes.Buffer
				err := json.Indent(&out, retBytes, "", "  ")
				if err == nil {
					return out.String(), string(retstaBytes), ""
				} else {
					return string(retBytes), string(retstaBytes), ""
				}
			} else {
				return err.Error(), string(retstaBytes), ""
			}
		} else {
			if retBytes, retstaBytes, err := e.GetHttpResponse(apiMethod, restUri, askBytes); err == nil {
				var out bytes.Buffer
				err := json.Indent(&out, retBytes, "", "  ")
				if err == nil {
					return out.String(), string(retstaBytes), ""
				} else {
					return string(retBytes), string(retstaBytes), ""
				}
			} else {
				return err.Error(), string(retstaBytes), ""
			}
		}
	} else {
		return "", "", e.MyErr("XCVZDSFGQWERDZ-Unable to get GetXConfig()", nil, true).Error()
	}
	// return "", "", "asfklsjljfad-Reached to end of RunEndRequest !"
}

func RunEndRequest(params string, body string) string {

	if err := GetXConfig(); err == nil {

		e.InitLog(XConfig["LogFilePath"], XConfig["ShowLogStdout"])
		log.Print("============ RunEndRequest Begins ==============")
		askfile := e.GetAskName()
		arrask := strings.Split(askfile, "@") // login@post 앞의 문자를 askname으로 설정
		askname := arrask[0]

		jsonsend := XConfig["JsonSendDir"] + askname + ".json"

		var err error
		if body, err = e.FileToStr(jsonsend); err != nil {
			return e.MyErr("WERZDSVCZSRE-JsonSendFile Not Found: ", err, true).Error()
		}

		if XConfig["ApiType"] == "Kafka" {
			// return RunRequest(KafkaRequest, &params, &body)
		} else if XConfig["ApiType"] == "gRpc" {
			// return RunRequest(GrpcRequest, &params, &body)
		} else if XConfig["ApiType"] == "Rest" {
			return RunRequest(RestRequest, &params, &body)
		} else {
			return e.MyErr("QREWFGARTEGF-Wrong ApiType in RunEndRequest()", nil, true).Error()
		}
	} else {
		return e.MyErr("XCVZDSFGQWERDZ-Unable to get GetXConfig()", nil, true).Error()
	}
	return "Reached to end of RunEndRequest !"
}

func RunRequest(MsgHandler func(v *AbangoAsk) (string, string, error), params *string, body *string) string {

	var v AbangoAsk
	v.UniqueId = e.RandString(20)
	v.Body = []byte(*body)

	jsonsvrparams := XConfig["JsonServerParamsPath"]
	if file, err := os.Open(jsonsvrparams); err == nil {
		if err = json.NewDecoder(file).Decode(&v.ServerParams); err != nil {
			return e.MyErr("LAAFDFDFERHYWE", err, true).Error()
		}
	} else {
		return e.MyErr("LAAFDFDWDERHYWE-"+jsonsvrparams+" File not found", err, true).Error()
	}

	if *params != "" { //User Params 있을 경우 해당을 가져온다.
		var askparmas []Param
		if err := json.Unmarshal([]byte(*params), &askparmas); err == nil {
			for _, j := range askparmas {
				for _, s := range v.ServerParams {
					if s.Key == j.Key {
						s.Value = j.Value
					} // 여기서 api-method 도 처리됨.
				}
				if j.Key == "ApiType" { // Ask Params 에 ApiType 이 지정되어 있다면
					v.ApiType = j.Value
				}
				if j.Key == "AskName" {
					v.AskName = j.Value
				}
			}
		} else {
			return e.MyErr("WERITOGFSERFDH-AskParams Format mismatched:", nil, true).Error()
		}

	} else {
		askfile := e.GetAskName()
		arrask := strings.Split(askfile, "@") // @앞의 문자를 askname으로 설정
		askname := arrask[0]
		apimethod := ""
		if len(arrask) >= 2 { //만약 argv[1] 이 login@Kafka 형태라면
			apimethod = arrask[1]
		}
		for i := 0; i < len(v.ServerParams); i++ {
			if v.ServerParams[i].Key == "api_method" { //GET, POST
				v.ServerParams[i].Value = apimethod
			}
		}

		v.ApiType = XConfig["ApiType"]
		v.AskName = askname
	}

	if v.ApiType == "" || v.AskName == "" {
		return e.MyErr("QWERDSFAERQRDA-ApiType or AskName was not specified:", nil, true).Error()
	}

	if retstr, retsta, err := MsgHandler(&v); err == nil {

		jsonreceive := XConfig["JsonReceiveDir"] + v.AskName + ".json"
		if XConfig["SaveReceivedJson"] == "Yes" {
			e.StrToFile(jsonreceive, retstr)
		}
		if XConfig["ShowReceivedJson"] == "Yes" {
			fmt.Println("Status: " + retsta + "  ReturnJsonFile: " + jsonreceive)
			fmt.Println(retstr)
		}

		return retstr
	} else {
		return e.MyErr("QWERDSFAERQRDA-MsgHandler", err, true).Error()
	}
}

///쓰이고 있지 않음.
// func GetEnvConf() error { // Kangan only

// 	conf := "conf/"
// 	RunFilename := conf + "run_conf.json"

// 	var run RunConf

// 	if file, err := os.Open(RunFilename); err != nil {
// 		e.MyErr("SDFLJDSAFJA", nil, true)
// 		return err
// 	} else {
// 		decoder := json.NewDecoder(file)
// 		if err = decoder.Decode(&run); err != nil {
// 			e.MyErr("LASJLDFJASFJ", err, true)
// 			return err
// 		}
// 	}

// 	filename := conf + run.RunMode + run.ConfPostFix
// 	if file, err := os.Open(filename); err != nil {
// 		e.MyErr("QERTRRTRRW", err, true)
// 		return err
// 	} else {
// 		decoder := json.NewDecoder(file)
// 		if err = decoder.Decode(&XEnv); err != nil {
// 			e.MyErr("LAAFDFERHY", err, true)
// 			return err
// 		}
// 	}

// 	if XEnv.DbType == "mysql" {
// 		XEnv.DbStr = XEnv.DbUser + ":" + XEnv.DbPassword + "@tcp(" + XEnv.DbHost + ":" + XEnv.DbPort + ")/" + XEnv.DbPrefix + XEnv.DbName + "?charset=utf8"
// 	} else if XEnv.DbType == "mssql" {
// 		// Add on more DbStr of Db types
// 	}

// 	return nil
// }

func MyLinkXDB() { //   항상 연결될 수 있는 MySQL  DB 사전 연결

	dbtype := XConfig["DbType"]
	connstr := XConfig["XDBConnString"] + XConfig["DBOptionString"]
	// connstr := XConfig["DbUser"] + ":" + XConfig["DbPassword"] + "@tcp(" + XConfig["DbHost"] + ":" + XConfig["DbPort"] + ")/" + XConfig["DbName"] + "?charset=utf8"
	var err error
	XDB, err = xorm.NewEngine(dbtype, connstr)

	strArr := strings.Split(connstr, "@tcp")
	if len(strArr) != 2 {
		e.MyErr(strArr[1], err, true)
		return
	}

	XDB.ShowSQL(false)
	XDB.SetMaxOpenConns(100)
	XDB.SetMaxIdleConns(20)
	XDB.SetConnMaxLifetime(60 * time.Second)
	if _, err := XDB.IsTableExist("aaa"); err != nil { //Connect Check
		e.MyErr("ASDFAERAFE-DATABASE DISCONNECTED", err, true)
	} else {
		fmt.Println("XDB CONNECTED :" + strArr[1])
	}

}

func MyLinkCrystalDB() { // Crystal Report Server

	dbtype := XConfig["DbType"]
	connstr := XConfig["CrystalDBConnString"] + XConfig["DBOptionString"]
	// connstr := XConfig["DbUser"] + ":" + XConfig["DbPassword"] + "@tcp(" + XConfig["DbHost"] + ":" + XConfig["DbPort"] + ")/" + XConfig["DbName"] + "?charset=utf8"
	var err error
	CrystalDB, err = xorm.NewEngine(dbtype, connstr)

	strArr := strings.Split(connstr, "@tcp")
	if len(strArr) != 2 {
		e.MyErr(strArr[1], err, true)
		return
	}

	CrystalDB.ShowSQL(false)
	CrystalDB.SetMaxOpenConns(100)
	CrystalDB.SetMaxIdleConns(20)
	CrystalDB.SetConnMaxLifetime(60 * time.Second)

}
