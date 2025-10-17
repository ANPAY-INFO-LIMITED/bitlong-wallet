package api

import (
	"encoding/hex"
	"encoding/json"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/service/rpcclient"
	"strconv"
	"strings"
)

type ListAssetBalanceInfo2 struct {
	GenesisPoint string `json:"genesis_point"`
	Name         string `json:"name"`
	MetaHash     string `json:"meta_hash"`
	AssetID      string `json:"asset_id"`
	AssetType    string `json:"asset_type"`
	OutputIndex  int    `json:"output_index"`
	Version      int    `json:"version"`
	Balance      string `json:"balance"`
	Decimal      int    `json:"decimal"`
}

func ProcessListBalancesResponse2(response *taprpc.ListBalancesResponse, assetsDecimal *map[string]int) *[]ListAssetBalanceInfo2 {
	var listAssetBalanceInfos []ListAssetBalanceInfo2
	for _, balance := range response.AssetBalances {
		listAssetBalanceInfos = append(listAssetBalanceInfos, ListAssetBalanceInfo2{
			GenesisPoint: balance.AssetGenesis.GenesisPoint,
			Name:         balance.AssetGenesis.Name,
			MetaHash:     hex.EncodeToString(balance.AssetGenesis.MetaHash),
			AssetID:      hex.EncodeToString(balance.AssetGenesis.AssetId),
			AssetType:    balance.AssetGenesis.AssetType.String(),
			OutputIndex:  int(balance.AssetGenesis.OutputIndex),
			Version:      -1,
			Balance:      strconv.FormatUint(balance.Balance, 10),
			Decimal:      (*assetsDecimal)[hex.EncodeToString(balance.AssetGenesis.AssetId)],
		})
	}
	return &listAssetBalanceInfos
}

func ExcludeListBalancesResponseCollectible2(listAssetBalanceInfos *[]ListAssetBalanceInfo2) *[]ListAssetBalanceInfo2 {
	var listAssetBalances []ListAssetBalanceInfo2
	for _, balance := range *listAssetBalanceInfos {
		if balance.AssetType == taprpc.AssetType_NORMAL.String() {
			listAssetBalances = append(listAssetBalances, balance)
		}
	}
	return &listAssetBalances
}

func FilterListAssetsNullGroupKey2(listAssetsResponse []*ListAssetsResponse) []*ListAssetsResponse {
	if listAssetsResponse == nil {
		return nil
	}
	var results []*ListAssetsResponse
	for _, asset := range listAssetsResponse {
		if asset.AssetGroup.TweakedGroupKey == "" {
			results = append(results, asset)
		}
	}
	return results
}

func PcListNormalBalances2() (*[]ListAssetBalanceInfo2, error) {
	response, err := listBalances(false, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "listBalances")
	}
	assetsDecimal, err := getAssetsDecimal()
	if err != nil {
		return nil, errors.Wrap(err, "getAssetsDecimal")
	}
	processed := ProcessListBalancesResponse2(response, assetsDecimal)
	filtered := ExcludeListBalancesResponseCollectible2(processed)
	return filtered, nil
}

