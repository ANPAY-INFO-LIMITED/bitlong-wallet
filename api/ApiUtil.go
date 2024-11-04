package api

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/lightninglabs/taproot-assets/tapdb"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/protobuf/proto"
	"log"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type ErrCode int

// Err type:Normal
const (
	DefaultErr   ErrCode = -1
	SUCCESS      ErrCode = 200
	NotFoundData ErrCode = iota + 299
	RequestError
)

const (
	GetBtcTransferOutInfosErr ErrCode = iota + 501
	ListTransfersAndGetProcessedResponseErr
	PostToSetAssetTransferErr
	PostToGetAssetTransferAndGetResponseErr
	BatchTxidToAssetIdErr
	AddrReceivesAndGetEventsErr
	PostToSetAddrReceivesEventsErr
	PostToGetAddrReceivesEventsErr
	jsonAddrsToAddrSliceErr
	DecodeAddrErr
	sendAssetsErr
	UploadBatchTransfersErr
	PostToGetBatchTransfersErr
	PostToSetAssetAddrErr
	PostToGetAssetAddrErr
	ListUtxosAndGetResponseErr
	ListUnspentAndGetResponseErr
	ListNftAssetsAndGetResponseErr
	PostToSetAssetLockErr
	PostToGetAssetLockErr
	IsTokenValidErr
	JsonUnmarshalErr
	ListBalancesAndProcessErr
	PostToSetAssetBalanceInfoErr
	FeeRateExceedMaxErr
	QueryAllAddrAndGetResponseErr
	UpdateAllAddrByAccountWithAddressesErr
	PostToGetAssetTransferByAssetIdAndGetResponseErr
	QueryAssetTransferSimplifiedErr
	RequestToGetNonZeroAssetBalanceErr
	GetZeroBalanceAssetBalanceSliceErr
	GetAssetHolderNumberByAssetBalancesInfoErr
	GetAssetHolderBalanceByAssetBalancesInfoErr
	AddrReceivesErr
	BurnAssetErr
	fetchAssetMetaErr
	GetInfoErr
	GetConnectionErr
	syncUniverseErr
	ProcessListAllAssetsSimplifiedErr
	allAssetBalancesErr
	allAssetGroupBalancesErr
	assetKeysTransferErr
	AssetLeavesSpecifiedErr
	assetLeavesIssuanceErr
	DecodeRawProofStringErr
	allAssetListErr
	GetAssetHoldInfosIncludeSpentErr
	GetAssetHoldInfosExcludeSpentErr
	GetAssetTransactionInfosErr
	GetTimeForManagedUtxoByBitcoindErr
	subServerStatusErr
	NewAddressP2trErr
	NewAddressP2wkhErr
	NewAddressNp2wkhErr
	CreateOrUpdateAddrErr
	ReadAddrErr
	DeleteAddrErr
	AllAddressesErr
	ListAddressesErr
	GetAccountWithAddressesErr
	UnmarshalErr
	resultIsNotSuccessErr
	GetAllAccountsErr
	InvalidAddressTypeErr
	GetBlockErr
	GetBlockHashErr
	getWalletBalanceErr
	ProcessGetWalletBalanceResultErr
	getInfoOfLndErr
	DecodePayReqErr
	ListChannelsErr
	ListInvoicesErr
	PendingChannelsErr
	ClosedChannelsErr
	NoFindChannelErr
	sendCoinsErr
	SendPaymentSyncErr
	AddrsLenZeroErr
	sendManyErr
	TrackPaymentV2Err
	streamRecvInfoErr
	streamRecvErr
	listAccountsErr
	AccountNotFoundErr
	BumpFeeErr
	HttpGetErr
	GetAddressTransferOutErr
	GetAddressTransactionsErr
	listAssetsErr
	responseNotSuccessErr
	assetNotFoundErr
	ListGroupsErr
	ListTransfersErr
	NewAddrErr
	QueryAddrErr
	listBalancesErr
	assetLeafKeysErr
	ListBatchesAndGetResponseErr
	assetLeavesErr
	GetTransactionsAndGetResponseErr
	GetAssetInfoErr
	ListAssetsProcessedErr
	responseAssetKeysZeroErr
	responseLeavesNullErr
	QueryAssetRootsErr
	blobLenZeroErr
	DecodeProofErr
	clientInfoErr
	queryAssetRootErr
	queryAssetStatsErr
	GetAllUserOwnServerAndLocalTapdIssuanceHistoryInfosErr
	GetIssuanceTransactionCalculatedFeeErr
	GetMintTransactionCalculatedFeeErr
	FinalizeBatchErr
	DecodeStringErr
	MintAssetErr
	getTransactionByMempoolErr
	deliverIssuanceProofErr
	deliverProofErr
	receiveProofErr
	readProofErr
	queryAssetProofsErr
	UploadAssetBurnErr
	GetAssetBurnTotalAmountByAssetIdErr
	GetOwnFairLaunchInfoIssuedSimplifiedAndExecuteMintReservedErr
	PostToSetFollowFairLaunchInfoErr
	PostToSetUnfollowFairLaunchInfoErr
	RequestToQueryIsFairLaunchFollowedErr
	ListBatchesAndPostToSetAssetLocalMintHistoriesErr
	ListUtxosAndPostToSetAssetManagedUtxosErr
	GetWalletBalanceCalculatedTotalValueErr
	UploadLogFileAndGetJsonResultErr
	GetAccountAssetBalanceByAssetIdAndGetResponseErr
	GetAccountAssetTransferByAssetIdAndGetResponseErr
	GetAssetHolderBalanceWithPageSizeAndPageNumberErr
	GetAccountAssetTransferPageNumberByPageSizeErr
	GetAccountAssetTransferWithPageSizeAndPageNumberErr
	GetAccountAssetBalancePageNumberByPageSizeErr
	GetAccountAssetBalanceWithPageSizeAndPageNumberErr
	GetAssetManagedUtxoWithPageSizeAndPageNumberErr
	GetAssetManagedUtxoPageNumberByPageSizeErr
	GetGroupFirstAssetMetaAndGetResponseErr
	SetGroupFirstAssetMetaAndGetResponseErr
	GetGroupFirstAssetIdAndGetResponseErr
	QueryAssetTransferSimplifiedOfAllNftErr
	GetDeliverProofNeedInfoAndGetResponseErr
	GetNftTransferByAssetIdAndGetResponseErr
	GetSpentNftAssetsAndGetResponseErr
	SchnorrSignErr
	SchnorrVerifyErr
	GetAccountAssetBalanceUserHoldTotalAmountByAssetIdErr
	LndGetInfoAndGetResponseErr
	GetZeroAmountAssetListSliceErr
	PostToSetAssetListInfoErr
)

