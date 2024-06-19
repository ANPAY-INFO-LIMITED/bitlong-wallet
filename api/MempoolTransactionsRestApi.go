package api

import (
	"encoding/json"
	"fmt"
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
	}
	response, err := http.Get(targetUrl)
	if err != nil {
		fmt.Printf("%s http.Get :%v\n", GetTimeNow(), err)
		return nil, err
	}
	bodyBytes, _ := io.ReadAll(response.Body)
	var transactionsResponse TransactionsResponse
	if err := json.Unmarshal(bodyBytes, &transactionsResponse); err != nil {
		fmt.Printf("%s GTBM json.Unmarshal :%v\n", GetTimeNow(), err)
		return nil, err
	}
	return &transactionsResponse, nil
}

func GetTransactionByMempool(transaction string) string {
	response, err := getTransactionByMempool(transaction)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, "Unmarshal response body fail.", "")
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
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
