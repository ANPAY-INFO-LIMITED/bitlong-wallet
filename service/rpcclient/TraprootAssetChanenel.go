package rpcclient

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/lightninglabs/taproot-assets/taprpc/priceoraclerpc"
	tchrpc "github.com/lightninglabs/taproot-assets/taprpc/tapchannelrpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	"github.com/wallet/service/apiConnect"
)

func getTaprootChannelClient() (tchrpc.TaprootAssetChannelsClient, func(), error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	client := tchrpc.NewTaprootAssetChannelsClient(conn)
	return client, clearUp, nil
}

func AddAssetInvoice(amount uint64, assetId string, rfqPeer string, memo string) (*tchrpc.AddInvoiceResponse, error) {
	client, clearUp, err := getTaprootChannelClient()
	if err != nil {
		return nil, err
	}
	defer clearUp()
	assetIDBytes, err := hex.DecodeString(assetId)
	if err != nil {
		return nil, fmt.Errorf("error hex decoding asset "+
			"ID: %w", err)
	}
	rfqPeerKey, err := hex.DecodeString(rfqPeer)
	if err != nil {
		return nil, fmt.Errorf("unable to decode RFQ peer public key: "+
			"%w", err)
	}
	req := &tchrpc.AddInvoiceRequest{
		AssetId:     assetIDBytes,
		AssetAmount: amount,
		PeerPubkey:  rfqPeerKey,
		InvoiceRequest: &lnrpc.Invoice{
			Memo: memo,
		},
	}
	resp, err := client.AddInvoice(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func PayAssetInvoice(invoice string, assetId string, rfqPeerPubkey string, outgoingChanId uint64, feeLimitSat int64, allowSelfPayment bool) (*lnrpc.Payment, error) {
	client, clearUp, err := getTaprootChannelClient()
	if err != nil {
		return nil, err
	}
	defer clearUp()
	assetIdByte, err := hex.DecodeString(assetId)
	if err != nil {
		return nil, err
	}
	peerPubkey, err := hex.DecodeString(rfqPeerPubkey)
	if err != nil {
		return nil, err
	}
	req := &tchrpc.SendPaymentRequest{
		AssetId:    assetIdByte,
		PeerPubkey: peerPubkey,
		PaymentRequest: &routerrpc.SendPaymentRequest{
			PaymentRequest:   invoice,
			TimeoutSeconds:   60,
			AllowSelfPayment: allowSelfPayment,
		},
	}
	if outgoingChanId != 0 {
		req.PaymentRequest.OutgoingChanIds = []uint64{outgoingChanId}
	}
	if feeLimitSat != 0 {
		req.PaymentRequest.FeeLimitSat = feeLimitSat
	} else {
		req.PaymentRequest.FeeLimitSat = 10000
	}
	stream, err := client.SendPayment(context.Background(), req)
	if err != nil {
		return nil, err
	}
	for {
		resp, err := stream.Recv()
		if err != nil {
			return nil, err
		} else if resp != nil {
			res := resp.Result
			switch r := res.(type) {
			case *tchrpc.SendPaymentResponse_AcceptedSellOrder:
			case *tchrpc.SendPaymentResponse_PaymentResult:
				if resp == nil || resp.Result == nil ||
					resp.GetPaymentResult() == nil {

					return nil, fmt.Errorf("unexpected response: %v", resp)
				}
				result := resp.GetPaymentResult()
				if result.Status == 2 || result.Status == 3 {
					return result, nil
				}
			default:
				return nil, fmt.Errorf("unexpected response type: %T", r)
			}
		}
	}
}
func QueryAssetRates(transType priceoraclerpc.TransactionType, assetIDStr string, maxAssetAmount uint64,
	maxsatamount uint64) (*priceoraclerpc.AssetRates, error) {
	assetIDBytes, err := hex.DecodeString(assetIDStr)
	if err != nil {
		return nil, fmt.Errorf("unable to decode assetID: %v", err)
	}
	var paymentAssetId = make([]byte, 32)
	request := &priceoraclerpc.QueryAssetRatesRequest{
		TransactionType: transType,
		SubjectAsset: &priceoraclerpc.AssetSpecifier{
			Id: &priceoraclerpc.AssetSpecifier_AssetId{
				AssetId: assetIDBytes,
			},
		},
		SubjectAssetMaxAmount: maxAssetAmount,
		PaymentAsset: &priceoraclerpc.AssetSpecifier{
			Id: &priceoraclerpc.AssetSpecifier_AssetId{
				AssetId: paymentAssetId,
			},
		},
		PaymentAssetMaxAmount: maxsatamount,
	}

	conn, closeConn, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return nil, err
	}
	defer closeConn()
	client := priceoraclerpc.NewPriceOracleClient(conn)
	resp, err := client.QueryAssetRates(context.Background(), request)
	if err != nil {
		return nil, err
	}

	switch result := resp.GetResult().(type) {
	case *priceoraclerpc.QueryAssetRatesResponse_Ok:
		if result.Ok.AssetRates == nil {
			return nil, fmt.Errorf("QueryAssetRates response is " +
				"successful but asset rates is nil")
		}
		return result.Ok.AssetRates, nil
	case *priceoraclerpc.QueryAssetRatesResponse_Error:
		if result.Error == nil {
			return nil, fmt.Errorf("QueryAssetRates response is " +
				"an error but error is nil")
		}
		return nil, fmt.Errorf("%s", result.Error.Message)
	default:
		return nil, fmt.Errorf("unexpected response type: %T", result)
	}
}
