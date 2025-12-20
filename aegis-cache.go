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

// View -> abango.Init() DB 접속시에만 Request 정보를 받는다.
// Update -> GateTokenGenerate 에서만 Request 정보를 받는다.
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
		return "", e.LogErr("QWVGAVAEFV-AegisView", "ageisRequest failed, status: "+retsta, err)
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
		return e.LogErr("QWVGAVAEFV-MDB.Update Error in Status: "+retsta, "", err)
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
		return e.LogErr("QWVGAVAEFV-MDB.Delete Error in Status: "+retsta, "", err)
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
		return "", e.LogErr("QWVGAVAEFV-MDB.Delete Error in Status: "+retsta, "", err)
	}

	return retstr, nil
}

const (
	aegisConnErrStatus = "909"
	aegisRpcErrStatus  = "907"
)

func ageisRequest(askname string, askstr string) (string, string, error) {

	addr := XConfig["AegisCacheConn"]

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                20 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return "", "", e.LogErr("QWVGAOYFV", "Aegis Connection Error in Status: "+aegisConnErrStatus, err)
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
		return "", "", e.LogErr("QWVGAOYFV", "gClient.StdRpc gRpc Error in Status: "+aegisRpcErrStatus, err)
	}

	return r.RetSta, r.RetStr, nil
}
