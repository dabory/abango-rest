package abg

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	e "github.com/dabory/abango-rest/etc"
)

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