var ErrCodeMapInfo = map[ErrCode]string{
	GetBtcTransferOutInfosErr:                        "获取BTC转出记录错误",
	ListTransfersAndGetProcessedResponseErr:          "列出转账记录并获取处理的响应错误",
	PostToSetAssetTransferErr:                        "请求发送资产转账记录错误",
	PostToGetAssetTransferAndGetResponseErr:          "请求获取资产转账记录和响应错误",
	BatchTxidToAssetIdErr:                            "批量交易ID转换资产ID错误",
	AddrReceivesAndGetEventsErr:                      "资产地址接收并获取事件错误",
	PostToSetAddrReceivesEventsErr:                   "请求发送资产接收事件错误",
	PostToGetAddrReceivesEventsErr:                   "请求获取资产接收事件错误",
	jsonAddrsToAddrSliceErr:                          "资产地址数组JSON字符串转换资产地址切片错误",
	DecodeAddrErr:                                    "解码资产地址错误",
	sendAssetsErr:                                    "发送资产错误",
	UploadBatchTransfersErr:                          "上传批量转账记录错误",
	PostToGetBatchTransfersErr:                       "请求获取批量转账记录错误",
	PostToSetAssetAddrErr:                            "请求发送资产地址记录错误",
	PostToGetAssetAddrErr:                            "请求获取资产地址记录错误",
	ListUtxosAndGetResponseErr:                       "列出资产UTXO并获取响应错误",
	ListUnspentAndGetResponseErr:                     "列出BTC的UTXO并获取相应错误",
	ListNftAssetsAndGetResponseErr:                   "列出NFT资产并获取响应错误",
	PostToSetAssetLockErr:                            "请求发送资产锁定信息错误",
	PostToGetAssetLockErr:                            "请求获取资产锁定信息错误",
	IsTokenValidErr:                                  "Token无效错误",
	JsonUnmarshalErr:                                 "JSON解码错误",
	ListBalancesAndProcessErr:                        "列出BTC余额并获取处理结果错误",
	PostToSetAssetBalanceInfoErr:                     "请求发送资产余额信息错误",
	FeeRateExceedMaxErr:                              "费率超出最大值错误",
	QueryAllAddrAndGetResponseErr:                    "查询所有资产地址并获取响应错误",
	UpdateAllAddrByAccountWithAddressesErr:           "通过账户更新所有资产地址错误",
	PostToGetAssetTransferByAssetIdAndGetResponseErr: "请求通过资产ID获取资产转账记录并获取响应错误",
	QueryAssetTransferSimplifiedErr:                  "查询简化资产转账记录错误",
	RequestToGetNonZeroAssetBalanceErr:               "请求获取非零资产余额信息错误",
	GetZeroBalanceAssetBalanceSliceErr:               "获取零余额资产余额信息切片错误",
	GetAssetHolderNumberByAssetBalancesInfoErr:       "通过资产余额信息获取资产持有人数量错误",
	GetAssetHolderBalanceByAssetBalancesInfoErr:      "通过资产余额信息获取资产持有人持有信息错误",
	AddrReceivesErr:                                  "资产接收记录错误",
	BurnAssetErr:                                     "销毁资产错误",
	fetchAssetMetaErr:                                "提取资产元数据错误",
	GetInfoErr:                                       "获取信息错误",
	GetConnectionErr:                                 "获取连接错误",
	syncUniverseErr:                                  "同步宇宙错误",
	ProcessListAllAssetsSimplifiedErr:                "处理列出的所有简化资产信息错误",
	allAssetBalancesErr:                              "获取所有资产余额信息错误",
	allAssetGroupBalancesErr:                         "获取所有资产组余额信息错误",
	assetKeysTransferErr:                             "获取资产转账Key错误",
	AssetLeavesSpecifiedErr:                          "获取指定类型资产叶子错误",
	assetLeavesIssuanceErr:                           "获取发行资产叶子错误",
	DecodeRawProofStringErr:                          "解码原始证明字符串错误",
	allAssetListErr:                                  "获取列出所有资产信息错误",
	GetAssetHoldInfosIncludeSpentErr:                 "获取包含已花费的资产持有信息错误",
	GetAssetHoldInfosExcludeSpentErr:                 "获取不包含已花费的资产持有信息错误",
	GetAssetTransactionInfosErr:                      "获取资产交易信息错误",
	GetTimeForManagedUtxoByBitcoindErr:               "通过Bitcoind为BTC的UTXO获取时间错误",
	subServerStatusErr:                               "获取子服务状态错误",
	NewAddressP2trErr:                                "生成P2TR地址错误",
	NewAddressP2wkhErr:                               "生成P2WKH地址错误",
	NewAddressNp2wkhErr:                              "生成NP2WKH地址错误",
	CreateOrUpdateAddrErr:                            "创建或更新资产地址信息错误",
	ReadAddrErr:                                      "读取资产地址错误",
	DeleteAddrErr:                                    "删除资产地址错误",
	AllAddressesErr:                                  "获取所有资产地址错误",
	ListAddressesErr:                                 "列出资产地址错误",
	GetAccountWithAddressesErr:                       "获取账户与地址错误",
	UnmarshalErr:                                     "解码错误",
	resultIsNotSuccessErr:                            "结果未成功错误",
	GetAllAccountsErr:                                "获取所有账户错误",
	InvalidAddressTypeErr:                            "无效的资产地址类型错误",
	GetBlockErr:                                      "获取区块错误",
	GetBlockHashErr:                                  "获取区块哈希值错误",
	getWalletBalanceErr:                              "获取钱包余额信息错误",
	ProcessGetWalletBalanceResultErr:                 "处理获取钱包余额信息的结果错误",
	getInfoOfLndErr:                                  "获取LND的信息错误",
	DecodePayReqErr:                                  "解码支付请求错误",
	ListChannelsErr:                                  "获取列出通道错误",
	ListInvoicesErr:                                  "获取列出发票错误",
	PendingChannelsErr:                               "获取等待中的通道错误",
	ClosedChannelsErr:                                "关闭通道错误",
	NoFindChannelErr:                                 "没有找到通道错误",
	sendCoinsErr:                                     "发送BTC币错误",
	SendPaymentSyncErr:                               "同步发起支付错误",
	AddrsLenZeroErr:                                  "资产地址为零错误",
	sendManyErr:                                      "发送多笔支付错误",
	TrackPaymentV2Err:                                "跟踪支付V2错误",
	streamRecvInfoErr:                                "流接收信息错误",
	streamRecvErr:                                    "流接收错误",
	listAccountsErr:                                  "获取列出账户错误",
	AccountNotFoundErr:                               "账户未找到错误",
	BumpFeeErr:                                       "碰撞费率错误",
	HttpGetErr:                                       "HTTP的GET请求错误",
	GetAddressTransferOutErr:                         "获取BTC地址转出错误",
	GetAddressTransactionsErr:                        "通过Mempool获取BTC地址交易记录错误",
	listAssetsErr:                                    "获取列出资产错误",
	responseNotSuccessErr:                            "响应未成功错误",
	assetNotFoundErr:                                 "资产未找到错误",
	ListGroupsErr:                                    "获取列出资产组错误",
	ListTransfersErr:                                 "获取列出BTC转账记录错误",
	NewAddrErr:                                       "新资产地址错误",
	QueryAddrErr:                                     "查询资产地址错误",
	listBalancesErr:                                  "获取列出BTC余额错误",
	assetLeafKeysErr:                                 "获取资产叶子Key错误",
	ListBatchesAndGetResponseErr:                     "获取列出批次并获取响应错误",
	assetLeavesErr:                                   "获取资产叶子错误",
	GetTransactionsAndGetResponseErr:                 "获取交易记录和响应错误",
	GetAssetInfoErr:                                  "获取资产信息错误",
	ListAssetsProcessedErr:                           "获取列出已处理的资产信息错误",
	responseAssetKeysZeroErr:                         "资产Key响应长度为零错误",
	responseLeavesNullErr:                            "资产叶子响应为空错误",
	QueryAssetRootsErr:                               "查询资产根错误",
	blobLenZeroErr:                                   "Blob长度为零错误",
	DecodeProofErr:                                   "解码证明错误",
	clientInfoErr:                                    "客户端信息错误",
	queryAssetRootErr:                                "查询资产根错误",
	queryAssetStatsErr:                               "查询资产统计错误",
	GetAllUserOwnServerAndLocalTapdIssuanceHistoryInfosErr: "获取用户所有服务器和本地发行历史记录错误",
	GetIssuanceTransactionCalculatedFeeErr:                 "获取发行交易计算费用错误",
	GetMintTransactionCalculatedFeeErr:                     "获取铸造交易计算费用错误",
	FinalizeBatchErr:                                       "提交批次错误",
	DecodeStringErr:                                        "解码字符串错误",
	MintAssetErr:                                           "本地铸造资产错误",
	getTransactionByMempoolErr:                             "通过Mempool获取交易信息错误",
	deliverIssuanceProofErr:                                "递送资产发行证明错误",
	deliverProofErr:                                        "递送资产证明错误",
	receiveProofErr:                                        "接收资产证明错误",
	readProofErr:                                           "读取资产证明错误",
	queryAssetProofsErr:                                    "查询资产证明错误",
	UploadAssetBurnErr:                                     "上传资产销毁信息错误",
	GetAssetBurnTotalAmountByAssetIdErr:                    "通过资产ID获取资产销毁总量错误",
	GetOwnFairLaunchInfoIssuedSimplifiedAndExecuteMintReservedErr: "获取简化的自己发行的公平发射资产信息并执行取回保留部分错误",
	PostToSetFollowFairLaunchInfoErr:                              "请求关注公平发射信息错误",
	PostToSetUnfollowFairLaunchInfoErr:                            "请求取消关注公平发射信息错误",
	RequestToQueryIsFairLaunchFollowedErr:                         "请求查询是否已关注公平发射信息错误",
	ListBatchesAndPostToSetAssetLocalMintHistoriesErr:             "获取列出批次并请求发送资产本地铸造历史记录错误",
	ListUtxosAndPostToSetAssetManagedUtxosErr:                     "请求列出UTXO并请求发送资产UTXO信息错误",
	GetWalletBalanceCalculatedTotalValueErr:                       "获取钱包余额计算总价值错误",
	UploadLogFileAndGetJsonResultErr:                              "上传日志文件并获取JSON结果响应错误",
	GetAccountAssetBalanceByAssetIdAndGetResponseErr:              "通过资产ID获取账户资产余额信息并获取响应错误",
	GetAccountAssetTransferByAssetIdAndGetResponseErr:             "通过资产ID获取账户资产转账记录并获取响应错误",
	GetAssetHolderBalanceWithPageSizeAndPageNumberErr:             "通过页面大小和页号获取资产持有信息错误",
	GetAccountAssetTransferPageNumberByPageSizeErr:                "通过页面大小获取账户资产转账记录页数错误",
	GetAccountAssetTransferWithPageSizeAndPageNumberErr:           "通过页面大小和页号获取资产转账信息错误",
	GetAccountAssetBalancePageNumberByPageSizeErr:                 "通过页面大小获取账户资产余额信息页数错误",
	GetAccountAssetBalanceWithPageSizeAndPageNumberErr:            "通过页面大小和页号获取账户资产余额信息错误",
	GetAssetManagedUtxoWithPageSizeAndPageNumberErr:               "通过页面大小和页号获取资产UTXO信息错误",
	GetAssetManagedUtxoPageNumberByPageSizeErr:                    "通过页面大小获取资产UTXO信息页数错误",
	GetGroupFirstAssetMetaAndGetResponseErr:                       "获取资产组的首个资产元数据并获取响应错误",
	SetGroupFirstAssetMetaAndGetResponseErr:                       "请求上传资产组首个资产的元数据并获取响应错误",
	GetGroupFirstAssetIdAndGetResponseErr:                         "获取资产组的首个资产ID并获取响应错误",
	QueryAssetTransferSimplifiedOfAllNftErr:                       "查询所有NFT简化资产转账记录错误",
	GetDeliverProofNeedInfoAndGetResponseErr:                      "获取发送证明文件操作所需前置信息并获取响应错误",
	GetNftTransferByAssetIdAndGetResponseErr:                      "通过资产ID查询NFT转账记录并获取响应错误",
	GetSpentNftAssetsAndGetResponseErr:                            "获取已转出的NFT资产并获取响应错误",
	SchnorrSignErr:                                                "Schnorr签名错误",
	SchnorrVerifyErr:                                              "Schnorr验证错误",
	GetAccountAssetBalanceUserHoldTotalAmountByAssetIdErr:         "通过资产ID获取托管资产余额用户持有总量错误",
	LndGetInfoAndGetResponseErr:                                   "Lnd获取信息并获取相应错误",
	GetZeroAmountAssetListSliceErr:                                "获取零余额资产列表信息切片错误",
	PostToSetAssetListInfoErr:                                     "请求发送资产列表信息错误",
}

