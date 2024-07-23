package api

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/protobuf/proto"
	"log"
	"math"
	"os"
	"os/exec"
	"path"
	"reflect"
	"runtime"
	"strings"
	"time"
)

type ErrCode int

// Errtype:Normal
const (
	DefaultErr   ErrCode = -1
	SUCCESS      ErrCode = 200
	NotFoundData ErrCode = iota + 299
	RequestError
)

// Errtype:Unkonwn
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
	streamRecvIoEofErr
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
)

var ErrMsgMap = map[ErrCode]error{
	NotFoundData: errors.New("not found Data"),
	SUCCESS:      errors.New(""),
	RequestError: errors.New("request error"),
}

func GetErrMsg(code ErrCode) string {
	return ErrMsgMap[code].Error()
}

var (
	SuccessErr   = errors.New("")
	SuccessError = SuccessErr.Error()
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
	if errorString == "" {
		errorString = GetErrMsg(code)
	}
	jsr := JsonResult{
		Error: errorString,
		Code:  code,
		Data:  data,
	}
	if code == SUCCESS {
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
// @param feeRateSatPerKw
// @return feeRateSatPerB
func FeeRateSatPerKwToSatPerB(feeRateSatPerKw int) (feeRateSatPerB int) {
	return feeRateSatPerKw * 4 / 1000
}

// FeeRateSatPerBToBtcPerKb
// @Description: sat/b to BTC/Kb
// @param feeRateSatPerB
// @return feeRateBtcPerKb
func FeeRateSatPerBToBtcPerKb(feeRateSatPerB int) (feeRateBtcPerKb float64) {
	return RoundToDecimalPlace(float64(feeRateSatPerB)/100000, 8)
}

// FeeRateSatPerBToSatPerKw
// @Description: sat/b to sat/kw
// @param feeRateSatPerB
// @return feeRateSatPerKw
func FeeRateSatPerBToSatPerKw(feeRateSatPerB int) (feeRateSatPerKw int) {
	return feeRateSatPerB * 1000 / 4
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
	dirPath := path.Join(testPath, ToLowerWordsWithHyphens(testFuncName))
	filePath := path.Join(dirPath, "main.go")
	executableFileName := testFuncName + ".exe"
	cmd := exec.Command("go", "build", "-o", executableFileName, filePath)
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
