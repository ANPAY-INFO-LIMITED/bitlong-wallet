package pcapi

import "github.com/wallet/api"

func GetAddressTransactionsByMempool(address string) ([]*api.TransactionsSimplified, error) {
	return api.PcGetAddressTransactionsByMempool(address)
}

func GetTransactionByMempool(txid string) (*api.TransactionsSimplified, error) {
	return api.PcGetTransactionByMempool(txid)
}
