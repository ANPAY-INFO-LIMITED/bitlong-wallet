package rpc

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/lightninglabs/taproot-assets/rfqmath"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/lightninglabs/taproot-assets/taprpc/btlchannelrpc"
	"github.com/lightninglabs/taproot-assets/taprpc/tapchannelrpc"
	"github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	"github.com/lightningnetwork/lnd/lntypes"
	"github.com/lightningnetwork/lnd/record"
	"github.com/pkg/errors"
	"github.com/wallet/api"
	"github.com/wallet/box/sc"
)

type Tap struct{}

func (t Tap) GetInfo() (*taprpc.GetInfoResponse, error) {
	conn, err := GetConn(tap, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	tac := taprpc.NewTaprootAssetsClient(conn)
	req := &taprpc.GetInfoRequest{}

	resp, err := tac.GetInfo(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "tac.GetInfo")
	}
	return resp, nil
}

func (t Tap) SyncUniverse(universeHost string, assetID string) (*universerpc.SyncResponse, error) {
	conn, err := GetConn(tap, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	uc := universerpc.NewUniverseClient(conn)

	var syncTargets []*universerpc.SyncTarget
	universeID := &universerpc.ID{
		Id: &universerpc.ID_AssetIdStr{
			AssetIdStr: assetID,
		},
		ProofType: universerpc.ProofType_PROOF_TYPE_ISSUANCE,
	}
	syncTargets = append(syncTargets, &universerpc.SyncTarget{
		Id: universeID,
	})

	req := &universerpc.SyncRequest{
		UniverseHost: universeHost,
		SyncMode:     universerpc.UniverseSyncMode_SYNC_FULL,
		SyncTargets:  syncTargets,
	}

	resp, err := uc.SyncUniverse(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "uc.SyncUniverse")
	}
	return resp, nil
}

func (t Tap) ListUtxos() (*taprpc.ListUtxosResponse, error) {
	conn, err := GetConn(tap, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	tac := taprpc.NewTaprootAssetsClient(conn)
	req := &taprpc.ListUtxosRequest{}

	resp, err := tac.ListUtxos(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "tac.ListUtxos")
	}
	return resp, nil
}

func GetListUtxos(token string) (*[]api.ManagedUtxo, error) {
	var t Tap
	resp, err := t.ListUtxos()
	if err != nil {
		return nil, err
	}
	managedUtxos := api.ListUtxosResponseToManagedUtxos(resp)
	err = api.RemoveNotLocalAssetManagedUtxos2(sc.BaseUrl, token, managedUtxos)
	if err != nil {
		return nil, errors.Wrap(err, "api.RemoveNotLocalAssetManagedUtxos")
	}
	managedUtxos, err = api.GetTimeForManagedUtxoByBitcoind2(sc.BaseUrl, token, managedUtxos)
	if err != nil {
		return nil, err
	}
	return managedUtxos, nil
}

func (t Tap) NewAddr(assetId []byte, amt uint64) (*taprpc.Addr, error) {
	conn, err := GetConn(tap, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	tac := taprpc.NewTaprootAssetsClient(conn)
	req := &taprpc.NewAddrRequest{
		AssetId:          assetId,
		Amt:              amt,
		ProofCourierAddr: proofCourierAddr,
	}

	resp, err := tac.NewAddr(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "tac.NewAddr")
	}
	return resp, nil
}

func (t Tap) SendAsset(tapAddrs []string, feeRate uint32) (*taprpc.SendAssetResponse, error) {
	conn, err := GetConn(tap, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	tac := taprpc.NewTaprootAssetsClient(conn)
	req := &taprpc.SendAssetRequest{
		TapAddrs: tapAddrs,
		FeeRate:  feeRate,
	}

	resp, err := tac.SendAsset(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "tac.SendAsset")
	}
	return resp, nil
}

func (t Tap) DecodeAddr(addr string) (*taprpc.Addr, error) {
	conn, err := GetConn(tap, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	tac := taprpc.NewTaprootAssetsClient(conn)
	req := &taprpc.DecodeAddrRequest{Addr: addr}

	resp, err := tac.DecodeAddr(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "tac.DecodeAddr")
	}
	return resp, nil
}

