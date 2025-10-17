package api

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/joho/godotenv"
	"github.com/lightninglabs/taproot-assets/tapdb"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/protobuf/proto"
)

type ErrCode int

const (
	DefaultErr   ErrCode = -1
	SUCCESS      ErrCode = 200
	NotFoundData ErrCode = iota + 299
	RequestError
	SUCCESS_2 ErrCode = 0
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
	UploadLogFileAndGetResponseErr
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
	UploadBigFileAndGetResponseErr
	DuplicateAddrErr
	GetListBalancesSimpleInfoHashAndUpdateAssetBalanceBackupErr
	CheckIfBackupIsRequiredErr
	GetAndUploadAssetBalanceHistoriesErr
	uploadBtcListUnspentUtxosErr
	FundChannelErr
	AddInvoiceErr
	SendPaymentErr
	ConnectPeerErr
	GetIdentityPubkeyErr
	CloseChannelErr
	DecodeAssetPayReqErr
	QueryListAssetsByAssetIdErr
	getSubServerStatusInfoErr
	lndSyncToChainErr
	InvalidParamsErr
	ListPaymentsErr
	OpenBtcChannelErr
	SendPaymentV2Err
	DeriveKeysErr
	CreateVirtualPSBTErr
	FundVirtualPSBTErr
	PrepareOutputAssetsErr
	SignVirtualPSBTErr
	PrepareBitcoinPSBTErr
	GeneratePaymentScriptErr
	GetAddressBip32DerivationErr
	SerializeErr
	CommitVirtualPsbtsErr
	NewFromRawBytesErr
	SignBitcoinPSBTErr
	SignPsbtErr
	EncodePSBTtoBase64Err
	FinalizePacketErr
	QueryAssetRatesErr
	getListEligibleCoinsErr
	psbtTrustlessSwapCreateSellOrderSignWithOneFilterErr
	FromStrErr
	psbtTrustlessSwapBuySOrderSignErr
	bytesToPsbtPacketErr
	tappsbtDecodeErr
	psbtTrustlessSwapPublishTxErr
	base64StdEncodingDecodeStringErr
	psbtTrustlessSwapBuySOrderProofErr
	LastProofFromStrErr
	fixVersionDirtyErr
	requestToGetLastProofErr
	GetAssetsDecimalErr
	InvalidDecimalDisplay
	NullInvoice
	lnurlPayErr
	PostServerToRequestInvoiceErr
	GetServerRequestAvailablePortErr
	PostServerRequestIsPortListeningErr
	ServerPortNotAvailable
	PostServerToRequestLnurlErr
	ServerRequestPortIsListening
	FrpcConfErr
	syncUniverseAssetErr
	pushProofErr
)

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
	case errors.Is(ec, SUCCESS_2):
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
	Success bool    `json:"success"`
	Error   string  `json:"error"`
	Code    ErrCode `json:"code"`
	Data    any     `json:"data"`
}

