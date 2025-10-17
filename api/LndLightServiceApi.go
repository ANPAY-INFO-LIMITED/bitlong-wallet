package api

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/walletrpc"
	"github.com/pkg/errors"
	"github.com/wallet/base"
	"github.com/wallet/service/apiConnect"
	"github.com/wallet/service/rpcclient"
)

func getWalletBalance() (*lnrpc.WalletBalanceResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.WalletBalanceRequest{}
	response, err := client.WalletBalance(context.Background(), request)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("%s lnrpc WalletBalance response: %v\n", GetTimeNow(), response.String())
	return response, nil

}

func getInfoOfLnd() (*lnrpc.GetInfoResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.GetInfoRequest{}
	response, err := client.GetInfo(context.Background(), request)
	return response, err
}

type GetInfoFeature struct {
	Key        uint32 `json:"key"`
	Name       string `json:"name"`
	IsRequired bool   `json:"is_required"`
	IsKnown    bool   `json:"is_known"`
}

type GetInfoChain struct {
	Chain   string `json:"chain"`
	Network string `json:"network"`
}

type GetInfoResponse struct {
	Version                   string                       `json:"version"`
	CommitHash                string                       `json:"commit_hash"`
	IdentityPubkey            string                       `json:"identity_pubkey"`
	Alias                     string                       `json:"alias"`
	Color                     string                       `json:"color"`
	NumPendingChannels        uint32                       `json:"num_pending_channels"`
	NumActiveChannels         uint32                       `json:"num_active_channels"`
	NumInactiveChannels       uint32                       `json:"num_inactive_channels"`
	NumPeers                  uint32                       `json:"num_peers"`
	BlockHeight               uint32                       `json:"block_height"`
	BlockHash                 string                       `json:"block_hash"`
	BestHeaderTimestamp       int64                        `json:"best_header_timestamp"`
	SyncedToChain             bool                         `json:"synced_to_chain"`
	SyncedToGraph             bool                         `json:"synced_to_graph"`
	Testnet                   bool                         `json:"testnet"`
	Chains                    []GetInfoChain               `json:"chains"`
	Uris                      []string                     `json:"uris"`
	Features                  []GetInfoFeature             `json:"features"`
	RequireHtlcInterceptor    bool                         `json:"require_htlc_interceptor"`
	StoreFinalHtlcResolutions bool                         `json:"store_final_htlc_resolutions"`
	BitcoindGetBlockchainInfo *PostGetBlockchainInfoResult `json:"bitcoind_get_blockchain_info"`
}

type PostGetBlockchainInfoResult struct {
	Chain                string  `json:"chain"`
	Blocks               int     `json:"blocks"`
	Headers              int     `json:"headers"`
	Bestblockhash        string  `json:"bestblockhash"`
	Difficulty           float64 `json:"difficulty"`
	Time                 int     `json:"time"`
	Mediantime           int     `json:"mediantime"`
	Verificationprogress float64 `json:"verificationprogress"`
	Initialblockdownload bool    `json:"initialblockdownload"`
	Chainwork            string  `json:"chainwork"`
	SizeOnDisk           int     `json:"size_on_disk"`
	Pruned               bool    `json:"pruned"`
	Warnings             string  `json:"warnings"`
}

type GetBlockchainInfoResponse struct {
	Success bool                         `json:"success"`
	Error   string                       `json:"error"`
	Code    ErrCode                      `json:"code"`
	Data    *PostGetBlockchainInfoResult `json:"data"`
}

func ProcessGetInfoResponse(getInfoResponse *lnrpc.GetInfoResponse, getBlockchainInfo *PostGetBlockchainInfoResult) *GetInfoResponse {
	if getInfoResponse == nil {
		return nil
	}
	var chains []GetInfoChain
	var features []GetInfoFeature
	if getInfoResponse.Features != nil {
		for k, f := range getInfoResponse.Features {
			features = append(features, GetInfoFeature{
				Key:        k,
				Name:       f.Name,
				IsRequired: f.IsRequired,
				IsKnown:    f.IsKnown,
			})
		}
	}
	if getInfoResponse.Chains != nil {
		for _, c := range getInfoResponse.Chains {
			chains = append(chains, GetInfoChain{
				Chain:   c.Chain,
				Network: c.Network,
			})
		}
	}
	return &GetInfoResponse{
		Version:                   getInfoResponse.Version,
		CommitHash:                getInfoResponse.CommitHash,
		IdentityPubkey:            getInfoResponse.IdentityPubkey,
		Alias:                     getInfoResponse.Alias,
		Color:                     getInfoResponse.Color,
		NumPendingChannels:        getInfoResponse.NumPendingChannels,
		NumActiveChannels:         getInfoResponse.NumActiveChannels,
		NumInactiveChannels:       getInfoResponse.NumInactiveChannels,
		NumPeers:                  getInfoResponse.NumPeers,
		BlockHeight:               getInfoResponse.BlockHeight,
		BlockHash:                 getInfoResponse.BlockHash,
		BestHeaderTimestamp:       getInfoResponse.BestHeaderTimestamp,
		SyncedToChain:             getInfoResponse.SyncedToChain,
		SyncedToGraph:             getInfoResponse.SyncedToGraph,
		Testnet:                   getInfoResponse.Testnet,
		Chains:                    chains,
		Uris:                      getInfoResponse.Uris,
		Features:                  features,
		RequireHtlcInterceptor:    getInfoResponse.RequireHtlcInterceptor,
		StoreFinalHtlcResolutions: getInfoResponse.StoreFinalHtlcResolutions,
		BitcoindGetBlockchainInfo: getBlockchainInfo,
	}
}

