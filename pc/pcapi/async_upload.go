package pcapi

import "github.com/wallet/api"

func UploadWalletBalance(token string, deviceId string) error {
	return api.PcUploadWalletBalance(token, deviceId)
}

func UploadAssetManagedUtxos(token string, deviceId string) error {
	return api.PcUploadAssetManagedUtxos(token, deviceId)
}

func UploadAssetLocalMintHistory(token string, deviceId string) error {
	return api.PcUploadAssetLocalMintHistory(token, deviceId)
}

func UploadAssetListInfo(token string, deviceId string) error {
	return api.PcUploadAssetListInfo(token, deviceId)
}

func UploadAddrReceives(token string, deviceId string) error {
	return api.PcUploadAddrReceives(token, deviceId)
}

func UploadAssetTransfer(token string, deviceId string) error {
	return api.PcUploadAssetTransfer(token, deviceId)
}

func UploadAssetBalanceInfo(token string, deviceId string) error {
	return api.PcUploadAssetBalanceInfo(token, deviceId)
}

func UploadAssetBalanceHistories(token string) error {
	return api.PcUploadAssetBalanceHistories(token)
}

func UploadBtcListUnspentUtxos(token string) error {
	return api.PcUploadBtcListUnspentUtxos(token)
}

func AutoMintReserved(token string, deviceId string) ([]string, error) {
	return api.PcAutoMintReserved(token, deviceId)
}
