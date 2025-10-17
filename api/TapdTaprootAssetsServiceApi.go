package api

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/lightninglabs/taproot-assets/tapfreighter"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	"github.com/pkg/errors"
	"github.com/wallet/service/apiConnect"
	"github.com/wallet/service/rpcclient"
	"sort"
	"strconv"
	"strings"
)

type AddrEvent struct {
	CreationTimeUnixSeconds int64        `json:"creation_time_unix_seconds"`
	Addr                    *QueriedAddr `json:"addr"`
	Status                  string       `json:"status"`
	Outpoint                string       `json:"outpoint"`
	Txid                    string       `json:"txid"`
	UtxoAmtSat              int64        `json:"utxo_amt_sat"`
	TaprootSibling          string       `json:"taproot_sibling"`
	ConfirmationHeight      int64        `json:"confirmation_height"`
	HasProof                bool         `json:"has_proof"`
}

func SortAddrEvents(addrEvents *[]AddrEvent) *[]AddrEvent {
	if addrEvents == nil {
		return nil
	}
	SortTimeDescInAssetTransfers := func(i, j int) bool {
		return (*addrEvents)[i].CreationTimeUnixSeconds > (*addrEvents)[j].CreationTimeUnixSeconds
	}
	sort.Slice(*addrEvents, SortTimeDescInAssetTransfers)
	return addrEvents
}

func AddrReceives(assetId string) string {
	response, err := rpcclient.AddrReceives()
	if err != nil {
		return MakeJsonErrorResult(AddrReceivesErr, err.Error(), nil)
	}
	var addrEvents []AddrEvent
	for _, event := range response.Events {
		if assetId != "" && assetId != hex.EncodeToString(event.Addr.AssetId) {
			continue
		}
		e := AddrEvent{}
		e.CreationTimeUnixSeconds = int64(event.CreationTimeUnixSeconds)
		a := QueriedAddr{}
		a.GetData(event.Addr)
		e.Addr = &a
		e.Status = event.Status.String()
		e.Outpoint = event.Outpoint
		e.Txid, _ = outpointToTransactionAndIndex(event.Outpoint)
		e.UtxoAmtSat = int64(event.UtxoAmtSat)
		e.TaprootSibling = hex.EncodeToString(event.TaprootSibling)
		e.ConfirmationHeight = int64(event.ConfirmationHeight)
		e.HasProof = event.HasProof
		addrEvents = append(addrEvents, e)
	}
	if len(addrEvents) == 0 {
		return MakeJsonErrorResult(SUCCESS, "NOT_FOUND", nil)
	}
	result := SortAddrEvents(&addrEvents)
	return MakeJsonErrorResult(SUCCESS, "", result)
}