func GetIntErrCodeString(intErrCode int) string {
	return GetErrCodeString(ErrCode(intErrCode))
}

func GetErrCodeString(errCode ErrCode) string {
	return errCode.Error()
}

func (ec ErrCode) Error() string {
	switch {
	case errors.Is(ec, NotFoundData):
		return "not found Data"
	case errors.Is(ec, SUCCESS):
		return ""
	case errors.Is(ec, RequestError):
		return "request error"
	default:
	}
	info, ok := ErrCodeMapInfo[ec]
	if !ok {
		info = "[无错误信息]"
	} else {
		info += strconv.Itoa(int(ec))
	}
	return info
}

var (
	SuccessError = SUCCESS.Error()
)

type JsonResult struct {
	// Deprecated: Use Code instead
	Success bool    `json:"success"`
	Error   string  `json:"error"`
	Code    ErrCode `json:"code"`
	Data    any     `json:"data"`
}

// Deprecated: Use MakeJsonErrorResult instead
func MakeJsonResult(success bool, error string, data any) string {
	jsr := JsonResult{
		Success: success,
		Error:   error,
		Code:    -1,
		Data:    data,
	}

	if success {
		jsr.Code = SUCCESS
	} else {
		jsr.Code = DefaultErr
	}
	jstr, err := json.Marshal(jsr)
	if err != nil {
		return MakeJsonResult(false, err.Error(), nil)
	}
	return string(jstr)
}

