package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/lightninglabs/taproot-assets/taprpc/mintrpc"
	"github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/base"
	"github.com/wallet/models"
	"github.com/wallet/service/apiConnect"
	"github.com/wallet/service/rpcclient"
	"gorm.io/gorm"
)

const (
	UniverseSocketMainnet = "132.232.109.84:8444"
	UniverseSocketRegtest = "132.232.109.84:8443"
)

type SimplifiedAssetsTransfer struct {
	TransferTimestamp  int                     `json:"transfer_timestamp"`
	AnchorTxHash       string                  `json:"anchor_tx_hash"`
	AnchorTxHeightHint int                     `json:"anchor_tx_height_hint"`
	AnchorTxChainFees  int                     `json:"anchor_tx_chain_fees"`
	Inputs             []AssetsTransfersInput  `json:"inputs"`
	Outputs            []AssetsTransfersOutput `json:"outputs"`
}

type AssetsTransfersInput struct {
	AnchorPoint string `json:"anchor_point"`
	AssetID     string `json:"asset_id"`
	Amount      int    `json:"amount"`
}

type AssetsTransfersOutputAnchor struct {
	Outpoint string `json:"outpoint"`
	Value    int    `json:"value"`
}

type AssetsTransfersOutput struct {
	Anchor           AssetsTransfersOutputAnchor
	ScriptKeyIsLocal bool   `json:"script_key_is_local"`
	Amount           int    `json:"amount"`
	OutputType       string `json:"output_type"`
	AssetVersion     string `json:"asset_version"`
}

func SimplifyAssetsTransfer() *[]SimplifiedAssetsTransfer {
	var simpleTransfers []SimplifiedAssetsTransfer
	response, _ := rpcclient.ListTransfers()
	for _, transfers := range response.Transfers {
		var inputs []AssetsTransfersInput
		for _, _input := range transfers.Inputs {
			inputs = append(inputs, AssetsTransfersInput{
				AnchorPoint: _input.AnchorPoint,
				AssetID:     hex.EncodeToString(_input.AssetId),
				Amount:      int(_input.Amount),
			})
		}
		var outputs []AssetsTransfersOutput
		for _, _output := range transfers.Outputs {
			outputs = append(outputs, AssetsTransfersOutput{
				Anchor: AssetsTransfersOutputAnchor{
					Outpoint: _output.Anchor.Outpoint,
					Value:    int(_output.Anchor.Value),
				},
				ScriptKeyIsLocal: _output.ScriptKeyIsLocal,
				Amount:           int(_output.Amount),
				OutputType:       _output.OutputType.String(),
				AssetVersion:     _output.AssetVersion.String(),
			})
		}
		simpleTransfers = append(simpleTransfers, SimplifiedAssetsTransfer{
			TransferTimestamp:  int(transfers.TransferTimestamp),
			AnchorTxHash:       hex.EncodeToString(transfers.AnchorTxHash),
			AnchorTxHeightHint: int(transfers.AnchorTxHeightHint),
			AnchorTxChainFees:  int(transfers.AnchorTxChainFees),
			Inputs:             inputs,
			Outputs:            outputs,
		})
	}
	return &simpleTransfers
}

type SimplifiedAssetsList struct {
	Version          string                 `json:"version"`
	AssetGenesis     AssetsListAssetGenesis `json:"asset_genesis"`
	Amount           int                    `json:"amount"`
	LockTime         int                    `json:"lock_time"`
	ScriptKeyIsLocal bool                   `json:"script_key_is_local"`
	ChainAnchor      AssetsListChainAnchor  `json:"chain_anchor"`
	IsSpent          bool                   `json:"is_spent"`
	LeaseOwner       string                 `json:"lease_owner"`
	LeaseExpiry      int                    `json:"lease_expiry"`
	IsBurn           bool                   `json:"is_burn"`
}

type AssetsListAssetGenesis struct {
	GenesisPoint string `json:"genesis_point"`
	Name         string `json:"name"`
	MetaHash     string `json:"meta_hash"`
	AssetID      string `json:"asset_id"`
	AssetType    string `json:"asset_type"`
	OutputIndex  int    `json:"output_index"`
	Version      int    `json:"version"`
}

type AssetsListChainAnchor struct {
	AnchorTx         string `json:"anchor_tx"`
	AnchorBlockHash  string `json:"anchor_block_hash"`
	AnchorOutpoint   string `json:"anchor_outpoint"`
	InternalKey      string `json:"internal_key"`
	MerkleRoot       string `json:"merkle_root"`
	TapscriptSibling string `json:"tapscript_sibling"`
	BlockHeight      int    `json:"block_height"`
}

func SimplifyAssetsList(assets []*taprpc.Asset) *[]SimplifiedAssetsList {
	var simpleAssetsList []SimplifiedAssetsList
	for _, _asset := range assets {
		simpleAssetsList = append(simpleAssetsList, SimplifiedAssetsList{
			Version: _asset.Version.String(),
			AssetGenesis: AssetsListAssetGenesis{
				GenesisPoint: _asset.AssetGenesis.GenesisPoint,
				Name:         _asset.AssetGenesis.Name,
				MetaHash:     hex.EncodeToString(_asset.AssetGenesis.MetaHash),
				AssetID:      hex.EncodeToString(_asset.AssetGenesis.AssetId),
				AssetType:    _asset.AssetGenesis.AssetType.String(),
				OutputIndex:  int(_asset.AssetGenesis.OutputIndex),
				Version:      int(_asset.Version),
			},
			Amount:           int(_asset.Amount),
			LockTime:         int(_asset.LockTime),
			ScriptKeyIsLocal: _asset.ScriptKeyIsLocal,
			ChainAnchor: AssetsListChainAnchor{
				AnchorTx:         hex.EncodeToString(_asset.ChainAnchor.AnchorTx),
				AnchorBlockHash:  _asset.ChainAnchor.AnchorBlockHash,
				AnchorOutpoint:   _asset.ChainAnchor.AnchorOutpoint,
				InternalKey:      hex.EncodeToString(_asset.ChainAnchor.InternalKey),
				MerkleRoot:       hex.EncodeToString(_asset.ChainAnchor.MerkleRoot),
				TapscriptSibling: hex.EncodeToString(_asset.ChainAnchor.TapscriptSibling),
				BlockHeight:      int(_asset.ChainAnchor.BlockHeight),
			},
			IsSpent:     _asset.IsSpent,
			LeaseOwner:  hex.EncodeToString(_asset.LeaseOwner),
			LeaseExpiry: int(_asset.LeaseExpiry),
			IsBurn:      _asset.IsBurn,
		})
	}
	return &simpleAssetsList
}

type AssetsBalanceAssetGenesis struct {
	GenesisPoint string `json:"genesis_point"`
	Name         string `json:"name"`
	MetaHash     string `json:"meta_hash"`
	AssetID      string `json:"asset_id"`
	AssetType    string `json:"asset_type"`
	OutputIndex  int    `json:"output_index"`
	Version      int    `json:"version"`
}

type AssetsBalanceGroupBalance struct {
	GroupKey string `json:"group_key"`
	Balance  int    `json:"balance"`
}

