package api

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/vincent-petithory/dataurl"
	"github.com/wallet/models"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type (
	ByteSize float64
)

const (
	BaseTransactionByteSize = 170
)

type IssuanceHistoryInfo struct {
	IsFairLaunchIssuance bool   `json:"isFairLaunchIssuance"`
	AssetName            string `json:"asset_name"`
	AssetID              string `json:"asset_id"`
	ReservedTotal        int    `json:"reserved_total"`
	AssetType            int    `json:"asset_type"`
	IssuanceTime         int    `json:"issuance_time"`
	IssuanceAmount       int    `json:"issuance_amount"`
	State                int    `json:"state"`
}

// GetUserOwnIssuanceHistoryInfos
// @Description: Get User Own Issuance History Infos
func GetUserOwnIssuanceHistoryInfos(token string) string {
	result, err := GetAllUserOwnServerAndLocalTapdIssuanceHistoryInfos(token)
	if err != nil {
		LogError("", err)
		return MakeJsonErrorResult(GetAllUserOwnServerAndLocalTapdIssuanceHistoryInfosErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", result)
}

type GetIssuanceTransactionFeeResponse struct {
	Success bool    `json:"success"`
	Error   string  `json:"error"`
	Code    ErrCode `json:"code"`
	Data    int     `json:"data"`
}

func RequestToGetIssuanceTransactionFee(token string, feeRate int) (fee int, err error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := "http://" + serverDomainOrSocket + "/v1/fee/query/fair_launch/issuance?fee_rate=" + strconv.Itoa(feeRate)
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
	var response GetIssuanceTransactionFeeResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}
	if response.Error != "" {
		return 0, errors.New(response.Error)
	}
	return response.Data, nil
}

type GetMintTransactionFeeResponse struct {
	Success bool    `json:"success"`
	Error   string  `json:"error"`
	Code    ErrCode `json:"code"`
	Data    int     `json:"data"`
}

func RequestToGetMintTransactionFee(token string, feeRate int) (fee int, err error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := "http://" + serverDomainOrSocket + "/v1/fee/query/fair_launch/mint?fee_rate=" + strconv.Itoa(feeRate)
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
	var response GetMintTransactionFeeResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}
	if response.Error != "" {
		return 0, errors.New(response.Error)
	}
	return response.Data, nil
}