func MakeJsonErrorResult(code ErrCode, errorString string, data any) string {
	jsr := JsonResult{
		Error: errorString,
		Code:  code,
		Data:  data,
	}
	if errors.Is(code, SUCCESS) {
		jsr.Success = true
	} else {
		jsr.Success = false
	}
	jstr, err := json.Marshal(jsr)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	return string(jstr)
}

func LnMarshalRespString(resp proto.Message) string {
	jsonBytes, err := lnrpc.ProtoJSONMarshalOpts.Marshal(resp)
	if err != nil {
		fmt.Printf("%s unable to decode response: %v\n", GetTimeNow(), err)
		return ""
	}
	return string(jsonBytes)
}

func TapMarshalRespString(resp proto.Message) string {
	jsonBytes, err := taprpc.ProtoJSONMarshalOpts.Marshal(resp)
	if err != nil {
		fmt.Printf("%s unable to decode response: %v\n", GetTimeNow(), err)
		return ""
	}
	return string(jsonBytes)
}

func B64DecodeToHex(s string) string {
	byte1, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "DECODE_ERROR"
	}
	return hex.EncodeToString(byte1)
}

func GetTimeNow() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

func GetTimeSuffixString() string {
	return time.Now().Format("20060102150405")
}

func RoundToDecimalPlace(number float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Round(number*shift) / shift
}

