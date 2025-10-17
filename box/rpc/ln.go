package rpc

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/wallet/box/loggers"
	"github.com/wallet/box/models"
	"github.com/wallet/box/sc"
	"gopkg.in/resty.v1"

	"github.com/lightninglabs/taproot-assets/rfqmsg"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	"github.com/lightningnetwork/lnd/lnrpc/walletrpc"
	"github.com/lightningnetwork/lnd/lntypes"
	"github.com/lightningnetwork/lnd/record"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/api"
	btlmassagerpc "github.com/wallet/box/btlmassage"
	"google.golang.org/grpc/status"
)

type Ln struct{}

const (
	Default = "default"
)

func OpenChan(pubkey string, host string, localFundingAmount int64, satPerVByte uint64, pushSat int64, memo string) (string, error) {

	var l Ln
	if _, err := l.ConnectPeer(pubkey, host); err != nil && !strings.Contains(err.Error(), "already connected to peer") {
		return "", errors.Wrap(err, "l.ConnectPeer")
	}

	cp, err := l.OpenChannel(pubkey, localFundingAmount, satPerVByte, pushSat, memo)
	if err != nil {
		return "", errors.Wrap(err, "l.OpenChannel")
	}

	t, o := cp.GetTxid(), cp.GetOutputIndex()

	return fmt.Sprintf("%s:%d", api.TxHashEncodeToString(t), o), nil

}

func (l Ln) GetInfo() (*lnrpc.GetInfoResponse, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	lc := lnrpc.NewLightningClient(conn)
	req := &lnrpc.GetInfoRequest{}

	resp, err := lc.GetInfo(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "lc.GetInfo")
	}
	return resp, nil
}

func (l Ln) GetState() (*lnrpc.GetStateResponse, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	sc := lnrpc.NewStateClient(conn)
	req := &lnrpc.GetStateRequest{}

	resp, err := sc.GetState(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "sc.GetState")
	}
	return resp, nil
}

func (l Ln) WalletBalance() (*lnrpc.WalletBalanceResponse, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	lc := lnrpc.NewLightningClient(conn)
	req := &lnrpc.WalletBalanceRequest{}

	resp, err := lc.WalletBalance(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "lc.WalletBalance")
	}
	return resp, nil
}

func GetWalletBalance() (*api.WalletBalanceResponse, error) {
	var l Ln
	resp, err := l.WalletBalance()
	if err != nil {
		return nil, errors.Wrap(err, "l.ChannelBalance")
	}
	resp, err = api.ProcessGetWalletBalanceResult(resp)
	if err != nil {
		return nil, errors.Wrap(err, "api.ProcessGetWalletBalanceResult")
	}
	return &api.WalletBalanceResponse{
		TotalBalance:       int(resp.TotalBalance),
		ConfirmedBalance:   int(resp.ConfirmedBalance),
		UnconfirmedBalance: int(resp.UnconfirmedBalance),
		LockedBalance:      int(resp.LockedBalance),
	}, nil

}

func GetWalletBalanceResponse() (*models.WalletBalanceResponse, error) {
	var l Ln
	resp, err := l.WalletBalance()
	if err != nil {
		return nil, errors.Wrap(err, "l.WalletBalance")
	}

	var df models.Default
	ab := resp.AccountBalance
	if d, ok := ab[Default]; ok {
		df = models.Default{
			ConfirmedBalance:   d.ConfirmedBalance,
			UnconfirmedBalance: d.UnconfirmedBalance,
		}
	}

	return &models.WalletBalanceResponse{
		TotalBalance:              resp.TotalBalance,
		ConfirmedBalance:          resp.ConfirmedBalance,
		UnconfirmedBalance:        resp.UnconfirmedBalance,
		LockedBalance:             resp.LockedBalance,
		ReservedBalanceAnchorChan: resp.ReservedBalanceAnchorChan,
		AccountBalance: models.AccountBalance{
			Default: df,
		},
	}, nil

}

func (l Ln) ConnectPeer(pubkey string, host string) (*lnrpc.ConnectPeerResponse, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	lc := lnrpc.NewLightningClient(conn)
	req := &lnrpc.ConnectPeerRequest{
		Addr: &lnrpc.LightningAddress{
			Pubkey: pubkey,
			Host:   host,
		},
	}

	resp, err := lc.ConnectPeer(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "lc.ConnectPeer")
	}
	return resp, nil
}

