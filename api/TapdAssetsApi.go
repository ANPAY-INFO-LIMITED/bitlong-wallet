package api

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/lightninglabs/taproot-assets/taprpc/mintrpc"
	"github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/wallet/api/connect"
	"github.com/wallet/api/rpcclient"
	"github.com/wallet/base"
	"io"
	"net/http"
	"strconv"
	"strings"
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
	//ScriptKey   string `json:"script_key"`
}

type AssetsTransfersOutputAnchor struct {
	Outpoint string `json:"outpoint"`
	Value    int    `json:"value"`
	//TaprootAssetRoot string `json:"taproot_asset_root"`
	//MerkleRoot       string `json:"merkle_root"`
	//TapscriptSibling string `json:"tapscript_sibling"`
	//NumPassiveAssets int    `json:"num_passive_assets"`
}

type AssetsTransfersOutput struct {
	Anchor           AssetsTransfersOutputAnchor
	ScriptKeyIsLocal bool `json:"script_key_is_local"`
	Amount           int  `json:"amount"`
	//SplitCommitRootHash string `json:"split_commit_root_hash"`
	OutputType   string `json:"output_type"`
	AssetVersion string `json:"asset_version"`
}

// @dev: May be deprecated
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
				//ScriptKey:   hex.EncodeToString(_input.ScriptKey),
			})
		}
		var outputs []AssetsTransfersOutput
		for _, _output := range transfers.Outputs {
			outputs = append(outputs, AssetsTransfersOutput{
				Anchor: AssetsTransfersOutputAnchor{
					Outpoint: _output.Anchor.Outpoint,
					Value:    int(_output.Anchor.Value),
					//TaprootAssetRoot: hex.EncodeToString(_output.anchor.TaprootAssetRoot),
					//MerkleRoot:       hex.EncodeToString(_output.anchor.MerkleRoot),
					//TapscriptSibling: hex.EncodeToString(_output.anchor.TapscriptSibling),
					//NumPassiveAssets: int(_output.anchor.NumPassiveAssets),
				},
				ScriptKeyIsLocal: _output.ScriptKeyIsLocal,
				Amount:           int(_output.Amount),
				//SplitCommitRootHash: hex.EncodeToString(_output.SplitCommitRootHash),
				OutputType:   _output.OutputType.String(),
				AssetVersion: _output.AssetVersion.String(),
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
	Version      string                 `json:"version"`
	AssetGenesis AssetsListAssetGenesis `json:"asset_genesis"`
	Amount       int                    `json:"amount"`
	LockTime     int                    `json:"lock_time"`
	//RelativeLockTime int    `json:"relative_lock_time"`
	//ScriptVersion    int    `json:"script_version"`
	//ScriptKey        string `json:"script_key"`
	ScriptKeyIsLocal bool `json:"script_key_is_local"`
	//RawGroupKey      string `json:"raw_group_key"`
	//AssetGroup       struct {
	//	RawGroupKey     string `json:"raw_group_key"`
	//	TweakedGroupKey string `json:"tweaked_group_key"`
	//	AssetWitness    string `json:"asset_witness"`
	//} `json:"asset_group"`
	ChainAnchor AssetsListChainAnchor `json:"chain_anchor"`
	//PrevWitnesses []interface{} `json:"prev_witnesses"`
	IsSpent     bool   `json:"is_spent"`
	LeaseOwner  string `json:"lease_owner"`
	LeaseExpiry int    `json:"lease_expiry"`
	IsBurn      bool   `json:"is_burn"`
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

// @dev: May be deprecated
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
				Version:      int(_asset.AssetGenesis.Version),
			},
			Amount:           int(_asset.Amount),
			LockTime:         int(_asset.LockTime),
			ScriptKeyIsLocal: _asset.ScriptKeyIsLocal,
			//RawGroupKey:      hex.EncodeToString(_asset.AssetGroup.RawGroupKey),
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

// SyncUniverseFullSpecified @dev
func SyncUniverseFullSpecified(universeHost string, id string, proofType string) string {
	if universeHost == "" {
		switch base.NetWork {
		case base.UseTestNet:
			universeHost = "testnet.universe.lightning.finance:10029"
		case base.UseMainNet:
			universeHost = "mainnet.universe.lightning.finance:10029"
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
		return MakeJsonErrorResult(DefaultErr, err.Error(), "")
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

// SyncAssetIssuance @dev
func SyncAssetIssuance(id string) string {
	return SyncUniverseFullSpecified("", id, universerpc.ProofType_PROOF_TYPE_ISSUANCE.String())
}

// SyncAssetTransfer @dev
func SyncAssetTransfer(id string) string {
	return SyncUniverseFullSpecified("", id, universerpc.ProofType_PROOF_TYPE_TRANSFER.String())
}

// SyncAssetAll @dev
func SyncAssetAll(id string) {
	fmt.Println(SyncAssetIssuance(id))
	fmt.Println(SyncAssetTransfer(id))
}

// SyncAssetAllSlice
// @dev
func SyncAssetAllSlice(ids []string) {
	if len(ids) == 0 {
		return
	}
	for _, _id := range ids {
		fmt.Println("Sync issuance:", _id, ".", SyncAssetIssuance(_id))
		fmt.Println("Sync transfer:", _id, ".", SyncAssetTransfer(_id))
	}
}

// SyncAssetAllWithAssets @dev
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

// GetAllAssetBalances
// @note: Get all balance of assets info
func GetAllAssetBalances() string {
	result := allAssetBalances()
	if result == nil {
		return MakeJsonErrorResult(DefaultErr, "Null Balances", nil)
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
		return MakeJsonErrorResult(DefaultErr, "Null Group Balances", nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", result)
}

// @dev: May be deprecated
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

// SyncAllAssetsByAssetBalance
// @note: Sync all assets of non-zero-balance to public universe
// @dev: May be deprecated
func SyncAllAssetsByAssetBalance() string {
	ids := GetAllAssetIdByAssetBalance(allAssetBalances())
	if ids != nil {
		SyncAssetAllSlice(*ids)
	}
	return MakeJsonErrorResult(SUCCESS, "", ids)
}

// GetAllAssetsIdSlice
// @dev: 3
// @note: Get an array including all assets ids
// @dev: May be deprecated
func GetAllAssetsIdSlice() string {
	ids := GetAllAssetIdByAssetBalance(allAssetBalances())
	return MakeJsonErrorResult(SUCCESS, "", ids)
}

// assetKeysTransfer
// @dev
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
		return MakeJsonErrorResult(DefaultErr, "Null Asset Keys", nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", result)
}

// AssetLeavesSpecified
// @dev: Need To Complete
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
	Name string `json:"name"`
	//MetaHash     string `json:"meta_hash"`
	AssetID   string `json:"asset_id"`
	Amount    int    `json:"amount"`
	ScriptKey string `json:"script_key"`
	//PrevWitnesses []struct {
	//	PrevID struct {
	//		AnchorPoint string `json:"anchor_point"`
	//		AssetID     string `json:"asset_id"`
	//		ScriptKey   string `json:"script_key"`
	//	} `json:"prev_id"`
	//	SplitCommitment struct {
	//		RootAsset struct {
	//			AssetGenesis struct {
	//				GenesisPoint string `json:"genesis_point"`
	//				Name         string `json:"name"`
	//				MetaHash     string `json:"meta_hash"`
	//				AssetID      string `json:"asset_id"`
	//			} `json:"asset_genesis"`
	//			Amount        int    `json:"amount"`
	//			ScriptKey     string `json:"script_key"`
	//			PrevWitnesses []struct {
	//				PrevID struct {
	//					AnchorPoint string `json:"anchor_point"`
	//					AssetID     string `json:"asset_id"`
	//					ScriptKey   string `json:"script_key"`
	//				} `json:"prev_id"`
	//				TxWitness []string `json:"tx_witness"`
	//			} `json:"prev_witnesses"`
	//		} `json:"root_asset"`
	//	} `json:"split_commitment"`
	//} `json:"prev_witnesses"`
	Proof string `json:"proof"`
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
		fmt.Printf("%s universerpc AssetLeaves Error.\n", GetTimeNow())
		return MakeJsonErrorResult(DefaultErr, errors.New("null asset leaves").Error(), nil)
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

// @dev: Not-exported same copy of AssetLeavesTransfer_ONLY_FOR_TEST
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

// GetAssetInfoByIssuanceLeaf @dev
func GetAssetInfoByIssuanceLeaf(id string) string {
	response := assetLeavesIssuance(id)
	if response == nil {
		fmt.Printf("%s Universerpc asset leaves issuance error.\n", GetTimeNow())
		return MakeJsonErrorResult(DefaultErr, errors.New("Null asset leaves").Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func DecodeRawProofByte(rawProof []byte) *taprpc.DecodeProofResponse {
	result, err := rpcclient.DecodeProof(rawProof, 0, false, false)
	if err != nil {
		return nil
	}
	return result
}

// DecodeRawProofString
// @dev
func DecodeRawProofString(proof string) *taprpc.DecodeProofResponse {
	decodeString, err := hex.DecodeString(proof)
	if err != nil {
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
		return nil
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
		return MakeJsonErrorResult(DefaultErr, "null raw proof", nil)
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
		return MakeJsonErrorResult(DefaultErr, "null asset list", nil)
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

// GetAllAssetListSimplified
// @dev
func GetAllAssetListSimplified() string {
	result := ProcessListAllAssetsSimplified(ProcessListAllAssets(allAssetList()))
	if result == nil {
		return MakeJsonErrorResult(DefaultErr, "null asset list", nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", result)
}

func GetAllAssetIdByListAll() []string {
	id := make(map[string]bool)
	var ids []string
	result := ProcessListAllAssetsSimplified(ProcessListAllAssets(allAssetList()))
	//var index int
	if result == nil || len(*result) == 0 {
		return nil
	}
	for _, asset := range *result {
		//index++
		//fmt.Println(index, asset.GenesisAssetID)
		id[asset.GenesisAssetID] = true
	}
	for k, _ := range id {
		ids = append(ids, k)
	}
	if len(ids) == 0 {
		return nil
	}
	//fmt.Println(len(ids))
	return ids
}

// SyncUniverseFullIssuanceByIdSlice
// @dev
// @note: Deprecated
// @dev: May be deprecated
func SyncUniverseFullIssuanceByIdSlice(ids []string) string {
	var universeHost string
	switch base.NetWork {
	case base.UseTestNet:
		universeHost = "testnet.universe.lightning.finance:10029"
	case base.UseMainNet:
		universeHost = "mainnet.universe.lightning.finance:10029"
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
		return MakeJsonErrorResult(DefaultErr, err.Error(), "")
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

// SyncUniverseFullTransferByIdSlice
// @dev
// @note: Deprecated
// @dev: May be deprecated
func SyncUniverseFullTransferByIdSlice(ids []string) string {
	var universeHost string
	switch base.NetWork {
	case base.UseTestNet:
		universeHost = "testnet.universe.lightning.finance:10029"
	case base.UseMainNet:
		universeHost = "mainnet.universe.lightning.finance:10029"
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
		return MakeJsonErrorResult(DefaultErr, err.Error(), "")
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

// SyncUniverseFullNoSlice
// @dev
// @note: Sync all assets
func SyncUniverseFullNoSlice() string {
	var universeHost string
	switch base.NetWork {
	case base.UseTestNet:
		universeHost = "testnet.universe.lightning.finance:10029"
	case base.UseMainNet:
		universeHost = "mainnet.universe.lightning.finance:10029"
	}
	var targets []*universerpc.SyncTarget
	response, err := syncUniverse(universeHost, targets, universerpc.UniverseSyncMode_SYNC_FULL)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), "")
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
	//Proof     string `json:"proof"`
	IsSpent bool `json:"isSpent"`
}

// OutpointToAddress
// @dev
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

// getTransactionAndIndexByOutpoint
// @dev: Split outpoint
func getTransactionAndIndexByOutpoint(outpoint string) (transaction string, index string) {
	result := strings.Split(outpoint, ":")
	return result[0], result[1]
}

// CompareScriptKey
// @dev
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

// GetAssetHoldInfosIncludeSpent
// @dev
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
			//Proof:     leaf.Proof,
			IsSpent: AddressIsSpentAll(address),
		})
	}
	return &idToAssetHoldInfo
}

// GetAssetHoldInfosExcludeSpent
// @Description: This function uses multiple http requests to call mempool's api during processing,
// and it is recommended to store the data in a database and update it manually
// @dev: Get hold info of asset
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
		return MakeJsonErrorResult(DefaultErr, "Get asset hold infos include spent fail, null response.", nil)
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
		return nil
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
		//Name:              "",
		//AssetId:           "",
		//Amount:            0,
		BlockTime:       response.Status.BlockTime,
		FeeRate:         RoundToDecimalPlace(float64(response.Fee)/(float64(response.Weight)/4), 2),
		ConfirmedBlocks: BlockTipHeight() - response.Status.BlockHeight,
		//IsSpent:           false,
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

// SyncAllAssetByList
// @note: Call this api to sync all
func SyncAllAssetByList() string {
	SyncAssetAllSlice(GetAllAssetIdByListAll())
	return MakeJsonErrorResult(SUCCESS, "", "Sync Completed.")
}

// GetAssetInfoById
// @note: Call this api to get asset info
func GetAssetInfoById(id string) string {
	return GetAssetInfoByIssuanceLeaf(id)
}

// GetAssetHoldInfosExcludeSpentSlow
// @note: Call this api to get asset hold info which is not be spent
// @dev: Wrap to call GetAssetHoldInfosExcludeSpent
// @notice: THIS COST A LOT OF TIME
func GetAssetHoldInfosExcludeSpentSlow(id string) string {
	response := GetAssetHoldInfosExcludeSpent(id)
	if response == nil {
		return MakeJsonErrorResult(DefaultErr, "Get asset hold infos exclude spent fail, null response.", nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

// GetAssetTransactionInfoSlow
// @note: Call this api to get asset transaction info
// @notice: THIS COST A LOT OF TIME
func GetAssetTransactionInfoSlow(id string) string {
	response := GetAssetTransactionInfos(id)
	if response == nil {
		return MakeJsonErrorResult(DefaultErr, "Get asset transaction infos fail, null response.", nil)
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

// GetAllAssetListWithoutProcession
// ONLY_FOR_TEST
// @dev: Need to look for the change transaction anchored outpoint, amount, and is_spent in previous witness.
// @dev: Returns exclude spent
func GetAllAssetListWithoutProcession() string {
	response := allAssetList()
	if response == nil {
		return MakeJsonErrorResult(DefaultErr, "Null list asset response.", nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func ListBatchesAndGetResponse() (*mintrpc.ListBatchResponse, error) {
	conn, clearUp, err := connect.GetConnection("tapd", false)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
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
	conn, clearUp, err := connect.GetConnection("lnd", false)
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
	conn, clearUp, err := connect.GetConnection("lnd", false)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
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

//@dev

func ListBatchesAndGetCustomResponse() (*[]ListBatchesResponse, error) {
	response, err := ListBatchesAndGetResponse()
	if err != nil {
		LogError("", err)
		return nil, err
	}
	var listBatchesResponse []ListBatchesResponse
	for _, batch := range (*response).Batches {
		var assets []ListBatchesAsset
		for _, _asset := range batch.Assets {
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
			BatchKey:  hex.EncodeToString(batch.BatchKey),
			BatchTxid: batch.BatchTxid,
			State:     batch.State.String(),
			Assets:    assets,
			//Amount:    0,
			//NewGroupedAsset: false,
			//GroupKey:        "",
			//GroupAnchor:     "",
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
				Version:      int(_asset.AssetGenesis.Version),
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
		LogError("", err)
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

// DecodeTransactionsWhoseLabelIsNotTapdAssetMinting
// @dev: Call to decode transactions
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
	serverDomainOrSocket := "132.232.109.84:8090"
	network := base.NetWork
	url := "http://" + serverDomainOrSocket + "/bitcoind/" + network + "/decode/transactions"
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
		err := Body.Close()
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
	serverDomainOrSocket := "132.232.109.84:8090"
	network := base.NetWork
	url := "http://" + serverDomainOrSocket + "/bitcoind/" + network + "/decode/query/transactions"
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
		err := Body.Close()
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

func GetThenDecodeAndQueryTransactionsWhoseLabelIsNotTapdAssetMinting(token string) (*[]PostGetRawTransactionResultSat, error) {
	getTransactions, err := GetTransactionsWhoseLabelIsNotTapdAssetMinting()
	if err != nil {
		return nil, err
	}
	var rawTransactions []string
	for _, transaction := range *getTransactions {
		rawTransactions = append(rawTransactions, transaction.RawTxHex)
	}
	decodedAndQueryTransactions, err := DecodeAndQueryTransactionsWhoseLabelIsNotTapdAssetMinting(token, rawTransactions)
	if err != nil {
		return nil, err
	}
	btcUesult := ProcessDecodedAndQueryTransactionsData(decodedAndQueryTransactions.Data)
	result := ProcessPostGetRawTransactionResultToUseSat(btcUesult)
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
	transactions, err := GetThenDecodeAndQueryTransactionsWhoseLabelIsNotTapdAssetMinting(token)
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

func GetBtcTransferOutInfosJsonResult(token string) string {
	response, err := GetBtcTransferOutInfos(token)
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
	serverDomainOrSocket := "132.232.109.84:8090"
	url := "http://" + serverDomainOrSocket + "/asset_transfer/set"
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
		err := Body.Close()
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
		return &response, errors.New(response.Error)
	}
	return &response, nil
}

type PostToGetAssetTransferTxidsResponse struct {
	Success bool     `json:"success"`
	Error   string   `json:"error"`
	Code    ErrCode  `json:"code"`
	Data    []string `json:"data"`
}

func PostToGetAssetTransferTxids(token string) (txids []string, err error) {
	serverDomainOrSocket := "132.232.109.84:8090"
	url := "http://" + serverDomainOrSocket + "/asset_transfer/get/txids"
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
		err := Body.Close()
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

func PostCallBitcoindToQueryAddressByOutpoints(token string, outpoints []string) (*GetAddressesByOutpointSliceResponse, error) {
	serverDomainOrSocket := "132.232.109.84:8090"
	network := base.NetWork
	url := "http://" + serverDomainOrSocket + "/bitcoind/" + network + "/address/outpoints"
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
		err := Body.Close()
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

func ProcessListTransfersResponse(token string, listTransfersResponse *taprpc.ListTransfersResponse) *[]AssetTransferProcessed {
	var assetTransferProcessed []AssetTransferProcessed
	allOutpoints := GetAllOutPointsOfListTransfersResponse(listTransfersResponse)
	response, err := PostCallBitcoindToQueryAddressByOutpoints(token, allOutpoints)
	if err != nil {
		return nil
	}
	addressMap := response.Data
	for _, listTransfer := range listTransfersResponse.Transfers {
		txid, err := GetTxidFromOutpoint(listTransfer.Outputs[0].Anchor.Outpoint)
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
				OutputType:             output.OutputType.String(),
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
		})
	}
	return &assetTransferProcessed
}

func ListTransfersAndGetProcessedResponse(token string) (*[]AssetTransferProcessed, error) {
	listTransfers, err := ListTransfersAndGetResponse()
	if err != nil {
		return nil, err
	}
	processedListTransfers := ProcessListTransfersResponse(token, listTransfers)
	return processedListTransfers, nil
}

func ListAndPostToSetAssetTransfers(token string) string {
	transfers, err := ListTransfersAndGetProcessedResponse(token)
	if err != nil {
		return MakeJsonErrorResult(ListTransfersAndGetProcessedResponseErr, err.Error(), nil)
	}
	_, err = PostToSetAssetTransfer(token, transfers)
	if err != nil {
		return MakeJsonErrorResult(PostToSetAssetTransferErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", nil)
}

type GetAssetTransferResponse struct {
	// Deprecated: Use Code instead
	Success bool                      `json:"success"`
	Error   string                    `json:"error"`
	Code    ErrCode                   `json:"code"`
	Data    *[]AssetTransferProcessed `json:"data"`
}

func PostToGetAssetTransferAndGetResponse(token string) (*GetAssetTransferResponse, error) {
	serverDomainOrSocket := "132.232.109.84:8090"
	url := "http://" + serverDomainOrSocket + "/asset_transfer/get"
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
		err := Body.Close()
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
		return &response, errors.New(response.Error)
	}
	return &response, nil
}

func PostToGetAssetTransfer(token string) string {
	response, err := PostToGetAssetTransferAndGetResponse(token)
	if err != nil {
		return MakeJsonErrorResult(PostToGetAssetTransferAndGetResponseErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, response.Data)
}

// UploadAssetTransfer
// @Description: Upload assets transfer info
func UploadAssetTransfer(token string) string {
	return ListAndPostToSetAssetTransfers(token)
}

// GetAssetTransfer
// @Description: Get assets transfer info
func GetAssetTransfer(token string) string {
	return PostToGetAssetTransfer(token)
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