func GetEnv(key string, filename ...string) string {
	err := godotenv.Load(filename...)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	value := os.Getenv(key)
	return value
}

func ToBTC(sat int) float64 {
	return float64(sat / 1e8)
}

func ToSat(btc float64) int {
	return int(btc * 1e8)
}

func LogInfo(info string) {
	fmt.Printf("%s %s\n", GetTimeNow(), info)
}

func LogInfos(infos ...string) {
	var info string
	for i, _info := range infos {
		if i != 0 {
			info += " "
		}
		info += _info
	}
	fmt.Printf("%s %s\n", GetTimeNow(), info)
}

func LogError(description string, err error) {
	fmt.Printf("%s %s :%v\n", GetTimeNow(), description, err)
}

// FeeRateBtcPerKbToSatPerKw
// @Description: BTC/Kb to sat/kw
// 1 sat/vB = 0.25 sat/wu
// https://bitcoin.stackexchange.com/questions/106333/different-fee-rate-units-sat-vb-sat-perkw-sat-perkb
func FeeRateBtcPerKbToSatPerKw(btcPerKb float64) (satPerKw int) {
	// @dev: 1 BTC/kB = 1e8 sat/kB 1e5 sat/B = 0.25e5 sat/w = 0.25e8 sat/kw
	return int(0.25e8 * btcPerKb)
}

