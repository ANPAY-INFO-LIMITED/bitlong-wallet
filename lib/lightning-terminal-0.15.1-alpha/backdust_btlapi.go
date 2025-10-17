//go:build btlapi
// +build btlapi

package terminal

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/lightninglabs/taproot-assets/rfq"
	"github.com/lightninglabs/taproot-assets/rfqmsg"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	"github.com/lightningnetwork/lnd/lntypes"
	"github.com/lightningnetwork/lnd/record"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"math/rand"
	"os"
	"time"
)

type litHandler struct {
	LitServer *LightningTerminal
	IsStarted bool
	LndCoon   *grpc.ClientConn
}

var defaultHandler litHandler

func setLitHandler(handler *LightningTerminal) {
	defaultHandler.LitServer = handler
	defaultHandler.IsStarted = true
	if err := defaultHandler.getClientConn(); err != nil {
		fmt.Println("getClientConn Error:", err)
		return
	}
	go defaultHandler.daemon()
}
func (l *litHandler) daemon() {
	if l.LndCoon == nil {
		return
	}
	defer l.LndCoon.Close()

	getStatus := func() bool {
		status, err := l.LitServer.statusMgr.SubServerStatus(nil, nil)
		if err != nil {
			return false
		}
		if status.SubServers["lit"].Running && status.SubServers["lnd"].Running {
			return true
		}
		return false
	}
	for {
		if getStatus() {
			l.SubscribeHtlcEvents(l.ReBackDust)
		}
		time.Sleep(time.Second * 2)
	}
}

func (l *litHandler) ReBackDust(e *routerrpc.HtlcEvent) {
	switch e.EventType {
	case routerrpc.HtlcEvent_RECEIVE:
		if settle, ok := e.Event.(*routerrpc.HtlcEvent_SettleEvent); ok {
			if l.checkInvoice(settle, e.IncomingChannelId) {
				l.sendDust(e.IncomingChannelId)
			}
		}
	case routerrpc.HtlcEvent_UNKNOWN:
		if final, ok := e.Event.(*routerrpc.HtlcEvent_FinalHtlcEvent); ok {
			if rfq.CheckIncomingHtlcId(e.IncomingHtlcId) {
				if final.FinalHtlcEvent.Settled && final.FinalHtlcEvent.Offchain {
					l.sendDust(e.IncomingChannelId)
				}
			}
		}
	default:
	}
}

func (l *litHandler) checkInvoice(settle *routerrpc.HtlcEvent_SettleEvent, channelId uint64) bool {
	if settle.SettleEvent.Preimage != nil {
		invoices, err := l.ListInvoices(20, true)
		if err != nil {
			fmt.Println("ReBackDust Error:", err)
			return false
		}
		for _, i := range invoices.Invoices {
			p := hex.EncodeToString(i.RPreimage)
			if p == hex.EncodeToString(settle.SettleEvent.Preimage) {
				if l.checkHtlcIsAsset(i.Htlcs, channelId) {
					return true
				}
				return false
			}
		}
	}
	return false
}
func (l *litHandler) checkHtlcIsAsset(htlcs []*lnrpc.InvoiceHTLC, channelId uint64) bool {
	for _, htlc := range htlcs {
		if htlc.ChanId == channelId {
			if htlc.CustomChannelData != nil {
				data := rfqmsg.JsonHtlc{}
				err := json.Unmarshal(htlc.CustomChannelData, &data)
				if err != nil {
					fmt.Println("ReBackDust Error:", err)
					return false
				}
				if data.Balances != nil && len(data.Balances) > 0 {
					return true
				}
				return false
			}
		}
	}
	return false
}
func (l *litHandler) sendDust(ChanId uint64) {
	var dust string
	response, err := l.ListChannels(true, true)
	if err != nil {
		return
	}
	if response == nil || len(response.Channels) == 0 {
		return
	}
	for _, channel := range response.Channels {
		if channel.ChanId == ChanId {
			if channel.Initiator {
				if channel.LocalBalance < 5540 {
					return
				}
			} else {
				if channel.LocalBalance < channel.LocalChanReserveSat+400 {
					return
				}
			}
			dust = channel.RemotePubkey
		}
	}
	if dust == "" {
		return
	}
	result, err := l.SendPaymentV2ByKeySend(dust, 354, 10, ChanId)
	if err != nil {
		fmt.Println("backDust Error:", err)
		return
	}
	fmt.Println("backDust Result:", result)
}

