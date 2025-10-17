package api

import (
	"encoding/hex"
	"github.com/lightninglabs/taproot-assets/proof"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	"github.com/pkg/errors"
	"github.com/wallet/base"
	"github.com/wallet/service/rpcclient"
)

type AssetInfo = struct {
	AssetId        string `json:"asset_Id"`
	Name           string `json:"name"`
	Point          string `json:"point"`
	AssetType      string `json:"assetType"`
	GroupName      string `json:"group_name"`
	GroupKey       string `json:"group_key"`
	Amount         uint64 `json:"amount"`
	Meta           string `json:"meta"`
	CreateHeight   int64  `json:"create_height"`
	CreateTime     int64  `json:"create_time"`
	Universe       string `json:"universe"`
	DecimalDisplay uint32 `json:"decimal_display"`
}

func PcGetAssetInfo(id string) (*AssetInfo, error) {
	root, err := rpcclient.PcQueryAssetRoots(id)
	if err != nil {
		return nil, errors.Wrap(err, "PcQueryAssetRoots")
	}
	queryId := id
	isGroup := false
	if groupKey, ok := root.IssuanceRoot.Id.Id.(*universerpc.ID_GroupKey); ok {
		isGroup = true
		queryId = hex.EncodeToString(groupKey.GroupKey)
	}
	response, err := assetLeaves(isGroup, queryId, universerpc.ProofType_PROOF_TYPE_ISSUANCE)
	if err != nil {
		return nil, errors.Wrap(err, "assetLeaves")
	}
	if response.Leaves == nil {
		return nil, errors.New("response.Leaves is nil")
	}
	var assetinfo *taprpc.Asset
	var blob proof.Blob
	for index, leaf := range response.Leaves {
		if hex.EncodeToString(leaf.Asset.AssetGenesis.GetAssetId()) == id {
			blob = response.Leaves[index].Proof
			assetinfo = leaf.Asset
			break
		}
	}
	if len(blob) == 0 {
		return nil, errors.New("response.Leaves[index].Proof length is 0")
	}
	p, _ := blob.AsSingleProof()
	assetId := p.Asset.ID().String()
	assetName := p.Asset.Tag
	assetPoint := p.Asset.FirstPrevOut.String()
	assetType := p.Asset.Type.String()
	amount := p.Asset.Amount
	createHeight := p.BlockHeight
	createTime := p.BlockHeader.Timestamp
	var (
		newMeta Meta
		m       = ""
	)
	if p.MetaReveal != nil {
		m = string(p.MetaReveal.Data)
	}
	newMeta.GetMetaFromStr(m)
	var assetInfo = AssetInfo{
		AssetId:      assetId,
		Name:         assetName,
		Point:        assetPoint,
		AssetType:    assetType,
		GroupName:    newMeta.GroupName,
		Amount:       amount,
		Meta:         newMeta.Description,
		CreateHeight: int64(createHeight),
		CreateTime:   createTime.Unix(),
		Universe:     "localhost",
	}
	if isGroup {
		assetInfo.GroupKey = queryId
	}
	if assetinfo != nil && assetinfo.DecimalDisplay != nil {
		assetInfo.DecimalDisplay = assetinfo.DecimalDisplay.DecimalDisplay
	} else {
		assetInfo.DecimalDisplay = 0
	}

	return &assetInfo, nil
}

func PcSyncUniverse(universeHost string, assetId string, isTransfer bool) (*universerpc.SyncResponse, error) {
	var targets []*universerpc.SyncTarget
	var proofType universerpc.ProofType
	if isTransfer {
		proofType = universerpc.ProofType_PROOF_TYPE_TRANSFER
	} else {
		proofType = universerpc.ProofType_PROOF_TYPE_ISSUANCE
	}
	universeID := &universerpc.ID{
		Id: &universerpc.ID_AssetIdStr{
			AssetIdStr: assetId,
		},
		ProofType: proofType,
	}
	if universeID != nil {
		targets = append(targets, &universerpc.SyncTarget{
			Id: universeID,
		})
	}
	var defaultHost string
	switch base.NetWork {
	case base.UseMainNet:
		defaultHost = "universe.lightning.finance:10029"
	case base.UseTestNet:
		defaultHost = "testnet.universe.lightning.finance:10029"
	}
	if universeHost == "" {
		universeHost = defaultHost
	}
	resp, err := syncUniverse(universeHost, targets, universerpc.UniverseSyncMode_SYNC_FULL)
	if err != nil {
		return nil, errors.Wrap(err, "syncUniverse")
	}
	return resp, nil
}