// FeeRateBtcPerKbToSatPerB
// @Description: BTC/Kb to sat/b
// @param btcPerKb
// @return satPerB
func FeeRateBtcPerKbToSatPerB(btcPerKb float64) (satPerB int) {
	return int(1e5 * btcPerKb)
}

// FeeRateSatPerKwToBtcPerKb
// @Description: sat/kw to BTC/Kb
// @param feeRateSatPerKw
// @return feeRateBtcPerKb
func FeeRateSatPerKwToBtcPerKb(feeRateSatPerKw int) (feeRateBtcPerKb float64) {
	return RoundToDecimalPlace(float64(feeRateSatPerKw)/0.25e8, 8)
}

// FeeRateSatPerKwToSatPerB
// @Description: sat/kw to sat/b
func FeeRateSatPerKwToSatPerB(feeRateSatPerKw int) (feeRateSatPerB int) {
	return int(math.Ceil(float64(feeRateSatPerKw) * 4 / 1000))
}

// FeeRateSatPerBToBtcPerKb
// @Description: sat/b to BTC/Kb
func FeeRateSatPerBToBtcPerKb(feeRateSatPerB int) (feeRateBtcPerKb float64) {
	return RoundToDecimalPlace(math.Ceil(float64(feeRateSatPerB)/100000), 8)
}