func RequestToGetBlockchainInfo(token string) (*PostGetBlockchainInfoResult, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	network := base.NetWork
	url := serverDomainOrSocket + "/bitcoind/" + network + "/blockchain/get_blockchain_info"
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
	var response GetBlockchainInfoResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Data, nil
}

func GetBlockchainInfoAndGetResponse(token string) (*PostGetBlockchainInfoResult, error) {
	return RequestToGetBlockchainInfo(token)
}

func LndGetInfoAndGetResponse(token string) (*GetInfoResponse, error) {
	response, err := getInfoOfLnd()
	if err != nil {
		return nil, AppendErrorInfo(err, "getInfoOfLnd")
	}
	getBlockchainInfoResult, err := GetBlockchainInfoAndGetResponse(token)
	if err != nil {
		// @dev: Do not return
		LogError("GetBlockchainInfoAndGetResponse", err)
		getBlockchainInfoResult = &PostGetBlockchainInfoResult{}
	}
	result := ProcessGetInfoResponse(response, getBlockchainInfoResult)
	return result, nil
}

func LndGetInfo(token string) string {
	response, err := LndGetInfoAndGetResponse(token)
	if err != nil {
		return MakeJsonErrorResult(LndGetInfoAndGetResponseErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SUCCESS.Error(), response)
}

func lndSyncToChain() (syncedToChain bool, err error) {
	response, err := getInfoOfLnd()
	if err != nil {
		return false, AppendErrorInfo(err, "getInfoOfLnd")
	}
	syncedToChain = response.SyncedToChain
	return syncedToChain, nil
}

func LndSyncToChain() string {
	syncedToChain, err := lndSyncToChain()
	if err != nil {
		return MakeJsonErrorResult(lndSyncToChainErr, err.Error(), false)
	}
	return MakeJsonErrorResult(SUCCESS, SUCCESS.Error(), syncedToChain)
}

func sendCoins(addr string, amount int64, feeRate uint64, all bool) (*lnrpc.SendCoinsResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.SendCoinsRequest{
		Addr: addr,
	}
	if feeRate > 0 {
		request.SatPerVbyte = feeRate
	}
	if all {
		request.SendAll = true
	} else {
		request.Amount = amount
	}
	response, err := client.SendCoins(context.Background(), request)
	return response, err
}

type WalletBalanceResponse struct {
	TotalBalance       int `json:"total_balance"`
	ConfirmedBalance   int `json:"confirmed_balance"`
	UnconfirmedBalance int `json:"unconfirmed_balance"`
	LockedBalance      int `json:"locked_balance"`
}

func GetWalletBalance() string {
	response, err := getWalletBalance()
	if err != nil {
		fmt.Printf("%s lnrpc WalletBalance err: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(getWalletBalanceErr, err.Error(), nil)
	}
	// @dev: mark imported tap addresses as locked
	response, err = ProcessGetWalletBalanceResult(response)
	if err != nil {
		return MakeJsonErrorResult(ProcessGetWalletBalanceResultErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", WalletBalanceResponse{
		TotalBalance:       int(response.TotalBalance),
		ConfirmedBalance:   int(response.ConfirmedBalance),
		UnconfirmedBalance: int(response.UnconfirmedBalance),
		LockedBalance:      int(response.LockedBalance),
	})
}

func ProcessGetWalletBalanceResult(walletBalanceResponse *lnrpc.WalletBalanceResponse) (*lnrpc.WalletBalanceResponse, error) {
	imported, ok := walletBalanceResponse.AccountBalance["imported"]
	if !ok {
		return walletBalanceResponse, nil
	}
	importedConfirmedBalance := imported.ConfirmedBalance
	if importedConfirmedBalance == 0 {
		return walletBalanceResponse, nil
	}
	walletBalanceResponse.ConfirmedBalance -= importedConfirmedBalance
	walletBalanceResponse.LockedBalance += importedConfirmedBalance
	return walletBalanceResponse, nil
}

func CalculateImportedTapAddressBalanceAmount(listAddressesResponse *walletrpc.ListAddressesResponse) (imported int64) {
	if listAddressesResponse == nil {
		return 0
	}
	for _, addresses := range (*listAddressesResponse).AccountWithAddresses {
		if addresses.Name == "imported" {
			for _, address := range addresses.Addresses {
				if address.Balance == 1000 {
					imported += address.Balance
				}
			}
		}
	}
	return imported
}

func GetInfoOfLnd() string {
	response, err := getInfoOfLnd()
	if err != nil {
		fmt.Printf("%s lnrpc GetInfo err: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(getInfoOfLndErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func GetIdentityPubkey() string {
	response, err := getInfoOfLnd()
	if err != nil {
		fmt.Printf("%s lnrpc GetInfo.IdentityPubkey err: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(GetIdentityPubkeyErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response.GetIdentityPubkey())
}

func GetNewAddress() string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return ""
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.NewAddressRequest{
		Type: lnrpc.AddressType_TAPROOT_PUBKEY,
	}
	response, err := client.NewAddress(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc NewAddress err: %v\n", GetTimeNow(), err)
		return ""
	}
	return response.Address
}

func AddInvoice(value int, memo string, private bool) string {
	AppendFileLog("/data/data/io.bitlong/files/bitlong_api_log.txt", "[sshBtcChannelAddInvoice01]", ValueJsonString(fmt.Sprintf("memo:%s,value:%d", memo, value)))
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return MakeJsonErrorResult(GetConnectionErr, err.Error(), nil)
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.Invoice{
		Value:   int64(value),
		Memo:    memo,
		Private: private,
	}
	response, err := client.AddInvoice(context.Background(), request)
	if err != nil {
		fmt.Printf("%s client.AddInvoice :%v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(AddInvoiceErr, err.Error(), nil)
	}
	AppendFileLog("/data/data/io.bitlong/files/bitlong_api_log.txt", "[sshBtcChannelAddInvoice01]", ValueJsonString(fmt.Sprintf("response:%s", response)))
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func ListInvoices() string {
	response, err := rpcclient.ListInvoices(0, false)
	if err != nil {
		fmt.Printf("%s client.ListInvoice :%v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(ListInvoicesErr, err.Error(), nil)
	}
	invoices := SimplifyInvoice(response)
	return MakeJsonErrorResult(SUCCESS, "", invoices)
}

type InvoiceSimplified struct {
	PaymentRequest string `json:"payment_request"`
	Value          int    `json:"value"`
	State          string `json:"state"`
	CreationDate   int    `json:"creation_date"`
}

func SimplifyInvoice(invoice *lnrpc.ListInvoiceResponse) *[]InvoiceSimplified {
	var invoices []InvoiceSimplified
	for _, invoice := range invoice.Invoices {
		invoices = append(invoices, InvoiceSimplified{
			PaymentRequest: invoice.PaymentRequest,
			Value:          int(invoice.Value),
			State:          invoice.State.String(),
			CreationDate:   int(invoice.CreationDate),
		})
	}
	return &invoices
}

type InvoiceAll struct {
	Invoices []struct {
		Memo            string        `json:"memo"`
		RPreimage       string        `json:"r_preimage"`
		RHash           string        `json:"r_hash"`
		Value           string        `json:"value"`
		ValueMsat       string        `json:"value_msat"`
		Settled         bool          `json:"settled"`
		CreationDate    string        `json:"creation_date"`
		SettleDate      string        `json:"settle_date"`
		PaymentRequest  string        `json:"payment_request"`
		DescriptionHash string        `json:"description_hash"`
		Expiry          string        `json:"expiry"`
		FallbackAddr    string        `json:"fallback_addr"`
		CltvExpiry      string        `json:"cltv_expiry"`
		RouteHints      []interface{} `json:"route_hints"`
		Private         bool          `json:"private"`
		AddIndex        string        `json:"add_index"`
		SettleIndex     string        `json:"settle_index"`
		AmtPaid         string        `json:"amt_paid"`
		AmtPaidSat      string        `json:"amt_paid_sat"`
		AmtPaidMsat     string        `json:"amt_paid_msat"`
		State           string        `json:"state"`
		Htlcs           []interface{} `json:"htlcs"`
		Features        struct {
			Num9 struct {
				Name       string `json:"name"`
				IsRequired bool   `json:"is_required"`
				IsKnown    bool   `json:"is_known"`
			} `json:"9"`
			Num14 struct {
				Name       string `json:"name"`
				IsRequired bool   `json:"is_required"`
				IsKnown    bool   `json:"is_known"`
			} `json:"14"`
			Num17 struct {
				Name       string `json:"name"`
				IsRequired bool   `json:"is_required"`
				IsKnown    bool   `json:"is_known"`
			} `json:"17"`
		} `json:"features"`
		IsKeysend       bool   `json:"is_keysend"`
		PaymentAddr     string `json:"payment_addr"`
		IsAmp           bool   `json:"is_amp"`
		AmpInvoiceState struct {
		} `json:"amp_invoice_state"`
	} `json:"invoices"`
	LastIndexOffset  string `json:"last_index_offset"`
	FirstIndexOffset string `json:"first_index_offset"`
}

func LookupInvoice(rhash string) string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return ""
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	b_rhash, _ := hex.DecodeString(rhash)
	request := &lnrpc.PaymentHash{
		RHash: b_rhash,
	}
	response, err := client.LookupInvoice(context.Background(), request)
	if err != nil {
		fmt.Printf("%s client.LookupInvoice :%v\n", GetTimeNow(), err)
		return ""
	}
	return response.String()
}

func AbandonChannel() bool {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return false
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.AbandonChannelRequest{}
	response, err := client.AbandonChannel(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc AbandonChannel err: %v\n", GetTimeNow(), err)
		return false
	}
	fmt.Printf("%s %v\n", GetTimeNow(), response)
	return true
}

func BatchOpenChannel() bool {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return false
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.BatchOpenChannelRequest{}
	response, err := client.BatchOpenChannel(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc BatchOpenChannel err: %v\n", GetTimeNow(), err)
		return false
	}
	fmt.Printf("%s %v\n", GetTimeNow(), response)
	return true
}

func ChannelAcceptor() bool {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return false
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	stream, err := client.ChannelAcceptor(context.Background())
	if err != nil {
		fmt.Printf("%s lnrpc ChannelAcceptor err: %v\n", GetTimeNow(), err)
		return false
	}
	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Printf("%s err == io.EOF, err: %v\n", GetTimeNow(), err)
				return false
			}
			fmt.Printf("%s stream Recv err: %v\n", GetTimeNow(), err)
			return false
		}
		fmt.Printf("%s %v\n", GetTimeNow(), response)
		return true
	}
}

func ChannelBalance() string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return ""
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.ChannelBalanceRequest{}
	response, err := client.ChannelBalance(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc ChannelBalance err: %v\n", GetTimeNow(), err)
		return ""
	}
	fmt.Printf("%s %v\n", GetTimeNow(), response)
	return response.String()
}

func CheckMacaroonPermissions() bool {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return false
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.CheckMacPermRequest{}
	response, err := client.CheckMacaroonPermissions(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc CheckMacaroonPermissions err: %v\n", GetTimeNow(), err)
		return false
	}
	fmt.Printf("%s %v\n", GetTimeNow(), response)
	return true
}

func CloseChannel(channelPoint string) string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return MakeJsonErrorResult(GetConnectionErr, err.Error(), nil)
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)

	parts := strings.Split(channelPoint, ":")
	if len(parts) != 2 {
		return MakeJsonErrorResult(CloseChannelErr, "Invalid channelPoint format", nil)
	}

	// 提取 FundingTxidStr
	fundingTxidStr := parts[0]

	// 提取 OutputIndex
	outputIndex, err := strconv.Atoi(parts[1])
	if err != nil {
		return MakeJsonErrorResult(CloseChannelErr, "Invalid output index", nil)
	}

	request := &lnrpc.CloseChannelRequest{
		ChannelPoint: &lnrpc.ChannelPoint{
			FundingTxid: &lnrpc.ChannelPoint_FundingTxidStr{FundingTxidStr: fundingTxidStr},
			OutputIndex: uint32(outputIndex),
		},
		SatPerVbyte: 5,
	}
	stream, err := client.CloseChannel(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc CloseChannel err: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(CloseChannelErr, err.Error(), nil)
	}
	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Printf("%s err == io.EOF, err: %v\n", GetTimeNow(), err)
				return MakeJsonErrorResult(CloseChannelErr, err.Error(), nil)
			}
			fmt.Printf("%s stream Recv err: %v\n", GetTimeNow(), err)
			return MakeJsonErrorResult(CloseChannelErr, err.Error(), nil)
		} else if response != nil {
			fmt.Printf("%s %v\n", GetTimeNow(), response)
			return MakeJsonErrorResult(SUCCESS, "", "通道进入关闭状态")
		}
	}
}