type Result2 struct {
	Errnos int         `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

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

func MakeJsonErrorResult2(code ErrCode, errorString string, data any) string {
	if int(code) == int(SUCCESS) {
		code = SUCCESS_2
	}
	jsr := Result2{
		Errnos: int(code),
		ErrMsg: errorString,
		Data:   data,
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

func FeeRateBtcPerKbToSatPerKw(btcPerKb float64) (satPerKw int) {
	return int(0.25e8 * btcPerKb)
}

func FeeRateBtcPerKbToSatPerB(btcPerKb float64) (satPerB int) {
	return int(1e5 * btcPerKb)
}

func FeeRateSatPerKwToBtcPerKb(feeRateSatPerKw int) (feeRateBtcPerKb float64) {
	return RoundToDecimalPlace(float64(feeRateSatPerKw)/0.25e8, 8)
}

func FeeRateSatPerKwToSatPerB(feeRateSatPerKw int) (feeRateSatPerB int) {
	return int(math.Ceil(float64(feeRateSatPerKw) * 4 / 1000))
}

func FeeRateSatPerBToBtcPerKb(feeRateSatPerB int) (feeRateBtcPerKb float64) {
	return RoundToDecimalPlace(math.Ceil(float64(feeRateSatPerB)/100000), 8)
}

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
	content := []byte("package main\n\n/*\n-tags \"litd autopilotrpc signrpc walletrpc chainrpc invoicesrpc watchtowerrpc neutrinorpc peersrpc btlapi\"\n*/\n\nfunc main() {\n\n}\n")
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
	tags := "litd autopilotrpc signrpc walletrpc chainrpc invoicesrpc watchtowerrpc neutrinorpc peersrpc btlapi"
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

func TxHashEncodeToString(h []byte) string {
	slices.Reverse(h)
	return hex.EncodeToString(h)
}

func TxHashStringReverse(h string) (d string) {
	for i := len(h) - 2; i > -1; i -= 2 {
		d = d + h[i:i+2]
	}
	return d
}

func FixAsset(output string) string {
	return MakeJsonErrorResult(SUCCESS, "", "str")
}

func GetRandomNumber(maxValue int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	var randNumber int
	for randNumber == 0 {
		randNumber = rand.Intn(maxValue)
	}
	return randNumber
}

func GetRandomNumberSlice(maxValue int, n int) ([]int, error) {
	const retryTimes = 2
	if maxValue < 1 || n < 0 {
		return nil, errors.New("max value or slice length is negative")
	}
	var slice []int
	randNumberMapGeneratedTimes := make(map[int]int)
	for len(slice) < n {
		randNumber := GetRandomNumber(maxValue + 1)
		randNumberMapGeneratedTimes[randNumber]++
		if randNumberMapGeneratedTimes[randNumber] == 1 || randNumberMapGeneratedTimes[randNumber] > 1+retryTimes {
			slice = append(slice, randNumber)
			randNumberMapGeneratedTimes[randNumber] = 1
		} else {
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

func Sha256(data any) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(jsonData)
	hashString := fmt.Sprintf("%x", hash)
	return hashString, nil
}

func EncodeDataToBase64(data any) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(jsonData), nil
}

func DecodeBase64ToData(encoded string, v any) error {
	byteData, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return err
	}
	return json.Unmarshal(byteData, v)
}

func DbSetTempStoreMemory() {
	db := tapdb.GetTestDB()
	_, err := db.Exec("PRAGMA temp_store = 2;")
	if err != nil {
		fmt.Println("设置PRAGMA失败:", err)
		return
	}
	fmt.Println("成功设置PRAGMA temp_store=2")
}

func DbSetTempStoreDefault() {
	db := tapdb.GetTestDB()
	_, err := db.Exec("PRAGMA temp_store = 0;")
	if err != nil {
		fmt.Println("设置PRAGMA失败:", err)
		return
	}
	fmt.Println("成功设置PRAGMA temp_store=0")
}

func AppendFileLog(filePath string, prefix string, content string) error {
	dir := filepath.Dir(filePath)
	if dir != "." {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				return errors.Wrap(err, "os.MkdirAll")
			}
		}
	}

	fileInfo, err := os.Stat(filePath)
	if err == nil && fileInfo.Size() > 3*1024*1024 {
		err := os.Truncate(filePath, 0)
		if err != nil {
			return errors.Wrap(err, "os.Truncate")
		}
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		return errors.Wrap(err, "os.OpenFile")
	}
	defer file.Close()
	logger := log.New(file, prefix, log.Lshortfile|log.Ldate|log.Ltime)
	logger.Println(content)
	return nil
}

func WriteToFile(path, content string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return errors.Wrap(err, "os.MkdirAll")
		}
	}
	filename := filepath.Base(path)
	filePath := dir + "/" + filename
	err := os.WriteFile(filePath, []byte(content), 0777)
	if err != nil {
		return errors.Wrap(err, "os.WriteFile")
	}
	return nil
}

func RespJsonStr(resp proto.Message) (string, error) {
	jsonBytes, err := lnrpc.ProtoJSONMarshalOpts.Marshal(resp)
	if err != nil {
		return "", errors.Wrap(err, "lnrpc.ProtoJSONMarshalOpts.Marshal")
	}
	return string(jsonBytes), nil
}
