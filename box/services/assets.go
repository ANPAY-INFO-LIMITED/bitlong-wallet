package services

import (
	"encoding/hex"
	"strconv"

	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/pkg/errors"
	"github.com/wallet/api"
	"github.com/wallet/box/rpc"
	"github.com/wallet/box/st"
)

func AssetTransferIn(assetId string) ([]*api.AddrEvent, error) {
	var t rpc.Tap
	receive, err := t.AddrReceives()
	if err != nil {
		return nil, errors.Wrap(err, "t.AddrReceives")
	}

	return api.AssetTransferIn(receive, assetId)
}

func AssetTransferOut(assetId string) ([]*api.AssetTransferSimplified, error) {
	var t rpc.Tap
	transfer, err := t.ListTransfers()
	if err != nil {
		return nil, errors.Wrap(err, "t.ListTransfers")
	}

	return api.AssetTransferOut(transfer, assetId)
}

func AssetUtxo(assetId string) ([]*api.ManagedUtxo, error) {
	err := updateToken()
	if err != nil {
		return nil, errors.Wrap(err, "updateToken")
	}
	token := st.Token()

	var t rpc.Tap
	utxo, err := t.ListUtxos()
	if err != nil {
		return nil, errors.Wrap(err, "t.ListUtxos")
	}

	return api.AssetUtxo(utxo, token, assetId)
}

type ListAssetBalanceInfo struct {
	GenesisPoint string `json:"genesis_point"`
	Name         string `json:"name"`
	MetaHash     string `json:"meta_hash"`
	AssetID      string `json:"asset_id"`
	AssetType    string `json:"asset_type"`
	OutputIndex  int    `json:"output_index"`
	Version      int    `json:"version"`
	Balance      string `json:"balance"`
}

func ProcessListBalancesResponse(response *taprpc.ListBalancesResponse) []*ListAssetBalanceInfo {
	var listAssetBalanceInfos []*ListAssetBalanceInfo
	for _, balance := range response.AssetBalances {
		listAssetBalanceInfos = append(listAssetBalanceInfos, &ListAssetBalanceInfo{
			GenesisPoint: balance.AssetGenesis.GenesisPoint,
			Name:         balance.AssetGenesis.Name,
			MetaHash:     hex.EncodeToString(balance.AssetGenesis.MetaHash),
			AssetID:      hex.EncodeToString(balance.AssetGenesis.AssetId),
			AssetType:    balance.AssetGenesis.AssetType.String(),
			OutputIndex:  int(balance.AssetGenesis.OutputIndex),
			Version:      -1,
			Balance:      strconv.FormatUint(balance.Balance, 10),
		})
	}
	return listAssetBalanceInfos
}

func ExcludeListBalancesResponseCollectible(listAssetBalanceInfos []*ListAssetBalanceInfo) []*ListAssetBalanceInfo {
	var listAssetBalances []*ListAssetBalanceInfo
	for _, balance := range listAssetBalanceInfos {
		if balance.AssetType == taprpc.AssetType_NORMAL.String() {
			listAssetBalances = append(listAssetBalances, balance)
		}
	}
	return listAssetBalances
}