func ClosedChannels() string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return ""
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.ClosedChannelsRequest{}
	response, err := client.ClosedChannels(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc ClosedChannels err: %v\n", GetTimeNow(), err)
		return err.Error()
	}
	return response.String()
}

func DecodePayReq(payReq string) string {
	req, err := rpcclient.DecodePayReq(payReq)
	if err != nil {
		return MakeJsonErrorResult(DecodePayReqErr, err.Error(), nil)
	}
	result := struct {
		Description string `json:"description"`
		Amount      int64  `json:"amount"`
		Timestamp   int64  `json:"timestamp"`
		Expiry      int64  `json:"expiry"`
		PaymentHash string `json:"payment_hash"`
		Destination string `json:"destination"`
	}{
		Description: req.Description,
		Amount:      req.NumSatoshis,
		Timestamp:   req.Timestamp,
		Expiry:      req.Expiry,
		PaymentHash: req.PaymentHash,
		Destination: req.Destination,
	}
	return MakeJsonErrorResult(SUCCESS, "", result)
}

func ExportAllChannelBackups() bool {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return false
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.ChanBackupExportRequest{}
	response, err := client.ExportAllChannelBackups(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc ChanBackupExportRequest err: %v\n", GetTimeNow(), err)
		return false
	}
	fmt.Printf("%s %v\n", GetTimeNow(), response)
	return true
}

