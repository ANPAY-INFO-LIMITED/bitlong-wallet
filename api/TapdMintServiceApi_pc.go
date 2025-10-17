package api

import (
	"context"
	"encoding/hex"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/lightninglabs/taproot-assets/taprpc/mintrpc"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/service/apiConnect"
)

func mintAsset2(assetVersionIsV1 bool, assetTypeIsCollectible bool, name string, assetMetaData string, AssetMetaTypeIsJsonNotOpaque bool, amount int, newGroupedAsset bool, groupedAsset bool, groupKey string, groupAnchor string, decimalDisplay int, shortResponse bool) (*PendingBatch, error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()
	mc := mintrpc.NewMintClient(conn)
	var _assetVersion taprpc.AssetVersion
	if assetVersionIsV1 {
		_assetVersion = taprpc.AssetVersion_ASSET_VERSION_V1
	} else {
		_assetVersion = taprpc.AssetVersion_ASSET_VERSION_V0
	}
	var _assetType taprpc.AssetType
	if assetTypeIsCollectible {
		_assetType = taprpc.AssetType_COLLECTIBLE
	} else {
		_assetType = taprpc.AssetType_NORMAL
	}
	_assetMetaDataByteSlice := []byte(assetMetaData)
	var _assetMetaType taprpc.AssetMetaType
	if AssetMetaTypeIsJsonNotOpaque {
	} else {
		_assetMetaType = taprpc.AssetMetaType_META_TYPE_OPAQUE
	}
	_groupKeyByteSlices, err := hex.DecodeString(groupKey)
	if err != nil {
		return nil, errors.Wrap(err, "hex.DecodeString")
	}
	if decimalDisplay < 0 || decimalDisplay > 100 {
		return nil, errors.New("invalid decimal display")
	}
	request := &mintrpc.MintAssetRequest{
		Asset: &mintrpc.MintAsset{
			AssetVersion: _assetVersion,
			AssetType:    _assetType,
			Name:         name,
			AssetMeta: &taprpc.AssetMeta{
				Data: _assetMetaDataByteSlice,
				Type: _assetMetaType,
			},
			Amount:          uint64(amount),
			NewGroupedAsset: newGroupedAsset,
			GroupedAsset:    groupedAsset,
			GroupKey:        _groupKeyByteSlices,
			GroupAnchor:     groupAnchor,
			DecimalDisplay:  uint32(decimalDisplay),
		},
		ShortResponse: shortResponse,
	}
	response, err := mc.MintAsset(context.Background(), request)
	if err != nil {
		return nil, errors.Wrap(err, "mc.MintAsset")
	}
	return MintAssetResponseToPendingBatch(response), nil
}

func finalizeBatch2(shortResponse bool, feeRate int, token string, deviceId string) (*PendingBatch, error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()
	mc := mintrpc.NewMintClient(conn)
	request := &mintrpc.FinalizeBatchRequest{
		ShortResponse: shortResponse,
		FeeRate:       uint32(feeRate),
	}
	response, err := mc.FinalizeBatch(context.Background(), request)
	if err != nil {
		return nil, errors.Wrap(err, "mc.FinalizeBatch")
	}
	err = UploadAssetLocalMints(token, deviceId, response)
	if err != nil {
		logrus.Infoln("", err)
	}
	err = SetGroupFirstAssetMeta(token, deviceId, response)
	if err != nil {
		logrus.Infoln("", err)
	}
	return FinalizeBatchResponseToPendingBatch(response), nil
}

func PcMintAsset(name string, assetTypeIsCollectible bool, description string, imagePath string, groupName string, amount int, decimalDisplay int, newGroupedAsset bool) (*PendingBatch, error) {
	assetMetaData := NewMeta(description)
	err := assetMetaData.LoadImageFile(imagePath)
	if err != nil {
		return nil, errors.Wrap(err, "LoadImageFile")
	}
	if groupName != "" {
		assetMetaData.GroupName = groupName
	}
	Metastr := assetMetaData.ToJsonStr()
	return mintAsset2(false, assetTypeIsCollectible, name, Metastr, false, amount, newGroupedAsset, false, "", "", decimalDisplay, false)
}

func PcFinalizeBatch(feeRate int, token string, deviceId string) (*PendingBatch, error) {
	if feeRate > FeeRateSatPerBToSatPerKw(500) {
		return nil, errors.New("fee rate exceeds max(500)")
	}
	return finalizeBatch2(false, feeRate, token, deviceId)
}

func PcAddGroupAsset(name string, assetTypeIsCollectible bool, description string, imagePath string, groupName string, amount int, groupKey string) (*PendingBatch, error) {
	assetMetaData := NewMeta(description)
	err := assetMetaData.LoadImageFile(imagePath)
	if err != nil {
		return nil, errors.Wrap(err, "LoadImageFile")
	}
	if groupName != "" {
		assetMetaData.GroupName = groupName
	}
	Metastr := assetMetaData.ToJsonStr()
	return mintAsset2(false, assetTypeIsCollectible, name, Metastr, false, amount, false, true, groupKey, "", 0, false)
}

func PcCancelBatch() error {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()
	mc := mintrpc.NewMintClient(conn)
	request := &mintrpc.CancelBatchRequest{}
	_, err = mc.CancelBatch(context.Background(), request)
	if err != nil {
		return errors.Wrap(err, "mc.CancelBatch")
	}
	return nil
}
