package api

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lightninglabs/taproot-assets/rfqmath"
	"github.com/lightninglabs/taproot-assets/rfqmsg"
	"github.com/lightninglabs/taproot-assets/taprpc/btlchannelrpc"
	"github.com/lightninglabs/taproot-assets/taprpc/priceoraclerpc"
	"github.com/lightninglabs/taproot-assets/taprpc/tapchannelrpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	"github.com/lightningnetwork/lnd/lntypes"
	"github.com/lightningnetwork/lnd/record"
	"github.com/wallet/service/apiConnect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	BtlPubKey  string = "027208d43d94fa830417a0e730d12cd11996d5cd62df2e210cdc48464feaafd3bc"
	BtlHost    string = "132.232.109.84:9736"
	Btl2PubKey string = "0298978cfb4dc3cf5f9ba387c456a52c99f9fa79d91d1e986dbaf34b5f3270c07b"
	Btl2Host   string = "118.24.37.253:9735"
)

func OpenAssetChannel(assetId, pubkey, host string, amount, feeRate, pushSat, localAmt int) string {
	AppendFileLog("/data/data/io.bitlong/files/bitlong_api_log.txt", "[sshTapChannels01]", ValueJsonString(fmt.Sprintf("assetId:%s,pubkey:%s,host:%s,amount:%d,feeRate:%d,pushSat:%d", assetId, pubkey, host, amount, feeRate, pushSat)))
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return MakeJsonErrorResult(GetConnectionErr, err.Error(), nil)
	}
	defer clearUp()

	if localAmt < 20000 {
		return MakeJsonErrorResult(ListChannelsErr, "本地sats不能低于2万", nil)
	}

	ctx := context.Background()
	ctxt, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	client := lnrpc.NewLightningClient(conn)

	request := &lnrpc.ConnectPeerRequest{
		Addr: &lnrpc.LightningAddress{
			Pubkey: pubkey,
			Host:   host,
		},
	}

	response, err := client.ConnectPeer(ctxt, request)
	if err != nil {
		fmt.Printf("%s lnrpc ConnectPeer err: %v\n", GetTimeNow(), err)
		if strings.Contains(err.Error(), "already connected to peer") {
			fmt.Printf("%s Already connected to peer, skipping error.\n", GetTimeNow())
		} else {
			return MakeJsonErrorResult(ConnectPeerErr, err.Error(), nil)
		}
	}
	AppendFileLog("/data/data/io.bitlong/files/bitlong_api_log.txt", "[sshTapChannels06]", ValueJsonString(response))
	if pubkey == BtlPubKey && host == BtlHost {
		pushSat = 5000
	} else if pubkey == Btl2PubKey && host == Btl2Host {
		pushSat = 5000
	}

	if feeRate > 10 {
		err := errors.New("fee rate exceeds max(10)")
		return MakeJsonErrorResult(FeeRateExceedMaxErr, err.Error(), nil)
	}
	return openAssetChannel(assetId, amount, pubkey, feeRate, pushSat, localAmt)
}

func openAssetChannel(assetId string, amount int, pubkey string, feeRate int, pushSat int, localAmt int) string {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return MakeJsonErrorResult(DecodeStringErr, err.Error(), nil)
	}

	defer clearUp()

	assetIdStr, err := hex.DecodeString(assetId)
	if err != nil {
		return MakeJsonErrorResult(DecodeStringErr, err.Error(), nil)
	}

	if len(assetIdStr) != sha256.Size {
		fmt.Errorf("asset id must be 32 bytes")
		return MakeJsonErrorResult(DecodeStringErr, "asset id must be 32 bytes", nil)
	}

	pubKeyStr, err := hex.DecodeString(pubkey)
	if err != nil {
		return MakeJsonErrorResult(DecodeStringErr, err.Error(), nil)
	}

	client := btlchannelrpc.NewBtlChannelsClient(conn)
	ctxt, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req := &btlchannelrpc.FundBtlChannelRequest{
		AssetAmount:        uint64(amount),
		AssetId:            assetIdStr,
		PeerPubkey:         pubKeyStr,
		FeeRateSatPerVbyte: uint32(feeRate),
		PushSat:            int64(pushSat),
		LocalAmt:           uint64(localAmt),
	}
	resp, err := client.FundBtlChannel(ctxt, req)
	if err != nil {
		fmt.Printf("%s tapchannelrpc FundChannel Error: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(FundChannelErr, err.Error(), nil)
	}

	resp1 := fmt.Sprintf("%s:%d", resp.Txid, resp.OutputIndex)
	AppendFileLog("/data/data/io.bitlong/files/bitlong_api_log.txt", "[sshTapChannels08]", ValueJsonString(resp1))

	return MakeJsonErrorResult(SUCCESS, "", resp1)
}

type AssetChannelDecodeAssetPayResp struct {
	AssetId      string `json:"asset_id"`
	AssetAmount  uint64 `json:"asset_amount"`
	AssetName    string `json:"asset_name"`
	GenesisPoint string `json:"genesis_point"`
	NumSatoshis  int64  `json:"num_satoshis"`
	Description  string `json:"description"`
	Timestamp    int64  `json:"timestamp"`
	Expiry       int64  `json:"expiry"`
	PaymentHash  string `json:"payment_hash"`
	Destination  string `json:"destination"`
	PaymentAddr  []byte `json:"payment_addr"`
}

