package api

import (
	"bytes"
	"encoding/json"
	"github.com/wallet/base"
	"io"
	"net/http"
)

func GetChildrenPayforParentByMempool() {}

type TransactionsResponse struct {
	Txid     string `json:"txid"`
	Version  int    `json:"version"`
	Locktime int    `json:"locktime"`
	Vin      []struct {
		Txid    string `json:"txid"`
		Vout    int    `json:"vout"`
		Prevout struct {
			Scriptpubkey        string `json:"scriptpubkey"`
			ScriptpubkeyAsm     string `json:"scriptpubkey_asm"`
			ScriptpubkeyType    string `json:"scriptpubkey_type"`
			ScriptpubkeyAddress string `json:"scriptpubkey_address"`
			Value               int    `json:"value"`
		} `json:"prevout"`
		Scriptsig    string   `json:"scriptsig"`
		ScriptsigAsm string   `json:"scriptsig_asm"`
		Witness      []string `json:"witness"`
		IsCoinbase   bool     `json:"is_coinbase"`
		Sequence     int64    `json:"sequence"`
	} `json:"vin"`
	Vout []struct {
		Scriptpubkey        string `json:"scriptpubkey"`
		ScriptpubkeyAsm     string `json:"scriptpubkey_asm"`
		ScriptpubkeyType    string `json:"scriptpubkey_type"`
		ScriptpubkeyAddress string `json:"scriptpubkey_address"`
		Value               int    `json:"value"`
	} `json:"vout"`
	Size   int `json:"size"`
	Weight int `json:"weight"`
	Sigops int `json:"sigops"`
	Fee    int `json:"fee"`
	Status struct {
		Confirmed   bool   `json:"confirmed"`
		BlockHeight int    `json:"block_height"`
		BlockHash   string `json:"block_hash"`
		BlockTime   int    `json:"block_time"`
	} `json:"status"`
}

func getTransactionByMempool(transaction string) (*TransactionsResponse, error) {
	var targetUrl string
	switch base.NetWork {
	case base.UseMainNet:
		targetUrl = "https://mempool.space/api/tx/" + transaction
	case base.UseTestNet:
		targetUrl = "https://mempool.space/testnet/api/tx/" + transaction
	default:
		targetUrl = "https://mempool.space/api/tx/" + transaction
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

	var transactionsResponse TransactionsResponse
	if err := json.Unmarshal(body, &transactionsResponse); err != nil {
		return nil, err
	}
	return &transactionsResponse, nil
}

func TransactionsResponseToTransactionsSimplified(transactionsResponse *TransactionsResponse) *TransactionsSimplified {
	var transactionsSimplified TransactionsSimplified
	confirmedBlocks := BlockTipHeight() - transactionsResponse.Status.BlockHeight
	if confirmedBlocks < 0 {
		confirmedBlocks = 0
	}
	var vin []TransactionsSimplifiedVin
	for _, in := range transactionsResponse.Vin {
		vin = append(vin, TransactionsSimplifiedVin{
			ScriptpubkeyAddress: in.Prevout.ScriptpubkeyAddress,
			Value:               in.Prevout.Value,
		})
	}
	var vout []TransactionsSimplifiedVout
	var balanceResult int
	for _, out := range transactionsResponse.Vout {
		balanceResult += out.Value
		vout = append(vout, TransactionsSimplifiedVout{
			ScriptpubkeyAddress: out.ScriptpubkeyAddress,
			Value:               out.Value,
		})
	}
	transactionsSimplified.Txid = transactionsResponse.Txid
	transactionsSimplified.Vin = vin
	transactionsSimplified.Vout = vout
	transactionsSimplified.BlockTime = transactionsResponse.Status.BlockTime
	transactionsSimplified.BalanceResult = balanceResult
	transactionsSimplified.FeeRate = RoundToDecimalPlace(float64(transactionsResponse.Fee)/(float64(transactionsResponse.Weight)/4), 2)
	transactionsSimplified.Fee = transactionsResponse.Fee
	transactionsSimplified.ConfirmedBlocks = confirmedBlocks
	return &transactionsSimplified
}

// GetTransactionByMempool
// @Description: Get transactions simplified info by txid
func GetTransactionByMempool(txid string) string {
	response, err := getTransactionByMempool(txid)
	if err != nil {
		return MakeJsonErrorResult(getTransactionByMempoolErr, err.Error(), nil)
	}
	result := TransactionsResponseToTransactionsSimplified(response)
	return MakeJsonErrorResult(SUCCESS, "", result)
}

func GetTransactionHexByMempool() {}

func GetTransactionMerkleblockProofByMempool() {}

func GetTransactionMerkleProofByMempool() {}

func GetTransactionOutspendByMempool() {}

func GetTransactionOutspendsByMempool() {}

func GetTransactionRawByMempool() {}

func GetTransactionRBFHistoryByMempool() {}

func GetTransactionStatusByMempool() {}

func GetTransactionTimesByMempool() {}

func PostTransactionByMempool() {}