func (l Ln) OpenChannel(nodePubkey string, localFundingAmount int64, satPerVByte uint64, pushSat int64, memo string) (*lnrpc.PendingUpdate, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	lc := lnrpc.NewLightningClient(conn)

	_nodePubkey, err := hex.DecodeString(nodePubkey)
	if err != nil {
		return nil, errors.Wrap(err, "hex.DecodeString")
	}
	req := &lnrpc.OpenChannelRequest{
		SatPerVbyte:        satPerVByte,
		NodePubkey:         _nodePubkey,
		LocalFundingAmount: localFundingAmount,
		PushSat:            pushSat,
		Memo:               memo,
	}

	stream, err := lc.OpenChannel(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "lc.OpenChannel")
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return nil, nil
		}

		if err != nil {
			if st, ok := status.FromError(err); ok {
				return nil, errors.Wrap(err, fmt.Sprintf("stream.Recv: code=%v, message=%v", st.Code(), st.Message()))
			} else {
				return nil, errors.Wrap(err, "stream.Recv")
			}
		}

		if resp.GetChanPending() != nil {
			return resp.GetChanPending(), nil
		}
	}

}

func (l Ln) ListPeers() (*lnrpc.ListPeersResponse, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	lc := lnrpc.NewLightningClient(conn)
	req := &lnrpc.ListPeersRequest{}

	resp, err := lc.ListPeers(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "lc.ListPeers")
	}
	return resp, nil
}

func (l Ln) ListChannels(activeOnly bool, private bool) (*lnrpc.ListChannelsResponse, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	lc := lnrpc.NewLightningClient(conn)
	req := &lnrpc.ListChannelsRequest{
		ActiveOnly:  activeOnly,
		PrivateOnly: private,
	}

	resp, err := lc.ListChannels(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "lc.ListChannels")
	}
	return resp, nil
}

func (l Ln) ListInvoices(reversed bool, numMaxInvoices uint64) (*lnrpc.ListInvoiceResponse, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	lc := lnrpc.NewLightningClient(conn)
	req := &lnrpc.ListInvoiceRequest{
		Reversed:       reversed,
		NumMaxInvoices: numMaxInvoices,
	}

	resp, err := lc.ListInvoices(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "lc.ListInvoices")
	}
	return resp, nil
}

func (l Ln) SendPaymentV2(timeoutSeconds int32, dest []byte, amt int64, feeLimitSat int64, destCustomRecords map[uint64][]byte, outgoingChanIds []uint64, paymentHash []byte) (*lnrpc.Payment, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	rc := routerrpc.NewRouterClient(conn)
	req := &routerrpc.SendPaymentRequest{
		TimeoutSeconds:    timeoutSeconds,
		Dest:              dest,
		Amt:               amt,
		FeeLimitSat:       feeLimitSat,
		DestCustomRecords: destCustomRecords,
		OutgoingChanIds:   outgoingChanIds,
		PaymentHash:       paymentHash,
	}

	ctxt, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	stream, err := rc.SendPaymentV2(ctxt, req)
	if err != nil {
		return nil, errors.Wrap(err, "rc.SendPaymentV2")
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return nil, nil
		}

		if err != nil {
			if st, ok := status.FromError(err); ok {
				return nil, errors.Wrap(err, fmt.Sprintf("stream.Recv: code=%v, message=%v", st.Code(), st.Message()))
			} else {
				return nil, errors.Wrap(err, "stream.Recv")
			}
		}

		return resp, nil
	}
}

var (
	invalidDestNodePubkey = errors.New("dest node pubkey must be exactly 33 bytes")
	invalidRespStatusCode = errors.New("invalid response status code, expected 2 or 3")
)