func SyncUniverseFullSpecified(universeHost string, id string, proofType string) string {
	if universeHost == "" {
		switch base.NetWork {
		case base.UseTestNet:
			universeHost = "testnet.universe.lightning.finance:10029"
		case base.UseMainNet:
			universeHost = UniverseSocketMainnet
		case base.UseRegTest:
			universeHost = UniverseSocketRegtest
		}
	}
	var _proofType universerpc.ProofType
	if proofType == "issuance" || proofType == "ISSUANCE" || proofType == "PROOF_TYPE_ISSUANCE" {
		_proofType = universerpc.ProofType_PROOF_TYPE_ISSUANCE
	} else if proofType == "transfer" || proofType == "TRANSFER" || proofType == "PROOF_TYPE_TRANSFER" {
		_proofType = universerpc.ProofType_PROOF_TYPE_TRANSFER
	} else {
		_proofType = universerpc.ProofType_PROOF_TYPE_UNSPECIFIED
	}
	var targets []*universerpc.SyncTarget
	universeID := &universerpc.ID{
		Id: &universerpc.ID_AssetIdStr{
			AssetIdStr: id,
		},
		ProofType: _proofType,
	}
	targets = append(targets, &universerpc.SyncTarget{
		Id: universeID,
	})
	response, err := syncUniverse(universeHost, targets, universerpc.UniverseSyncMode_SYNC_FULL)
	if err != nil {
		return MakeJsonErrorResult(syncUniverseErr, err.Error(), "")
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func SyncAssetIssuance(id string) string {
	return SyncUniverseFullSpecified("", id, universerpc.ProofType_PROOF_TYPE_ISSUANCE.String())
}

func SyncAssetTransfer(id string) string {
	return SyncUniverseFullSpecified("", id, universerpc.ProofType_PROOF_TYPE_TRANSFER.String())
}

func SyncAssetAll(id string) {
	fmt.Println(SyncAssetIssuance(id))
	fmt.Println(SyncAssetTransfer(id))
}

func SyncAssetAllSlice(ids []string) {
	if len(ids) == 0 {
		return
	}
	for _, _id := range ids {
		fmt.Println("Sync issuance:", _id, ".", SyncAssetIssuance(_id))
		fmt.Println("Sync transfer:", _id, ".", SyncAssetTransfer(_id))
	}
}

func SyncAssetAllWithAssets(ids ...string) {
	if len(ids) == 0 {
		return
	}
	for _, _id := range ids {
		fmt.Println(SyncAssetIssuance(_id))
		fmt.Println(SyncAssetTransfer(_id))
	}
}

type AssetBalance struct {
	Name      string `json:"name"`
	MetaHash  string `json:"meta_hash"`
	AssetID   string `json:"asset_id"`
	AssetType string `json:"asset_type"`
	Balance   int    `json:"balance"`
}

type AssetGroupBalance struct {
	ID       string `json:"id"`
	Balance  int    `json:"balance"`
	GroupKey string `json:"group_key"`
}

func allAssetBalances() *[]AssetBalance {
	response, _ := listBalances(false, nil, nil)
	var assetBalances []AssetBalance
	for _, v := range response.AssetBalances {
		assetBalances = append(assetBalances, AssetBalance{
			Name:      v.AssetGenesis.Name,
			MetaHash:  hex.EncodeToString(v.AssetGenesis.MetaHash),
			AssetID:   hex.EncodeToString(v.AssetGenesis.AssetId),
			AssetType: v.AssetGenesis.AssetType.String(),
			Balance:   int(v.Balance),
		})
	}
	if len(assetBalances) == 0 {
		return nil
	}
	return &assetBalances
}

func GetAllAssetBalances() string {
	result := allAssetBalances()
	if result == nil {
		return MakeJsonErrorResult(allAssetBalancesErr, "Null Balances", nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", result)
}

func allAssetGroupBalances() *[]AssetGroupBalance {
	response, _ := listBalances(false, nil, nil)
	var assetGroupBalances []AssetGroupBalance
	for k, v := range response.AssetGroupBalances {
		assetGroupBalances = append(assetGroupBalances, AssetGroupBalance{
			ID:       k,
			Balance:  int(v.Balance),
			GroupKey: hex.EncodeToString(v.GroupKey),
		})
	}
	if len(assetGroupBalances) == 0 {
		return nil
	}
	return &assetGroupBalances
}

func GetAllAssetGroupBalances() string {
	result := allAssetGroupBalances()
	if result == nil {
		return MakeJsonErrorResult(allAssetGroupBalancesErr, "Null Asset Group Balances", nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", result)
}

func GetAllAssetIdByAssetBalance(assetBalance *[]AssetBalance) *[]string {
	if assetBalance == nil {
		return nil
	}
	var ids []string
	for _, v := range *assetBalance {
		ids = append(ids, v.AssetID)
	}
	return &ids
}

func SyncAllAssetsByAssetBalance() string {
	ids := GetAllAssetIdByAssetBalance(allAssetBalances())
	if ids != nil {
		SyncAssetAllSlice(*ids)
	}
	return MakeJsonErrorResult(SUCCESS, "", ids)
}

func GetAllAssetsIdSlice() string {
	ids := GetAllAssetIdByAssetBalance(allAssetBalances())
	return MakeJsonErrorResult(SUCCESS, "", ids)
}

func assetKeysTransfer(id string) *[]AssetKey {
	var _proofType universerpc.ProofType
	_proofType = universerpc.ProofType_PROOF_TYPE_TRANSFER
	response, err := assetLeafKeys(id, _proofType)
	if err != nil {
		fmt.Printf("%s universerpc AssetLeafKeys Error: %v\n", GetTimeNow(), err)
		return nil
	}
	if len(response.AssetKeys) == 0 {
		return nil
	}
	return processAssetKey(response)
}

func AssetKeysTransfer(id string) string {
	result := assetKeysTransfer(id)
	if result == nil {
		return MakeJsonErrorResult(assetKeysTransferErr, "Null Asset Keys", nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", result)
}

func AssetLeavesSpecified(id string, proofType string) *universerpc.AssetLeafResponse {
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
		fmt.Printf("%s universerpc AssetLeaves Error: %v\n", GetTimeNow(), err)
		return nil
	}
	return response
}

type AssetTransferLeave struct {
	Name      string `json:"name"`
	AssetID   string `json:"asset_id"`
	Amount    int    `json:"amount"`
	ScriptKey string `json:"script_key"`
	Proof     string `json:"proof"`
}

func ProcessAssetTransferLeave(response *universerpc.AssetLeafResponse) *[]AssetTransferLeave {
	var assetTransferLeaves []AssetTransferLeave
	for _, leave := range response.Leaves {
		assetTransferLeaves = append(assetTransferLeaves, AssetTransferLeave{
			Name:      leave.Asset.AssetGenesis.Name,
			AssetID:   hex.EncodeToString(leave.Asset.AssetGenesis.AssetId),
			Amount:    int(leave.Asset.Amount),
			ScriptKey: hex.EncodeToString(leave.Asset.ScriptKey),
			Proof:     hex.EncodeToString(leave.Proof),
		})
	}
	return &assetTransferLeaves
}

func AssetLeavesTransfer(id string) string {
	response := AssetLeavesSpecified(id, universerpc.ProofType_PROOF_TYPE_TRANSFER.String())
	if response == nil {
		return MakeJsonErrorResult(AssetLeavesSpecifiedErr, errors.New("null asset leaves").Error(), nil)
	}
	assetTransferLeaves := ProcessAssetTransferLeave(response)
	return MakeJsonErrorResult(SUCCESS, "", assetTransferLeaves)
}

func AssetLeavesTransfer_ONLY_FOR_TEST(id string) *[]AssetTransferLeave {
	response := AssetLeavesSpecified(id, universerpc.ProofType_PROOF_TYPE_TRANSFER.String())
	if response == nil {
		fmt.Printf("%s universerpc AssetLeaves Error.\n", GetTimeNow())
		return nil
	}
	return ProcessAssetTransferLeave(response)
}

func assetLeavesTransfer(id string) *[]AssetTransferLeave {
	response := AssetLeavesSpecified(id, universerpc.ProofType_PROOF_TYPE_TRANSFER.String())
	if response == nil {
		fmt.Printf("%s universerpc AssetLeaves Error.\n", GetTimeNow())
		return nil
	}
	return ProcessAssetTransferLeave(response)
}

type AssetIssuanceLeave struct {
	Version            string `json:"version"`
	GenesisPoint       string `json:"genesis_point"`
	Name               string `json:"name"`
	MetaHash           string `json:"meta_hash"`
	AssetID            string `json:"asset_id"`
	AssetType          string `json:"asset_type"`
	GenesisOutputIndex int    `json:"genesis_output_index"`
	Amount             int    `json:"amount"`
	LockTime           int    `json:"lock_time"`
	RelativeLockTime   int    `json:"relative_lock_time"`
	ScriptVersion      int    `json:"script_version"`
	ScriptKey          string `json:"script_key"`
	ScriptKeyIsLocal   bool   `json:"script_key_is_local"`
	IsSpent            bool   `json:"is_spent"`
	LeaseOwner         string `json:"lease_owner"`
	LeaseExpiry        int    `json:"lease_expiry"`
	IsBurn             bool   `json:"is_burn"`
	Proof              string `json:"proof"`
}

func ProcessAssetIssuanceLeave(response *universerpc.AssetLeafResponse) *AssetIssuanceLeave {
	if response == nil {
		return nil
	}
	return &AssetIssuanceLeave{
		Version:            response.Leaves[0].Asset.Version.String(),
		GenesisPoint:       response.Leaves[0].Asset.AssetGenesis.GenesisPoint,
		Name:               response.Leaves[0].Asset.AssetGenesis.Name,
		MetaHash:           hex.EncodeToString(response.Leaves[0].Asset.AssetGenesis.MetaHash),
		AssetID:            hex.EncodeToString(response.Leaves[0].Asset.AssetGenesis.AssetId),
		AssetType:          response.Leaves[0].Asset.AssetGenesis.AssetType.String(),
		GenesisOutputIndex: int(response.Leaves[0].Asset.AssetGenesis.OutputIndex),
		Amount:             int(response.Leaves[0].Asset.Amount),
		LockTime:           int(response.Leaves[0].Asset.LockTime),
		RelativeLockTime:   int(response.Leaves[0].Asset.RelativeLockTime),
		ScriptVersion:      int(response.Leaves[0].Asset.ScriptVersion),
		ScriptKey:          hex.EncodeToString(response.Leaves[0].Asset.ScriptKey),
		ScriptKeyIsLocal:   response.Leaves[0].Asset.ScriptKeyIsLocal,
		IsSpent:            response.Leaves[0].Asset.IsSpent,
		LeaseOwner:         hex.EncodeToString(response.Leaves[0].Asset.LeaseOwner),
		LeaseExpiry:        int(response.Leaves[0].Asset.LeaseExpiry),
		IsBurn:             response.Leaves[0].Asset.IsBurn,
		Proof:              hex.EncodeToString(response.Leaves[0].Proof),
	}
}

func assetLeavesIssuance(id string) *AssetIssuanceLeave {
	response := AssetLeavesSpecified(id, universerpc.ProofType_PROOF_TYPE_ISSUANCE.String())
	if response == nil {
		fmt.Printf("%s Universerpc asset leaves issuance error.\n", GetTimeNow())
		return nil
	}
	return ProcessAssetIssuanceLeave(response)
}

func GetAssetInfoByIssuanceLeaf(id string) string {
	response := assetLeavesIssuance(id)
	if response == nil {
		return MakeJsonErrorResult(assetLeavesIssuanceErr, errors.New("Null asset leaves").Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func DecodeRawProofByte(rawProof []byte) *taprpc.DecodeProofResponse {
	result, err := rpcclient.DecodeProof(rawProof, 0, false, false)
	if err != nil {
		logrus.Errorln(errors.Wrap(err, "rpcclient.DecodeProof"))
		return nil
	}
	return result
}

func DecodeRawProofString(proof string) *taprpc.DecodeProofResponse {
	decodeString, err := hex.DecodeString(proof)
	if err != nil {
		logrus.Errorln(errors.Wrap(err, "hex.DecodeString"))
		return nil
	}
	return DecodeRawProofByte(decodeString)
}

type DecodedProof struct {
	NumberOfProofs  int    `json:"number_of_proofs"`
	Name            string `json:"name"`
	AssetID         string `json:"asset_id"`
	Amount          int    `json:"amount"`
	ScriptKey       string `json:"script_key"`
	AnchorTx        string `json:"anchor_tx"`
	AnchorBlockHash string `json:"anchor_block_hash"`
	AnchorOutpoint  string `json:"anchor_outpoint"`
	InternalKey     string `json:"internal_key"`
	MerkleRoot      string `json:"merkle_root"`
	BlockHeight     int    `json:"block_height"`
}

func ProcessProof(response *taprpc.DecodeProofResponse) *DecodedProof {
	if response == nil {
		return &DecodedProof{}
	}
	return &DecodedProof{
		NumberOfProofs:  int(response.DecodedProof.NumberOfProofs),
		Name:            response.DecodedProof.Asset.AssetGenesis.Name,
		AssetID:         hex.EncodeToString(response.DecodedProof.Asset.AssetGenesis.AssetId),
		Amount:          int(response.DecodedProof.Asset.Amount),
		ScriptKey:       hex.EncodeToString(response.DecodedProof.Asset.ScriptKey),
		AnchorTx:        hex.EncodeToString(response.DecodedProof.Asset.ChainAnchor.AnchorTx),
		AnchorBlockHash: response.DecodedProof.Asset.ChainAnchor.AnchorBlockHash,
		AnchorOutpoint:  response.DecodedProof.Asset.ChainAnchor.AnchorOutpoint,
		InternalKey:     hex.EncodeToString(response.DecodedProof.Asset.ChainAnchor.InternalKey),
		MerkleRoot:      hex.EncodeToString(response.DecodedProof.Asset.ChainAnchor.MerkleRoot),
		BlockHeight:     int(response.DecodedProof.Asset.ChainAnchor.BlockHeight),
	}
}

func DecodeRawProof(proof string) string {
	response := DecodeRawProofString(proof)
	if response == nil {
		return MakeJsonErrorResult(DecodeRawProofStringErr, "null raw proof", nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", ProcessProof(response))
}

func allAssetList() *taprpc.ListAssetResponse {
	response, err := listAssets(false, true, false)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return response
}

type ListAllAsset struct {
	Version            string `json:"version"`
	GenesisPoint       string `json:"genesis_point"`
	GenesisName        string `json:"genesis_name"`
	GenesisMetaHash    string `json:"genesis_meta_hash"`
	GenesisAssetID     string `json:"genesis_asset_id"`
	GenesisAssetType   string `json:"genesis_asset_type"`
	GenesisOutputIndex int    `json:"genesis_output_index"`
	Amount             int    `json:"amount"`
	LockTime           int    `json:"lock_time"`
	RelativeLockTime   int    `json:"relative_lock_time"`
	ScriptVersion      int    `json:"script_version"`
	ScriptKey          string `json:"script_key"`
	ScriptKeyIsLocal   bool   `json:"script_key_is_local"`
	AnchorTx           string `json:"anchor_tx"`
	AnchorBlockHash    string `json:"anchor_block_hash"`
	AnchorOutpoint     string `json:"anchor_outpoint"`
	AnchorInternalKey  string `json:"anchor_internal_key"`
	AnchorBlockHeight  int    `json:"anchor_block_height"`
	IsSpent            bool   `json:"is_spent"`
	LeaseOwner         string `json:"lease_owner"`
	LeaseExpiry        int    `json:"lease_expiry"`
	IsBurn             bool   `json:"is_burn"`
}

func ProcessListAllAssets(response *taprpc.ListAssetResponse) *[]ListAllAsset {
	if response == nil || response.Assets == nil || len(response.Assets) == 0 {
		return nil
	}
	var listAllAssets []ListAllAsset
	for _, asset := range response.Assets {
		listAllAssets = append(listAllAssets, ListAllAsset{
			Version:            asset.Version.String(),
			GenesisPoint:       asset.AssetGenesis.GenesisPoint,
			GenesisName:        asset.AssetGenesis.Name,
			GenesisMetaHash:    hex.EncodeToString(asset.AssetGenesis.MetaHash),
			GenesisAssetID:     hex.EncodeToString(asset.AssetGenesis.AssetId),
			GenesisAssetType:   asset.AssetGenesis.AssetType.String(),
			GenesisOutputIndex: int(asset.AssetGenesis.OutputIndex),
			Amount:             int(asset.Amount),
			LockTime:           int(asset.LockTime),
			RelativeLockTime:   int(asset.RelativeLockTime),
			ScriptVersion:      int(asset.ScriptVersion),
			ScriptKey:          hex.EncodeToString(asset.ScriptKey),
			ScriptKeyIsLocal:   asset.ScriptKeyIsLocal,
			AnchorTx:           hex.EncodeToString(asset.ChainAnchor.AnchorTx),
			AnchorBlockHash:    asset.ChainAnchor.AnchorBlockHash,
			AnchorOutpoint:     asset.ChainAnchor.AnchorOutpoint,
			AnchorInternalKey:  hex.EncodeToString(asset.ChainAnchor.InternalKey),
			AnchorBlockHeight:  int(asset.ChainAnchor.BlockHeight),
			IsSpent:            asset.IsSpent,
			LeaseOwner:         hex.EncodeToString(asset.LeaseOwner),
			LeaseExpiry:        int(asset.LeaseExpiry),
			IsBurn:             asset.IsBurn,
		})
	}
	if len(listAllAssets) == 0 {
		return nil
	}
	return &listAllAssets
}

func GetAllAssetList() string {
	response := allAssetList()
	if response == nil {
		return MakeJsonErrorResult(allAssetListErr, "null asset list", nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", ProcessListAllAssets(response))
}

type ListAllAssetSimplified struct {
	GenesisName      string `json:"genesis_name"`
	GenesisAssetID   string `json:"genesis_asset_id"`
	GenesisAssetType string `json:"genesis_asset_type"`
	Amount           int    `json:"amount"`
	AnchorOutpoint   string `json:"anchor_outpoint"`
	IsSpent          bool   `json:"is_spent"`
}

func ProcessListAllAssetsSimplified(result *[]ListAllAsset) *[]ListAllAssetSimplified {
	if result == nil || len(*result) == 0 {
		return nil
	}
	var listAllAssetsSimplified []ListAllAssetSimplified
	for _, asset := range *result {
		listAllAssetsSimplified = append(listAllAssetsSimplified, ListAllAssetSimplified{
			GenesisName:      asset.GenesisName,
			GenesisAssetID:   asset.GenesisAssetID,
			GenesisAssetType: asset.GenesisAssetType,
			Amount:           asset.Amount,
			AnchorOutpoint:   asset.AnchorOutpoint,
			IsSpent:          asset.IsSpent,
		})
	}
	if len(listAllAssetsSimplified) == 0 {
		return nil
	}
	return &listAllAssetsSimplified
}

func GetAllAssetListSimplified() string {
	result := ProcessListAllAssetsSimplified(ProcessListAllAssets(allAssetList()))
	if result == nil {
		return MakeJsonErrorResult(ProcessListAllAssetsSimplifiedErr, "null asset list", nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", result)
}

func GetAllAssetIdByListAll() []string {
	id := make(map[string]bool)
	var ids []string
	result := ProcessListAllAssetsSimplified(ProcessListAllAssets(allAssetList()))
	if result == nil || len(*result) == 0 {
		return nil
	}
	for _, asset := range *result {
		id[asset.GenesisAssetID] = true
	}
	for k, _ := range id {
		ids = append(ids, k)
	}
	if len(ids) == 0 {
		return nil
	}
	return ids
}

func SyncUniverseFullIssuanceByIdSlice(ids []string) string {
	var universeHost string
	switch base.NetWork {
	case base.UseTestNet:
		universeHost = "testnet.universe.lightning.finance:10029"
	case base.UseMainNet:
		universeHost = UniverseSocketMainnet
	case base.UseRegTest:
		universeHost = UniverseSocketRegtest
	}
	var targets []*universerpc.SyncTarget
	for _, id := range ids {
		targets = append(targets, &universerpc.SyncTarget{
			Id: &universerpc.ID{
				Id: &universerpc.ID_AssetIdStr{
					AssetIdStr: id,
				},
				ProofType: universerpc.ProofType_PROOF_TYPE_ISSUANCE,
			},
		})
	}
	response, err := syncUniverse(universeHost, targets, universerpc.UniverseSyncMode_SYNC_FULL)
	if err != nil {
		return MakeJsonErrorResult(syncUniverseErr, err.Error(), "")
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func SyncUniverseFullTransferByIdSlice(ids []string) string {
	var universeHost string
	switch base.NetWork {
	case base.UseTestNet:
		universeHost = "testnet.universe.lightning.finance:10029"
	case base.UseMainNet:
		universeHost = UniverseSocketMainnet
	case base.UseRegTest:
		universeHost = UniverseSocketRegtest
	}
	var targets []*universerpc.SyncTarget
	for _, id := range ids {
		targets = append(targets, &universerpc.SyncTarget{
			Id: &universerpc.ID{
				Id: &universerpc.ID_AssetIdStr{
					AssetIdStr: id,
				},
				ProofType: universerpc.ProofType_PROOF_TYPE_TRANSFER,
			},
		})
	}
	response, err := syncUniverse(universeHost, targets, universerpc.UniverseSyncMode_SYNC_FULL)
	if err != nil {
		return MakeJsonErrorResult(syncUniverseErr, err.Error(), "")
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func SyncUniverseFullNoSlice() string {
	var universeHost string
	switch base.NetWork {
	case base.UseTestNet:
		universeHost = "testnet.universe.lightning.finance:10029"
	case base.UseMainNet:
		universeHost = "universe.lightning.finance:10029"
	}
	var targets []*universerpc.SyncTarget
	response, err := syncUniverse(universeHost, targets, universerpc.UniverseSyncMode_SYNC_FULL)
	if err != nil {
		return MakeJsonErrorResult(syncUniverseErr, err.Error(), "")
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

type AssetHoldInfo struct {
	Name      string `json:"name"`
	AssetId   string `json:"assetId"`
	Amount    int    `json:"amount"`
	Outpoint  string `json:"outpoint"`
	Address   string `json:"address"`
	ScriptKey string `json:"scriptKey"`
	IsSpent   bool   `json:"isSpent"`
}

func OutpointToAddress(outpoint string) string {
	transaction, indexStr := getTransactionAndIndexByOutpoint(outpoint)
	index, _ := strconv.Atoi(indexStr)
	response, err := getTransactionByMempool(transaction)
	if err != nil {
		return ""
	}
	return response.Vout[index].ScriptpubkeyAddress
}

func TransactionAndIndexToAddress(transaction string, indexStr string) string {
	index, _ := strconv.Atoi(indexStr)
	response, err := getTransactionByMempool(transaction)
	if err != nil {
		return ""
	}
	return response.Vout[index].ScriptpubkeyAddress
}

func TransactionAndIndexToValue(transaction string, indexStr string) int {
	index, _ := strconv.Atoi(indexStr)
	response, err := getTransactionByMempool(transaction)
	if err != nil {
		return 0
	}
	return response.Vout[index].Value
}

func getTransactionAndIndexByOutpoint(outpoint string) (transaction string, index string) {
	result := strings.Split(outpoint, ":")
	return result[0], result[1]
}

func CompareScriptKey(scriptKey1 string, scriptKey2 string) string {
	if scriptKey1 == scriptKey2 {
		return scriptKey1
	} else if len(scriptKey1) == len(scriptKey2) {
		return ""
	} else if len(scriptKey1) > len(scriptKey2) {
		if scriptKey1 == "0"+scriptKey2 || scriptKey1 == "02"+scriptKey2 {
			return scriptKey1
		} else if scriptKey1 == "2"+scriptKey2 {
			return "02" + scriptKey2
		} else {
			return ""
		}
	} else if len(scriptKey1) < len(scriptKey2) {
		if "0"+scriptKey1 == scriptKey2 || "02"+scriptKey1 == scriptKey2 {
			return scriptKey2
		} else if "2"+scriptKey1 == scriptKey2 {
			return "02" + scriptKey1
		} else {
			return ""
		}
	}
	return ""
}

func GetAssetHoldInfosIncludeSpent(id string) *[]AssetHoldInfo {
	assetLeavesTransfers := assetLeavesTransfer(id)
	var idToAssetHoldInfo []AssetHoldInfo
	for _, leaf := range *assetLeavesTransfers {
		outpoint := ProcessProof(DecodeRawProofString(leaf.Proof)).AnchorOutpoint
		address := OutpointToAddress(outpoint)
		idToAssetHoldInfo = append(idToAssetHoldInfo, AssetHoldInfo{
			Name:      leaf.Name,
			AssetId:   leaf.AssetID,
			Amount:    leaf.Amount,
			Outpoint:  outpoint,
			Address:   address,
			ScriptKey: leaf.ScriptKey,
			IsSpent:   AddressIsSpentAll(address),
		})
	}
	return &idToAssetHoldInfo
}

func GetAssetHoldInfosExcludeSpent(id string) *[]AssetHoldInfo {
	assetLeavesTransfers := assetLeavesTransfer(id)
	var idToAssetHoldInfo []AssetHoldInfo
	for _, leaf := range *assetLeavesTransfers {
		outpoint := ProcessProof(DecodeRawProofString(leaf.Proof)).AnchorOutpoint
		address := OutpointToAddress(outpoint)
		isSpent := AddressIsSpentAll(address)
		if !isSpent {
			idToAssetHoldInfo = append(idToAssetHoldInfo, AssetHoldInfo{
				Name:      leaf.Name,
				AssetId:   leaf.AssetID,
				Amount:    leaf.Amount,
				Outpoint:  outpoint,
				Address:   address,
				ScriptKey: leaf.ScriptKey,
				IsSpent:   isSpent,
			})
		}
	}
	return &idToAssetHoldInfo
}

func GetAssetHoldInfosIncludeSpentSlow(id string) string {
	response := GetAssetHoldInfosIncludeSpent(id)
	if response == nil {
		return MakeJsonErrorResult(GetAssetHoldInfosIncludeSpentErr, "Get asset hold infos include spent fail, null response.", nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func AddressIsSpent(address string) bool {
	addressInfo := getAddressInfoByMempool(address)
	if addressInfo.ChainStats.SpentTxoSum == 0 {
		return false
	}
	return true

}

func AddressIsSpentAll(address string) bool {
	if !AddressIsSpent(address) {
		return false
	}
	addressInfo := getAddressInfoByMempool(address)
	if int(addressInfo.ChainStats.FundedTxoSum) == addressInfo.ChainStats.SpentTxoSum {
		return true
	}
	return false
}

func OutpointToTransactionInfo(outpoint string) *AssetTransactionInfo {
	transaction, indexStr := getTransactionAndIndexByOutpoint(outpoint)
	index, _ := strconv.Atoi(indexStr)
	response, err := getTransactionByMempool(transaction)
	if err != nil {
		logrus.Errorln(errors.Wrap(err, "getTransactionByMempool"))
		return &AssetTransactionInfo{}
	}
	var inputs []string
	for _, input := range response.Vin {
		if input.Prevout.Value == 1000 {
			inputs = append(inputs, input.Prevout.ScriptpubkeyAddress)
		}
	}
	result := AssetTransactionInfo{
		AnchorTransaction: response.Txid,
		From:              inputs,
		To:                response.Vout[index].ScriptpubkeyAddress,
		BlockTime:         response.Status.BlockTime,
		FeeRate:           RoundToDecimalPlace(float64(response.Fee)/(float64(response.Weight)/4), 2),
		ConfirmedBlocks:   BlockTipHeight() - response.Status.BlockHeight,
	}
	return &result
}

type AssetTransactionInfo struct {
	AnchorTransaction string   `json:"anchor_transaction"`
	From              []string `json:"from"`
	To                string   `json:"to"`
	Name              string   `json:"name"`
	AssetId           string   `json:"assetId"`
	Amount            int      `json:"amount"`
	BlockTime         int      `json:"block_time"`
	FeeRate           float64  `json:"fee_rate"`
	ConfirmedBlocks   int      `json:"confirmed_blocks"`
	IsSpent           bool     `json:"isSpent"`
}

func GetAssetTransactionInfos(id string) *[]AssetTransactionInfo {
	assetLeavesTransfers := assetLeavesTransfer(id)
	var idToAssetTransactionInfos []AssetTransactionInfo
	for _, leaf := range *assetLeavesTransfers {
		outpoint := ProcessProof(DecodeRawProofString(leaf.Proof)).AnchorOutpoint
		transactionInfo := OutpointToTransactionInfo(outpoint)
		transactionInfo.Name = leaf.Name
		transactionInfo.AssetId = leaf.AssetID
		transactionInfo.Amount = leaf.Amount
		transactionInfo.IsSpent = AddressIsSpentAll(transactionInfo.To)
		idToAssetTransactionInfos = append(idToAssetTransactionInfos, *transactionInfo)
	}
	return &idToAssetTransactionInfos
}

func SyncAllAssetByList() string {
	SyncAssetAllSlice(GetAllAssetIdByListAll())
	return MakeJsonErrorResult(SUCCESS, "", "Sync Completed.")
}

func GetAssetInfoById(id string) string {
	return GetAssetInfoByIssuanceLeaf(id)
}

func GetAssetHoldInfosExcludeSpentSlow(id string) string {
	response := GetAssetHoldInfosExcludeSpent(id)
	if response == nil {
		return MakeJsonErrorResult(GetAssetHoldInfosExcludeSpentErr, "Get asset hold infos exclude spent fail, null response.", nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func GetAssetTransactionInfoSlow(id string) string {
	response := GetAssetTransactionInfos(id)
	if response == nil {
		return MakeJsonErrorResult(GetAssetTransactionInfosErr, "Get asset transaction infos fail, null response.", nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func AssetIDAndTransferScriptKeyToOutpoint(id string, scriptKey string) string {
	keys := assetKeysTransfer(id)
	for _, key := range *keys {
		cs := CompareScriptKey(scriptKey, key.ScriptKeyBytes)
		if scriptKey == cs {
			return key.OpStr
		}
	}
	return ""
}

func GetAllAssetListWithoutProcession() string {
	response := allAssetList()
	if response == nil {
		return MakeJsonErrorResult(allAssetListErr, "Null list asset response.", nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func ListBatchesAndGetResponse() (*mintrpc.ListBatchResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()
	client := mintrpc.NewMintClient(conn)
	request := &mintrpc.ListBatchRequest{}
	response, err := client.ListBatches(context.Background(), request)
	if err != nil {
		fmt.Printf("%s mintrpc ListBatches Error: %v\n", GetTimeNow(), err)
		return nil, err
	}
	return response, nil
}

type ListBatchesResponse struct {
	BatchKey        string             `json:"batch_key"`
	BatchTxid       string             `json:"batch_txid"`
	State           string             `json:"state"`
	Assets          []ListBatchesAsset `json:"asset_meta"`
	Amount          int                `json:"amount"`
	NewGroupedAsset bool               `json:"new_grouped_asset"`
	GroupKey        string             `json:"group_key"`
	GroupAnchor     string             `json:"group_anchor"`
}

type ListBatchesAsset struct {
	AssetVersion string               `json:"asset_version"`
	AssetType    string               `json:"asset_type"`
	Name         string               `json:"name"`
	AssetMeta    ListBatchesAssetMeta `json:"asset_meta"`
}

type ListBatchesAssetMeta struct {
	Data     string `json:"data"`
	Type     string `json:"type"`
	MetaHash string `json:"meta_hash"`
}

func GetTransactionsAndGetResponse() (*lnrpc.TransactionDetails, error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return nil, err
	}
	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.GetTransactionsRequest{}
	response, err := client.GetTransactions(context.Background(), request)
	return response, err
}

func GetTransactionsExcludeLabelTapdAssetMinting() (*[]*lnrpc.Transaction, error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.GetTransactionsRequest{}
	response, err := client.GetTransactions(context.Background(), request)
	if err != nil {
		return nil, err
	}
	transactions := ExcludeLabelIsTapdAssetMinting(response)
	return &transactions, err
}

func ExcludeLabelIsTapdAssetMinting(response *lnrpc.TransactionDetails) []*lnrpc.Transaction {
	var transactions []*lnrpc.Transaction
	for _, transaction := range response.Transactions {
		if transaction.Label != "tapd-asset-minting" {
			transactions = append(transactions, transaction)
		}
	}
	return transactions
}

type GetTransactionsResponse struct {
	TxHash            string                             `json:"tx_hash"`
	Amount            int                                `json:"amount"`
	NumConfirmations  int                                `json:"num_confirmations"`
	BlockHash         string                             `json:"block_hash"`
	BlockHeight       int                                `json:"block_height"`
	TimeStamp         int                                `json:"time_stamp"`
	TotalFees         int                                `json:"total_fees"`
	DestAddresses     []string                           `json:"dest_addresses"`
	OutputDetails     []GetTransactionsOutputDetails     `json:"output_details"`
	RawTxHex          string                             `json:"raw_tx_hex"`
	Label             string                             `json:"label"`
	PreviousOutpoints []GetTransactionsPreviousOutpoints `json:"previous_outpoints"`
}

type GetTransactionsOutputDetails struct {
	OutputType   string `json:"output_type"`
	Address      string `json:"address"`
	PkScript     string `json:"pk_script"`
	OutputIndex  int    `json:"output_index"`
	Amount       int    `json:"amount"`
	IsOurAddress bool   `json:"is_our_address"`
}

type GetTransactionsPreviousOutpoints struct {
	Outpoint    string `json:"outpoint"`
	IsOurOutput bool   `json:"is_our_output"`
}

type AssetGenesisStruct struct {
	GenesisPoint string `json:"genesis_point"`
	Name         string `json:"name"`
	MetaHash     string `json:"meta_hash"`
	AssetID      string `json:"asset_id"`
	AssetType    int    `json:"asset_type"`
	OutputIndex  int    `json:"output_index"`
	Version      int    `json:"version"`
}

type ChainAnchorStruct struct {
	AnchorTx         string `json:"anchor_tx"`
	AnchorBlockHash  string `json:"anchor_block_hash"`
	AnchorOutpoint   string `json:"anchor_outpoint"`
	InternalKey      string `json:"internal_key"`
	MerkleRoot       string `json:"merkle_root"`
	TapscriptSibling string `json:"tapscript_sibling"`
	BlockHeight      int    `json:"block_height"`
}

type ListAssetResponse struct {
	Version          string             `json:"version"`
	AssetGenesis     AssetGenesisStruct `json:"asset_genesis"`
	Amount           int                `json:"amount"`
	LockTime         int                `json:"lock_time"`
	RelativeLockTime int                `json:"relative_lock_time"`
	ScriptVersion    int                `json:"script_version"`
	ScriptKey        string             `json:"script_key"`
	ScriptKeyIsLocal bool               `json:"script_key_is_local"`
	ChainAnchor      ChainAnchorStruct  `json:"chain_anchor"`
	IsSpent          bool               `json:"is_spent"`
	LeaseOwner       string             `json:"lease_owner"`
	LeaseExpiry      int                `json:"lease_expiry"`
	IsBurn           bool               `json:"is_burn"`
}

func ListAssetAndGetResponse() (*taprpc.ListAssetResponse, error) {
	return listAssets(false, true, false)
}

func ListAssetAndGetResponseByFlags(withWitness, includeSpent, includeLeased bool) (*taprpc.ListAssetResponse, error) {
	return listAssets(withWitness, includeSpent, includeLeased)
}

func ListBatchesAndGetCustomResponse() (*[]ListBatchesResponse, error) {
	response, err := ListBatchesAndGetResponse()
	if err != nil {
		LogError("", err)
		return nil, err
	}
	var listBatchesResponse []ListBatchesResponse
	for _, batch := range (*response).Batches {
		var assets []ListBatchesAsset
		for _, _asset := range batch.Batch.Assets {
			assets = append(assets, ListBatchesAsset{
				AssetVersion: _asset.AssetVersion.String(),
				AssetType:    _asset.AssetType.String(),
				Name:         _asset.Name,
				AssetMeta: ListBatchesAssetMeta{
					Data:     hex.EncodeToString(_asset.AssetMeta.Data),
					Type:     _asset.AssetMeta.Type.String(),
					MetaHash: hex.EncodeToString(_asset.AssetMeta.MetaHash),
				},
			})
		}
		listBatchesResponse = append(listBatchesResponse, ListBatchesResponse{
			BatchKey:  hex.EncodeToString(batch.Batch.BatchKey),
			BatchTxid: batch.Batch.BatchTxid,
			State:     batch.Batch.State.String(),
			Assets:    assets,
		})
	}
	return &listBatchesResponse, nil
}

func ListAssetAndGetCustomResponse() (*[]ListAssetResponse, error) {
	response, err := listAssets(false, true, false)
	if err != nil {
		LogError("", err)
		return nil, err
	}
	var listAssetResponses []ListAssetResponse
	for _, _asset := range (*response).Assets {
		listAssetResponses = append(listAssetResponses, ListAssetResponse{
			Version: _asset.Version.String(),
			AssetGenesis: AssetGenesisStruct{
				GenesisPoint: _asset.AssetGenesis.GenesisPoint,
				Name:         _asset.AssetGenesis.Name,
				MetaHash:     hex.EncodeToString(_asset.AssetGenesis.MetaHash),
				AssetID:      hex.EncodeToString(_asset.AssetGenesis.AssetId),
				AssetType:    int(_asset.AssetGenesis.AssetType),
				OutputIndex:  int(_asset.AssetGenesis.OutputIndex),
				Version:      int(_asset.Version),
			},
			Amount:           int(_asset.Amount),
			LockTime:         int(_asset.LockTime),
			RelativeLockTime: int(_asset.RelativeLockTime),
			ScriptVersion:    int(_asset.ScriptVersion),
			ScriptKey:        hex.EncodeToString(_asset.ScriptKey),
			ScriptKeyIsLocal: _asset.ScriptKeyIsLocal,
			ChainAnchor: ChainAnchorStruct{
				AnchorTx:         hex.EncodeToString(_asset.ChainAnchor.AnchorTx),
				AnchorBlockHash:  _asset.ChainAnchor.AnchorBlockHash,
				AnchorOutpoint:   _asset.ChainAnchor.AnchorOutpoint,
				InternalKey:      hex.EncodeToString(_asset.ChainAnchor.InternalKey),
				MerkleRoot:       hex.EncodeToString(_asset.ChainAnchor.MerkleRoot),
				TapscriptSibling: hex.EncodeToString(_asset.ChainAnchor.TapscriptSibling),
				BlockHeight:      int(_asset.ChainAnchor.BlockHeight),
			},
			IsSpent:     _asset.IsSpent,
			LeaseOwner:  hex.EncodeToString(_asset.LeaseOwner),
			LeaseExpiry: int(_asset.LeaseExpiry),
			IsBurn:      _asset.IsBurn,
		})
	}
	return &listAssetResponses, nil
}

func GetTransactionsAndGetCustomResponse() (*[]GetTransactionsResponse, error) {
	response, err := GetTransactionsAndGetResponse()
	if err != nil {
		LogError("", err)
		return nil, err
	}
	var getTransactionsResponse []GetTransactionsResponse
	for _, transaction := range response.Transactions {
		var outputDetails []GetTransactionsOutputDetails
		for _, output := range transaction.OutputDetails {
			outputDetails = append(outputDetails, GetTransactionsOutputDetails{
				OutputType:   output.OutputType.String(),
				Address:      output.Address,
				PkScript:     output.PkScript,
				OutputIndex:  int(output.OutputIndex),
				Amount:       int(output.Amount),
				IsOurAddress: output.IsOurAddress,
			})
		}
		var previousOutpoints []GetTransactionsPreviousOutpoints
		for _, previousOutpoint := range transaction.PreviousOutpoints {
			previousOutpoints = append(previousOutpoints, GetTransactionsPreviousOutpoints{
				Outpoint:    previousOutpoint.Outpoint,
				IsOurOutput: previousOutpoint.IsOurOutput,
			})
		}
		getTransactionsResponse = append(getTransactionsResponse, GetTransactionsResponse{
			TxHash:            transaction.TxHash,
			Amount:            int(transaction.Amount),
			NumConfirmations:  int(transaction.NumConfirmations),
			BlockHash:         transaction.BlockHash,
			BlockHeight:       int(transaction.BlockHeight),
			TimeStamp:         int(transaction.TimeStamp),
			TotalFees:         int(transaction.TotalFees),
			DestAddresses:     transaction.DestAddresses,
			OutputDetails:     outputDetails,
			RawTxHex:          transaction.RawTxHex,
			Label:             transaction.Label,
			PreviousOutpoints: previousOutpoints,
		})
	}
	return &getTransactionsResponse, nil
}

func AssetLeafKeysIssuance(assetId string) (*universerpc.AssetLeafKeyResponse, error) {
	proofType := universerpc.ProofType_PROOF_TYPE_ISSUANCE
	return AssetLeafKeysAndGetResponse(assetId, proofType)
}

func AssetLeavesIssuance(assetId string) (*universerpc.AssetLeafResponse, error) {
	proofType := universerpc.ProofType_PROOF_TYPE_ISSUANCE
	return AssetLeavesAndGetResponse(false, assetId, proofType)
}

func GetTransactionsWhoseLabelIsTapdAssetMinting() (*[]GetTransactionsResponse, error) {
	response, err := GetTransactionsAndGetCustomResponse()
	if err != nil {
		LogError("", err)
		return nil, err
	}
	var getTransactionsResponse []GetTransactionsResponse
	for _, transaction := range *response {
		if transaction.Label == "tapd-asset-minting" {
			getTransactionsResponse = append(getTransactionsResponse, transaction)
		}
	}
	return &getTransactionsResponse, nil
}

func GetTransactionsWhoseLabelIsNotTapdAssetMinting() (*[]GetTransactionsResponse, error) {
	response, err := GetTransactionsAndGetCustomResponse()
	if err != nil {
		return nil, err
	}
	var getTransactionsResponse []GetTransactionsResponse
	for _, transaction := range *response {
		if transaction.Label != "tapd-asset-minting" {
			getTransactionsResponse = append(getTransactionsResponse, transaction)
		}
	}
	return &getTransactionsResponse, nil
}

type PostGetRawTransactionResponse struct {
	Result *PostGetRawTransactionResult `json:"result"`
	Error  *BitcoindRpcResponseError    `json:"error"`
	ID     string                       `json:"id"`
}

type BitcoindRpcResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type PostGetRawTransactionResult struct {
	Txid          string                     `json:"txid"`
	Hash          string                     `json:"hash"`
	Version       int                        `json:"version"`
	Size          int                        `json:"size"`
	Vsize         int                        `json:"vsize"`
	Weight        int                        `json:"weight"`
	Locktime      int                        `json:"locktime"`
	Vin           []RawTransactionResultVin  `json:"vin"`
	Vout          []RawTransactionResultVout `json:"vout"`
	Fee           float64                    `json:"fee"`
	Hex           string                     `json:"hex"`
	Blockhash     string                     `json:"blockhash"`
	Confirmations int                        `json:"confirmations"`
	Time          int                        `json:"time"`
	Blocktime     int                        `json:"blocktime"`
}

type RawTransactionResultVin struct {
	Txid        string                           `json:"txid"`
	Vout        int                              `json:"vout"`
	ScriptSig   RawTransactionResultVinScriptSig `json:"scriptSig"`
	Txinwitness []string                         `json:"txinwitness"`
	Prevout     RawTransactionResultVinPrevout   `json:"prevout"`
	Sequence    int                              `json:"sequence"`
}

type RawTransactionResultVinPrevout struct {
	Generated    bool                                       `json:"generated"`
	Height       int                                        `json:"height"`
	Value        float64                                    `json:"value"`
	ScriptPubKey RawTransactionResultVinPrevoutScriptPubKey `json:"scriptPubKey"`
}

type RawTransactionResultVinPrevoutScriptPubKey struct {
	Asm_    string `json:"asm"`
	Desc    string `json:"desc"`
	Hex     string `json:"hex"`
	Address string `json:"address"`
	Type    string `json:"type"`
}

type RawTransactionResultVinScriptSig struct {
	Asm_ string `json:"asm"`
	Hex  string `json:"hex"`
}

type RawTransactionResultVout struct {
	Value        float64                              `json:"value"`
	N            int                                  `json:"n"`
	ScriptPubKey RawTransactionResultVoutScriptPubKey `json:"scriptPubKey"`
}

type RawTransactionResultVoutScriptPubKey struct {
	Asm_    string `json:"asm"`
	Desc    string `json:"desc"`
	Hex     string `json:"hex"`
	Address string `json:"address"`
	Type    string `json:"type"`
}

type PostGetRawTransactionResultSat struct {
	Txid          string                        `json:"txid"`
	Hash          string                        `json:"hash"`
	Version       int                           `json:"version"`
	Size          int                           `json:"size"`
	Vsize         int                           `json:"vsize"`
	Weight        int                           `json:"weight"`
	Locktime      int                           `json:"locktime"`
	Vin           []RawTransactionResultVinSat  `json:"vin"`
	Vout          []RawTransactionResultVoutSat `json:"vout"`
	Fee           int                           `json:"fee"`
	Hex           string                        `json:"hex"`
	Blockhash     string                        `json:"blockhash"`
	Confirmations int                           `json:"confirmations"`
	Time          int                           `json:"time"`
	Blocktime     int                           `json:"blocktime"`
}

type RawTransactionResultVinSat struct {
	Txid        string                            `json:"txid"`
	Vout        int                               `json:"vout"`
	ScriptSig   RawTransactionResultVinScriptSig  `json:"scriptSig"`
	Txinwitness []string                          `json:"txinwitness"`
	Prevout     RawTransactionResultVinPrevoutSat `json:"prevout"`
	Sequence    int                               `json:"sequence"`
}

type RawTransactionResultVinPrevoutSat struct {
	Generated    bool                                       `json:"generated"`
	Height       int                                        `json:"height"`
	Value        int                                        `json:"value"`
	ScriptPubKey RawTransactionResultVinPrevoutScriptPubKey `json:"scriptPubKey"`
}

type RawTransactionResultVoutSat struct {
	Value        int                                  `json:"value"`
	N            int                                  `json:"n"`
	ScriptPubKey RawTransactionResultVoutScriptPubKey `json:"scriptPubKey"`
}

func DecodeTransactionsWhoseLabelIsNotTapdAssetMinting(token string, rawTransactions []string) (*DecodeRawTransactionsResponse, error) {
	decodedRawTransactions, err := PostCallBitcoindToDecodeRawTransaction(token, rawTransactions)
	if err != nil {
		return nil, err
	}
	return decodedRawTransactions, nil
}

func DecodeAndQueryTransactionsWhoseLabelIsNotTapdAssetMinting(token string, rawTransactions []string) (*DecodeAndQueryTransactionsResponse, error) {
	decodedRawTransactions, err := PostCallBitcoindToDecodeAndQueryTransaction(token, rawTransactions)
	if err != nil {
		return nil, err
	}
	return decodedRawTransactions, nil
}

func RawTransactionHexSliceToRequestBodyRawString(rawTransactions []string) (request string) {
	request = "{\"transactions\":["
	for i, transaction := range rawTransactions {
		element := fmt.Sprintf("\"%s\"", transaction)
		request += element
		if i != len(rawTransactions)-1 {
			request += ","
		}
	}
	request += "]}"
	return request
}

func OutpointSliceToRequestBodyRawString(outpoints []string) (request string) {
	request = "{\"outpoints\":["
	for i, transaction := range outpoints {
		element := fmt.Sprintf("\"%s\"", transaction)
		request += element
		if i != len(outpoints)-1 {
			request += ","
		}
	}
	request += "]}"
	return request
}

type DecodeRawTransactionsResponse struct {
	Success bool                                `json:"success"`
	Error   string                              `json:"error"`
	Code    int                                 `json:"code"`
	Data    *[]PostDecodeRawTransactionResponse `json:"data"`
}

func PostCallBitcoindToDecodeRawTransaction(token string, rawTransactions []string) (*DecodeRawTransactionsResponse, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	network := base.NetWork
	url := serverDomainOrSocket + "/bitcoind/" + network + "/decode/transactions"
	requestStr := RawTransactionHexSliceToRequestBodyRawString(rawTransactions)
	payload := strings.NewReader(requestStr)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response DecodeRawTransactionsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func PostCallBitcoindToDecodeAndQueryTransaction(token string, rawTransactions []string) (*DecodeAndQueryTransactionsResponse, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	network := base.NetWork
	url := serverDomainOrSocket + "/bitcoind/" + network + "/decode/query/transactions"
	requestStr := RawTransactionHexSliceToRequestBodyRawString(rawTransactions)
	payload := strings.NewReader(requestStr)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response DecodeAndQueryTransactionsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

type PostDecodeRawTransactionResponse struct {
	Result *PostDecodeRawTransactionResult `json:"result"`
	Error  *BitcoindRpcResponseError       `json:"error"`
	ID     string                          `json:"id"`
}

type PostDecodeRawTransactionResult struct {
	Txid     string                           `json:"txid"`
	Hash     string                           `json:"hash"`
	Version  int                              `json:"version"`
	Size     int                              `json:"size"`
	Vsize    int                              `json:"vsize"`
	Weight   int                              `json:"weight"`
	Locktime int                              `json:"locktime"`
	Vin      []DecodeRawTransactionResultVin  `json:"vin"`
	Vout     []DecodeRawTransactionResultVout `json:"vout"`
}

type DecodeRawTransactionResultVin struct {
	Txid        string                                 `json:"txid"`
	Vout        int                                    `json:"vout"`
	ScriptSig   DecodeRawTransactionResultVinScriptSig `json:"scriptSig"`
	Txinwitness []string                               `json:"txinwitness"`
	Sequence    int64                                  `json:"sequence"`
}

type DecodeRawTransactionResultVinScriptSig struct {
	Asm_ string `json:"asm"`
	Hex  string `json:"hex"`
}

type DecodeRawTransactionResultVout struct {
	Value        float64                                    `json:"value"`
	N            int                                        `json:"n"`
	ScriptPubKey DecodeRawTransactionResultVoutScriptPubKey `json:"scriptPubKey"`
}

type DecodeRawTransactionResultVoutScriptPubKey struct {
	Asm_    string `json:"asm"`
	Desc    string `json:"desc"`
	Hex     string `json:"hex"`
	Address string `json:"address"`
	Type    string `json:"type"`
}

func ProcessDecodedTransactionsData(decodedRawTransactions *[]PostDecodeRawTransactionResponse) *[]PostDecodeRawTransactionResponse {
	var result []PostDecodeRawTransactionResponse
	for _, rawTransaction := range *decodedRawTransactions {
		if rawTransaction.Error == nil {
			result = append(result, rawTransaction)
		}
	}
	return &result
}

func ProcessDecodedAndQueryTransactionsData(decodedRawTransactions *[]PostGetRawTransactionResponse) *[]PostGetRawTransactionResult {
	var result []PostGetRawTransactionResult
	if decodedRawTransactions == nil {
		return &result
	}
	for _, rawTransaction := range *decodedRawTransactions {
		if rawTransaction.Error == nil && rawTransaction.Result != nil {
			result = append(result, *(rawTransaction.Result))
		}
	}
	return &result
}

func GetThenDecodeTransactionsWhoseLabelIsNotTapdAssetMinting(token string) (*[]PostDecodeRawTransactionResponse, error) {
	getTransactions, err := GetTransactionsWhoseLabelIsNotTapdAssetMinting()
	if err != nil {
		return nil, err
	}
	var rawTransactions []string
	for _, transaction := range *getTransactions {
		rawTransactions = append(rawTransactions, transaction.RawTxHex)
	}
	decodedTransactions, err := DecodeTransactionsWhoseLabelIsNotTapdAssetMinting(token, rawTransactions)
	if err != nil {
		return nil, err
	}
	result := ProcessDecodedTransactionsData(decodedTransactions.Data)
	return result, nil
}

type DecodeAndQueryTransactionsResponse struct {
	Success bool                             `json:"success"`
	Error   string                           `json:"error"`
	Code    int                              `json:"code"`
	Data    *[]PostGetRawTransactionResponse `json:"data"`
}

func ProcessPostGetRawTransactionResultToUseSat(btcUesult *[]PostGetRawTransactionResult) *[]PostGetRawTransactionResultSat {
	var result []PostGetRawTransactionResultSat
	for _, transaction := range *btcUesult {
		var rawTransactionResultVinSats []RawTransactionResultVinSat
		for _, vin := range transaction.Vin {
			rawTransactionResultVinSats = append(rawTransactionResultVinSats, RawTransactionResultVinSat{
				Txid:        vin.Txid,
				Vout:        vin.Vout,
				ScriptSig:   vin.ScriptSig,
				Txinwitness: vin.Txinwitness,
				Prevout: RawTransactionResultVinPrevoutSat{
					Generated:    vin.Prevout.Generated,
					Height:       vin.Prevout.Height,
					Value:        ToSat(vin.Prevout.Value),
					ScriptPubKey: vin.Prevout.ScriptPubKey,
				},
				Sequence: vin.Sequence,
			})
		}
		var rawTransactionResultVoutSats []RawTransactionResultVoutSat
		for _, vout := range transaction.Vout {
			rawTransactionResultVoutSats = append(rawTransactionResultVoutSats, RawTransactionResultVoutSat{
				Value:        ToSat(vout.Value),
				N:            vout.N,
				ScriptPubKey: vout.ScriptPubKey,
			})
		}
		result = append(result, PostGetRawTransactionResultSat{
			Txid:          transaction.Txid,
			Hash:          transaction.Hash,
			Version:       transaction.Version,
			Size:          transaction.Size,
			Vsize:         transaction.Vsize,
			Weight:        transaction.Weight,
			Locktime:      transaction.Locktime,
			Vin:           rawTransactionResultVinSats,
			Vout:          rawTransactionResultVoutSats,
			Fee:           ToSat(transaction.Fee),
			Hex:           transaction.Hex,
			Blockhash:     transaction.Blockhash,
			Confirmations: transaction.Confirmations,
			Time:          transaction.Time,
			Blocktime:     transaction.Blocktime,
		})
	}
	return &result
}

type GetAddressesByOutpointSliceResponse struct {
	Success bool              `json:"success"`
	Error   string            `json:"error"`
	Code    ErrCode           `json:"code"`
	Data    map[string]string `json:"data"`
}

func GetThenDecodeAndQueryTransactionsWhoseLabelIsNotTapdAssetMintingOut(token string) (*[]PostGetRawTransactionResultSat, error) {
	getTransactions, err := GetTransactionsWhoseLabelIsNotTapdAssetMinting()
	if err != nil {
		return nil, err
	}
	var rawTransactions []string
	for _, transaction := range *getTransactions {
		if transaction.Amount < 0 {
			rawTransactions = append(rawTransactions, transaction.RawTxHex)
		}
	}
	decodedAndQueryTransactions, err := DecodeAndQueryTransactionsWhoseLabelIsNotTapdAssetMinting(token, rawTransactions)
	if err != nil {
		return nil, err
	}
	if decodedAndQueryTransactions.Error != "" {
		return nil, errors.New(decodedAndQueryTransactions.Error)
	}
	btcResult := ProcessDecodedAndQueryTransactionsData(decodedAndQueryTransactions.Data)
	result := ProcessPostGetRawTransactionResultToUseSat(btcResult)
	return result, nil
}

func GetThenDecodeAndQueryTransactionsWhoseLabelIsNotTapdAssetMintingIn(token string) (*[]PostGetRawTransactionResultSat, error) {
	getTransactions, err := GetTransactionsWhoseLabelIsNotTapdAssetMinting()
	if err != nil {
		return nil, err
	}
	var rawTransactions []string
	for _, transaction := range *getTransactions {
		if transaction.Amount >= 0 {
			rawTransactions = append(rawTransactions, transaction.RawTxHex)
		}
	}
	decodedAndQueryTransactions, err := DecodeAndQueryTransactionsWhoseLabelIsNotTapdAssetMinting(token, rawTransactions)
	if err != nil {
		return nil, err
	}
	if decodedAndQueryTransactions.Error != "" {
		return nil, errors.New(decodedAndQueryTransactions.Error)
	}
	btcResult := ProcessDecodedAndQueryTransactionsData(decodedAndQueryTransactions.Data)
	result := ProcessPostGetRawTransactionResultToUseSat(btcResult)
	return result, nil
}

type BtcTransferOutInfo struct {
	Address string                          `json:"address"`
	Value   int                             `json:"value"`
	Time    int                             `json:"time"`
	Detail  *PostGetRawTransactionResultSat `json:"detail"`
}

type BtcTransferOutInfoSimplified struct {
	Address string                  `json:"address"`
	Value   int                     `json:"value"`
	Time    int                     `json:"time"`
	Detail  *TransactionsSimplified `json:"detail"`
}

type BtcTransferInInfo struct {
	Address string                          `json:"address"`
	Value   int                             `json:"value"`
	Time    int                             `json:"time"`
	Detail  *PostGetRawTransactionResultSat `json:"detail"`
}

type BtcTransferInInfoSimplified struct {
	Address string                  `json:"address"`
	Value   int                     `json:"value"`
	Time    int                     `json:"time"`
	Detail  *TransactionsSimplified `json:"detail"`
}

func GetAllAddresses() ([]string, error) {
	var result []string
	listAddress, err := ListAddressesAndGetResponse()
	if err != nil {
		return nil, err
	}
	for _, accountWithAddresse := range listAddress.AccountWithAddresses {
		addresses := accountWithAddresse.Addresses
		for _, address := range addresses {
			result = append(result, address.Address)
		}
	}
	return result, nil
}

func GetBtcTransferOutInfos(token string) (*[]BtcTransferOutInfoSimplified, error) {
	var btcTransferOutInfos []BtcTransferOutInfo
	addresses, err := GetAllAddresses()
	if err != nil {
		return nil, err
	}
	transactions, err := GetThenDecodeAndQueryTransactionsWhoseLabelIsNotTapdAssetMintingOut(token)
	if err != nil {
		return nil, err
	}
	for _, transaction := range *transactions {
		for _, vin := range transaction.Vin {
			vinAddress := vin.Prevout.ScriptPubKey.Address
			for _, address := range addresses {
				if vinAddress == address {
					btcTransferOutInfos = append(btcTransferOutInfos, BtcTransferOutInfo{
						Address: vinAddress,
						Value:   vin.Prevout.Value,
						Time:    transaction.Time,
						Detail:  &transaction,
					})
				}
			}
		}
	}
	transactionsSimplified := BtcTransferOutInfoToBtcTransferOutInfoSimplified(&btcTransferOutInfos)
	return transactionsSimplified, nil
}

func GetBtcTransferInInfos(token string) (*[]BtcTransferInInfoSimplified, error) {
	var btcTransferInInfos []BtcTransferInInfo
	addresses, err := GetAllAddresses()
	if err != nil {
		return nil, err
	}
	transactions, err := GetThenDecodeAndQueryTransactionsWhoseLabelIsNotTapdAssetMintingIn(token)
	if err != nil {
		return nil, err
	}
	for _, transaction := range *transactions {
		for _, out := range transaction.Vout {
			voutAddress := out.ScriptPubKey.Address
			for _, address := range addresses {
				if voutAddress == address {
					btcTransferInInfos = append(btcTransferInInfos, BtcTransferInInfo{
						Address: voutAddress,
						Value:   out.Value,
						Time:    transaction.Time,
						Detail:  &transaction,
					})
				}
			}
		}
	}
	transactionsSimplified := BtcTransferOutInfoToBtcTransferInInfoSimplified(&btcTransferInInfos)
	return transactionsSimplified, nil
}

func BtcTransferOutInfoToTransactionsSimplified(btcTransferOutInfos *[]BtcTransferOutInfo) *[]TransactionsSimplified {
	var transactionsSimplified []TransactionsSimplified
	for _, btcTransferOutInfo := range *btcTransferOutInfos {
		feeRate := RoundToDecimalPlace(float64(btcTransferOutInfo.Detail.Fee)/float64(btcTransferOutInfo.Detail.Vsize), 2)
		var transactionsSimplifiedVin []TransactionsSimplifiedVin
		var transactionsSimplifiedVout []TransactionsSimplifiedVout
		for _, vin := range btcTransferOutInfo.Detail.Vin {
			transactionsSimplifiedVin = append(transactionsSimplifiedVin, TransactionsSimplifiedVin{
				ScriptpubkeyAddress: vin.Prevout.ScriptPubKey.Address,
				Value:               vin.Prevout.Value,
			})
		}
		for _, vout := range btcTransferOutInfo.Detail.Vout {
			transactionsSimplifiedVout = append(transactionsSimplifiedVout, TransactionsSimplifiedVout{
				ScriptpubkeyAddress: vout.ScriptPubKey.Address,
				Value:               vout.Value,
			})
		}
		transactionsSimplified = append(transactionsSimplified, TransactionsSimplified{
			Txid:            btcTransferOutInfo.Detail.Txid,
			Vin:             transactionsSimplifiedVin,
			Vout:            transactionsSimplifiedVout,
			BlockTime:       btcTransferOutInfo.Detail.Blocktime,
			BalanceResult:   -(btcTransferOutInfo.Value),
			FeeRate:         feeRate,
			Fee:             btcTransferOutInfo.Detail.Fee,
			ConfirmedBlocks: btcTransferOutInfo.Detail.Confirmations,
		})
	}
	return &transactionsSimplified
}

func BtcTransferOutInfoToBtcTransferOutInfoSimplified(btcTransferOutInfos *[]BtcTransferOutInfo) *[]BtcTransferOutInfoSimplified {
	var btcTransferOutInfoSimplified []BtcTransferOutInfoSimplified
	for _, btcTransferOutInfo := range *btcTransferOutInfos {
		var transactionsSimplified TransactionsSimplified
		var postGetRawTransactionResultSat PostGetRawTransactionResultSat
		postGetRawTransactionResultSat = *btcTransferOutInfo.Detail
		feeRate := RoundToDecimalPlace(float64(postGetRawTransactionResultSat.Fee)/float64(postGetRawTransactionResultSat.Vsize), 2)
		var transactionsSimplifiedVin []TransactionsSimplifiedVin
		var transactionsSimplifiedVout []TransactionsSimplifiedVout
		for _, vin := range postGetRawTransactionResultSat.Vin {
			transactionsSimplifiedVin = append(transactionsSimplifiedVin, TransactionsSimplifiedVin{
				ScriptpubkeyAddress: vin.Prevout.ScriptPubKey.Address,
				Value:               vin.Prevout.Value,
			})
		}
		for _, vout := range postGetRawTransactionResultSat.Vout {
			transactionsSimplifiedVout = append(transactionsSimplifiedVout, TransactionsSimplifiedVout{
				ScriptpubkeyAddress: vout.ScriptPubKey.Address,
				Value:               vout.Value,
			})
		}
		transactionsSimplified = TransactionsSimplified{
			Txid:            postGetRawTransactionResultSat.Txid,
			Vin:             transactionsSimplifiedVin,
			Vout:            transactionsSimplifiedVout,
			BlockTime:       postGetRawTransactionResultSat.Blocktime,
			BalanceResult:   -(btcTransferOutInfo.Value),
			FeeRate:         feeRate,
			Fee:             postGetRawTransactionResultSat.Fee,
			ConfirmedBlocks: postGetRawTransactionResultSat.Confirmations,
		}
		btcTransferOutInfoSimplified = append(btcTransferOutInfoSimplified, BtcTransferOutInfoSimplified{
			Address: btcTransferOutInfo.Address,
			Value:   btcTransferOutInfo.Value,
			Time:    btcTransferOutInfo.Time,
			Detail:  &transactionsSimplified,
		})
	}
	return &btcTransferOutInfoSimplified
}

func BtcTransferOutInfoToBtcTransferInInfoSimplified(btcTransferInInfos *[]BtcTransferInInfo) *[]BtcTransferInInfoSimplified {
	var btcTransferOutInfoSimplified []BtcTransferInInfoSimplified
	for _, btcTransferOutInfo := range *btcTransferInInfos {
		var transactionsSimplified TransactionsSimplified
		var postGetRawTransactionResultSat PostGetRawTransactionResultSat
		postGetRawTransactionResultSat = *btcTransferOutInfo.Detail
		feeRate := RoundToDecimalPlace(float64(postGetRawTransactionResultSat.Fee)/float64(postGetRawTransactionResultSat.Vsize), 2)
		var transactionsSimplifiedVin []TransactionsSimplifiedVin
		var transactionsSimplifiedVout []TransactionsSimplifiedVout
		for _, vin := range postGetRawTransactionResultSat.Vin {
			transactionsSimplifiedVin = append(transactionsSimplifiedVin, TransactionsSimplifiedVin{
				ScriptpubkeyAddress: vin.Prevout.ScriptPubKey.Address,
				Value:               vin.Prevout.Value,
			})
		}
		for _, vout := range postGetRawTransactionResultSat.Vout {
			transactionsSimplifiedVout = append(transactionsSimplifiedVout, TransactionsSimplifiedVout{
				ScriptpubkeyAddress: vout.ScriptPubKey.Address,
				Value:               vout.Value,
			})
		}
		transactionsSimplified = TransactionsSimplified{
			Txid:            postGetRawTransactionResultSat.Txid,
			Vin:             transactionsSimplifiedVin,
			Vout:            transactionsSimplifiedVout,
			BlockTime:       postGetRawTransactionResultSat.Blocktime,
			BalanceResult:   +(btcTransferOutInfo.Value),
			FeeRate:         feeRate,
			Fee:             postGetRawTransactionResultSat.Fee,
			ConfirmedBlocks: postGetRawTransactionResultSat.Confirmations,
		}
		btcTransferOutInfoSimplified = append(btcTransferOutInfoSimplified, BtcTransferInInfoSimplified{
			Address: btcTransferOutInfo.Address,
			Value:   btcTransferOutInfo.Value,
			Time:    btcTransferOutInfo.Time,
			Detail:  &transactionsSimplified,
		})
	}
	return &btcTransferOutInfoSimplified
}

func GetBtcTransferOutInfosJsonResult(token string) string {
	response, err := GetBtcTransferOutInfos(token)
	if err != nil {
		return MakeJsonErrorResult(GetBtcTransferOutInfosErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func GetBtcTransferInInfosJsonResult(token string) string {
	response, err := GetBtcTransferInInfos(token)
	if err != nil {
		return MakeJsonErrorResult(GetBtcTransferOutInfosErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

type AssetTransferType int

const (
	AssetTransferTypeOut AssetTransferType = iota
	AssetTransferTypeIn
)

type AssetTransferSetRequest struct {
	AssetID           string            `json:"asset_id" gorm:"type:varchar(255)"`
	AssetAddressFrom  string            `json:"address_from" gorm:"type:varchar(255)"`
	AssetAddressTo    string            `json:"address_to" gorm:"type:varchar(255)"`
	Amount            int               `json:"amount"`
	TransferType      AssetTransferType `json:"transfer_type"`
	TransactionID     string            `json:"transaction_id" gorm:"type:varchar(255)"`
	TransferTimestamp int               `json:"transfer_timestamp"`
	AnchorTxChainFees int               `json:"anchor_tx_chain_fees"`
}

func ListTransfersAndGetResponse() (*taprpc.ListTransfersResponse, error) {
	response, err := rpcclient.ListTransfers()
	if err != nil {
		return nil, err
	}
	return response, nil
}

type AssetTransferProcessed struct {
	Txid               string                         `json:"txid"`
	AssetID            string                         `json:"asset_id"`
	TransferTimestamp  int                            `json:"transfer_timestamp"`
	AnchorTxHash       string                         `json:"anchor_tx_hash"`
	AnchorTxHeightHint int                            `json:"anchor_tx_height_hint"`
	AnchorTxChainFees  int                            `json:"anchor_tx_chain_fees"`
	Inputs             []AssetTransferProcessedInput  `json:"inputs"`
	Outputs            []AssetTransferProcessedOutput `json:"outputs"`
	DeviceID           string                         `json:"device_id"`
}

type AssetTransferProcessedInput struct {
	Address     string `json:"address"`
	Amount      int    `json:"amount"`
	AnchorPoint string `json:"anchor_point"`
	ScriptKey   string `json:"script_key"`
}

type AssetTransferProcessedOutput struct {
	Address                string `json:"address"`
	Amount                 int    `json:"amount"`
	AnchorOutpoint         string `json:"anchor_outpoint"`
	AnchorValue            int    `json:"anchor_value"`
	AnchorInternalKey      string `json:"anchor_internal_key"`
	AnchorTaprootAssetRoot string `json:"anchor_taproot_asset_root"`
	AnchorMerkleRoot       string `json:"anchor_merkle_root"`
	AnchorTapscriptSibling string `json:"anchor_tapscript_sibling"`
	AnchorNumPassiveAssets int    `json:"anchor_num_passive_assets"`
	ScriptKey              string `json:"script_key"`
	ScriptKeyIsLocal       bool   `json:"script_key_is_local"`
	NewProofBlob           string `json:"new_proof_blob"`
	SplitCommitRootHash    string `json:"split_commit_root_hash"`
	OutputType             string `json:"output_type"`
	AssetVersion           string `json:"asset_version"`
}

func PostToSetAssetTransfer(token string, assetTransferSetRequest *[]AssetTransferProcessed) (*JsonResult, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_transfer/set"
	requestJsonBytes, err := json.Marshal(assetTransferSetRequest)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response JsonResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return &response, nil
}

type PostToGetAssetTransferTxidsResponse struct {
	Success bool     `json:"success"`
	Error   string   `json:"error"`
	Code    ErrCode  `json:"code"`
	Data    []string `json:"data"`
}

func RequestToGetAssetTransferTxids(token string) (txids []string, err error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_transfer/get/txids"
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response PostToGetAssetTransferTxidsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

func GetTxidFromOutpoint(outpoint string) (string, error) {
	txid, indexStr := getTransactionAndIndexByOutpoint(outpoint)
	if txid == "" || indexStr == "" {
		return "", errors.New("txid or index is empty")
	}
	return txid, nil
}

func GetAllOutPointsOfListTransfersResponse(listTransfersResponse *taprpc.ListTransfersResponse) []string {
	var allOutPoints []string
	for _, listTransfer := range listTransfersResponse.Transfers {
		for _, input := range listTransfer.Inputs {
			allOutPoints = append(allOutPoints, input.AnchorPoint)
		}
		for _, output := range listTransfer.Outputs {
			allOutPoints = append(allOutPoints, output.Anchor.Outpoint)
		}
	}
	return allOutPoints
}

func GetAllOutPointsOfAssetTransfersResponse(assetTransfersResponse *[]Transfer) []string {
	var allOutPoints []string
	if assetTransfersResponse == nil {
		return allOutPoints
	}
	for _, assetTransfer := range *assetTransfersResponse {
		for _, input := range assetTransfer.Inputs {
			allOutPoints = append(allOutPoints, input.AnchorPoint)
		}
		for _, output := range assetTransfer.Outputs {
			allOutPoints = append(allOutPoints, output.Anchor.Outpoint)
		}
	}
	return allOutPoints
}

func PostCallBitcoindToQueryAddressByOutpoints(token string, outpoints []string) (*GetAddressesByOutpointSliceResponse, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	network := base.NetWork
	url := serverDomainOrSocket + "/bitcoind/" + network + "/address/outpoints"
	requestStr := OutpointSliceToRequestBodyRawString(outpoints)
	payload := strings.NewReader(requestStr)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetAddressesByOutpointSliceResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func ProcessListTransfersResponse(token string, listTransfersResponse *taprpc.ListTransfersResponse, deviceId string) *[]AssetTransferProcessed {
	var assetTransferProcessed []AssetTransferProcessed
	allOutpoints := GetAllOutPointsOfListTransfersResponse(listTransfersResponse)
	response, err := PostCallBitcoindToQueryAddressByOutpoints(token, allOutpoints)
	if err != nil {
		return nil
	}
	addressMap := response.Data
	for _, listTransfer := range listTransfersResponse.Transfers {
		var txid string
		txid, err = GetTxidFromOutpoint(listTransfer.Outputs[0].Anchor.Outpoint)
		if err != nil {
			return nil
		}
		var assetTransferProcessedInput []AssetTransferProcessedInput
		for _, input := range listTransfer.Inputs {
			inOp := input.AnchorPoint
			assetTransferProcessedInput = append(assetTransferProcessedInput, AssetTransferProcessedInput{
				Address:     addressMap[inOp],
				Amount:      int(input.Amount),
				AnchorPoint: inOp,
				ScriptKey:   hex.EncodeToString(input.ScriptKey),
			})
		}
		var assetTransferProcessedOutput []AssetTransferProcessedOutput
		for _, output := range listTransfer.Outputs {
			outOp := output.Anchor.Outpoint
			outputType, ok := taprpc.OutputType_name[int32(output.OutputType)]
			if !ok {
				outputType = strconv.Itoa(int(int32(output.OutputType)))
			}
			assetTransferProcessedOutput = append(assetTransferProcessedOutput, AssetTransferProcessedOutput{
				Address:                addressMap[outOp],
				Amount:                 int(output.Amount),
				AnchorOutpoint:         outOp,
				AnchorValue:            int(output.Anchor.Value),
				AnchorInternalKey:      hex.EncodeToString(output.Anchor.InternalKey),
				AnchorTaprootAssetRoot: hex.EncodeToString(output.Anchor.TaprootAssetRoot),
				AnchorMerkleRoot:       hex.EncodeToString(output.Anchor.MerkleRoot),
				AnchorTapscriptSibling: hex.EncodeToString(output.Anchor.TapscriptSibling),
				AnchorNumPassiveAssets: int(output.Anchor.NumPassiveAssets),
				ScriptKey:              hex.EncodeToString(output.ScriptKey),
				ScriptKeyIsLocal:       output.ScriptKeyIsLocal,
				NewProofBlob:           hex.EncodeToString(output.NewProofBlob),
				SplitCommitRootHash:    hex.EncodeToString(output.SplitCommitRootHash),
				OutputType:             outputType,
				AssetVersion:           output.AssetVersion.String(),
			})
		}
		assetTransferProcessed = append(assetTransferProcessed, AssetTransferProcessed{
			Txid:               txid,
			AssetID:            hex.EncodeToString(listTransfer.Inputs[0].AssetId),
			TransferTimestamp:  int(listTransfer.TransferTimestamp),
			AnchorTxHash:       hex.EncodeToString(listTransfer.AnchorTxHash),
			AnchorTxHeightHint: int(listTransfer.AnchorTxHeightHint),
			AnchorTxChainFees:  int(listTransfer.AnchorTxChainFees),
			Inputs:             assetTransferProcessedInput,
			Outputs:            assetTransferProcessedOutput,
			DeviceID:           deviceId,
		})
	}
	return &assetTransferProcessed
}

func ListTransfersAndGetProcessedResponse(token string, deviceId string) (*[]AssetTransferProcessed, error) {
	listTransfers, err := ListTransfersAndGetResponse()
	if err != nil {
		return nil, AppendErrorInfo(err, "ListTransfersAndGetResponse")
	}
	processedListTransfers := ProcessListTransfersResponse(token, listTransfers, deviceId)
	return processedListTransfers, nil
}

func ListAndPostToSetAssetTransfers(token string, deviceId string) string {
	transfers, err := ListTransfersAndGetProcessedResponse(token, deviceId)
	if err != nil {
		return MakeJsonErrorResult(ListTransfersAndGetProcessedResponseErr, err.Error(), nil)
	}
	if transfers == nil || len(*transfers) == 0 {
		return MakeJsonErrorResult(SUCCESS, "", nil)
	}
	_, err = PostToSetAssetTransfer(token, transfers)
	if err != nil {
		return MakeJsonErrorResult(PostToSetAssetTransferErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", nil)
}

type GetAssetTransferResponse struct {
	Success bool                      `json:"success"`
	Error   string                    `json:"error"`
	Code    ErrCode                   `json:"code"`
	Data    *[]AssetTransferProcessed `json:"data"`
}

func RequestToGetAssetTransferAndGetResponse(token string) (*GetAssetTransferResponse, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_transfer/get"
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetAssetTransferResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return &response, nil
}

func RequestToGetAssetTransferByAssetIdAndGetResponse(token string, assetId string) (*GetAssetTransferResponse, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_transfer/get/" + assetId
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetAssetTransferResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return &response, nil
}

func PostToGetAssetTransfer(token string) string {
	response, err := RequestToGetAssetTransferAndGetResponse(token)
	if err != nil {
		return MakeJsonErrorResult(PostToGetAssetTransferAndGetResponseErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, response.Data)
}

type AssetTransferProcessedSimplified struct {
	Txid               string                                    `json:"txid"`
	AssetID            string                                    `json:"asset_id"`
	TransferTimestamp  int                                       `json:"transfer_timestamp"`
	AnchorTxHeightHint int                                       `json:"anchor_tx_height_hint"`
	AnchorTxChainFees  int                                       `json:"anchor_tx_chain_fees"`
	Inputs             *[]AssetTransferProcessedInputSimplified  `json:"inputs"`
	Outputs            *[]AssetTransferProcessedOutputSimplified `json:"outputs"`
	DeviceID           string                                    `json:"device_id"`
}

type AssetTransferProcessedInputSimplified struct {
	Address     string `json:"address"`
	Amount      int    `json:"amount"`
	AnchorPoint string `json:"anchor_point"`
	ScriptKey   string `json:"script_key"`
}

type AssetTransferProcessedOutputSimplified struct {
	Address          string `json:"address"`
	Amount           int    `json:"amount"`
	AnchorOutpoint   string `json:"anchor_outpoint"`
	AnchorValue      int    `json:"anchor_value"`
	ScriptKey        string `json:"script_key"`
	ScriptKeyIsLocal bool   `json:"script_key_is_local"`
	OutputType       string `json:"output_type"`
	AssetVersion     string `json:"asset_version"`
}

type AssetTransferProcessedOutputSimplifiedAndTotalAmount struct {
	OutputSimplified *[]AssetTransferProcessedOutputSimplified `json:"output_simplified"`
	TotalAmount      int                                       `json:"total_amount"`
}

type AssetTransferProcessedSimplifiedResponse struct {
	AssetID     string                            `json:"asset_id"`
	Txid        string                            `json:"txid"`
	TotalAmount int                               `json:"totalAmount"`
	Time        int                               `json:"time"`
	Detail      *AssetTransferProcessedSimplified `json:"detail"`
}

func SimplifyAssetTransferProcessedInput(assetTransferProcessedInput *[]AssetTransferProcessedInput) *[]AssetTransferProcessedInputSimplified {
	var assetTransferProcessedInputSimplified []AssetTransferProcessedInputSimplified
	for _, processedInput := range *assetTransferProcessedInput {
		assetTransferProcessedInputSimplified = append(assetTransferProcessedInputSimplified, AssetTransferProcessedInputSimplified{
			Address:     processedInput.Address,
			Amount:      processedInput.Amount,
			AnchorPoint: processedInput.AnchorPoint,
			ScriptKey:   processedInput.ScriptKey,
		})
	}
	return &assetTransferProcessedInputSimplified
}

func SimplifyAssetTransferProcessedOutput(assetTransferProcessedOutput *[]AssetTransferProcessedOutput) *[]AssetTransferProcessedOutputSimplified {
	var assetTransferProcessedOutputSimplified []AssetTransferProcessedOutputSimplified
	for _, processedOutput := range *assetTransferProcessedOutput {
		assetTransferProcessedOutputSimplified = append(assetTransferProcessedOutputSimplified, AssetTransferProcessedOutputSimplified{
			Address:          processedOutput.Address,
			Amount:           processedOutput.Amount,
			AnchorOutpoint:   processedOutput.AnchorOutpoint,
			AnchorValue:      processedOutput.AnchorValue,
			ScriptKey:        processedOutput.ScriptKey,
			ScriptKeyIsLocal: processedOutput.ScriptKeyIsLocal,
			OutputType:       processedOutput.OutputType,
			AssetVersion:     processedOutput.AssetVersion,
		})
	}
	return &assetTransferProcessedOutputSimplified
}

func SimplifyAssetTransferProcessedOutputAndTotalAmount(assetTransferProcessedOutput *[]AssetTransferProcessedOutput) *AssetTransferProcessedOutputSimplifiedAndTotalAmount {
	var assetTransferProcessedOutputSimplified []AssetTransferProcessedOutputSimplified
	var totalAmount int
	for _, processedOutput := range *assetTransferProcessedOutput {
		assetTransferProcessedOutputSimplified = append(assetTransferProcessedOutputSimplified, AssetTransferProcessedOutputSimplified{
			Address:          processedOutput.Address,
			Amount:           processedOutput.Amount,
			AnchorOutpoint:   processedOutput.AnchorOutpoint,
			AnchorValue:      processedOutput.AnchorValue,
			ScriptKey:        processedOutput.ScriptKey,
			ScriptKeyIsLocal: processedOutput.ScriptKeyIsLocal,
			OutputType:       processedOutput.OutputType,
			AssetVersion:     processedOutput.AssetVersion,
		})
		totalAmount += processedOutput.Amount
	}
	return &AssetTransferProcessedOutputSimplifiedAndTotalAmount{
		OutputSimplified: &assetTransferProcessedOutputSimplified,
		TotalAmount:      totalAmount,
	}
}

func AssetTransferProcessedToSimplifiedResponse(assetTransferProcessed *[]AssetTransferProcessed) *[]AssetTransferProcessedSimplifiedResponse {
	if assetTransferProcessed == nil {
		return nil
	}
	var assetTransferProcessedSimplifiedResponse []AssetTransferProcessedSimplifiedResponse
	for _, transferProcessed := range *assetTransferProcessed {
		inputs := SimplifyAssetTransferProcessedInput(&(transferProcessed.Inputs))
		outputs := SimplifyAssetTransferProcessedOutputAndTotalAmount(&(transferProcessed.Outputs))
		assetTransferProcessedSimplified := AssetTransferProcessedSimplified{
			Txid:               transferProcessed.Txid,
			AssetID:            transferProcessed.AssetID,
			TransferTimestamp:  transferProcessed.TransferTimestamp,
			AnchorTxHeightHint: transferProcessed.AnchorTxHeightHint,
			AnchorTxChainFees:  transferProcessed.AnchorTxChainFees,
			Inputs:             inputs,
			Outputs:            outputs.OutputSimplified,
			DeviceID:           transferProcessed.DeviceID,
		}
		assetTransferProcessedSimplifiedResponse = append(assetTransferProcessedSimplifiedResponse, AssetTransferProcessedSimplifiedResponse{
			AssetID:     assetTransferProcessedSimplified.AssetID,
			Txid:        assetTransferProcessedSimplified.Txid,
			TotalAmount: outputs.TotalAmount,
			Time:        assetTransferProcessedSimplified.TransferTimestamp,
			Detail:      &assetTransferProcessedSimplified,
		})
	}
	return &assetTransferProcessedSimplifiedResponse
}

func PostToGetAssetTransferByAssetId(token string, assetId string) string {
	response, err := RequestToGetAssetTransferByAssetIdAndGetResponse(token, assetId)
	if err != nil {
		return MakeJsonErrorResult(PostToGetAssetTransferByAssetIdAndGetResponseErr, err.Error(), nil)
	}
	var result *[]AssetTransferProcessedSimplifiedResponse
	if response == nil {
		result = nil
	} else {
		result = AssetTransferProcessedToSimplifiedResponse(response.Data)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, result)
}

func UploadAssetTransfer(token string, deviceId string) string {
	return ListAndPostToSetAssetTransfers(token, deviceId)
}

func GetAssetTransfer(token string) string {
	return PostToGetAssetTransfer(token)
}

func GetAssetTransferByAssetIdFromServer(token string, assetId string) string {
	return PostToGetAssetTransferByAssetId(token, assetId)
}

func outpointToTransactionAndIndex(outpoint string) (transaction string, index string) {
	result := strings.Split(outpoint, ":")
	return result[0], result[1]
}

func BatchTxidToAssetId(batchTxid string) (string, error) {
	assets, _ := listAssets(true, true, false)
	for _, asset := range assets.Assets {
		txid, _ := outpointToTransactionAndIndex(asset.GetChainAnchor().GetAnchorOutpoint())
		if batchTxid == txid {
			return hex.EncodeToString(asset.GetAssetGenesis().AssetId), nil
		}
	}
	err := errors.New("no asset found for batch txid")
	return "", err
}

func QueryAssetIdByBatchTxid(batchTxid string) string {
	assetId, err := BatchTxidToAssetId(batchTxid)
	if err != nil {
		return MakeJsonErrorResult(BatchTxidToAssetIdErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, assetId)
}

func AddrReceivesAndGetResponse() (*taprpc.AddrReceivesResponse, error) {
	return rpcclient.AddrReceives()
}

type AddrReceiveEvent struct {
	CreationTimeUnixSeconds int                  `json:"creation_time_unix_seconds"`
	Addr                    AddrReceiveEventAddr `json:"addr"`
	Status                  string               `json:"status"`
	Outpoint                string               `json:"outpoint"`
	UtxoAmtSat              int                  `json:"utxo_amt_sat"`
	ConfirmationHeight      int                  `json:"confirmation_height"`
	HasProof                bool                 `json:"has_proof,omitempty"`
	DeviceID                string               `json:"device_id"`
}

type AddrReceiveEventAddr struct {
	Encoded          string `json:"encoded"`
	AssetID          string `json:"asset_id"`
	Amount           int    `json:"amount"`
	ScriptKey        string `json:"script_key"`
	InternalKey      string `json:"internal_key"`
	TaprootOutputKey string `json:"taproot_output_key"`
	ProofCourierAddr string `json:"proof_courier_addr"`
}

func AddrReceivesResponseToAddrReceiveEvents(addrReceivesResponse *taprpc.AddrReceivesResponse, deviceId string) *[]AddrReceiveEvent {
	var addrReceiveEvents []AddrReceiveEvent
	for _, event := range addrReceivesResponse.Events {
		addrReceiveEvents = append(addrReceiveEvents, AddrReceiveEvent{
			CreationTimeUnixSeconds: int(event.CreationTimeUnixSeconds),
			Addr: AddrReceiveEventAddr{
				Encoded:          event.Addr.Encoded,
				AssetID:          hex.EncodeToString(event.Addr.AssetId),
				Amount:           int(event.Addr.Amount),
				ScriptKey:        hex.EncodeToString(event.Addr.ScriptKey),
				InternalKey:      hex.EncodeToString(event.Addr.InternalKey),
				TaprootOutputKey: hex.EncodeToString(event.Addr.TaprootOutputKey),
				ProofCourierAddr: event.Addr.ProofCourierAddr,
			},
			Status:             event.Status.String(),
			Outpoint:           event.Outpoint,
			UtxoAmtSat:         int(event.UtxoAmtSat),
			ConfirmationHeight: int(event.ConfirmationHeight),
			HasProof:           event.HasProof,
			DeviceID:           deviceId,
		})
	}
	return &addrReceiveEvents
}

func AddrReceivesAndGetEvents(deviceId string) (*[]AddrReceiveEvent, error) {
	response, err := AddrReceivesAndGetResponse()
	if err != nil {
		return nil, err
	}
	return AddrReceivesResponseToAddrReceiveEvents(response, deviceId), nil
}

func PostToSetAddrReceivesEvents(token string, addrReceiveEvents *[]AddrReceiveEvent) error {
	if addrReceiveEvents == nil || len(*addrReceiveEvents) == 0 {
		return nil
	}
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/addr_receive/set"
	requestJsonBytes, err := json.Marshal(addrReceiveEvents)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return errors.Wrap(err, "http.NewRequest")
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, " http.DefaultClient.Do")
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, " io.ReadAll")
	}
	var response JsonResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		return errors.Wrap(err, " json.Unmarshal")
	}
	if response.Error != "" {
		return errors.New(response.Error)
	}
	return nil
}

type GetAddrReceivesEventsResponse struct {
	Success bool                `json:"success"`
	Error   string              `json:"error"`
	Code    ErrCode             `json:"code"`
	Data    *[]AddrReceiveEvent `json:"data"`
}

func RequestToGetAddrReceivesEvents(token string) (*[]AddrReceiveEvent, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/addr_receive/get/origin"
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetAddrReceivesEventsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

func UploadAddrReceives(token string, deviceId string) string {
	events, err := AddrReceivesAndGetEvents(deviceId)
	if err != nil {
		return MakeJsonErrorResult(AddrReceivesAndGetEventsErr, err.Error(), nil)
	}
	err = PostToSetAddrReceivesEvents(token, events)
	if err != nil {
		return MakeJsonErrorResult(PostToSetAddrReceivesEventsErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, nil)
}

func GetAddrReceives(token string) string {
	response, err := RequestToGetAddrReceivesEvents(token)
	if err != nil {
		return MakeJsonErrorResult(PostToGetAddrReceivesEventsErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, response)
}

type BatchTransferRequest struct {
	Encoded            string `json:"encoded"`
	AssetID            string `json:"asset_id"`
	Amount             int    `json:"amount"`
	ScriptKey          string `json:"script_key"`
	InternalKey        string `json:"internal_key"`
	TaprootOutputKey   string `json:"taproot_output_key"`
	ProofCourierAddr   string `json:"proof_courier_addr"`
	Txid               string `json:"txid"`
	TxTotalAmount      int    `json:"tx_total_amount"`
	Index              int    `json:"index"`
	TransferTimestamp  int    `json:"transfer_timestamp"`
	AnchorTxHash       string `json:"anchor_tx_hash"`
	AnchorTxHeightHint int    `json:"anchor_tx_height_hint"`
	AnchorTxChainFees  int    `json:"anchor_tx_chain_fees"`
	DeviceID           string `json:"device_id"`
}

func PostToSetBatchTransfers(token string, batchTransfers *[]BatchTransferRequest) (err error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/batch_transfer/set_slice"
	requestJsonBytes, err := json.Marshal(batchTransfers)
	if err != nil {
		return err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var response JsonResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	if response.Error != "" {
		return errors.New(response.Error)
	}
	return nil
}

func RequestToGetBatchTransfers(token string) (*[]BatchTransfer, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/batch_transfer/get"
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetBatchTransfersResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

type GetBatchTransfersResponse struct {
	Success bool             `json:"success"`
	Error   string           `json:"error"`
	Code    ErrCode          `json:"code"`
	Data    *[]BatchTransfer `json:"data"`
}

type BatchTransfer struct {
	gorm.Model
	Encoded            string `json:"encoded"`
	AssetID            string `json:"asset_id" gorm:"type:varchar(255)"`
	Amount             int    `json:"amount"`
	ScriptKey          string `json:"script_key" gorm:"type:varchar(255)"`
	InternalKey        string `json:"internal_key" gorm:"type:varchar(255)"`
	TaprootOutputKey   string `json:"taproot_output_key" gorm:"type:varchar(255)"`
	ProofCourierAddr   string `json:"proof_courier_addr" gorm:"type:varchar(255)"`
	Txid               string `json:"txid" gorm:"type:varchar(255)"`
	TxTotalAmount      int    `json:"tx_total_amount"`
	Index              int    `json:"index"`
	TransferTimestamp  int    `json:"transfer_timestamp"`
	AnchorTxHash       string `json:"anchor_tx_hash" gorm:"type:varchar(255)"`
	AnchorTxHeightHint int    `json:"anchor_tx_height_hint"`
	AnchorTxChainFees  int    `json:"anchor_tx_chain_fees"`
	DeviceID           string `json:"device_id" gorm:"type:varchar(255)"`
	UserID             int    `json:"user_id"`
	Status             int    `json:"status" gorm:"default:1"`
}

func UploadBatchTransfers(token string, batchTransfers *[]BatchTransferRequest) (err error) {
	return PostToSetBatchTransfers(token, batchTransfers)
}

func GetUserAllBatchTransfers(token string) string {
	response, err := RequestToGetBatchTransfers(token)
	if err != nil {
		return MakeJsonErrorResult(PostToGetBatchTransfersErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, response)
}

type BatchTransferInfo struct {
	AssetId           string                     `json:"asset_id"`
	Txid              string                     `json:"txid"`
	TotalAmount       int                        `json:"total_amount"`
	TransferTimestamp int                        `json:"transfer_timestamp"`
	Details           *[]BatchTransferInfoDetail `json:"details"`
}

type BatchTransferInfoDetail struct {
	EncodedAddr string `json:"encoded_addr"`
	Amount      int    `json:"amount"`
	Index       int    `json:"index"`
}

func BatchTransfersToBatchTransferInfoDetails(batchTransfers *[]BatchTransfer) *[]BatchTransferInfoDetail {
	var result []BatchTransferInfoDetail
	for _, batchTransfer := range *batchTransfers {
		result = append(result, BatchTransferInfoDetail{
			EncodedAddr: batchTransfer.Encoded,
			Amount:      batchTransfer.Amount,
			Index:       batchTransfer.Index,
		})
	}
	return &result
}

func BatchTransfersToBatchTransferInfos(batchTransfers *[]BatchTransfer) *[]BatchTransferInfo {
	var result []BatchTransferInfo
	assetIdToBatchTransfers := SplitBatchTransfersByTxid(batchTransfers)
	for txid, transfers := range *assetIdToBatchTransfers {
		details := BatchTransfersToBatchTransferInfoDetails(transfers)
		result = append(result, BatchTransferInfo{
			AssetId:           (*transfers)[0].AssetID,
			Txid:              txid,
			TotalAmount:       (*transfers)[0].TxTotalAmount,
			TransferTimestamp: (*transfers)[0].TransferTimestamp,
			Details:           details,
		})
	}
	return &result
}

func SplitBatchTransfersByTxid(batchTransfers *[]BatchTransfer) *map[string]*[]BatchTransfer {
	result := make(map[string]*[]BatchTransfer)
	for _, batchTransfer := range *batchTransfers {
		txid := batchTransfer.Txid
		batchTransferSlice, ok := result[txid]
		if !ok {
			newBatchTransfers := &[]BatchTransfer{batchTransfer}
			result[txid] = newBatchTransfers
		} else {
			*batchTransferSlice = append(*batchTransferSlice, batchTransfer)
		}
	}
	return &result
}

func GetBatchTransfers(token string) string {
	response, err := RequestToGetBatchTransfers(token)
	result := BatchTransfersToBatchTransferInfos(response)
	if err != nil {
		return MakeJsonErrorResult(PostToGetBatchTransfersErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, result)
}

type AssetAddr struct {
	gorm.Model
	Encoded          string `json:"encoded"`
	AssetId          string `json:"asset_id" gorm:"type:varchar(255)"`
	AssetType        int    `json:"asset_type"`
	Amount           int    `json:"amount"`
	GroupKey         string `json:"group_key" gorm:"type:varchar(255)"`
	ScriptKey        string `json:"script_key" gorm:"type:varchar(255)"`
	InternalKey      string `json:"internal_key" gorm:"type:varchar(255)"`
	TapscriptSibling string `json:"tapscript_sibling" gorm:"type:varchar(255)"`
	TaprootOutputKey string `json:"taproot_output_key" gorm:"type:varchar(255)"`
	ProofCourierAddr string `json:"proof_courier_addr" gorm:"type:varchar(255)"`
	AssetVersion     int    `json:"asset_version"`
	DeviceID         string `json:"device_id" gorm:"type:varchar(255)"`
	UserId           int    `json:"user_id"`
	Username         string `json:"username" gorm:"type:varchar(255)"`
	Status           int    `json:"status" gorm:"default:1"`
}

type AssetAddrSetRequest struct {
	Encoded          string `json:"encoded"`
	AssetId          string `json:"asset_id"`
	AssetType        int    `json:"asset_type"`
	Amount           int    `json:"amount"`
	GroupKey         string `json:"group_key"`
	ScriptKey        string `json:"script_key"`
	InternalKey      string `json:"internal_key"`
	TapscriptSibling string `json:"tapscript_sibling"`
	TaprootOutputKey string `json:"taproot_output_key"`
	ProofCourierAddr string `json:"proof_courier_addr"`
	AssetVersion     int    `json:"asset_version"`
	DeviceID         string `json:"device_id"`
}

type GetAssetAddrResponse struct {
	Success bool         `json:"success"`
	Error   string       `json:"error"`
	Code    ErrCode      `json:"code"`
	Data    *[]AssetAddr `json:"data"`
}

func PostToSetAssetAddr(token string, assetAddrSetRequest *AssetAddrSetRequest) (err error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_addr/set"
	requestJsonBytes, err := json.Marshal(assetAddrSetRequest)
	if err != nil {
		return err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var response JsonResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	if response.Error != "" {
		return errors.New(response.Error)
	}
	return nil
}

func RequestToGetAssetAddr(token string) (*[]AssetAddr, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_addr/get"
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetAssetAddrResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

func RequestToGetAssetAddrByScriptKey(token string, scriptKey string) (*[]AssetAddr, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_addr/get/script_key/" + scriptKey
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetAssetAddrResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

func UploadAssetAddr(token string, assetAddrSetRequest *AssetAddrSetRequest) string {
	err := PostToSetAssetAddr(token, assetAddrSetRequest)
	if err != nil {
		return MakeJsonErrorResult(PostToSetAssetAddrErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, nil)
}

func GetAssetAddrs(token string) string {
	response, err := RequestToGetAssetAddr(token)
	if err != nil {
		return MakeJsonErrorResult(PostToGetAssetAddrErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, response)
}

func GetAssetAddrsByScriptKey(token string, scriptKey string) string {
	response, err := RequestToGetAssetAddrByScriptKey(token, scriptKey)
	if err != nil {
		return MakeJsonErrorResult(PostToGetAssetAddrErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, response)
}

type AssetLock struct {
	gorm.Model
	AssetId          string `json:"asset_id" gorm:"type:varchar(255)"`
	AssetName        string `json:"asset_name" gorm:"type:varchar(255)"`
	AssetType        string `json:"asset_type" gorm:"type:varchar(255)"`
	LockAmount       int    `json:"lock_amount"`
	LockTime         int    `json:"lock_time"`
	RelativeLockTime int    `json:"relative_lock_time"`
	HashLock         string `json:"hash_lock" gorm:"type:varchar(255)"`
	Invoice          string `json:"invoice" gorm:"type:varchar(255)"`
	DeviceId         string `json:"device_id" gorm:"type:varchar(255)"`
	UserId           int    `json:"user_id"`
	Status           int    `json:"status" gorm:"default:1"`
}

type AssetLockSetRequest struct {
	AssetId          string `json:"asset_id" gorm:"type:varchar(255)"`
	AssetName        string `json:"asset_name" gorm:"type:varchar(255)"`
	AssetType        string `json:"asset_type" gorm:"type:varchar(255)"`
	LockAmount       int    `json:"lock_amount"`
	LockTime         int    `json:"lock_time"`
	RelativeLockTime int    `json:"relative_lock_time"`
	HashLock         string `json:"hash_lock" gorm:"type:varchar(255)"`
	Invoice          string `json:"invoice" gorm:"type:varchar(255)"`
	DeviceId         string `json:"device_id" gorm:"type:varchar(255)"`
}

type GetAssetLockResponse struct {
	Success bool         `json:"success"`
	Error   string       `json:"error"`
	Code    ErrCode      `json:"code"`
	Data    *[]AssetLock `json:"data"`
}

func PostToSetAssetLock(token string, assetLockSetRequest *AssetLockSetRequest) (err error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_lock/set"
	requestJsonBytes, err := json.Marshal(assetLockSetRequest)
	if err != nil {
		return err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var response JsonResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	if response.Error != "" {
		return errors.New(response.Error)
	}
	return nil
}

func RequestToGetAssetLock(token string) (*[]AssetLock, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_lock/get"
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetAssetLockResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

func UploadAssetLock(token string, assetLockSetRequest *AssetLockSetRequest) string {
	err := PostToSetAssetLock(token, assetLockSetRequest)
	if err != nil {
		return MakeJsonErrorResult(PostToSetAssetLockErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, nil)
}

func GetAssetLocks(token string) string {
	response, err := RequestToGetAssetLock(token)
	if err != nil {
		return MakeJsonErrorResult(PostToGetAssetLockErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, response)
}

type ValidateTokenResponse struct {
	Error string `json:"error"`
}

func GetValidateTokenResult(token string) (*ValidateTokenResponse, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/validate_token/ping"
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response ValidateTokenResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func IsTokenValid(token string) (bool, error) {
	response, err := GetValidateTokenResult(token)
	if err != nil {
		return false, err
	}
	if response.Error != "" {
		return false, errors.New(response.Error)
	}
	return true, nil
}

func ListBalancesAndGetResponse() (*taprpc.ListBalancesResponse, error) {
	return listBalances(false, nil, nil)
}

type ListBalanceInfo struct {
	GenesisPoint string `json:"genesis_point"`
	Name         string `json:"name"`
	MetaHash     string `json:"meta_hash"`
	AssetID      string `json:"asset_id"`
	AssetType    string `json:"asset_type"`
	OutputIndex  int    `json:"output_index"`
	Version      int    `json:"version"`
	Balance      int    `json:"balance"`
}

type ListBalanceSimpleInfo struct {
	AssetID string `json:"asset_id"`
	Balance int    `json:"balance"`
}

func ListBalancesResponseToListBalanceInfos(listBalancesResponse *taprpc.ListBalancesResponse) *[]ListBalanceInfo {
	var listBalanceInfos []ListBalanceInfo
	for _, balance := range listBalancesResponse.AssetBalances {
		listBalanceInfos = append(listBalanceInfos, ListBalanceInfo{
			GenesisPoint: balance.AssetGenesis.GenesisPoint,
			Name:         balance.AssetGenesis.Name,
			MetaHash:     hex.EncodeToString(balance.AssetGenesis.MetaHash),
			AssetID:      hex.EncodeToString(balance.AssetGenesis.AssetId),
			AssetType:    balance.AssetGenesis.AssetType.String(),
			OutputIndex:  int(balance.AssetGenesis.OutputIndex),
			Version:      -1,
			Balance:      int(balance.Balance),
		})
	}
	return &listBalanceInfos
}

func ListBalancesResponseToListBalanceSimpleInfos(listBalancesResponse *taprpc.ListBalancesResponse) *[]ListBalanceSimpleInfo {
	if listBalancesResponse == nil {
		return new([]ListBalanceSimpleInfo)
	}
	var listBalanceSimpleInfos []ListBalanceSimpleInfo
	for _, balance := range listBalancesResponse.AssetBalances {
		listBalanceSimpleInfos = append(listBalanceSimpleInfos, ListBalanceSimpleInfo{
			AssetID: hex.EncodeToString(balance.AssetGenesis.AssetId),
			Balance: int(balance.Balance),
		})
	}
	return &listBalanceSimpleInfos
}

func ListBalancesAndProcess() (*[]ListBalanceInfo, error) {
	response, err := ListBalancesAndGetResponse()
	if err != nil {
		return nil, err
	}
	processed := ListBalancesResponseToListBalanceInfos(response)
	return processed, nil
}

func GetListBalancesSimpleInfo() (*[]ListBalanceSimpleInfo, error) {
	response, err := ListBalancesAndGetResponse()
	if err != nil {
		return new([]ListBalanceSimpleInfo), err
	}
	processed := ListBalancesResponseToListBalanceSimpleInfos(response)
	sort.Slice(*processed, func(i, j int) bool {
		return (*processed)[i].AssetID < (*processed)[j].AssetID
	})
	return processed, nil
}

func GetListBalancesSimpleInfoHash() (string, error) {
	response, err := GetListBalancesSimpleInfo()
	if err != nil {
		return "", AppendErrorInfo(err, "GetListBalancesSimpleInfo")
	}
	hash, err := Sha256(response)
	if err != nil {
		return "", err
	}
	return hash, nil
}

type GetAssetBalanceBackupResponse struct {
	Success bool    `json:"success"`
	Error   string  `json:"error"`
	Code    ErrCode `json:"code"`
	Data    string  `json:"data"`
}

func RequestToGetAssetBalanceBackupHash(token string) (string, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_balance_backup/get"
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return "", err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var response GetAssetBalanceBackupResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}
	if response.Error != "" {
		return "", errors.New(response.Error)
	}
	return response.Data, nil
}

func GetAssetBalanceBackupHash(token string) (string, error) {
	hash, err := RequestToGetAssetBalanceBackupHash(token)
	if err != nil {
		return "", err
	}
	return hash, nil
}

type UpdateAssetBalanceBackupResponse struct {
	Success bool    `json:"success"`
	Error   string  `json:"error"`
	Code    ErrCode `json:"code"`
	Data    string  `json:"data"`
}

func PostToUpdateAssetBalanceBackup(token string, hash string) (string, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_balance_backup/update" + "?hash=" + hash
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return "", err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var response UpdateAssetBalanceBackupResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}
	if response.Error != "" {
		return "", errors.New(response.Error)
	}
	return response.Data, nil
}

func UpdateAssetBalanceBackupHash(token string, hash string) (string, error) {
	return PostToUpdateAssetBalanceBackup(token, hash)
}

func GetListBalancesSimpleInfoHashAndUpdateAssetBalanceBackup(token string) (string, error) {
	hash, err := GetListBalancesSimpleInfoHash()
	if err != nil {
		return "", AppendErrorInfo(err, "GetListBalancesSimpleInfoHash")
	}
	response, err := UpdateAssetBalanceBackupHash(token, hash)
	if err != nil {
		return "", AppendErrorInfo(err, "UpdateAssetBalanceBackupHash")
	}
	return response, nil
}

func UploadAssetBalanceBackupHash(token string) string {
	hash, err := GetListBalancesSimpleInfoHashAndUpdateAssetBalanceBackup(token)
	if err != nil {
		return MakeJsonErrorResult(GetListBalancesSimpleInfoHashAndUpdateAssetBalanceBackupErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, hash)
}

func CheckIfBackupIsRequired(token string) (bool, error) {
	hash, err := GetAssetBalanceBackupHash(token)
	if err != nil {
		return true, AppendErrorInfo(err, "GetAssetBalanceBackupHash")
	}
	if hash == "" {
		return true, nil
	}
	hashLocal, err := GetListBalancesSimpleInfoHash()
	return hash != hashLocal, nil
}

func CheckBackup(token string) string {
	isRequired, err := CheckIfBackupIsRequired(token)
	if err != nil {
		return MakeJsonErrorResult(CheckIfBackupIsRequiredErr, err.Error(), isRequired)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, isRequired)
}

type AssetBalanceInfo struct {
	gorm.Model
	GenesisPoint  string `json:"genesis_point" gorm:"type:varchar(255)"`
	Name          string `json:"name" gorm:"type:varchar(255)"`
	MetaHash      string `json:"meta_hash" gorm:"type:varchar(255)"`
	AssetID       string `json:"asset_id" gorm:"type:varchar(255)"`
	AssetType     string `json:"asset_type" gorm:"type:varchar(255)"`
	OutputIndex   int    `json:"output_index"`
	Version       int    `json:"version"`
	Balance       int    `json:"balance"`
	DeviceId      string `json:"device_id" gorm:"type:varchar(255)"`
	UserId        int    `json:"user_id"`
	Username      string `json:"username" gorm:"type:varchar(255)"`
	Status        int    `json:"status" gorm:"default:1"`
	FromListAsset bool   `json:"from_list_asset"`
}

type AssetBalanceSetRequest struct {
	GenesisPoint  string `json:"genesis_point"`
	Name          string `json:"name"`
	MetaHash      string `json:"meta_hash"`
	AssetID       string `json:"asset_id"`
	AssetType     string `json:"asset_type"`
	OutputIndex   int    `json:"output_index"`
	Version       int    `json:"version"`
	Balance       int    `json:"balance"`
	DeviceId      string `json:"device_id" gorm:"type:varchar(255)"`
	FromListAsset bool   `json:"from_list_asset"`
}

func PostToSetAssetBalanceInfo(assetBalanceSetRequest *[]AssetBalanceSetRequest, token string) (*JsonResult, error) {
	if assetBalanceSetRequest == nil || len(*assetBalanceSetRequest) == 0 {
		return &JsonResult{}, nil
	}
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_balance/set_slice"
	requestJsonBytes, err := json.Marshal(assetBalanceSetRequest)
	if err != nil {
		return nil, errors.Wrap(err, " json.Marshal")
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, errors.Wrap(err, " http.NewRequest")
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, " http.DefaultClient.Do")
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, " io.ReadAll")
	}
	var response JsonResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.Wrap(err, " json.Unmarshal")
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return &response, nil
}

func ListBalanceInfosToAssetBalanceSetRequests(listBalanceInfos *[]ListBalanceInfo, deviceId string, fromListAsset bool) *[]AssetBalanceSetRequest {
	var result []AssetBalanceSetRequest
	for _, listBalanceInfo := range *listBalanceInfos {
		result = append(result, AssetBalanceSetRequest{
			GenesisPoint:  listBalanceInfo.GenesisPoint,
			Name:          listBalanceInfo.Name,
			MetaHash:      listBalanceInfo.MetaHash,
			AssetID:       listBalanceInfo.AssetID,
			AssetType:     listBalanceInfo.AssetType,
			OutputIndex:   listBalanceInfo.OutputIndex,
			Version:       listBalanceInfo.Version,
			Balance:       listBalanceInfo.Balance,
			DeviceId:      deviceId,
			FromListAsset: fromListAsset,
		})
	}
	return &result
}

func UploadListBalancesProcessedInfo(token string, deviceId string) string {
	isTokenValid, err := IsTokenValid(token)
	if err != nil {
		return MakeJsonErrorResult(IsTokenValidErr, "server "+err.Error()+"; token is invalid, did not send.", nil)
	} else if !isTokenValid {
		return MakeJsonErrorResult(IsTokenValidErr, "token is invalid, did not send.", nil)
	}
	balances, err := ListBalancesAndProcess()
	if err != nil {
		return MakeJsonErrorResult(ListBalancesAndProcessErr, err.Error(), nil)
	}
	zeroBalances, err := GetZeroBalanceAssetBalanceSlice(token, balances)
	if err != nil {
		return MakeJsonErrorResult(GetZeroBalanceAssetBalanceSliceErr, err.Error(), nil)
	}
	zeroListBalance := AssetBalanceInfosToListBalanceInfos(zeroBalances)
	setBalances := append(*balances, *zeroListBalance...)
	requests := ListBalanceInfosToAssetBalanceSetRequests(&setBalances, deviceId, false)
	result, err := PostToSetAssetBalanceInfo(requests, token)
	if err != nil {
		return MakeJsonErrorResult(PostToSetAssetBalanceInfoErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, result.Data)
}

func ListAssetsResponseSliceToListBalanceInfoSlice(listAssetsResponseSlice *[]ListAssetsResponse) *[]ListBalanceInfo {
	if listAssetsResponseSlice == nil {
		return nil
	}
	var listBalanceInfos []ListBalanceInfo
	for _, listAssetResponse := range *listAssetsResponseSlice {
		listBalanceInfos = append(listBalanceInfos, ListBalanceInfo{
			GenesisPoint: listAssetResponse.AssetGenesis.GenesisPoint,
			Name:         listAssetResponse.AssetGenesis.Name,
			MetaHash:     listAssetResponse.AssetGenesis.MetaHash,
			AssetID:      listAssetResponse.AssetGenesis.AssetID,
			AssetType:    listAssetResponse.AssetGenesis.AssetType,
			OutputIndex:  listAssetResponse.AssetGenesis.OutputIndex,
			Version:      listAssetResponse.AssetGenesis.Version,
			Balance:      listAssetResponse.Amount,
		})
	}
	return &listBalanceInfos
}

func UploadListBalancesProcessedInfoFromListAsset(token string, deviceId string) string {
	isTokenValid, err := IsTokenValid(token)
	if err != nil {
		return MakeJsonErrorResult(IsTokenValidErr, "server "+err.Error()+"; token is invalid, did not send.", nil)
	} else if !isTokenValid {
		return MakeJsonErrorResult(IsTokenValidErr, "token is invalid, did not send.", nil)
	}
	response, err := ListAssetsProcessed(true, false, false)
	if err != nil {
		return MakeJsonErrorResult(ListAssetsProcessedErr, err.Error(), nil)
	}
	balances := ListAssetsResponseSliceToListBalanceInfoSlice(response)
	zeroBalances, err := GetZeroBalanceAssetBalanceSlice(token, balances)
	if err != nil {
		return MakeJsonErrorResult(GetZeroBalanceAssetBalanceSliceErr, err.Error(), nil)
	}
	zeroListBalance := AssetBalanceInfosToListBalanceInfos(zeroBalances)
	setBalances := append(*balances, *zeroListBalance...)
	requests := ListBalanceInfosToAssetBalanceSetRequests(&setBalances, deviceId, true)
	result, err := PostToSetAssetBalanceInfo(requests, token)
	if err != nil {
		return MakeJsonErrorResult(PostToSetAssetBalanceInfoErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, result.Data)
}

func UploadAssetBalanceInfo(token string, deviceId string) string {
	return UploadListBalancesProcessedInfo(token, deviceId)
}

type GetAssetBalanceInfoResponse struct {
	Success bool                `json:"success"`
	Error   string              `json:"error"`
	Code    ErrCode             `json:"code"`
	Data    *[]AssetBalanceInfo `json:"data"`
}

func RequestToGetNonZeroAssetBalance(token string) (*[]AssetBalanceInfo, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_balance/get"
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetAssetBalanceInfoResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

func GetNonZeroAssetBalanceInfo(token string) string {
	response, err := RequestToGetNonZeroAssetBalance(token)
	if err != nil {
		return MakeJsonErrorResult(RequestToGetNonZeroAssetBalanceErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, response)
}

func GetNonZeroBalanceAssetBalanceSlice(token string) (*[]AssetBalanceInfo, error) {
	var assetBalances []AssetBalanceInfo
	response, err := RequestToGetNonZeroAssetBalance(token)
	if err != nil {
		return nil, err
	}
	for _, assetBalance := range *response {
		assetBalances = append(assetBalances, assetBalance)
	}
	return &assetBalances, nil
}

func CompareToGetZeroBalanceAssetIdWithListBalance(listBalanceInfos []ListBalanceInfo, assetBalances []AssetBalanceInfo) *[]AssetBalanceInfo {
	isAssetIdsOfListBalanceInfosExists := make(map[string]bool)
	for _, listBalanceInfo := range listBalanceInfos {
		isAssetIdsOfListBalanceInfosExists[listBalanceInfo.AssetID] = true
	}
	var zeroAssetBalances []AssetBalanceInfo
	for _, assetBalance := range assetBalances {
		isExists, ok := isAssetIdsOfListBalanceInfosExists[assetBalance.AssetID]
		if !ok || isExists == false {
			assetBalance.Balance = 0
			zeroAssetBalances = append(zeroAssetBalances, assetBalance)
		}
	}
	return &zeroAssetBalances
}

func GetZeroBalanceAssetBalanceSlice(token string, listBalanceInfos *[]ListBalanceInfo) (*[]AssetBalanceInfo, error) {
	assetBalances, err := GetNonZeroBalanceAssetBalanceSlice(token)
	if err != nil {
		return nil, err
	}
	zeroAssetIds := CompareToGetZeroBalanceAssetIdWithListBalance(*listBalanceInfos, *assetBalances)
	return zeroAssetIds, nil
}

func GetZeroAmountAssetListSlice(token string, listAssetsResponse *[]ListAssetsResponse) (*[]AssetListInfo, error) {
	assetLists, err := GetNonZeroAmountAssetListSlice(token)
	if err != nil {
		return nil, err
	}
	zeroAssetIds := CompareToGetZeroAmountAssetIdWithListAssets(*listAssetsResponse, *assetLists)
	return zeroAssetIds, nil
}

func AssetBalanceInfosToListAssetsResponseSlice(assetListInfos *[]AssetListInfo) *[]ListAssetsResponse {
	var listAssetsResponse []ListAssetsResponse
	for _, assetListInfo := range *assetListInfos {
		listAssetsResponse = append(listAssetsResponse, *AssetListInfoToListAssetsResponse(&assetListInfo))
	}
	return &listAssetsResponse
}

func ListAssetsResponseSliceToAssetListSetRequests(listAssetsResponseSlice *[]ListAssetsResponse, deviceId string) *[]AssetListSetRequest {
	var result []AssetListSetRequest
	for _, listAssetsResponse := range *listAssetsResponseSlice {
		result = append(result, AssetListSetRequest{
			Version:          listAssetsResponse.Version,
			GenesisPoint:     listAssetsResponse.AssetGenesis.GenesisPoint,
			Name:             listAssetsResponse.AssetGenesis.Name,
			MetaHash:         listAssetsResponse.AssetGenesis.MetaHash,
			AssetID:          listAssetsResponse.AssetGenesis.AssetID,
			AssetType:        listAssetsResponse.AssetGenesis.AssetType,
			OutputIndex:      listAssetsResponse.AssetGenesis.OutputIndex,
			Amount:           listAssetsResponse.Amount,
			LockTime:         listAssetsResponse.LockTime,
			RelativeLockTime: listAssetsResponse.RelativeLockTime,
			ScriptKey:        listAssetsResponse.ScriptKey,
			AnchorOutpoint:   listAssetsResponse.ChainAnchor.AnchorOutpoint,
			TweakedGroupKey:  listAssetsResponse.AssetGroup.TweakedGroupKey,
			DeviceId:         deviceId,
		})
	}
	return &result
}

func PostToSetAssetListInfo(assetListSetRequests *[]AssetListSetRequest, token string) (*JsonResult, error) {
	if assetListSetRequests == nil || len(*assetListSetRequests) == 0 {
		return &JsonResult{}, nil
	}
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_list/set_slice"
	requestJsonBytes, err := json.Marshal(assetListSetRequests)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, errors.Wrap(err, "http.NewRequest")
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "http.DefaultClient.Do")
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "io.ReadAll")
	}
	var response JsonResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return &response, nil
}

func AssetListInfoToListAssetsResponse(assetListInfo *AssetListInfo) *ListAssetsResponse {
	if assetListInfo == nil {
		return nil
	}
	return &ListAssetsResponse{
		Version: assetListInfo.Version,
		AssetGenesis: ListAssetsResponseAssetGenesis{
			GenesisPoint: assetListInfo.GenesisPoint,
			Name:         assetListInfo.Name,
			MetaHash:     assetListInfo.MetaHash,
			AssetID:      assetListInfo.AssetID,
			AssetType:    assetListInfo.AssetType,
			OutputIndex:  assetListInfo.OutputIndex,
		},
		Amount:           assetListInfo.Amount,
		LockTime:         assetListInfo.LockTime,
		RelativeLockTime: assetListInfo.RelativeLockTime,
		ScriptKey:        assetListInfo.ScriptKey,
		ChainAnchor: ListAssetsResponseChainAnchor{
			AnchorOutpoint: assetListInfo.AnchorOutpoint,
		},
		AssetGroup: ListAssetsResponseAssetGroup{
			TweakedGroupKey: assetListInfo.TweakedGroupKey,
		},
	}
}

func CompareToGetZeroAmountAssetIdWithListAssets(listAssetsResponse []ListAssetsResponse, assetLists []AssetListInfo) *[]AssetListInfo {
	isAssetIdsOfListAssetsResponseExists := make(map[string]bool)
	for _, listAsset := range listAssetsResponse {
		isAssetIdsOfListAssetsResponseExists[listAsset.AssetGenesis.AssetID] = true
	}
	var zeroAmountAssetLists []AssetListInfo
	for _, assetList := range assetLists {
		isExists, ok := isAssetIdsOfListAssetsResponseExists[assetList.AssetID]
		if !ok || isExists == false {
			assetList.Amount = 0
			zeroAmountAssetLists = append(zeroAmountAssetLists, assetList)
		}
	}
	return &zeroAmountAssetLists
}

func GetNonZeroAmountAssetListSlice(token string) (*[]AssetListInfo, error) {
	var assetLists []AssetListInfo
	response, err := RequestToGetNonZeroAmountAssetList(token)
	if err != nil {
		return nil, err
	}
	for _, assetList := range *response {
		assetLists = append(assetLists, assetList)
	}
	return &assetLists, nil
}

type GetNonZeroAmountAssetListResponse struct {
	Success bool             `json:"success"`
	Error   string           `json:"error"`
	Code    ErrCode          `json:"code"`
	Data    *[]AssetListInfo `json:"data"`
}

func RequestToGetNonZeroAmountAssetList(token string) (*[]AssetListInfo, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_list/get"
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetNonZeroAmountAssetListResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

func AssetBalanceInfoToListBalanceInfo(assetBalanceInfo *AssetBalanceInfo) *ListBalanceInfo {
	return &ListBalanceInfo{
		GenesisPoint: assetBalanceInfo.GenesisPoint,
		Name:         assetBalanceInfo.Name,
		MetaHash:     assetBalanceInfo.MetaHash,
		AssetID:      assetBalanceInfo.AssetID,
		AssetType:    assetBalanceInfo.AssetType,
		OutputIndex:  assetBalanceInfo.OutputIndex,
		Version:      assetBalanceInfo.Version,
		Balance:      assetBalanceInfo.Balance,
	}
}

func AssetBalanceInfosToListBalanceInfos(assetBalanceInfos *[]AssetBalanceInfo) *[]ListBalanceInfo {
	var istBalanceInfos []ListBalanceInfo
	for _, assetBalanceInfo := range *assetBalanceInfos {
		istBalanceInfos = append(istBalanceInfos, *AssetBalanceInfoToListBalanceInfo(&assetBalanceInfo))
	}
	return &istBalanceInfos
}

func QueryAssetTransfersByAssetIdFromServer(token string, assetId string) string {
	return GetAssetTransferByAssetIdFromServer(token, assetId)
}

func QueryAssetTransfersAndGetResponse(assetId string) (*[]Transfer, error) {
	response, err := rpcclient.ListTransfers()
	if err != nil {
		return nil, err
	}
	var transfers []Transfer
	for _, t := range response.Transfers {
		if assetId != "" && assetId != hex.EncodeToString(t.Inputs[0].AssetId) {
			continue
		}
		newTransfer := Transfer{}
		newTransfer.GetData(t)
		transfers = append(transfers, newTransfer)
	}
	if len(transfers) == 0 {
		return nil, err
	}
	return &transfers, nil
}

func GetAllNftAssetIdOfListAssets() (map[string]bool, error) {
	listAssetResponse, err := ListAssetAndGetResponse()
	if err != nil {
		return nil, err
	}
	assetIdExists := make(map[string]bool)
	for _, asset := range listAssetResponse.Assets {
		if asset.AssetGenesis.AssetType == taprpc.AssetType_COLLECTIBLE {
			assetId := hex.EncodeToString(asset.AssetGenesis.AssetId)
			assetIdExists[assetId] = true
		}
	}
	return assetIdExists, nil
}

func QueryAssetTransfersAndGetResponseOfAllNft(assetIdExists map[string]bool) (*[]Transfer, error) {
	response, err := rpcclient.ListTransfers()
	if err != nil {
		return nil, err
	}
	var transfers []Transfer
	for _, t := range response.Transfers {
		assetId := hex.EncodeToString(t.Inputs[0].AssetId)
		if !assetIdExists[assetId] {
			continue
		}
		newTransfer := Transfer{}
		newTransfer.GetData(t)
		transfers = append(transfers, newTransfer)
	}
	if len(transfers) == 0 {
		return nil, err
	}
	return &transfers, nil
}

type AssetTransferSimplified struct {
	AssetID     string    `json:"asset_id"`
	Txid        string    `json:"txid"`
	TotalAmount int       `json:"totalAmount"`
	Time        int       `json:"time"`
	Detail      *Transfer `json:"detail"`
}

func ProcessAssetTransferByBitcoind(token string, allOutpoints []string, assetTransfers *[]Transfer) (*[]AssetTransferSimplified, error) {
	var assetTransferSimplified []AssetTransferSimplified
	response, err := PostCallBitcoindToQueryAddressByOutpoints(token, allOutpoints)
	if err != nil {
		return nil, err
	}
	addressMap := response.Data
	for _, assetTransfer := range *assetTransfers {
		for _, input := range assetTransfer.Inputs {
			(*input).Address = addressMap[input.AnchorPoint]
		}
		var totalAmount int
		for _, output := range assetTransfer.Outputs {
			(*output).Anchor.Address = addressMap[output.Anchor.Outpoint]
			totalAmount += int(output.Amount)
		}
		assetTransferSimplified = append(assetTransferSimplified, AssetTransferSimplified{
			Txid:        assetTransfer.Txid,
			AssetID:     assetTransfer.Inputs[0].AssetID,
			TotalAmount: totalAmount,
			Time:        int(assetTransfer.TransferTimestamp),
			Detail:      &assetTransfer,
		})
	}
	return &assetTransferSimplified, nil
}

func ProcessAssetTransfer(assetTransfers *[]Transfer) (*[]AssetTransferSimplified, error) {
	var assetTransferSimplified []AssetTransferSimplified
	for _, assetTransfer := range *assetTransfers {
		var totalAmount int
		for _, output := range assetTransfer.Outputs {
			totalAmount += int(output.Amount)
		}
		assetTransferSimplified = append(assetTransferSimplified, AssetTransferSimplified{
			Txid:        assetTransfer.Txid,
			AssetID:     assetTransfer.Inputs[0].AssetID,
			TotalAmount: totalAmount,
			Time:        int(assetTransfer.TransferTimestamp),
			Detail:      &assetTransfer,
		})
	}
	return &assetTransferSimplified, nil
}

func QueryAssetTransferSimplified(token string, assetId string) (*[]AssetTransferSimplified, error) {
	var assetTransferSimplified *[]AssetTransferSimplified
	assetTransfers, err := QueryAssetTransfersAndGetResponse(assetId)
	if err != nil {
		return nil, err
	}
	if assetTransfers == nil {
		return nil, nil
	}
	_ = token
	assetTransferSimplified, err = ProcessAssetTransfer(assetTransfers)
	return assetTransferSimplified, nil
}

func QueryAssetTransferSimplifiedOfAllNft(token string) (*[]AssetTransferSimplified, error) {
	var assetTransferSimplified *[]AssetTransferSimplified
	assetIdExists, err := GetAllNftAssetIdOfListAssets()
	if err != nil {
		return nil, err
	}
	assetTransfers, err := QueryAssetTransfersAndGetResponseOfAllNft(assetIdExists)
	if err != nil {
		return nil, err
	}
	if assetTransfers == nil {
		return nil, nil
	}
	_ = token
	assetTransferSimplified, err = ProcessAssetTransfer(assetTransfers)
	return assetTransferSimplified, nil
}

type GetAssetHolderNumberByAssetBalancesInfoResponse struct {
	Success bool    `json:"success"`
	Error   string  `json:"error"`
	Code    ErrCode `json:"code"`
	Data    int     `json:"data"`
}

func RequestToGetAssetHolderNumberByAssetBalancesInfo(token string, assetId string) (int, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_balance/get/holder/number/" + assetId
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return 0, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	var response GetAssetHolderNumberByAssetBalancesInfoResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}
	if response.Error != "" {
		return 0, errors.New(response.Error)
	}
	return response.Data, nil
}

func GetAssetHolderNumberByAssetBalancesInfo(token string, assetId string) (int, error) {
	holderNumber, err := RequestToGetAssetHolderNumberByAssetBalancesInfo(token, assetId)
	if err != nil {
		return 0, err
	}
	return holderNumber, nil
}

func GetAssetHolderNumber(token string, assetId string) string {
	holderNumber, err := GetAssetHolderNumberByAssetBalancesInfo(token, assetId)
	if err != nil {
		return MakeJsonErrorResult(GetAssetHolderNumberByAssetBalancesInfoErr, err.Error(), 0)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, holderNumber)
}

type AssetIdAndBalance struct {
	AssetId       string              `json:"asset_id"`
	AssetBalances *[]AssetBalanceInfo `json:"asset_balances"`
}

type GetAssetHolderBalanceByAssetBalancesInfoResponse struct {
	Success bool               `json:"success"`
	Error   string             `json:"error"`
	Code    ErrCode            `json:"code"`
	Data    *AssetIdAndBalance `json:"data"`
}

func RequestToGetAssetHolderBalanceByAssetBalancesInfo(token string, assetId string) (*AssetIdAndBalance, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_balance/get/holder/balance/all/" + assetId
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetAssetHolderBalanceByAssetBalancesInfoResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

type GetAssetHolderBalanceRecordsLengthResponse struct {
	Success bool    `json:"success"`
	Error   string  `json:"error"`
	Code    ErrCode `json:"code"`
	Data    int     `json:"data"`
}

func RequestToGetAssetHolderBalanceRecordsLengthByAssetBalancesInfo(token string, assetId string) (int, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_balance/get/holder/balance/records/" + assetId
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return 0, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	var response GetAssetHolderBalanceRecordsLengthResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}
	if response.Error != "" {
		return 0, errors.New(response.Error)
	}
	return response.Data, nil
}

func GetAssetHolderBalanceRecordsLengthNumber(token string, assetId string) (int, error) {
	return RequestToGetAssetHolderBalanceRecordsLengthByAssetBalancesInfo(token, assetId)
}

func GetAssetHolderBalanceByAssetBalancesInfo(token string, assetId string) (*AssetIdAndBalance, error) {
	holderBalance, err := RequestToGetAssetHolderBalanceByAssetBalancesInfo(token, assetId)
	if err != nil {
		return nil, err
	}
	return holderBalance, nil
}

type AssetBalanceInfoSimplified struct {
	Version  int    `json:"version"`
	Balance  int    `json:"balance"`
	DeviceId string `json:"device_id"`
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
}

type AssetIdAndBalanceSimplified struct {
	AssetId       string                        `json:"asset_id"`
	AssetBalances *[]AssetBalanceInfoSimplified `json:"asset_balances"`
}

func AssetIdAndBalanceToAssetIdAndBalanceSimplified(assetIdAndBalance *AssetIdAndBalance) *AssetIdAndBalanceSimplified {
	if assetIdAndBalance == nil {
		return nil
	}
	assetIdAndBalanceSimplified := &AssetIdAndBalanceSimplified{}
	assetIdAndBalanceSimplified.AssetId = assetIdAndBalance.AssetId
	var assetBalanceInfoSimplified []AssetBalanceInfoSimplified
	if assetIdAndBalance.AssetBalances == nil {
		assetIdAndBalanceSimplified.AssetBalances = &[]AssetBalanceInfoSimplified{}
		return assetIdAndBalanceSimplified
	}
	for _, assetBalanceInfo := range *(assetIdAndBalance.AssetBalances) {
		assetBalanceInfoSimplified = append(assetBalanceInfoSimplified, AssetBalanceInfoSimplified{
			Version:  assetBalanceInfo.Version,
			Balance:  assetBalanceInfo.Balance,
			DeviceId: assetBalanceInfo.DeviceId,
			UserId:   assetBalanceInfo.UserId,
			Username: assetBalanceInfo.Username,
		})
	}
	assetIdAndBalanceSimplified.AssetBalances = &assetBalanceInfoSimplified
	return assetIdAndBalanceSimplified
}

func GetAssetHolderBalance(token string, assetId string) string {
	holderBalance, err := GetAssetHolderBalanceByAssetBalancesInfo(token, assetId)
	if err != nil {
		return MakeJsonErrorResult(GetAssetHolderBalanceByAssetBalancesInfoErr, err.Error(), nil)
	}
	result := AssetIdAndBalanceToAssetIdAndBalanceSimplified(holderBalance)
	return MakeJsonErrorResult(SUCCESS, SuccessError, result)
}

type GetTimesByOutpointSliceResponse struct {
	Success bool           `json:"success"`
	Error   string         `json:"error"`
	Code    ErrCode        `json:"code"`
	Data    map[string]int `json:"data"`
}

func PostCallBitcoindToQueryTimeByOutpoints(token string, outpoints []string) (*GetTimesByOutpointSliceResponse, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	network := base.NetWork
	url := serverDomainOrSocket + "/bitcoind/" + network + "/time/outpoints"
	requestStr := OutpointSliceToRequestBodyRawString(outpoints)
	payload := strings.NewReader(requestStr)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetTimesByOutpointSliceResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func PostCallBitcoindToQueryTimeByOutpoints2(host, token string, outpoints []string) (*GetTimesByOutpointSliceResponse, error) {
	network := base.NetWork
	url := host + "/bitcoind/" + network + "/time/outpoints"
	requestStr := OutpointSliceToRequestBodyRawString(outpoints)
	payload := strings.NewReader(requestStr)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetTimesByOutpointSliceResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func GetAllOutpointsOfListUnspentUtxos(listUnspentUtxo *[]ListUnspentUtxo) []string {
	var ops []string
	for _, utxo := range *listUnspentUtxo {
		ops = append(ops, utxo.Outpoint)
	}
	return ops
}

func GetTimeForListUnspentUtxoByBitcoind(token string, listUnspentUtxo *[]ListUnspentUtxo) (*[]ListUnspentUtxo, error) {
	ops := GetAllOutpointsOfListUnspentUtxos(listUnspentUtxo)
	opMapTime, err := PostCallBitcoindToQueryTimeByOutpoints(token, ops)
	if err != nil {
		return nil, err
	}
	for i, utxo := range *listUnspentUtxo {
		(*listUnspentUtxo)[i].Time = opMapTime.Data[utxo.Outpoint]
	}
	return listUnspentUtxo, nil
}

type AssetHolderBalanceLimitAndOffsetRequest struct {
	AssetId string `json:"asset_id"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
}

func PostToGetAssetHolderBalanceLimitAndOffsetByAssetBalancesInfo(token string, assetId string, limit int, offset int) (*AssetIdAndBalance, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	assetIdLimitAndOffset := AssetHolderBalanceLimitAndOffsetRequest{
		AssetId: assetId,
		Limit:   limit,
		Offset:  offset,
	}
	url := serverDomainOrSocket + "/asset_balance/get/holder/balance/limit_offset"
	requestJsonBytes, err := json.Marshal(assetIdLimitAndOffset)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetAssetHolderBalanceByAssetBalancesInfoResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

type GetAssetAddrByEncodedResponse struct {
	Success bool       `json:"success"`
	Error   string     `json:"error"`
	Code    ErrCode    `json:"code"`
	Data    *AssetAddr `json:"data"`
}

func RequestToGetAssetAddrByEncoded(token string, encoded string) (*AssetAddr, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_addr/get/encoded/" + encoded
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetAssetAddrByEncodedResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

func GetAssetAddrByEncoded(token string, encoded string) (*AssetAddr, error) {
	return RequestToGetAssetAddrByEncoded(token, encoded)
}

func GetUsernameByEncoded(token string, encoded string) (string, error) {
	assetAddr, err := GetAssetAddrByEncoded(token, encoded)
	if err != nil {
		return "", err
	}
	return assetAddr.Username, nil
}

type AssetBurn struct {
	gorm.Model
	AssetId  string `json:"asset_id" gorm:"type:varchar(255)"`
	Amount   int    `json:"amount"`
	DeviceId string `json:"device_id" gorm:"type:varchar(255)"`
	UserId   int    `json:"user_id"`
	Username string `json:"username" gorm:"type:varchar(255)"`
	Status   int    `json:"status" gorm:"default:1"`
}

type AssetBurnSetRequest struct {
	AssetId  string `json:"asset_id"`
	Amount   int    `json:"amount"`
	DeviceId string `json:"device_id"`
}

func PostToSetAssetBurn(token string, assetBurnSetRequest *AssetBurnSetRequest) (*JsonResult, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_burn/set"
	requestJsonBytes, err := json.Marshal(assetBurnSetRequest)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response JsonResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return &response, nil
}

func UploadAssetBurn(token string, assetId string, amount int, deviceId string) error {
	assetBurnSetRequest := &AssetBurnSetRequest{
		AssetId:  assetId,
		Amount:   amount,
		DeviceId: deviceId,
	}
	_, err := PostToSetAssetBurn(token, assetBurnSetRequest)
	return err
}

type GetAssetBurnTotalAmountByAssetIdResponse struct {
	Success bool    `json:"success"`
	Error   string  `json:"error"`
	Code    ErrCode `json:"code"`
	Data    int     `json:"data"`
}

func RequestToGetAssetBurnTotalAmountByAssetId(token string, assetId string) (int, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_burn/get/asset_id/" + assetId
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return 0, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	var response GetAssetBurnTotalAmountByAssetIdResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}
	if response.Error != "" {
		return 0, errors.New(response.Error)
	}
	return response.Data, nil
}

func GetAssetBurnTotalAmountByAssetId(token string, assetId string) (int, error) {
	return RequestToGetAssetBurnTotalAmountByAssetId(token, assetId)
}

func GetAssetBurnTotalAmount(token string, assetId string) string {
	totalAmount, err := GetAssetBurnTotalAmountByAssetId(token, assetId)
	if err != nil {
		return MakeJsonErrorResult(GetAssetBurnTotalAmountByAssetIdErr, err.Error(), 0)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, totalAmount)
}

type FairLaunchInfoSimplified struct {
	ID                    int                    `json:"id"`
	Name                  string                 `json:"name"`
	ReserveTotal          int                    `json:"reserve_total"`
	CalculationExpression string                 `json:"calculation_expression"`
	AssetID               string                 `json:"asset_id"`
	State                 models.FairLaunchState `json:"state"`
}

type GetOwnFairLaunchInfoIssuedSimplifiedResponse struct {
	Success bool                        `json:"success"`
	Error   string                      `json:"error"`
	Code    ErrCode                     `json:"code"`
	Data    *[]FairLaunchInfoSimplified `json:"data"`
}

func RequestToGetOwnFairLaunchInfoIssuedSimplified(token string) (*[]FairLaunchInfoSimplified, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/v1/fair_launch/query/own_set/issued/simplified"
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetOwnFairLaunchInfoIssuedSimplifiedResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

func GetOwnFairLaunchInfoIssuedSimplified(token string) (*[]FairLaunchInfoSimplified, error) {
	return RequestToGetOwnFairLaunchInfoIssuedSimplified(token)
}

type MintReservedRequest struct {
	AssetID     string `json:"asset_id"`
	EncodedAddr string `json:"encoded_addr"`
}

type MintReservedResponse struct {
	Success bool    `json:"success"`
	Error   string  `json:"error"`
	Code    ErrCode `json:"code"`
	Data    struct {
		AnchorOutpoint string `json:"anchor_outpoint"`
	} `json:"data"`
}

func PostToFairLaunchMintReserved(token string, mintReservedRequest *MintReservedRequest) (string, error) {
	if mintReservedRequest == nil {
		return "", errors.New("invalid request")
	}
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/v1/fair_launch/mint_reserved"
	requestJsonBytes, err := json.Marshal(mintReservedRequest)
	if err != nil {
		return "", err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var response MintReservedResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}
	if response.Error != "" {
		return "", errors.New(response.Error)
	}
	return response.Data.AnchorOutpoint, nil
}

func FairLaunchMintReserved(token string, assetId string, encodedAddr string) (string, error) {
	mintReservedRequest := MintReservedRequest{
		AssetID:     assetId,
		EncodedAddr: encodedAddr,
	}
	return PostToFairLaunchMintReserved(token, &mintReservedRequest)
}

func GetOwnFairLaunchInfoIssuedSimplifiedAndExecuteMintReserved(token string, deviceId string) ([]string, error) {
	fairLaunchInfoIssuedSimplified, err := GetOwnFairLaunchInfoIssuedSimplified(token)
	if err != nil {
		return nil, err
	}
	if fairLaunchInfoIssuedSimplified == nil {
		return []string{}, nil
	}
	var outpoints []string
	for _, fairLaunchInfo := range *fairLaunchInfoIssuedSimplified {
		assetId := fairLaunchInfo.AssetID
		amount := fairLaunchInfo.ReserveTotal
		var encoded string
		encoded, err = NewAddrAndGetResponseEncoded(assetId, amount, token, deviceId)
		if err != nil {
			return nil, err
		}
		var op string
		op, err = FairLaunchMintReserved(token, fairLaunchInfo.AssetID, encoded)
		if err != nil {
			return nil, err
		}
		outpoints = append(outpoints, op)
	}
	return outpoints, nil
}

func AutoMintReserved(token string, deviceId string) string {
	result, err := GetOwnFairLaunchInfoIssuedSimplifiedAndExecuteMintReserved(token, deviceId)
	if err != nil {
		return MakeJsonErrorResult(GetOwnFairLaunchInfoIssuedSimplifiedAndExecuteMintReservedErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, result)
}

type AssetLocalMint struct {
	gorm.Model
	AssetVersion    string `json:"asset_version" gorm:"type:varchar(255)"`
	AssetType       string `json:"asset_type" gorm:"type:varchar(255)"`
	Name            string `json:"name" gorm:"type:varchar(255)"`
	AssetMetaData   string `json:"asset_meta_data"`
	AssetMetaType   string `json:"asset_meta_type" gorm:"type:varchar(255)"`
	AssetMetaHash   string `json:"asset_meta_hash" gorm:"type:varchar(255)"`
	Amount          int    `json:"amount"`
	NewGroupedAsset bool   `json:"new_grouped_asset"`
	GroupKey        string `json:"group_key" gorm:"type:varchar(255)"`
	GroupAnchor     string `json:"group_anchor" gorm:"type:varchar(255)"`
	GroupedAsset    bool   `json:"grouped_asset"`
	BatchKey        string `json:"batch_key" gorm:"type:varchar(255)"`
	BatchTxid       string `json:"batch_txid" gorm:"type:varchar(255)"`
	AssetId         string `json:"asset_id" gorm:"type:varchar(255)"`
	DeviceId        string `json:"device_id" gorm:"type:varchar(255)"`
	UserId          int    `json:"user_id"`
	Username        string `json:"username" gorm:"type:varchar(255)"`
	Status          int    `json:"status" gorm:"default:1"`
}

type AssetLocalMintSetRequest struct {
	AssetVersion    string `json:"asset_version"`
	AssetType       string `json:"asset_type"`
	Name            string `json:"name"`
	AssetMetaData   string `json:"asset_meta_data"`
	AssetMetaType   string `json:"asset_meta_type"`
	AssetMetaHash   string `json:"asset_meta_hash"`
	Amount          int    `json:"amount"`
	NewGroupedAsset bool   `json:"new_grouped_asset"`
	GroupKey        string `json:"group_key"`
	GroupAnchor     string `json:"group_anchor"`
	GroupedAsset    bool   `json:"grouped_asset"`
	BatchKey        string `json:"batch_key"`
	BatchTxid       string `json:"batch_txid"`
	AssetId         string `json:"asset_id"`
	DeviceId        string `json:"device_id"`
}

func BatchTxidAnchorToAssetId(batchTxidAnchor string) (string, error) {
	assets, _ := listAssets(true, true, false)
	for _, asset := range assets.Assets {
		txid, _ := outpointToTransactionAndIndex(asset.GetChainAnchor().GetAnchorOutpoint())
		if batchTxidAnchor == txid {
			return hex.EncodeToString(asset.GetAssetGenesis().AssetId), nil
		}
	}
	err := errors.New("no asset found for batch txid")
	return "", err
}

func BatchTxidAndAssetMintInfoToAssetId(batchTxid string, pendingAsset *mintrpc.PendingAsset) (string, error) {
	assets, _ := listAssets(true, true, false)
	for _, asset := range assets.Assets {
		txid, _ := outpointToTransactionAndIndex(asset.GetChainAnchor().GetAnchorOutpoint())
		var isMetaHashEqual bool
		if pendingAsset.AssetMeta == nil || asset.AssetGenesis == nil {
			isMetaHashEqual = true
		} else if pendingAsset.AssetMeta.MetaHash == nil || asset.AssetGenesis.MetaHash == nil {
			isMetaHashEqual = true
		} else {
			isMetaHashEqual = hex.EncodeToString(pendingAsset.AssetMeta.MetaHash) == hex.EncodeToString(asset.AssetGenesis.MetaHash)
		}
		if batchTxid == txid &&
			pendingAsset.Name == asset.AssetGenesis.Name &&
			pendingAsset.Amount == asset.Amount &&
			isMetaHashEqual &&
			pendingAsset.AssetType == asset.AssetGenesis.AssetType {
			return hex.EncodeToString(asset.GetAssetGenesis().AssetId), nil
		}
	}
	err := errors.New("no asset found for batch txid")
	return "", err
}

func BatchAssetToAssetLocalMintSetRequest(batchKey string, batchTxid string, deviceId string, asset *mintrpc.PendingAsset) *AssetLocalMintSetRequest {
	if asset == nil {
		return nil
	}
	groupKey := hex.EncodeToString(asset.GroupKey)
	groupedAsset := asset.NewGroupedAsset || groupKey != ""
	assetId, err := BatchTxidAndAssetMintInfoToAssetId(batchTxid, asset)
	if err != nil {
	}
	return &AssetLocalMintSetRequest{
		AssetVersion:    asset.AssetVersion.String(),
		AssetType:       asset.AssetType.String(),
		Name:            asset.Name,
		AssetMetaData:   hex.EncodeToString(asset.AssetMeta.Data),
		AssetMetaType:   asset.AssetMeta.Type.String(),
		AssetMetaHash:   hex.EncodeToString(asset.AssetMeta.MetaHash),
		Amount:          int(asset.Amount),
		NewGroupedAsset: asset.NewGroupedAsset,
		GroupKey:        groupKey,
		GroupAnchor:     asset.GroupAnchor,
		GroupedAsset:    groupedAsset,
		BatchKey:        batchKey,
		BatchTxid:       batchTxid,
		AssetId:         assetId,
		DeviceId:        deviceId,
	}
}

func FinalizeBatchResponseToAssetLocalMintSetRequests(deviceId string, finalizeBatchResponse *mintrpc.FinalizeBatchResponse) *[]AssetLocalMintSetRequest {
	batch := finalizeBatchResponse.GetBatch()
	if batch == nil {
		return nil
	}
	var assetLocalMintSetRequests []AssetLocalMintSetRequest
	batchKey := hex.EncodeToString(batch.BatchKey)
	batchTxid := batch.BatchTxid
	for _, asset := range (*batch).Assets {
		assetLocalMintSetRequests = append(assetLocalMintSetRequests, *BatchAssetToAssetLocalMintSetRequest(batchKey, batchTxid, deviceId, asset))
	}
	return &assetLocalMintSetRequests
}

func PostToSetAssetLocalMints(token string, assetLocalMintSetRequests *[]AssetLocalMintSetRequest) error {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_local_mint/set/slice"
	requestJsonBytes, err := json.Marshal(assetLocalMintSetRequests)
	if err != nil {
		return err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var response JsonResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	if response.Error != "" {
		return errors.New(response.Error)
	}
	return nil
}

func UploadAssetLocalMints(token string, deviceId string, finalizeBatchResponse *mintrpc.FinalizeBatchResponse) error {
	assetLocalMintSetRequests := FinalizeBatchResponseToAssetLocalMintSetRequests(deviceId, finalizeBatchResponse)
	return PostToSetAssetLocalMints(token, assetLocalMintSetRequests)
}

type PendingBatch struct {
	BatchKey  string              `json:"batch_key"`
	BatchTxid string              `json:"batch_txid"`
	State     string              `json:"state"`
	Assets    []PendingBatchAsset `json:"assets"`
}

type PendingBatchAsset struct {
	AssetVersion      string `json:"asset_version"`
	AssetType         string `json:"asset_type"`
	Name              string `json:"name"`
	AssetMetaData     string `json:"asset_meta_data"`
	AssetMetaType     string `json:"asset_meta_type"`
	AssetMetaMetaHash string `json:"asset_meta_meta_hash"`
	Amount            int    `json:"amount"`
	NewGroupedAsset   bool   `json:"new_grouped_asset"`
	GroupKey          string `json:"group_key"`
	GroupAnchor       string `json:"group_anchor"`
}

func BatchPendingAssetToPendingBatchAsset(pendingAsset *mintrpc.PendingAsset) *PendingBatchAsset {
	if pendingAsset == nil {
		return nil
	}
	return &PendingBatchAsset{
		AssetVersion:      pendingAsset.AssetVersion.String(),
		AssetType:         pendingAsset.AssetType.String(),
		Name:              pendingAsset.Name,
		AssetMetaData:     "",
		AssetMetaType:     pendingAsset.AssetMeta.Type.String(),
		AssetMetaMetaHash: hex.EncodeToString(pendingAsset.AssetMeta.MetaHash),
		Amount:            int(pendingAsset.Amount),
		NewGroupedAsset:   pendingAsset.NewGroupedAsset,
		GroupKey:          hex.EncodeToString(pendingAsset.GroupKey),
		GroupAnchor:       pendingAsset.GroupAnchor,
	}
}

func BatchPendingAssetSliceToPendingBatchAssetSlice(pendingAssets []*mintrpc.PendingAsset) []PendingBatchAsset {
	var finalizedBatchAssets []PendingBatchAsset
	if len(pendingAssets) == 0 {
		return finalizedBatchAssets
	}
	for _, pendingAsset := range pendingAssets {
		finalizedBatchAsset := BatchPendingAssetToPendingBatchAsset(pendingAsset)
		if finalizedBatchAsset == nil {
			continue
		}
		finalizedBatchAssets = append(finalizedBatchAssets, *finalizedBatchAsset)
	}
	return finalizedBatchAssets
}

func MintAssetResponseToPendingBatch(mintAssetResponse *mintrpc.MintAssetResponse) *PendingBatch {
	if mintAssetResponse == nil {
		return nil
	}
	return &PendingBatch{
		BatchKey:  hex.EncodeToString(mintAssetResponse.PendingBatch.BatchKey),
		BatchTxid: mintAssetResponse.PendingBatch.BatchTxid,
		State:     mintAssetResponse.PendingBatch.State.String(),
		Assets:    BatchPendingAssetSliceToPendingBatchAssetSlice(mintAssetResponse.PendingBatch.Assets),
	}
}

func FinalizeBatchResponseToPendingBatch(finalizeBatchResponse *mintrpc.FinalizeBatchResponse) *PendingBatch {
	if finalizeBatchResponse == nil {
		return nil
	}
	return &PendingBatch{
		BatchKey:  hex.EncodeToString(finalizeBatchResponse.Batch.BatchKey),
		BatchTxid: finalizeBatchResponse.Batch.BatchTxid,
		State:     finalizeBatchResponse.Batch.State.String(),
		Assets:    BatchPendingAssetSliceToPendingBatchAssetSlice(finalizeBatchResponse.Batch.Assets),
	}
}

type AssetRecommend struct {
	gorm.Model
	AssetId           string `json:"asset_id" gorm:"type:varchar(255)"`
	AssetFromAddr     string `json:"asset_from_addr" gorm:"type:varchar(255)"`
	RecommendUserId   int    `json:"recommend_user_id"`
	RecommendUsername string `json:"recommend_username" gorm:"type:varchar(255)"`
	RecommendTime     int    `json:"recommend_time"`
	DeviceId          string `json:"device_id" gorm:"type:varchar(255)"`
	UserId            int    `json:"user_id"`
	Username          string `json:"username" gorm:"type:varchar(255)"`
	Status            int    `json:"status" gorm:"default:1"`
}

type AssetRecommendSetRequest struct {
	AssetId           string `json:"asset_id"`
	AssetFromAddr     string `json:"asset_from_addr"`
	RecommendUserId   int    `json:"recommend_user_id"`
	RecommendUsername string `json:"recommend_username"`
	RecommendTime     int    `json:"recommend_time"`
	DeviceId          string `json:"device_id"`
}

type GetAssetRecommendsByUserIdAndAssetId struct {
	Success bool            `json:"success"`
	Error   string          `json:"error"`
	Code    ErrCode         `json:"code"`
	Data    *AssetRecommend `json:"data"`
}

func RequestToGetUserAssetRecommendByAssetId(token string, assetId string) (*AssetRecommend, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_recommend/get/user/asset_id/" + assetId
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetAssetRecommendsByUserIdAndAssetId
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

func GetUserAssetRecommendByAssetId(token string, assetId string) (*AssetRecommend, error) {
	return RequestToGetUserAssetRecommendByAssetId(token, assetId)
}

type UserIdAndAssetId struct {
	UserId  int    `json:"user_id"`
	AssetId string `json:"asset_id"`
}

func PostToGetAssetRecommendByUserIdAndAssetId(token string, userIdAndAssetId UserIdAndAssetId) (*AssetRecommend, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_recommend/get/user_id_and_asset_id"
	requestJsonBytes, err := json.Marshal(userIdAndAssetId)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetAssetRecommendsByUserIdAndAssetId
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		if response.Error == gorm.ErrRecordNotFound.Error() {
			err = gorm.ErrRecordNotFound
		} else {
			err = errors.New(response.Error)
		}
		return nil, err
	}
	return response.Data, nil
}

func GetAssetRecommendByUserIdAndAssetId(token string, userIdAndAssetId UserIdAndAssetId) (*AssetRecommend, error) {
	return PostToGetAssetRecommendByUserIdAndAssetId(token, userIdAndAssetId)
}

func PostToSetAssetRecommendByAssetId(token string, assetRecommendSetRequest *AssetRecommendSetRequest) (*JsonResult, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_recommend/set"
	requestJsonBytes, err := json.Marshal(assetRecommendSetRequest)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response JsonResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return &response, nil
}

func SetAssetRecommendByAssetId(token string, assetId string, assetFromAddr string, recommendUserId int, recommendUsername string, recommendTime int, deviceId string) error {
	assetRecommendSetRequest := AssetRecommendSetRequest{
		AssetId:           assetId,
		AssetFromAddr:     assetFromAddr,
		RecommendUserId:   recommendUserId,
		RecommendUsername: recommendUsername,
		RecommendTime:     recommendTime,
		DeviceId:          deviceId,
	}
	_, err := PostToSetAssetRecommendByAssetId(token, &assetRecommendSetRequest)
	return err
}

func CheckIfAssetIsIssuedLocally(assetId string) (bool, error) {
	keys, err := assetLeafKeys(assetId, universerpc.ProofType_PROOF_TYPE_ISSUANCE)
	if err != nil || len(keys.AssetKeys) == 0 {
		errorAppendInfo := ErrorAppendInfo(err)
		return false, errorAppendInfo("failed to get asset info")
	}
	result := struct {
		IsLocal   bool   `json:"is_local"`
		AssetId   string `json:"asset_id"`
		BatchTxid string `json:"batch_txid"`
		Amount    int64  `json:"amount"`
		Timestamp int64  `json:"timestamp"`
		ScriptKey string `json:"script_key"`
	}{
		IsLocal: false,
		AssetId: assetId,
	}
	Outpoint := keys.AssetKeys[0].Outpoint
	if o, ok := Outpoint.(*universerpc.AssetKey_OpStr); ok {
		opStr := strings.Split(o.OpStr, ":")
		var listBatch *mintrpc.ListBatchResponse
		listBatch, err = ListBatchesAndGetResponse()
		if err != nil {
			errorAppendInfo := ErrorAppendInfo(err)
			return false, errorAppendInfo("failed to get mint info")
		}
		for _, b := range listBatch.Batches {
			if b.Batch.BatchTxid == opStr[0] {
				var leaves *universerpc.AssetLeafResponse
				leaves, err = assetLeaves(false, assetId, universerpc.ProofType_PROOF_TYPE_ISSUANCE)
				if err != nil {
					errorAppendInfo := ErrorAppendInfo(err)
					return false, errorAppendInfo("failed to get mint info")
				}
				result.Amount = int64(leaves.Leaves[0].Asset.Amount)
				var transactions *lnrpc.TransactionDetails
				transactions, err = GetTransactionsAndGetResponse()
				if err != nil {
					errorAppendInfo := ErrorAppendInfo(err)
					return false, errorAppendInfo("failed to get mint info")
				}
				for _, tx := range transactions.Transactions {
					if tx.TxHash == opStr[0] {
						result.Timestamp = tx.TimeStamp
						break
					}
				}
				result.IsLocal = true
				result.BatchTxid = o.OpStr
				if s, _ok := keys.AssetKeys[0].ScriptKey.(*universerpc.AssetKey_ScriptKeyBytes); _ok {
					result.ScriptKey = "02" + hex.EncodeToString(s.ScriptKeyBytes)
				}
				break
			}
		}
		return true, nil
	}
	errorAppendInfo := ErrorAppendInfo(err)
	return false, errorAppendInfo("failed to get mint info")
}

func IsAssetLocalIssuance(assetId string) bool {
	isLocal, err := CheckIfAssetIsIssuedLocally(assetId)
	if err != nil {
		LogError("", err)
		return false
	}
	return isLocal
}

type AssetLocalMintHistory struct {
	gorm.Model
	AssetVersion    string `json:"asset_version" gorm:"type:varchar(255)"`
	AssetType       string `json:"asset_type" gorm:"type:varchar(255)"`
	Name            string `json:"name" gorm:"type:varchar(255)"`
	AssetMetaData   string `json:"asset_meta_data"`
	AssetMetaType   string `json:"asset_meta_type" gorm:"type:varchar(255)"`
	AssetMetaHash   string `json:"asset_meta_hash" gorm:"type:varchar(255)"`
	Amount          int    `json:"amount"`
	NewGroupedAsset bool   `json:"new_grouped_asset"`
	GroupKey        string `json:"group_key" gorm:"type:varchar(255)"`
	GroupAnchor     string `json:"group_anchor" gorm:"type:varchar(255)"`
	GroupedAsset    bool   `json:"grouped_asset"`
	BatchKey        string `json:"batch_key" gorm:"type:varchar(255)"`
	BatchTxid       string `json:"batch_txid" gorm:"type:varchar(255)"`
	AssetId         string `json:"asset_id" gorm:"type:varchar(255)"`
	DeviceId        string `json:"device_id" gorm:"type:varchar(255)"`
	UserId          int    `json:"user_id"`
	Username        string `json:"username" gorm:"type:varchar(255)"`
	Status          int    `json:"status" gorm:"default:1"`
}

type AssetLocalMintHistorySetRequest struct {
	AssetVersion    string `json:"asset_version"`
	AssetType       string `json:"asset_type"`
	Name            string `json:"name"`
	AssetMetaData   string `json:"asset_meta_data"`
	AssetMetaType   string `json:"asset_meta_type"`
	AssetMetaHash   string `json:"asset_meta_hash"`
	Amount          int    `json:"amount"`
	NewGroupedAsset bool   `json:"new_grouped_asset"`
	GroupKey        string `json:"group_key"`
	GroupAnchor     string `json:"group_anchor"`
	GroupedAsset    bool   `json:"grouped_asset"`
	BatchKey        string `json:"batch_key"`
	BatchTxid       string `json:"batch_txid"`
	AssetId         string `json:"asset_id"`
	DeviceId        string `json:"device_id"`
}

func BatchAssetToAssetLocalMintHistorySetRequest(batchKey string, batchTxid string, deviceId string, asset *mintrpc.PendingAsset) *AssetLocalMintHistorySetRequest {
	if asset == nil {
		return nil
	}
	groupKey := hex.EncodeToString(asset.GroupKey)
	groupedAsset := asset.NewGroupedAsset || groupKey != ""
	assetId, err := BatchTxidAndAssetMintInfoToAssetId(batchTxid, asset)
	if err != nil {
	}
	var assetMetaData string
	var assetMetaHash string
	var assetMetaType string
	if asset.AssetMeta != nil {
		if asset.AssetMeta.Data != nil {
			assetMetaData = hex.EncodeToString(asset.AssetMeta.Data)
		}
		if asset.AssetMeta.MetaHash != nil {
			assetMetaHash = hex.EncodeToString(asset.AssetMeta.MetaHash)
		}
		assetMetaType = asset.AssetMeta.Type.String()
	}
	return &AssetLocalMintHistorySetRequest{
		AssetVersion:    asset.AssetVersion.String(),
		AssetType:       asset.AssetType.String(),
		Name:            asset.Name,
		AssetMetaData:   assetMetaData,
		AssetMetaType:   assetMetaType,
		AssetMetaHash:   assetMetaHash,
		Amount:          int(asset.Amount),
		NewGroupedAsset: asset.NewGroupedAsset,
		GroupKey:        groupKey,
		GroupAnchor:     asset.GroupAnchor,
		GroupedAsset:    groupedAsset,
		BatchKey:        batchKey,
		BatchTxid:       batchTxid,
		AssetId:         assetId,
		DeviceId:        deviceId,
	}
}

func MintingBatchToAssetLocalMintHistorySetRequests(deviceId string, mintingBatch *mintrpc.MintingBatch) *[]AssetLocalMintHistorySetRequest {
	if mintingBatch == nil {
		return nil
	}
	var assetLocalMintHistorySetRequests []AssetLocalMintHistorySetRequest
	batchKey := hex.EncodeToString(mintingBatch.BatchKey)
	batchTxid := mintingBatch.BatchTxid
	for _, asset := range (*mintingBatch).Assets {
		request := BatchAssetToAssetLocalMintHistorySetRequest(batchKey, batchTxid, deviceId, asset)
		if request.AssetId != "" {
			assetLocalMintHistorySetRequests = append(assetLocalMintHistorySetRequests, *request)
		}
	}
	return &assetLocalMintHistorySetRequests
}

func ListBatchResponseToAssetLocalMintHistorySetRequests(deviceId string, listBatchResponse *mintrpc.ListBatchResponse) *[]AssetLocalMintHistorySetRequest {
	if listBatchResponse == nil {
		return nil
	}
	var assetLocalMintHistorySetRequests []AssetLocalMintHistorySetRequest
	for _, b := range listBatchResponse.Batches {
		requests := MintingBatchToAssetLocalMintHistorySetRequests(deviceId, b.Batch)
		if requests != nil {
			assetLocalMintHistorySetRequests = append(assetLocalMintHistorySetRequests, *requests...)
		}
	}
	return &assetLocalMintHistorySetRequests
}

func PostToSetAssetLocalMintHistories(token string, assetLocalMintHistorySetRequests *[]AssetLocalMintHistorySetRequest) (*JsonResult, error) {
	if assetLocalMintHistorySetRequests == nil {
		return &JsonResult{
			Success: true,
		}, nil
	}
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_local_mint_history/set"
	requestJsonBytes, err := json.Marshal(assetLocalMintHistorySetRequests)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response JsonResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return &response, nil
}

func ListBatchesAndPostToSetAssetLocalMintHistories(token string, deviceId string) error {
	listBatchResponse, err := ListBatchesAndGetResponse()
	if err != nil {
		return err
	}
	assetLocalMintHistorySetRequests := ListBatchResponseToAssetLocalMintHistorySetRequests(deviceId, listBatchResponse)
	_, err = PostToSetAssetLocalMintHistories(token, assetLocalMintHistorySetRequests)
	return err
}

func UploadAssetLocalMintHistory(token string, deviceId string) string {
	err := ListBatchesAndPostToSetAssetLocalMintHistories(token, deviceId)
	if err != nil {
		return MakeJsonErrorResult(ListBatchesAndPostToSetAssetLocalMintHistoriesErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, nil)
}

type GetAssetManagedUtxoIdsResponse struct {
	Success bool                `json:"success"`
	Error   string              `json:"error"`
	Code    ErrCode             `json:"code"`
	Data    *[]AssetManagedUtxo `json:"data"`
}

func RequestToGetAssetManagedUtxos(token string) (*[]AssetManagedUtxo, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_managed_utxo/get/user"
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetAssetManagedUtxoIdsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

func RequestToGetAssetManagedUtxos2(host, token string) (*[]AssetManagedUtxo, error) {
	url := host + "/asset_managed_utxo/get/user"
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetAssetManagedUtxoIdsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

func ManagedUtxosToAssetIds(managedUtxos *[]ManagedUtxo) *[]string {
	if managedUtxos == nil {
		return nil
	}
	var assetIds []string
	for _, managedUtxo := range *managedUtxos {
		for _, asset := range managedUtxo.ManagedUtxosAssets {
			assetIds = append(assetIds, asset.AssetGenesis.AssetID)
		}
	}
	return &assetIds
}

func AssetIdsToAssetIdMapIsExist(assetIds *[]string) *map[string]bool {
	assetIdMapIsExist := make(map[string]bool)
	if assetIds == nil {
		return &assetIdMapIsExist
	}
	for _, assetId := range *assetIds {
		assetIdMapIsExist[assetId] = true
	}
	return &assetIdMapIsExist
}

func PostToRemoveAssetManagedUtxos(token string, managedUtxos *[]int) (*JsonResult, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_managed_utxo/remove"
	requestJsonBytes, err := json.Marshal(managedUtxos)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response JsonResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return &response, nil
}

func RemoveNotLocalAssetManagedUtxos(token string, managedUtxos *[]ManagedUtxo) error {
	assetManagedUtxos, err := RequestToGetAssetManagedUtxos(token)
	if err != nil {
		errorAppendInfo := ErrorAppendInfo(err)
		return errorAppendInfo("Request To Get Asset Managed Utxos")
	}
	localAssetIds := ManagedUtxosToAssetIds(managedUtxos)
	localAssetIdMapIsExist := AssetIdsToAssetIdMapIsExist(localAssetIds)
	var idsRemove []int
	for _, assetManagedUtxo := range *assetManagedUtxos {
		if !(*localAssetIdMapIsExist)[assetManagedUtxo.AssetGenesisAssetID] {
			idsRemove = append(idsRemove, int(assetManagedUtxo.ID))
		}
	}
	if len(idsRemove) != 0 {
		_, err = PostToRemoveAssetManagedUtxos(token, &idsRemove)
		if err != nil {
			errorAppendInfo := ErrorAppendInfo(err)
			return errorAppendInfo("Post To Remove Asset Managed Utxos")
		}
	}
	return nil
}

func RemoveNotLocalAssetManagedUtxos2(host, token string, managedUtxos *[]ManagedUtxo) error {
	assetManagedUtxos, err := RequestToGetAssetManagedUtxos2(host, token)
	if err != nil {
		errorAppendInfo := ErrorAppendInfo(err)
		return errorAppendInfo("Request To Get Asset Managed Utxos")
	}
	localAssetIds := ManagedUtxosToAssetIds(managedUtxos)
	localAssetIdMapIsExist := AssetIdsToAssetIdMapIsExist(localAssetIds)
	var idsRemove []int
	for _, assetManagedUtxo := range *assetManagedUtxos {
		if !(*localAssetIdMapIsExist)[assetManagedUtxo.AssetGenesisAssetID] {
			idsRemove = append(idsRemove, int(assetManagedUtxo.ID))
		}
	}
	if len(idsRemove) != 0 {
		_, err = PostToRemoveAssetManagedUtxos(token, &idsRemove)
		if err != nil {
			errorAppendInfo := ErrorAppendInfo(err)
			return errorAppendInfo("Post To Remove Asset Managed Utxos")
		}
	}
	return nil
}

func ListUtxosAndGetProcessedManagedUtxos(token string) (*[]ManagedUtxo, error) {
	response, err := ListUtxosAndGetResponse(true)
	if err != nil {
		return nil, err
	}
	managedUtxos := ListUtxosResponseToManagedUtxos(response)
	err = RemoveNotLocalAssetManagedUtxos(token, managedUtxos)
	if err != nil {
		errorAppendInfo := ErrorAppendInfo(err)
		return nil, errorAppendInfo("Remove Not Local Asset Managed Utxos")
	}
	managedUtxos, err = GetTimeForManagedUtxoByBitcoind(token, managedUtxos)
	if err != nil {
		return nil, err
	}
	return managedUtxos, nil
}

type AssetManagedUtxo struct {
	gorm.Model
	Op                          string `json:"op" gorm:"type:varchar(255)"`
	OutPoint                    string `json:"out_point" gorm:"type:varchar(255)"`
	Time                        int    `json:"time"`
	AmtSat                      int    `json:"amt_sat"`
	InternalKey                 string `json:"internal_key" gorm:"type:varchar(255)"`
	TaprootAssetRoot            string `json:"taproot_asset_root" gorm:"type:varchar(255)"`
	MerkleRoot                  string `json:"merkle_root" gorm:"type:varchar(255)"`
	Version                     string `json:"version" gorm:"type:varchar(255)"`
	AssetGenesisPoint           string `json:"asset_genesis_point" gorm:"type:varchar(255)"`
	AssetGenesisName            string `json:"asset_genesis_name" gorm:"type:varchar(255)"`
	AssetGenesisMetaHash        string `json:"asset_genesis_meta_hash" gorm:"type:varchar(255)"`
	AssetGenesisAssetID         string `json:"asset_genesis_asset_id" gorm:"type:varchar(255)"`
	AssetGenesisAssetType       string `json:"asset_genesis_asset_type" gorm:"type:varchar(255)"`
	AssetGenesisOutputIndex     int    `json:"asset_genesis_output_index"`
	AssetGenesisVersion         int    `json:"asset_genesis_version"`
	Amount                      int    `json:"amount"`
	LockTime                    int    `json:"lock_time"`
	RelativeLockTime            int    `json:"relative_lock_time"`
	ScriptVersion               int    `json:"script_version"`
	ScriptKey                   string `json:"script_key" gorm:"type:varchar(255)"`
	ScriptKeyIsLocal            bool   `json:"script_key_is_local"`
	AssetGroupRawGroupKey       string `json:"asset_group_raw_group_key" gorm:"type:varchar(255)"`
	AssetGroupTweakedGroupKey   string `json:"asset_group_tweaked_group_key" gorm:"type:varchar(255)"`
	AssetGroupAssetWitness      string `json:"asset_group_asset_witness"`
	ChainAnchorTx               string `json:"chain_anchor_tx"`
	ChainAnchorBlockHash        string `json:"chain_anchor_block_hash" gorm:"type:varchar(255)"`
	ChainAnchorOutpoint         string `json:"chain_anchor_outpoint" gorm:"type:varchar(255)"`
	ChainAnchorInternalKey      string `json:"chain_anchor_internal_key" gorm:"type:varchar(255)"`
	ChainAnchorMerkleRoot       string `json:"chain_anchor_merkle_root" gorm:"type:varchar(255)"`
	ChainAnchorTapscriptSibling string `json:"chain_anchor_tapscript_sibling"`
	ChainAnchorBlockHeight      int    `json:"chain_anchor_block_height"`
	IsSpent                     bool   `json:"is_spent"`
	LeaseOwner                  string `json:"lease_owner" gorm:"type:varchar(255)"`
	LeaseExpiry                 int    `json:"lease_expiry"`
	IsBurn                      bool   `json:"is_burn"`
	DeviceId                    string `json:"device_id" gorm:"type:varchar(255)"`
	UserId                      int    `json:"user_id"`
	Username                    string `json:"username" gorm:"type:varchar(255)"`
	Status                      int    `json:"status" gorm:"default:1"`
}

type AssetManagedUtxoSetRequest struct {
	Op                          string `json:"op"`
	OutPoint                    string `json:"out_point"`
	Time                        int    `json:"time"`
	AmtSat                      int    `json:"amt_sat"`
	InternalKey                 string `json:"internal_key"`
	TaprootAssetRoot            string `json:"taproot_asset_root"`
	MerkleRoot                  string `json:"merkle_root"`
	Version                     string `json:"version"`
	AssetGenesisPoint           string `json:"asset_genesis_point"`
	AssetGenesisName            string `json:"asset_genesis_name"`
	AssetGenesisMetaHash        string `json:"asset_genesis_meta_hash"`
	AssetGenesisAssetID         string `json:"asset_genesis_asset_id"`
	AssetGenesisAssetType       string `json:"asset_genesis_asset_type"`
	AssetGenesisOutputIndex     int    `json:"asset_genesis_output_index"`
	AssetGenesisVersion         int    `json:"asset_genesis_version"`
	Amount                      int    `json:"amount"`
	LockTime                    int    `json:"lock_time"`
	RelativeLockTime            int    `json:"relative_lock_time"`
	ScriptVersion               int    `json:"script_version"`
	ScriptKey                   string `json:"script_key"`
	ScriptKeyIsLocal            bool   `json:"script_key_is_local"`
	AssetGroupRawGroupKey       string `json:"asset_group_raw_group_key"`
	AssetGroupTweakedGroupKey   string `json:"asset_group_tweaked_group_key"`
	AssetGroupAssetWitness      string `json:"asset_group_asset_witness"`
	ChainAnchorTx               string `json:"chain_anchor_tx"`
	ChainAnchorBlockHash        string `json:"chain_anchor_block_hash"`
	ChainAnchorOutpoint         string `json:"chain_anchor_outpoint"`
	ChainAnchorInternalKey      string `json:"chain_anchor_internal_key"`
	ChainAnchorMerkleRoot       string `json:"chain_anchor_merkle_root"`
	ChainAnchorTapscriptSibling string `json:"chain_anchor_tapscript_sibling"`
	ChainAnchorBlockHeight      int    `json:"chain_anchor_block_height"`
	IsSpent                     bool   `json:"is_spent"`
	LeaseOwner                  string `json:"lease_owner"`
	LeaseExpiry                 int    `json:"lease_expiry"`
	IsBurn                      bool   `json:"is_burn"`
	DeviceId                    string `json:"device_id" gorm:"type:varchar(255)"`
}

func ManagedUtxoAssetToAssetManagedUtxoSetRequest(deviceId string, op string, outPoint string, time int, amtSat int, internalKey string, taprootAssetRoot string, merkleRoot string, managedUtxoAsset ManagedUtxoAsset) AssetManagedUtxoSetRequest {
	return AssetManagedUtxoSetRequest{
		Op:                          op,
		OutPoint:                    outPoint,
		Time:                        time,
		AmtSat:                      amtSat,
		InternalKey:                 internalKey,
		TaprootAssetRoot:            taprootAssetRoot,
		MerkleRoot:                  merkleRoot,
		Version:                     managedUtxoAsset.Version,
		AssetGenesisPoint:           managedUtxoAsset.AssetGenesis.GenesisPoint,
		AssetGenesisName:            managedUtxoAsset.AssetGenesis.Name,
		AssetGenesisMetaHash:        managedUtxoAsset.AssetGenesis.MetaHash,
		AssetGenesisAssetID:         managedUtxoAsset.AssetGenesis.AssetID,
		AssetGenesisAssetType:       managedUtxoAsset.AssetGenesis.AssetType,
		AssetGenesisOutputIndex:     managedUtxoAsset.AssetGenesis.OutputIndex,
		AssetGenesisVersion:         managedUtxoAsset.AssetGenesis.Version,
		Amount:                      managedUtxoAsset.Amount,
		LockTime:                    managedUtxoAsset.LockTime,
		RelativeLockTime:            managedUtxoAsset.RelativeLockTime,
		ScriptVersion:               managedUtxoAsset.ScriptVersion,
		ScriptKey:                   managedUtxoAsset.ScriptKey,
		ScriptKeyIsLocal:            managedUtxoAsset.ScriptKeyIsLocal,
		AssetGroupRawGroupKey:       managedUtxoAsset.AssetGroup.RawGroupKey,
		AssetGroupTweakedGroupKey:   managedUtxoAsset.AssetGroup.TweakedGroupKey,
		AssetGroupAssetWitness:      managedUtxoAsset.AssetGroup.AssetWitness,
		ChainAnchorTx:               managedUtxoAsset.ChainAnchor.AnchorTx,
		ChainAnchorBlockHash:        managedUtxoAsset.ChainAnchor.AnchorBlockHash,
		ChainAnchorOutpoint:         managedUtxoAsset.ChainAnchor.AnchorOutpoint,
		ChainAnchorInternalKey:      managedUtxoAsset.ChainAnchor.InternalKey,
		ChainAnchorMerkleRoot:       managedUtxoAsset.ChainAnchor.MerkleRoot,
		ChainAnchorTapscriptSibling: managedUtxoAsset.ChainAnchor.TapscriptSibling,
		ChainAnchorBlockHeight:      managedUtxoAsset.ChainAnchor.BlockHeight,
		IsSpent:                     managedUtxoAsset.IsSpent,
		LeaseOwner:                  managedUtxoAsset.LeaseOwner,
		LeaseExpiry:                 managedUtxoAsset.LeaseExpiry,
		IsBurn:                      managedUtxoAsset.IsBurn,
		DeviceId:                    deviceId,
	}
}

func ManagedUtxoToAssetManagedUtxoSetRequests(deviceId string, managedUtxo *ManagedUtxo) []AssetManagedUtxoSetRequest {
	var assetManagedUtxoSetRequests []AssetManagedUtxoSetRequest
	if managedUtxo == nil || len((*managedUtxo).ManagedUtxosAssets) == 0 {
		return assetManagedUtxoSetRequests
	}
	for _, asset := range (*managedUtxo).ManagedUtxosAssets {
		request := ManagedUtxoAssetToAssetManagedUtxoSetRequest(deviceId, managedUtxo.Op, managedUtxo.OutPoint, managedUtxo.Time, managedUtxo.AmtSat, managedUtxo.InternalKey, managedUtxo.TaprootAssetRoot, managedUtxo.MerkleRoot, asset)
		assetManagedUtxoSetRequests = append(assetManagedUtxoSetRequests, request)
	}
	return assetManagedUtxoSetRequests
}

func ManagedUtxosToAssetManagedUtxoSetRequests(deviceId string, managedUtxos *[]ManagedUtxo) *[]AssetManagedUtxoSetRequest {
	if managedUtxos == nil {
		return nil
	}
	var assetManagedUtxoSetRequests []AssetManagedUtxoSetRequest
	for _, managedUtxo := range *managedUtxos {
		assetManagedUtxoSetRequests = append(assetManagedUtxoSetRequests, ManagedUtxoToAssetManagedUtxoSetRequests(deviceId, &managedUtxo)...)
	}
	return &assetManagedUtxoSetRequests
}

func PostToSetAssetManagedUtxos(token string, assetManagedUtxoSetRequests *[]AssetManagedUtxoSetRequest) (*JsonResult, error) {
	if assetManagedUtxoSetRequests == nil {
		return &JsonResult{
			Success: true,
		}, nil
	}
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_managed_utxo/set"
	requestJsonBytes, err := json.Marshal(assetManagedUtxoSetRequests)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response JsonResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return &response, nil
}

func ListUtxosAndPostToSetAssetManagedUtxos(token string, deviceId string) error {
	managedUtxos, err := ListUtxosAndGetProcessedManagedUtxos(token)
	if err != nil {
		errorAppendInfo := ErrorAppendInfo(err)
		return errorAppendInfo("List Utxos And Get Processed Managed Utxos")
	}
	assetManagedUtxoSetRequests := ManagedUtxosToAssetManagedUtxoSetRequests(deviceId, managedUtxos)
	_, err = PostToSetAssetManagedUtxos(token, assetManagedUtxoSetRequests)
	return err
}

func UploadAssetManagedUtxos(token string, deviceId string) string {
	err := ListUtxosAndPostToSetAssetManagedUtxos(token, deviceId)
	if err != nil {
		return MakeJsonErrorResult(ListUtxosAndPostToSetAssetManagedUtxosErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, nil)
}

func AssetIssuanceIsLocal(assetId string) (bool, error) {
	keys, err := assetLeafKeys(assetId, universerpc.ProofType_PROOF_TYPE_ISSUANCE)
	if err != nil {
		return false, err
	}
	if len(keys.AssetKeys) == 0 {
		return false, errors.New("asset keys is zero")
	}
	result := struct {
		IsLocal   bool   `json:"is_local"`
		AssetId   string `json:"asset_id"`
		BatchTxid string `json:"batch_txid"`
		Amount    int64  `json:"amount"`
		Timestamp int64  `json:"timestamp"`
		ScriptKey string `json:"script_key"`
	}{
		IsLocal: false,
		AssetId: assetId,
	}
	Outpoint := keys.AssetKeys[0].Outpoint
	if o, ok := Outpoint.(*universerpc.AssetKey_OpStr); ok {
		opStr := strings.Split(o.OpStr, ":")
		listBatch, err := ListBatchesAndGetResponse()
		if err != nil {
			return false, err
		}
		for _, b := range listBatch.Batches {
			if b.Batch.BatchTxid == opStr[0] {
				leaves, err := assetLeaves(false, assetId, universerpc.ProofType_PROOF_TYPE_ISSUANCE)
				if err != nil {
					return false, err
				}
				result.Amount = int64(leaves.Leaves[0].Asset.Amount)
				transactions, err := GetTransactionsAndGetResponse()
				if err != nil {
					return false, err
				}
				for _, tx := range transactions.Transactions {
					if tx.TxHash == opStr[0] {
						result.Timestamp = tx.TimeStamp
						break
					}
				}
				result.IsLocal = true
				result.BatchTxid = o.OpStr
				if s, _ok := keys.AssetKeys[0].ScriptKey.(*universerpc.AssetKey_ScriptKeyBytes); _ok {
					result.ScriptKey = "02" + hex.EncodeToString(s.ScriptKeyBytes)
				}
				break
			}
		}
		return true, nil
	}
	return false, errors.New("fail to get asset info")
}

func IsLocalMintAsset(assetId string) bool {
	isLocal, err := AssetIssuanceIsLocal(assetId)
	if err != nil {
		return false
	}
	return isLocal
}

type GetAssetLocalMintHistoryAssetIdResponse struct {
	Success bool                   `json:"success"`
	Error   string                 `json:"error"`
	Code    ErrCode                `json:"code"`
	Data    *AssetLocalMintHistory `json:"data"`
}

func RequestToGetAssetLocalMintHistoryAssetId(token string, assetId string) (*AssetLocalMintHistory, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_local_mint_history/get/asset_id/" + assetId
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetAssetLocalMintHistoryAssetIdResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

func GetAssetLocalMintHistoryAssetId(token string, assetId string) (*AssetLocalMintHistory, error) {
	return RequestToGetAssetLocalMintHistoryAssetId(token, assetId)
}

type GetCustodyAccountBalanceResponse struct {
	Balance int `json:"balance"`
}

func PostToGetCustodyAccountBalance(token string) (int, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/custodyAccount/invoice/querybalance"
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return 0, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	var response GetCustodyAccountBalanceResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}
	return response.Balance, nil
}

func GetCustodyAccountBalance(token string) (int, error) {
	return PostToGetCustodyAccountBalance(token)
}

type BtcOrAssetsValueRequest struct {
	Ids     string `json:"ids"`
	Numbers int    `json:"numbers"`
}

type BtcOrAssetsValueResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Data    struct {
		List []BtcOrAssetsValue `json:"list"`
	} `json:"data"`
}

type BtcOrAssetsValue struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
}

func ProcessBtcOrAssetsValueRequest(btcOrAssetsValueRequest []BtcOrAssetsValueRequest) string {
	var result string
	for i, request := range btcOrAssetsValueRequest {
		if request.Ids == "" || request.Numbers == 0 {
			continue
		}
		if i == 0 {
			result += "?"
		} else {
			result += "&"
		}
		result += "ids="
		result += request.Ids
		result += "&numbers="
		result += strconv.Itoa(request.Numbers)
	}
	return result
}

func RequestToGetBtcOrAssetsValue(btcOrAssetsValueRequest []BtcOrAssetsValueRequest) ([]BtcOrAssetsValue, error) {
	url := "http://api.nostr.microlinktoken.com/realtime/one_price" + ProcessBtcOrAssetsValueRequest(btcOrAssetsValueRequest)
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response BtcOrAssetsValueResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data.List, nil
}

func GetBtcOrAssetsValue(btcOrAssetsValueRequest []BtcOrAssetsValueRequest) ([]BtcOrAssetsValue, error) {
	return RequestToGetBtcOrAssetsValue(btcOrAssetsValueRequest)
}

func ListNormalBalancesAndGetResponse() (*[]ListAssetBalanceInfo, error) {
	response, err := listBalances(false, nil, nil)
	if err != nil {
		return nil, err
	}
	processed := ProcessListBalancesResponse(response)
	filtered := ExcludeListBalancesResponseCollectible(processed)
	return filtered, nil
}

func GetWalletBalanceTotalValueAndGetResponse(token string) ([]BtcOrAssetsValue, error) {
	var btcOrAssetsValueRequest []BtcOrAssetsValueRequest
	walletBalance, err := getWalletBalance()
	if err != nil {
		return nil, err
	}
	custodyAccountBalance, err := GetCustodyAccountBalance(token)
	if err != nil {
		return nil, err
	}
	btc := int(walletBalance.TotalBalance) + custodyAccountBalance
	btcOrAssetsValueRequest = append(btcOrAssetsValueRequest, BtcOrAssetsValueRequest{
		Ids:     "btc",
		Numbers: btc,
	})
	normalBalances, err := ListNormalBalancesAndGetResponse()
	if err != nil {
		return nil, err
	}
	for _, normalBalance := range *normalBalances {
		var balance int
		balance, err = strconv.Atoi(normalBalance.Balance)
		if err != nil {
			continue
		}
		btcOrAssetsValueRequest = append(btcOrAssetsValueRequest, BtcOrAssetsValueRequest{
			Ids:     normalBalance.AssetID,
			Numbers: balance,
		})
	}
	nftAssets, err := ListNftAssetsAndGetResponse()
	if err != nil {
		return nil, err
	}
	for _, nftAsset := range *nftAssets {
		btcOrAssetsValueRequest = append(btcOrAssetsValueRequest, BtcOrAssetsValueRequest{
			Ids:     nftAsset.AssetGenesis.AssetID,
			Numbers: nftAsset.Amount,
		})
	}
	btcOrAssetsValue, err := GetBtcOrAssetsValue(btcOrAssetsValueRequest)
	if err != nil {
		return nil, err
	}
	return btcOrAssetsValue, nil
}

func GetWalletBalanceCalculatedTotalValue(token string) (float64, error) {
	btcOrAssetsValues, err := GetWalletBalanceTotalValueAndGetResponse(token)
	if err != nil {
		return 0, err
	}
	var totalValue float64
	for _, btcOrAssetsValue := range btcOrAssetsValues {
		totalValue += btcOrAssetsValue.Price
	}
	return totalValue, nil
}

func GetWalletBalanceTotalValue(token string) string {
	totalValue, err := GetWalletBalanceCalculatedTotalValue(token)
	if err != nil {
		return MakeJsonErrorResult(GetWalletBalanceCalculatedTotalValueErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SUCCESS.Error(), totalValue)
}

type GetAssetBalanceByUserIdAndAssetIdRequest struct {
	UserId  int    `json:"user_id"`
	AssetId string `json:"asset_id"`
}

type GetAssetBalanceByUserIdAndAssetIdResponse struct {
	Success bool          `json:"success"`
	Error   string        `json:"error"`
	Code    ErrCode       `json:"code"`
	Data    *AssetBalance `json:"data"`
}

func PostToGetAssetBalanceByUserIdAndAssetId(token string, getAssetBalanceByUserIdAndAssetIdRequest GetAssetBalanceByUserIdAndAssetIdRequest) (*AssetBalance, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_balance/get/balance/asset_id_and_user_id"
	requestJsonBytes, err := json.Marshal(getAssetBalanceByUserIdAndAssetIdRequest)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetAssetBalanceByUserIdAndAssetIdResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

func GetAssetBalanceByUserIdAndAssetId(token string, getAssetBalanceByUserIdAndAssetIdRequest GetAssetBalanceByUserIdAndAssetIdRequest) (*AssetBalance, error) {
	return PostToGetAssetBalanceByUserIdAndAssetId(token, getAssetBalanceByUserIdAndAssetIdRequest)
}

func GetAssetRecommendUser(token string, assetId string, encoded string, deviceId string) (string, error) {
	assetRecipientIsAssetIssuer := errors.New("asset recipient is asset issuer")
	addr, err := GetAssetAddrByEncoded(token, encoded)
	if err != nil {
		return "", err
	}
	addrUserId := addr.UserId
	assetLocalMintHistory, err := GetAssetLocalMintHistoryAssetId(token, assetId)
	if err != nil {
		return "", err
	}
	if assetLocalMintHistory.UserId == addrUserId {
		return "", assetRecipientIsAssetIssuer
	}
	request := GetAssetBalanceByUserIdAndAssetIdRequest{
		UserId:  addrUserId,
		AssetId: addr.AssetId,
	}
	assetBalance, err := GetAssetBalanceByUserIdAndAssetId(token, request)
	var userHadAsset bool
	if err == nil && assetBalance != nil {
		userHadAsset = true
	}
	userHadAssetErr := errors.New("user had asset")
	var assetRecommend *AssetRecommend
	assetRecommend, err = GetAssetRecommendByUserIdAndAssetId(token, UserIdAndAssetId{
		UserId:  addrUserId,
		AssetId: assetId,
	})
	if err != nil {
		if userHadAsset {
			return "", userHadAssetErr
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {

			err = SetAssetRecommendByAssetId(token, assetId, encoded, 0, "", 0, deviceId)
			if err != nil {
				return "", err
			}
			return GetNPublicKey(), nil
		}
		return "", err
	}
	return assetRecommend.RecommendUsername, nil
}

func GetAssetRecommendUserByJsonAddrs(token string, assetId string, jsonAddrs string, deviceId string) (*map[string]string, error) {
	var addrs []string
	err := json.Unmarshal([]byte(jsonAddrs), &addrs)
	if err != nil {
		return nil, err
	}
	addrMapRecommendUser := make(map[string]string)
	var recommendUser string
	for _, addr := range addrs {
		recommendUser, err = GetAssetRecommendUser(token, assetId, addr, deviceId)
		if err != nil {
			continue
		}
		addrMapRecommendUser[addr] = recommendUser
	}
	return &addrMapRecommendUser, nil
}

func UploadLogFileAndGetResponse(filePath string, deviceId string, info string, auth string) (*uint, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/log_file_upload/upload"
	stat, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	if stat.Size() > 15*1024*1024 {
		return nil, errors.New("file too large, its size is more than 15MB")
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			LogError("", err)
		}
	}(file)
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	_ = writer.WriteField("device_id", deviceId)
	_ = writer.WriteField("info", info)
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, requestBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))
	req.Header.Add("accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			LogError("", err)
		}
	}(resp.Body)
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response UploadLogFileResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

type UploadBigFileResponse struct {
	Success bool    `json:"success"`
	Error   string  `json:"error"`
	Code    ErrCode `json:"code"`
	Data    *uint   `json:"data"`
}

type UploadLogFileResponse struct {
	Success bool    `json:"success"`
	Error   string  `json:"error"`
	Code    ErrCode `json:"code"`
	Data    *uint   `json:"data"`
}

func UploadBigFileAndGetResponse(filePath string, deviceId string, info string, auth string) (*uint, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/log_file_upload/upload_big"
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			LogError("", err)
		}
	}(file)
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	_ = writer.WriteField("device_id", deviceId)
	_ = writer.WriteField("info", info)
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, requestBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))
	req.Header.Add("accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			LogError("", err)
		}
	}(resp.Body)
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response UploadBigFileResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

func UploadLogFile(filePath string, deviceId string, info string, auth string) string {
	id, err := UploadLogFileAndGetResponse(filePath, deviceId, info, auth)
	if err != nil {
		return MakeJsonErrorResult(UploadLogFileAndGetResponseErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SUCCESS.Error(), id)
}

func UploadBigFile(filePath string, deviceId string, info string, auth string) string {
	id, err := UploadBigFileAndGetResponse(filePath, deviceId, info, auth)
	if err != nil {
		return MakeJsonErrorResult(UploadBigFileAndGetResponseErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SUCCESS.Error(), id)
}

type AccountAssetBalanceExtend struct {
	AccountID uint   ` json:"account_id"`
	AssetId   string ` json:"asset_id"`
	Amount    int    ` json:"amount"`
	UserID    int    ` json:"user_id"`
	Username  string ` json:"username"`
}

type GetAccountAssetBalanceByAssetIdResponse struct {
	Success bool                         `json:"success"`
	Error   string                       `json:"error"`
	Code    ErrCode                      `json:"code"`
	Data    *[]AccountAssetBalanceExtend `json:"data"`
}

func RequestToGetAccountAssetBalanceByAssetId(token string, assetId string) (*[]AccountAssetBalanceExtend, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/account_asset/balance/get/asset_id/" + assetId
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetAccountAssetBalanceByAssetIdResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

func GetAccountAssetBalanceByAssetIdAndGetResponse(token string, assetId string) (*[]AccountAssetBalanceExtend, error) {
	return RequestToGetAccountAssetBalanceByAssetId(token, assetId)
}

type GetAccountAssetBalanceUserHoldByAssetIdResponse struct {
	Success bool    `json:"success"`
	Error   string  `json:"error"`
	Code    ErrCode `json:"code"`
	Data    int     `json:"data"`
}

func RequestToGetAccountAssetBalanceUserHoldTotalAmountByAssetId(token string, assetId string) (int, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/account_asset/balance/query/total_amount?asset_id=" + assetId
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return 0, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	var response GetAccountAssetBalanceUserHoldByAssetIdResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}
	if response.Error != "" {
		return 0, errors.New(response.Error)
	}
	return response.Data, nil
}

func GetAccountAssetBalanceUserHoldTotalAmountByAssetId(token string, assetId string) (int, error) {
	return RequestToGetAccountAssetBalanceUserHoldTotalAmountByAssetId(token, assetId)
}

func GetAccountAssetBalanceUserHoldTotalAmount(token string, assetId string) string {
	totalAmount, err := GetAccountAssetBalanceUserHoldTotalAmountByAssetId(token, assetId)
	if err != nil {
		return MakeJsonErrorResult(GetAccountAssetBalanceUserHoldTotalAmountByAssetIdErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SUCCESS.Error(), totalAmount)
}

type AssetIdAndAccountAssetBalanceExtends struct {
	AssetId                    string                       `json:"asset_id"`
	AccountAssetBalanceExtends *[]AccountAssetBalanceExtend `json:"account_asset_balance_extends"`
}

func AccountAssetBalanceExtendsToAssetIdAndAccountAssetBalanceExtends(assetId string, accountAssetBalanceExtends *[]AccountAssetBalanceExtend) *AssetIdAndAccountAssetBalanceExtends {
	if accountAssetBalanceExtends == nil {
		return nil
	}
	return &AssetIdAndAccountAssetBalanceExtends{
		AssetId:                    assetId,
		AccountAssetBalanceExtends: accountAssetBalanceExtends,
	}
}

func GetAccountAssetBalances(token string, assetId string) string {
	accountAssetBalanceExtends, err := GetAccountAssetBalanceByAssetIdAndGetResponse(token, assetId)
	if err != nil {
		return MakeJsonErrorResult(GetAccountAssetBalanceByAssetIdAndGetResponseErr, err.Error(), nil)
	}
	assetIdAndAccountAssetBalanceExtends := AccountAssetBalanceExtendsToAssetIdAndAccountAssetBalanceExtends(assetId, accountAssetBalanceExtends)
	return MakeJsonErrorResult(SUCCESS, SUCCESS.Error(), assetIdAndAccountAssetBalanceExtends)
}

type AccountAssetTransfer struct {
	BillBalanceId int    `json:"bill_balance_id"`
	AccountId     int    `json:"account_id"`
	Username      string `json:"username"`
	BillType      string `json:"bill_type"`
	Away          string `json:"away"`
	Amount        int    `json:"amount"`
	ServerFee     int    `json:"server_fee"`
	AssetId       string `json:"asset_id"`
	Invoice       string `json:"invoice"`
	Outpoint      string `json:"outpoint"`
	Time          int    `json:"time"`
}

type GetAccountAssetTransferByAssetId struct {
	Success bool                    `json:"success"`
	Error   string                  `json:"error"`
	Code    ErrCode                 `json:"code"`
	Data    *[]AccountAssetTransfer `json:"data"`
}

func RequestToGetAccountAssetTransferByAssetId(token string, assetId string) (*[]AccountAssetTransfer, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/account_asset/transfer/get/asset_id/" + assetId
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetAccountAssetTransferByAssetId
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

func GetAccountAssetTransferByAssetIdAndGetResponse(token string, assetId string) (*[]AccountAssetTransfer, error) {
	return RequestToGetAccountAssetTransferByAssetId(token, assetId)
}

func GetAccountAssetTransfers(token string, assetId string) string {
	accountAssetTransfers, err := GetAccountAssetTransferByAssetIdAndGetResponse(token, assetId)
	if err != nil {
		return MakeJsonErrorResult(GetAccountAssetTransferByAssetIdAndGetResponseErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SUCCESS.Error(), accountAssetTransfers)
}

type GetAccountAssetBalanceLimitAndOffsetRequest struct {
	AssetId string `json:"asset_id"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
}

type GetAccountAssetBalanceLimitAndOffsetResponse struct {
	Success bool                         `json:"success"`
	Error   string                       `json:"error"`
	Code    ErrCode                      `json:"code"`
	Data    *[]AccountAssetBalanceExtend `json:"data"`
}

func PostToGetAccountAssetBalanceLimitAndOffset(token string, assetId string, limit int, offset int) (*[]AccountAssetBalanceExtend, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	assetIdLimitAndOffset := GetAccountAssetBalanceLimitAndOffsetRequest{
		AssetId: assetId,
		Limit:   limit,
		Offset:  offset,
	}
	url := serverDomainOrSocket + "/account_asset/balance/get/limit_offset"
	requestJsonBytes, err := json.Marshal(assetIdLimitAndOffset)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetAccountAssetBalanceLimitAndOffsetResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

func GetAccountAssetBalanceLimitAndOffset(token string, assetId string, limit int, offset int) (*[]AccountAssetBalanceExtend, error) {
	return PostToGetAccountAssetBalanceLimitAndOffset(token, assetId, limit, offset)
}

type GetAccountAssetBalancePageNumberByPageSizeRequest struct {
	AssetId  string `json:"asset_id"`
	PageSize int    `json:"page_size"`
}

type GetAccountAssetBalancePageNumberByPageSizeResponse struct {
	Success bool    `json:"success"`
	Error   string  `json:"error"`
	Code    ErrCode `json:"code"`
	Data    int     `json:"data"`
}

func PostToGetAccountAssetBalancePageNumberByPageSize(token string, assetId string, pageSize int) (int, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	getAssetHolderBalancePageNumberRequest := GetAccountAssetBalancePageNumberByPageSizeRequest{
		AssetId:  assetId,
		PageSize: pageSize,
	}
	url := serverDomainOrSocket + "/account_asset/balance/get/page_number"
	requestJsonBytes, err := json.Marshal(getAssetHolderBalancePageNumberRequest)
	if err != nil {
		return 0, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	var response GetAccountAssetBalancePageNumberByPageSizeResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}
	if response.Error != "" {
		return 0, errors.New(response.Error)
	}
	return response.Data, nil
}

func GetAccountAssetBalancePageNumberByPageSize(token string, assetId string, pageSize int) (int, error) {
	pageNumber, err := PostToGetAccountAssetBalancePageNumberByPageSize(token, assetId, pageSize)
	if err != nil {
		return 0, err
	}
	return pageNumber, nil
}

func GetAccountAssetBalancePageNumber(token string, assetId string, pageSize int) string {
	pageNumber, err := GetAccountAssetBalancePageNumberByPageSize(token, assetId, pageSize)
	if err != nil {
		return MakeJsonErrorResult(GetAccountAssetBalancePageNumberByPageSizeErr, err.Error(), 0)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, pageNumber)
}

func GetAccountAssetBalanceWithPageSizeAndPageNumber(token string, assetId string, pageSize int, pageNumber int) (*[]AccountAssetBalanceExtend, error) {
	if !(pageSize > 0 && pageNumber > 0) {
		return nil, errors.New("page size and page number must be greater than 0")
	}
	var limit int
	var offset int
	limit = pageSize
	if pageNumber > 1 {
		offset = (pageNumber - 1) * pageSize
	}
	return GetAccountAssetBalanceLimitAndOffset(token, assetId, limit, offset)
}

func GetAccountAssetBalancePage(token string, assetId string, pageSize int, pageNumber int) string {
	accountAssetTransfers, err := GetAccountAssetBalanceWithPageSizeAndPageNumber(token, assetId, pageSize, pageNumber)
	if err != nil {
		return MakeJsonErrorResult(GetAccountAssetBalanceWithPageSizeAndPageNumberErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, accountAssetTransfers)
}

type GetAccountAssetTransferLimitAndOffsetRequest struct {
	AssetId string `json:"asset_id"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
}

type GetAccountAssetTransferLimitAndOffsetResponse struct {
	Success bool                    `json:"success"`
	Error   string                  `json:"error"`
	Code    ErrCode                 `json:"code"`
	Data    *[]AccountAssetTransfer `json:"data"`
}

func PostToGetAccountAssetTransferLimitAndOffset(token string, assetId string, limit int, offset int) (*[]AccountAssetTransfer, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	assetIdLimitAndOffset := GetAccountAssetTransferLimitAndOffsetRequest{
		AssetId: assetId,
		Limit:   limit,
		Offset:  offset,
	}
	url := serverDomainOrSocket + "/account_asset/transfer/get/limit_offset"
	requestJsonBytes, err := json.Marshal(assetIdLimitAndOffset)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetAccountAssetTransferLimitAndOffsetResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

func GetAccountAssetTransferLimitAndOffset(token string, assetId string, limit int, offset int) (*[]AccountAssetTransfer, error) {
	return PostToGetAccountAssetTransferLimitAndOffset(token, assetId, limit, offset)
}

type GetAccountAssetTransferPageNumberByPageSizeRequest struct {
	AssetId  string `json:"asset_id"`
	PageSize int    `json:"page_size"`
}

type GetAccountAssetTransferPageNumberByPageSizeResponse struct {
	Success bool    `json:"success"`
	Error   string  `json:"error"`
	Code    ErrCode `json:"code"`
	Data    int     `json:"data"`
}

func PostToGetAccountAssetTransferPageNumberByPageSize(token string, assetId string, pageSize int) (int, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	getAssetHolderBalancePageNumberRequest := GetAccountAssetTransferPageNumberByPageSizeRequest{
		AssetId:  assetId,
		PageSize: pageSize,
	}
	url := serverDomainOrSocket + "/account_asset/transfer/get/page_number"
	requestJsonBytes, err := json.Marshal(getAssetHolderBalancePageNumberRequest)
	if err != nil {
		return 0, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	var response GetAccountAssetTransferPageNumberByPageSizeResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}
	if response.Error != "" {
		return 0, errors.New(response.Error)
	}
	return response.Data, nil
}

func GetAccountAssetTransferPageNumberByPageSize(token string, assetId string, pageSize int) (int, error) {
	pageNumber, err := PostToGetAccountAssetTransferPageNumberByPageSize(token, assetId, pageSize)
	if err != nil {
		return 0, err
	}
	return pageNumber, nil
}

func GetAccountAssetTransferPageNumber(token string, assetId string, pageSize int) string {
	pageNumber, err := GetAccountAssetTransferPageNumberByPageSize(token, assetId, pageSize)
	if err != nil {
		return MakeJsonErrorResult(GetAccountAssetTransferPageNumberByPageSizeErr, err.Error(), 0)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, pageNumber)
}

func GetAccountAssetTransferWithPageSizeAndPageNumber(token string, assetId string, pageSize int, pageNumber int) (*[]AccountAssetTransfer, error) {
	if !(pageSize > 0 && pageNumber > 0) {
		return nil, errors.New("page size and page number must be greater than 0")
	}
	var limit int
	var offset int
	limit = pageSize
	if pageNumber > 1 {
		offset = (pageNumber - 1) * pageSize
	}
	return GetAccountAssetTransferLimitAndOffset(token, assetId, limit, offset)
}

func GetAccountAssetTransfersPage(token string, assetId string, pageSize int, pageNumber int) string {
	accountAssetTransfers, err := GetAccountAssetTransferWithPageSizeAndPageNumber(token, assetId, pageSize, pageNumber)
	if err != nil {
		return MakeJsonErrorResult(GetAccountAssetTransferWithPageSizeAndPageNumberErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, accountAssetTransfers)
}

func GetAssetHolderBalanceByAssetBalancesInfoLimitAndOffset(token string, assetId string, limit int, offset int) (*AssetIdAndBalance, error) {
	holderBalance, err := PostToGetAssetHolderBalanceLimitAndOffsetByAssetBalancesInfo(token, assetId, limit, offset)
	if err != nil {
		return nil, err
	}
	return holderBalance, nil
}

type GetAssetHolderBalancePageNumberByPageSizeResponse struct {
	Success bool    `json:"success"`
	Error   string  `json:"error"`
	Code    ErrCode `json:"code"`
	Data    int     `json:"data"`
}

type GetAssetHolderBalancePageNumberRequest struct {
	AssetId  string `json:"asset_id"`
	PageSize int    `json:"page_size"`
}

func PostToGetAssetHolderBalancePageNumberByPageSize(token string, assetId string, pageSize int) (int, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	getAssetHolderBalancePageNumberRequest := GetAssetHolderBalancePageNumberRequest{
		AssetId:  assetId,
		PageSize: pageSize,
	}
	url := serverDomainOrSocket + "/asset_balance/get/holder/balance/page_number"
	requestJsonBytes, err := json.Marshal(getAssetHolderBalancePageNumberRequest)
	if err != nil {
		return 0, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	var response GetAssetHolderBalancePageNumberByPageSizeResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}
	if response.Error != "" {
		return 0, errors.New(response.Error)
	}
	return response.Data, nil
}

func GetAssetHolderBalancePageNumberByPageSize(token string, assetId string, pageSize int) (int, error) {
	pageNumber, err := PostToGetAssetHolderBalancePageNumberByPageSize(token, assetId, pageSize)
	if err != nil {
		return 0, err
	}
	return pageNumber, nil
}

func GetAssetHolderBalancePageNumber(token string, assetId string, pageSize int) string {
	pageNumber, err := GetAssetHolderBalancePageNumberByPageSize(token, assetId, pageSize)
	if err != nil {
		return MakeJsonErrorResult(GetAssetHolderBalanceWithPageSizeAndPageNumberErr, err.Error(), 0)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, pageNumber)
}

func GetAssetHolderBalanceWithPageSizeAndPageNumber(token string, assetId string, pageSize int, pageNumber int) (*AssetIdAndBalance, error) {
	if !(pageSize > 0 && pageNumber > 0) {
		return nil, errors.New("page size and page number must be greater than 0")
	}
	var limit int
	var offset int
	limit = pageSize
	if pageNumber > 1 {
		offset = (pageNumber - 1) * pageSize
	}
	return GetAssetHolderBalanceByAssetBalancesInfoLimitAndOffset(token, assetId, limit, offset)
}

func GetAssetHolderBalancePage(token string, assetId string, pageSize int, pageNumber int) string {
	holderBalance, err := GetAssetHolderBalanceWithPageSizeAndPageNumber(token, assetId, pageSize, pageNumber)
	if err != nil {
		return MakeJsonErrorResult(GetAssetHolderBalanceWithPageSizeAndPageNumberErr, err.Error(), nil)
	}
	result := AssetIdAndBalanceToAssetIdAndBalanceSimplified(holderBalance)
	return MakeJsonErrorResult(SUCCESS, SuccessError, result)
}

type GetAssetManagedUtxoPageNumberByPageSizeResponse struct {
	Success bool    `json:"success"`
	Error   string  `json:"error"`
	Code    ErrCode `json:"code"`
	Data    int     `json:"data"`
}

type GetAssetManagedUtxoPageNumberRequest struct {
	AssetId  string `json:"asset_id"`
	PageSize int    `json:"page_size"`
}

func PostToGetAssetManagedUtxoPageNumberByPageSize(token string, assetId string, pageSize int) (int, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	getAssetHolderBalancePageNumberRequest := GetAssetManagedUtxoPageNumberRequest{
		AssetId:  assetId,
		PageSize: pageSize,
	}
	url := serverDomainOrSocket + "/asset_managed_utxo/get/page_number"
	requestJsonBytes, err := json.Marshal(getAssetHolderBalancePageNumberRequest)
	if err != nil {
		return 0, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	var response GetAssetManagedUtxoPageNumberByPageSizeResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}
	if response.Error != "" {
		return 0, errors.New(response.Error)
	}
	return response.Data, nil
}

func GetAssetManagedUtxoPageNumberByPageSize(token string, assetId string, pageSize int) (int, error) {
	pageNumber, err := PostToGetAssetManagedUtxoPageNumberByPageSize(token, assetId, pageSize)
	if err != nil {
		return 0, err
	}
	return pageNumber, nil
}

type GetAssetManagedUtxoLimitAndOffsetRequest struct {
	AssetId string `json:"asset_id"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
}

type GetAssetManagedUtxoLimitAndOffsetResponse struct {
	Success bool                `json:"success"`
	Error   string              `json:"error"`
	Code    ErrCode             `json:"code"`
	Data    *[]AssetManagedUtxo `json:"data"`
}

func PostToGetAssetManagedUtxoLimitAndOffset(token string, assetId string, limit int, offset int) (*[]AssetManagedUtxo, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	assetIdLimitAndOffset := GetAssetManagedUtxoLimitAndOffsetRequest{
		AssetId: assetId,
		Limit:   limit,
		Offset:  offset,
	}
	url := serverDomainOrSocket + "/asset_managed_utxo/get/limit_offset"
	requestJsonBytes, err := json.Marshal(assetIdLimitAndOffset)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetAssetManagedUtxoLimitAndOffsetResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

func GetAssetManagedUtxoLimitAndOffset(token string, assetId string, limit int, offset int) (*[]AssetManagedUtxo, error) {
	assetManagedUtxo, err := PostToGetAssetManagedUtxoLimitAndOffset(token, assetId, limit, offset)
	if err != nil {
		return nil, err
	}
	return assetManagedUtxo, nil
}

func GetAssetManagedUtxoWithPageSizeAndPageNumber(token string, assetId string, pageSize int, pageNumber int) (*[]AssetManagedUtxo, error) {
	if !(pageSize > 0 && pageNumber > 0) {
		return nil, errors.New("page size and page number must be greater than 0")
	}
	var limit int
	var offset int
	limit = pageSize
	if pageNumber > 1 {
		offset = (pageNumber - 1) * pageSize
	}
	return GetAssetManagedUtxoLimitAndOffset(token, assetId, limit, offset)
}

func GetAssetManagedUtxoPage(token string, assetId string, pageSize int, pageNumber int) string {
	assetManagedUtxo, err := GetAssetManagedUtxoWithPageSizeAndPageNumber(token, assetId, pageSize, pageNumber)
	if err != nil {
		return MakeJsonErrorResult(GetAssetManagedUtxoWithPageSizeAndPageNumberErr, err.Error(), nil)
	}
	result := AssetManagedUtxoSliceToAssetManagedUtxoSimplifiedSlice(assetManagedUtxo)
	return MakeJsonErrorResult(SUCCESS, SuccessError, result)
}

func GetAssetManagedUtxoPageNumber(token string, assetId string, pageSize int) string {
	pageNumber, err := GetAssetManagedUtxoPageNumberByPageSize(token, assetId, pageSize)
	if err != nil {
		return MakeJsonErrorResult(GetAssetManagedUtxoPageNumberByPageSizeErr, err.Error(), 0)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, pageNumber)
}

type AssetManagedUtxoSimplified struct {
	UpdatedAt             time.Time `json:"updated_at"`
	OutPoint              string    `json:"out_point"`
	Time                  int       `json:"time"`
	AmtSat                int       `json:"amt_sat"`
	AssetGenesisPoint     string    `json:"asset_genesis_point"`
	AssetGenesisName      string    `json:"asset_genesis_name"`
	AssetGenesisMetaHash  string    `json:"asset_genesis_meta_hash"`
	AssetGenesisAssetID   string    `json:"asset_genesis_asset_id"`
	AssetGenesisAssetType string    `json:"asset_genesis_asset_type"`
	Amount                int       `json:"amount"`
	LockTime              int       `json:"lock_time"`
	RelativeLockTime      int       `json:"relative_lock_time"`
	ScriptKey             string    `json:"script_key"`
	AssetGroupRawGroupKey string    `json:"asset_group_raw_group_key"`
	ChainAnchorOutpoint   string    `json:"chain_anchor_outpoint"`
	IsSpent               bool      `json:"is_spent"`
	IsBurn                bool      `json:"is_burn"`
	DeviceId              string    `json:"device_id"`
	Username              string    `json:"username"`
}

func AssetManagedUtxoToAssetManagedUtxoSimplified(assetManagedUtxo AssetManagedUtxo) AssetManagedUtxoSimplified {
	return AssetManagedUtxoSimplified{
		UpdatedAt:             assetManagedUtxo.UpdatedAt,
		OutPoint:              assetManagedUtxo.OutPoint,
		Time:                  assetManagedUtxo.Time,
		AmtSat:                assetManagedUtxo.AmtSat,
		AssetGenesisPoint:     assetManagedUtxo.AssetGenesisPoint,
		AssetGenesisName:      assetManagedUtxo.AssetGenesisName,
		AssetGenesisMetaHash:  assetManagedUtxo.AssetGenesisMetaHash,
		AssetGenesisAssetID:   assetManagedUtxo.AssetGenesisAssetID,
		AssetGenesisAssetType: assetManagedUtxo.AssetGenesisAssetType,
		Amount:                assetManagedUtxo.Amount,
		LockTime:              assetManagedUtxo.LockTime,
		RelativeLockTime:      assetManagedUtxo.RelativeLockTime,
		ScriptKey:             assetManagedUtxo.ScriptKey,
		AssetGroupRawGroupKey: assetManagedUtxo.AssetGroupRawGroupKey,
		ChainAnchorOutpoint:   assetManagedUtxo.ChainAnchorOutpoint,
		IsSpent:               assetManagedUtxo.IsSpent,
		IsBurn:                assetManagedUtxo.IsBurn,
		DeviceId:              assetManagedUtxo.DeviceId,
		Username:              assetManagedUtxo.Username,
	}
}

func AssetManagedUtxoSliceToAssetManagedUtxoSimplifiedSlice(assetManagedUtxos *[]AssetManagedUtxo) *[]AssetManagedUtxoSimplified {
	if assetManagedUtxos == nil {
		return nil
	}
	var assetManagedUtxoSimplified []AssetManagedUtxoSimplified
	for _, assetManagedUtxo := range *assetManagedUtxos {
		assetManagedUtxoSimplified = append(assetManagedUtxoSimplified, AssetManagedUtxoToAssetManagedUtxoSimplified(assetManagedUtxo))
	}
	return &assetManagedUtxoSimplified
}

type AssetGroupSetRequest struct {
	TweakedGroupKey string `json:"tweaked_group_key"`
	FirstAssetMeta  string `json:"first_asset_meta"`
	FirstAssetId    string `json:"first_asset_id" gorm:"type:varchar(255)"`
	DeviceId        string `json:"device_id" gorm:"type:varchar(255)"`
}

type GetGroupFirstAssetMetaResponse struct {
	Success bool    `json:"success"`
	Error   string  `json:"error"`
	Code    ErrCode `json:"code"`
	Data    string  `json:"data"`
}

func RequestToGetGroupFirstAssetMeta(token string, groupKey string) (string, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_group/get/first_meta/group_key/" + groupKey
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return "", err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var response GetGroupFirstAssetMetaResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}
	if response.Error != "" {
		return "", errors.New(response.Error)
	}
	return response.Data, nil
}

func GetGroupFirstAssetMetaAndGetResponse(token string, groupKey string) (string, error) {
	assetMeta, err := RequestToGetGroupFirstAssetMeta(token, groupKey)
	if err != nil {
		return "", err
	}
	return assetMeta, nil
}

func PostToSetGroupFirstAssetMeta(token string, assetGroupSetRequest *AssetGroupSetRequest) (*JsonResult, error) {
	if assetGroupSetRequest == nil {
		return &JsonResult{
			Success: true,
		}, nil
	}
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_group/set/first_meta/"
	requestJsonBytes, err := json.Marshal(assetGroupSetRequest)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response JsonResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return &response, nil
}

func SetGroupFirstAssetMetaAndGetResponse(token string, tweakedGroupKey string, firstAssetMeta string, firstAssetId string, deviceId string) (*JsonResult, error) {
	assetGroupSetRequest := &AssetGroupSetRequest{
		TweakedGroupKey: tweakedGroupKey,
		FirstAssetMeta:  firstAssetMeta,
		FirstAssetId:    firstAssetId,
		DeviceId:        deviceId,
	}
	return PostToSetGroupFirstAssetMeta(token, assetGroupSetRequest)
}

type GetGroupFirstAssetIdResponse struct {
	Success bool    `json:"success"`
	Error   string  `json:"error"`
	Code    ErrCode `json:"code"`
	Data    string  `json:"data"`
}

func RequestToGetGroupFirstAssetId(token string, groupKey string) (string, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_group/get/first_asset_id/group_key/" + groupKey
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return "", err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var response GetGroupFirstAssetMetaResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}
	if response.Error != "" {
		return "", errors.New(response.Error)
	}
	return response.Data, nil
}

func GetGroupFirstAssetIdAndGetResponse(token string, groupKey string) (string, error) {
	assetMeta, err := RequestToGetGroupFirstAssetId(token, groupKey)
	if err != nil {
		return "", err
	}
	return assetMeta, nil
}

func GetGroupFirstAssetMeta(token string, groupKey string) string {
	assetMeta, err := GetGroupFirstAssetMetaAndGetResponse(token, groupKey)
	if err != nil {
		return MakeJsonErrorResult(GetGroupFirstAssetMetaAndGetResponseErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, assetMeta)
}

func GetGroupFirstAssetId(token string, groupKey string) string {
	assetMeta, err := GetGroupFirstAssetIdAndGetResponse(token, groupKey)
	if err != nil {
		return MakeJsonErrorResult(GetGroupFirstAssetIdAndGetResponseErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, assetMeta)
}

func SetGroupFirstAssetMeta(token string, deviceId string, finalizeBatchResponse *mintrpc.FinalizeBatchResponse) error {
	assetLocalMintSetRequests := FinalizeBatchResponseToAssetLocalMintSetRequests(deviceId, finalizeBatchResponse)
	var firstErr error
	for _, request := range *assetLocalMintSetRequests {
		tweakedGroupKey := request.GroupKey
		if tweakedGroupKey == "" {
			continue
		}
		firstAssetMeta := request.AssetMetaData
		firstAssetId := request.AssetId
		_, err := SetGroupFirstAssetMetaAndGetResponse(token, tweakedGroupKey, firstAssetMeta, firstAssetId, deviceId)
		if err != nil {
			fmt.Println(err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

type DeliverProofNeedInfo struct {
	AssetId   string `json:"asset_id"`
	GroupKey  string `json:"group_key"`
	ScriptKey string `json:"script_key"`
	Outpoint  string `json:"outpoint"`
}

func GetDeliverProofNeedInfoAndGetResponse(assetId string) (*DeliverProofNeedInfo, error) {
	response, err := listAssets(true, true, false)
	if err != nil {
		return nil, err
	}
	if len(response.Assets) == 0 {
		return nil, errors.New("asset list response null")
	}
	for _, asset := range response.Assets {
		assetGenesisAssetId := hex.EncodeToString(asset.AssetGenesis.AssetId)
		if assetGenesisAssetId == assetId {
			deliverProofNeedInfo := DeliverProofNeedInfo{
				AssetId: assetId,
			}
			if asset.AssetGroup != nil && asset.AssetGroup.TweakedGroupKey != nil {
				deliverProofNeedInfo.GroupKey = hex.EncodeToString(asset.AssetGroup.TweakedGroupKey)
			}
			if asset.ScriptKey != nil {
				deliverProofNeedInfo.ScriptKey = hex.EncodeToString(asset.ScriptKey)
			}
			deliverProofNeedInfo.Outpoint = asset.ChainAnchor.AnchorOutpoint
			return &deliverProofNeedInfo, nil
		}
	}
	return nil, errors.New("asset not found")
}

func GetDeliverProofNeedInfo(assetId string) string {
	deliverProofNeedInfo, err := GetDeliverProofNeedInfoAndGetResponse(assetId)
	if err != nil {
		MakeJsonErrorResult(GetDeliverProofNeedInfoAndGetResponseErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, deliverProofNeedInfo)
}

type NftTransfer struct {
	gorm.Model
	Txid     string `json:"txid" gorm:"type:varchar(255)"`
	AssetId  string `json:"asset_id" gorm:"type:varchar(255);index"`
	Time     int    `json:"time"`
	FromAddr string `json:"from_addr"`
	ToAddr   string `json:"to_addr"`
	FromInfo string `json:"from_info"`
	ToInfo   string `json:"to_info"`
	DeviceId string `json:"device_id" gorm:"type:varchar(255);index"`
	UserId   int    `json:"user_id" gorm:"index"`
	Username string `json:"username" gorm:"type:varchar(255);index"`
}

type NftTransferSetRequest struct {
	Txid     string `json:"txid" gorm:"type:varchar(255)"`
	AssetId  string `json:"asset_id" gorm:"type:varchar(255);index"`
	Time     int    `json:"time"`
	FromAddr string `json:"from_addr"`
	ToAddr   string `json:"to_addr"`
	FromInfo string `json:"from_info"`
	ToInfo   string `json:"to_info"`
	DeviceId string `json:"device_id" gorm:"type:varchar(255);index"`
}

type NftTransferSimplified struct {
	ID       uint   `gorm:"primarykey"`
	Txid     string `json:"txid" gorm:"type:varchar(255)"`
	AssetId  string `json:"asset_id" gorm:"type:varchar(255);index"`
	Time     int    `json:"time"`
	FromAddr string `json:"from_addr"`
	ToAddr   string `json:"to_addr"`
}

func PostToSetNftTransfer(token string, nftTransferSetRequest *NftTransferSetRequest) (*JsonResult, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/nft_transfer/set"
	requestJsonBytes, err := json.Marshal(nftTransferSetRequest)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response JsonResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return &response, nil
}

func SetNftTransferAndGetResponse(token string, txid string, assetId string, time int, fromAddr string, toAddr string, fromInfo string, toInfo string, deviceId string) (*JsonResult, error) {
	assetGroupSetRequest := &NftTransferSetRequest{
		Txid:     txid,
		AssetId:  assetId,
		Time:     time,
		FromAddr: fromAddr,
		ToAddr:   toAddr,
		FromInfo: fromInfo,
		ToInfo:   toInfo,
		DeviceId: deviceId,
	}
	return PostToSetNftTransfer(token, assetGroupSetRequest)
}

type GetNftTransferByAssetIdResponse struct {
	Success bool           `json:"success"`
	Error   string         `json:"error"`
	Code    ErrCode        `json:"code"`
	Data    *[]NftTransfer `json:"data"`
}

func RequestToGetNftTransferByAssetId(token string, assetId string) (*[]NftTransfer, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/nft_transfer/get/asset_id/" + assetId
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetNftTransferByAssetIdResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

func GetNftTransferByAssetIdAndGetResponse(token string, assetId string) (*[]NftTransfer, error) {
	nftTransfers, err := RequestToGetNftTransferByAssetId(token, assetId)
	if err != nil {
		return nil, err
	}
	return nftTransfers, nil
}

func NftTransferToNftTransferSimplified(nftTransfer NftTransfer) NftTransferSimplified {
	return NftTransferSimplified{
		ID:       nftTransfer.ID,
		Txid:     nftTransfer.Txid,
		AssetId:  nftTransfer.AssetId,
		Time:     nftTransfer.Time,
		FromAddr: nftTransfer.FromAddr,
		ToAddr:   nftTransfer.ToAddr,
	}
}

func NftTransferSliceToNftTransferSimplifiedSlice(nftTransfers *[]NftTransfer) *[]NftTransferSimplified {
	if nftTransfers == nil {
		return nil
	}
	var NftTransferSimplifiedSlice []NftTransferSimplified
	for _, nftTransfer := range *nftTransfers {
		NftTransferSimplifiedSlice = append(NftTransferSimplifiedSlice, NftTransferToNftTransferSimplified(nftTransfer))
	}
	return &NftTransferSimplifiedSlice
}

func GetReceiveAddrByAssetId(assetId string) (string, error) {
	addrEvents, err := AddrReceivesAndGetEvents("")
	if err != nil {
		return "", err
	}
	for _, event := range *addrEvents {
		if event.Addr.AssetID == assetId {
			return event.Addr.Encoded, nil
		}
	}
	return "", errors.New("not found match addr by asset id")
}

func SetNftTransferWithoutInfo(token string, txid string, assetId string, time int, fromAddr string, toAddr string, deviceId string) error {
	_, err := SetNftTransferAndGetResponse(token, txid, assetId, time, fromAddr, toAddr, "", "", deviceId)
	if err != nil {
		return err
	}
	return nil
}

func UploadNftTransfer(token string, deviceId string, txid string, assetId string, _time int, fromAddr string, toAddr string) error {
	return SetNftTransferWithoutInfo(token, txid, assetId, _time, fromAddr, toAddr, deviceId)
}

func GetNftTransferByAssetId(token string, assetId string) string {
	nftTransfers, err := GetNftTransferByAssetIdAndGetResponse(token, assetId)
	result := NftTransferSliceToNftTransferSimplifiedSlice(nftTransfers)
	if err != nil {
		return MakeJsonErrorResult(GetNftTransferByAssetIdAndGetResponseErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, result)
}

type AssetListInfo struct {
	Version      string `json:"version"`
	GenesisPoint string `json:"genesis_point"`
	Name         string `json:"name"`
	MetaHash     string `json:"meta_hash"`
	AssetID      string `json:"asset_id"`
	AssetType    string `json:"asset_type"`
	OutputIndex  int    `json:"output_index"`

	Amount           int    `json:"amount"`
	LockTime         int32  `json:"lock_time"`
	RelativeLockTime int32  `json:"relative_lock_time"`
	ScriptKey        string `json:"script_key"`

	AnchorOutpoint string `json:"anchor_outpoint"`

	TweakedGroupKey string `json:"tweaked_group_key"`

	DeviceId string `json:"device_id" gorm:"type:varchar(255)"`
	UserId   int    `json:"user_id"`
	Username string `json:"username" gorm:"type:varchar(255)"`
}

type AssetListSetRequest struct {
	Version          string `json:"version" gorm:"type:varchar(255);index"`
	GenesisPoint     string `json:"genesis_point" gorm:"type:varchar(255)"`
	Name             string `json:"name" gorm:"type:varchar(255);index"`
	MetaHash         string `json:"meta_hash" gorm:"type:varchar(255);index"`
	AssetID          string `json:"asset_id" gorm:"type:varchar(255);index"`
	AssetType        string `json:"asset_type" gorm:"type:varchar(255);index"`
	OutputIndex      int    `json:"output_index"`
	Amount           int    `json:"amount"`
	LockTime         int32  `json:"lock_time"`
	RelativeLockTime int32  `json:"relative_lock_time"`
	ScriptKey        string `json:"script_key" gorm:"type:varchar(255);index"`
	AnchorOutpoint   string `json:"anchor_outpoint" gorm:"type:varchar(255);index"`
	TweakedGroupKey  string `json:"tweaked_group_key" gorm:"type:varchar(255);index"`
	DeviceId         string `json:"device_id" gorm:"type:varchar(255);index"`
}

func UploadAssetListInfo(token string, deviceId string) string {
	return UploadAssetListProcessedInfo(token, deviceId)
}

func UploadAssetListProcessedInfo(token string, deviceId string) string {
	isTokenValid, err := IsTokenValid(token)
	if err != nil {
		return MakeJsonErrorResult(IsTokenValidErr, "server "+err.Error()+"; token is invalid, did not send.", nil)
	} else if !isTokenValid {
		return MakeJsonErrorResult(IsTokenValidErr, "token is invalid, did not send.", nil)
	}
	assets, err := ListAssetsProcessed(true, false, false)
	if err != nil {
		return MakeJsonErrorResult(ListAssetsProcessedErr, err.Error(), nil)
	}
	zeroAmountAssetLists, err := GetZeroAmountAssetListSlice(token, assets)
	if err != nil {
		return MakeJsonErrorResult(GetZeroAmountAssetListSliceErr, err.Error(), nil)
	}
	listAssetsResponseSlice := AssetBalanceInfosToListAssetsResponseSlice(zeroAmountAssetLists)
	setListAssetsResponseSlice := append(*assets, *listAssetsResponseSlice...)
	requests := ListAssetsResponseSliceToAssetListSetRequests(&setListAssetsResponseSlice, deviceId)
	result, err := PostToSetAssetListInfo(requests, token)
	if err != nil {
		return MakeJsonErrorResult(PostToSetAssetListInfoErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, result.Data)
}

type AssetBalanceHistorySetRequest struct {
	AssetId string `json:"asset_id" gorm:"type:varchar(255);index"`
	Balance int    `json:"balance" gorm:"index"`
}

type AssetBalanceHistoryRecord struct {
	ID       uint   `json:"id" gorm:"primarykey"`
	AssetId  string `json:"asset_id" gorm:"type:varchar(255);index"`
	Balance  int    `json:"balance" gorm:"index"`
	Username string `json:"username" gorm:"type:varchar(255);index"`
}

func PostToCreateAssetBalanceHistories(token string, requests *[]AssetBalanceHistorySetRequest) (*JsonResult, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_balance_history/create"
	requestJsonBytes, err := json.Marshal(requests)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response JsonResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return &response, nil
}

type GetLatestAssetBalanceHistoriesResponse struct {
	Success bool                         `json:"success"`
	Error   string                       `json:"error"`
	Code    ErrCode                      `json:"code"`
	Data    *[]AssetBalanceHistoryRecord `json:"data"`
}

func RequestToGetLatestAssetBalanceHistories(token string) (*[]AssetBalanceHistoryRecord, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := serverDomainOrSocket + "/asset_balance_history/get/latest"
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response GetLatestAssetBalanceHistoriesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

func GetAndUploadAssetBalanceHistories(token string) error {
	records, err := RequestToGetLatestAssetBalanceHistories(token)
	if err != nil {
		return AppendErrorInfo(err, "RequestToGetLatestAssetBalanceHistories")
	}
	recordsMapBalance := make(map[string]int)
	if records != nil {
		for _, record := range *records {
			recordsMapBalance[record.AssetId] = record.Balance
		}
	}
	balances, err := GetListBalancesSimpleInfo()
	if err != nil {
		return AppendErrorInfo(err, "GetListBalancesSimpleInfo")
	}
	var changedBalances []ListBalanceSimpleInfo
	for _, balance := range *balances {
		mapBalance, ok := recordsMapBalance[balance.AssetID]
		if !ok {
			changedBalances = append(changedBalances, balance)
		} else if balance.Balance != mapBalance {
			changedBalances = append(changedBalances, balance)
		}
	}
	if changedBalances == nil {
		return nil
	}
	var requests []AssetBalanceHistorySetRequest
	for _, changedBalance := range changedBalances {
		requests = append(requests, AssetBalanceHistorySetRequest{
			AssetId: changedBalance.AssetID,
			Balance: changedBalance.Balance,
		})
	}
	if requests == nil || len(requests) == 0 {
		return nil
	}
	_, err = PostToCreateAssetBalanceHistories(token, &requests)
	if err != nil {
		return AppendErrorInfo(err, "PostToCreateAssetBalanceHistories")
	}
	return nil
}

func UploadAssetBalanceHistories(token string) string {
	err := GetAndUploadAssetBalanceHistories(token)
	if err != nil {
		return MakeJsonErrorResult(GetAndUploadAssetBalanceHistoriesErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, nil)
}

type AssetKeys struct {
	OpStr          string `json:"op_str"`
	ScriptKeyBytes string `json:"script_key_bytes"`
}

func _assetLeafKeys(isGroup bool, id string, proofType universerpc.ProofType) (*universerpc.AssetLeafKeyResponse, error) {

	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()
	client := universerpc.NewUniverseClient(conn)

	request := &universerpc.AssetLeafKeysRequest{
		Id: &universerpc.ID{
			ProofType: proofType,
		},
	}
	if isGroup {
		groupKey := &universerpc.ID_GroupKeyStr{
			GroupKeyStr: id,
		}
		request.Id.Id = groupKey
	} else {
		AssetId := &universerpc.ID_AssetIdStr{
			AssetIdStr: id,
		}
		request.Id.Id = AssetId
	}

	response, err := client.AssetLeafKeys(context.Background(), request)
	if err != nil {
		return nil, AppendErrorInfo(err, "AssetLeafKeys")
	}
	return response, nil
}

func AssetLeafKeyResponseToAssetKeys(response *universerpc.AssetLeafKeyResponse) *[]AssetKeys {
	if response == nil {
		return nil
	}
	var assetKeys []AssetKeys
	for _, key := range response.AssetKeys {
		assetKeys = append(assetKeys, AssetKeys{
			OpStr:          key.Outpoint.(*universerpc.AssetKey_OpStr).OpStr,
			ScriptKeyBytes: hex.EncodeToString(key.GetScriptKeyBytes()),
		})
	}
	return &assetKeys
}

func _queryProof(isGroup bool, id string, outpoint string, scriptKey string, proofType universerpc.ProofType) (*universerpc.AssetProofResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()
	client := universerpc.NewUniverseClient(conn)

	request := &universerpc.UniverseKey{
		Id: &universerpc.ID{
			ProofType: proofType,
		},
		LeafKey: &universerpc.AssetKey{
			Outpoint:  &universerpc.AssetKey_OpStr{OpStr: outpoint},
			ScriptKey: &universerpc.AssetKey_ScriptKeyStr{ScriptKeyStr: scriptKey},
		},
	}
	if isGroup {
		groupKey := &universerpc.ID_GroupKeyStr{
			GroupKeyStr: id,
		}
		request.Id.Id = groupKey
	} else {
		AssetId := &universerpc.ID_AssetIdStr{
			AssetIdStr: id,
		}
		request.Id.Id = AssetId
	}
	response, err := client.QueryProof(context.Background(), request)
	if err != nil {
		return nil, AppendErrorInfo(err, "QueryProof")
	}
	return response, nil
}

func QueryProofToGetAssetId(groupKey string, outpoint string, scriptKey string) (string, error) {
	response, err := _queryProof(true, groupKey, outpoint, scriptKey, universerpc.ProofType_PROOF_TYPE_ISSUANCE)
	if err != nil {
		return "", err
	}
	assetId := hex.EncodeToString(response.AssetLeaf.Asset.AssetGenesis.AssetId)
	return assetId, nil
}

type AssetMeta struct {
	Data     string `json:"data"`
	Type     string `json:"type"`
	MetaHash string `json:"meta_hash"`
}

func _fetchAssetMetaByAssetId(assetId string) (*taprpc.AssetMeta, error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()
	client := taprpc.NewTaprootAssetsClient(conn)
	request := &taprpc.FetchAssetMetaRequest{
		Asset: &taprpc.FetchAssetMetaRequest_AssetIdStr{
			AssetIdStr: assetId,
		},
	}
	response, err := client.FetchAssetMeta(context.Background(), request)
	return response, err
}

func FetchAssetMetaByAssetId(assetId string) (*AssetMeta, error) {
	response, err := _fetchAssetMetaByAssetId(assetId)
	if err != nil {
		return nil, err
	}
	assetMeta := AssetMeta{
		Data:     string(response.Data),
		Type:     response.Type.String(),
		MetaHash: hex.EncodeToString(response.MetaHash),
	}
	return &assetMeta, nil
}

func GetGroupNamesByGroupKeys(groupKeys []string) (*map[string]string, error) {
	var totalOutpoints []string
	groupKeyMapName := make(map[string]string)
	groupKeyMapOps := make(map[string][]string)
	opMapScriptKey := make(map[string]string)
	for _, groupKey := range groupKeys {
		assetKeys, err := func(isGroup bool, id string, proofType universerpc.ProofType) (*[]AssetKeys, error) {
			response, err := _assetLeafKeys(isGroup, id, proofType)
			if err != nil {
				return nil, err
			}
			var assetKeys *[]AssetKeys
			assetKeys = AssetLeafKeyResponseToAssetKeys(response)
			return assetKeys, nil
		}(true, groupKey, universerpc.ProofType_PROOF_TYPE_ISSUANCE)
		if err != nil {
			LogError("api AssetLeafKeys err:%v", err)
		}
		if len(*assetKeys) == 0 {
			err = errors.New("length of assetKeys(" + strconv.Itoa(len(*assetKeys)) + ") is zero, not fount AssetLeafKey")
			if err != nil {
				LogError("%v", err)
			}
		}
		var outpoints []string
		for _, assetKey := range *assetKeys {
			outpoints = append(outpoints, assetKey.OpStr)
			opMapScriptKey[assetKey.OpStr] = assetKey.ScriptKeyBytes
		}
		totalOutpoints = append(totalOutpoints, outpoints...)
		groupKeyMapOps[groupKey] = outpoints
	}
	type timeAndAssetKey struct {
		Time           int    `json:"time"`
		OpStr          string `json:"op_str"`
		ScriptKeyBytes string `json:"script_key_bytes"`
	}
	for _, groupKey := range groupKeys {
		var timeAndAssetKeys []timeAndAssetKey
		ops := groupKeyMapOps[groupKey]
		for _, op := range ops {
			timeAndAssetKeys = append(timeAndAssetKeys, timeAndAssetKey{
				Time:           0,
				OpStr:          op,
				ScriptKeyBytes: opMapScriptKey[op],
			})
		}
		firstAssetKey := timeAndAssetKeys[0]
		assetId, err := QueryProofToGetAssetId(groupKey, firstAssetKey.OpStr, firstAssetKey.ScriptKeyBytes)
		if err != nil {
			LogError("api QueryProofToGetAssetId err:%v", err)
			continue
		}
		assetMeta, err := FetchAssetMetaByAssetId(assetId)
		if err != nil {
			LogError("api FetchAssetMetaByAssetId err:%v", err)
			continue
		}
		var meta Meta
		meta.GetMetaFromStr(assetMeta.Data)
		groupKeyMapName[groupKey] = meta.GroupName
	}
	return &groupKeyMapName, nil
}

func queryListAssetsByAssetId(assetId string) (assets []ListAssetsResponse, err error) {
	listAssetsProcessed, err := ListAssetsProcessed(false, false, false)
	if err != nil {
		return assets, AppendErrorInfo(err, "ListAssetsProcessed")
	}
	for _, asset := range *listAssetsProcessed {
		if asset.AssetGenesis.AssetID == assetId && asset.ScriptKeyIsLocal {
			assets = append(assets, asset)
		}
	}
	return assets, nil
}

type ListAssetAmountInfo struct {
	Version      string                         `json:"version"`
	AssetGenesis ListAssetsResponseAssetGenesis `json:"asset_genesis"`
	Amount       int                            `json:"amount"`
	AssetGroup   ListAssetsResponseAssetGroup   `json:"asset_group"`
}

func ListAssetsResponseSliceToListAssetAmountInfo(assets []ListAssetsResponse) (listAssetAmountInfo ListAssetAmountInfo) {
	for i, asset := range assets {
		if i == 0 {
			listAssetAmountInfo = ListAssetAmountInfo{
				Version:      asset.Version,
				AssetGenesis: asset.AssetGenesis,
				AssetGroup:   asset.AssetGroup,
			}
		}
		listAssetAmountInfo.Amount += asset.Amount
	}
	return listAssetAmountInfo
}

func QueryListAssetsByAssetId(assetId string) string {
	assets, err := queryListAssetsByAssetId(assetId)
	if err != nil {
		return MakeJsonErrorResult(QueryListAssetsByAssetIdErr, err.Error(), ListAssetAmountInfo{})
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, ListAssetsResponseSliceToListAssetAmountInfo(assets))
}