// GetIssuanceTransactionFee
// @Description: Get Issuance Transaction Fee
func GetIssuanceTransactionFee(token string, feeRate int) string {
	result, err := RequestToGetIssuanceTransactionFee(token, feeRate)
	if err != nil {
		LogError("", err)
		return MakeJsonErrorResult(GetIssuanceTransactionCalculatedFeeErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", result)
}

// GetMintTransactionFee
// @Description: Get Mint Transaction Fee
func GetMintTransactionFee(token string, feeRate int) string {
	result, err := RequestToGetMintTransactionFee(token, feeRate)
	if err != nil {
		LogError("", err)
		return MakeJsonErrorResult(GetMintTransactionCalculatedFeeErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", result)
}

func GetLocalIssuanceTransactionFee(feeRate int) string {
	result := int(float64(GetLocalIssuanceTransactionByteSize()) * float64(feeRate))
	return MakeJsonErrorResult(SUCCESS, "", result)
}

func GetLocalIssuanceTransactionByteSize() ByteSize {
	// TODO: need to complete
	byteSize := BaseTransactionByteSize * (1 + 0e0)
	return ByteSize(byteSize)
}

func GetIssuanceTransactionCalculatedFee(token string) (fee int, err error) {
	size := GetIssuanceTransactionByteSize()
	serverFeeRateResponse, err := GetServerFeeRate(token)
	if err != nil {
		LogError("", err)
		return 0, err
	}
	feeRate := serverFeeRateResponse.Data.SatPerB
	return feeRate*size + 3000, err
}

func GetMintTransactionCalculatedFee(token string, id int, number int) (fee int, err error) {
	size := int(GetMintTransactionByteSize())
	serverQueryMintResponse, err := GetServerQueryMint(token, id, number)
	if err != nil {
		LogError("", err)
		return 0, err
	}
	feeRate := serverQueryMintResponse.Data.CalculatedFeeRateSatPerB
	return feeRate*size + 1500, err
}

func GetIssuanceTransactionByteSize() int {
	// TODO: need to complete
	return int(float64(GetTapdMintAssetAndFinalizeTransactionByteSize()) + float64(GetTapdSendReservedAssetTransactionByteSize()))
}

func GetTapdMintAssetAndFinalizeTransactionByteSize() ByteSize {
	// TODO: need to complete
	byteSize := BaseTransactionByteSize * (1 + 0x1p-2)
	return ByteSize(byteSize)
}

func GetTapdSendReservedAssetTransactionByteSize() ByteSize {
	// TODO: need to complete
	byteSize := BaseTransactionByteSize * (1 + 0x1p-2)
	return ByteSize(byteSize)
}

func GetMintTransactionByteSize() ByteSize {
	// TODO: need to complete
	byteSize := BaseTransactionByteSize * (1 + 0e0)
	return ByteSize(byteSize)
}

type ServerOwnSetFairLaunchInfoResponse struct {
	Success bool                    `json:"success"`
	Error   string                  `json:"error"`
	Data    []models.FairLaunchInfo `json:"data"`
}

func GetServerOwnSetFairLaunchInfos(token string) (fairLaunchInfos *[]models.FairLaunchInfo, err error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := "http://" + serverDomainOrSocket + "/v1/fair_launch/query/own_set"
	client := &http.Client{}
	var jsonData []byte
	request, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	request.Header.Add("Authorization", "Bearer "+token)
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			LogError("", err)
		}
	}(response.Body)
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var ownSetFairLaunchInfos ServerOwnSetFairLaunchInfoResponse
	if err := json.Unmarshal(bodyBytes, &ownSetFairLaunchInfos); err != nil {
		return nil, err
	}
	return &ownSetFairLaunchInfos.Data, nil
}

func ProcessOwnSetFairLaunchResponseToIssuanceHistoryInfo(fairLaunchInfos *[]models.FairLaunchInfo) (*[]IssuanceHistoryInfo, error) {
	var err error
	var issuanceHistoryInfos []IssuanceHistoryInfo
	if fairLaunchInfos == nil {
		err = errors.New("fairLaunchInfos is null")
		LogError("", err)
		return nil, err
	}
	if len(*(fairLaunchInfos)) == 0 {
		//LogInfo("fairLaunchInfos length is zero")
		return &issuanceHistoryInfos, nil
	}
	for _, fairLaunchInfo := range *fairLaunchInfos {
		issuanceHistoryInfos = append(issuanceHistoryInfos, IssuanceHistoryInfo{
			IsFairLaunchIssuance: true,
			AssetName:            fairLaunchInfo.Name,
			AssetID:              fairLaunchInfo.AssetID,
			ReservedTotal:        fairLaunchInfo.ReserveTotal,
			AssetType:            int(fairLaunchInfo.AssetType),
			IssuanceTime:         fairLaunchInfo.SetTime,
			State:                int(fairLaunchInfo.State),
		})
	}
	return &issuanceHistoryInfos, nil

}

type ServerFeeRateResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Data    struct {
		SatPerKw int     `json:"sat_per_kw"`
		SatPerB  int     `json:"sat_per_b"`
		BtcPerKb float64 `json:"btc_per_kb"`
	} `json:"data"`
}

func GetServerFeeRate(token string) (*ServerFeeRateResponse, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := "http://" + serverDomainOrSocket + "/v1/fee/query/rate"
	client := &http.Client{}
	var jsonData []byte
	request, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", "Bearer "+token)
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			LogError("", err)
		}
	}(response.Body)
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var serverFeeRateResponse ServerFeeRateResponse
	if err = json.Unmarshal(bodyBytes, &serverFeeRateResponse); err != nil {
		return nil, err
	}
	return &serverFeeRateResponse, nil
}

type ServerQueryMintResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Data    struct {
		CalculatedFeeRateSatPerB  int  `json:"calculated_fee_rate_sat_per_b"`
		CalculatedFeeRateSatPerKw int  `json:"calculated_fee_rate_sat_per_kw"`
		InventoryAmount           int  `json:"inventory_amount"`
		IsMintAvailable           bool `json:"is_mint_available"`
	} `json:"data"`
}

func GetServerQueryMint(token string, id int, number int) (*ServerQueryMintResponse, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := "http://" + serverDomainOrSocket + "/v1/fair_launch/query/mint"
	client := &http.Client{}
	requestJson := struct {
		FairLaunchInfoId int `json:"fair_launch_info_id"`
		MintedNumber     int `json:"minted_number"`
	}{
		FairLaunchInfoId: id,
		MintedNumber:     number,
	}
	requestJsonBytes, _ := json.Marshal(requestJson)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestJsonBytes))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", "Bearer "+token)
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			LogError("", err)
		}
	}(response.Body)
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var serverQueryMintResponse ServerQueryMintResponse
	if err = json.Unmarshal(bodyBytes, &serverQueryMintResponse); err != nil {
		return nil, err
	}
	return &serverQueryMintResponse, nil
}

// GetServerIssuanceHistoryInfos
// @Description: Get Server Issuance History Info
func GetServerIssuanceHistoryInfos(token string) (*[]IssuanceHistoryInfo, error) {
	fairLaunchInfos, err := GetServerOwnSetFairLaunchInfos(token)
	if err != nil {
		LogError("", err)
		return nil, err
	}
	issuanceHistoryInfos, err := ProcessOwnSetFairLaunchResponseToIssuanceHistoryInfo(fairLaunchInfos)
	if err != nil {
		LogError("", err)
		return nil, err
	}
	return issuanceHistoryInfos, nil
}

func GetLocalTapdIssuanceHistoryInfos() (*[]IssuanceHistoryInfo, error) {
	var err error
	var issuanceHistoryInfos []IssuanceHistoryInfo
	batchs, err := ListBatchesAndGetResponse()
	if err != nil {
		LogError("", err)
		return nil, err
	}
	transactions, err := GetTransactionsAndGetResponse()
	if err != nil {
		LogError("", err)
		return nil, err
	}
	listAssetResponse, err := ListAssetAndGetResponse()
	if err != nil {
		LogError("", err)
		return nil, err
	}
	var timestamp int
	var assetId string
	var outpoint string
	var assets *[]ListAssetResponse
	var leave *universerpc.AssetLeafResponse
	for _, batch := range (*batchs).Batches {
		var transaction *lnrpc.Transaction
		transaction, err = GetTransactionByBatchTxid(transactions, batch.Batch.BatchTxid)
		// transaction not found
		if err != nil {
			//LogError("", err)
			continue
			//	@dev: Do not return
		}
		timestamp = int(transaction.TimeStamp)
		outpoint = transaction.PreviousOutpoints[0].Outpoint
		assets, err = GetAssetsByOutpointWithListAssetResponse(listAssetResponse, outpoint)
		if err != nil {
			return nil, err
		}
		for _, asset := range *assets {
			assetId = asset.AssetGenesis.AssetID
			leave, err = AssetLeavesIssuance(assetId)
			if err != nil {
				return nil, err
			}
			leaveAsset := leave.Leaves[0].Asset
			issuanceHistoryInfos = append(issuanceHistoryInfos, IssuanceHistoryInfo{
				IsFairLaunchIssuance: false,
				AssetName:            asset.AssetGenesis.Name,
				AssetID:              assetId,
				AssetType:            asset.AssetGenesis.AssetType,
				IssuanceTime:         timestamp,
				IssuanceAmount:       int(leaveAsset.Amount),
				State:                int(batch.Batch.State),
			})
		}
	}
	return &issuanceHistoryInfos, nil
}