func PcCheckAssetIssuanceIsLocal(assetId string) (*IsLocalResult, error) {
	keys, err := assetLeafKeys(assetId, universerpc.ProofType_PROOF_TYPE_ISSUANCE)
	if err != nil {
		return nil, errors.Wrap(err, "assetLeafKeys")
	}

	if len(keys.AssetKeys) == 0 {
		return nil, errors.New("assetLeafKeys is empty")
	}

	result := IsLocalResult{
		IsLocal: false,
		AssetId: assetId,
	}

	Outpoint := keys.AssetKeys[0].Outpoint
	if o, ok := Outpoint.(*universerpc.AssetKey_OpStr); ok {
		opStr := strings.Split(o.OpStr, ":")
		listBatch, err := ListBatchesAndGetResponse()
		if err != nil {
			return nil, errors.Wrap(err, "ListBatchesAndGetResponse")
		}
		for _, batch := range listBatch.Batches {
			if batch.Batch.BatchTxid == opStr[0] {
				leaves, err := assetLeaves(false, assetId, universerpc.ProofType_PROOF_TYPE_ISSUANCE)
				if err != nil {
					return nil, errors.Wrap(err, "assetLeaves")
				}
				result.Amount = int64(leaves.Leaves[0].Asset.Amount)
				transactions, err := GetTransactionsAndGetResponse()
				if err != nil {
					return nil, errors.Wrap(err, "GetTransactionsAndGetResponse")
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
		return &result, nil
	}
	return nil, errors.New("failed to get asset info")
}

func PcAddrReceives(assetId string) (*[]AddrEvent, error) {
	response, err := rpcclient.AddrReceives()
	if err != nil {
		return nil, errors.Wrap(err, "rpcclient.AddrReceives")
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
		return &addrEvents, nil
	}
	result := SortAddrEvents(&addrEvents)
	return result, nil
}

func PcQueryAssetTransfers(assetId string) (*[]AssetTransferSimplified, error) {
	token := ""
	assetTransfers, err := QueryAssetTransferSimplified(token, assetId)
	if err != nil {
		return nil, errors.Wrap(err, "QueryAssetTransferSimplified")
	}
	assetTransfers = SortAssetTransferSimplified(assetTransfers)
	return assetTransfers, nil
}

func PcAssetUtxos(token string, assetId string) (*[]ManagedUtxo, error) {
	response, err := ListUtxosAndGetResponse(true)
	if err != nil {
		return nil, errors.Wrap(err, "ListUtxosAndGetResponse")
	}
	managedUtxos := ListUtxosResponseToManagedUtxos(response)
	managedUtxos = ManagedUtxosFilterByAssetId(managedUtxos, assetId)
	managedUtxos, err = GetTimeForManagedUtxoByBitcoind(token, managedUtxos)
	if err != nil {
		logrus.Infoln("GetTimeForManagedUtxoByBitcoind", err)
	}
	managedUtxos = SortAssetUtxos(managedUtxos)
	return managedUtxos, nil
}

func PcNewAddr(assetId string, amt int, token string, deviceId string) (*QueriedAddr, error) {
	response, err := rpcclient.NewAddr(assetId, amt)
	if err != nil {
		return nil, errors.Wrap(err, "rpcclient.NewAddr")
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
	return &result, nil
}

func PcSendAssets(jsonAddrs string, feeRate int64, token string, deviceId string) (string, error) {
	if int(feeRate) > FeeRateSatPerBToSatPerKw(500) {
		return "", errors.New("fee rate exceeds max(500)")
	}
	_, err := IsTokenValid(token)
	if err != nil {
		logrus.Infoln("token is invalid", err)
	}
	var addrs []string
	err = json.Unmarshal([]byte(jsonAddrs), &addrs)
	if err != nil {
		return "", errors.Wrap(err, "json.Unmarshal")
	}
	{
		addrMap := make(map[string]bool)
		for _, addr := range addrs {
			if addrMap[addr] {
				return "", errors.New("Duplicate addr(" + addr + ")")
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
				logrus.Infoln("Decode Addr[0] (before send assets)", err)
			} else {
				if decodedAddr.AssetType == taprpc.AssetType_COLLECTIBLE {
					assetId := hex.EncodeToString(decodedAddr.AssetId)
					formAddr, err = GetReceiveAddrByAssetId(assetId)
					if err != nil {
						logrus.Infoln("Get Receive Addr By AssetId (before send assets)", err)
					}
				}
			}
		}
	}

	response, err := sendAssets(addrs, uint32(feeRate))
	if err != nil {
		return "", errors.Wrap(err, "sendAssets")
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
			logrus.Infoln("; Assets sent, but upload failed.", err)
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
				logrus.Infoln("GetAssetRecommendUserByJsonAddrs failed", err)
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
				logrus.Infoln("Decode Addr[0]", err)
			} else {
				if decodedAddr.AssetType == taprpc.AssetType_COLLECTIBLE {
					assetId := hex.EncodeToString(decodedAddr.AssetId)
					_time := GetTimestamp()
					err = UploadNftTransfer(token, deviceId, txid, assetId, _time, formAddr, addr)
					if err != nil {
						logrus.Infoln("Upload NftTransfer", err)
					}
				}
			}
		}
	}

	return txid, nil
}

func PcListNftGroups() ([]*NftGroup, error) {
	resResponse, err := rpcclient.ListGroups()
	if err != nil {
		return nil, errors.Wrap(err, "rpcclient.ListGroups")
	}

	var groups []*NftGroup
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
			groups = append(groups, &NftGroup{
				GroupKey:  key,
				GroupName: meta.GroupName,
				Supply:    len(group.Assets),
				NftIds:    &nftIds,
			})
		}
	}

	spentAndLeasedGroups, err := GetSpentAndLeasedGroups(gotKeys)
	if err != nil {
		logrus.Infoln("GetSpentAndLeasedGroups", err)
		spentAndLeasedGroups = []NftGroup{}
	}

	for _, g := range spentAndLeasedGroups {
		groups = append(groups, &g)
	}
	return groups, nil
}

func PcListNonGroupNftAssets() ([]*ListAssetsResponse, error) {
	processed, err := ListAssetsProcessed(false, false, false)
	if err != nil {
		return nil, errors.Wrap(err, "ListAssetsProcessed")
	}
	var result []*ListAssetsResponse
	if processed == nil {
		return result, nil
	}
	for index, pr := range *processed {
		if pr.AssetGenesis.AssetType == "COLLECTIBLE" {
			result = append(result, &((*processed)[index]))
		}
	}
	resultFiltered := FilterListAssetsNullGroupKey2(result)
	return resultFiltered, nil
}

func PcGetSpentNftAssets() (*[]ListAssetsSimplifiedResponse, error) {
	response, err := GetSpentNftAssetsAndGetResponse()
	if err != nil {
		return response, errors.Wrap(err, "GetSpentNftAssetsAndGetResponse")
	}
	return response, nil
}

func PcQueryAddrs(assetId string) ([]*QueriedAddr, error) {
	addrRcv, err := rpcclient.AddrReceives()
	if err != nil {
		return nil, errors.Wrap(err, "rpcclient.AddrReceives")
	}
	addrMap := make(map[string]int)
	for _, events := range addrRcv.Events {
		addrMap[events.Addr.Encoded]++
	}
	_addrs, err := rpcclient.QueryAddr()
	if err != nil {
		return nil, errors.Wrap(err, "rpcclient.QueryDbAddr")
	}
	var addrs []*QueriedAddr
	for _, a := range _addrs.Addrs {
		if assetId != "" && assetId != hex.EncodeToString(a.AssetId) {
			continue
		}
		addrTemp := QueriedAddr{}
		addrTemp.GetData(a)
		addrTemp.ReceiveNum = addrMap[addrTemp.Encoded]
		addrs = append(addrs, &addrTemp)
	}
	return addrs, nil
}