func AssetChannelDecodeAssetPayReq(assetId string, paymentReq string) string {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return MakeJsonErrorResult(GetConnectionErr, err.Error(), nil)
	}

	defer clearUp()
	client := tapchannelrpc.NewTaprootAssetChannelsClient(conn)

	assetIdStr, err := hex.DecodeString(assetId)
	if err != nil {
		return MakeJsonErrorResult(DecodeStringErr, err.Error(), nil)
	}
	resp, err := client.DecodeAssetPayReq(context.Background(), &tapchannelrpc.AssetPayReq{
		AssetId:      assetIdStr,
		PayReqString: paymentReq,
	})
	if err != nil {
		fmt.Printf("%s tapchannelrpc DecodeAssetPayReq Error: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(DecodeAssetPayReqErr, err.Error(), nil)
	}

	return MakeJsonErrorResult(SUCCESS, "", &AssetChannelDecodeAssetPayResp{
		AssetId:      hex.EncodeToString(resp.GenesisInfo.AssetId),
		AssetAmount:  resp.AssetAmount,
		AssetName:    resp.GenesisInfo.Name,
		GenesisPoint: resp.GenesisInfo.GenesisPoint,
		NumSatoshis:  resp.PayReq.NumSatoshis,
		Description:  resp.PayReq.Description,
		Timestamp:    resp.PayReq.Timestamp,
		Expiry:       resp.PayReq.Expiry,
		PaymentHash:  resp.PayReq.PaymentHash,
		Destination:  resp.PayReq.Description,
		PaymentAddr:  resp.PayReq.PaymentAddr,
	})
}

type Memo struct {
	Acronym     string `json:"acronym"`
	Description string `json:"description"`
	Name        string `json:"name"`
	AssetId     string `json:"asset_id"`
	Amount      string `json:"amount"`
}

func (m *Memo) ToJsonStr() string {
	metastr, _ := json.Marshal(m)
	return string(metastr)
}

func FromJsonStr(jsonStr string) (*Memo, error) {
	memo := Memo{}
	jsonBytes := []byte(jsonStr)
	err := json.Unmarshal(jsonBytes, memo)
	if err != nil {
		return nil, err
	}
	return &memo, nil
}

func MemoToString(memo *Memo) string {
	memoString := memo.ToJsonStr()
	return memoString
}

func StringToMemo(memoString string) (*Memo, error) {
	memo := Memo{}
	jsonBytes := []byte(memoString)
	err := json.Unmarshal(jsonBytes, memo)
	if err != nil {
		return &Memo{
			Acronym:     "",
			Description: memoString,
			Name:        "",
			AssetId:     "",
			Amount:      "",
		}, nil
	}
	return &memo, nil
}

func AssetChannelAddInvoice(assetId string, amount int, pubkey string, memo string) string {
	AppendFileLog("/data/data/io.bitlong/files/bitlong_api_log.txt", "[sshAssetChannelAddInvoice01]", ValueJsonString(fmt.Sprintf("assetId:%s,pubkey:%s,amount:%d", assetId, pubkey, amount)))
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return MakeJsonErrorResult(GetConnectionErr, err.Error(), nil)
	}

	defer clearUp()
	client := tapchannelrpc.NewTaprootAssetChannelsClient(conn)

	assetIdStr, err := hex.DecodeString(assetId)
	if err != nil {
		return MakeJsonErrorResult(DecodeStringErr, err.Error(), nil)
	}

	peerPubkey, err := hex.DecodeString(pubkey)
	if err != nil {
		return MakeJsonErrorResult(DecodeStringErr, err.Error(), nil)
	}

	memoStr, _ := StringToMemo(memo)
	memoStr.AssetId = assetId
	memoStr.Amount = strconv.Itoa(amount)
	memo = MemoToString(memoStr)

	resp, err := client.AddInvoice(context.Background(), &tapchannelrpc.AddInvoiceRequest{
		AssetId:     assetIdStr,
		AssetAmount: uint64(amount),
		PeerPubkey:  peerPubkey,
		InvoiceRequest: &lnrpc.Invoice{
			Memo: memo,
		},
	})
	if err != nil {
		fmt.Printf("%s tapchannelrpc AddInvoice Error: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(AddInvoiceErr, err.Error(), nil)
	}
	AppendFileLog("/data/data/io.bitlong/files/bitlong_api_log.txt", "[sshAssetChannelAddInvoice02]", ValueJsonString(fmt.Sprintf("resp:%s,invoice:%s", resp, resp.InvoiceResult.PaymentRequest)))
	return MakeJsonErrorResult(SUCCESS, "", resp)
}

func AssetChannelSendPayment(assetId string, pubkey string, paymentReq string, timeoutSeconds int, outgoingChanId int, feeLimitSat int, allowSelfPayment bool) string {
	AppendFileLog("/data/data/io.bitlong/files/bitlong_api_log.txt", "[sshAssetChannelSendPayment]", ValueJsonString(fmt.Sprintf("assetId:%s,pubkey:%s,paymentReq:%s,outgoingChanId:%d", assetId, pubkey, paymentReq, outgoingChanId)))
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return MakeJsonErrorResult(GetConnectionErr, err.Error(), nil)
	}

	defer clearUp()
	client := tapchannelrpc.NewTaprootAssetChannelsClient(conn)

	assetIdStr, err := hex.DecodeString(assetId)
	if err != nil {
		return MakeJsonErrorResult(DecodeStringErr, err.Error(), nil)
	}

	peerPubkey, err := hex.DecodeString(pubkey)
	if err != nil {
		return MakeJsonErrorResult(DecodeStringErr, err.Error(), nil)
	}

	qes := tapchannelrpc.SendPaymentRequest{
		AssetId:    assetIdStr,
		PeerPubkey: peerPubkey,
		PaymentRequest: &routerrpc.SendPaymentRequest{
			PaymentRequest:   paymentReq,
			TimeoutSeconds:   int32(timeoutSeconds),
			AllowSelfPayment: allowSelfPayment,
		},
		AllowOverpay: true,
	}
	if outgoingChanId != 0 {
		qes.PaymentRequest.OutgoingChanIds = []uint64{uint64(outgoingChanId)}
	}
	if feeLimitSat == 0 {
		qes.PaymentRequest.FeeLimitSat = 20
	} else if feeLimitSat != 0 {
		qes.PaymentRequest.FeeLimitSat = int64(feeLimitSat)
	}

	resp, err := client.SendPayment(context.Background(), &qes)
	if err != nil {
		fmt.Printf("%s tapchannelrpc SendPayment Error: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(SendPaymentErr, err.Error(), nil)
	}
	for {
		resp1, err := resp.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Printf("%s err == io.EOF, err: %v\n", GetTimeNow(), err)
				return MakeJsonErrorResult(SendPaymentErr, err.Error(), nil)
			}
			fmt.Printf("%s stream Recv err: %v\n", GetTimeNow(), err)
			return MakeJsonErrorResult(SendPaymentErr, err.Error(), nil)
		} else if resp1 != nil {
			resp2 := resp1.GetPaymentResult()
			if resp2 != nil {
				if resp2.Status == 2 {
					return MakeJsonErrorResult(SUCCESS, "", resp2)
				} else if resp2.Status == 3 {
					return MakeJsonErrorResult(SendPaymentErr, resp2.FailureReason.String(), resp2)
				}
			}
		}
	}
}