func AddrReceivesOfAllNft() string {
	response, err := rpcclient.AddrReceives()
	if err != nil {
		return MakeJsonErrorResult(AddrReceivesErr, err.Error(), nil)
	}
	var addrEvents []AddrEvent
	for _, event := range response.Events {
		if event.Addr.AssetType != taprpc.AssetType_COLLECTIBLE {
			continue
		}
		a := QueriedAddr{}
		a.GetData(event.Addr)
		txid, _ := outpointToTransactionAndIndex(event.Outpoint)
		e := AddrEvent{
			CreationTimeUnixSeconds: int64(event.CreationTimeUnixSeconds),
			Addr:                    &a,
			Status:                  event.Status.String(),
			Outpoint:                event.Outpoint,
			Txid:                    txid,
			UtxoAmtSat:              int64(event.UtxoAmtSat),
			TaprootSibling:          hex.EncodeToString(event.TaprootSibling),
			ConfirmationHeight:      int64(event.ConfirmationHeight),
			HasProof:                event.HasProof,
		}
		addrEvents = append(addrEvents, e)
	}
	if len(addrEvents) == 0 {
		return MakeJsonErrorResult(SUCCESS, "NOT_FOUND", nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", addrEvents)
}

func BurnAsset(token string, AssetIdStr string, amountToBurn int64, deviceId string) string {
	response, err := rpcclient.BurnAsset(AssetIdStr, uint64(amountToBurn))
	if err != nil {
		return MakeJsonErrorResult(BurnAssetErr, err.Error(), nil)
	}
	err = UploadAssetBurn(token, AssetIdStr, int(amountToBurn), deviceId)
	if err != nil {
		LogError("Upload asset burn", err)
	}
	txHash := hex.EncodeToString(response.BurnTransfer.AnchorTxHash)
	return MakeJsonErrorResult(SUCCESS, "", txHash)
}

func DebugLevel() {

}

func DecodeAddr(addr string) string {
	response, err := rpcclient.DecodeAddr(addr)
	if err != nil {
		return MakeJsonErrorResult(DecodeAddrErr, err.Error(), nil)
	}
	result := QueriedAddr{}
	result.GetData(response)
	return MakeJsonErrorResult(SUCCESS, "", result)
}

func DecodeProof(rawProof string) {

}

func ExportProof() {

}

func FetchAssetMeta(isHash bool, data string) string {
	response, err := fetchAssetMeta(isHash, data)
	if err != nil {
		return MakeJsonErrorResult(fetchAssetMetaErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", string(response.Data))
}

func GetInfoOfTap() string {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return MakeJsonErrorResult(GetConnectionErr, err.Error(), nil)
	}
	defer clearUp()
	client := taprpc.NewTaprootAssetsClient(conn)
	request := &taprpc.GetInfoRequest{}
	response, err := client.GetInfo(context.Background(), request)
	if err != nil {
		return MakeJsonErrorResult(GetInfoErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func ListAssets(withWitness, includeSpent, includeLeased bool) string {
	response, err := listAssets(withWitness, includeSpent, includeLeased)
	if err != nil {
		return MakeJsonErrorResult(listAssetsErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func ListSimpleAssets(withWitness, includeSpent, includeLeased bool) string {
	response, err := listAssets(withWitness, includeSpent, includeLeased)
	if err != nil {
		return MakeJsonErrorResult(listAssetsErr, err.Error(), nil)
	}
	var (
		simpleAssets []struct {
			AssetId string `json:"asset_id"`
			Name    string `json:"name"`
			Amount  uint64 `json:"amount"`
			Type    string `json:"type"`
		}
	)
	for _, asset := range response.Assets {
		set := true
		for index, s := range simpleAssets {
			if s.AssetId == hex.EncodeToString(asset.AssetGenesis.GetAssetId()) {
				simpleAssets[index].Amount = asset.Amount + s.Amount
				set = false
				break
			}
		}
		if !set {
			continue
		}
		simpleAssets = append(simpleAssets, struct {
			AssetId string `json:"asset_id"`
			Name    string `json:"name"`
			Amount  uint64 `json:"amount"`
			Type    string `json:"type"`
		}{
			AssetId: hex.EncodeToString(asset.AssetGenesis.GetAssetId()),
			Name:    asset.AssetGenesis.Name,
			Amount:  asset.Amount,
			Type:    asset.AssetGenesis.AssetType.String(),
		})
	}

	return MakeJsonErrorResult(SUCCESS, "", simpleAssets)
}

func FindAssetByAssetName(assetName string) string {
	var response = struct {
		Success bool                     `json:"success"`
		Error   string                   `json:"error"`
		Data    taprpc.ListAssetResponse `json:"data"`
	}{}
	list := ListAssets(false, false, false)
	err := json.Unmarshal([]byte(list), &response)
	if err != nil {
		return MakeJsonErrorResult(UnmarshalErr, err.Error(), nil)
	}
	if response.Success == false {
		return MakeJsonErrorResult(responseNotSuccessErr, response.Error, nil)
	}
	var assets []*taprpc.Asset
	for _, asset := range response.Data.Assets {
		if asset.AssetGenesis.Name == assetName {
			assets = append(assets, asset)
		}
	}
	if len(assets) == 0 {
		return MakeJsonErrorResult(assetNotFoundErr, "NOT_FOUND", nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", assets)
}

func ListGroups() string {
	response, err := rpcclient.ListGroups()
	if err != nil {
		return MakeJsonErrorResult(ListGroupsErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func SortAssetTransferSimplified(assetTransfers *[]AssetTransferSimplified) *[]AssetTransferSimplified {
	if assetTransfers == nil {
		return nil
	}
	SortTimeDescInAssetTransfers := func(i, j int) bool {
		return (*assetTransfers)[i].Time > (*assetTransfers)[j].Time
	}
	sort.Slice(*assetTransfers, SortTimeDescInAssetTransfers)
	return assetTransfers
}

func QueryAssetTransfers(token string, assetId string) string {
	assetTransfers, err := QueryAssetTransferSimplified(token, assetId)
	if err != nil {
		return MakeJsonErrorResult(QueryAssetTransferSimplifiedErr, err.Error(), nil)
	}
	assetTransfers = SortAssetTransferSimplified(assetTransfers)
	return MakeJsonErrorResult(SUCCESS, SuccessError, assetTransfers)
}

func QueryAssetTransfersOfAllNft(token string) string {
	response, err := QueryAssetTransferSimplifiedOfAllNft(token)
	if err != nil {
		return MakeJsonErrorResult(QueryAssetTransferSimplifiedOfAllNftErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, response)
}

func queryAssetTransfers(assetId string) string {
	response, err := rpcclient.ListTransfers()
	if err != nil {
		return MakeJsonErrorResult(ListTransfersErr, err.Error(), nil)
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
		return MakeJsonErrorResult(SUCCESS, "NOT_FOUND", transfers)
	}
	return MakeJsonErrorResult(SUCCESS, "", transfers)
}

type ManagedUtxo struct {
	Op                 string             `json:"op"`
	OutPoint           string             `json:"out_point"`
	Time               int                `json:"time"`
	AmtSat             int                `json:"amt_sat"`
	InternalKey        string             `json:"internal_key"`
	TaprootAssetRoot   string             `json:"taproot_asset_root"`
	MerkleRoot         string             `json:"merkle_root"`
	ManagedUtxosAssets []ManagedUtxoAsset `json:"assets"`
}

type ManagedUtxoAsset struct {
	Version          string                      `json:"version"`
	AssetGenesis     ManagedUtxoAssetGenesis     `json:"asset_genesis"`
	Amount           int                         `json:"amount"`
	LockTime         int                         `json:"lock_time"`
	RelativeLockTime int                         `json:"relative_lock_time"`
	ScriptVersion    int                         `json:"script_version"`
	ScriptKey        string                      `json:"script_key"`
	ScriptKeyIsLocal bool                        `json:"script_key_is_local"`
	AssetGroup       ManagedUtxoAssetGroup       `json:"asset_group"`
	ChainAnchor      ManagedUtxoAssetChainAnchor `json:"chain_anchor"`
	IsSpent          bool                        `json:"is_spent"`
	LeaseOwner       string                      `json:"lease_owner"`
	LeaseExpiry      int                         `json:"lease_expiry"`
	IsBurn           bool                        `json:"is_burn"`
}

type ManagedUtxoAssetGenesis struct {
	GenesisPoint string `json:"genesis_point"`
	Name         string `json:"name"`
	MetaHash     string `json:"meta_hash"`
	AssetID      string `json:"asset_id"`
	AssetType    string `json:"asset_type"`
	OutputIndex  int    `json:"output_index"`
	Version      int    `json:"version"`
}

type ManagedUtxoAssetGroup struct {
	RawGroupKey     string `json:"raw_group_key"`
	TweakedGroupKey string `json:"tweaked_group_key"`
	AssetWitness    string `json:"asset_witness"`
}

type ManagedUtxoAssetChainAnchor struct {
	AnchorTx         string `json:"anchor_tx"`
	AnchorBlockHash  string `json:"anchor_block_hash"`
	AnchorOutpoint   string `json:"anchor_outpoint"`
	InternalKey      string `json:"internal_key"`
	MerkleRoot       string `json:"merkle_root"`
	TapscriptSibling string `json:"tapscript_sibling"`
	BlockHeight      int    `json:"block_height"`
}

func ListUtxosResponseToManagedUtxos(listUtxosResponse *taprpc.ListUtxosResponse) *[]ManagedUtxo {
	var managedUtxos []ManagedUtxo
	for op, utxo := range listUtxosResponse.ManagedUtxos {
		var managedUtxo ManagedUtxo
		var managedUtxosAssets []ManagedUtxoAsset
		for _, asset := range utxo.Assets {
			var managedUtxoAssetGroup ManagedUtxoAssetGroup
			if asset.AssetGroup == nil {
				managedUtxoAssetGroup = ManagedUtxoAssetGroup{}
			} else {
				managedUtxoAssetGroup = ManagedUtxoAssetGroup{
					RawGroupKey:     hex.EncodeToString(asset.AssetGroup.RawGroupKey),
					TweakedGroupKey: hex.EncodeToString(asset.AssetGroup.TweakedGroupKey),
					AssetWitness:    hex.EncodeToString(asset.AssetGroup.AssetWitness),
				}
			}
			managedUtxosAssets = append(managedUtxosAssets, ManagedUtxoAsset{
				Version: asset.Version.String(),
				AssetGenesis: ManagedUtxoAssetGenesis{
					GenesisPoint: asset.AssetGenesis.GenesisPoint,
					Name:         asset.AssetGenesis.Name,
					MetaHash:     hex.EncodeToString(asset.AssetGenesis.MetaHash),
					AssetID:      hex.EncodeToString(asset.AssetGenesis.AssetId),
					AssetType:    asset.AssetGenesis.AssetType.String(),
					OutputIndex:  int(asset.AssetGenesis.OutputIndex),
					Version:      int(asset.Version),
				},
				Amount:           int(asset.Amount),
				LockTime:         int(asset.LockTime),
				RelativeLockTime: int(asset.RelativeLockTime),
				ScriptVersion:    int(asset.ScriptVersion),
				ScriptKey:        hex.EncodeToString(asset.ScriptKey),
				ScriptKeyIsLocal: asset.ScriptKeyIsLocal,
				AssetGroup:       managedUtxoAssetGroup,
				ChainAnchor: ManagedUtxoAssetChainAnchor{
					AnchorTx:         hex.EncodeToString(asset.ChainAnchor.AnchorTx),
					AnchorBlockHash:  asset.ChainAnchor.AnchorBlockHash,
					AnchorOutpoint:   asset.ChainAnchor.AnchorOutpoint,
					InternalKey:      hex.EncodeToString(asset.ChainAnchor.InternalKey),
					MerkleRoot:       hex.EncodeToString(asset.ChainAnchor.MerkleRoot),
					TapscriptSibling: hex.EncodeToString(asset.ChainAnchor.TapscriptSibling),
					BlockHeight:      int(asset.ChainAnchor.BlockHeight),
				},
				IsSpent:     asset.IsSpent,
				LeaseOwner:  hex.EncodeToString(asset.LeaseOwner),
				LeaseExpiry: int(asset.LeaseExpiry),
				IsBurn:      asset.IsBurn,
			})
		}
		managedUtxo = ManagedUtxo{
			Op:                 op,
			OutPoint:           utxo.OutPoint,
			Time:               0,
			AmtSat:             int(utxo.AmtSat),
			InternalKey:        hex.EncodeToString(utxo.InternalKey),
			TaprootAssetRoot:   hex.EncodeToString(utxo.TaprootAssetRoot),
			MerkleRoot:         hex.EncodeToString(utxo.MerkleRoot),
			ManagedUtxosAssets: managedUtxosAssets,
		}
		managedUtxos = append(managedUtxos, managedUtxo)
	}
	return &managedUtxos
}

func ManagedUtxosFilterByAssetId(utxos *[]ManagedUtxo, assetId string) *[]ManagedUtxo {
	var managedUtxos []ManagedUtxo
	for _, utxo := range *utxos {
		var assets []ManagedUtxoAsset
		for _, asset := range utxo.ManagedUtxosAssets {
			if assetId == asset.AssetGenesis.AssetID {
				assets = append(assets, asset)
			}
		}
		if len(assets) == 0 {
			continue
		}
		utxo.ManagedUtxosAssets = assets
		managedUtxos = append(managedUtxos, utxo)
	}
	return &managedUtxos
}

func GetAllOutpointsOfManagedUtxos(managedUtxos *[]ManagedUtxo) []string {
	var ops []string
	for _, utxo := range *managedUtxos {
		ops = append(ops, utxo.OutPoint)
	}
	return ops
}

func GetTimeForManagedUtxoByBitcoind(token string, managedUtxos *[]ManagedUtxo) (*[]ManagedUtxo, error) {
	ops := GetAllOutpointsOfManagedUtxos(managedUtxos)
	opMapTime, err := PostCallBitcoindToQueryTimeByOutpoints(token, ops)
	if err != nil {
		return nil, err
	}
	for i, utxo := range *managedUtxos {
		(*managedUtxos)[i].Time = opMapTime.Data[utxo.OutPoint]
	}
	return managedUtxos, nil
}

func GetTimeForManagedUtxoByBitcoind2(host, token string, managedUtxos *[]ManagedUtxo) (*[]ManagedUtxo, error) {
	ops := GetAllOutpointsOfManagedUtxos(managedUtxos)
	opMapTime, err := PostCallBitcoindToQueryTimeByOutpoints2(host, token, ops)
	if err != nil {
		return nil, err
	}
	for i, utxo := range *managedUtxos {
		(*managedUtxos)[i].Time = opMapTime.Data[utxo.OutPoint]
	}
	return managedUtxos, nil
}

func SortAssetUtxos(managedUtxos *[]ManagedUtxo) *[]ManagedUtxo {
	if managedUtxos == nil {
		return nil
	}
	SortTimeDescInAssetUtxos := func(i, j int) bool {
		return (*managedUtxos)[i].Time > (*managedUtxos)[j].Time
	}
	sort.Slice(*managedUtxos, SortTimeDescInAssetUtxos)
	return managedUtxos
}

func AssetUtxos(token string, assetId string) string {
	response, err := ListUtxosAndGetResponse(true)
	if err != nil {
		return MakeJsonErrorResult(ListUtxosAndGetResponseErr, err.Error(), nil)
	}
	managedUtxos := ListUtxosResponseToManagedUtxos(response)
	managedUtxos = ManagedUtxosFilterByAssetId(managedUtxos, assetId)
	managedUtxos, err = GetTimeForManagedUtxoByBitcoind(token, managedUtxos)
	if err != nil {
		LogError("GetTimeForManagedUtxoByBitcoind", err)
	}
	managedUtxos = SortAssetUtxos(managedUtxos)
	return MakeJsonErrorResult(SUCCESS, SuccessError, managedUtxos)
}

func ListUtxosAndGetResponse(includeLeased bool) (*taprpc.ListUtxosResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()
	client := taprpc.NewTaprootAssetsClient(conn)
	request := &taprpc.ListUtxosRequest{
		IncludeLeased: includeLeased,
	}
	return client.ListUtxos(context.Background(), request)
}

func ListUtxos(includeLeased bool) string {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return ""
	}
	defer clearUp()
	client := taprpc.NewTaprootAssetsClient(conn)
	request := &taprpc.ListUtxosRequest{
		IncludeLeased: includeLeased,
	}
	response, err := client.ListUtxos(context.Background(), request)
	if err != nil {
		fmt.Printf("%s taprpc ListUtxos Error: %v\n", GetTimeNow(), err)
		return ""
	}
	return response.String()
}

func NewAddr(assetId string, amt int, token string, deviceId string) string {
	response, err := rpcclient.NewAddr(assetId, amt)
	if err != nil {
		return MakeJsonErrorResult(NewAddrErr, err.Error(), "")
	}
	result := QueriedAddr{}
	result.GetData(response)
	UploadAssetAddr(token, &AssetAddrSetRequest{
		Encoded:          result.Encoded,
		AssetId:          result.AssetId,
		AssetType:        result.AssetType,
		Amount:           result.Amount,
		GroupKey:         result.GroupKey,
		ScriptKey:        result.ScriptKey,
		InternalKey:      result.InternalKey,
		TapscriptSibling: result.TapscriptSibling,
		TaprootOutputKey: result.TaprootOutputKey,
		ProofCourierAddr: result.ProofCourierAddr,
		AssetVersion:     result.AssetVersion,
		DeviceID:         deviceId,
	})
	return MakeJsonErrorResult(SUCCESS, "", result)
}

func NewAddrAndGetResponseEncoded(assetId string, amt int, token string, deviceId string) (string, error) {
	response, err := rpcclient.NewAddr(assetId, amt)
	if err != nil {
		return "", err
	}
	result := QueriedAddr{}
	result.GetData(response)
	UploadAssetAddr(token, &AssetAddrSetRequest{
		Encoded:          result.Encoded,
		AssetId:          result.AssetId,
		AssetType:        result.AssetType,
		Amount:           result.Amount,
		GroupKey:         result.GroupKey,
		ScriptKey:        result.ScriptKey,
		InternalKey:      result.InternalKey,
		TapscriptSibling: result.TapscriptSibling,
		TaprootOutputKey: result.TaprootOutputKey,
		ProofCourierAddr: result.ProofCourierAddr,
		AssetVersion:     result.AssetVersion,
		DeviceID:         deviceId,
	})
	return response.Encoded, nil
}

func QueryAddrs(assetId string) string {
	addrRcv, err := rpcclient.AddrReceives()
	if err != nil {
		return MakeJsonErrorResult(AddrReceivesErr, err.Error(), "")
	}
	addrMap := make(map[string]int)
	for _, events := range addrRcv.Events {
		addrMap[events.Addr.Encoded]++
	}
	_addrs, err := rpcclient.QueryAddr()
	if err != nil {
		return MakeJsonErrorResult(QueryAddrErr, err.Error(), "")
	}
	var addrs []QueriedAddr
	for _, a := range _addrs.Addrs {
		if assetId != "" && assetId != hex.EncodeToString(a.AssetId) {
			continue
		}
		addrTemp := QueriedAddr{}
		addrTemp.GetData(a)
		addrTemp.ReceiveNum = addrMap[addrTemp.Encoded]
		addrs = append(addrs, addrTemp)
	}
	if len(addrs) == 0 {
		return MakeJsonErrorResult(SUCCESS, "NOT_FOUND", addrs)
	}
	return MakeJsonErrorResult(SUCCESS, "", addrs)
}

func SendAssets(jsonAddrs string, feeRate int64, token string, deviceId string) string {
	if int(feeRate) > FeeRateSatPerBToSatPerKw(500) {
		err := errors.New("fee rate exceeds max(500)")
		return MakeJsonErrorResult(FeeRateExceedMaxErr, err.Error(), nil)
	}
	_, err := IsTokenValid(token)
	if err != nil {
		LogError("token is invalid", err)
	}
	var addrs []string
	err = json.Unmarshal([]byte(jsonAddrs), &addrs)
	if err != nil {
		return MakeJsonErrorResult(JsonUnmarshalErr, "Please use the correct json format", nil)
	}
	{
		addrMap := make(map[string]bool)
		for _, addr := range addrs {
			if addrMap[addr] {
				return MakeJsonErrorResult(DuplicateAddrErr, "Duplicate addr("+addr+")", nil)
			}
			addrMap[addr] = true
		}
	}
	var formAddr string
	{
		if len(addrs) == 1 {
			addr := addrs[0]
			var decodedAddr *taprpc.Addr
			decodedAddr, err = rpcclient.DecodeAddr(addr)
			if err != nil {
				LogError("Decode Addr[0] (before send assets)", err)
			} else {
				if decodedAddr.AssetType == taprpc.AssetType_COLLECTIBLE {
					assetId := hex.EncodeToString(decodedAddr.AssetId)
					formAddr, err = GetReceiveAddrByAssetId(assetId)
					if err != nil {
						LogError("Get Receive Addr By AssetId (before send assets)", err)
					}
				}
			}
		}
	}

	response, err := sendAssets(addrs, uint32(feeRate))
	if err != nil {
		return MakeJsonErrorResult(sendAssetsErr, err.Error(), nil)
	}

	{
		var batchTransfersRequest []BatchTransferRequest
		var decodedAddr *taprpc.Addr
		var totalAmount int
		for index, addr := range addrs {
			decodedAddr, err = rpcclient.DecodeAddr(addr)
			if err != nil {
				continue
			}
			totalAmount += int(decodedAddr.Amount)
			txid, _ := getTransactionAndIndexByOutpoint(response.Transfer.Outputs[0].Anchor.Outpoint)
			batchTransfersRequest = append(batchTransfersRequest, BatchTransferRequest{
				Encoded:            decodedAddr.Encoded,
				AssetID:            hex.EncodeToString(decodedAddr.AssetId),
				Amount:             int(decodedAddr.Amount),
				ScriptKey:          hex.EncodeToString(decodedAddr.ScriptKey),
				InternalKey:        hex.EncodeToString(decodedAddr.InternalKey),
				TaprootOutputKey:   hex.EncodeToString(decodedAddr.TaprootOutputKey),
				ProofCourierAddr:   decodedAddr.ProofCourierAddr,
				Txid:               txid,
				Index:              index,
				TransferTimestamp:  GetTimestamp(),
				AnchorTxHash:       hex.EncodeToString(response.Transfer.AnchorTxHash),
				AnchorTxHeightHint: int(response.Transfer.AnchorTxHeightHint),
				AnchorTxChainFees:  int(response.Transfer.AnchorTxChainFees),
				DeviceID:           deviceId,
			})
		}
		for i, _ := range batchTransfersRequest {
			batchTransfersRequest[i].TxTotalAmount = totalAmount
		}
		err = UploadBatchTransfers(token, &batchTransfersRequest)
		if err != nil {
			LogError("; Assets sent, but upload failed.", err)
		}
	}

	{
		var decodedAddr *taprpc.Addr
		for _, addr := range addrs {
			decodedAddr, err = rpcclient.DecodeAddr(addr)
			if decodedAddr != nil {
				break
			}
		}
		if decodedAddr != nil {
			assetId := hex.EncodeToString(decodedAddr.AssetId)
			_, err = GetAssetRecommendUserByJsonAddrs(token, assetId, jsonAddrs, deviceId)
			if err != nil {
				LogError("GetAssetRecommendUserByJsonAddrs failed", err)
			} else {
			}
		}
	}

	txid, _ := getTransactionAndIndexByOutpoint(response.Transfer.Outputs[0].Anchor.Outpoint)

	{
		if len(addrs) == 1 {
			addr := addrs[0]
			var decodedAddr *taprpc.Addr
			decodedAddr, err = rpcclient.DecodeAddr(addr)
			if err != nil {
				LogError("Decode Addr[0]", err)
			} else {
				if decodedAddr.AssetType == taprpc.AssetType_COLLECTIBLE {
					assetId := hex.EncodeToString(decodedAddr.AssetId)
					_time := GetTimestamp()
					err = UploadNftTransfer(token, deviceId, txid, assetId, _time, formAddr, addr)
					if err != nil {
						LogError("Upload NftTransfer", err)
					}
				}
			}
		}
	}

	return MakeJsonErrorResult(SUCCESS, SuccessError, txid)
}

func sendAssets(addrs []string, feeRate uint32) (*taprpc.SendAssetResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()
	client := taprpc.NewTaprootAssetsClient(conn)

	request := &taprpc.SendAssetRequest{
		TapAddrs: addrs,
	}
	if feeRate > 0 {
		request.FeeRate = feeRate
	}

	response, err := client.SendAsset(context.Background(), request)
	if err != nil {
		if strings.Contains(err.Error(), tapfreighter.ErrMatchingAssetsNotFound.Error()) {
			return nil, fmt.Errorf("无可使用的资产（资产余额不足 或 资产锁定中）")
		}
		if strings.Contains(err.Error(), "on total output value") {
			return nil, fmt.Errorf("fee 比例过高，大于20%%，请调整 feeRate 参数 %w", err)
		}
		return nil, err
	}
	return response, nil
}

func SubscribeReceiveAssetEventNtfns() {

}

func SubscribeSendAssetEventNtfns() {

}

func VerifyProof() {

}

func TapStopDaemon() bool {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return false
	}
	defer clearUp()
	client := taprpc.NewTaprootAssetsClient(conn)
	request := &taprpc.StopRequest{}
	_, err = client.StopDaemon(context.Background(), request)
	if err != nil {
		fmt.Printf("%s taprpc TapStopDaemon Error: %v\n", GetTimeNow(), err)
		return false
	}
	return true
}

func fetchAssetMeta(isHash bool, data string) (*taprpc.AssetMeta, error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()

	client := taprpc.NewTaprootAssetsClient(conn)
	request := &taprpc.FetchAssetMetaRequest{}
	if isHash {
		request.Asset = &taprpc.FetchAssetMetaRequest_MetaHashStr{
			MetaHashStr: data,
		}
	} else {
		request.Asset = &taprpc.FetchAssetMetaRequest_AssetIdStr{
			AssetIdStr: data,
		}
	}
	response, err := client.FetchAssetMeta(context.Background(), request)
	return response, err
}

func listBalances(useGroupKey bool, assetFilter, groupKeyFilter []byte) (*taprpc.ListBalancesResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()
	client := taprpc.NewTaprootAssetsClient(conn)
	request := &taprpc.ListBalancesRequest{
		AssetFilter:    assetFilter,
		GroupKeyFilter: groupKeyFilter,
	}
	if useGroupKey {
		request.GroupBy = &taprpc.ListBalancesRequest_GroupKey{GroupKey: true}
	} else {
		request.GroupBy = &taprpc.ListBalancesRequest_AssetId{AssetId: true}
	}
	response, err := client.ListBalances(context.Background(), request)
	return response, err
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

type ListAssetGroupBalanceInfo struct {
	GroupKey string `json:"group_key"`
	Balance  string `json:"balance"`
}

func ProcessListBalancesResponse(response *taprpc.ListBalancesResponse) *[]ListAssetBalanceInfo {
	var listAssetBalanceInfos []ListAssetBalanceInfo
	for _, balance := range response.AssetBalances {
		listAssetBalanceInfos = append(listAssetBalanceInfos, ListAssetBalanceInfo{
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
	return &listAssetBalanceInfos
}

func ExcludeListBalancesResponseCollectible(listAssetBalanceInfos *[]ListAssetBalanceInfo) *[]ListAssetBalanceInfo {
	var listAssetBalances []ListAssetBalanceInfo
	for _, balance := range *listAssetBalanceInfos {
		if balance.AssetType == taprpc.AssetType_NORMAL.String() {
			listAssetBalances = append(listAssetBalances, balance)
		}
	}
	return &listAssetBalances
}

func ProcessListBalancesByGroupKeyResponse(response *taprpc.ListBalancesResponse) *[]ListAssetGroupBalanceInfo {
	var listAssetBalanceInfos []ListAssetGroupBalanceInfo
	for _, balance := range response.AssetGroupBalances {
		listAssetBalanceInfos = append(listAssetBalanceInfos, ListAssetGroupBalanceInfo{
			GroupKey: hex.EncodeToString(balance.GroupKey),
			Balance:  strconv.FormatUint(balance.Balance, 10),
		})
	}
	return &listAssetBalanceInfos
}

func ListBalances() string {
	response, err := listBalances(false, nil, nil)
	if err != nil {
		return MakeJsonErrorResult(listBalancesErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", ProcessListBalancesResponse(response))
}

func ListNormalBalances() string {
	response, err := listBalances(false, nil, nil)
	if err != nil {
		return MakeJsonErrorResult(listBalancesErr, err.Error(), nil)
	}
	processed := ProcessListBalancesResponse(response)
	filtered := ExcludeListBalancesResponseCollectible(processed)
	return MakeJsonErrorResult(SUCCESS, "", filtered)
}

func getAssetsDecimal() (*map[string]int, error) {
	assetsDecimal := make(map[string]int)

	assets, err := listAssets(false, true, false)
	if err != nil {
		return nil, AppendErrorInfo(err, "listAssets")
	} else {
		for _, asset := range assets.Assets {
			if asset.AssetGenesis == nil {
				continue
			}
			decimal := 0
			if asset.DecimalDisplay != nil {
				decimal = int(asset.DecimalDisplay.DecimalDisplay)
			}
			assetsDecimal[hex.EncodeToString(asset.AssetGenesis.AssetId)] = decimal
		}
	}
	assets2, err := listAssets(false, false, true)
	if err != nil {
		return nil, AppendErrorInfo(err, "listAssets2")
	} else {
		for _, asset := range assets2.Assets {
			if asset.AssetGenesis == nil {
				continue
			}
			decimal := 0
			if asset.DecimalDisplay != nil {
				decimal = int(asset.DecimalDisplay.DecimalDisplay)
			}
			assetsDecimal[hex.EncodeToString(asset.AssetGenesis.AssetId)] = decimal
		}
	}
	return &assetsDecimal, nil
}

func ListNormalBalances2() string {
	response, err := listBalances(false, nil, nil)
	if err != nil {
		return MakeJsonErrorResult(listBalancesErr, err.Error(), nil)
	}
	assetsDecimal, err := getAssetsDecimal()
	if err != nil {
		return MakeJsonErrorResult(GetAssetsDecimalErr, err.Error(), nil)
	}
	processed := ProcessListBalancesResponse2(response, assetsDecimal)
	filtered := ExcludeListBalancesResponseCollectible2(processed)
	return MakeJsonErrorResult(SUCCESS, "", filtered)
}

func ListBalancesByGroupKey() string {
	response, err := listBalances(true, nil, nil)
	if err != nil {
		return MakeJsonErrorResult(listBalancesErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", ProcessListBalancesByGroupKeyResponse(response))
}

func listAssets(withWitness, includeSpent, includeLeased bool) (*taprpc.ListAssetResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()
	client := taprpc.NewTaprootAssetsClient(conn)
	request := &taprpc.ListAssetRequest{
		WithWitness:             withWitness,
		IncludeSpent:            includeSpent,
		IncludeLeased:           includeLeased,
		IncludeUnconfirmedMints: true,
	}
	response, err := client.ListAssets(context.Background(), request)
	return response, err
}

func CheckAssetIssuanceIsLocal(assetId string) string {
	keys, err := assetLeafKeys(assetId, universerpc.ProofType_PROOF_TYPE_ISSUANCE)
	if err != nil || len(keys.AssetKeys) == 0 {
		return MakeJsonErrorResult(assetLeafKeysErr, fmt.Errorf("failed to get asset info: %v", err).Error(), "")
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
			return MakeJsonErrorResult(ListBatchesAndGetResponseErr, fmt.Errorf("failed to get mint info: %v", err).Error(), "")
		}
		for _, batch := range listBatch.Batches {
			if batch.Batch.BatchTxid == opStr[0] {
				leaves, err := assetLeaves(false, assetId, universerpc.ProofType_PROOF_TYPE_ISSUANCE)
				if err != nil {
					return MakeJsonErrorResult(assetLeavesErr, fmt.Errorf("failed to get asset info: %v", err).Error(), "")
				}
				result.Amount = int64(leaves.Leaves[0].Asset.Amount)
				transactions, err := GetTransactionsAndGetResponse()
				if err != nil {
					return MakeJsonErrorResult(GetTransactionsAndGetResponseErr, fmt.Errorf("failed to get asset info: %v", err).Error(), "")
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
		return MakeJsonErrorResult(SUCCESS, "", result)
	}
	return MakeJsonErrorResult(GetAssetInfoErr, fmt.Errorf("failed to get asset info: %v", err).Error(), "")
}

type ListAssetsResponse struct {
	Version          string                            `json:"version"`
	AssetGenesis     ListAssetsResponseAssetGenesis    `json:"asset_genesis"`
	Amount           int                               `json:"amount"`
	LockTime         int32                             `json:"lock_time"`
	RelativeLockTime int32                             `json:"relative_lock_time"`
	ScriptVersion    int32                             `json:"script_version"`
	ScriptKey        string                            `json:"script_key"`
	ScriptKeyIsLocal bool                              `json:"script_key_is_local"`
	ChainAnchor      ListAssetsResponseChainAnchor     `json:"chain_anchor"`
	PrevWitnesses    []ListAssetsResponsePrevWitnesses `json:"prev_witnesses"`
	AssetGroup       ListAssetsResponseAssetGroup      `json:"asset_group"`
	IsSpent          bool                              `json:"is_spent"`
	LeaseOwner       string                            `json:"lease_owner"`
	LeaseExpiry      int                               `json:"lease_expiry"`
	IsBurn           bool                              `json:"is_burn"`
}

type ListAssetsResponseAssetGenesis struct {
	GenesisPoint string `json:"genesis_point"`
	Name         string `json:"name"`
	MetaHash     string `json:"meta_hash"`
	AssetID      string `json:"asset_id"`
	AssetType    string `json:"asset_type"`
	OutputIndex  int    `json:"output_index"`
	Version      int    `json:"version"`
}

type ListAssetsResponseChainAnchor struct {
	AnchorTx         string `json:"anchor_tx"`
	AnchorBlockHash  string `json:"anchor_block_hash"`
	AnchorOutpoint   string `json:"anchor_outpoint"`
	InternalKey      string `json:"internal_key"`
	MerkleRoot       string `json:"merkle_root"`
	TapscriptSibling string `json:"tapscript_sibling"`
	BlockHeight      int    `json:"block_height"`
}

type ListAssetsResponsePrevWitnesses struct {
	PrevID    ListAssetsResponsePrevWitnessesPrevID `json:"prev_id"`
	TxWitness []string                              `json:"tx_witness"`
}

type ListAssetsResponsePrevWitnessesPrevID struct {
	AnchorPoint string `json:"anchor_point"`
	AssetID     string `json:"asset_id"`
	ScriptKey   string `json:"script_key"`
	Amount      int    `json:"amount"`
}

type ListAssetsResponseAssetGroup struct {
	RawGroupKey     string `json:"raw_group_key"`
	TweakedGroupKey string `json:"tweaked_group_key"`
	AssetWitness    string `json:"asset_witness"`
}

func ListAssetsProcessed(withWitness, includeSpent, includeLeased bool) (*[]ListAssetsResponse, error) {
	var listAssetsResponse []ListAssetsResponse
	response, err := listAssets(withWitness, includeSpent, includeLeased)
	if err != nil {
		return nil, err
	}
	for _, asset := range response.Assets {
		var listAssetsResponsePrevWitnesses []ListAssetsResponsePrevWitnesses
		for _, witness := range asset.PrevWitnesses {
			var txWitness []string
			for _, txWit := range witness.TxWitness {
				txWitness = append(txWitness, hex.EncodeToString(txWit))
			}
			listAssetsResponsePrevWitnesses = append(listAssetsResponsePrevWitnesses, ListAssetsResponsePrevWitnesses{
				PrevID: ListAssetsResponsePrevWitnessesPrevID{
					AnchorPoint: witness.PrevId.AnchorPoint,
					AssetID:     hex.EncodeToString(witness.PrevId.AssetId),
					ScriptKey:   hex.EncodeToString(witness.PrevId.ScriptKey),
					Amount:      int(witness.PrevId.Amount),
				},
				TxWitness: txWitness,
			})
		}
		var listAssetsResponseAssetGroup ListAssetsResponseAssetGroup
		if asset.AssetGroup != nil {
			listAssetsResponseAssetGroup = ListAssetsResponseAssetGroup{
				RawGroupKey:     hex.EncodeToString(asset.AssetGroup.RawGroupKey),
				TweakedGroupKey: hex.EncodeToString(asset.AssetGroup.TweakedGroupKey),
				AssetWitness:    hex.EncodeToString(asset.AssetGroup.AssetWitness),
			}
		}
		listAssetsResponse = append(listAssetsResponse, ListAssetsResponse{
			Version: asset.Version.String(),
			AssetGenesis: ListAssetsResponseAssetGenesis{
				GenesisPoint: asset.AssetGenesis.GenesisPoint,
				Name:         asset.AssetGenesis.Name,
				MetaHash:     hex.EncodeToString(asset.AssetGenesis.MetaHash),
				AssetID:      hex.EncodeToString(asset.AssetGenesis.AssetId),
				AssetType:    asset.AssetGenesis.AssetType.String(),
				OutputIndex:  int(asset.AssetGenesis.OutputIndex),
				Version:      int(asset.Version),
			},
			Amount:           int(asset.Amount),
			LockTime:         asset.LockTime,
			RelativeLockTime: asset.RelativeLockTime,
			ScriptVersion:    asset.ScriptVersion,
			ScriptKey:        hex.EncodeToString(asset.ScriptKey),
			ScriptKeyIsLocal: asset.ScriptKeyIsLocal,
			ChainAnchor: ListAssetsResponseChainAnchor{
				AnchorTx:         hex.EncodeToString(asset.ChainAnchor.AnchorTx),
				AnchorBlockHash:  asset.ChainAnchor.AnchorBlockHash,
				AnchorOutpoint:   asset.ChainAnchor.AnchorOutpoint,
				InternalKey:      hex.EncodeToString(asset.ChainAnchor.InternalKey),
				MerkleRoot:       hex.EncodeToString(asset.ChainAnchor.MerkleRoot),
				TapscriptSibling: hex.EncodeToString(asset.ChainAnchor.TapscriptSibling),
				BlockHeight:      int(asset.ChainAnchor.BlockHeight),
			},
			PrevWitnesses: listAssetsResponsePrevWitnesses,
			AssetGroup:    listAssetsResponseAssetGroup,
			IsSpent:       asset.IsSpent,
			LeaseOwner:    hex.EncodeToString(asset.LeaseOwner),
			LeaseExpiry:   int(asset.LeaseExpiry),
			IsBurn:        asset.IsBurn,
		})
	}
	return &listAssetsResponse, nil
}

func FilterListAssetsNullGroupKey(listAssetsResponse *[]ListAssetsResponse) *[]ListAssetsResponse {
	if listAssetsResponse == nil {
		return nil
	}
	var results []ListAssetsResponse
	for _, asset := range *listAssetsResponse {
		if asset.AssetGroup.TweakedGroupKey == "" {
			results = append(results, asset)
		}
	}
	return &results
}

func ListNftAssetsAndGetResponse() (*[]ListAssetsResponse, error) {
	processed, err := ListAssetsProcessed(false, false, false)
	if err != nil {
		return nil, err
	}
	var result []ListAssetsResponse
	for index, pr := range *processed {
		if pr.AssetGenesis.AssetType == "COLLECTIBLE" {
			result = append(result, (*processed)[index])
		}
	}
	return &result, nil
}

func ListNftAssetsIncludeSpentAndGetResponse() (*[]ListAssetsResponse, error) {
	processed, err := ListAssetsProcessed(false, true, false)
	if err != nil {
		return nil, err
	}
	var result []ListAssetsResponse
	for index, pr := range *processed {
		if pr.AssetGenesis.AssetType == "COLLECTIBLE" {
			result = append(result, (*processed)[index])
		}
	}
	return &result, nil
}

func ListSpentNftAssetsAndGetResponse() (*[]ListAssetsResponse, error) {
	response, err := ListNftAssetsIncludeSpentAndGetResponse()
	if err != nil {
		return nil, err
	}
	assetNotZero := make(map[string]bool)
	var spentAssets []ListAssetsResponse
	for _, asset := range *response {
		if asset.Amount != 0 {
			assetNotZero[asset.AssetGenesis.AssetType] = true
		}
		if asset.IsSpent {
			spentAssets = append(spentAssets, asset)
		}
	}
	zeroAsset := make(map[string]ListAssetsResponse)
	for _, asset := range spentAssets {
		if !(assetNotZero[asset.AssetGenesis.AssetID]) {
			zeroAsset[asset.AssetGenesis.AssetID] = asset
		}
	}
	var result []ListAssetsResponse
	for _, asset := range zeroAsset {
		result = append(result, asset)
	}
	return &result, nil
}

func GetSpentNftAssets() string {
	response, err := GetSpentNftAssetsAndGetResponse()
	if err != nil {
		return MakeJsonErrorResult(GetSpentNftAssetsAndGetResponseErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func GetSpentNftAssetsAndGetResponse() (*[]ListAssetsSimplifiedResponse, error) {
	response, err := ListSpentNftAssetsAndGetResponse()
	if err != nil {
		return nil, err
	}
	result := ListAssetsResponseSliceToListAssetsSimplifiedResponseSlice(response)
	return result, nil
}

type ListAssetsSimplifiedResponse struct {
	AssetID         string `json:"asset_id"`
	Name            string `json:"name"`
	AssetType       string `json:"asset_type"`
	Amount          int    `json:"amount"`
	IsSpent         bool   `json:"is_spent"`
	TweakedGroupKey string `json:"tweaked_group_key"`
}

func ListAssetsResponseToListAssetsSimplifiedResponse(listAssetsResponse ListAssetsResponse) ListAssetsSimplifiedResponse {
	return ListAssetsSimplifiedResponse{
		AssetID:         listAssetsResponse.AssetGenesis.AssetID,
		Name:            listAssetsResponse.AssetGenesis.Name,
		AssetType:       listAssetsResponse.AssetGenesis.AssetType,
		Amount:          listAssetsResponse.Amount,
		IsSpent:         listAssetsResponse.IsSpent,
		TweakedGroupKey: listAssetsResponse.AssetGroup.TweakedGroupKey,
	}
}

func ListAssetsResponseSliceToListAssetsSimplifiedResponseSlice(listAssetsResponseSlice *[]ListAssetsResponse) *[]ListAssetsSimplifiedResponse {
	if listAssetsResponseSlice == nil {
		return nil
	}
	var listAssetsSimplifiedResponseSlice []ListAssetsSimplifiedResponse
	for _, asset := range *listAssetsResponseSlice {
		listAssetsSimplifiedResponseSlice = append(listAssetsSimplifiedResponseSlice, ListAssetsResponseToListAssetsSimplifiedResponse(asset))
	}
	return &listAssetsSimplifiedResponseSlice
}

func GetGroupAssets(groupKey string) (*[]ListAssetsResponse, error) {
	listNftAssets, err := ListNftAssetsAndGetResponse()
	if err != nil {
		return nil, err
	}
	var result []ListAssetsResponse
	for _, asset := range *listNftAssets {
		tweakedGroupKey := asset.AssetGroup.TweakedGroupKey

		if len(tweakedGroupKey) == 66 && len(groupKey) == 64 {
			tweakedGroupKey = tweakedGroupKey[2:]
		} else if len(tweakedGroupKey) == 64 && len(groupKey) == 66 {
			groupKey = groupKey[2:]
		}

		if tweakedGroupKey == groupKey {
			result = append(result, asset)
		}
	}
	return &result, nil
}

func getGroupAssetsIncludingSpentAndLeased(groupKey string) (listAssets *[]ListAssetsResponse, err error) {

	spentListAssetsResponse, err := ListAssetsProcessed(false, true, false)
	if err != nil {
		return &[]ListAssetsResponse{}, AppendErrorInfo(err, "ListAssetsProcessed spent")
	}

	leasedListAssetsResponse, err := ListAssetsProcessed(false, false, true)
	if err != nil {
		LogError("ListAssetsProcessed leased", err)
		leasedListAssetsResponse = &[]ListAssetsResponse{}
	}

	var gotAssets []ListAssetsResponse
	var zeroAssets []ListAssetsResponse
	zeroAssetExist := make(map[string]bool)

	for _, asset := range *spentListAssetsResponse {
		if asset.AssetGenesis.AssetType != "COLLECTIBLE" {
			continue
		}
		if asset.IsSpent || asset.LeaseOwner != "" || asset.LeaseExpiry != 0 {

			tweakedGroupKey := asset.AssetGroup.TweakedGroupKey

			if len(tweakedGroupKey) == 66 && len(groupKey) == 64 {
				tweakedGroupKey = tweakedGroupKey[2:]
			} else if len(tweakedGroupKey) == 64 && len(groupKey) == 66 {
				groupKey = groupKey[2:]
			}

			if tweakedGroupKey == groupKey {

				if zeroAssetExist[asset.AssetGenesis.AssetID] {
					continue
				}

				asset.Amount = 0
				zeroAssets = append(zeroAssets, asset)

				zeroAssetExist[asset.AssetGenesis.AssetID] = true
			}

		} else {
			tweakedGroupKey := asset.AssetGroup.TweakedGroupKey

			if len(tweakedGroupKey) == 66 && len(groupKey) == 64 {
				tweakedGroupKey = tweakedGroupKey[2:]
			} else if len(tweakedGroupKey) == 64 && len(groupKey) == 66 {
				groupKey = groupKey[2:]
			}

			if tweakedGroupKey == groupKey {
				gotAssets = append(gotAssets, asset)
			}
		}
	}

	for _, asset := range *leasedListAssetsResponse {
		if asset.AssetGenesis.AssetType != "COLLECTIBLE" {
			continue
		}
		if asset.IsSpent || asset.LeaseOwner != "" || asset.LeaseExpiry != 0 {

			tweakedGroupKey := asset.AssetGroup.TweakedGroupKey

			if len(tweakedGroupKey) == 66 && len(groupKey) == 64 {
				tweakedGroupKey = tweakedGroupKey[2:]
			} else if len(tweakedGroupKey) == 64 && len(groupKey) == 66 {
				groupKey = groupKey[2:]
			}

			if tweakedGroupKey == groupKey {

				if zeroAssetExist[asset.AssetGenesis.AssetID] {
					continue
				}

				asset.Amount = 0
				zeroAssets = append(zeroAssets, asset)

				zeroAssetExist[asset.AssetGenesis.AssetID] = true
			}

		} else {
		}
	}

	var listAssetsResponse []ListAssetsResponse
	listAssetsResponse = append(listAssetsResponse, gotAssets...)
	listAssetsResponse = append(listAssetsResponse, zeroAssets...)
	listAssets = &listAssetsResponse

	return listAssets, nil
}

func ListAssetsAll() string {
	response, err := ListAssetsProcessed(true, true, false)
	if err != nil {
		return MakeJsonErrorResult(ListAssetsProcessedErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

type NftGroupIdTag struct {
	Id  string `json:"id"`
	Tag string `json:"tag"`
}

type NftGroup struct {
	GroupKey  string           `json:"group_key"`
	GroupName string           `json:"group_name"`
	Supply    int              `json:"supply"`
	NftIds    *[]NftGroupIdTag `json:"nft_ids"`
}

func ListNftGroups() string {
	resResponse, err := rpcclient.ListGroups()
	if err != nil {
		return MakeJsonErrorResult(ListGroupsErr, err.Error(), nil)
	}

	var groups []NftGroup
	var gotKeys []string
	if resResponse.Groups != nil {
		for key, group := range resResponse.Groups {
			if group.Assets[0].Type != taprpc.AssetType_COLLECTIBLE {
				break
			}
			var nftIds []NftGroupIdTag
			for _, asset := range group.Assets {
				nftIds = append(nftIds, NftGroupIdTag{
					Id:  hex.EncodeToString(asset.Id),
					Tag: asset.Tag,
				})
			}
			meta := Meta{}
			meta.FetchAssetMeta(false, hex.EncodeToString(group.Assets[0].Id))
			{
				assets, err := GetGroupAssets(key)
				if err != nil || assets == nil || len(*assets) == 0 {
					continue
				}
			}
			gotKeys = append(gotKeys, key)
			groups = append(groups, NftGroup{
				GroupKey:  key,
				GroupName: meta.GroupName,
				Supply:    len(group.Assets),
				NftIds:    &nftIds,
			})
		}
	}

	spentAndLeasedGroups, err := GetSpentAndLeasedGroups(gotKeys)
	if err != nil {
		LogError("GetSpentAndLeasedGroups", err)
		spentAndLeasedGroups = []NftGroup{}
	}

	groups = append(groups, spentAndLeasedGroups...)

	return MakeJsonErrorResult(SUCCESS, "", groups)
}

func ListNftAssets() string {
	processed, err := ListAssetsProcessed(false, false, false)
	if err != nil {
		return MakeJsonErrorResult(ListAssetsProcessedErr, err.Error(), nil)
	}
	var result []ListAssetsResponse
	for index, pr := range *processed {
		if pr.AssetGenesis.AssetType == "COLLECTIBLE" {
			result = append(result, (*processed)[index])
		}
	}
	return MakeJsonErrorResult(SUCCESS, "", &result)
}

func ListNonGroupNftAssets() string {
	processed, err := ListAssetsProcessed(false, false, false)
	if err != nil {
		return MakeJsonErrorResult(ListAssetsProcessedErr, err.Error(), nil)
	}
	var result []ListAssetsResponse
	for index, pr := range *processed {
		if pr.AssetGenesis.AssetType == "COLLECTIBLE" {
			result = append(result, (*processed)[index])
		}
	}
	resultFiltered := FilterListAssetsNullGroupKey(&result)
	return MakeJsonErrorResult(SUCCESS, "", resultFiltered)
}

func QueryAllNftByGroup(groupKey string) string {
	response, err := getGroupAssetsIncludingSpentAndLeased(groupKey)
	if err != nil {
		return MakeJsonErrorResult(ListNftAssetsAndGetResponseErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func GetSpentAndLeasedGroups(gotKeys []string) (nftGroups []NftGroup, err error) {

	spentListAssetsResponse, err := ListAssetsProcessed(false, true, false)
	if err != nil {
		LogError("ListAssetsProcessed spent", err)
		spentListAssetsResponse = &[]ListAssetsResponse{}
	}
	leasedListAssetsResponse, err := ListAssetsProcessed(false, false, true)
	if err != nil {
		LogError("ListAssetsProcessed leased", err)
		leasedListAssetsResponse = &[]ListAssetsResponse{}
	}

	var groupKeys []string
	groupKeyExist := make(map[string]bool)

	for _, gotKey := range gotKeys {
		groupKeyExist[gotKey] = true
	}

	for _, asset := range *spentListAssetsResponse {
		if asset.AssetGenesis.AssetType != "COLLECTIBLE" {
			continue
		}
		if asset.IsSpent || asset.LeaseOwner != "" || asset.LeaseExpiry != 0 {

			if asset.AssetGroup.TweakedGroupKey == "" {
				continue
			}

			groupKey := asset.AssetGroup.TweakedGroupKey
			if groupKeyExist[groupKey] {
				continue
			}

			groupKeys = append(groupKeys, groupKey)
			groupKeyExist[groupKey] = true

			nftGroups = append(nftGroups, NftGroup{
				GroupKey:  groupKey,
				GroupName: "",
				Supply:    0,
				NftIds:    &[]NftGroupIdTag{},
			})

		}
	}

	for _, asset := range *leasedListAssetsResponse {
		if asset.AssetGenesis.AssetType != "COLLECTIBLE" {
			continue
		}
		if asset.IsSpent || asset.LeaseOwner != "" || asset.LeaseExpiry != 0 {

			if asset.AssetGroup.TweakedGroupKey == "" {
				continue
			}

			groupKey := asset.AssetGroup.TweakedGroupKey
			if groupKeyExist[groupKey] {
				continue
			}

			groupKeys = append(groupKeys, groupKey)
			groupKeyExist[groupKey] = true

			nftGroups = append(nftGroups, NftGroup{
				GroupKey:  groupKey,
				GroupName: "",
				Supply:    0,
				NftIds:    &[]NftGroupIdTag{},
			})

		}
	}

	keysMapNames, err := GetGroupNamesByGroupKeys(groupKeys)
	if err != nil {
		LogError("GetGroupNamesByGroupKeys", err)
		mapNames := make(map[string]string)
		keysMapNames = &mapNames
	}

	for i, group := range nftGroups {
		if name, ok := (*keysMapNames)[group.GroupKey]; ok {
			nftGroups[i].GroupName = name
		}
	}

	return nftGroups, nil
}

type IsLocalResult struct {
	IsLocal   bool   `json:"is_local"`
	AssetId   string `json:"asset_id"`
	BatchTxid string `json:"batch_txid"`
	Amount    int64  `json:"amount"`
	Timestamp int64  `json:"timestamp"`
	ScriptKey string `json:"script_key"`
}
