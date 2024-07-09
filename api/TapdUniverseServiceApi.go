package api

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/wire"
	"github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	"github.com/wallet/base"
	"github.com/wallet/service/apiConnect"
	rpcclient2 "github.com/wallet/service/rpcclient"
)

func AddFederationServer() {}

func assetLeafKeys(id string, proofType universerpc.ProofType) (*universerpc.AssetLeafKeyResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
	}
	defer clearUp()
	client := universerpc.NewUniverseClient(conn)
	request := &universerpc.AssetLeafKeysRequest{
		Id: &universerpc.ID{
			Id: &universerpc.ID_AssetIdStr{
				AssetIdStr: id,
			},
			ProofType: proofType,
		},
		//Offset:    0,
		//Limit:     0,
		//Direction: 0,
	}
	response, err := client.AssetLeafKeys(context.Background(), request)
	if err != nil {
		fmt.Printf("%s universerpc Info Error: %v\n", GetTimeNow(), err)
		return nil, err
	}
	return response, nil
}

func AssetLeafKeysAndGetResponse(assetId string, proofType universerpc.ProofType) (*universerpc.AssetLeafKeyResponse, error) {
	return assetLeafKeys(assetId, proofType)
}

func AssetLeafKeys(id string, proofType string) string {
	var _proofType universerpc.ProofType
	if proofType == "issuance" || proofType == "ISSUANCE" || proofType == "PROOF_TYPE_ISSUANCE" {
		_proofType = universerpc.ProofType_PROOF_TYPE_ISSUANCE
	} else if proofType == "transfer" || proofType == "TRANSFER" || proofType == "PROOF_TYPE_TRANSFER" {
		_proofType = universerpc.ProofType_PROOF_TYPE_TRANSFER
	} else {
		_proofType = universerpc.ProofType_PROOF_TYPE_UNSPECIFIED
	}
	response, err := assetLeafKeys(id, _proofType)
	if err != nil {
		fmt.Printf("%s universerpc AssetLeafKeys Error: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	if len(response.AssetKeys) == 0 {
		return MakeJsonErrorResult(DefaultErr, "Result length is zero.", nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", processAssetKey(response))
}

type AssetKey struct {
	OpStr          string `json:"op_str"`
	ScriptKeyBytes string `json:"script_key_bytes"`
}

func processAssetKey(response *universerpc.AssetLeafKeyResponse) *[]AssetKey {
	var assetKey []AssetKey
	for _, keys := range response.AssetKeys {
		assetKey = append(assetKey, AssetKey{
			OpStr:          keys.GetOpStr(),
			ScriptKeyBytes: hex.EncodeToString(keys.GetScriptKeyBytes()),
		})
	}
	return &assetKey
}

func AssetLeaves(id string) string {
	response, err := assetLeaves(false, id, universerpc.ProofType_PROOF_TYPE_ISSUANCE)
	if err != nil {
		fmt.Printf("%s universerpc AssetLeaves Error: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}

	if response.Leaves == nil {
		return MakeJsonErrorResult(DefaultErr, "NOT_FOUND", nil)
	}

	return MakeJsonErrorResult(SUCCESS, "", response)
}

func GetAssetInfo(id string) string {
	response, err := assetLeaves(false, id, universerpc.ProofType_PROOF_TYPE_ISSUANCE)
	if err != nil {
		fmt.Printf("%s universerpc AssetLeaves Error: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	if response.Leaves == nil {
		return MakeJsonErrorResult(DefaultErr, "NOT_FOUND", nil)
	}

	proof, err := rpcclient2.DecodeProof(response.Leaves[0].Proof, 0, false, false)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	block, err := rpcclient2.GetBlock(proof.DecodedProof.Asset.ChainAnchor.AnchorBlockHash)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}

	msgBlock := &wire.MsgBlock{}
	blockReader := bytes.NewReader(block.RawBlock)
	err = msgBlock.Deserialize(blockReader)
	timeStamp := msgBlock.Header.Timestamp
	createTime := timeStamp.Unix()
	createHeight := proof.DecodedProof.Asset.ChainAnchor.BlockHeight

	assetId := hex.EncodeToString(proof.DecodedProof.Asset.AssetGenesis.GetAssetId())
	assetType := proof.DecodedProof.Asset.AssetGenesis.AssetType.String()
	assetPoint := proof.DecodedProof.Asset.AssetGenesis.GenesisPoint
	amount := proof.DecodedProof.Asset.Amount
	assetName := proof.DecodedProof.Asset.AssetGenesis.Name
	fmt.Println(proof)
	var newMeta Meta
	newMeta.FetchAssetMeta(false, id)

	var assetInfo = struct {
		AssetId      string  `json:"asset_Id"`
		Name         string  `json:"name"`
		Point        string  `json:"point"`
		AssetType    string  `json:"assetType"`
		GroupName    *string `json:"group_name"`
		GroupKey     *string `json:"group_key"`
		Amount       uint64  `json:"amount"`
		Meta         *string `json:"meta"`
		CreateHeight int64   `json:"create_height"`
		CreateTime   int64   `json:"create_time"`
		Universe     string  `json:"universe"`
	}{
		AssetId:      assetId,
		Name:         assetName,
		Point:        assetPoint,
		AssetType:    assetType,
		GroupName:    &newMeta.GroupName,
		Amount:       amount,
		Meta:         &newMeta.Description,
		CreateHeight: int64(createHeight),
		CreateTime:   createTime,
		Universe:     "localhost",
	}
	if proof.DecodedProof.Asset.AssetGroup != nil {
		groupKey := hex.EncodeToString(proof.DecodedProof.Asset.AssetGroup.RawGroupKey)
		assetInfo.GroupKey = &groupKey
	}

	return MakeJsonErrorResult(SUCCESS, "", assetInfo)
}

func AssetRoots() {}

func DeleteAssetRoot() {}

func DeleteFederationServer() {}

// UniverseInfo
//
//	@Description: Info returns a set of information about the current state of the Universe.
//	@return string
func UniverseInfo() string {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
	}
	defer clearUp()

	client := universerpc.NewUniverseClient(conn)
	request := &universerpc.InfoRequest{}
	response, err := client.Info(context.Background(), request)
	if err != nil {
		fmt.Printf("%s universerpc Info Error: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func InsertProof() {}

// ListFederationServers
//
//	@Description: ListFederationServers lists the set of servers that make up the federation of the local Universe server.
//	This servers are used to push out new proofs, and also periodically call sync new proofs from the remote server.
//	@return string
func ListFederationServers() string {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
	}
	defer clearUp()
	client := universerpc.NewUniverseClient(conn)
	request := &universerpc.ListFederationServersRequest{}
	response, err := client.ListFederationServers(context.Background(), request)
	if err != nil {
		fmt.Printf("%s universerpc ListFederationServers Error: %v\n", GetTimeNow(), err)
		return ""
	}
	return response.String()
}

func MultiverseRoot() {}

func QueryAssetRoots(id string) string {
	response, err := queryAssetRoot(id)
	if err != nil {
		fmt.Printf("%s universerpc AssetRoots Error: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func QueryAssetStats(assetId string) string {
	response, err := queryAssetStats(assetId)
	if err != nil {
		fmt.Printf("%s universerpc QueryAssetStats Error: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(DefaultErr, err.Error(), "")
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func QueryEvents() {}

func QueryFederationSyncConfig() {}

func QueryProof() {}

func SetFederationSyncConfig() {}

func SyncUniverse(universeHost string, assetId string) string {
	var targets []*universerpc.SyncTarget
	universeID := &universerpc.ID{
		Id: &universerpc.ID_AssetIdStr{
			AssetIdStr: assetId,
		},
		ProofType: universerpc.ProofType_PROOF_TYPE_ISSUANCE,
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
	response, err := syncUniverse(universeHost, targets, universerpc.UniverseSyncMode_SYNC_FULL)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), "")
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func UniverseStats() {}

func queryAssetRoot(id string) (*universerpc.QueryRootResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
	}
	defer clearUp()

	requst := &universerpc.AssetRootQuery{
		Id: &universerpc.ID{
			Id: &universerpc.ID_AssetIdStr{
				AssetIdStr: id,
			},
		},
	}
	client := universerpc.NewUniverseClient(conn)
	response, err := client.QueryAssetRoots(context.Background(), requst)
	return response, err
}

func assetLeaves(isGroup bool, id string, proofType universerpc.ProofType) (*universerpc.AssetLeafResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
	}
	defer clearUp()
	request := &universerpc.ID{
		ProofType: proofType,
	}

	if isGroup {
		groupKey := &universerpc.ID_GroupKeyStr{
			GroupKeyStr: id,
		}
		request.Id = groupKey
	} else {
		AssetId := &universerpc.ID_AssetIdStr{
			AssetIdStr: id,
		}
		request.Id = AssetId
	}

	client := universerpc.NewUniverseClient(conn)
	response, err := client.AssetLeaves(context.Background(), request)
	return response, err
}

func AssetLeavesAndGetResponse(isGroup bool, id string, proofType universerpc.ProofType) (*universerpc.AssetLeafResponse, error) {
	return assetLeaves(isGroup, id, proofType)
}

func queryAssetStats(assetId string) (*universerpc.UniverseAssetStats, error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
	}
	defer clearUp()
	id, err := hex.DecodeString(assetId)
	client := universerpc.NewUniverseClient(conn)
	request := &universerpc.AssetStatsQuery{
		AssetIdFilter: id,
	}
	response, err := client.QueryAssetStats(context.Background(), request)
	return response, err
}

func syncUniverse(universeHost string, syncTargets []*universerpc.SyncTarget, syncMode universerpc.UniverseSyncMode) (*universerpc.SyncResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
	}
	defer clearUp()
	request := &universerpc.SyncRequest{
		UniverseHost: universeHost,
		SyncMode:     syncMode,
		SyncTargets:  syncTargets,
	}
	client := universerpc.NewUniverseClient(conn)
	response, err := client.SyncUniverse(context.Background(), request)
	return response, err
}