func (l *litHandler) ListInvoices(maxNum int, reversed bool) (*lnrpc.ListInvoiceResponse, error) {
	client := lnrpc.NewLightningClient(l.LndCoon)
	request := &lnrpc.ListInvoiceRequest{
		Reversed: reversed,
	}
	if maxNum > 0 {
		request.NumMaxInvoices = uint64(maxNum)
	}
	response, err := client.ListInvoices(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return response, err
}
func (l *litHandler) ListChannels(activeOnly bool, private bool) (*lnrpc.ListChannelsResponse, error) {
	client := lnrpc.NewLightningClient(l.LndCoon)
	request := &lnrpc.ListChannelsRequest{
		ActiveOnly:  activeOnly,
		PrivateOnly: private,
	}
	response, err := client.ListChannels(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (l *litHandler) SendPaymentV2ByKeySend(dest string, amt int, feelimit int, outgoingChanId uint64) (*lnrpc.Payment, error) {
	client := routerrpc.NewRouterClient(l.LndCoon)
	destNode, _ := hex.DecodeString(dest)
	if len(destNode) != 33 {
		return nil, fmt.Errorf("dest node pubkey must be exactly 33 bytes, is "+
			"instead: %v", len(dest))
	}
	req := &routerrpc.SendPaymentRequest{
		TimeoutSeconds:    30,
		Dest:              destNode,
		Amt:               int64(amt),
		FeeLimitSat:       int64(feelimit),
		DestCustomRecords: make(map[uint64][]byte),
		OutgoingChanIds:   make([]uint64, 1),
	}
	req.OutgoingChanIds[0] = outgoingChanId

	var rHash []byte
	var preimage lntypes.Preimage
	if _, err := rand.Read(preimage[:]); err != nil {
		return nil, err
	}
	// Set the preimage. If the user supplied a preimage with the
	// data flag, the preimage that is set here will be overwritten
	// later.
	req.DestCustomRecords[record.KeySendType] = preimage[:]
	hash := preimage.Hash()
	rHash = hash[:]
	req.PaymentHash = rHash

	ctx := context.Background()
	ctxt, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	stream, err := client.SendPaymentV2(ctxt, req)
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
func (l *litHandler) SubscribeHtlcEvents(f func(e *routerrpc.HtlcEvent)) {
	client := routerrpc.NewRouterClient(l.LndCoon)
	if f == nil {
		f = func(e *routerrpc.HtlcEvent) {
			fmt.Println(e)
		}
	}
	request := &routerrpc.SubscribeHtlcEventsRequest{}
	stream, err := client.SubscribeHtlcEvents(context.Background(), request)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		event, err := stream.Recv()
		if err != nil {
			fmt.Println(err)
			return
		}
		if event != nil {
			f(event)
		}
	}
}

func (l *litHandler) getClientConn() error {
	var err error
	host, _, _, _, macData := l.LitServer.cfg.lndConnectParams()
	// get tls cert and macaroon
	//creds, err := NewTlsCert(tlsPath)
	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true, // nolint:gosec
	})
	//if err != nil {
	//	return err
	//}
	l.LndCoon, err = grpc.Dial(host, grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(NewMacaroonCredential(macData)), grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(10*1024*1024), // 10 MB
			grpc.MaxCallSendMsgSize(10*1024*1024), // 10 MB
		))
	if err != nil {
		return err
	}
	return nil
}

func NewTlsCert(tlsCertPath string) (credentials.TransportCredentials, error) {
	cert, err := os.ReadFile(tlsCertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read tls cert: %v", err)
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(cert) {
		return nil, fmt.Errorf("failed to append cert: %v", err)
	}
	config := &tls.Config{
		MinVersion: tls.VersionTLS12,
		RootCAs:    certPool,
	}
	creds := credentials.NewTLS(config)
	return creds, nil
}

type MacaroonCredential struct {
	macaroon string
}

func NewMacaroonCredential(macaroonBytes []byte) *MacaroonCredential {
	macaroon := hex.EncodeToString(macaroonBytes)
	return &MacaroonCredential{macaroon: macaroon}
}

func (c *MacaroonCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"macaroon": c.macaroon}, nil
}

func (c *MacaroonCredential) RequireTransportSecurity() bool {
	return true
}