type getChannelByChanPointResp struct {
	ChannelTypeOne *channelTypeOneResp `json:"channel_type_one"`
	ChannelTypeTwo *channelTypeTwo     `json:"channel_type_two"`
	ChannelId      uint64              `json:"channel_id"`
}

type channelTypeTwo struct {
	TotalLimboBalance          int64                                              `json:"total_limbo_balance"`
	PendingOpenChannel         *lnrpc.PendingChannelsResponse_PendingOpenChannel  `json:"pending_open_channel"`
	PendingForceClosingChannel *lnrpc.PendingChannelsResponse_ForceClosedChannel  `json:"pending_force_closing_channel"`
	WaitingCloseChannel        *lnrpc.PendingChannelsResponse_WaitingCloseChannel `json:"waiting_close_channel"`
	CloseChannel               *lnrpc.ChannelCloseSummary                         `json:"close_channel"`
}

type channelTypeOneResp struct {
	Active               bool                    `json:"active"`
	RemotePubkey         string                  `json:"remote_pubkey"`
	ChannelPoint         string                  `json:"channel_point"`
	ChanId               uint64                  `json:"chan_id"`
	Capacity             int64                   `json:"capacity"`
	LocalBalance         int64                   `json:"local_balance"`
	RemoteBalance        int64                   `json:"remote_balance"`
	CommitFee            int64                   `json:"commit_fee"`
	CommitWeight         int64                   `json:"commit_weight"`
	FeePerKw             int64                   `json:"fee_per_kw"`
	UnsettledBalance     int64                   `json:"unsettled_balance"`
	NumUpdates           uint64                  `json:"num_updates"`
	Initiator            bool                    `json:"initiator"`
	ChanStatusFlags      string                  `json:"chan_status_flags"`
	LocalChanReserveSat  int64                   `json:"local_chan_reserve_sat"`
	RemoteChanReserveSat int64                   `json:"remote_chan_reserve_sat"`
	Lifetime             int64                   `json:"lifetime"`
	Uptime               int64                   `json:"uptime"`
	CloseAddress         string                  `json:"close_address"`
	PushAmountSat        uint64                  `json:"push_amount_sat"`
	Memo                 string                  `json:"memo"`
	CustomChannelData    rfqmsg.JsonAssetChannel `json:"custom_channel_data"`
}

func GetChannelByChanPoint(chanPoint string) string {
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
			var result rfqmsg.JsonAssetChannel
			err := json.Unmarshal(channel.CustomChannelData, &result)
			if err != nil {
				return MakeJsonErrorResult(DecodeStringErr, err.Error(), nil)
			}
			return MakeJsonErrorResult(SUCCESS, "", getChannelByChanPointResp{
				ChannelTypeOne: &channelTypeOneResp{
					Active:            channel.Active,
					RemotePubkey:      channel.RemotePubkey,
					ChannelPoint:      channel.ChannelPoint,
					ChanId:            channel.ChanId,
					Capacity:          channel.Capacity,
					LocalBalance:      channel.LocalBalance,
					RemoteBalance:     channel.RemoteBalance,
					CommitFee:         channel.CommitFee,
					CommitWeight:      channel.CommitWeight,
					FeePerKw:          channel.FeePerKw,
					UnsettledBalance:  channel.UnsettledBalance,
					NumUpdates:        channel.NumUpdates,
					Initiator:         channel.Initiator,
					ChanStatusFlags:   channel.ChanStatusFlags,
					Lifetime:          channel.Lifetime,
					Uptime:            channel.Uptime,
					CloseAddress:      channel.CloseAddress,
					PushAmountSat:     channel.PushAmountSat,
					Memo:              channel.Memo,
					CustomChannelData: result,
				},
				ChannelId: channel.ChanId,
			})
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
			return MakeJsonErrorResult(SUCCESS, "", getChannelByChanPointResp{
				ChannelTypeTwo: &channelTypeTwo{
					TotalLimboBalance:  pendingresponse.TotalLimboBalance,
					PendingOpenChannel: channel,
				},
			})
		}
	}
	for _, channel := range pendingresponse.WaitingCloseChannels {
		if channel.Channel.ChannelPoint == chanPoint {
			return MakeJsonErrorResult(SUCCESS, "", getChannelByChanPointResp{
				ChannelTypeTwo: &channelTypeTwo{
					TotalLimboBalance:   pendingresponse.TotalLimboBalance,
					WaitingCloseChannel: channel,
				},
			})
		}
	}
	for _, channel := range pendingresponse.PendingForceClosingChannels {
		if channel.Channel.ChannelPoint == chanPoint {
			return MakeJsonErrorResult(SUCCESS, "", getChannelByChanPointResp{
				ChannelTypeTwo: &channelTypeTwo{
					TotalLimboBalance:          pendingresponse.TotalLimboBalance,
					PendingForceClosingChannel: channel,
				},
			})
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
			return MakeJsonErrorResult(SUCCESS, "", getChannelByChanPointResp{
				ChannelTypeTwo: &channelTypeTwo{
					TotalLimboBalance: pendingresponse.TotalLimboBalance,
					CloseChannel:      channel,
				},
				ChannelId: channel.ChanId,
			})
		}
	}

	return MakeJsonErrorResult(NoFindChannelErr, "NO_FIND_CHANNEL", nil)
}

