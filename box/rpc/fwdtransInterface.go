package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/masterzen/azure-sdk-for-go/core/http"
	"github.com/wallet/box/st"
)

var host = "https://api.btc.microlinktoken.com:28095"

func FwdtApplyInvoice(assetId string, invoice string) (string, error) {

	url := host + "/v1/fwdt/applyInvoice"
	method := "POST"

	invoiceInfo := struct {
		AssetId       string
		TargetInvoice string
	}{
		AssetId:       assetId,
		TargetInvoice: invoice,
	}
	jsonData, err := json.Marshal(invoiceInfo)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return "", err
	}
	payload := bytes.NewReader(jsonData)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+st.Token())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	r := struct {
		Invoice string
		Error   string
	}{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return "", err
	}
	if r.Error != "" {
		return "", fmt.Errorf(r.Error)
	}
	return r.Invoice, nil
}

func FwdtPayInvoice(invoice string) error {
	url := host + "/v1/fwdt/fwdtPayment"
	method := "POST"

	creds := struct {
		MappingInvoice string
	}{
		MappingInvoice: invoice,
	}
	jsonData, err := json.Marshal(creds)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return err
	}
	payload := bytes.NewReader(jsonData)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+st.Token())

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if res.StatusCode != 200 {
		r := struct {
			Error string
		}{}
		err = json.Unmarshal(body, &r)
		if err != nil {
			return err
		}
		if r.Error != "" {
			return fmt.Errorf(r.Error)
		}
	}
	return nil
}

type InvoiceInfo struct {
	Invoice string  `json:"invoice"`
	AssetId string  `json:"assetId"`
	Amount  float64 `json:"amount"`
}

func CheckInvoiceIsCustody(invoice string) (*InvoiceInfo, error) {
	url := host + "/v1/fwdt/checkToCustodyInvoice"
	method := "POST"

	creds := struct {
		Invoice string
	}{
		Invoice: invoice,
	}
	jsonData, err := json.Marshal(creds)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil, err
	}
	payload := bytes.NewReader(jsonData)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+st.Token())

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if res.StatusCode != 200 {
		r := struct {
			Error string
		}{}
		err = json.Unmarshal(body, &r)
		if err != nil {
			return nil, err
		}
		if r.Error != "" {
			return nil, fmt.Errorf(r.Error)
		}
	}
	r := &InvoiceInfo{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func PayToCustody(keySendHash string, invoice string) error {
	url := host + "/v1/fwdt/payToCustodyInvoice"
	method := "POST"

	creds := struct {
		KeySendHash string `json:"keySendHash"`
		Invoice     string `json:"invoice"`
	}{
		KeySendHash: keySendHash,
		Invoice:     invoice,
	}
	jsonData, err := json.Marshal(creds)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return err
	}
	payload := bytes.NewReader(jsonData)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+st.Token())

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if res.StatusCode != 200 {
		r := struct {
			Error string
		}{}
		err = json.Unmarshal(body, &r)
		if err != nil {
			return err
		}
		if r.Error != "" {
			return fmt.Errorf(r.Error)
		}
	}
	return nil
}
