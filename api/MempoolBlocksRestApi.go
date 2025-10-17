package api

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
		logrus.Errorln(errors.Wrap(err, "http.Get"))
		return 0
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Errorln(errors.Wrap(err, "io.ReadAll"))
		return 0
	}
	height, err := strconv.Atoi(string(bodyBytes))
	if err != nil {
		logrus.Errorln(errors.Wrap(err, "strconv.Atoi"))
		return 0
	}
	return height
}

func GetBlockTipHashByMempool() {}

func GetBlockTransactionIDByMempool() {}

func GetBlockTransactionIDsByMempool() {}

func GetBlockTransactionsByMempool() {}

func GetBlocksByMempool() {}

func GetBlocksBulkByMempool() {}