func ExportChannelBackup() bool {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return false
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.ExportChannelBackupRequest{}
	response, err := client.ExportChannelBackup(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc ExportChannelBackup err: %v\n", GetTimeNow(), err)
		return false
	}
	fmt.Printf("%s %v\n", GetTimeNow(), response)
	return true
}

func GetChanInfo(chanId string) string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return ""
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	chainIdUint64, err := strconv.ParseUint(chanId, 10, 64)
	if err != nil {
		fmt.Printf("%s string to uint64 err: %v\n", GetTimeNow(), err)
	}
	request := &lnrpc.ChanInfoRequest{
		ChanId: chainIdUint64,
	}
	response, err := client.GetChanInfo(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc GetChanInfo err: %v\n", GetTimeNow(), err)
		return ""
	}
	return response.String()
}

func OpenChannelSync(nodePubkey string, localFundingAmount int64) string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return ""
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	_nodePubkeyByteSlice, _ := hex.DecodeString(nodePubkey)
	request := &lnrpc.OpenChannelRequest{
		NodePubkey:         _nodePubkeyByteSlice,
		LocalFundingAmount: localFundingAmount,
	}
	response, err := client.OpenChannelSync(context.Background(), request)
	if err != nil {
		return err.Error()
	}
	//deal with  the byte-reversed hash
	var txBytes []byte
	for i := len(response.GetFundingTxidBytes()) - 1; i >= 0; {
		txBytes = append(txBytes, response.GetFundingTxidBytes()[i])
		i--
	}

	txStr := hex.EncodeToString(txBytes)
	return txStr + ":" + strconv.Itoa(int(response.GetOutputIndex()))
}