type channelAllBalance struct {
	AllSatsBalance        int64 `json:"all_sats_balance"`
	AllAssetBalance       int64 `json:"all_asset_balance"`
	RemoteAllSatsBalance  int64 `json:"remote_all_sats_balance"`
	RemoteAllAssetBalance int64 `json:"remote_all_asset_balance"`
}

func GetAllChannelsBalanceByAssetId(assetId string) string {
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
	var balance channelAllBalance
	for _, channel := range response.Channels {
		if len(channel.CustomChannelData) == 0 {
			continue
		}
		var result rfqmsg.JsonAssetChannel
		err := json.Unmarshal(channel.CustomChannelData, &result)
		if err != nil {
			return MakeJsonErrorResult(DecodeStringErr, err.Error(), nil)
		}
		if result.FundingAssets[0].AssetGenesis.AssetID == assetId {
			balance.AllSatsBalance += channel.LocalBalance
			balance.AllAssetBalance += int64(result.LocalBalance)
			balance.RemoteAllSatsBalance += channel.RemoteBalance
			balance.RemoteAllAssetBalance += int64(result.RemoteBalance)
		}
	}
	return MakeJsonErrorResult(SUCCESS, "", balance)
}

func GetListPeers() string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return MakeJsonErrorResult(GetConnectionErr, err.Error(), nil)
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	resp, err := client.ListPeers(context.Background(), &lnrpc.ListPeersRequest{})
	if err != nil {
		return MakeJsonErrorResult(DecodeStringErr, err.Error(), nil)
	}

	return MakeJsonErrorResult(SUCCESS, "", resp)
}

type simplifyAssetChannelInvoices struct {
	Memo           string                     `json:"memo"`
	Value          int64                      `json:"value"`
	Amount         uint64                     `json:"amount"`
	ChanId         uint64                     `json:"chan_id"`
	CreationDate   int64                      `json:"creation_date"`
	SettleDate     int64                      `json:"settle_date"`
	PaymentRequest string                     `json:"payment_request"`
	Expiry         int64                      `json:"expiry"`
	AddIndex       uint64                     `json:"add_index"`
	SettleIndex    uint64                     `json:"settle_index"`
	State          lnrpc.Invoice_InvoiceState `json:"state"`
	PaymentAddr    []byte                     `json:"payment_addr"`
}

func newSimplifyInvoice(invoice *lnrpc.Invoice) *simplifyAssetChannelInvoices {
	return &simplifyAssetChannelInvoices{
		Memo:           invoice.Memo,
		Value:          invoice.Value,
		CreationDate:   invoice.CreationDate,
		SettleDate:     invoice.SettleDate,
		PaymentRequest: invoice.PaymentRequest,
		Expiry:         invoice.Expiry,
		AddIndex:       invoice.AddIndex,
		SettleIndex:    invoice.SettleIndex,
		State:          invoice.State,
		PaymentAddr:    invoice.PaymentAddr,
	}
}

func createFilterFunc(isBtc, pendingOnly, cancel bool) func(*lnrpc.Invoice) bool {
	switch {
	case pendingOnly && !cancel:
		return func(inv *lnrpc.Invoice) bool {
			return inv.Private != isBtc
		}
	case !pendingOnly && cancel:
		return func(inv *lnrpc.Invoice) bool {
			return inv.Private != isBtc && inv.State == lnrpc.Invoice_CANCELED
		}
	default:
		return func(inv *lnrpc.Invoice) bool { return true }
	}
}

func GetAssetAndBtcChannelListInvoices(pendingOnly bool, cancel bool) string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return MakeJsonErrorResult(GetConnectionErr, err.Error(), nil)
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	resp, err := client.ListInvoices(context.Background(), &lnrpc.ListInvoiceRequest{
		PendingOnly: pendingOnly,
	})
	if err != nil {
		return MakeJsonErrorResult(ListInvoicesErr, err.Error(), nil)
	}

	filteredInvoices := make([]*simplifyAssetChannelInvoices, 0, len(resp.Invoices))
	if pendingOnly && !cancel {
		for _, invoice := range resp.Invoices {
			if invoice.State == lnrpc.Invoice_OPEN {
				filteredInvoices = append(filteredInvoices, newSimplifyInvoice(invoice))
			}
		}
		return MakeJsonErrorResult(SUCCESS, "", filteredInvoices)
	} else if !pendingOnly && cancel {
		for _, invoice := range resp.Invoices {
			if invoice.State == lnrpc.Invoice_CANCELED {
				filteredInvoices = append(filteredInvoices, newSimplifyInvoice(invoice))
			}
		}
		return MakeJsonErrorResult(SUCCESS, "", filteredInvoices)
	}

	return MakeJsonErrorResult(SUCCESS, "", filteredInvoices)
}

type simplifyAssetChannelPayments struct {
	AssetId        string                      `json:"asset_id"`
	Amount         uint64                      `json:"amount"`
	RfqID          string                      `json:"rfq_id"`
	ChanId         uint64                      `json:"chan_id"`
	ValueSat       int64                       `json:"value_sat"`
	CreationTimeNs int64                       `json:"creation_time_ns"`
	PaymentRequest string                      `json:"payment_request"`
	PaymentIndex   uint64                      `json:"payment_index"`
	Status         lnrpc.Payment_PaymentStatus `json:"status"`
	FailureReason  lnrpc.PaymentFailureReason  `json:"failure_reason"`
}