// FeeRateSatPerBToSatPerKw
// @Description: sat/b to sat/kw
func FeeRateSatPerBToSatPerKw(feeRateSatPerB int) (feeRateSatPerKw int) {
	return int(math.Ceil(float64(feeRateSatPerB) * 1000 / 4))
}

func ValueJsonString(value any) string {
	resultJSON, err := json.MarshalIndent(value, "", "\t")
	if err != nil {
		LogError("MarshalIndent error", err)
		return ""
	}
	return string(resultJSON)
}

func AppendErrorInfo(err error, info string) error {
	if err == nil {
		err = errors.New("[nil err]")
	}
	return errors.New(err.Error() + ";" + info)
}

func AppendError(e error) func(error) error {
	return func(err error) error {
		if e == nil {
			e = errors.New("")
		}
		if err == nil {
			return e
		}
		if e.Error() == "" {
			e = err
			return e
		}
		e = errors.New(e.Error() + "; " + err.Error())
		return e
	}
}

func AppendInfo(s string) func(string) string {
	return func(info string) string {
		if info == "" {
			return s
		}
		if s == "" {
			s = info
			return s
		}
		s = s + "; " + info
		return s
	}
}

func InfoAppendError(i string) func(error) error {
	e := errors.New(i)
	return func(err error) error {
		if err == nil {
			return e
		}
		if e.Error() == "" {
			e = err
			return e
		}
		e = errors.New(e.Error() + "; " + err.Error())
		return e
	}
}

func ErrorAppendInfo(e error) func(string) error {
	return func(info string) error {
		if e == nil {
			e = errors.New("")
		}
		if info == "" {
			return e
		}
		if e.Error() == "" {
			e = errors.New(info)
			return e
		}
		info = e.Error() + "; " + info
		e = errors.New(info)
		return e
	}
}

func IsPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func IsHexString(s string) bool {
	_, err := hex.DecodeString(s)
	return err == nil
}

func SwapValue[T any](a *T, b *T) {
	temp := *a
	*a = *b
	*b = temp
}

func SwapInt(a *int, b *int) {
	*a ^= *b
	*b ^= *a
	*a ^= *b
}

func ToLowerWords(s string) string {
	var result strings.Builder
	for i, char := range s {
		if i > 0 && char >= 'A' && char <= 'Z' {
			temp := result.String()
			if len(temp) > 0 && temp[len(temp)-1] != ' ' {
				result.WriteRune(' ')
			}
		}
		result.WriteRune(char)
	}
	return strings.ToLower(result.String())
}

func ToLowerWordsWithHyphens(s string) string {
	var result strings.Builder
	for i, char := range s {
		if char == ' ' {
			continue
		}
		if i > 0 && char >= 'A' && char <= 'Z' {
			temp := result.String()
			if len(temp) > 0 && temp[len(temp)-1] != ' ' {
				result.WriteRune('-')
			}
		}
		result.WriteRune(char)
	}
	return strings.ToLower(result.String())
}

func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func ToCamelWord(s string, isByUnderline bool, isLowerCaseInitial bool) string {
	var sli []string
	if isByUnderline {
		sli = strings.Split(s, "_")
	} else {
		sli = strings.Split(s, " ")
	}
	var result strings.Builder
	for _, word := range sli {
		if result.String() == "" && isLowerCaseInitial {
			result.WriteString(word)
		} else {
			result.WriteString(FirstUpper(word))
		}
	}
	return result.String()
}

