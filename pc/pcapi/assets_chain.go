package pcapi

import (
	"github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	"github.com/wallet/api"
)

func ListNormalBalances() ([]*api.ListAssetBalanceInfo2, error) {
	resp, err := api.PcListNormalBalances2()
	if err != nil {
		return nil, err
	}
	var balances []*api.ListAssetBalanceInfo2
	if resp == nil {
		return balances, nil
	}
	for _, b := range *resp {
		balances = append(balances, &b)
	}
	return balances, nil
}

func CheckAssetIssuanceIsLocal(assetId string) (*api.IsLocalResult, error) {
	return api.PcCheckAssetIssuanceIsLocal(assetId)
}

func AddrReceives(assetId string) ([]*api.AddrEvent, error) {
	resp, err := api.PcAddrReceives(assetId)
	if err != nil {
		return nil, err
	}
	var addrReceives []*api.AddrEvent
	if resp == nil {
		return addrReceives, nil
	}
	for _, r := range *resp {
		addrReceives = append(addrReceives, &r)
	}
	return addrReceives, nil
}

func QueryAssetTransfers(assetId string) ([]*api.AssetTransferSimplified, error) {
	resp, err := api.PcQueryAssetTransfers(assetId)
	if err != nil {
		return nil, err
	}
	var assetTransfers []*api.AssetTransferSimplified
	if resp == nil {
		return assetTransfers, nil
	}
	for _, t := range *resp {
		assetTransfers = append(assetTransfers, &t)
	}
	return assetTransfers, nil
}

func AssetUtxos(token string, assetId string) ([]*api.ManagedUtxo, error) {
	resp, err := api.PcAssetUtxos(token, assetId)
	if err != nil {
		return nil, err
	}
	var assetUtxos []*api.ManagedUtxo
	if resp == nil {
		return assetUtxos, nil
	}
	for _, u := range *resp {
		assetUtxos = append(assetUtxos, &u)
	}
	return assetUtxos, nil
}

func NewAddr(assetId string, amt int, token string, deviceId string) (*api.QueriedAddr, error) {
	return api.PcNewAddr(assetId, amt, token, deviceId)
}

func QueryAddrs(assetId string) ([]*api.QueriedAddr, error) {
	return api.PcQueryAddrs(assetId)
}

func SendAssets(jsonAddrs string, feeRate int64, token string, deviceId string) (string, error) {
	return api.PcSendAssets(jsonAddrs, feeRate, token, deviceId)
}

func ListNftGroups() ([]*api.NftGroup, error) {
	return api.PcListNftGroups()
}

func ListNonGroupNftAssets() ([]*api.ListAssetsResponse, error) {
	return api.PcListNonGroupNftAssets()
}

func GetSpentNftAssets() ([]*api.ListAssetsSimplifiedResponse, error) {
	resp, err := api.PcGetSpentNftAssets()
	if err != nil {
		return nil, err
	}
	var assets []*api.ListAssetsSimplifiedResponse
	if resp == nil {
		return assets, nil
	}
	for _, a := range *resp {
		assets = append(assets, &a)
	}
	return assets, nil
}

func MintAsset(name string, assetTypeIsCollectible bool, description string, imagePath string, groupName string, amount int, decimalDisplay int, newGroupedAsset bool) (*api.PendingBatch, error) {
	return api.PcMintAsset(name, assetTypeIsCollectible, description, imagePath, groupName, amount, decimalDisplay, newGroupedAsset)
}

func AddGroupAsset(name string, assetTypeIsCollectible bool, description string, imagePath string, groupName string, amount int, groupKey string) (*api.PendingBatch, error) {
	return api.PcAddGroupAsset(name, assetTypeIsCollectible, description, imagePath, groupName, amount, groupKey)
}

func CancelBatch() error {
	return api.PcCancelBatch()
}

func FinalizeBatch(feeRate int, token string, deviceId string) (*api.PendingBatch, error) {
	return api.PcFinalizeBatch(feeRate, token, deviceId)
}

func GetIssuanceTransactionFee(token string, feeRate int) (int, error) {
	return api.PcGetIssuanceTransactionFee(token, feeRate)
}

func GetAssetInfo(assetId string) (*api.AssetInfo, error) {
	return api.PcGetAssetInfo(assetId)
}

func GetWalletBalanceTotalValue(token string) (float64, error) {
	return api.PcGetWalletBalanceTotalValue(token)
}

func SyncUniverse(universeHost string, assetId string, isTransfer bool) (*universerpc.SyncResponse, error) {
	return api.PcSyncUniverse(universeHost, assetId, isTransfer)
}
