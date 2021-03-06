package etc

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	time "time"
)

// func GetHttpResponseSimplePost(method string, apiurl string, jsBytes []byte) (retbody []byte, retsta []byte, reterr error) {

// 	response, err := http.Post(apiurl, "application/json", bytes.NewBuffer(jsBytes))
// 	if err != nil {
// 		fmt.Printf("The HTTP request failed with error %s\n", err)
// 	} else {
// 		retbody, _ = ioutil.ReadAll(response.Body)
// 	}
// 	return retbody, []byte(strconv.Itoa(response.StatusCode)), nil

// }

func GetHttpResponse(method string, apiurl string, jsBytes []byte) ([]byte, []byte, error) {

	reader := bytes.NewBuffer(jsBytes)
	req, err := http.NewRequest(method, apiurl, reader)
	if err != nil {
		return nil, []byte("909"), MyErr("WERZDSVADFZ-http.NewRequest", err, false)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Endpoint-Agent", "abango-rest-api-v1.0")
	req.Header.Add("Accept-Language", "en-US")
	req.Header.Add("User-Agent", runtime.GOOS+"-"+runtime.Version()) // for checking OS Type in Server

	i := len(os.Args)
	if i != 1 { // 1일 경우는 go function call 의 경우 이므로  memory fault 가 난다.
		gateToken := os.Args[i-2]
		if len(gateToken) == 20 { // Argument 뒤에서 2번째 Arg가 20자리이면 GateToken 이라고 간주
			req.Header.Add("GateToken", gateToken)
		}
	}

	req.Body = ioutil.NopCloser(bytes.NewReader(jsBytes))

	// Client객체에서 Request 실행
	client := &http.Client{
		Timeout: time.Second * 20, //Otherwirse, it can cause crash without this line. Must Must.
	} // Normal is 10 but extend 20 on 1 Dec 2018

	// fmt.Println(reflect.TypeOf(respo))
	resp, err := client.Do(req)
	if err != nil {
		return nil, []byte("909"), MyErr("WERZDSVXBDCZSRE-client.Do "+apiurl, err, false)
	}
	defer resp.Body.Close()

	byteRtn, _ := ioutil.ReadAll(resp.Body)
	return byteRtn, []byte(strconv.Itoa(resp.StatusCode)), nil

}

func GetHttpResponseOLd(method string, apiurl string, jsBytes []byte) ([]byte, []byte, error) {
	// func GetHttpResponseOld(method string, apiurl string, jsBytes []byte) ([]byte, []byte, error) {
	form := url.Values{}
	// form.Add("postvalues", string(kkk))
	// Values.Encode() encodes the values into "URL encoded" form sorted by key.
	// eForm := v.Encode()
	// fmt.Printf("v.Encode(): %v\n", s)
	reader := strings.NewReader(form.Encode()) // This causes 411 error

	// req, err := http.NewRequest("GET", apiurl, reader)
	req, err := http.NewRequest(method, apiurl, reader)
	if err != nil {
		return nil, []byte("909"), MyErr("WERZDSVADFZ-http.NewRequest", err, true)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Endpoint-Agent", "abango-rest-api-v1.0")
	req.Header.Add("Accept-Language", "en-US")
	req.Header.Add("User-Agent", runtime.GOOS+"-"+runtime.Version()) // for checking OS Type in Server

	req.Body = ioutil.NopCloser(bytes.NewReader(jsBytes))

	// Client객체에서 Request 실행
	client := &http.Client{
		Timeout: time.Second * 20, //Otherwirse, it can cause crash without this line. Must Must.
	} // Normal is 10 but extend 20 on 1 Dec 2018

	// fmt.Println(reflect.TypeOf(respo))
	resp, err := client.Do(req)
	if err != nil {
		return nil, []byte("909"), MyErr("REWTAVDEQWFAF-client.Do "+apiurl, err, true)
	}
	defer resp.Body.Close()

	byteRtn, _ := ioutil.ReadAll(resp.Body)
	return byteRtn, []byte(strconv.Itoa(resp.StatusCode)), nil

}