func OpenBtcChannel(pubkey string, host string, localFundingAmount int, satPerVbyte int, pushSat int, memo string) string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return MakeJsonErrorResult(GetConnectionErr, err.Error(), nil)
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.ConnectPeerRequest{
		Addr: &lnrpc.LightningAddress{
			Pubkey: pubkey,
			Host:   host,
		},
	}
	response, err := client.ConnectPeer(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc ConnectPeer err: %v\n", GetTimeNow(), err)
		if strings.Contains(err.Error(), "already connected to peer") {
			fmt.Printf("%s Already connected to peer, skipping error.\n", GetTimeNow())
		} else {
			return MakeJsonErrorResult(ConnectPeerErr, err.Error(), nil)
		}
	}
	if pubkey == BtlPubKey && host == BtlHost {
		pushSat = 5000
	} else if pubkey == Btl2PubKey && host == Btl2Host {
		pushSat = 5000
	}

	AppendFileLog("/data/data/io.bitlong/files/bitlong_api_log.txt", "[sshBtcChannel-pushSat]", ValueJsonString(pushSat))
	fmt.Printf("%s %v\n", GetTimeNow(), response)

	if satPerVbyte > 10 {
		err := errors.New("fee rate exceeds max(10)")
		return MakeJsonErrorResult(FeeRateExceedMaxErr, err.Error(), nil)
	}
	return openBtcChannel(pubkey, localFundingAmount, satPerVbyte, pushSat, memo)
}

func openBtcChannel(nodePubkey string, localFundingAmount int, satPerVbyte int, pushSat int, memo string) string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return MakeJsonErrorResult(GetConnectionErr, err.Error(), nil)
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	_nodePubkeyByteSlice, _ := hex.DecodeString(nodePubkey)
	request := &lnrpc.OpenChannelRequest{
		SatPerVbyte:        uint64(satPerVbyte),
		NodePubkey:         _nodePubkeyByteSlice,
		LocalFundingAmount: int64(localFundingAmount),
		PushSat:            int64(pushSat),
		Memo:               memo,
	}
	stream, err := client.OpenChannel(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc OpenChannel err: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(OpenBtcChannelErr, err.Error(), nil)
	}

	conversion := func(b []byte) string {
		for i := 0; i < len(b)/2; i++ {
			temp := b[i]
			b[i] = b[len(b)-i-1]
			b[len(b)-i-1] = temp
		}
		txHash := hex.EncodeToString(b)
		return txHash
	}
	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Printf("%s err == io.EOF, err: %v\n", GetTimeNow(), err)
				return MakeJsonErrorResult(OpenBtcChannelErr, err.Error(), nil)
			}
			fmt.Printf("%s stream Recv err: %v\n", GetTimeNow(), err)
			return MakeJsonErrorResult(OpenBtcChannelErr, err.Error(), nil)
		} else if response.PendingChanId != nil {
			err := AppendFileLog("/data/data/io.bitlong/files/bitlong_log.txt", "[sshBtcChannels]", ValueJsonString(response))
			if err != nil {
				fmt.Printf("%s AppendFileLog err: %v\n", GetTimeNow(), err)
			}
			txid := response.GetChanPending().Txid
			hash := conversion(txid)
			OutputIndex := response.GetChanPending().OutputIndex
			resp := fmt.Sprintf("%s:%d", hash, OutputIndex)
			return MakeJsonErrorResult(SUCCESS, "", resp)
		}
	}
}

type AllBtcChannelBalance struct {
	AllSatsBalance          int64 `json:"all_sats_balance"`
	RemoteAllSatsBalance    int64 `json:"remote_all_sats_balance"`
	TotalCombinedSats       int64 `json:"total_combined_sats"`
	RemoteTotalCombinedSats int64 `json:"remote_total_combined_sats"`
}

func SumAllBtcChannelBalances() string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return MakeJsonErrorResult(GetConnectionErr, err.Error(), nil)
	}
	defer clearUp()

	var balance AllBtcChannelBalance

	client := lnrpc.NewLightningClient(conn)
	response, err := client.ListChannels(context.Background(), &lnrpc.ListChannelsRequest{})
	if err != nil {
		fmt.Printf("%s lnrpc ListChannels err: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(ListChannelsErr, err.Error(), nil)
	}
	for _, channel := range response.Channels {
		if len(channel.CustomChannelData) == 0 {
			balance.AllSatsBalance += channel.LocalBalance
			balance.RemoteAllSatsBalance += channel.RemoteBalance
			balance.TotalCombinedSats += channel.LocalBalance
			balance.RemoteTotalCombinedSats += channel.RemoteBalance
		} else {
			balance.TotalCombinedSats += channel.LocalBalance
			balance.RemoteTotalCombinedSats += channel.RemoteBalance
		}
	}

	return MakeJsonErrorResult(SUCCESS, "", balance)
}

