package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/wallet/base"
)

type AvailablePortResponse struct {
	Code Code   `json:"code"`
	Msg  string `json:"msg"`
	Data uint16 `json:"data"`
}

type IsPortListeningResponse struct {
	Code Code   `json:"code"`
	Msg  string `json:"msg"`
	Data bool   `json:"data"`
}

func RequestServerGetPortAvailable(host string) int {
	targetUrl := fmt.Sprintf("http://%s/api/v1/lnurl/available_port", host)

	request, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		fmt.Printf("%s http.NewRequest err :%v\n", GetTimeNow(), err)
		return 0
	}

	// 设置基本认证
	username := base.QueryConfigByKey("BasicAuthUser")
	password := base.QueryConfigByKey("BasicAuthPass")
	request.SetBasicAuth(username, password)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("%s http Get err :%v\n", GetTimeNow(), err)
		return 0
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("%s Body.Close err :%v\n", GetTimeNow(), err)
		}
	}(response.Body)

	bodyBytes, _ := io.ReadAll(response.Body)
	var resp AvailablePortResponse
	if err := json.Unmarshal(bodyBytes, &resp); err != nil {
		fmt.Printf("%s RSGPA json.Unmarshal :%v\n", GetTimeNow(), err)
		return 0
	}
	return int(resp.Data)
}

func RequestPostServerIsPortListening(remotePort string) bool {
	host := base.QueryConfigByKey("LnurlServerHost")
	targetUrl := fmt.Sprintf("http://%s/api/v1/lnurl/is_port_listening?remote_port=%s", host, remotePort)

	request, err := http.NewRequest("POST", targetUrl, nil)
	if err != nil {
		fmt.Printf("%s http.NewRequest :%v\n", GetTimeNow(), err)
		return true
	}

	// 设置基本认证
	username := base.QueryConfigByKey("BasicAuthUser")
	password := base.QueryConfigByKey("BasicAuthPass")
	request.SetBasicAuth(username, password)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("%s http.PostForm :%v\n", GetTimeNow(), err)
		return true
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("%s Body.Close err :%v\n", GetTimeNow(), err)
		}
	}(response.Body)

	bodyBytes, _ := io.ReadAll(response.Body)
	var resp IsPortListeningResponse
	if err := json.Unmarshal(bodyBytes, &resp); err != nil {
		fmt.Printf("%s RPSIPL json.Unmarshal :%v\n", GetTimeNow(), err)
		return true
	}
	return resp.Data
}

func GetServerRequestAvailablePort(host string) (int, error) {
	targetUrl := fmt.Sprintf("http://%s/api/v1/lnurl/available_port", host)

	request, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		return 0, errors.Wrap(err, "http.NewRequest")
	}

	// 设置基本认证
	username := base.QueryConfigByKey("BasicAuthUser")
	password := base.QueryConfigByKey("BasicAuthPass")
	request.SetBasicAuth(username, password)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return 0, errors.Wrap(err, "client.Do")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("%s Body.Close err :%v\n", GetTimeNow(), err)
		}
	}(response.Body)

	bodyBytes, _ := io.ReadAll(response.Body)
	var resp AvailablePortResponse
	if err = json.Unmarshal(bodyBytes, &resp); err != nil {
		return 0, errors.Wrap(err, "json.Unmarshal")
	}
	return int(resp.Data), nil
}

func PostServerRequestIsPortListening(remotePort string) (bool, error) {
	host := base.QueryConfigByKey("LnurlServerHost")
	targetUrl := fmt.Sprintf("http://%s/api/v1/lnurl/is_port_listening?remote_port=%s", host, remotePort)

	request, err := http.NewRequest("POST", targetUrl, nil)
	if err != nil {
		return false, errors.Wrap(err, "http.NewRequest")
	}

	// 设置基本认证
	username := base.QueryConfigByKey("BasicAuthUser")
	password := base.QueryConfigByKey("BasicAuthPass")
	request.SetBasicAuth(username, password)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return false, errors.Wrap(err, "client.Do")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("%s Body.Close err :%v\n", GetTimeNow(), err)
		}
	}(response.Body)

	bodyBytes, _ := io.ReadAll(response.Body)
	var resp IsPortListeningResponse
	if err = json.Unmarshal(bodyBytes, &resp); err != nil {
		return false, errors.Wrap(err, "json.Unmarshal")
	}
	return resp.Data, nil
}
