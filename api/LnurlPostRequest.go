package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/wallet/base"
)

type InvoiceResponse struct {
	Code Code   `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

type UploadUserResp struct {
	Code Code   `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

type UploadUserReq struct {
	LID        string `json:"lid"`
	Name       string `json:"name"`
	Socket     string `json:"socket"`
	RemotePort string `json:"remote_port"`
}

func PostServerToUploadUserInfo(id, name, localPort, remotePort string) string {

	host := base.QueryConfigByKey("LnurlServerHost")
	targetUrl := fmt.Sprintf("http://%s/api/v1/lnurl/upload/user", host)
	requestJsonBytes, err := json.Marshal(UploadUserReq{
		LID:        id,
		Name:       name,
		Socket:     localPort,
		RemotePort: remotePort,
	})
	if err != nil {
		return ""
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", targetUrl, payload)
	if err != nil {
		return ""
	}

	username := base.QueryConfigByKey("BasicAuthUser")
	password := base.QueryConfigByKey("BasicAuthPass")
	req.SetBasicAuth(username, password)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ""
	}
	var resp UploadUserResp
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return ""
	}
	if resp.Msg != "" {
		return errors.New(resp.Msg).Error()
	}
	return resp.Data
}

func PostPhoneToAddInvoice(remotePort, amount string) string {

	return ""
}

type PayInvoiceReq struct {
	InvoiceType InvoiceType `json:"invoice_type"`
	AssetID     string      `json:"asset_id"`
	Amount      uint64      `json:"amount"`
	PubKey      string      `json:"pub_key"`
	Memo        string      `json:"memo"`
}

func PostServerToPayByPhoneAddInvoice(lnu string, invoiceType int, assetID string, amount int, pubkey string, memo string) string {
	targetUrl := Decode(lnu)

	requestJsonBytes, err := json.Marshal(PayInvoiceReq{
		InvoiceType: InvoiceType(invoiceType),
		AssetID:     assetID,
		Amount:      uint64(amount),
		PubKey:      pubkey,
		Memo:        memo,
	})
	if err != nil {
		return ""
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", targetUrl, payload)
	if err != nil {
		return ""
	}

	username := base.QueryConfigByKey("BasicAuthUser")
	password := base.QueryConfigByKey("BasicAuthPass")
	req.SetBasicAuth(username, password)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ""
	}
	var resp InvoiceResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return ""
	}
	if resp.Msg != "" {
		return errors.New(resp.Msg).Error()
	}
	return resp.Data
}

func PostServerToRequestInvoice(lnu string, invoiceType int, assetID string, amount int, pubkey string, memo string) (string, error) {
	targetUrl := Decode(lnu)

	requestJsonBytes, err := json.Marshal(PayInvoiceReq{
		InvoiceType: InvoiceType(invoiceType),
		AssetID:     assetID,
		Amount:      uint64(amount),
		PubKey:      pubkey,
		Memo:        memo,
	})
	if err != nil {
		return "", errors.Wrap(err, "json.Marshal")
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", targetUrl, payload)
	if err != nil {
		return "", errors.Wrap(err, "http.NewRequest")
	}

	username := base.QueryConfigByKey("BasicAuthUser")
	password := base.QueryConfigByKey("BasicAuthPass")
	req.SetBasicAuth(username, password)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "http.DefaultClient.Do")
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "io.ReadAll")
	}
	var resp InvoiceResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", errors.Wrap(err, "json.Unmarshal")
	}
	if resp.Msg != "" {
		return "", errors.New(resp.Msg)
	}
	return resp.Data, nil
}

func PostServerToRequestLnurl(id, name, localPort, remotePort string) (string, error) {

	host := base.QueryConfigByKey("LnurlServerHost")
	targetUrl := fmt.Sprintf("http://%s/api/v1/lnurl/upload/user", host)
	requestJsonBytes, err := json.Marshal(UploadUserReq{
		LID:        id,
		Name:       name,
		Socket:     localPort,
		RemotePort: remotePort,
	})
	if err != nil {
		return "", errors.Wrap(err, "json.Marshal")
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", targetUrl, payload)
	if err != nil {
		return "", errors.Wrap(err, "http.NewRequest")
	}

	username := base.QueryConfigByKey("BasicAuthUser")
	password := base.QueryConfigByKey("BasicAuthPass")
	req.SetBasicAuth(username, password)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "http.DefaultClient.Do")
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "io.ReadAll")
	}
	var resp UploadUserResp
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", errors.Wrap(err, "json.Unmarshal")
	}
	if resp.Msg != "" {
		return "", errors.New(resp.Msg)
	}
	return resp.Data, nil
}