func SendPaymentV2ByKeySend(dest string, amt int64, feeLimitSat int64, outgoingChanId uint64) (*lnrpc.Payment, error) {

	destNode, err := hex.DecodeString(dest)
	if err != nil {
		return nil, errors.Wrap(err, "hex.DecodeString")
	}
	if len(destNode) != 33 {
		return nil, errors.Wrap(invalidDestNodePubkey, strconv.Itoa(len(dest)))
	}
	destCustomRecords := make(map[uint64][]byte)
	outgoingChanIds := make([]uint64, 1)

	outgoingChanIds[0] = outgoingChanId

	var rHash []byte
	var preimage lntypes.Preimage

	if _, err := rand.Read(preimage[:]); err != nil {
		return nil, err
	}

	destCustomRecords[record.KeySendType] = preimage[:]
	hash := preimage.Hash()
	rHash = hash[:]
	paymentHash := rHash

	var l Ln
	resp, err := l.SendPaymentV2(30, destNode, amt, feeLimitSat, destCustomRecords, outgoingChanIds, paymentHash)
	if err != nil {
		return nil, errors.Wrap(err, "l.SendPaymentV2")
	}

	if resp.Status == 2 || resp.Status == 3 || resp.Status == 1 {
		return resp, nil
	}

	return nil, errors.Wrap(invalidRespStatusCode, strconv.Itoa(int(resp.Status)))
}

func (l Ln) CheckHTLCIsAsset(htlcid uint64) error {
	conn, err := GetConn(tap, false)
	if err != nil {
		return errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	client := btlmassagerpc.NewBtlMassageClient(conn)
	req := btlmassagerpc.BackDustRequest{
		HtlcIndex: htlcid,
	}
	_, err = client.BackDust(context.Background(), &req)
	if err != nil {
		return errors.Wrap(err, "BackDust")
	}
	return nil
}

func (l Ln) SubscribeHtlcEvents(f func(e *routerrpc.HtlcEvent)) error {
	conn, err := GetConn(ln, false)
	if err != nil {
		return errors.Wrap(err, "GetConn")
	}
	defer Close(conn)
	client := routerrpc.NewRouterClient(conn)

	if f == nil {
		f = func(e *routerrpc.HtlcEvent) {
			logrus.Infoln(e)
		}
	}
	request := &routerrpc.SubscribeHtlcEventsRequest{}
	stream, err := client.SubscribeHtlcEvents(context.Background(), request)
	if err != nil {
		return errors.Wrap(err, "SubscribeHtlcEvents")
	}
	for {
		event, err := stream.Recv()
		if err != nil {
			return errors.Wrap(err, "stream.Recv")
		}
		if event != nil {
			f(event)
		}
	}
}

func (l Ln) ChannelBalance() (*lnrpc.ChannelBalanceResponse, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	lc := lnrpc.NewLightningClient(conn)
	req := &lnrpc.ChannelBalanceRequest{}

	resp, err := lc.ChannelBalance(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "lc.ChannelBalance")
	}
	return resp, nil
}

func (l Ln) ListUnspent() (*walletrpc.ListUnspentResponse, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	wkc := walletrpc.NewWalletKitClient(conn)
	req := &walletrpc.ListUnspentRequest{}

	resp, err := wkc.ListUnspent(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "wkc.ListUnspent")
	}
	return resp, nil
}

func GetListUnspent() (Utxos []*api.UnspentUtxo, err error) {
	var l Ln
	utxos, err := l.ListUnspent()
	if err != nil {
		return nil, errors.Wrap(err, "l.ListUnspent")
	}
	return api.ListUnspentResponseToUnspentUtxos(utxos), nil
}

type BoxListInvoicesRecords struct {
	Invoices    []*lnrpc.Invoice
	TotalAmount int64
}

