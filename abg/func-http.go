package abg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	e "github.com/dabory/abango-rest/etc"
)

func HttpResponseOwnerPost(method string, apiurl string, jsBytes []byte,
	ownerKey string) (retbody []byte, retsta int, reterr error) {
	// 1. 요청 생성
	req, err := http.NewRequest(method, apiurl, bytes.NewBuffer(jsBytes))
	if err != nil {
		return nil, 0, e.LogErr("JYUHIGHA", e.FuncNameErr()+": ", fmt.Errorf("NewRequest error: %w", err))
	}

	// fmt.Println("ownerKey: ", ownerKey)
	// 2. 필수 헤더
	req.Header.Set("Content-Type", "application/json")
	// 3. 사용자 정의 헤더 추가 (예: Authorization, OwnerKey 등)
	if ownerKey != "" {
		req.Header.Set("OwnerKey", ownerKey) // 사용자 정의 헤더
	}

	// 4. 요청 실행
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, 0, e.LogErr("JYUHEGHA", e.FuncNameErr()+": ", fmt.Errorf("HTTP request failed: %w", err))
	}
	defer response.Body.Close()

	// 5. Body 읽기
	retbody, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, 0, e.LogErr("JYUHEGHA", e.FuncNameErr()+": ", fmt.Errorf("HTTP request failed: %w", err))
	}

	return retbody, response.StatusCode, nil
}

func HttpResponseSimplePost(method string, apiurl string, jsBytes []byte) (retbody []byte, retsta int, reterr error) {

	response, err := http.Post(apiurl, "application/json", bytes.NewBuffer(jsBytes))
	if err != nil {
		return nil, 0, errors.New(e.FuncRunErr("65rfg0csdew", "The HTTP request failed with error "+e.FuncNameErr()+err.Error()))
	} else {
		retbody, err = io.ReadAll(response.Body)
		if err != nil {
			return nil, 0, errors.New(e.FuncRunErr("kjda89382", "ReadAll error "+e.FuncNameErr()+err.Error()))
		}
	}
	return retbody, response.StatusCode, nil
}

func HttpResponseSimpleGet(apiurl string) (retbody []byte, reterr error) {

	response, err := http.Get(apiurl)
	if err == nil {
		retbody, err = io.ReadAll(response.Body)
		if err != nil {
			return nil, errors.New(e.FuncRunErr("09t43vda3rf", "ReadAll error "+e.FuncNameErr()+err.Error()))
		}
	} else {
		return nil, errors.New(e.FuncRunErr("567wgq34r", "The HTTP request failed with error "+e.FuncNameErr()+err.Error()))
	}
	return retbody, nil
}
