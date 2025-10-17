package pcapi

import (
	"github.com/wallet/api"
)

func GetListEligibleCoins(assetId string) ([]*api.EligibleCoin, error) {
	return api.PcGetListEligibleCoins(assetId)
}

func CreateSellOrderSign(assetId string, assetNum uint64, price int64, anchorPoint string, internalKey string) (string, error) {
	return api.PcCreateSellOrderSign(assetId, assetNum, price, anchorPoint, internalKey)
}

func BuySOrderSign(signedSellOrder string, feeRate uint64) (string, error) {
	return api.PcBuySOrderSign(signedSellOrder, feeRate)
}

func PublishSOrderTx(signedBoughtSOrder string) (string, error) {
	return api.PcPublishSOrderTx(signedBoughtSOrder)
}

func AllowFederationSyncInsertAndExport() error {
	return api.PcAllowFederationSyncInsertAndExport()
}

func InsertProofAndRegisterTransfer(assetId string, signedBoughtSOrder string, lastProofStr string) (string, error) {
	return api.PcInsertProofAndRegisterTransfer(assetId, signedBoughtSOrder, lastProofStr)
}