var (
	invoiceResultCache = make(map[string][]*simplifyAssetChannelInvoices)
	paymentResultCache = make(map[string][]*simplifyAssetChannelPayments)
	cacheMutex         sync.RWMutex
)

func getCacheKey(isAllChannel bool, assetId string, chanId int, maxPayments int) string {
	if isAllChannel {
		return fmt.Sprintf("ALL_%s_%d", assetId, maxPayments)
	}
	return fmt.Sprintf("SINGLE_%d_%d", chanId, maxPayments)
}

func paginate(payments []*simplifyAssetChannelPayments, offset, limit uint64) []*simplifyAssetChannelPayments {
	if offset >= uint64(len(payments)) {
		return nil
	}
	end := offset + limit
	if end > uint64(len(payments)) {
		end = uint64(len(payments))
	}
	return payments[offset:end]
}

func GetAssetChannelListPayments(isUpdate bool, isAllChannel bool, assetId string, chanId int, indexOffset int, maxPayments int) string {
	if maxPayments <= 0 {
		maxPayments = 10000
	}
	if indexOffset < 0 {
		return MakeJsonErrorResult(InvalidParamsErr, "invalid negative offset", nil)
	}

	if isAllChannel {
		if assetId == "" {
			return MakeJsonErrorResult(InvalidParamsErr, "assetId is required for all channels query", nil)
		}
		if chanId != 0 {
			return MakeJsonErrorResult(InvalidParamsErr, "channelId must be 0 for all channels query", nil)
		}
	} else {
		if chanId <= 0 {
			return MakeJsonErrorResult(InvalidParamsErr, "invalid channel ID", nil)
		}
		assetId = ""
	}

	cacheKey := getCacheKey(isAllChannel, assetId, chanId, maxPayments)

	if !isUpdate {
		cacheMutex.RLock()
		cached, ok := paymentResultCache[cacheKey]
		if ok {
			paginated := paginate(cached, uint64(indexOffset), uint64(maxPayments))
			cacheMutex.RUnlock()
			return MakeJsonErrorResult(SUCCESS, "", paginated)
		}
		cacheMutex.RUnlock()
	}

	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return MakeJsonErrorResult(ListPaymentsErr, err.Error(), nil)
	}
	defer clearUp()

	client := lnrpc.NewLightningClient(conn)
	resp, err := client.ListPayments(context.Background(), &lnrpc.ListPaymentsRequest{})
	if err != nil {
		return MakeJsonErrorResult(ListPaymentsErr, err.Error(), nil)
	}

	payments := make([]*simplifyAssetChannelPayments, 0, len(resp.Payments))
	for _, payment := range resp.Payments {
		if inv, ok := filterPayments(payment, assetId, uint64(chanId), isAllChannel); ok {
			payments = append(payments, inv)
		}
	}

	sort.Slice(payments, func(i, j int) bool {
		return payments[i].CreationTimeNs > payments[j].CreationTimeNs
	})

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if isUpdate || paymentResultCache[cacheKey] == nil {
		paymentResultCache[cacheKey] = payments
	} else {
		merged := append(paymentResultCache[cacheKey], payments...)
		seen := make(map[string]struct{})
		var unique []*simplifyAssetChannelPayments
		for _, p := range merged {
			key := fmt.Sprintf("%d|%d", p.ChanId, p.CreationTimeNs)
			if _, exists := seen[key]; !exists {
				seen[key] = struct{}{}
				unique = append(unique, p)
			}
		}
		paymentResultCache[cacheKey] = unique
	}

	result := paginate(paymentResultCache[cacheKey], uint64(indexOffset), uint64(maxPayments))
	if len(result) == 0 {
		return MakeJsonErrorResult(ListPaymentsErr, "no more results", nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", result)
}

func filterPayments(payment *lnrpc.Payment, assetId string, chanId uint64, isAllChannel bool) (*simplifyAssetChannelPayments, bool) {
	if payment.FailureReason != 0 {
		return nil, false
	}

	if len(payment.Htlcs) == 0 {
		return nil, false
	}

	for _, htlc := range payment.Htlcs {
		if htlc.Route == nil || len(htlc.Route.Hops) == 0 {
			continue
		}
		targetChanId := htlc.Route.Hops[0].ChanId

		var result rfqmsg.JsonHtlc
		if err := json.Unmarshal(htlc.Route.CustomChannelData, &result); err != nil {
			continue // 静默处理解析错误
		}

		if len(result.Balances) == 0 {
			continue
		}

		currentAssetId := result.Balances[0].AssetID

		if isAllChannel {
			if currentAssetId == assetId {
				return buildPaymentStruct(payment, result, targetChanId), true
			}
		} else {
			if targetChanId == chanId {
				return buildPaymentStruct(payment, result, targetChanId), true
			}
		}
	}
	return nil, false
}

func buildPaymentStruct(payment *lnrpc.Payment, result rfqmsg.JsonHtlc, chanId uint64) *simplifyAssetChannelPayments {
	return &simplifyAssetChannelPayments{
		AssetId:        result.Balances[0].AssetID,
		Amount:         result.Balances[0].Amount,
		RfqID:          result.RfqID,
		ChanId:         chanId,
		ValueSat:       payment.ValueSat,
		CreationTimeNs: payment.CreationTimeNs,
		PaymentRequest: payment.PaymentRequest,
		PaymentIndex:   payment.PaymentIndex,
		Status:         payment.Status,
		FailureReason:  payment.FailureReason,
	}
}

func paginateInvoice(invoices []*simplifyAssetChannelInvoices, offset, limit uint64) []*simplifyAssetChannelInvoices {
	if offset >= uint64(len(invoices)) {
		return nil
	}
	end := offset + limit
	if end > uint64(len(invoices)) {
		end = uint64(len(invoices))
	}
	return invoices[offset:end]
}

func GetAssetChannelListInvoicesSettled(isUpdate bool, isAllChannel bool, assetId string, chanId int, indexOffset int, maxInvoice int) string {
	if maxInvoice <= 0 {
		maxInvoice = 6000
	}
	if indexOffset < 0 {
		return MakeJsonErrorResult(InvalidParamsErr, "invalid negative offset", nil)
	}

	if isAllChannel {
		if assetId == "" {
			return MakeJsonErrorResult(InvalidParamsErr, "assetId is required for all channels query", nil)
		}
		if chanId != 0 {
			return MakeJsonErrorResult(InvalidParamsErr, "channelId must be 0 for all channels query", nil)
		}
	} else {
		if chanId <= 0 {
			return MakeJsonErrorResult(InvalidParamsErr, "invalid channel ID", nil)
		}
		assetId = ""
	}

	cacheKey := getCacheKey(isAllChannel, assetId, chanId, maxInvoice)

	if !isUpdate {
		cacheMutex.RLock()
		cached, ok := invoiceResultCache[cacheKey]
		if ok {
			paginated := paginateInvoice(cached, uint64(indexOffset), uint64(maxInvoice))
			cacheMutex.RUnlock()
			return MakeJsonErrorResult(SUCCESS, "", paginated)
		}
		cacheMutex.RUnlock()
	}
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return MakeJsonErrorResult(GetConnectionErr, err.Error(), nil)
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	resp, err := client.ListInvoices(context.Background(), &lnrpc.ListInvoiceRequest{})
	if err != nil {
		return MakeJsonErrorResult(ListInvoicesErr, err.Error(), nil)
	}

	invoices := make([]*simplifyAssetChannelInvoices, 0, len(resp.Invoices))
	for _, invoice := range resp.Invoices {
		if invoice.State == 1 {
			if inv, ok := filterInvoices(invoice, assetId, uint64(chanId), isAllChannel); ok {
				invoices = append(invoices, inv)
			}
		}
	}

	sort.Slice(invoices, func(i, j int) bool {
		return invoices[i].CreationDate > invoices[j].CreationDate
	})

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if isUpdate || invoiceResultCache[cacheKey] == nil {
		invoiceResultCache[cacheKey] = invoices
	} else {
		merged := append(invoiceResultCache[cacheKey], invoices...)
		seen := make(map[string]struct{})
		var unique []*simplifyAssetChannelInvoices
		for _, p := range merged {
			key := fmt.Sprintf("%d|%d", p.ChanId, p.CreationDate)
			if _, exists := seen[key]; !exists {
				seen[key] = struct{}{}
				unique = append(unique, p)
			}
		}
		invoiceResultCache[cacheKey] = unique
	}

	result := paginateInvoice(invoiceResultCache[cacheKey], uint64(indexOffset), uint64(maxInvoice))
	if len(result) == 0 {
		return MakeJsonErrorResult(ListInvoicesErr, "no more results", nil)
	}

	return MakeJsonErrorResult(SUCCESS, "", invoices)
}

func filterInvoices(invoice *lnrpc.Invoice, assetId string, chanId uint64, isAllChannel bool) (*simplifyAssetChannelInvoices, bool) {
	if len(invoice.Htlcs) == 0 {
		return nil, false
	}

	for _, htlc := range invoice.Htlcs {
		targetChanId := htlc.ChanId
		if assetId == "" && targetChanId == chanId && htlc.CustomChannelData == nil {
			return buildInvoiceStruct(invoice, rfqmsg.JsonHtlc{
				Balances: []*rfqmsg.JsonAssetTranche{
					{AssetID: "00", Amount: 0},
				},
				RfqID: ""}, targetChanId), true
		}

		var result rfqmsg.JsonHtlc
		if err := json.Unmarshal(htlc.CustomChannelData, &result); err != nil {
			continue // 静默处理解析错误
		}

		if len(result.Balances) == 0 {
			continue
		}

		currentAssetId := result.Balances[0].AssetID

		if isAllChannel {
			if currentAssetId == assetId {
				return buildInvoiceStruct(invoice, result, targetChanId), true
			}
		} else {
			if targetChanId == chanId {
				return buildInvoiceStruct(invoice, result, targetChanId), true
			}
		}
	}
	return nil, false
}

func buildInvoiceStruct(invoice *lnrpc.Invoice, result rfqmsg.JsonHtlc, chanId uint64) *simplifyAssetChannelInvoices {
	return &simplifyAssetChannelInvoices{
		Memo:           invoice.Memo,
		Value:          invoice.Value,
		Amount:         result.Balances[0].Amount,
		ChanId:         chanId,
		CreationDate:   invoice.CreationDate,
		SettleDate:     invoice.SettleDate,
		PaymentRequest: invoice.PaymentRequest,
		Expiry:         invoice.Expiry,
		AddIndex:       invoice.AddIndex,
		SettleIndex:    invoice.SettleIndex,
		State:          invoice.State,
		PaymentAddr:    invoice.PaymentAddr,
	}
}

type ChannelInfo struct {
	ChannelID    uint64 `json:"channel_id"`
	ChannelPoint string `json:"channel_point"`
}

type AssetChannelGroup struct {
	AssetID  string        `json:"asset_id"`
	Channels []ChannelInfo `json:"channels"`
}

type GetChannelsInfoResp struct {
	AssetIds []string            `json:"asset_ids"`
	Group    []AssetChannelGroup `json:"group"`
}

func GetChannelListInfo() string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return MakeJsonErrorResult(GetConnectionErr, err.Error(), nil)
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	resp, err := client.ListChannels(context.Background(), &lnrpc.ListChannelsRequest{})
	if err != nil {
		AppendFileLog("/data/data/io.bitlong/files/bitlong_api_log.txt", "[GetChannelIdsAndPoints01]", ValueJsonString(err))
		return MakeJsonErrorResult(ListChannelsErr, err.Error(), nil)
	}

	result := &GetChannelsInfoResp{
		AssetIds: make([]string, 0),
		Group:    make([]AssetChannelGroup, 0),
	}

	assetGroups := make(map[string][]ChannelInfo)
	assetIDSet := make(map[string]struct{}) // 用于去重asset_ids

	btcChannels := make([]ChannelInfo, 0)

	for _, channel := range resp.Channels {
		channelInfo := ChannelInfo{
			ChannelID:    channel.ChanId,
			ChannelPoint: channel.ChannelPoint,
		}

		if channel.CustomChannelData == nil {
			btcChannels = append(btcChannels, channelInfo)
			continue
		}

		var customData rfqmsg.JsonAssetChannel
		if err := json.Unmarshal(channel.CustomChannelData, &customData); err != nil {
			continue
		}

		if len(customData.FundingAssets) == 0 {
			continue
		}

		assetID := customData.FundingAssets[0].AssetGenesis.AssetID

		if _, exists := assetIDSet[assetID]; !exists {
			result.AssetIds = append(result.AssetIds, assetID)
			assetIDSet[assetID] = struct{}{}
		}

		assetGroups[assetID] = append(assetGroups[assetID], channelInfo)
	}

	if len(btcChannels) > 0 {
		result.Group = append(result.Group, AssetChannelGroup{
			AssetID:  "00",
			Channels: btcChannels,
		})
	}

	for assetID, channels := range assetGroups {
		result.Group = append(result.Group, AssetChannelGroup{
			AssetID:  assetID,
			Channels: channels,
		})
	}

	sort.Strings(result.AssetIds)
	sort.Slice(result.Group, func(i, j int) bool {
		if result.Group[i].AssetID == "00" {
			return false // BTC放在最后
		}
		if result.Group[j].AssetID == "00" {
			return true
		}
		return result.Group[i].AssetID < result.Group[j].AssetID
	})

	return MakeJsonErrorResult(SUCCESS, "", result)
}

