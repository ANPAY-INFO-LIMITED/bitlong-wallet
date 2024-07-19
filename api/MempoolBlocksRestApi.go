package api

import (
	"github.com/wallet/base"
	"io"
	"net/http"
	"strconv"
)

func GetBlockByMempoolByMempool() {}

func GetBlockHeaderByMempool() {}

func GetBlockHeightByMempool() {}

func GetBlockTimestampByMempool() {}

func GetBlockRawByMempool() {}

func GetBlockStatusByMempool() {}

func GetBlockTipHeightByMempool() string {
	var targetUrl string
	switch base.NetWork {
	case base.UseMainNet:
		targetUrl = "https://mempool.space/api/blocks/tip/height"
	case base.UseTestNet:
		targetUrl = "https://mempool.space/testnet/api/blocks/tip/height"
	}
	response, err := http.Get(targetUrl)
	if err != nil {
		return MakeJsonErrorResult(HttpGetErr, "http get fail.", "")
	}
	bodyBytes, _ := io.ReadAll(response.Body)
	var height string
	height = string(bodyBytes)
	return MakeJsonErrorResult(SUCCESS, "", height)
}

// BlockTipHeight
// @dev: NOT STANDARD RESULT RETURN
func BlockTipHeight() int {
	var targetUrl string
	switch base.NetWork {
	case base.UseMainNet:
		targetUrl = "https://mempool.space/api/blocks/tip/height"
	case base.UseTestNet:
		targetUrl = "https://mempool.space/testnet/api/blocks/tip/height"
	}
	response, err := http.Get(targetUrl)
	if err != nil {
		return 0
	}
	bodyBytes, _ := io.ReadAll(response.Body)
	height, _ := strconv.Atoi(string(bodyBytes))
	return height
}

func GetBlockTipHashByMempool() {}

func GetBlockTransactionIDByMempool() {}

func GetBlockTransactionIDsByMempool() {}

func GetBlockTransactionsByMempool() {}

func GetBlocksByMempool() {}

func GetBlocksBulkByMempool() {}