func ListChannels() string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return MakeJsonErrorResult(GetConnectionErr, err.Error(), nil)
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.ListChannelsRequest{}
	response, err := client.ListChannels(context.Background(), request)
	if err != nil {
		return MakeJsonErrorResult(ListChannelsErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func PendingChannels() string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return MakeJsonErrorResult(GetConnectionErr, err.Error(), nil)
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.PendingChannelsRequest{}
	response, err := client.PendingChannels(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc PendingChannels err: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(PendingChannelsErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func GetChannelState(chanPoint string) string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return MakeJsonErrorResult(GetConnectionErr, err.Error(), nil)
	}

	defer clearUp()

	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.ListChannelsRequest{}
	response, err := client.ListChannels(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc ListChannels err: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(ListChannelsErr, err.Error(), nil)
	}

	var ChannelState string
	for _, channel := range response.Channels {
		if channel.ChannelPoint == chanPoint {
			if channel.Active {
				ChannelState = "ACTIVE"
			} else {
				ChannelState = "INACTIVE"
			}
			return MakeJsonErrorResult(SUCCESS, "", ChannelState)
		}
	}
	pendrequest := &lnrpc.PendingChannelsRequest{}
	pendingresponse, err := client.PendingChannels(context.Background(), pendrequest)
	if err != nil {
		fmt.Printf("%s lnrpc PendingChannels err: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(PendingChannelsErr, err.Error(), nil)
	}
	for _, channel := range pendingresponse.PendingOpenChannels {
		if channel.Channel.ChannelPoint == chanPoint {

			ChannelState = "PENDING_OPEN"
			return MakeJsonErrorResult(SUCCESS, "", ChannelState)
		}
	}
	for _, channel := range pendingresponse.WaitingCloseChannels {
		if channel.Channel.ChannelPoint == chanPoint {
			ChannelState = "PENDING_CLOSE"
			return MakeJsonErrorResult(SUCCESS, "", ChannelState)
		}
	}

	closerequest := &lnrpc.ClosedChannelsRequest{}
	closeresponse, err := client.ClosedChannels(context.Background(), closerequest)
	if err != nil {
		fmt.Printf("%s lnrpc ClosedChannels err: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(ClosedChannelsErr, err.Error(), nil)
	}
	for _, channel := range closeresponse.Channels {
		if channel.ChannelPoint == chanPoint {
			ChannelState = "CLOSED"
			return MakeJsonErrorResult(SUCCESS, "", ChannelState)
		}
	}

	return MakeJsonErrorResult(NoFindChannelErr, "NO_FIND_CHANNEL", nil)
}

func GetChannelInfo(chanPoint string) string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return MakeJsonErrorResult(GetConnectionErr, err.Error(), nil)
	}

	defer clearUp()

	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.ListChannelsRequest{}
	response, err := client.ListChannels(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc ListChannels err: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(ListChannelsErr, err.Error(), nil)
	}
	for _, channel := range response.Channels {
		if channel.ChannelPoint == chanPoint {
			return MakeJsonErrorResult(SUCCESS, "", channel)
		}
	}
	return MakeJsonErrorResult(NoFindChannelErr, "NO_FIND_CHANNEL", nil)
}

func RestoreChannelBackups() bool {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return false
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.RestoreChanBackupRequest{}
	response, err := client.RestoreChannelBackups(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc RestoreChannelBackups err: %v\n", GetTimeNow(), err)
		return false
	}
	fmt.Printf("%s %v\n", GetTimeNow(), response)
	return true
}

func SubscribeChannelBackups() bool {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return false
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.ChannelBackupSubscription{}
	stream, err := client.SubscribeChannelBackups(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc SubscribeChannelBackups err: %v\n", GetTimeNow(), err)
		return false
	}
	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Printf("%s err == io.EOF, err: %v\n", GetTimeNow(), err)
				return false
			}
			fmt.Printf("%s stream Recv err: %v\n", GetTimeNow(), err)
			return false
		}
		fmt.Printf("%s %v\n", GetTimeNow(), response)
		return true
	}

}

func SubscribeChannelEvents() bool {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return false
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.ChannelEventSubscription{}
	stream, err := client.SubscribeChannelEvents(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc  err: %v\n", GetTimeNow(), err)
		return false
	}
	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Printf("%s err == io.EOF, err: %v\n", GetTimeNow(), err)
				return false
			}
			fmt.Printf("%s stream Recv err: %v\n", GetTimeNow(), err)
			return false
		}
		fmt.Printf("%s %v\n", GetTimeNow(), response)
		return true
	}

}

func SubscribeChannelGraph() bool {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return false
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.GraphTopologySubscription{}
	stream, err := client.SubscribeChannelGraph(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc SubscribeChannelGraph err: %v\n", GetTimeNow(), err)
		return false
	}
	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Printf("%s err == io.EOF, err: %v\n", GetTimeNow(), err)
				return false
			}
			fmt.Printf("%s stream Recv err: %v\n", GetTimeNow(), err)
			return false
		}
		fmt.Printf("%s %v\n", GetTimeNow(), response)
		return true
	}

}

func UpdateChannelPolicy() bool {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return false
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.PolicyUpdateRequest{}
	response, err := client.UpdateChannelPolicy(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc UpdateChannelPolicy err: %v\n", GetTimeNow(), err)
		return false
	}
	fmt.Printf("%s %v\n", GetTimeNow(), response)
	return true
}

func VerifyChanBackup() bool {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return false
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.ChanBackupSnapshot{}
	response, err := client.VerifyChanBackup(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc VerifyChanBackup err: %v\n", GetTimeNow(), err)
		return false
	}
	fmt.Printf("%s %v\n", GetTimeNow(), response)
	return true
}

func ConnectPeer(pubkey, host string) string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		AppendFileLog("/data/data/io.bitlong/files/bitlong_api_log.txt", "[connectPeer]", ValueJsonString(err))
		return MakeJsonErrorResult(ConnectPeerErr, err.Error(), nil)
	}
	defer clearUp()

	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.ConnectPeerRequest{
		Addr: &lnrpc.LightningAddress{
			Pubkey: pubkey,
			Host:   host,
		},
	}
	response, err := client.ConnectPeer(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc ConnectPeer err: %v\n", GetTimeNow(), err)
		if strings.Contains(err.Error(), "already connected to peer") {
			return MakeJsonErrorResult(SUCCESS, "", "already connected to peer")
		} else {
			AppendFileLog("/data/data/io.bitlong/files/bitlong_api_log.txt", "[connectPeer02]", ValueJsonString(err))
			return MakeJsonErrorResult(ConnectPeerErr, err.Error(), nil)
		}
	}
	fmt.Printf("%s %v\n", GetTimeNow(), response)
	return MakeJsonErrorResult(SUCCESS, "", "connect peer success")
}