type PriceOracleClient struct {
	client priceoraclerpc.PriceOracleClient
	conn   *grpc.ClientConn
}

func NewPriceOracleClient(grpcAddr string) (*PriceOracleClient, error) {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	creds := credentials.NewTLS(tlsConfig)

	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("grpc connection failed: %v", err)
	}

	return &PriceOracleClient{
		client: priceoraclerpc.NewPriceOracleClient(conn),
		conn:   conn,
	}, nil
}

func (p *PriceOracleClient) Close() {
	p.conn.Close()
}
func GetAssetRates(assetId string) string {
	price, err := NewPriceOracleClient("118.24.37.253:10086")
	if err != nil {
		return MakeJsonErrorResult(QueryAssetRatesErr, err.Error(), nil)
	}
	defer price.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	subjectAssetIdBytes, _ := hex.DecodeString(assetId)

	paymentAssetId := "0000000000000000000000000000000000000000000000000000000000000000"
	paymentAssetIdBytes, _ := hex.DecodeString(paymentAssetId)

	req := &priceoraclerpc.QueryAssetRatesRequest{
		TransactionType: priceoraclerpc.TransactionType_SALE,
		SubjectAsset: &priceoraclerpc.AssetSpecifier{
			Id: &priceoraclerpc.AssetSpecifier_AssetId{
				AssetId: subjectAssetIdBytes,
			},
		},
		PaymentAsset: &priceoraclerpc.AssetSpecifier{
			Id: &priceoraclerpc.AssetSpecifier_AssetId{
				AssetId: paymentAssetIdBytes,
			},
		},
		PaymentAssetMaxAmount: 600000,
	}

	res, err := price.client.QueryAssetRates(ctx, req)
	if err != nil {
		return MakeJsonErrorResult(QueryAssetRatesErr, err.Error(), nil)
	}
	switch result := res.Result.(type) {
	case *priceoraclerpc.QueryAssetRatesResponse_Ok:
		if result.Ok.AssetRates.SubjectAssetRate.Coefficient == "" {
			result.Ok.AssetRates.SubjectAssetRate.Coefficient = result.Ok.AssetRates.PaymentAssetRate.Coefficient
		}
		subjectCoeff, err1 := strconv.ParseInt(result.Ok.AssetRates.SubjectAssetRate.Coefficient, 10, 64)
		paymentCoeff, err2 := strconv.ParseInt(result.Ok.AssetRates.PaymentAssetRate.Coefficient, 10, 64)
		if err1 != nil || err2 != nil {
			return MakeJsonErrorResult(QueryAssetRatesErr, "Invalid coefficient format", nil)
		}

		rateFloat := float64(subjectCoeff) / float64(paymentCoeff)
		return MakeJsonErrorResult(SUCCESS, "", rateFloat)
	case *priceoraclerpc.QueryAssetRatesResponse_Error:
		return MakeJsonErrorResult(QueryAssetRatesErr, result.Error.Message, nil)
	default:
		return MakeJsonErrorResult(QueryAssetRatesErr, "Unknown result type", nil)
	}
}

