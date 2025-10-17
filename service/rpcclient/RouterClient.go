package rpcclient

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	"github.com/lightningnetwork/lnd/lntypes"
	"github.com/lightningnetwork/lnd/record"
	"github.com/pkg/errors"
	"github.com/wallet/service/apiConnect"
)

func getRouterClient() (routerrpc.RouterClient, func(), error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	client := routerrpc.NewRouterClient(conn)
	return client, clearUp, nil
}

func SendPaymentV2(invoice string, amt int, feelimit int, outgoingChanId uint64, allowSelfPayment bool) (*lnrpc.Payment, error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()

	client := routerrpc.NewRouterClient(conn)
	request := &routerrpc.SendPaymentRequest{
		TimeoutSeconds:   30,
		AllowSelfPayment: allowSelfPayment,
	}
	if invoice != "" {
		request.PaymentRequest = invoice
	}
	if outgoingChanId != 0 {
		request.OutgoingChanIds = []uint64{outgoingChanId}
	}
	if amt != 0 {
		request.Amt = int64(amt)
	}
	if feelimit == 0 {
		request.FeeLimitSat = 1000
	} else {
		request.FeeLimitSat = int64(feelimit)
	}
	ctx := context.Background()
	ctxt, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	stream, err := client.SendPaymentV2(ctxt, request)
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

func SendPaymentV2ByKeySend(dest string, amt int, feelimit int, outgoingChanId uint64) (*lnrpc.Payment, error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()
	client := routerrpc.NewRouterClient(conn)

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

func SubscribeHtlcEvents(f func(e *routerrpc.HtlcEvent)) {
	client, clearUp, err := getRouterClient()
	if err != nil {
		fmt.Println(err)
	}
	defer clearUp()
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
