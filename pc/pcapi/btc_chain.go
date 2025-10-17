package pcapi

import "github.com/wallet/api"

func GetWalletBalance() (*api.WalletBalanceResponse, error) {
	return api.PcGetWalletBalance()
}

func GetBtcTransferInInfosJsonResult(token string) ([]*api.BtcTransferInInfoSimplified, error) {
	resp, err := api.PcGetBtcTransferInInfosJsonResult(token)
	if err != nil {
		return nil, err
	}
	var btcTransferInInfos []*api.BtcTransferInInfoSimplified
	if resp == nil {
		return btcTransferInInfos, nil
	}
	for _, i := range *resp {
		btcTransferInInfos = append(btcTransferInInfos, &i)
	}
	return btcTransferInInfos, nil
}

func GetBtcTransferOutInfosJsonResult(token string) ([]*api.BtcTransferOutInfoSimplified, error) {
	resp, err := api.PcGetBtcTransferOutInfosJsonResult(token)
	if err != nil {
		return nil, err
	}
	var btcTransferOutInfos []*api.BtcTransferOutInfoSimplified
	if resp == nil {
		return btcTransferOutInfos, nil
	}
	for _, o := range *resp {
		btcTransferOutInfos = append(btcTransferOutInfos, &o)
	}
	return btcTransferOutInfos, nil
}

func BtcUtxos(token string) ([]*api.ListUnspentUtxo, error) {
	resp, err := api.PcBtcUtxos(token)
	if err != nil {
		return nil, err
	}
	var btcUtxos []*api.ListUnspentUtxo
	if resp == nil {
		return btcUtxos, nil
	}
	for _, u := range *resp {
		btcUtxos = append(btcUtxos, &u)
	}
	return btcUtxos, nil
}

func GetNewAddress() (string, error) {
	return api.PcGetNewAddress()
}

func SendCoins(addr string, amount int64, feeRate int64, sendAll bool) (string, error) {
	return api.PcSendCoins(addr, amount, feeRate, sendAll)
}

func MergeUTXO(feeRate int64) (string, error) {
	return api.PcMergeUTXO(feeRate)
}