func KeySendToAssetChannel(assetId string, amount string, pubkey string, outgoingChanId string) string {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return MakeJsonErrorResult(GetConnectionErr, err.Error(), nil)
	}

	defer clearUp()

	if assetId == "" || amount == "" || amount == "0" || pubkey == "" || outgoingChanId == "" {
		return MakeJsonErrorResult(InvalidParamsErr, "invalid params", nil)
	}

	client := tapchannelrpc.NewTaprootAssetChannelsClient(conn)

	assetIdStr, err := hex.DecodeString(assetId)
	if err != nil {
		return MakeJsonErrorResult(DecodeStringErr, err.Error(), nil)
	}

	peerPubkey, err := hex.DecodeString(pubkey)
	if err != nil {
		return MakeJsonErrorResult(DecodeStringErr, err.Error(), nil)
	}

	amountUint, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		return MakeJsonErrorResult(InvalidParamsErr, err.Error(), nil)
	}

	outgoingChanIdUint, err := strconv.ParseUint(outgoingChanId, 10, 64)
	if err != nil {
		return MakeJsonErrorResult(InvalidParamsErr, err.Error(), nil)
	}

	req := &tapchannelrpc.SendPaymentRequest{
		AssetId:     assetIdStr,
		AssetAmount: amountUint,
		PaymentRequest: &routerrpc.SendPaymentRequest{
			Dest:              peerPubkey,
			Amt:               int64(rfqmath.DefaultOnChainHtlcSat),
			OutgoingChanId:    outgoingChanIdUint,
			TimeoutSeconds:    30,
			DestCustomRecords: make(map[uint64][]byte),
		},
	}

	destRecords := req.PaymentRequest.DestCustomRecords
	_, isKeysend := destRecords[record.KeySendType]
	var rHash []byte
	var preimage lntypes.Preimage
	if _, err := rand.Read(preimage[:]); err != nil {
		return MakeJsonErrorResult(SendPaymentErr, err.Error(), nil)
	}
	if !isKeysend {
		destRecords[record.KeySendType] = preimage[:]
		hash := preimage.Hash()
		rHash = hash[:]

		req.PaymentRequest.PaymentHash = rHash

	}
	resp, err := client.SendPayment(context.Background(), req)
	if err != nil {
		fmt.Printf("%s tapchannelrpc SendPayment Error: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(SendPaymentErr, err.Error(), nil)
	}
	for {
		resp1, err := resp.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Printf("%s err == io.EOF, err: %v\n", GetTimeNow(), err)
				return MakeJsonErrorResult(SendPaymentErr, err.Error(), nil)
			}
			fmt.Printf("%s stream Recv err: %v\n", GetTimeNow(), err)
			return MakeJsonErrorResult(SendPaymentErr, err.Error(), nil)
		} else if resp1 != nil {
			resp2 := resp1.GetPaymentResult()
			if resp2 != nil {
				if resp2.Status == 2 {
					return MakeJsonErrorResult(SUCCESS, "", resp2)
				} else if resp2.Status == 3 {
					return MakeJsonErrorResult(SendPaymentErr, resp2.FailureReason.String(), resp2)
				}
			}
		}
	}
}

