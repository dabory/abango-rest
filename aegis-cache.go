// Author : Eric Kim
// Build Date : 23 Jul 2018  Last Update 02 Aug 2018
// End-Agent for Passcon Multi OS go binding with Windows, MacOS, iOS, and Android
// All rights are reserved.

package abango

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	e "github.com/dabory/abango-rest/etc"
	grp1 "github.com/dabory/abango-rest/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

const (
	aegisConnErrStatus = "909"
	aegisRpcErrStatus  = "907"
)

var (
	AEGIS_MEMBER_ON  bool
	AEGIS_CACHE_CONN string
)

func AegisEncrypt(r *http.Request, plainText string) (encrypted string, err error) {

	req := AskNameController{
		CurrOtp:  OtpManager.CurrOTP,
		AccessIp: "20.23.324.314",
		Key:      "",
		Function: "encrypt",
	}

	aegBytes, _ := json.Marshal(req)
	retsta, retstr, err := ageisRequest(string(aegBytes), plainText)
	if err != nil {
		return "", e.LogErr("QWVGAEFV", e.FuncNameErr()+": "+retsta, err)
	}

	return retstr, nil // e
}

func AegisDecrypt(r *http.Request, encrypted string) (decrypted string, err error) {

	req := AskNameController{
		CurrOtp:  OtpManager.CurrOTP,
		AccessIp: "20.23.324.314",
		Key:      "",
		Function: "decrypt",
	}

	aegBytes, _ := json.Marshal(req)
	retsta, retstr, err := ageisRequest(string(aegBytes), encrypted)
	if err != nil {
		return "", e.LogErr("QWVGAEFV5", e.FuncNameErr()+": "+retsta, err)
	}

	return retstr, nil
}

func AegisView(r *http.Request, key string) (string, error) {

	req := AskNameController{
		CurrOtp:  OtpManager.CurrOTP,
		AccessIp: "20.23.324.314",
		Key:      key,
		Function: "view",
	}

	aegBytes, _ := json.Marshal(req)
	retsta, retstr, err := ageisRequest(string(aegBytes), "")
	if err != nil {
		return "", e.LogErr("QWVGAE8V", e.FuncNameErr()+": "+retsta, err)
	}

	return retstr, nil
}

func AegisUpdate(r *http.Request, key string, value string) error {

	req := AskNameController{
		CurrOtp:  OtpManager.CurrOTP,
		AccessIp: "20.23.324.314",
		Key:      key,
		Function: "update",
	}

	aegBytes, _ := json.Marshal(req)
	retsta, _, err := ageisRequest(string(aegBytes), value)
	if err != nil {
		return e.LogErr("QWV2GAEFV", e.FuncNameErr()+": "+retsta, err)
	}

	return nil
}

func AegisDelete(key string) error {

	req := AskNameController{
		CurrOtp:  OtpManager.CurrOTP,
		AccessIp: "20.23.324.314",
		Key:      key,
		Function: "delete",
	}

	aegBytes, _ := json.Marshal(req)
	retsta, _, err := ageisRequest(string(aegBytes), "")
	if err != nil {
		return e.LogErr("QWVGAEF7V", e.FuncNameErr()+": "+retsta, err)
	}

	return nil
}

func AegisStatus() (string, error) {

	req := AskNameController{
		CurrOtp:  OtpManager.CurrOTP,
		AccessIp: "20.23.324.314",
		Key:      "",
		Function: "status",
	}

	aegBytes, _ := json.Marshal(req)
	retsta, retstr, err := ageisRequest(string(aegBytes), "")
	if err != nil {
		return "", e.LogErr("QWVGAEF4V", e.FuncNameErr()+": "+retsta, err)
	}

	return retstr, nil
}

func ageisRequest(askname string, askstr string) (string, string, error) {

	// addr := XConfig["AegisCacheConn"]

	conn, err := grpc.NewClient(
		AEGIS_CACHE_CONN,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                20 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return "", "", e.LogErr("QWVGAOYF3V", e.FuncNameErr()+": Aegis Connection Error", err)
	}
	defer conn.Close()

	gClient := grp1.NewGrp1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	r, err := gClient.StdRpc(ctx, &grp1.StdAsk{
		AskName: askname,
		AskStr:  askstr,
	})
	if err != nil {
		return "", "", e.LogErr("QWVGPOYF2", e.FuncNameErr()+": gClient.StdRpc gRpc Error", err)
	}

	return r.RetSta, r.RetStr, nil
}

func InitializeOTP() error {
	sec, err := e.MacSecretGet()
	if err != nil {
		return e.LogErr("LOOUJGYT1", e.FuncNameErr()+": initializeOTP Error", err)
	}
	ServerSecret = sec

	manager := e.NewOTPManager(ServerSecret, e.OtpDigits, e.OtpPeriod)

	ctx, _ := context.WithCancel(context.Background())
	manager.Start(ctx)

	OtpManager = manager

	return nil
}