// GetAllUserOwnServerAndLocalTapdIssuanceHistoryInfos
// @Description: Get User Own Issuance History Infos
func GetAllUserOwnServerAndLocalTapdIssuanceHistoryInfos(token string) (*[]IssuanceHistoryInfo, error) {
	var issuanceHistoryInfos []IssuanceHistoryInfo
	serverResult, err := GetServerIssuanceHistoryInfos(token)
	if err != nil {
		LogError("", err)
		return nil, err
	}
	localTapdResult, err := GetLocalTapdIssuanceHistoryInfos()
	if err != nil {
		LogError("", err)
		return nil, err
	}
	issuanceHistoryInfos = append(issuanceHistoryInfos, *serverResult...)
	issuanceHistoryInfos = append(issuanceHistoryInfos, *localTapdResult...)
	return &issuanceHistoryInfos, nil
}

func GetTimestampByBatchTxidWithGetTransactionsResponse(transactionDetails *lnrpc.TransactionDetails, batchTxid string) (timestamp int, err error) {
	for _, transaction := range transactionDetails.Transactions {
		if batchTxid == transaction.TxHash {
			return int(transaction.TimeStamp), nil
		}
	}
	err = errors.New("transaction not found")
	return 0, err
}

func GetTransactionByBatchTxid(transactionDetails *lnrpc.TransactionDetails, batchTxid string) (transaction *lnrpc.Transaction, err error) {
	for _, transaction = range transactionDetails.Transactions {
		if batchTxid == transaction.TxHash && transaction.Label == "tapd-asset-minting" {
			return transaction, nil
		}
	}
	err = errors.New("transaction not found")
	return nil, err
}

func GetAssetIdByBatchTxidWithListAssetResponse(listAssetResponse *taprpc.ListAssetResponse, batchTxid string) (assetId string, err error) {
	for _, asset := range listAssetResponse.Assets {
		tx, _ := getTransactionAndIndexByOutpoint(asset.ChainAnchor.AnchorOutpoint)
		if batchTxid == tx {
			return hex.EncodeToString(asset.AssetGenesis.AssetId), nil
		}
	}
	err = errors.New("asset not found")
	return "", err
}

// GetAssetIdByOutpointAndNameWithListAssetResponse
// @dev: may be deprecated
func GetAssetIdByOutpointAndNameWithListAssetResponse(listAssetResponse *taprpc.ListAssetResponse, outpoint string, name string) (assetId string, err error) {
	for _, asset := range listAssetResponse.Assets {
		if outpoint == asset.AssetGenesis.GenesisPoint && name == asset.AssetGenesis.Name {
			return hex.EncodeToString(asset.AssetGenesis.AssetId), nil
		}
	}
	err = errors.New("asset not found")
	return "", err
}

