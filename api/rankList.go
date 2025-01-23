package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func GetCustodyAssetsRankList(AssetId, token string, page, pageSize int) string {
	return custodyRankList(AssetId, token, page, pageSize)
}

func custodyRankList(AssetId, token string, page, pageSize int) string {
	serverDomainOrSocket := Cfg.BtlServerHost
	quest := struct {
		AssetId  string `json:"assetId"`
		Page     int    `json:"page"`
		PageSize int    `json:"pageSize"`
	}{
		AssetId:  AssetId,
		Page:     page,
		PageSize: pageSize,
	}
	url := "http://" + serverDomainOrSocket + "/account_asset/balance/rankList"
	requestJsonBytes, err := json.Marshal(quest)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), "")
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), "")
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), "")
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), "")
	}
	return string(body)
}