func EstimateFee(addr string, amount int64) string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return ""
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	addrToAmount := make(map[string]int64)
	addrToAmount[addr] = amount
	request := &lnrpc.EstimateFeeRequest{
		AddrToAmount: addrToAmount,
	}
	response, err := client.EstimateFee(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc ConnectPeer err: %v\n", GetTimeNow(), err)
		return ""
	}
	return response.String()
}

func SendPaymentSync(invoice string) string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return MakeJsonErrorResult(GetConnectionErr, err.Error(), nil)
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.SendRequest{
		PaymentRequest: invoice,
	}
	response, err := client.SendPaymentSync(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc SendPaymentSync :%v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(SendPaymentSyncErr, err.Error(), nil)
	}
	paymentHash := hex.EncodeToString(response.PaymentHash)
	return MakeJsonErrorResult(SUCCESS, "", paymentHash)
}

func SendPaymentSync0amt(invoice string, amt int64) string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return ""
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.SendRequest{
		PaymentRequest: invoice,
		Amt:            amt,
	}
	stream, err := client.SendPaymentSync(context.Background(), request)
	if err != nil {
		fmt.Printf("%s client.SendPaymentSync :%v\n", GetTimeNow(), err)
		return "false"
	}
	fmt.Printf("%s %s", GetTimeNow(), stream.String())
	return hex.EncodeToString(stream.PaymentHash)
}

