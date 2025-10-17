package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/wallet/base"
	"io"
	"net/http"
)

func SimplifyTransactions2(address string, responses *GetAddressTransactionsResponse) []*TransactionsSimplified {
	var simplified []*TransactionsSimplified
	for _, transaction := range *responses {
		var simplifiedTx TransactionsSimplified
		simplifiedTx.Txid = transaction.Txid
		simplifiedTx.BlockTime = transaction.Status.BlockTime
		simplifiedTx.FeeRate = RoundToDecimalPlace(float64(transaction.Fee)/(float64(transaction.Weight)/4), 2)
		simplifiedTx.Fee = transaction.Fee
		blockHeight := BlockTipHeight()
		if !transaction.Status.Confirmed {
			simplifiedTx.ConfirmedBlocks = 0
		} else if blockHeight == 0 {
			fmt.Println("block height is zero")
			simplifiedTx.ConfirmedBlocks = 0
		} else {
			simplifiedTx.ConfirmedBlocks = BlockTipHeight() - transaction.Status.BlockHeight
		}
		for _, vin := range transaction.Vin {
			simplifiedTx.Vin = append(simplifiedTx.Vin, struct {
				ScriptpubkeyAddress string `json:"scriptpubkey_address"`
				Value               int    `json:"value"`
			}{
				ScriptpubkeyAddress: vin.Prevout.ScriptpubkeyAddress,
				Value:               vin.Prevout.Value,
			})
			if vin.Prevout.ScriptpubkeyAddress == address {
				simplifiedTx.BalanceResult -= vin.Prevout.Value
			}
		}
		for _, vout := range transaction.Vout {
			simplifiedTx.Vout = append(simplifiedTx.Vout, struct {
				ScriptpubkeyAddress string `json:"scriptpubkey_address"`
				Value               int    `json:"value"`
			}{
				ScriptpubkeyAddress: vout.ScriptpubkeyAddress,
				Value:               vout.Value,
			})
			if vout.ScriptpubkeyAddress == address {
				simplifiedTx.BalanceResult += vout.Value
			}
		}
		simplified = append(simplified, &simplifiedTx)
	}
	return simplified
}

func PcGetAddressTransactionsByMempool(address string) ([]*TransactionsSimplified, error) {
	var targetUrl string
	switch base.NetWork {
	case base.UseMainNet:
		targetUrl = "https://mempool.space/api/address/" + address + "/txs"
	case base.UseTestNet:
		targetUrl = "https://mempool.space/testnet/api/address/" + address + "/txs"
	default:
		targetUrl = "https://mempool.space/api/address/" + address + "/txs"
	}

	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", targetUrl, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, AppendErrorInfo(err, "http.Get")

	}
	var getAddressTransactionsResponse GetAddressTransactionsResponse
	if err = json.Unmarshal(body, &getAddressTransactionsResponse); err != nil {
		return nil, err
	}
	transactions := SimplifyTransactions2(address, &getAddressTransactionsResponse)
	return transactions, nil
}