func (l Ln) BoxListInvoicesRecords(assetId string, PendingOnly bool, NumMaxInvoices uint64, CreationDateStart uint64, CreationDateEnd uint64) (BoxListInvoicesRecords, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return BoxListInvoicesRecords{}, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)
	if NumMaxInvoices == 0 {
		NumMaxInvoices = 1000
	}

	ctx := context.Background()
	lc := lnrpc.NewLightningClient(conn)
	req := &lnrpc.ListInvoiceRequest{
		PendingOnly:       PendingOnly,
		NumMaxInvoices:    NumMaxInvoices,
		CreationDateStart: CreationDateStart,
		CreationDateEnd:   CreationDateEnd,
	}

	resp, err := lc.ListInvoices(ctx, req)
	if err != nil {
		return BoxListInvoicesRecords{}, errors.Wrap(err, "lc.ListInvoices")
	}

	var records []*lnrpc.Invoice
	var totalAmount int64

	if PendingOnly {
		for _, invoice := range resp.Invoices {
			if assetId == "" {
				records = append(records, invoice)
			} else if assetId == "00" && !invoice.Private {
				records = append(records, invoice)
			} else if assetId == "01" && invoice.Private {
				records = append(records, invoice)
			}
		}
		for i, j := 0, len(records)-1; i < j; i, j = i+1, j-1 {
			records[i], records[j] = records[j], records[i]
		}
		return BoxListInvoicesRecords{
			Invoices:    records,
			TotalAmount: totalAmount,
		}, nil
	}

	if assetId == "" {
		for _, invoice := range resp.Invoices {
			if invoice.State == lnrpc.Invoice_SETTLED {
				records = append(records, invoice)
				totalAmount += invoice.Value
			}
		}
		for i, j := 0, len(records)-1; i < j; i, j = i+1, j-1 {
			records[i], records[j] = records[j], records[i]
		}
		return BoxListInvoicesRecords{
			Invoices:    records,
			TotalAmount: totalAmount,
		}, nil
	}

	if assetId == "00" {
		for _, invoice := range resp.Invoices {
			if invoice.State == lnrpc.Invoice_SETTLED {
				if len(invoice.Htlcs) == 0 || invoice.Htlcs[0] == nil {
					records = append(records, invoice)
					totalAmount += invoice.Value
					continue
				}
				if invoice.Htlcs[0].CustomChannelData == nil {
					records = append(records, invoice)
					totalAmount += invoice.Value
					continue
				}
				var result rfqmsg.JsonHtlc
				err := json.Unmarshal(invoice.Htlcs[0].CustomChannelData, &result)
				if err != nil {
					continue
				}
				if len(result.Balances) == 0 || result.Balances[0].AssetID == "" {
					records = append(records, invoice)
					totalAmount += invoice.Value
					continue
				}
				continue
			}
		}
		for i, j := 0, len(records)-1; i < j; i, j = i+1, j-1 {
			records[i], records[j] = records[j], records[i]
		}
		return BoxListInvoicesRecords{
			Invoices:    records,
			TotalAmount: totalAmount,
		}, nil
	}

	for _, invoice := range resp.Invoices {
		if len(invoice.Htlcs) == 0 {
			continue
		}
		if invoice.Htlcs[0].CustomChannelData != nil && invoice.State == lnrpc.Invoice_SETTLED {
			var result rfqmsg.JsonHtlc
			err := json.Unmarshal(invoice.Htlcs[0].CustomChannelData, &result)
			if err != nil {
				continue
			}
			if len(result.Balances) > 0 && result.Balances[0].AssetID == assetId {
				records = append(records, invoice)
				if len(result.Balances) > 0 {
					totalAmount += int64(result.Balances[0].Amount)
				}
			}
		}
	}
	for i, j := 0, len(records)-1; i < j; i, j = i+1, j-1 {
		records[i], records[j] = records[j], records[i]
	}
	return BoxListInvoicesRecords{
		Invoices:    records,
		TotalAmount: totalAmount,
	}, nil
}

type BoxListPaymentsRecords struct {
	Payments    []*lnrpc.Payment
	TotalAmount int64
}