func KeySendToBtcChannel(amt string, pubkey string, outgoingChanId string) string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return MakeJsonErrorResult(GetConnectionErr, err.Error(), nil)
	}
	defer clearUp()

	client := routerrpc.NewRouterClient(conn)

	if amt == "" || amt == "0" {
		return MakeJsonErrorResult(InvalidParamsErr, "invalid amount", nil)
	}

	if pubkey == "" {
		return MakeJsonErrorResult(InvalidParamsErr, "invalid pubkey", nil)
	}

	if outgoingChanId == "" {
		return MakeJsonErrorResult(InvalidParamsErr, "invalid channel ID", nil)
	}

	amtInt, err := strconv.ParseInt(amt, 10, 64)
	if err != nil {
		return MakeJsonErrorResult(InvalidParamsErr, err.Error(), nil)
	}

	peerPubkey, err := hex.DecodeString(pubkey)
	if err != nil {
		return MakeJsonErrorResult(DecodeStringErr, err.Error(), nil)
	}

	outgoingChanIdUint, err := strconv.ParseUint(outgoingChanId, 10, 64)
	if err != nil {
		return MakeJsonErrorResult(InvalidParamsErr, err.Error(), nil)
	}

	req := &routerrpc.SendPaymentRequest{
		Dest:              peerPubkey,
		Amt:               amtInt,
		OutgoingChanId:    outgoingChanIdUint,
		TimeoutSeconds:    30,
		DestCustomRecords: make(map[uint64][]byte),
	}

	destRecords := req.DestCustomRecords
	_, isKeysend := destRecords[record.KeySendType]
	var rHash []byte
	var preimage lntypes.Preimage
	if _, err := rand.Read(preimage[:]); err != nil {
		return MakeJsonErrorResult(SendPaymentErr, err.Error(), nil)
	}
	if !isKeysend {
		destRecords[record.KeySendType] = preimage[:]
		hash := preimage.Hash()
		rHash = hash[:]

		req.PaymentHash = rHash
	}
	resp, err := client.SendPaymentV2(context.Background(), req)
	if err != nil {
		fmt.Printf("%s routerrpc SendPaymentV2 Error: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(SendPaymentV2Err, err.Error(), nil)
	}

	for {
		response, err := resp.Recv()
		if err != nil {
			return MakeJsonErrorResult(SendPaymentV2Err, err.Error(), nil)
		} else if response != nil {
			switch response.Status {
			case 1: // IN_FLIGHT
				continue
			case 2: // SUCCESS
				return MakeJsonErrorResult(SUCCESS, "", response)
			case 3: // FAILED
				return MakeJsonErrorResult(SendPaymentV2Err, response.FailureReason.String(), response)
			default:
				return MakeJsonErrorResult(SendPaymentV2Err, fmt.Sprintf("未知状态: %d", response.Status), response)
			}
		}
	}
}

func ForceCloseChannel(channelPoint string) string {
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

	fundingTxidStr := parts[0]

	outputIndex, err := strconv.Atoi(parts[1])
	if err != nil {
		return MakeJsonErrorResult(CloseChannelErr, "Invalid output index", nil)
	}

	request := &lnrpc.CloseChannelRequest{
		ChannelPoint: &lnrpc.ChannelPoint{
			FundingTxid: &lnrpc.ChannelPoint_FundingTxidStr{FundingTxidStr: fundingTxidStr},
			OutputIndex: uint32(outputIndex),
		},
		Force: true,
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
			return MakeJsonErrorResult(SUCCESS, "", "通道进入强制关闭状态")
		}
	}
}
