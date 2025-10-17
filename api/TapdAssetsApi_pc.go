package api

import (
	"github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	"github.com/pkg/errors"
)

func AssetLeavesSpecified2(id string, proofType string) (*universerpc.AssetLeafResponse, error) {
	var _proofType universerpc.ProofType
	if proofType == "issuance" || proofType == "ISSUANCE" || proofType == "PROOF_TYPE_ISSUANCE" {
		_proofType = universerpc.ProofType_PROOF_TYPE_ISSUANCE
	} else if proofType == "transfer" || proofType == "TRANSFER" || proofType == "PROOF_TYPE_TRANSFER" {
		_proofType = universerpc.ProofType_PROOF_TYPE_TRANSFER
	} else {
		_proofType = universerpc.ProofType_PROOF_TYPE_UNSPECIFIED
	}
	response, err := assetLeaves(false, id, _proofType)
	if err != nil {
		return nil, errors.Wrap(err, "assetLeaves")
	}
	return response, nil
}

func assetLeavesTransfer2(assetID string) (*[]AssetTransferLeave, error) {
	response, err := AssetLeavesSpecified2(assetID, universerpc.ProofType_PROOF_TYPE_TRANSFER.String())
	if err != nil {
		return nil, errors.Wrap(err, "AssetLeavesSpecified2")
	}
	return ProcessAssetTransferLeave(response), nil
}

func PcUploadAssetManagedUtxos(token string, deviceId string) error {
	return ListUtxosAndPostToSetAssetManagedUtxos(token, deviceId)
}

func PcUploadAssetLocalMintHistory(token string, deviceId string) error {
	return ListBatchesAndPostToSetAssetLocalMintHistories(token, deviceId)
}

func PcUploadAssetListInfo(token string, deviceId string) error {
	isTokenValid, err := IsTokenValid(token)
	if err != nil {
		return errors.Wrap(err, "IsTokenValid")
	} else if !isTokenValid {
		return errors.New("token is invalid")
	}
	assets, err := ListAssetsProcessed(true, false, false)
	if err != nil {
		return errors.Wrap(err, "ListAssetsProcessed")
	}
	zeroAmountAssetLists, err := GetZeroAmountAssetListSlice(token, assets)
	if err != nil {
		return errors.Wrap(err, "GetZeroAmountAssetListSlice")
	}
	listAssetsResponseSlice := AssetBalanceInfosToListAssetsResponseSlice(zeroAmountAssetLists)
	setListAssetsResponseSlice := append(*assets, *listAssetsResponseSlice...)
	requests := ListAssetsResponseSliceToAssetListSetRequests(&setListAssetsResponseSlice, deviceId)
	_, err = PostToSetAssetListInfo(requests, token)
	if err != nil {
		return errors.Wrap(err, "PostToSetAssetListInfo")
	}
	return nil
}

func PcUploadAddrReceives(token string, deviceId string) error {
	events, err := AddrReceivesAndGetEvents(deviceId)
	if err != nil {
		return errors.Wrap(err, "AddrReceivesAndGetEvents")
	}
	err = PostToSetAddrReceivesEvents(token, events)
	if err != nil {
		return errors.Wrap(err, "PostToSetAddrReceivesEvents")
	}
	return nil
}

func PcUploadAssetTransfer(token string, deviceId string) error {
	transfers, err := ListTransfersAndGetProcessedResponse(token, deviceId)
	if err != nil {
		return errors.Wrap(err, "ListTransfersAndGetProcessedResponse")
	}
	if transfers == nil || len(*transfers) == 0 {
		return nil
	}
	_, err = PostToSetAssetTransfer(token, transfers)
	if err != nil {
		return errors.Wrap(err, "PostToSetAssetTransfer")
	}
	return nil
}

func PcUploadAssetBalanceInfo(token string, deviceId string) error {
	isTokenValid, err := IsTokenValid(token)
	if err != nil {
		return errors.Wrap(err, "IsTokenValid")
	} else if !isTokenValid {
		return errors.New("token is invalid")
	}
	balances, err := ListBalancesAndProcess()
	if err != nil {
		return errors.Wrap(err, " ListBalancesAndProcess")
	}
	zeroBalances, err := GetZeroBalanceAssetBalanceSlice(token, balances)
	if err != nil {
		return errors.Wrap(err, "GetZeroBalanceAssetBalanceSlice")
	}
	zeroListBalance := AssetBalanceInfosToListBalanceInfos(zeroBalances)
	setBalances := append(*balances, *zeroListBalance...)
	requests := ListBalanceInfosToAssetBalanceSetRequests(&setBalances, deviceId, false)
	_, err = PostToSetAssetBalanceInfo(requests, token)
	if err != nil {
		return errors.Wrap(err, " PostToSetAssetBalanceInfo")
	}
	return nil
}

func PcGetBtcTransferInInfosJsonResult(token string) (*[]BtcTransferInInfoSimplified, error) {
	response, err := GetBtcTransferInInfos(token)
	if err != nil {
		return nil, errors.Wrap(err, "GetBtcTransferInInfos")
	}
	return response, nil
}

func PcGetBtcTransferOutInfosJsonResult(token string) (*[]BtcTransferOutInfoSimplified, error) {
	response, err := GetBtcTransferOutInfos(token)
	if err != nil {
		return nil, errors.Wrap(err, "GetBtcTransferOutInfos")
	}
	return response, nil
}

func PcUploadAssetBalanceHistories(token string) error {
	err := GetAndUploadAssetBalanceHistories(token)
	if err != nil {
		return errors.Wrap(err, "GetAndUploadAssetBalanceHistories")
	}
	return nil
}

func PcAutoMintReserved(token string, deviceId string) ([]string, error) {
	result, err := GetOwnFairLaunchInfoIssuedSimplifiedAndExecuteMintReserved(token, deviceId)
	if err != nil {
		return nil, errors.Wrap(err, "GetOwnFairLaunchInfoIssuedSimplifiedAndExecuteMintReserved")
	}
	return result, nil
}

func PcGetWalletBalanceTotalValue(token string) (float64, error) {
	return GetWalletBalanceCalculatedTotalValue(token)
}