func (l Ln) BoxListPaymentsRecords(assetId string, MaxPayments uint64, CreationDateStart uint64, CreationDateEnd uint64) (BoxListPaymentsRecords, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return BoxListPaymentsRecords{}, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	lc := lnrpc.NewLightningClient(conn)
	req := &lnrpc.ListPaymentsRequest{
		MaxPayments:       MaxPayments,
		CreationDateStart: CreationDateStart,
		CreationDateEnd:   CreationDateEnd,
	}

	resp, err := lc.ListPayments(ctx, req)
	if err != nil {
		return BoxListPaymentsRecords{}, errors.Wrap(err, "lc.ListPayments")
	}

	var records []*lnrpc.Payment
	var totalAmount int64

	if assetId == "" {
		for _, payment := range resp.Payments {
			if payment.Status == lnrpc.Payment_SUCCEEDED {
				records = append(records, payment)
				totalAmount += payment.ValueSat
			}
		}
		for i, j := 0, len(records)-1; i < j; i, j = i+1, j-1 {
			records[i], records[j] = records[j], records[i]
		}
		return BoxListPaymentsRecords{
			Payments:    records,
			TotalAmount: totalAmount,
		}, nil
	}

	if assetId == "00" {
		for _, payment := range resp.Payments {
			if len(payment.Htlcs) > 0 && payment.Htlcs[0].Route.CustomChannelData == nil && payment.Status == lnrpc.Payment_SUCCEEDED {
				records = append(records, payment)
				totalAmount += payment.ValueSat
			}
		}
		for i, j := 0, len(records)-1; i < j; i, j = i+1, j-1 {
			records[i], records[j] = records[j], records[i]
		}
		return BoxListPaymentsRecords{
			Payments:    records,
			TotalAmount: totalAmount,
		}, nil
	}

	for _, payment := range resp.Payments {
		if len(payment.Htlcs) == 0 || payment.Htlcs[0].Route.CustomChannelData == nil {
			continue
		}
		var result rfqmsg.JsonHtlc
		err := json.Unmarshal(payment.Htlcs[0].Route.CustomChannelData, &result)
		if err != nil {
			continue
		}
		if len(result.Balances) > 0 && result.Balances[0].AssetID == assetId && payment.Status == lnrpc.Payment_SUCCEEDED {
			records = append(records, payment)
			if len(result.Balances) > 0 {
				totalAmount += int64(result.Balances[0].Amount)
			}
		}
	}
	for i, j := 0, len(records)-1; i < j; i, j = i+1, j-1 {
		records[i], records[j] = records[j], records[i]
	}
	return BoxListPaymentsRecords{
		Payments:    records,
		TotalAmount: totalAmount,
	}, nil
}

func (l Ln) GetTransactions() (*lnrpc.TransactionDetails, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	lc := lnrpc.NewLightningClient(conn)
	req := &lnrpc.GetTransactionsRequest{}

	resp, err := lc.GetTransactions(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "lc.GetTransactions")
	}
	return resp, nil
}

func (l Ln) ListAddresses() (*walletrpc.ListAddressesResponse, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	wkc := walletrpc.NewWalletKitClient(conn)
	req := &walletrpc.ListAddressesRequest{}

	resp, err := wkc.ListAddresses(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "wkc.ListAddresses")
	}
	return resp, nil
}

func GetAllAddresses() ([]string, error) {
	var result []string
	var l Ln
	listAddress, err := l.ListAddresses()
	if err != nil {
		return nil, err
	}
	for _, accountWithAddresse := range listAddress.AccountWithAddresses {
		addresses := accountWithAddresse.Addresses
		for _, address := range addresses {
			result = append(result, address.Address)
		}
	}
	return result, nil
}

func (l Ln) GetClosedChannels() (*lnrpc.ClosedChannelsResponse, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	lc := lnrpc.NewLightningClient(conn)
	req := &lnrpc.ClosedChannelsRequest{}

	resp, err := lc.ClosedChannels(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "lc.ClosedChannels")
	}
	return resp, nil
}

func (l Ln) SendCoins(Addr string, Amount int64, SatPerVbyte uint64, SendAll bool) (*lnrpc.SendCoinsResponse, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	lc := lnrpc.NewLightningClient(conn)
	req := &lnrpc.SendCoinsRequest{
		Addr:        Addr,
		Amount:      Amount,
		SatPerVbyte: SatPerVbyte,
		SendAll:     SendAll,
	}

	resp, err := lc.SendCoins(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "lc.SendCoins")
	}
	return resp, nil
}

func (l Ln) NewAddress() (*lnrpc.NewAddressResponse, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	lc := lnrpc.NewLightningClient(conn)
	req := &lnrpc.NewAddressRequest{
		Type: lnrpc.AddressType_TAPROOT_PUBKEY,
	}

	resp, err := lc.NewAddress(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "lc.NewAddress")
	}
	return resp, nil
}

func (l Ln) ListChaintxns() (*lnrpc.TransactionDetails, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	lc := lnrpc.NewLightningClient(conn)

	resp, err := lc.GetTransactions(ctx, &lnrpc.GetTransactionsRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "lc.GetTransactions")
	}
	return resp, nil
}

func (l Ln) PendingChannels() (*lnrpc.PendingChannelsResponse, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	lc := lnrpc.NewLightningClient(conn)
	req := &lnrpc.PendingChannelsRequest{}

	resp, err := lc.PendingChannels(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "lc.PendingChannels")
	}
	return resp, nil
}