func GetAssetsByOutpointWithListAssetResponse(listAssetResponse *taprpc.ListAssetResponse, outpoint string) (*[]ListAssetResponse, error) {
	var assets []ListAssetResponse
	//var err error
	isAssetIdExist := make(map[string]bool)
	for _, asset := range listAssetResponse.Assets {
		if outpoint == asset.AssetGenesis.GenesisPoint {
			assetId := hex.EncodeToString(asset.AssetGenesis.AssetId)
			if !isAssetIdExist[assetId] {
				isAssetIdExist[assetId] = true
				assets = append(assets, ListAssetResponse{
					Version: asset.Version.String(),
					AssetGenesis: AssetGenesisStruct{
						GenesisPoint: asset.AssetGenesis.GenesisPoint,
						Name:         asset.AssetGenesis.Name,
						MetaHash:     hex.EncodeToString(asset.AssetGenesis.MetaHash),
						AssetID:      hex.EncodeToString(asset.AssetGenesis.AssetId),
						AssetType:    int(asset.AssetGenesis.AssetType),
						OutputIndex:  int(asset.AssetGenesis.OutputIndex),
						Version:      int(asset.Version),
					},
					Amount:           int(asset.Amount),
					LockTime:         int(asset.LockTime),
					RelativeLockTime: int(asset.RelativeLockTime),
					ScriptVersion:    int(asset.ScriptVersion),
					ScriptKey:        hex.EncodeToString(asset.ScriptKey),
					ScriptKeyIsLocal: asset.ScriptKeyIsLocal,
					ChainAnchor: ChainAnchorStruct{
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
		}
	}
	return &assets, nil
}

// GetImageByImageData
// @Description: Get Image By Image Data
func GetImageByImageData(imageData string) []byte {
	if imageData == "" {
		return nil
	}
	dataUrl, err := dataurl.DecodeString(imageData)
	if err != nil {
		return nil
	}
	ContentType := dataUrl.MediaType.ContentType()
	datatype := strings.Split(ContentType, "/")
	if datatype[0] != "image" {
		fmt.Println("is not image dataurl")
		return nil
	}
	return dataUrl.Data
}

type FairLaunchFollowSetRequest struct {
	FairLaunchInfoId int    `json:"fair_launch_info_id"`
	AssetId          string `json:"asset_id" gorm:"type:varchar(255)"`
	DeviceId         string `json:"device_id" gorm:"type:varchar(255)"`
}

func PostToSetFollowFairLaunchInfo(token string, fairLaunchFollowSetRequest *FairLaunchFollowSetRequest) (*JsonResult, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := "http://" + serverDomainOrSocket + "/fair_launch_follow/follow"
	requestJsonBytes, err := json.Marshal(fairLaunchFollowSetRequest)
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

func PostToSetUnfollowFairLaunchInfo(token string, assetId string) (*JsonResult, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := "http://" + serverDomainOrSocket + "/fair_launch_follow/unfollow/asset_id/" + assetId
	requestJsonBytes, err := json.Marshal(nil)
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

// FollowFairLaunchAsset
// @Description: Follow fair launch asset
func FollowFairLaunchAsset(token string, fairLaunchInfoId int, assetId string, deviceId string) string {
	_, err := PostToSetFollowFairLaunchInfo(token, &FairLaunchFollowSetRequest{
		FairLaunchInfoId: fairLaunchInfoId,
		AssetId:          assetId,
		DeviceId:         deviceId,
	})
	if err != nil {
		return MakeJsonErrorResult(PostToSetFollowFairLaunchInfoErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, assetId)
}

// UnfollowFairLaunchAsset
// @Description: Unfollow fair launch asset
func UnfollowFairLaunchAsset(token string, assetId string) string {
	_, err := PostToSetUnfollowFairLaunchInfo(token, assetId)
	if err != nil {
		return MakeJsonErrorResult(PostToSetUnfollowFairLaunchInfoErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, assetId)
}

type QueryIsFairLaunchFollowedResponse struct {
	Success bool    `json:"success"`
	Error   string  `json:"error"`
	Code    ErrCode `json:"code"`
	Data    bool    `json:"data"`
}

func RequestToQueryIsFairLaunchFollowed(token string, assetId string) (bool, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	url := "http://" + serverDomainOrSocket + "/fair_launch_follow/query/user/is_followed/asset_id/" + assetId
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return false, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return false, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}
	var response QueryIsFairLaunchFollowedResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, err
	}
	if response.Error != "" {
		return false, errors.New(response.Error)
	}
	return response.Data, nil

}

// QueryIsFairLaunchFollowed
// @Description: Query is fair launch followed
func QueryIsFairLaunchFollowed(token string, assetId string) string {
	response, err := RequestToQueryIsFairLaunchFollowed(token, assetId)
	if err != nil {
		return MakeJsonErrorResult(RequestToQueryIsFairLaunchFollowedErr, err.Error(), false)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, response)
}
