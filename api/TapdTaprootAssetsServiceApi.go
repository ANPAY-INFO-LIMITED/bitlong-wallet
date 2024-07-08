package api

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	"github.com/wallet/service/apiConnect"
	"github.com/wallet/service/rpcclient"
	"strconv"
	"strings"
)

func AddrReceives(assetId string) string {
	response, err := rpcclient.AddrReceives()
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	type addrEvent struct {
		CreationTimeUnixSeconds int64           `json:"creation_time_unix_seconds"`
		Addr                    *jsonResultAddr `json:"addr"`
		Status                  string          `json:"status"`
		Outpoint                string          `json:"outpoint"`
		UtxoAmtSat              int64           `json:"utxo_amt_sat"`
		TaprootSibling          string          `json:"taproot_sibling"`
		ConfirmationHeight      int64           `json:"confirmation_height"`
		HasProof                bool            `json:"has_proof"`
	}
	var addrEvents []addrEvent
	for _, event := range response.Events {
		if assetId != "" && assetId != hex.EncodeToString(event.Addr.AssetId) {
			continue
		}
		e := addrEvent{}
		e.CreationTimeUnixSeconds = int64(event.CreationTimeUnixSeconds)
		a := jsonResultAddr{}
		a.getData(event.Addr)
		e.Addr = &a
		e.Status = event.Status.String()
		e.Outpoint = event.Outpoint
		e.UtxoAmtSat = int64(event.UtxoAmtSat)
		e.TaprootSibling = hex.EncodeToString(event.TaprootSibling)
		e.ConfirmationHeight = int64(event.ConfirmationHeight)
		e.HasProof = event.HasProof
		addrEvents = append(addrEvents, e)
	}
	if len(addrEvents) == 0 {
		return MakeJsonErrorResult(SUCCESS, "NOT_FOUND", nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", addrEvents)
}

func BurnAsset(AssetIdStr string, amountToBurn int64) string {
	response, err := rpcclient.BurnAsset(AssetIdStr, uint64(amountToBurn))
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	txHash := hex.EncodeToString(response.BurnTransfer.AnchorTxHash)
	return MakeJsonErrorResult(SUCCESS, "", txHash)
}

func DebugLevel() {

}

func DecodeAddr(addr string) string {
	response, err := rpcclient.DecodeAddr(addr)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	// make result struct
	result := jsonResultAddr{}
	result.getData(response)

	return MakeJsonErrorResult(SUCCESS, "", result)
}

func DecodeProof(rawProof string) {

}

func ExportProof() {

}

func FetchAssetMeta(isHash bool, data string) string {
	response, err := fetchAssetMeta(isHash, data)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", string(response.Data))
}

// GetInfoOfTap
//
//	@Description: GetInfo returns the information for the node.
//	@return string
func GetInfoOfTap() string {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
	}
	defer clearUp()
	client := taprpc.NewTaprootAssetsClient(conn)
	request := &taprpc.GetInfoRequest{}
	response, err := client.GetInfo(context.Background(), request)
	if err != nil {
		fmt.Printf("%s taprpc GetInfo Error: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

// ListAssets
//
//	@Description: ListAssets lists the set of assets owned by the target daemon.
//	@return string
func ListAssets(withWitness, includeSpent, includeLeased bool) string {
	response, err := listAssets(withWitness, includeSpent, includeLeased)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func ListSimpleAssets(withWitness, includeSpent, includeLeased bool) string {
	response, err := listAssets(withWitness, includeSpent, includeLeased)
	if err != nil {
		fmt.Printf("%s taprpc ListAssets Error: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
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
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	if response.Success == false {
		return MakeJsonErrorResult(DefaultErr, response.Error, nil)
	}
	var assets []*taprpc.Asset
	for _, asset := range response.Data.Assets {
		//if hex.EncodeToString(asset.AssetGenesis.GetAssetId()) == assetName {
		if asset.AssetGenesis.Name == assetName {
			assets = append(assets, asset)
		}
	}
	if len(assets) == 0 {
		return MakeJsonErrorResult(DefaultErr, "NOT_FOUND", nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", assets)
}

// ListGroups
//
//	@Description: ListGroups lists the asset groups known to the target daemon, and the assets held in each group.
//	@return string
func ListGroups() string {
	response, err := rpcclient.ListGroups()
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

// ListTransfers
//
//	@Description: ListTransfers lists outbound asset transfer tracked by the target daemon.
//	@return string
func QueryAssetTransfers(assetId string) string {
	response, err := rpcclient.ListTransfers()
	if err != nil {
		fmt.Printf("%s taprpc ListTransfers Error: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	var transfers []transfer
	for _, t := range response.Transfers {
		if assetId != "" && assetId != hex.EncodeToString(t.Inputs[0].AssetId) {
			continue
		}
		newTransfer := transfer{}
		newTransfer.geData(t)
		transfers = append(transfers, newTransfer)
	}
	if len(transfers) == 0 {
		return MakeJsonErrorResult(SUCCESS, "NOT_FOUND", transfers)
	}
	return MakeJsonErrorResult(SUCCESS, "", transfers)
}

// ListUtxos
//
//	@Description: ListUtxos lists the UTXOs managed by the target daemon, and the assets they hold.
//	@return string
func ListUtxos(includeLeased bool) string {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
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

// NewAddr
//
//	@Description:NewAddr makes a new address from the set of request params.
//	@return string
func NewAddr(assetId string, amt int) string {
	response, err := rpcclient.NewAddr(assetId, amt)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), "")
	}
	result := jsonResultAddr{}
	result.getData(response)

	return MakeJsonErrorResult(SUCCESS, "", result)
}

func QueryAddrs(assetId string) string {
	addrRcv, err := rpcclient.QueryAddr()
	if err != nil {
		fmt.Printf("%s taprpc QueryAddrs Error: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(DefaultErr, err.Error(), "")
	}

	var addrs []jsonResultAddr
	for _, a := range addrRcv.Addrs {
		if assetId != "" && assetId != hex.EncodeToString(a.AssetId) {
			continue
		}
		addrTemp := jsonResultAddr{}
		addrTemp.getData(a)
		addrs = append(addrs, addrTemp)
	}
	if len(addrs) == 0 {
		return MakeJsonErrorResult(SUCCESS, "NOT_FOUND", addrs)
	}
	return MakeJsonErrorResult(SUCCESS, "", addrs)
}

// jsonAddrs : ["addrs1","addrs2",...]
func SendAssets(jsonAddrs string, feeRate int64, token string, deviceId string) string {
	var addrs []string
	err := json.Unmarshal([]byte(jsonAddrs), &addrs)
	if err != nil {
		fmt.Printf("%s json.Unmarshal Error: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(DefaultErr, "Please use the correct json format", "")
	}
	response, err := sendAssets(addrs, uint32(feeRate))
	if err != nil {
		return MakeJsonErrorResult(sendAssetsErr, err.Error(), nil)
	}
	// @dev: decode addrs
	var batchTransfersRequest []BatchTransferRequest
	var decodedAddr *taprpc.Addr
	for index, addr := range addrs {
		decodedAddr, err = rpcclient.DecodeAddr(addr)
		if err != nil {
			return MakeJsonErrorResult(DecodeAddrErr, err.Error(), "")
		}
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
	// @dev: Upload
	err = UploadBatchTransfers(token, &batchTransfersRequest)
	if err != nil {
		return MakeJsonErrorResult(UploadBatchTransfersErr, err.Error(), nil)
	}
	txid, _ := getTransactionAndIndexByOutpoint(response.Transfer.Outputs[0].Anchor.Outpoint)
	return MakeJsonErrorResult(SUCCESS, SuccessError, txid)
}

// SendAsset
// @Description:SendAsset uses one or multiple passed Taproot Asset address(es) to attempt to complete an asset send.
// The method returns information w.r.t the on chain send, as well as the proof file information the receiver needs to fully receive the asset.
// @return string
// skipped function SendAsset with unsupported parameter or return types
func sendAssets(addrs []string, feeRate uint32) (*taprpc.SendAssetResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
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
		fmt.Printf("%s taprpc SendAsset Error: %v\n", GetTimeNow(), err)
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

// TapStopDaemon
//
//	@Description: StopDaemon will send a shutdown request to the interrupt handler, triggering a graceful shutdown of the daemon.
//	@return bool
func TapStopDaemon() bool {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
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
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
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
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
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
			Version:      int(balance.AssetGenesis.Version),
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
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", ProcessListBalancesResponse(response))
}

func ListNormalBalances() string {
	response, err := listBalances(false, nil, nil)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	processed := ProcessListBalancesResponse(response)
	filtered := ExcludeListBalancesResponseCollectible(processed)
	return MakeJsonErrorResult(SUCCESS, "", filtered)
}

func ListBalancesByGroupKey() string {
	response, err := listBalances(true, nil, nil)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", ProcessListBalancesByGroupKeyResponse(response))
}

func listAssets(withWitness, includeSpent, includeLeased bool) (*taprpc.ListAssetResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
	}
	defer clearUp()
	client := taprpc.NewTaprootAssetsClient(conn)
	request := &taprpc.ListAssetRequest{
		WithWitness:   withWitness,
		IncludeSpent:  includeSpent,
		IncludeLeased: includeLeased,
	}
	response, err := client.ListAssets(context.Background(), request)
	return response, err
}

func CheckAssetIssuanceIsLocal(assetId string) string {
	keys, err := assetLeafKeys(assetId, universerpc.ProofType_PROOF_TYPE_ISSUANCE)
	if err != nil || len(keys.AssetKeys) == 0 {
		return MakeJsonErrorResult(DefaultErr, fmt.Errorf("failed to get asset info: %v", err).Error(), "")
	}

	result := struct {
		IsLocal   bool   `json:"is_local"`
		AssetId   string `json:"asset_id"`
		BatchTxid string `json:"batch_txid"`
		Amount    int64  `json:"amount"`
		Timestamp int64  `json:"timestamp"`
	}{
		IsLocal: false,
		AssetId: assetId,
	}

	Outpoint := keys.AssetKeys[0].Outpoint
	if o, ok := Outpoint.(*universerpc.AssetKey_OpStr); ok {
		opStr := strings.Split(o.OpStr, ":")
		listBatch, err := ListBatchesAndGetResponse()
		if err != nil {
			return MakeJsonErrorResult(DefaultErr, fmt.Errorf("failed to get mint info: %v", err).Error(), "")
		}
		for _, batch := range listBatch.Batches {
			if batch.BatchTxid == opStr[0] {
				leaves, err := assetLeaves(false, assetId, universerpc.ProofType_PROOF_TYPE_ISSUANCE)
				if err != nil {
					return MakeJsonErrorResult(DefaultErr, fmt.Errorf("failed to get asset info: %v", err).Error(), "")
				}
				result.Amount = int64(leaves.Leaves[0].Asset.Amount)
				transactions, err := GetTransactionsAndGetResponse()
				if err != nil {
					return MakeJsonErrorResult(DefaultErr, fmt.Errorf("failed to get asset info: %v", err).Error(), "")
				}
				for _, tx := range transactions.Transactions {
					if tx.TxHash == opStr[0] {
						result.Timestamp = tx.TimeStamp
						break
					}
				}
				result.IsLocal = true
				result.BatchTxid = batch.BatchTxid
				break
			}
		}
		return MakeJsonErrorResult(SUCCESS, "", result)
	}
	return MakeJsonErrorResult(DefaultErr, fmt.Errorf("failed to get asset info: %v", err).Error(), "")
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
	//SplitCommitment *SplitCommitment `protobuf:"bytes,3,opt,name=split_commitment,json=splitCommitment,proto3" json:"split_commitment,omitempty"`
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
				Version:      int(asset.AssetGenesis.Version),
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

func ListAssetsAll() string {
	response, err := ListAssetsProcessed(true, true, false)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func ListNFTGroups() string {
	resResponse, err := rpcclient.ListGroups()
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	type NFTId struct {
		Id  string `json:"id"`
		Tag string `json:"tag"`
	}
	type Group struct {
		GroupKey  string   `json:"group_key"`
		GroupName string   `json:"group_name"`
		Supply    int      `json:"supply"`
		NFTIds    *[]NFTId `json:"nft_ids"`
	}
	var Groups []Group
	if resResponse.Groups != nil {
		for key, group := range resResponse.Groups {
			if group.Assets[0].Type != taprpc.AssetType_COLLECTIBLE {
				break
			}
			var nftIds []NFTId
			for _, asset := range group.Assets {
				nftIds = append(nftIds, NFTId{
					Id:  hex.EncodeToString(asset.Id),
					Tag: asset.Tag,
				})
			}
			meta := Meta{}
			meta.FetchAssetMeta(false, hex.EncodeToString(group.Assets[0].Id))
			Groups = append(Groups, Group{
				GroupKey:  key,
				GroupName: meta.Name,
				Supply:    len(group.Assets),
				NFTIds:    &nftIds,
			})
		}
	}
	return MakeJsonErrorResult(SUCCESS, "", Groups)
}

func ListNFTAssets() string {
	processed, err := ListAssetsProcessed(false, false, false)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	var result []ListAssetsResponse
	for index, pr := range *processed {
		if pr.AssetGenesis.AssetType == "COLLECTIBLE" {
			result = append(result, (*processed)[index])
		}
	}
	return MakeJsonErrorResult(SUCCESS, "", &result)
}

func QueryAllNFTByGroup() string {

	return ""
}