func (l Ln) SubscribeInvoices(ctx context.Context) {
	request := &lnrpc.InvoiceSubscription{}
	baseDelay := 2 * time.Second
	maxDelay := 30 * time.Second
	curDelay := baseDelay
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		conn, err := GetConn(ln, false)
		if err != nil {
			loggers.Chan().Println(errors.Wrap(err, "SubscribeInvoices GetConn"))
			time.Sleep(curDelay)
			if curDelay < maxDelay {
				curDelay *= 2
				if curDelay > maxDelay {
					curDelay = maxDelay
				}
			}
			continue
		}

		loggers.Chan().Println("SubscribeInvoices GetConn success")

		client := lnrpc.NewLightningClient(conn)
		stream, err := client.SubscribeInvoices(ctx, request)
		if err != nil {
			loggers.Chan().Println(errors.Wrap(err, "SubscribeInvoices SubscribeInvoices"))
			Close(conn)
			time.Sleep(curDelay)
			if curDelay < maxDelay {
				curDelay *= 2
				if curDelay > maxDelay {
					curDelay = maxDelay
				}
			}
			continue
		}
		loggers.Chan().Println("SubscribeInvoices Prostream")

		for {
			invoice, err := stream.Recv()
			if err != nil {
				loggers.Chan().Println("SubscribeInvoices stream.Recv error")
				Close(conn)
				loggers.Chan().Println(errors.Wrap(err, "SubscribeInvoices stream.Recv"))
				time.Sleep(curDelay)
				if curDelay < maxDelay {
					curDelay *= 2
					if curDelay > maxDelay {
						curDelay = maxDelay
					}
				}
				break
			}
			if invoice == nil {
				loggers.Chan().Println("SubscribeInvoices invoice == nil")
				continue
			}
			if invoice.State == lnrpc.Invoice_SETTLED {
				loggers.Chan().Println("SubscribeInvoices invoice.State == lnrpc.Invoice_SETTLED")
				curDelay = baseDelay
				resp, err := listchannelsForTradeChanInfoData()
				if err != nil {
					loggers.Chan().Println(errors.Wrap(err, "listchannelsForTradeChanInfoData"))
					continue
				}
				loggers.Chan().Println("listchannelsForTradeChanInfoData success")
				if err := SubToUploadTradeChanInfo(resp); err != nil {
					loggers.Chan().Println(errors.Wrap(err, "SubToUploadTradeChanInfo"))
					continue
				}
				loggers.Chan().Println("SubToUploadTradeChanInfo success")
			}
		}
	}
}

type ListchannelsForTradeChanInfoReq struct {
	IdentityPubkey         string `json:"identity_pubkey"`
	BtcChannelPoint        string `json:"btc_channel_point"`
	AssetChannelPoint      string `json:"asset_channel_point"`
	LocalSatsBalance       int64  `json:"local_sats_balance"`
	RemoteSatsBalance      int64  `json:"remote_sats_balance"`
	AssetSatsBalance       int64  `json:"asset_sats_balance"`
	AssetRemoteSatsBalance int64  `json:"asset_remote_sats_balance"`
	AssetLocalBalance      int64  `json:"asset_local_balance"`
	AssetRemoteBalance     int64  `json:"asset_remote_balance"`
}