func MergeUTXO(feeRate int64) string {
	if feeRate > 500 {
		err := errors.New("fee rate exceeds max(500)")
		return MakeJsonErrorResult(FeeRateExceedMaxErr, err.Error(), nil)
	}
	//创建一个地址
	addrTarget := GetNewAddress()
	//将所有utxo合并到这个地址
	response, err := sendCoins(addrTarget, 0, uint64(feeRate), true)
	if err != nil {
		return MakeJsonErrorResult(sendCoinsErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func SendCoins(addr string, amount int64, feeRate int64, sendAll bool) string {
	if feeRate > 500 {
		err := errors.New("fee rate exceeds max(500)")
		return MakeJsonErrorResult(FeeRateExceedMaxErr, err.Error(), nil)
	}
	if !verifyBtcAddress(addr) {
		return MakeJsonErrorResult(sendCoinsErr, fmt.Sprintf("invalid address:%s", addr), nil)
	}
	response, err := sendCoins(addr, amount, uint64(feeRate), sendAll)
	if err != nil {
		return MakeJsonErrorResult(sendCoinsErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func SendMany(jsonAddr string, feeRate int64) string {
	if feeRate > 500 {
		err := errors.New("fee rate exceeds max(500)")
		return MakeJsonErrorResult(FeeRateExceedMaxErr, err.Error(), nil)
	}
	var addrs []struct {
		Address string `json:"address"`
		Amount  int64  `json:"btcSum"`
	}
	err := json.Unmarshal([]byte(jsonAddr), &addrs)
	if err != nil {
		return MakeJsonErrorResult(UnmarshalErr, "Please use the correct json format", nil)
	}
	if len(addrs) == 0 {
		return MakeJsonErrorResult(AddrsLenZeroErr, "Please input the correct address and amount", nil)
	}
	addrTo := make(map[string]int64)
	for _, addr := range addrs {
		if addr.Amount <= 0 {
			return MakeJsonErrorResult(sendManyErr, "cannot send 0 Satoshi", nil)
		}
		if !verifyBtcAddress(addr.Address) {
			if addr.Address == "" {
				return MakeJsonErrorResult(sendManyErr, "got a empty address", nil)
			}
			return MakeJsonErrorResult(sendManyErr, fmt.Sprintf("invalid address:%s", addr.Address), nil)
		}
		addrTo[addr.Address] = addr.Amount
	}
	response, err := sendMany(addrTo, uint64(feeRate))
	if err != nil {
		return MakeJsonErrorResult(sendManyErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}
func verifyBtcAddress(address string) bool {
	var params *chaincfg.Params
	// 解析并验证地址
	switch base.NetWork {
	case base.UseTestNet:
		params = &chaincfg.TestNet3Params
	case base.UseMainNet:
		params = &chaincfg.MainNetParams
	case base.UseRegTest:
		params = &chaincfg.RegressionNetParams
	default:
		log.Println("NetWork need set testnet, mainnet or regtest")
		return false
	}
	_, err := btcutil.DecodeAddress(address, params)
	if err != nil {
		return false
	}
	return true
}

func sendMany(addr map[string]int64, feerate uint64) (*lnrpc.SendManyResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.SendManyRequest{
		AddrToAmount: addr,
	}
	if feerate > 0 {
		request.SatPerVbyte = feerate
	}
	response, err := client.SendMany(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func SendAllCoins(addr string, feeRate int) string {
	if feeRate > 500 {
		err := errors.New("fee rate exceeds max(500)")
		return MakeJsonErrorResult(FeeRateExceedMaxErr, err.Error(), nil)
	}
	response, err := sendCoins(addr, 0, uint64(feeRate), true)
	if err != nil {
		return MakeJsonErrorResult(sendCoinsErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func LndStopDaemon() bool {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return false
	}

	defer clearUp()

	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.StopRequest{}
	response, err := client.StopDaemon(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc StopDaemon err: %v\n", GetTimeNow(), err)
		return false
	}
	fmt.Printf("%s %v\n", GetTimeNow(), response)
	return true
}

func ListPermissions() string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return ""
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.ListPermissionsRequest{}
	response, err := client.ListPermissions(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc ListPermissions err: %v\n", GetTimeNow(), err)
		return err.Error()
	}
	fmt.Printf("%s %v\n", GetTimeNow(), response)
	return response.String()
}

type BtcChanInvoiceSimplified struct {
	PaymentRequest string `json:"payment_request"`
	Value          int64  `json:"value"`
	ChanId         uint64 `json:"chan_id"`
	State          string `json:"state"`
	CreationDate   int64  `json:"creation_date"`
}

func GetBtcChannelListInvoicesSettled(isAllChannel bool, chanId int, indexOffset int, maxInvoices int) string {
	if !isAllChannel && chanId <= 0 {
		return MakeJsonErrorResult(InvalidParamsErr, "channel ID is required when isAllChannel is false", nil)
	}
	if maxInvoices <= 0 || maxInvoices > 10000 {
		maxInvoices = 100
	}

	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return MakeJsonErrorResult(GetConnectionErr, err.Error(), nil)
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.ListInvoiceRequest{}
	response, err := client.ListInvoices(context.Background(), request)
	if err != nil {
		return MakeJsonErrorResult(ListInvoicesErr, err.Error(), nil)
	}

	invoices := SimplifyBtcChanInvoices(response)
	sort.Slice(invoices, func(i, j int) bool {
		return invoices[i].CreationDate > invoices[j].CreationDate
	})
	if isAllChannel {
		start, end := calculatePagination(indexOffset, maxInvoices, len(invoices))
		result := invoices[start:end]
		return MakeJsonErrorResult(SUCCESS, "", result)
	}
	var InvoiceSimplified []BtcChanInvoiceSimplified
	for _, invoice := range invoices {
		if invoice.ChanId == uint64(chanId) {
			InvoiceSimplified = append(InvoiceSimplified, invoice)
		}
	}
	start, end := calculatePagination(indexOffset, maxInvoices, len(InvoiceSimplified))
	result := InvoiceSimplified[start:end]
	return MakeJsonErrorResult(SUCCESS, "", result)
}

func SimplifyBtcChanInvoices(invoice *lnrpc.ListInvoiceResponse) []BtcChanInvoiceSimplified {
	var invoices []BtcChanInvoiceSimplified
	for _, invoice := range invoice.Invoices {
		if invoice.State == 1 && invoice.Htlcs[0].CustomChannelData == nil {
			invoices = append(invoices, BtcChanInvoiceSimplified{
				PaymentRequest: invoice.PaymentRequest,
				Value:          invoice.Value,
				ChanId:         invoice.Htlcs[0].ChanId,
				State:          invoice.State.String(),
				CreationDate:   invoice.CreationDate,
			})
		}
	}
	return invoices
}

type BtcChanPaymentSimplified struct {
	AssetId        string                      `json:"asset_id"`
	Amount         int64                       `json:"amount"`
	RfqId          string                      `json:"rfq_id"`
	ChanId         uint64                      `json:"chan_id"`
	ValueSat       int64                       `json:"value_sat"`
	CreationTimeNs int64                       `json:"creation_time_ns"`
	PaymentRequest string                      `json:"payment_request"`
	PaymentIndex   uint64                      `json:"payment_index"`
	Status         lnrpc.Payment_PaymentStatus `json:"status"`
	FailureReason  lnrpc.PaymentFailureReason  `json:"failure_reason"`
}

func GetBtcListPayments(isAllChannel bool, chanId int, indexOffset int, maxPayments int) string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return MakeJsonErrorResult(GetConnectionErr, err.Error(), nil)
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	response, err := client.ListPayments(context.Background(), &lnrpc.ListPaymentsRequest{})
	if err != nil {
		return MakeJsonErrorResult(ListPaymentsErr, err.Error(), nil)
	}
	payments := SimplifyBtcChanPayments(response)
	sort.Slice(payments, func(i, j int) bool {
		return payments[i].CreationTimeNs > payments[j].CreationTimeNs
	})
	if isAllChannel {
		start, end := calculatePagination(indexOffset, maxPayments, len(payments))
		result := payments[start:end]
		return MakeJsonErrorResult(SUCCESS, "", result)
	}

	var PaymentSimplified []BtcChanPaymentSimplified
	for _, payment := range payments {
		if payment.ChanId == uint64(chanId) {
			PaymentSimplified = append(PaymentSimplified, payment)
		}
	}

	start, end := calculatePagination(indexOffset, maxPayments, len(PaymentSimplified))
	result := PaymentSimplified[start:end]
	return MakeJsonErrorResult(SUCCESS, "", result)
}

func SimplifyBtcChanPayments(payment *lnrpc.ListPaymentsResponse) []BtcChanPaymentSimplified {
	var payments []BtcChanPaymentSimplified
	for _, payment := range payment.Payments {
		if payment.FailureReason == 0 && payment.Htlcs[0].Route.CustomChannelData == nil {
			payments = append(payments, BtcChanPaymentSimplified{
				AssetId:        "00",
				Amount:         0,
				RfqId:          "",
				ChanId:         payment.Htlcs[0].Route.Hops[0].ChanId,
				ValueSat:       payment.ValueSat,
				CreationTimeNs: payment.CreationTimeNs,
				PaymentRequest: payment.PaymentRequest,
				PaymentIndex:   payment.PaymentIndex,
				Status:         payment.Status,
				FailureReason:  payment.FailureReason,
			})
		}
	}
	return payments
}

func calculatePagination(offset, limit, total int) (int, int) {
	if offset >= total {
		return 0, 0
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return offset, end
}
