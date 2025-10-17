package api

import "github.com/pkg/errors"

func PcGetTransactionByMempool(txid string) (*TransactionsSimplified, error) {
	response, err := getTransactionByMempool(txid)
	if err != nil {
		return nil, errors.Wrap(err, "getTransactionByMempool")
	}
	result := TransactionsResponseToTransactionsSimplified(response)
	return result, nil
}
