package api

import (
	"encoding/json"
	"github.com/wallet/base"
	"io"
	"net/http"
)

type DifficultyResponse struct {
	ProgressPercent       float64 `json:"progressPercent"`
	DifficultyChange      float64 `json:"difficultyChange"`
	EstimatedRetargetDate int64   `json:"estimatedRetargetDate"`
	RemainingBlocks       int     `json:"remainingBlocks"`
	RemainingTime         int     `json:"remainingTime"`
	PreviousRetarget      float64 `json:"previousRetarget"`
	PreviousTime          int     `json:"previousTime"`
	NextRetargetHeight    int     `json:"nextRetargetHeight"`
	TimeAvg               int     `json:"timeAvg"`
	AdjustedTimeAvg       int     `json:"adjustedTimeAvg"`
	TimeOffset            int     `json:"timeOffset"`
	ExpectedBlocks        float64 `json:"expectedBlocks"`
}

func GetDifficultyAdjustmentByMempool() string {
	var targetUrl string
	switch base.NetWork {
	case base.UseMainNet:
		targetUrl = "https://mempool.space/api/v1/difficulty-adjustment"

	case base.UseTestNet:
		targetUrl = "https://mempool.space/testnet/api/v1/difficulty-adjustment"
	}
	response, err := http.Get(targetUrl)
	if err != nil {
		return MakeJsonErrorResult(HttpGetErr, "http get fail.", "")
	}
	bodyBytes, _ := io.ReadAll(response.Body)
	var difficultyResponse DifficultyResponse
	if err := json.Unmarshal(bodyBytes, &difficultyResponse); err != nil {
		return MakeJsonErrorResult(UnmarshalErr, "Unmarshal response body fail.", "")
	}
	return MakeJsonErrorResult(SUCCESS, "", difficultyResponse)
}
