package api

import (
	"fmt"

	"github.com/wallet/base"
	"gopkg.in/resty.v1"
)

var (
	_token string
)

func SetToken(token string) {
	localhost := "127.0.0.1"
	targetUrl := fmt.Sprintf("http://%s:%s/lnurl/set_token", localhost, LnurlRouterPort)
	client := resty.New()

	var r RespStr

	resp, err := client.R().
		SetBasicAuth(base.QueryConfigByKey("BasicAuthUser"), base.QueryConfigByKey("BasicAuthPass")).
		SetFormData(map[string]string{
			"token": token,
		}).
		SetResult(&r).
		SetError(&r).
		Post(targetUrl)

	_ = resp

	if err != nil {
		fmt.Println("client.R() err:", err)
	}
	if r.Msg != "" {
		fmt.Println("err:", r.Msg)
	}
	if r.Data != "" {
		fmt.Println("SetToken resp:", r.Data)
	}

	return
}

func GetToken() string {
	localhost := "127.0.0.1"
	targetUrl := fmt.Sprintf("http://%s:%s/lnurl/get_token", localhost, LnurlRouterPort)
	client := resty.New()

	var r RespStr

	resp, err := client.R().
		SetBasicAuth(base.QueryConfigByKey("BasicAuthUser"), base.QueryConfigByKey("BasicAuthPass")).
		SetResult(&r).
		SetError(&r).
		Post(targetUrl)

	_ = resp
	if err != nil {
		fmt.Println("client.R() err:", err)
	}
	if r.Msg != "" {
		fmt.Println("err:", r.Msg)
	}

	return r.Data
}

func setToken(token string) {
	_token = token
}

func getToken() string {
	return _token
}