func listchannelsForTradeChanInfoData() (*ListchannelsForTradeChanInfoReq, error) {
	var l Ln

	info, err := l.GetInfo()
	if err != nil {
		return nil, errors.Wrap(err, "l.GetInfo")
	}

	chans, err := l.ListChannels(false, false)
	if err != nil {
		return nil, errors.Wrap(err, "l.ListChannels")
	}

	var btcChan *lnrpc.Channel
	var assetChan *lnrpc.Channel

	for _, c := range chans.Channels {
		if len(c.CustomChannelData) == 0 {
			if btcChan == nil {
				btcChan = c
			}
		} else {
			if assetChan == nil {
				assetChan = c
			}
		}
		if btcChan != nil && assetChan != nil {
			break
		}
	}

	var result *ListchannelsForTradeChanInfoReq = &ListchannelsForTradeChanInfoReq{
		IdentityPubkey:         info.IdentityPubkey,
		BtcChannelPoint:        "",
		AssetChannelPoint:      "",
		LocalSatsBalance:       0,
		RemoteSatsBalance:      0,
		AssetSatsBalance:       0,
		AssetRemoteSatsBalance: 0,
		AssetLocalBalance:      0,
		AssetRemoteBalance:     0,
	}

	if btcChan != nil {
		result.BtcChannelPoint = btcChan.ChannelPoint
		result.LocalSatsBalance = btcChan.LocalBalance
		result.RemoteSatsBalance = btcChan.RemoteBalance
	}

	if assetChan != nil {
		result.AssetChannelPoint = assetChan.ChannelPoint
		result.AssetSatsBalance = assetChan.LocalBalance
		result.AssetRemoteSatsBalance = assetChan.RemoteBalance

		var ac rfqmsg.JsonAssetChannel
		if err := json.Unmarshal(assetChan.CustomChannelData, &ac); err == nil {
			result.AssetLocalBalance = int64(ac.LocalBalance)
			result.AssetRemoteBalance = int64(ac.RemoteBalance)
		}
	}

	return result, nil
}

func SubToUploadTradeChanInfo(req *ListchannelsForTradeChanInfoReq) error {
	loggers.Chan().Println("SubToUploadTradeChanInfo")
	host := sc.BaseUrl
	targetUrl := fmt.Sprintf("%s/box_devices/update_box_channel_info", host)

	client := resty.New()

	var r models.JResult2
	var e models.ErrResp

	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(sc.BoxIpBasicUser, sc.BoxIpBasicPass).
		SetBody(req).
		SetResult(&r).
		SetError(&e).
		Post(targetUrl)

	if err != nil {
		loggers.Chan().Println(errors.Wrap(err, "client.R.Post"))
		return errors.Wrap(err, "client.R.Post")
	}

	if e.Error != "" {
		return errors.New(fmt.Sprintf("error: %s", e.Error))
	}

	if r.ErrMsg != "" {
		return errors.New(fmt.Sprintf("error: %s", r.ErrMsg))
	}

	return nil
}

func (l Ln) ListChannelsAll() (*lnrpc.ListChannelsResponse, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	lc := lnrpc.NewLightningClient(conn)
	req := &lnrpc.ListChannelsRequest{}

	resp, err := lc.ListChannels(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "lc.ListChannels")
	}
	return resp, nil
}

func (l Ln) ListAssetPaymentsAmount(assetId string) (int64, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return 0, err
	}

	defer Close(conn)
	totalAmount := int64(0)

	ctx := context.Background()
	lc := lnrpc.NewLightningClient(conn)
	req := &lnrpc.ListPaymentsRequest{}

	resp, err := lc.ListPayments(ctx, req)
	if err != nil {
		return 0, errors.Wrap(err, "lc.ListPayments")
	}
	for _, payment := range resp.Payments {
		if payment.Status != lnrpc.Payment_SUCCEEDED {
			continue
		}
		if len(payment.Htlcs) == 0 {
			continue
		}
		for _, h := range payment.Htlcs {
			if h.Route == nil || h.Route.CustomChannelData == nil {
				continue
			}
			var result rfqmsg.JsonHtlc
			if err := json.Unmarshal(h.Route.CustomChannelData, &result); err != nil {
				continue
			}
			if len(result.Balances) > 0 && result.Balances[0].AssetID == assetId {
				totalAmount += int64(result.Balances[0].Amount)
			}
		}
	}

	return totalAmount, nil
}

func (l Ln) BoxPendingChannelsAmount(assetId string) (int64, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return 0, err
	}

	defer Close(conn)
	totalAmount := int64(0)

	ctx := context.Background()
	lc := lnrpc.NewLightningClient(conn)
	req := &lnrpc.PendingChannelsRequest{}

	resp, err := lc.PendingChannels(ctx, req)
	if err != nil {
		return 0, errors.Wrap(err, "lc.PendingChannels")
	}
	for _, channel := range resp.PendingForceClosingChannels {
		if len(channel.Channel.CustomChannelData) == 0 {
			continue
		}
		if channel.Channel.CustomChannelData != nil {
			var result rfqmsg.JsonAssetChannel
			err := json.Unmarshal(channel.Channel.CustomChannelData, &result)
			if err != nil {
				continue
			}
			if result.FundingAssets[0].AssetGenesis.AssetID == assetId {
				totalAmount += int64(result.LocalBalance)
			}
		}
	}
	for _, channel := range resp.WaitingCloseChannels {
		if len(channel.Channel.CustomChannelData) == 0 {
			continue
		}
		if channel.Channel.CustomChannelData != nil {
			var result rfqmsg.JsonAssetChannel
			err := json.Unmarshal(channel.Channel.CustomChannelData, &result)
			if err != nil {
				continue
			}
			if result.FundingAssets[0].AssetGenesis.AssetID == assetId {
				totalAmount += int64(result.LocalBalance)
			}
		}
	}
	return totalAmount, nil
}