func CreateTestMainFile(testPath string, testFuncName string) {
	dirPath := path.Join(testPath, ToLowerWordsWithHyphens(testFuncName))
	err := os.Mkdir(dirPath, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	filePath := path.Join(dirPath, "main.go")
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(f)
	content := []byte("package main\n\nfunc main() {\n\n}\n")
	_, err = f.Write(content)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(filePath, "has been created successfully!")
}

func BuildTestMainFile(testPath string, testFuncName string) {
	if strings.HasPrefix(testFuncName, ".\\") && strings.HasSuffix(testFuncName, ".exe") {
		testFuncName, _ = strings.CutPrefix(testFuncName, ".\\")
		testFuncName, _ = strings.CutSuffix(testFuncName, ".exe")
	}
	dirPath := path.Join(testPath, ToLowerWordsWithHyphens(testFuncName))
	filePath := path.Join(dirPath, "main.go")
	executableFileName := testFuncName + ".exe"
	tags := "signrpc walletrpc chainrpc invoicesrpc autopilotrpc btlapi"
	cmd := exec.Command("go", "build", "-tags", tags, "-o", executableFileName, filePath)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(path.Join(testPath, executableFileName), "has been built successfully!")
}

func GetFunctionName(i any) string {
	completeName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	s := strings.Split(completeName, ".")
	return s[len(s)-1]
}

func GetTimestamp() int {
	return int(time.Now().Unix())
}

// TxHashConversion size-end conversion
func TxHashConversion(txHash string) string {
	b, err := hex.DecodeString(txHash)
	if err != nil {
		return err.Error()
	}
	for i := 0; i < len(b)/2; i++ {
		temp := b[i]
		b[i] = b[len(b)-i-1]
		b[len(b)-i-1] = temp
	}
	txHash = hex.EncodeToString(b)
	return txHash
}

func FixAsset(output string) string {
	//str, err := rpcclient.FixAsset(output, false)
	//if err != nil {
	//	fmt.Println(err)
	//	return MakeJsonErrorResult(DefaultErr, "FixAsset error, please check the output parameter and whether the asset needs to be repaired", nil)
	//}
	//fmt.Println(str)
	return MakeJsonErrorResult(SUCCESS, "", "str")
}

// GetRandomNumber
// @Description: Return a random number whose range is (0,maxValue).
func GetRandomNumber(maxValue int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	var randNumber int
	for randNumber == 0 {
		randNumber = rand.Intn(maxValue)
	}
	return randNumber
}

// GetRandomNumberSlice
// @Description: Use this to generate a slice of random number
func GetRandomNumberSlice(maxValue int, n int) ([]int, error) {
	// @dev: If the number of times the same random number is obtained exceeds this value,
	// add this random number to the slice of random number results,
	// so that duplicates can be obtained relatively randomly
	// when the range of random numbers is small
	// (the number of random numbers that can be obtained is less than
	// the number of random numbers that need to be generated).
	// The third time when a random number duplicates, it will be accepted.
	const retryTimes = 2
	if maxValue < 1 || n < 0 {
		return nil, errors.New("max value or slice length is negative")
	}
	var slice []int
	randNumberMapGeneratedTimes := make(map[int]int)
	for len(slice) < n {
		// @dev: Return a random number whose range is (0,maxValue].
		randNumber := GetRandomNumber(maxValue + 1)
		randNumberMapGeneratedTimes[randNumber]++
		if randNumberMapGeneratedTimes[randNumber] == 1 || randNumberMapGeneratedTimes[randNumber] > 1+retryTimes {
			slice = append(slice, randNumber)
			randNumberMapGeneratedTimes[randNumber] = 1
		} else {
			// @dev: length of slice will not add
			continue
		}
	}
	return slice, nil
}

func GetNowTimeStringWithHyphens() string {
	now := time.Now().Format("2006-01-02-15-04-05.000000")
	now = strings.ReplaceAll(now, ".", "-")
	return now
}

// todo:To be optimized
// 1.Memory risk
// 2.Use db

func DbSetTempStoreMemory() {
	db := tapdb.GetTestDB()
	// 设置PRAGMA temp_store=2，将临时文件存储在内存中
	_, err := db.Exec("PRAGMA temp_store = 2;")
	if err != nil {
		fmt.Println("设置PRAGMA失败:", err)
		return
	}
	fmt.Println("成功设置PRAGMA temp_store=2")
}

func DbSetTempStoreDefault() {
	db := tapdb.GetTestDB()
	// 设置PRAGMA temp_store=2，将临时文件存储在内存中
	_, err := db.Exec("PRAGMA temp_store = 0;")
	if err != nil {
		fmt.Println("设置PRAGMA失败:", err)
		return
	}
	fmt.Println("成功设置PRAGMA temp_store=0")
}