func (t Tap) FundBtlAssetChannel(assetAmount uint64, assetId string, peerPubkey string, feeRate uint32, pushSat int64, localAmt uint64) (*btlchannelrpc.FundBtlChannelResponse, error) {
	conn, err := GetConn(tap, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	assetIdBytes, err := hex.DecodeString(assetId)
	if err != nil {
		return nil, errors.Wrap(err, "hex.DecodeString")
	}

	pubkeyBytes, err := hex.DecodeString(peerPubkey)
	if err != nil {
		return nil, errors.Wrap(err, "hex.DecodeString")
	}

	ctx := context.Background()
	bc := btlchannelrpc.NewBtlChannelsClient(conn)
	resp, err := bc.FundBtlChannel(ctx, &btlchannelrpc.FundBtlChannelRequest{
		AssetAmount:        assetAmount,
		AssetId:            assetIdBytes,
		PeerPubkey:         pubkeyBytes,
		FeeRateSatPerVbyte: feeRate,
		PushSat:            pushSat,
		LocalAmt:           localAmt,
	})
	if err != nil {
		return nil, errors.Wrap(err, "bc.FundBtlChannel")
	}

	return resp, nil
}

func (t Tap) BoxAddAssetInvoice(assetId string, assetAmount uint64, memo string, peerPubkey string) (*tapchannelrpc.AddInvoiceResponse, error) {
	conn, err := GetConn(tap, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)
	assetIdBytes, err := hex.DecodeString(assetId)
	if err != nil {
		return nil, errors.Wrap(err, "hex.DecodeString")
	}

	pubkeyBytes, err := hex.DecodeString(peerPubkey)
	if err != nil {
		return nil, errors.Wrap(err, "hex.DecodeString")
	}
	ctx := context.Background()
	tac := tapchannelrpc.NewTaprootAssetChannelsClient(conn)
	resp, err := tac.AddInvoice(ctx, &tapchannelrpc.AddInvoiceRequest{
		AssetId:     assetIdBytes,
		AssetAmount: assetAmount,
		PeerPubkey:  pubkeyBytes,
		InvoiceRequest: &lnrpc.Invoice{
			Memo: memo,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "tac.AddInvoice")
	}
	return resp, nil
}

func (t Tap) FwdBoxAddAssetInvoice(assetId string, assetAmount uint64, memo string, peerPubkey string) (*tapchannelrpc.AddInvoiceResponse, string, error) {
	conn, err := GetConn(tap, false)
	if err != nil {
		return nil, "", errors.Wrap(err, "GetConn")
	}

	defer Close(conn)
	assetIdBytes, err := hex.DecodeString(assetId)
	if err != nil {
		return nil, "", errors.Wrap(err, "hex.DecodeString")
	}

	pubkeyBytes, err := hex.DecodeString(peerPubkey)
	if err != nil {
		return nil, "", errors.Wrap(err, "hex.DecodeString")
	}
	ctx := context.Background()
	tac := tapchannelrpc.NewTaprootAssetChannelsClient(conn)
	resp, err := tac.AddInvoice(ctx, &tapchannelrpc.AddInvoiceRequest{
		AssetId:     assetIdBytes,
		AssetAmount: assetAmount,
		PeerPubkey:  pubkeyBytes,
		InvoiceRequest: &lnrpc.Invoice{
			Memo: memo,
		},
	})
	if err != nil {
		return nil, "", errors.Wrap(err, "tac.AddInvoice")
	}
	mappingInvoice, err := FwdtApplyInvoice(assetId, resp.InvoiceResult.GetPaymentRequest())
	if err != nil {
		return nil, "", err
	}
	return resp, mappingInvoice, nil
}

func (t Tap) DecodeAssetPayReq(assetId string, payReqString string) (*tapchannelrpc.AssetPayReqResponse, error) {
	conn, err := GetConn(tap, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)
	assetIdBytes, err := hex.DecodeString(assetId)
	if err != nil {
		return nil, errors.Wrap(err, "hex.DecodeString")
	}
	ctx := context.Background()
	tac := tapchannelrpc.NewTaprootAssetChannelsClient(conn)
	resp, err := tac.DecodeAssetPayReq(ctx, &tapchannelrpc.AssetPayReq{
		AssetId:      assetIdBytes,
		PayReqString: payReqString,
	})
	if err != nil {
		return nil, errors.Wrap(err, "tac.DecodeAssetPayReq")
	}
	return resp, nil
}

func (t Tap) BoxAssetChannelSendPayment(assetId string, pubkey string, paymentReq string, outgoingChanId string, feeLimitSat int, allowSelfPayment bool) (*lnrpc.Payment, error) {
	conn, err := GetConn(tap, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	assetIdBytes, err := hex.DecodeString(assetId)
	if err != nil {
		return nil, errors.Wrap(err, "hex.DecodeString")
	}
	peerPubkey, err := hex.DecodeString(pubkey)
	if err != nil {
		return nil, errors.Wrap(err, "hex.DecodeString")
	}
	client := tapchannelrpc.NewTaprootAssetChannelsClient(conn)
	req := tapchannelrpc.SendPaymentRequest{
		AssetId:    assetIdBytes,
		PeerPubkey: peerPubkey,
		PaymentRequest: &routerrpc.SendPaymentRequest{
			PaymentRequest:   paymentReq,
			TimeoutSeconds:   60,
			AllowSelfPayment: allowSelfPayment,
		},
		AllowOverpay: true,
	}

	req.PaymentRequest.OutgoingChanIds, err = ParseChanIDs([]string{outgoingChanId})
	if err != nil {
		return nil, err
	}
	if feeLimitSat != 0 {
		req.PaymentRequest.FeeLimitSat = int64(feeLimitSat)
	} else {
		req.PaymentRequest.FeeLimitSat = 20
	}
	resp, err := client.SendPayment(context.Background(), &req)
	if err != nil {
		return nil, errors.Wrap(err, "client.SendPayment")
	}
	for {
		resp1, err := resp.Recv()
		if err != nil {
			if err == io.EOF {
				return nil, errors.Wrap(err, "io.EOF")
			}
			return nil, errors.Wrap(err, "stream Recv")
		} else if resp1 != nil {
			resp2 := resp1.GetPaymentResult()
			if resp2 != nil {
				if resp2.Status == 2 || resp2.Status == 3 {
					return resp2, nil
				}
			}
		}
	}
}

func (t Tap) FwdBoxAssetChannelSendPayment(assetId string, pubkey string, paymentReq string, outgoingChanId string, feeLimitSat int, allowSelfPayment bool) (*lnrpc.Payment, error) {
	invoice, _ := CheckInvoiceIsCustody(paymentReq)
	if invoice != nil {
		payment, hash, err := t.BoxKeySendasset(uint64(invoice.Amount))
		if err != nil {
			return nil, err
		}
		err = PayToCustody(hash, invoice.Invoice)
		if err != nil {
			return nil, err
		}
		return payment, nil
	}

	conn, err := GetConn(tap, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}
	defer Close(conn)

	assetIdBytes, err := hex.DecodeString(assetId)
	if err != nil {
		return nil, errors.Wrap(err, "hex.DecodeString")
	}
	peerPubkey, err := hex.DecodeString(pubkey)
	if err != nil {
		return nil, errors.Wrap(err, "hex.DecodeString")
	}
	client := tapchannelrpc.NewTaprootAssetChannelsClient(conn)
	req := tapchannelrpc.SendPaymentRequest{
		AssetId:    assetIdBytes,
		PeerPubkey: peerPubkey,
		PaymentRequest: &routerrpc.SendPaymentRequest{
			PaymentRequest:   paymentReq,
			TimeoutSeconds:   60,
			AllowSelfPayment: allowSelfPayment,
		},
		AllowOverpay: true,
	}

	req.PaymentRequest.OutgoingChanIds, err = ParseChanIDs([]string{outgoingChanId})
	if err != nil {
		return nil, err
	}
	if feeLimitSat != 0 {
		req.PaymentRequest.FeeLimitSat = int64(feeLimitSat)
	} else {
		req.PaymentRequest.FeeLimitSat = 20
	}
	resp, err := client.SendPayment(context.Background(), &req)
	if err != nil {
		return nil, errors.Wrap(err, "client.SendPayment")
	}
	for {
		resp1, err := resp.Recv()
		if err != nil {
			if err == io.EOF {
				return nil, errors.Wrap(err, "io.EOF")
			}
			return nil, errors.Wrap(err, "stream Recv")
		} else if resp1 != nil {
			resp2 := resp1.GetPaymentResult()
			if resp2 != nil {
				if resp2.Status == 2 || resp2.Status == 3 {
					if resp2.Status == 2 {
						err := FwdtPayInvoice(paymentReq)
						if err != nil {
							return nil, err
						}
					}
					return resp2, nil
				}
			}
		}
	}
}

func (t Tap) AddrReceives() (*taprpc.AddrReceivesResponse, error) {
	conn, err := GetConn(tap, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	tac := taprpc.NewTaprootAssetsClient(conn)
	req := &taprpc.AddrReceivesRequest{}

	resp, err := tac.AddrReceives(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "tac.AddrReceives")
	}
	return resp, nil
}

func (t Tap) ListTransfers() (*taprpc.ListTransfersResponse, error) {
	conn, err := GetConn(tap, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	tac := taprpc.NewTaprootAssetsClient(conn)
	req := &taprpc.ListTransfersRequest{}

	resp, err := tac.ListTransfers(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "tac.ListTransfers")
	}
	return resp, nil
}

func (t Tap) KeySendasset(assetId string, amount uint64, pubkey string, outgoingChanId uint64) error {
	conn, err := GetConn(tap, false)
	if err != nil {
		return err
	}

	defer Close(conn)
	client := tapchannelrpc.NewTaprootAssetChannelsClient(conn)

	assetIdStr, err := hex.DecodeString(assetId)
	if err != nil {
		return err
	}
	peerPubkey, err := hex.DecodeString(pubkey)
	if err != nil {
		return err
	}

	req := &tapchannelrpc.SendPaymentRequest{
		AssetId:     assetIdStr,
		AssetAmount: amount,
		PaymentRequest: &routerrpc.SendPaymentRequest{
			Dest:              peerPubkey,
			Amt:               int64(rfqmath.DefaultOnChainHtlcSat),
			OutgoingChanId:    outgoingChanId,
			TimeoutSeconds:    30,
			DestCustomRecords: make(map[uint64][]byte),
		},
	}
	destRecords := req.PaymentRequest.DestCustomRecords
	_, isKeysend := destRecords[record.KeySendType]
	var rHash []byte
	var preimage lntypes.Preimage
	if _, err := rand.Read(preimage[:]); err != nil {
		return nil
	}
	if !isKeysend {
		destRecords[record.KeySendType] = preimage[:]
		hash := preimage.Hash()
		rHash = hash[:]

		req.PaymentRequest.PaymentHash = rHash

	}
	resp, err := client.SendPayment(context.Background(), req)
	if err != nil {
		return err
	}
	for {
		resp1, err := resp.Recv()
		if err != nil {
			if err == io.EOF {
				return err
			}
			return err
		} else if resp1 != nil {
			resp2 := resp1.GetPaymentResult()
			if resp2 != nil {
				if resp2.Status == 2 {
					return nil
				} else if resp2.Status == 3 {
					return nil
				}
			}
		}
	}
}

func (t Tap) AssetsListBalances() (*taprpc.ListBalancesResponse, error) {
	conn, err := GetConn(tap, false)
	if err != nil {
		return nil, err
	}

	defer Close(conn)
	client := taprpc.NewTaprootAssetsClient(conn)
	req := &taprpc.ListBalancesRequest{
		GroupBy: &taprpc.ListBalancesRequest_AssetId{
			AssetId: true,
		},
	}
	resp, err := client.ListBalances(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (t Tap) BoxAssetDecodePayReq(assetId string, payReq string) (*tapchannelrpc.AssetPayReqResponse, error) {
	conn, err := GetConn(tap, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	assetIdStr, err := hex.DecodeString(assetId)
	if err != nil {
		return nil, err
	}

	defer Close(conn)
	ctx := context.Background()
	tapchan := tapchannelrpc.NewTaprootAssetChannelsClient(conn)
	resp, err := tapchan.DecodeAssetPayReq(ctx, &tapchannelrpc.AssetPayReq{
		AssetId:      assetIdStr,
		PayReqString: payReq,
	})
	if err != nil {
		return nil, errors.Wrap(err, "tapchannelrpc.DecodeAssetPayReq")
	}
	return resp, nil
}

func (t Tap) ListBalances() (*taprpc.ListBalancesResponse, error) {
	conn, err := GetConn(tap, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	tac := taprpc.NewTaprootAssetsClient(conn)
	req := &taprpc.ListBalancesRequest{
		GroupBy: &taprpc.ListBalancesRequest_AssetId{
			AssetId: true,
		},
	}

	resp, err := tac.ListBalances(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "tac.ListBalances")
	}
	return resp, nil
}

func (t Tap) BoxKeySendasset(amount uint64) (*lnrpc.Payment, string, error) {
	const assetId = "97b98f3c45f926057d430ef71f20a6d3e25d7a00fbd1d7b72b306a49d48c9d8c"
	lnConn, err := GetConn(ln, false)
	if err != nil {
		return nil, "", err
	}

	defer Close(lnConn)

	lnClient := lnrpc.NewLightningClient(lnConn)
	chans, err := lnClient.ListChannels(context.Background(), &lnrpc.ListChannelsRequest{})
	if err != nil {
		return nil, "", err
	}

	var assetChan *lnrpc.Channel
	for _, channel := range chans.Channels {
		if channel.Private {
			assetChan = channel
			break
		}
	}
	if assetChan == nil {
		return nil, "", errors.New("asset channel not found")
	}

	conn, err := GetConn(tap, false)
	if err != nil {
		return nil, "", err
	}

	defer Close(conn)
	client := tapchannelrpc.NewTaprootAssetChannelsClient(conn)

	assetIdStr, err := hex.DecodeString(assetId)
	if err != nil {
		return nil, "", err
	}
	peerPubkey, err := hex.DecodeString(assetChan.RemotePubkey)
	if err != nil {
		return nil, "", err
	}

	req := &tapchannelrpc.SendPaymentRequest{
		AssetId:     assetIdStr,
		AssetAmount: amount,
		PaymentRequest: &routerrpc.SendPaymentRequest{
			Dest:              peerPubkey,
			Amt:               int64(rfqmath.DefaultOnChainHtlcSat),
			OutgoingChanId:    assetChan.ChanId,
			TimeoutSeconds:    30,
			DestCustomRecords: make(map[uint64][]byte),
		},
	}
	destRecords := req.PaymentRequest.DestCustomRecords
	_, isKeysend := destRecords[record.KeySendType]
	var rHash []byte
	var preimage lntypes.Preimage
	if _, err := rand.Read(preimage[:]); err != nil {
		return nil, "", err
	}
	if !isKeysend {
		destRecords[record.KeySendType] = preimage[:]
		hash := preimage.Hash()
		rHash = hash[:]

		req.PaymentRequest.PaymentHash = rHash

	}
	resp, err := client.SendPayment(context.Background(), req)
	if err != nil {
		return nil, "", err
	}
	for {
		resp1, err := resp.Recv()
		if err != nil {
			if err == io.EOF {
				return nil, "", err
			}
			return nil, "", err
		} else if resp1 != nil {
			resp2 := resp1.GetPaymentResult()
			if resp2 != nil {
				if resp2.Status == 2 {
					return resp2, resp2.PaymentHash, nil
				} else if resp2.Status == 3 {
					return nil, "", errors.New("payment failed")
				}
			}
		}
	}
}