func (l Ln) BoxBtcAddInvoices(value int64, memo string) (string, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return "", errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	lc := lnrpc.NewLightningClient(conn)
	req := &lnrpc.Invoice{
		Memo:    memo,
		Value:   value,
		Private: true,
	}

	resp, err := lc.AddInvoice(ctx, req)
	if err != nil {
		return "", errors.Wrap(err, "lc.AddInvoice")
	}
	return resp.PaymentRequest, nil
}

func (l Ln) BoxBtcPayInvoice(invoice string, amt int, feelimit int, outgoingChanId string, allowSelfPayment bool) (*lnrpc.Payment, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return nil, err
	}

	defer Close(conn)

	ctx := context.Background()
	ctxt, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	lc := routerrpc.NewRouterClient(conn)
	req := &routerrpc.SendPaymentRequest{
		PaymentRequest:   invoice,
		Amt:              int64(amt),
		FeeLimitSat:      int64(feelimit),
		AllowSelfPayment: allowSelfPayment,
	}

	req.OutgoingChanIds, err = ParseChanIDs([]string{outgoingChanId})
	if err != nil {
		return nil, err
	}
	stream, err := lc.SendPaymentV2(ctxt, req)
	if err != nil {
		return nil, err
	}
	for {
		response, err := stream.Recv()
		if err != nil {
			return nil, err
		} else if response != nil {
			if response.Status == 2 || response.Status == 3 {
				return response, nil
			}
		}
	}
}

func ParseChanIDs(idStrings []string) ([]uint64, error) {
	if len(idStrings) == 0 {
		return nil, nil
	}

	chanIDs := make([]uint64, len(idStrings))
	for i, idStr := range idStrings {
		scid, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return nil, err
		}

		chanIDs[i] = scid
	}

	return chanIDs, nil
}

type BoxListChannelsResp struct {
	ChannelsInfo *lnrpc.ListChannelsResponse
	BtcChanId    string
	AssetChanId  string
}

func (l Ln) BoxListChannels() (*BoxListChannelsResp, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	lc := lnrpc.NewLightningClient(conn)
	req := &lnrpc.ListChannelsRequest{}

	chanInfo, err := lc.ListChannels(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "lc.ListChannels")
	}
	resp := &BoxListChannelsResp{ChannelsInfo: chanInfo}
	for _, c := range chanInfo.Channels {
		if len(c.CustomChannelData) == 0 {
			if resp.BtcChanId == "" {
				resp.BtcChanId = strconv.FormatUint(c.ChanId, 10)
			}
		} else {
			if resp.AssetChanId == "" {
				resp.AssetChanId = strconv.FormatUint(c.ChanId, 10)
			}
		}
		if resp.BtcChanId != "" && resp.AssetChanId != "" {
			break
		}
	}
	return resp, nil
}

func (l Ln) BoxBtcDecodePayReq(invoice string) (*lnrpc.PayReq, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)
	ctx := context.Background()
	lc := lnrpc.NewLightningClient(conn)
	req := &lnrpc.PayReqString{
		PayReq: invoice,
	}
	resp, err := lc.DecodePayReq(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "lc.DecodePayReq")
	}
	return resp, nil
}

func (l Ln) BoxGetInfo() (*lnrpc.GetInfoResponse, error) {
	conn, err := GetConn(ln, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)
	ctx := context.Background()
	lc := lnrpc.NewLightningClient(conn)
	req := &lnrpc.GetInfoRequest{}
	resp, err := lc.GetInfo(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "lc.GetInfo")
	}
	return resp, nil
}
