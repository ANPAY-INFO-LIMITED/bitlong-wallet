package api

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	taprootassets "github.com/lightninglabs/taproot-assets"
	"github.com/lightninglabs/taproot-assets/address"
	"github.com/lightninglabs/taproot-assets/asset"
	"github.com/lightninglabs/taproot-assets/proof"
	"github.com/lightninglabs/taproot-assets/rpcutils"
	"github.com/lightninglabs/taproot-assets/tappsbt"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/lightninglabs/taproot-assets/taprpc/assetwalletrpc"
	"github.com/lightninglabs/taproot-assets/taprpc/btlrpc"
	"github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	"github.com/lightninglabs/taproot-assets/tapsend"
	"github.com/lightninglabs/taproot-assets/universe"
	"github.com/lightningnetwork/lnd/keychain"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/walletrpc"
	"github.com/lightningnetwork/lnd/lntest"
	"github.com/lightningnetwork/lnd/lntest/wait"
	"github.com/pkg/errors"
	"github.com/wallet/base"
	"github.com/wallet/service/apiConnect"
)

type psbtTrustlessSwapResp string

func (p psbtTrustlessSwapResp) String() string {
	return string(p)
}

const (
	defaultTimeout = 600 * time.Second
	grpcTargetLnd  = "lnd"
	grpcTargetTapd = "tapd"

	nullString = ""
)

func getListEligibleCoins(assetId string) (coins []*btlrpc.Coin, err error) {
	ctxb := context.Background()
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()
	client := btlrpc.NewBtlClient(conn)

	var request = &btlrpc.GetListEligibleCoinsRequest{
		AssetId: assetId,
	}

	response, err := client.GetListEligibleCoins(ctxb, request)
	if err != nil {
		if strings.Contains(err.Error(), "failed to find coin(s) that satisfy given constraints") {
			LogError("client.GetListEligibleCoins", err)
			return nil, nil
		}

		return nil, errors.Wrap(err, "client.GetListEligibleCoins")
	}

	return response.Coins, nil
}

type EligibleCoin struct {
	AnchorPoint string `json:"anchor_point"`
	InternalKey string `json:"internal_key"`
	AssetAmount uint64 `json:"asset_amount"`
}

func btlrpcCoinsToGetListEligibleCoinsResults(coins []*btlrpc.Coin) (results []*EligibleCoin) {
	for _, coin := range coins {
		var anchorPoint = coin.AnchorPoint
		var internalKey = hex.EncodeToString(coin.InternalKey)
		var assetAmount = coin.AssetAmount
		results = append(results, &EligibleCoin{
			AnchorPoint: anchorPoint,
			InternalKey: internalKey,
			AssetAmount: assetAmount,
		})
	}
	return results
}

func psbtTrustlessSwapCreateSellOrderSign(assetId string, assetNum uint64, price int64, relativeLockTime uint64, feeRate uint64, checkRequire bool) (resp psbtTrustlessSwapResp, err error) {

	if feeRate > 500 {
		err = fmt.Errorf("feeRate too large, more than 500 (feeRate: %d)", feeRate)
		return nullString, err
	}

	ctxb := context.Background()

	var _chainParams *address.ChainParams
	var netParams *chaincfg.Params

	switch base.NetWork {
	case base.UseMainNet:
		_chainParams = &address.MainNetTap
		netParams = &chaincfg.MainNetParams
	case base.UseTestNet:
		_chainParams = &address.TestNet3Tap
		netParams = &chaincfg.TestNet3Params
	case base.UseRegTest:
		_chainParams = &address.RegressionNetTap
		netParams = &chaincfg.RegressionNetParams
	default:
		_chainParams = &address.RegressionNetTap
		netParams = &chaincfg.RegressionNetParams
	}

	var assetID asset.ID
	assetIdBytes, err := hex.DecodeString(assetId)
	if err != nil {
		return nullString, errors.Wrap(err, "hex.DecodeString")
	}
	if len(assetIdBytes) != sha256.Size {
		err = fmt.Errorf("assetId must be of length %d", sha256.Size)
		return nullString, errors.Wrap(err, "len(assetIdBytes) != sha256.Size")
	}
	copy(assetID[:], assetIdBytes)

	aliceDummyScriptKey, aliceAnchorInternalKey, err := deriveKeys()
	if err != nil {
		return nullString, errors.Wrap(err, "deriveKeys")
	}

	vPkt := tappsbt.ForInteractiveSend(assetID, assetNum, aliceDummyScriptKey, 0, relativeLockTime, 1, aliceAnchorInternalKey, asset.V0, _chainParams)

	fundResp, err := maybeFundPacket(vPkt)
	if err != nil {
		return nullString, errors.Wrap(err, "maybeFundPacket")
	}
	vPkt, err = tappsbt.Decode(fundResp.FundedPsbt)
	if err != nil {
		return nullString, errors.Wrap(err, "tappsbt.Decode")
	}
	if err = requireLenBe(vPkt.Inputs, 1); err != nil && checkRequire {
		return nullString, errors.Wrap(err, "requireLenBe(vPkt.Inputs, 1)")
	}
	if err = requireLengths(vPkt.Outputs, 1, 2); err != nil && checkRequire {
		return nullString, errors.Wrap(err, "requireLengths(vPkt.Outputs, 1, 2)")
	}

	for i := range vPkt.Outputs {
		vPkt.Inputs[i].SighashType = txscript.SigHashNone
	}

	if err = requireEqual(vPkt.Outputs[0].Type, tappsbt.TypeSimple); err != nil && checkRequire {
		return nullString, errors.Wrap(err, "requireEqual(vPkt.Outputs[0].Type, tappsbt.TypeSimple)")
	}
	err = tapsend.PrepareOutputAssets(ctxb, vPkt)
	if err != nil {
		return nullString, errors.Wrap(err, "tapsend.PrepareOutputAssets")
	}

	{
	}
	fundedPsbtBytes, err := tappsbt.Encode(vPkt)
	if err != nil {
		return nullString, errors.Wrap(err, "tappsbt.Encode")
	}

	var awc assetwalletrpc.AssetWalletClient
	conn, clearUp, err := apiConnect.GetConnection(grpcTargetTapd, false)
	if err != nil {
		return nullString, errors.Wrap(err, "apiConnect.GetConnection "+grpcTargetTapd)
	}
	defer clearUp()
	awc = assetwalletrpc.NewAssetWalletClient(conn)

	signedResp, err := awc.SignVirtualPsbt(ctxb, &assetwalletrpc.SignVirtualPsbtRequest{FundedPsbt: fundedPsbtBytes})
	if err != nil {
		return nullString, errors.Wrap(err, "awc.SignVirtualPsbt")
	}

	if err = requireContains(signedResp.SignedInputs, uint32(0)); err != nil && checkRequire {
		return nullString, errors.Wrap(err, "requireContains(signedResp.SignedInputs, uint32(0))")
	}

	vPkt, err = tappsbt.Decode(signedResp.SignedPsbt)
	if err != nil {
		return nullString, errors.Wrap(err, "tappsbt.Decode")
	}

	btcPacket, err := tapsend.PrepareAnchoringTemplate([]*tappsbt.VPacket{vPkt})
	if err != nil {
		return nullString, errors.Wrap(err, "tapsend.PrepareAnchoringTemplate")
	}

	if err = requireLen(btcPacket.Inputs, len(vPkt.Inputs)); err != nil && checkRequire {
		return nullString, errors.Wrap(err, "requireLen(btcPacket.Inputs, len(vPkt.Inputs))")
	}
	if err = requireLen(btcPacket.Outputs, len(vPkt.Outputs)+1); err != nil && checkRequire {
		return nullString, errors.Wrap(err, "requireLen(btcPacket.Outputs, len(vPkt.Outputs)+1)")
	}

	var lnc lnrpc.LightningClient
	LndConn, LndClearUp, err := apiConnect.GetConnection(grpcTargetLnd, false)
	if err != nil {
		return nullString, errors.Wrap(err, "apiConnect.GetConnection "+grpcTargetLnd)
	}
	defer LndClearUp()
	lnc = lnrpc.NewLightningClient(LndConn)

	addrResp, err := lnc.NewAddress(ctxb, &lnrpc.NewAddressRequest{Type: lnrpc.AddressType_TAPROOT_PUBKEY})
	if err != nil {
		return nullString, errors.Wrap(err, "lnc.NewAddress")
	}

	aliceP2TR, err := btcutil.DecodeAddress(addrResp.Address, netParams)
	if err != nil {
		return nullString, errors.Wrap(err, "btcutil.DecodeAddress")
	}

	alicePkScript, err := txscript.PayToAddrScript(aliceP2TR)
	if err != nil {
		return nullString, errors.Wrap(err, "txscript.PayToAddrScript")
	}

	btcPacket.UnsignedTx.TxOut[0].PkScript = alicePkScript
	btcPacket.UnsignedTx.TxOut[0].Value = price

	derivation, trDerivation, err := getAddressBip32Derivation(addrResp.Address)
	if err != nil {
		return nullString, errors.Wrap(err, "getAddressBip32Derivation")
	}

	btcPacket.Outputs[0].Bip32Derivation = []*psbt.Bip32Derivation{derivation}
	btcPacket.Outputs[0].TaprootBip32Derivation = []*psbt.TaprootBip32Derivation{trDerivation}
	btcPacket.Outputs[0].TaprootInternalKey = trDerivation.XOnlyPubKey

	var b bytes.Buffer
	err = btcPacket.Serialize(&b)
	if err != nil {
		return nullString, errors.Wrap(err, "btcPacket.Serialize")
	}

	cvpResp, err := awc.CommitVirtualPsbts(
		ctxb, &assetwalletrpc.CommitVirtualPsbtsRequest{
			VirtualPsbts:       [][]byte{signedResp.SignedPsbt},
			AnchorPsbt:         b.Bytes(),
			AnchorChangeOutput: &assetwalletrpc.CommitVirtualPsbtsRequest_Add{Add: true},
			Fees:               &assetwalletrpc.CommitVirtualPsbtsRequest_SatPerVbyte{SatPerVbyte: feeRate},
		},
	)
	if err != nil {
		return nullString, errors.Wrap(err, "awc.CommitVirtualPsbts")
	}

	btcPacket, err = psbt.NewFromRawBytes(bytes.NewReader(cvpResp.AnchorPsbt), false)
	if err != nil {
		return nullString, errors.Wrap(err, "psbt.NewFromRawBytes(bytes.NewReader(cvpResp.AnchorPsbt), false)")
	}

	prevInputLen := len(vPkt.Inputs)
	prevOutputLen := len(vPkt.Outputs)
	for i := range prevInputLen {
		btcPacket.Inputs[i].SighashType = txscript.SigHashSingle | txscript.SigHashAnyOneCanPay
	}

	nowInputLen := len(btcPacket.Inputs)

	btcPacket.Inputs = append(
		btcPacket.Inputs[:prevInputLen], btcPacket.Inputs[nowInputLen:]...,
	)
	btcPacket.UnsignedTx.TxIn = append(
		btcPacket.UnsignedTx.TxIn[:prevInputLen], btcPacket.UnsignedTx.TxIn[nowInputLen:]...,
	)

	btcPacket.Outputs = btcPacket.Outputs[:1+prevOutputLen]
	btcPacket.UnsignedTx.TxOut = btcPacket.UnsignedTx.TxOut[:1+prevOutputLen]

	b.Reset()
	err = btcPacket.Serialize(&b)
	if err != nil {
		return nullString, errors.Wrap(err, "btcPacket.Serialize")
	}

	wkc := walletrpc.NewWalletKitClient(LndConn)

	signPsbtResp, err := wkc.SignPsbt(ctxb, &walletrpc.SignPsbtRequest{
		FundedPsbt: b.Bytes(),
	})
	if err != nil {
		return nullString, errors.Wrap(err, "wkc.SignPsbt")
	}

	if err = requireLen(signPsbtResp.SignedInputs, prevInputLen); err != nil && checkRequire {
		return nullString, errors.Wrap(err, "requireLen(signPsbtResp.SignedInputs, prevInputLen)")
	}

	btcPacket, err = psbt.NewFromRawBytes(bytes.NewReader(signPsbtResp.SignedPsbt), false)
	if err != nil {
		return nullString, errors.Wrap(err, "psbt.NewFromRawBytes(bytes.NewReader(signPsbtResp.SignedPsbt), false)")
	}

	if err = requireLen(btcPacket.Inputs, prevInputLen); err != nil && checkRequire {
		return nullString, errors.Wrap(err, "requireLen(btcPacket.Inputs, prevInputLen)")
	}
	if err = requireLen(btcPacket.Outputs, 1+prevOutputLen); err != nil && checkRequire {
		return nullString, errors.Wrap(err, "requireLen(btcPacket.Outputs, 1+prevOutputLen)")
	}

	signedVPsbtBytes, err := tappsbt.Encode(vPkt)
	if err != nil {
		return nullString, errors.Wrap(err, "tappsbt.Encode")
	}

	signedPsbtBytes := signPsbtResp.SignedPsbt

	var _resp = psbtTrustlessSwapCreateSellOrderResponse{
		SignedVPsbtBytes: signedVPsbtBytes,
		SignedPsbtBytes:  signedPsbtBytes,
	}
	encoded, err := EncodeDataToBase64(_resp)
	if err != nil {
		return nullString, errors.Wrap(err, "EncodeDataToBase64")
	}
	return psbtTrustlessSwapResp(encoded), nil
}

func psbtTrustlessSwapCreateSellOrderSignWithOneFilter(assetId string, assetNum uint64, price int64, relativeLockTime uint64, feeRate uint64, checkRequire bool, anchorPoint string, internalKey string) (resp psbtTrustlessSwapResp, err error) {

	if feeRate > 500 {
		err = fmt.Errorf("feeRate too large, more than 500 (feeRate: %d)", feeRate)
		return nullString, err
	}

	ctxb := context.Background()

	var _chainParams *address.ChainParams
	var netParams *chaincfg.Params

	switch base.NetWork {
	case base.UseMainNet:
		_chainParams = &address.MainNetTap
		netParams = &chaincfg.MainNetParams
	case base.UseTestNet:
		_chainParams = &address.TestNet3Tap
		netParams = &chaincfg.TestNet3Params
	case base.UseRegTest:
		_chainParams = &address.RegressionNetTap
		netParams = &chaincfg.RegressionNetParams
	default:
		_chainParams = &address.RegressionNetTap
		netParams = &chaincfg.RegressionNetParams
	}

	var assetID asset.ID
	assetIdBytes, err := hex.DecodeString(assetId)
	if err != nil {
		return nullString, errors.Wrap(err, "hex.DecodeString")
	}
	if len(assetIdBytes) != sha256.Size {
		err = fmt.Errorf("assetId must be of length %d", sha256.Size)
		return nullString, errors.Wrap(err, "len(assetIdBytes) != sha256.Size")
	}
	copy(assetID[:], assetIdBytes)

	aliceDummyScriptKey, aliceAnchorInternalKey, err := deriveKeys()
	if err != nil {
		return nullString, errors.Wrap(err, "deriveKeys")
	}

	vPkt := tappsbt.ForInteractiveSend(assetID, assetNum, aliceDummyScriptKey, 0, relativeLockTime, 1, aliceAnchorInternalKey, asset.V0, _chainParams)

	var coinsFilter *btlrpc.CoinsFilter
	if anchorPoint != "" {
		var key []byte
		key, err = hex.DecodeString(internalKey)
		if err != nil {
			return nullString, errors.Wrap(err, "hex.DecodeString(internalKey)")
		}
		coinsFilter = &btlrpc.CoinsFilter{Coins: []*btlrpc.Coin{
			{
				AnchorPoint: anchorPoint,
				InternalKey: key,
			},
		}}
	}
	fundResp, err := fundPacket2(vPkt, coinsFilter)
	if err != nil {
		return nullString, errors.Wrap(err, "maybeFundPacket")
	}
	vPkt, err = tappsbt.Decode(fundResp.FundedPsbt)
	if err != nil {
		return nullString, errors.Wrap(err, "tappsbt.Decode")
	}
	if err = requireLenBe(vPkt.Inputs, 1); err != nil && checkRequire {
		return nullString, errors.Wrap(err, "requireLenBe(vPkt.Inputs, 1)")
	}
	if err = requireLengths(vPkt.Outputs, 1, 2); err != nil && checkRequire {
		return nullString, errors.Wrap(err, "requireLengths(vPkt.Outputs, 1, 2)")
	}

	for i := range vPkt.Inputs {
		vPkt.Inputs[i].SighashType = txscript.SigHashNone
	}

	if err = requireEqual(vPkt.Outputs[0].Type, tappsbt.TypeSimple); err != nil && checkRequire {
		return nullString, errors.Wrap(err, "requireEqual(vPkt.Outputs[0].Type, tappsbt.TypeSimple)")
	}
	err = tapsend.PrepareOutputAssets(ctxb, vPkt)
	if err != nil {
		return nullString, errors.Wrap(err, "tapsend.PrepareOutputAssets")
	}

	{
	}
	fundedPsbtBytes, err := tappsbt.Encode(vPkt)
	if err != nil {
		return nullString, errors.Wrap(err, "tappsbt.Encode")
	}

	var awc assetwalletrpc.AssetWalletClient
	conn, clearUp, err := apiConnect.GetConnection(grpcTargetTapd, false)
	if err != nil {
		return nullString, errors.Wrap(err, "apiConnect.GetConnection "+grpcTargetTapd)
	}
	defer clearUp()
	awc = assetwalletrpc.NewAssetWalletClient(conn)

	signedResp, err := awc.SignVirtualPsbt(ctxb, &assetwalletrpc.SignVirtualPsbtRequest{FundedPsbt: fundedPsbtBytes})
	if err != nil {
		return nullString, errors.Wrap(err, "awc.SignVirtualPsbt")
	}

	if err = requireContains(signedResp.SignedInputs, uint32(0)); err != nil && checkRequire {
		return nullString, errors.Wrap(err, "requireContains(signedResp.SignedInputs, uint32(0))")
	}

	vPkt, err = tappsbt.Decode(signedResp.SignedPsbt)
	if err != nil {
		return nullString, errors.Wrap(err, "tappsbt.Decode")
	}

	btcPacket, err := tapsend.PrepareAnchoringTemplate([]*tappsbt.VPacket{vPkt})
	if err != nil {
		return nullString, errors.Wrap(err, "tapsend.PrepareAnchoringTemplate")
	}

	if err = requireLen(btcPacket.Inputs, len(vPkt.Inputs)); err != nil && checkRequire {
		return nullString, errors.Wrap(err, "requireLen(btcPacket.Inputs, len(vPkt.Inputs))")
	}
	if err = requireLen(btcPacket.Outputs, len(vPkt.Outputs)+1); err != nil && checkRequire {
		return nullString, errors.Wrap(err, "requireLen(btcPacket.Outputs, len(vPkt.Outputs)+1)")
	}

	var lnc lnrpc.LightningClient
	LndConn, LndClearUp, err := apiConnect.GetConnection(grpcTargetLnd, false)
	if err != nil {
		return nullString, errors.Wrap(err, "apiConnect.GetConnection "+grpcTargetLnd)
	}
	defer LndClearUp()
	lnc = lnrpc.NewLightningClient(LndConn)

	addrResp, err := lnc.NewAddress(ctxb, &lnrpc.NewAddressRequest{Type: lnrpc.AddressType_TAPROOT_PUBKEY})
	if err != nil {
		return nullString, errors.Wrap(err, "lnc.NewAddress")
	}

	aliceP2TR, err := btcutil.DecodeAddress(addrResp.Address, netParams)
	if err != nil {
		return nullString, errors.Wrap(err, "btcutil.DecodeAddress")
	}

	alicePkScript, err := txscript.PayToAddrScript(aliceP2TR)
	if err != nil {
		return nullString, errors.Wrap(err, "txscript.PayToAddrScript")
	}

	btcPacket.UnsignedTx.TxOut[0].PkScript = alicePkScript
	btcPacket.UnsignedTx.TxOut[0].Value = price

	derivation, trDerivation, err := getAddressBip32Derivation(addrResp.Address)
	if err != nil {
		return nullString, errors.Wrap(err, "getAddressBip32Derivation")
	}

	btcPacket.Outputs[0].Bip32Derivation = []*psbt.Bip32Derivation{derivation}
	btcPacket.Outputs[0].TaprootBip32Derivation = []*psbt.TaprootBip32Derivation{trDerivation}
	btcPacket.Outputs[0].TaprootInternalKey = trDerivation.XOnlyPubKey

	var b bytes.Buffer
	err = btcPacket.Serialize(&b)
	if err != nil {
		return nullString, errors.Wrap(err, "btcPacket.Serialize")
	}

	cvpResp, err := awc.CommitVirtualPsbts(
		ctxb, &assetwalletrpc.CommitVirtualPsbtsRequest{
			VirtualPsbts:       [][]byte{signedResp.SignedPsbt},
			AnchorPsbt:         b.Bytes(),
			AnchorChangeOutput: &assetwalletrpc.CommitVirtualPsbtsRequest_Add{Add: true},
			Fees:               &assetwalletrpc.CommitVirtualPsbtsRequest_SatPerVbyte{SatPerVbyte: feeRate},
		},
	)
	if err != nil {
		return nullString, errors.Wrap(err, "awc.CommitVirtualPsbts")
	}

	btcPacket, err = psbt.NewFromRawBytes(bytes.NewReader(cvpResp.AnchorPsbt), false)
	if err != nil {
		return nullString, errors.Wrap(err, "psbt.NewFromRawBytes(bytes.NewReader(cvpResp.AnchorPsbt), false)")
	}

	prevInputLen := len(vPkt.Inputs)
	prevOutputLen := len(vPkt.Outputs)
	for i := range prevInputLen {
		btcPacket.Inputs[i].SighashType = txscript.SigHashSingle | txscript.SigHashAnyOneCanPay
	}

	nowInputLen := len(btcPacket.Inputs)

	btcPacket.Inputs = append(
		btcPacket.Inputs[:prevInputLen], btcPacket.Inputs[nowInputLen:]...,
	)
	btcPacket.UnsignedTx.TxIn = append(
		btcPacket.UnsignedTx.TxIn[:prevInputLen], btcPacket.UnsignedTx.TxIn[nowInputLen:]...,
	)

	btcPacket.Outputs = btcPacket.Outputs[:1+prevOutputLen]
	btcPacket.UnsignedTx.TxOut = btcPacket.UnsignedTx.TxOut[:1+prevOutputLen]

	b.Reset()
	err = btcPacket.Serialize(&b)
	if err != nil {
		return nullString, errors.Wrap(err, "btcPacket.Serialize")
	}

	wkc := walletrpc.NewWalletKitClient(LndConn)

	signPsbtResp, err := wkc.SignPsbt(ctxb, &walletrpc.SignPsbtRequest{
		FundedPsbt: b.Bytes(),
	})
	if err != nil {
		return nullString, errors.Wrap(err, "wkc.SignPsbt")
	}

	if err = requireLen(signPsbtResp.SignedInputs, prevInputLen); err != nil && checkRequire {
		return nullString, errors.Wrap(err, "requireLen(signPsbtResp.SignedInputs, prevInputLen)")
	}

	btcPacket, err = psbt.NewFromRawBytes(bytes.NewReader(signPsbtResp.SignedPsbt), false)
	if err != nil {
		return nullString, errors.Wrap(err, "psbt.NewFromRawBytes(bytes.NewReader(signPsbtResp.SignedPsbt), false)")
	}

	if err = requireLen(btcPacket.Inputs, prevInputLen); err != nil && checkRequire {
		return nullString, errors.Wrap(err, "requireLen(btcPacket.Inputs, prevInputLen)")
	}
	if err = requireLen(btcPacket.Outputs, 1+prevOutputLen); err != nil && checkRequire {
		return nullString, errors.Wrap(err, "requireLen(btcPacket.Outputs, 1+prevOutputLen)")
	}

	signedVPsbtBytes, err := tappsbt.Encode(vPkt)
	if err != nil {
		return nullString, errors.Wrap(err, "tappsbt.Encode")
	}

	signedPsbtBytes := signPsbtResp.SignedPsbt

	var _resp = psbtTrustlessSwapCreateSellOrderResponse{
		SignedVPsbtBytes: signedVPsbtBytes,
		SignedPsbtBytes:  signedPsbtBytes,
	}
	encoded, err := EncodeDataToBase64(_resp)
	if err != nil {
		return nullString, errors.Wrap(err, "EncodeDataToBase64")
	}
	return psbtTrustlessSwapResp(encoded), nil
}

type psbtTrustlessSwapCreateSellOrderResponse struct {
	SignedVPsbtBytes []byte `json:"signed_v_psbt_bytes"`
	SignedPsbtBytes  []byte `json:"signed_psbt_bytes"`
}

func (r *psbtTrustlessSwapCreateSellOrderResponse) FromStr(resp string) (err error) {
	return DecodeBase64ToData(resp, r)
}

func psbtTrustlessSwapBuySOrderSign(signedVPsbtBytes []byte, signedPsbtBytes []byte, feeRate uint64, deliveryAddr string, checkRequire bool) (resp psbtTrustlessSwapResp, err error) {

	if feeRate > 500 {
		err = fmt.Errorf("feeRate too large, more than 500 (feeRate: %d)", feeRate)
		return nullString, err
	}

	ctxb := context.Background()

	var netParams *chaincfg.Params

	switch base.NetWork {
	case base.UseMainNet:
		netParams = &chaincfg.MainNetParams
	case base.UseTestNet:
		netParams = &chaincfg.TestNet3Params
	case base.UseRegTest:
		netParams = &chaincfg.RegressionNetParams
	default:
		netParams = &chaincfg.RegressionNetParams
	}

	bobVPsbt, err := tappsbt.Decode(signedVPsbtBytes)
	if err != nil {
		return nullString, errors.Wrap(err, "tappsbt.Decode")
	}

	if err = requireLengths(bobVPsbt.Outputs, 1, 2); err != nil && checkRequire {
		return nullString, errors.Wrap(err, "requireLengths(bobVPsbt.Outputs, 1, 2)")
	}

	var btcPacket *psbt.Packet
	btcPacket, err = psbt.NewFromRawBytes(bytes.NewReader(signedPsbtBytes), false)
	if err != nil {
		return nullString, errors.Wrap(err, "psbt.NewFromRawBytes(bytes.NewReader(signedPsbtBytes), false)")
	}

	bobScriptKey, bobAnchorInternalKey, err := deriveKeys()
	if err != nil {
		return nullString, errors.Wrap(err, "deriveKeys")
	}

	bobVOut := bobVPsbt.Outputs[0]
	bobVOut.ScriptKey = bobScriptKey
	bobVOut.AnchorOutputBip32Derivation = nil
	bobVOut.AnchorOutputTaprootBip32Derivation = nil

	bobVOut.SetAnchorInternalKey(bobAnchorInternalKey, netParams.HDCoinType)
	deliveryAddrStr := fmt.Sprintf("%s://%s", proof.UniverseRpcCourierType, deliveryAddr)
	_deliveryAddr, err := url.Parse(deliveryAddrStr)
	if err != nil {
		return nullString, errors.Wrap(err, "url.Parse(deliveryAddrStr)")
	}
	for i := range bobVPsbt.Outputs {
		bobVPsbt.Outputs[i].ProofDeliveryAddress = _deliveryAddr
	}

	btcPacket.Outputs[1].TaprootInternalKey = schnorr.SerializePubKey(bobAnchorInternalKey.PubKey)
	btcPacket.Outputs[1].Bip32Derivation = bobVOut.AnchorOutputBip32Derivation
	btcPacket.Outputs[1].TaprootBip32Derivation = bobVOut.AnchorOutputTaprootBip32Derivation

	witnessBackupMap := make(map[int][]asset.Witness)
	for i := range bobVPsbt.Outputs {
		witnessBackupMap[i] = bobVPsbt.Outputs[i].Asset.PrevWitnesses
	}

	err = tapsend.PrepareOutputAssets(ctxb, bobVPsbt)
	if err != nil {
		return nullString, errors.Wrap(err, "tapsend.PrepareOutputAssets")
	}

	if err = requireLengths(bobVPsbt.Outputs, 1, 2); err != nil && checkRequire {
		return nullString, errors.Wrap(err, "requireLengths(bobVPsbt.Outputs, 1, 2)")
	}
	if err = requireEqual(bobVPsbt.Outputs[0].ScriptKey, bobVPsbt.Outputs[0].Asset.ScriptKey); err != nil && checkRequire {
		return nullString, errors.Wrap(err, "requireEqual(bobVPsbt.Outputs[0].ScriptKey, bobVPsbt.Outputs[0].Asset.ScriptKey)")
	}

	for i := range bobVPsbt.Outputs {
		bobVPsbt.Outputs[i].Asset.PrevWitnesses = witnessBackupMap[i]
	}

	bobVPsbtBytes, err := tappsbt.Encode(bobVPsbt)
	if err != nil {
		return nullString, errors.Wrap(err, "tappsbt.Encode")
	}

	var b bytes.Buffer
	err = btcPacket.Serialize(&b)
	if err != nil {
		return nullString, errors.Wrap(err, "btcPacket.Serialize")
	}

	var awc assetwalletrpc.AssetWalletClient
	conn, clearUp, err := apiConnect.GetConnection(grpcTargetTapd, false)
	if err != nil {
		return nullString, errors.Wrap(err, "apiConnect.GetConnection "+grpcTargetTapd)
	}
	defer clearUp()
	awc = assetwalletrpc.NewAssetWalletClient(conn)

	var cvpResp *assetwalletrpc.CommitVirtualPsbtsResponse
	cvpResp, err = awc.CommitVirtualPsbts(
		ctxb, &assetwalletrpc.CommitVirtualPsbtsRequest{
			VirtualPsbts:       [][]byte{bobVPsbtBytes},
			AnchorPsbt:         b.Bytes(),
			AnchorChangeOutput: &assetwalletrpc.CommitVirtualPsbtsRequest_Add{Add: true},
			Fees:               &assetwalletrpc.CommitVirtualPsbtsRequest_SatPerVbyte{SatPerVbyte: feeRate},
		},
	)
	if err != nil {
		return nullString, errors.Wrap(err, "awc.CommitVirtualPsbts")
	}

	bobVPsbt, err = tappsbt.Decode(cvpResp.VirtualPsbts[0])
	if err != nil {
		return nullString, errors.Wrap(err, "tappsbt.Decode")
	}

	LndConn, LndClearUp, err := apiConnect.GetConnection(grpcTargetLnd, false)
	if err != nil {
		return nullString, errors.Wrap(err, "apiConnect.GetConnection "+grpcTargetLnd)
	}
	defer LndClearUp()
	wkc := walletrpc.NewWalletKitClient(LndConn)

	signResp, err := wkc.SignPsbt(ctxb, &walletrpc.SignPsbtRequest{FundedPsbt: cvpResp.AnchorPsbt})
	if err != nil {
		return nullString, errors.Wrap(err, "wkc.SignPsbt")
	}

	finalPsbt, err := psbt.NewFromRawBytes(bytes.NewReader(signResp.SignedPsbt), false)
	if err != nil {
		return nullString, errors.Wrap(err, "psbt.NewFromRawBytes(bytes.NewReader(signResp.SignedPsbt), false)")
	}

	if err = requireLenBe(finalPsbt.Inputs, 2); err != nil && checkRequire {
		return nullString, errors.Wrap(err, "requireLenBe(finalPsbt.Inputs, 2)")
	}

	signedPkt, err := finalizePacket(finalPsbt)
	if err != nil {
		return nullString, errors.Wrap(err, "finalizePacket")
	}

	if err = requireEqual(true, signedPkt.IsComplete()); err != nil && checkRequire {
		return nullString, errors.Wrap(err, "requireEqual(true, signedPkt.IsComplete())")
	}

	bobScriptKeyBytes := bobScriptKey.PubKey.SerializeCompressed()

	bobOutputIndex := uint32(1)
	transferTXID := finalPsbt.UnsignedTx.TxHash()
	bobAssetOutpoint := fmt.Sprintf("%s:%d", transferTXID.String(), bobOutputIndex)

	scriptKeyBytes, assetOutpoint := bobScriptKeyBytes, bobAssetOutpoint

	signedPktBytes, err := psbtPacketToBytes(signedPkt)
	if err != nil {
		return nullString, errors.Wrap(err, "psbtPacketToBytes")
	}

	_bobVPsbtBytes, err := tappsbt.Encode(bobVPsbt)
	if err != nil {
		return nullString, errors.Wrap(err, "tappsbt.Encode")
	}

	var _resp = psbtTrustlessSwapBuySOrderResponse{
		SignedPktBytes: signedPktBytes,
		BobVPsbtBytes:  _bobVPsbtBytes,
		CvpResp:        cvpResp,
		ScriptKeyBytes: scriptKeyBytes,
		AssetOutpoint:  assetOutpoint,
	}
	encoded, err := EncodeDataToBase64(_resp)
	if err != nil {
		return nullString, errors.Wrap(err, "EncodeDataToBase64")
	}
	return psbtTrustlessSwapResp(encoded), nil
}

type psbtTrustlessSwapBuySOrderResponse struct {
	SignedPktBytes []byte                                     `json:"signed_pkt_bytes"`
	BobVPsbtBytes  []byte                                     `json:"bob_v_psbt_bytes"`
	CvpResp        *assetwalletrpc.CommitVirtualPsbtsResponse `json:"cvp_resp"`
	ScriptKeyBytes []byte                                     `json:"script_key_bytes"`
	AssetOutpoint  string                                     `json:"asset_outpoint"`
}

func (r *psbtTrustlessSwapBuySOrderResponse) FromStr(resp string) (err error) {
	return DecodeBase64ToData(resp, r)
}

func psbtTrustlessSwapPublishTx(signedPkt *psbt.Packet, bobVPsbt *tappsbt.VPacket, cvpResp *assetwalletrpc.CommitVirtualPsbtsResponse) (saResp *taprpc.SendAssetResponse, err error) {
	logResp, err := logAndPublish(signedPkt, []*tappsbt.VPacket{bobVPsbt}, nil, cvpResp)
	if err != nil {
		return nil, errors.Wrap(err, "logAndPublish")
	}
	return logResp, nil
}

type getLastProofResponse struct {
	Errnos int    `json:"errno"`
	ErrMsg string `json:"errmsg"`
	Data   string `json:"data"`
}

func requestToGetLastProof(token string, scriptKey string, outpoint string, assetId string) (string, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	_url := serverDomainOrSocket + "/proof/get_last_proof?" + "script_key=" + scriptKey + "&outpoint=" + outpoint + "&asset_id=" + assetId
	requestJsonBytes, err := json.Marshal(nil)
	if err != nil {
		return "", err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("GET", _url, payload)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var response getLastProofResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}
	if response.ErrMsg != "" {
		return "", errors.New(response.ErrMsg)
	}
	return response.Data, nil
}

func LastProofFromStr(lastProofStr string) (lastProof []byte, err error) {
	return base64.StdEncoding.DecodeString(lastProofStr)
}

func psbtTrustlessSwapBuySOrderProof(scriptKeyBytes []byte, assetOutpoint string, assetId string, lastProof []byte) (err error) {

	bobScriptKeyBytes, bobAssetOutpoint := scriptKeyBytes, assetOutpoint

	genInfo, assetGroup, err := getAssetGenInfoAndAssetGroup(assetId)
	if err != nil {
		return errors.Wrap(err, "getAssetGenInfoAndAssetGroup")
	}

	importResp, err := transferProofUniRPC(lastProof)
	if err != nil {
		return errors.Wrap(err, "transferProofUniRPC")
	}

	ctxb := context.Background()
	ctxt, cancel := context.WithTimeout(ctxb, defaultTimeout)
	defer cancel()
	conn, clearUp, err := apiConnect.GetConnection(grpcTargetTapd, false)
	if err != nil {
		return errors.Wrap(err, "apiConnect.GetConnection "+grpcTargetTapd)
	}
	defer clearUp()
	tac := taprpc.NewTaprootAssetsClient(conn)

	var (
		assetID         asset.ID
		tweakedGroupKey []byte
		txid            string
		idx             uint32
	)
	copy(assetID[:], genInfo.AssetId)
	if assetGroup == nil {
		tweakedGroupKey = nil
	} else {
		tweakedGroupKey = assetGroup.TweakedGroupKey
	}

	txid, idxStr, err := _outpointToTransactionAndIndex(bobAssetOutpoint)
	if err != nil {
		return errors.Wrap(err, "OutpointToTransactionAndIndex")
	}
	index, err := strconv.ParseUint(idxStr, 10, 32)
	if err != nil {
		return errors.Wrap(err, "strconv.ParseUint")
	}
	idx = uint32(index)

	transferTXID, err := chainhash.NewHashFromStr(txid)
	if err != nil {
		return errors.Wrap(err, "chainhash.NewHashFromStr")
	}
	bobOutputIndex := idx

	registerResp, err := tac.RegisterTransfer(ctxt, &taprpc.RegisterTransferRequest{
		AssetId:   assetID[:],
		GroupKey:  tweakedGroupKey,
		ScriptKey: bobScriptKeyBytes,
		Outpoint: &taprpc.OutPoint{
			Txid:        transferTXID[:],
			OutputIndex: bobOutputIndex,
		},
	})

	_, _ = importResp, registerResp

	return nil
}

func psbtTrustlessSwapCreateBuyOrderSign() {

}

func psbtTrustlessSwapSellBOrderSign() {

}

func deriveKeys() (scriptKey asset.ScriptKey, internalKey keychain.KeyDescriptor, err error) {
	ctx := context.Background()
	ctxt, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	conn, clearUp, err := apiConnect.GetConnection(grpcTargetTapd, false)
	if err != nil {
		return asset.ScriptKey{}, keychain.KeyDescriptor{}, errors.Wrap(err, "apiConnect.GetConnection "+grpcTargetTapd)
	}
	defer clearUp()

	awc := assetwalletrpc.NewAssetWalletClient(conn)

	scriptKeyDesc, err := awc.NextScriptKey(
		ctxt, &assetwalletrpc.NextScriptKeyRequest{
			KeyFamily: uint32(asset.TaprootAssetsKeyFamily),
		},
	)
	if err != nil {
		return asset.ScriptKey{}, keychain.KeyDescriptor{}, errors.Wrap(err, "awc.NextScriptKey")
	}

	_scriptKey, err := rpcutils.UnmarshalScriptKey(scriptKeyDesc.ScriptKey)
	if err != nil {
		return asset.ScriptKey{}, keychain.KeyDescriptor{}, errors.Wrap(err, "taprpc.UnmarshalScriptKey")
	}

	internalKeyDesc, err := awc.NextInternalKey(
		ctxt, &assetwalletrpc.NextInternalKeyRequest{
			KeyFamily: uint32(asset.TaprootAssetsKeyFamily),
		},
	)
	if err != nil {
		return asset.ScriptKey{}, keychain.KeyDescriptor{}, errors.Wrap(err, "awc.NextInternalKey")
	}

	internalKeyLnd, err := rpcutils.UnmarshalKeyDescriptor(
		internalKeyDesc.InternalKey,
	)
	if err != nil {
		return asset.ScriptKey{}, keychain.KeyDescriptor{}, errors.Wrap(err, "taprpc.UnmarshalKeyDescriptor")
	}

	return *_scriptKey, internalKeyLnd, nil
}

func maybeFundPacket(vPkg *tappsbt.VPacket) (*assetwalletrpc.FundVirtualPsbtResponse, error) {

	var buf bytes.Buffer
	err := vPkg.Serialize(&buf)
	if err != nil {
		return nil, errors.Wrap(err, "vPkg.serialize")
	}
	ctx := context.Background()
	ctxt, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	conn, clearUp, err := apiConnect.GetConnection(grpcTargetTapd, false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection "+grpcTargetTapd)
	}
	defer clearUp()
	awc := assetwalletrpc.NewAssetWalletClient(conn)

	return awc.FundVirtualPsbt(ctxt, &assetwalletrpc.FundVirtualPsbtRequest{
		Template: &assetwalletrpc.FundVirtualPsbtRequest_Psbt{Psbt: buf.Bytes()},
	})
}

func fundPacket2(vPkg *tappsbt.VPacket, coinsFilter *btlrpc.CoinsFilter) (*assetwalletrpc.FundVirtualPsbtResponse, error) {

	if coinsFilter != nil {
		if len(coinsFilter.Coins) > 1 {
			return nil, errors.New("coinsFilter.Coins only support one coin now")
		}
	}

	var buf bytes.Buffer
	err := vPkg.Serialize(&buf)
	if err != nil {
		return nil, errors.Wrap(err, "vPkg.serialize")
	}
	ctx := context.Background()
	ctxt, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	conn, clearUp, err := apiConnect.GetConnection(grpcTargetTapd, false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection "+grpcTargetTapd)
	}
	defer clearUp()

	bc := btlrpc.NewBtlClient(conn)

	return bc.FundVirtualPsbt2(ctxt, &btlrpc.FundVirtualPsbtRequest2{
		Template:    &btlrpc.FundVirtualPsbtRequest2_Psbt{Psbt: buf.Bytes()},
		CoinsFilter: coinsFilter,
	})
}

func requireLen(object any, length int) (err error) {
	l, ok := getLen(object)
	if !ok {
		err = errors.New("object does not have a length")
		return errors.Wrap(err, "getLen")
	}
	if l != length {
		err = fmt.Errorf("length mismatch(l: %d, length: %d)", l, length)
		return errors.Wrap(err, "l != length")
	}
	return nil
}

func requireLengths(object any, length int, lengths ...int) (err error) {
	l, ok := getLen(object)
	if !ok {
		err = errors.New("object does not have a length")
		return errors.Wrap(err, "getLen")
	}
	if l != length {
		for _, ls := range lengths {
			if l == ls {
				return nil
			}
			continue
		}
		err = fmt.Errorf("length mismatch(l: %d, length: %d, lengths: %v)", l, length, lengths)
		return errors.Wrap(err, "l != length or lengths")
	}
	return nil
}

func requireLenBe(object any, length int) (err error) {
	l, ok := getLen(object)
	if !ok {
		err = errors.New("object does not have a length")
		return errors.Wrap(err, "getLen")
	}
	if !(l >= length) {
		err = fmt.Errorf("l not bigger equal than length (l: %d, length: %d)", l, length)
		return errors.Wrap(err, "!(l >= length)")
	}
	return nil
}

func getLen(x any) (length int, ok bool) {
	v := reflect.ValueOf(x)
	defer func() {
		ok = recover() == nil
	}()
	return v.Len(), true
}

func requireEqual(expected, actual any) (err error) {
	if err = validateEqualArgs(expected, actual); err != nil {
		err = errors.New("invalid expected and actual arguments")
		return errors.Wrap(err, "validateEqualArgs")
	}
	if !objectsAreEqual(expected, actual) {
		err = errors.New("expected and actual arguments mismatch")
		return errors.Wrap(err, "objectsAreEqual")
	}
	return nil
}

func validateEqualArgs(expected, actual any) error {
	if expected == nil && actual == nil {
		return nil
	}

	if isFunction(expected) || isFunction(actual) {
		return errors.New("cannot take func type as argument")
	}
	return nil
}

func isFunction(arg any) bool {
	if arg == nil {
		return false
	}
	return reflect.TypeOf(arg).Kind() == reflect.Func
}

func objectsAreEqual(expected, actual any) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, actual)
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}
	return bytes.Equal(exp, act)
}

func requireNil(object any) (err error) {
	if !isNil(object) {
		err = errors.New("object is not nil")
		return errors.Wrap(err, "isNil")
	}
	return nil
}

func isNil(object any) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	switch value.Kind() {
	case
		reflect.Chan, reflect.Func,
		reflect.Interface, reflect.Map,
		reflect.Ptr, reflect.Slice, reflect.UnsafePointer:

		return value.IsNil()
	default:
		panic("unhandled default case")
	}

	return false
}

func requireContains(s, contains any) (err error) {
	ok, found := containsElement(s, contains)
	if !ok {
		err = errors.New("could not be applied builtin len()")
		return errors.Wrap(err, "containsElement")
	}
	if !found {
		err = errors.New("s does not contain elements")
		return errors.Wrap(err, "containsElement")
	}
	return nil
}

func containsElement(list any, element any) (ok, found bool) {

	listValue := reflect.ValueOf(list)
	listType := reflect.TypeOf(list)
	if listType == nil {
		return false, false
	}
	listKind := listType.Kind()
	defer func() {
		if e := recover(); e != nil {
			ok = false
			found = false
		}
	}()

	if listKind == reflect.String {
		elementValue := reflect.ValueOf(element)
		return true, strings.Contains(listValue.String(), elementValue.String())
	}

	if listKind == reflect.Map {
		mapKeys := listValue.MapKeys()
		for i := 0; i < len(mapKeys); i++ {
			if objectsAreEqual(mapKeys[i].Interface(), element) {
				return true, true
			}
		}
		return true, false
	}

	for i := 0; i < listValue.Len(); i++ {
		if objectsAreEqual(listValue.Index(i).Interface(), element) {
			return true, true
		}
	}
	return true, false

}

func getAddressBip32Derivation(addr string) (*psbt.Bip32Derivation, *psbt.TaprootBip32Derivation, error) {

	var (
		path        []uint32
		pubKeyBytes []byte
		err         error
	)

	ctx := context.Background()
	ctxt, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	conn, clearUp, err := apiConnect.GetConnection(grpcTargetLnd, false)
	if err != nil {
		return nil, nil, errors.Wrap(err, "apiConnect.GetConnection "+grpcTargetLnd)
	}
	defer clearUp()

	wkc := walletrpc.NewWalletKitClient(conn)

	addresses, err := wkc.ListAddresses(ctxt, &walletrpc.ListAddressesRequest{})
	if err != nil {
		return nil, nil, errors.Wrap(err, "wkc.ListAddresses")
	}

	for _, account := range addresses.AccountWithAddresses {
		for _, _address := range account.Addresses {
			if _address.Address == addr {
				path, err = lntest.ParseDerivationPath(
					_address.DerivationPath,
				)
				if err != nil {
					return nil, nil, errors.Wrap(err, "lntest.ParseDerivationPath")
				}
				pubKeyBytes = _address.PublicKey
			}
		}
	}

	if len(path) != 5 || len(pubKeyBytes) == 0 {
		err = fmt.Errorf("derivation path for address %s not found or invalid", addr)
		return nil, nil, errors.Wrap(err, "")
	}

	path[0] += hdkeychain.HardenedKeyStart
	path[1] += hdkeychain.HardenedKeyStart
	path[2] += hdkeychain.HardenedKeyStart

	return &psbt.Bip32Derivation{
			PubKey:    pubKeyBytes,
			Bip32Path: path,
		}, &psbt.TaprootBip32Derivation{
			XOnlyPubKey: pubKeyBytes[1:],
			Bip32Path:   path,
		}, nil
}

func finalizePacket(pkt *psbt.Packet) (*psbt.Packet, error) {
	var buf bytes.Buffer
	err := pkt.Serialize(&buf)
	if err != nil {
		return nil, errors.Wrap(err, "pkt.Serialize")
	}

	ctx := context.Background()
	ctxt, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	conn, clearUp, err := apiConnect.GetConnection(grpcTargetLnd, false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection "+grpcTargetLnd)
	}
	defer clearUp()
	wkc := walletrpc.NewWalletKitClient(conn)

	finalizeResp, err := wkc.FinalizePsbt(ctxt, &walletrpc.FinalizePsbtRequest{FundedPsbt: buf.Bytes()})
	if err != nil {
		return nil, errors.Wrap(err, "wkc.FinalizePsbt")
	}

	signedPacket, err := psbt.NewFromRawBytes(
		bytes.NewReader(finalizeResp.SignedPsbt), false,
	)
	if err != nil {
		return nil, errors.Wrap(err, "psbt.NewFromRawBytes")
	}

	return signedPacket, nil
}

func logAndPublish(btcPkt *psbt.Packet, activeAssets []*tappsbt.VPacket, passiveAssets []*tappsbt.VPacket, commitResp *assetwalletrpc.CommitVirtualPsbtsResponse) (*taprpc.SendAssetResponse, error) {

	ctxb := context.Background()
	ctxt, cancel := context.WithTimeout(ctxb, defaultTimeout)
	defer cancel()

	conn, clearUp, err := apiConnect.GetConnection(grpcTargetTapd, false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection "+grpcTargetTapd)
	}
	defer clearUp()
	awc := assetwalletrpc.NewAssetWalletClient(conn)

	var buf bytes.Buffer
	err = btcPkt.Serialize(&buf)
	if err != nil {
		return nil, errors.Wrap(err, "btcPkt.Serialize")
	}

	request := &assetwalletrpc.PublishAndLogRequest{
		AnchorPsbt:        buf.Bytes(),
		VirtualPsbts:      make([][]byte, len(activeAssets)),
		PassiveAssetPsbts: make([][]byte, len(passiveAssets)),
		ChangeOutputIndex: commitResp.ChangeOutputIndex,
		LndLockedUtxos:    commitResp.LndLockedUtxos,
	}

	for idx := range activeAssets {
		request.VirtualPsbts[idx], err = tappsbt.Encode(activeAssets[idx])
		if err != nil {
			return nil, errors.Wrap(err, "tappsbt.Encode")
		}
	}
	for idx := range passiveAssets {
		request.PassiveAssetPsbts[idx], err = tappsbt.Encode(passiveAssets[idx])
		if err != nil {
			return nil, errors.Wrap(err, "tappsbt.Encode")
		}
	}

	resp, err := awc.PublishAndLogTransfer(ctxt, request)
	if err != nil {
		return nil, errors.Wrap(err, "awc.PublishAndLogTransfer")
	}

	return resp, nil
}

func GetRawLastProof(scriptKey []byte, genInfo *taprpc.GenesisInfo, group *taprpc.AssetGroup, outpoint string) ([]byte, error) {
	proofFile, err := exportProofFileFromUniverse(genInfo.AssetId, scriptKey, outpoint, group)
	if err != nil {
		return nil, errors.Wrap(err, "exportProofFileFromUniverse")
	}

	lastProof, err := proofFile.RawLastProof()
	if err != nil {
		return nil, errors.Wrap(err, "proofFile.rawLastProof")
	}

	return lastProof, nil
}

func transferProofUniRPC(lastProof []byte) (*universerpc.AssetProofResponse, error) {
	importResp, err := insertProofIntoUniverse(lastProof)
	if err != nil {
		return nil, errors.Wrap(err, "insertProofIntoUniverse")
	}

	return importResp, nil
}

func exportProofFileFromUniverse(assetIDBytes, scriptKey []byte, outpoint string, group *taprpc.AssetGroup) (*proof.File, error) {

	ctxb := context.Background()
	ctxt, cancel := context.WithTimeout(ctxb, defaultTimeout)
	defer cancel()

	conn, clearUp, err := apiConnect.GetConnection(grpcTargetTapd, false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection "+grpcTargetTapd)
	}
	defer clearUp()

	unc := universerpc.NewUniverseClient(conn)

	var assetID asset.ID
	copy(assetID[:], assetIDBytes)

	scriptPubKey, err := btcec.ParsePubKey(scriptKey)
	if err != nil {
		return nil, errors.Wrap(err, "btcec.ParsePubKey")
	}

	op, err := wire.NewOutPointFromString(outpoint)
	if err != nil {
		return nil, errors.Wrap(err, "wire.NewOutPointFromString")
	}

	loc := proof.Locator{
		AssetID:   &assetID,
		ScriptKey: *scriptPubKey,
		OutPoint:  op,
	}

	if group != nil {
		groupKey, err := btcec.ParsePubKey(group.TweakedGroupKey)
		if err != nil {
			return nil, errors.Wrap(err, "btcec.ParsePubKey")
		}

		loc.GroupKey = groupKey
	}

	fetchUniProof := func(ctx context.Context,
		loc proof.Locator) (proof.Blob, error) {

		uniID := universe.Identifier{
			AssetID: *loc.AssetID,
		}
		if loc.GroupKey != nil {
			uniID.GroupKey = loc.GroupKey
		}

		rpcUniID, err := taprootassets.MarshalUniID(uniID)
		if err != nil {
			return nil, errors.Wrap(err, "taprootassets.MarshalUniID")
		}

		op := &universerpc.Outpoint{
			HashStr: loc.OutPoint.Hash.String(),
			Index:   int32(loc.OutPoint.Index),
		}
		scriptKeyBytes := loc.ScriptKey.SerializeCompressed()

		uniProof, err := unc.QueryProof(ctx, &universerpc.UniverseKey{
			Id: rpcUniID,
			LeafKey: &universerpc.AssetKey{
				Outpoint: &universerpc.AssetKey_Op{
					Op: op,
				},
				ScriptKey: &universerpc.AssetKey_ScriptKeyBytes{
					ScriptKeyBytes: scriptKeyBytes,
				},
			},
		})
		if err != nil {
			return nil, err
		}

		return uniProof.AssetLeaf.Proof, nil
	}

	var proofFile *proof.File
	err = wait.NoError(func() error {
		proofFile, err = proof.FetchProofProvenance(
			ctxt, nil, loc, fetchUniProof,
		)
		return err
	}, defaultTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "proof.FetchProofProvenance")
	}

	return proofFile, nil
}

func (f *_file) rawLastProof() ([]byte, error) {
	if err := f.isValid(); err != nil {
		return nil, err
	}

	return f.rawProofAt(uint32(len(f.proofs)) - 1)
}

type _file struct {
	Version _version

	proofs []*hashedProof
}
type hashedProof struct {
	proofBytes []byte

	hash [sha256.Size]byte
}

func (f *_file) isUnknownVersion() bool {
	switch f.Version {
	case _v0:
		return false
	default:
		return true
	}
}

type _version uint32

const (
	_v0 _version = 0
)

func (f *_file) isEmpty() bool {
	return len(f.proofs) == 0
}

func (f *_file) isValid() error {
	if f.isEmpty() {
		return errNoProofAvailable
	}

	if f.isUnknownVersion() {
		return errUnknownVersion
	}

	return nil
}

var (
	errNoProofAvailable = errors.New("no proof available")

	errUnknownVersion = errors.New("proof: unknown proof version")
)

func (f *_file) rawProofAt(index uint32) ([]byte, error) {
	if err := f.isValid(); err != nil {
		return nil, err
	}

	if index > uint32(len(f.proofs))-1 {
		return nil, fmt.Errorf("invalid index %d", index)
	}

	proofCopy := make([]byte, len(f.proofs[index].proofBytes))
	copy(proofCopy, f.proofs[index].proofBytes)

	return proofCopy, nil
}

func insertProofIntoUniverse(proofBytes proof.Blob) (*universerpc.AssetProofResponse, error) {

	ctxb := context.Background()
	ctxt, cancel := context.WithTimeout(ctxb, defaultTimeout)
	defer cancel()

	conn, clearUp, err := apiConnect.GetConnection(grpcTargetTapd, false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection "+grpcTargetTapd)
	}
	defer clearUp()

	tac := taprpc.NewTaprootAssetsClient(conn)

	resp, err := tac.DecodeProof(ctxt, &taprpc.DecodeProofRequest{
		RawProof:          proofBytes,
		WithMetaReveal:    true,
		WithPrevWitnesses: true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "tac.DecodeProof")
	}

	rpcProof := resp.DecodedProof
	rpcAsset := rpcProof.Asset
	rpcAnchor := rpcAsset.ChainAnchor

	uniID := universe.Identifier{
		ProofType: universe.ProofTypeTransfer,
	}
	if rpcProof.GenesisReveal != nil {
		uniID.ProofType = universe.ProofTypeIssuance
	}

	copy(uniID.AssetID[:], rpcAsset.AssetGenesis.AssetId)
	if rpcAsset.AssetGroup != nil {
		uniID.GroupKey, err = btcec.ParsePubKey(
			rpcAsset.AssetGroup.TweakedGroupKey,
		)
		if err != nil {
			return nil, errors.Wrap(err, "btcec.ParsePubKey")
		}
	}

	rpcUniID, err := taprootassets.MarshalUniID(uniID)
	if err != nil {
		return nil, errors.Wrap(err, "taprootassets.MarshalUniID")
	}

	uc := universerpc.NewUniverseClient(conn)

	importResp, err := uc.InsertProof(ctxt, &universerpc.AssetProof{
		Key: &universerpc.UniverseKey{
			Id: rpcUniID,
			LeafKey: &universerpc.AssetKey{
				Outpoint: &universerpc.AssetKey_OpStr{
					OpStr: rpcAnchor.AnchorOutpoint,
				},
				ScriptKey: &universerpc.AssetKey_ScriptKeyBytes{
					ScriptKeyBytes: rpcAsset.ScriptKey,
				},
			},
		},
		AssetLeaf: &universerpc.AssetLeaf{
			Proof: proofBytes,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "uc.InsertProof")
	}

	return importResp, nil
}

func getAssetGenInfoAndAssetGroup(id string) (genInfo *taprpc.GenesisInfo, assetGroup *taprpc.AssetGroup, err error) {
	universeHost := base.QueryConfigByKey("universeHost")
	var targets = func() []*universerpc.SyncTarget {
		universeID := &universerpc.ID{
			Id: &universerpc.ID_AssetIdStr{
				AssetIdStr: id,
			},
			ProofType: universerpc.ProofType_PROOF_TYPE_ISSUANCE,
		}
		var _targets []*universerpc.SyncTarget
		_targets = append(_targets, &universerpc.SyncTarget{
			Id: universeID,
		})
		return _targets
	}()
	_, err = syncUniverse(universeHost, targets, universerpc.UniverseSyncMode_SYNC_FULL)
	if err != nil {
		LogError("syncUniverse ISSUANCE FULL", err)
	}
	{
	}
	root, err := queryAssetRoot(id)
	if err != nil {
		return nil, nil, errors.Wrap(err, "queryAssetRoot")
	}
	if root == nil || root.IssuanceRoot.Id == nil {
		return nil, nil, errors.New("query asset roots err")
	}
	queryId := id
	isGroup := false
	if groupKey, ok := root.IssuanceRoot.Id.Id.(*universerpc.ID_GroupKey); ok {
		isGroup = true
		queryId = hex.EncodeToString(groupKey.GroupKey)
	}
	response, err := assetLeaves(isGroup, queryId, universerpc.ProofType_PROOF_TYPE_ISSUANCE)
	if err != nil {
		return nil, nil, err
	}
	if response.Leaves == nil {
		return nil, nil, errors.New("response leaves null err")
	}

	for _, leaf := range response.Leaves {
		if hex.EncodeToString(leaf.Asset.AssetGenesis.GetAssetId()) == id {
			genInfo = leaf.Asset.AssetGenesis
			assetGroup = leaf.Asset.AssetGroup
			return genInfo, assetGroup, nil
		}
	}
	err = errors.New("query asset leaves not found")
	return nil, nil, err
}

func psbtPacketToBytes(psbt *psbt.Packet) ([]byte, error) {
	var b bytes.Buffer
	err := psbt.Serialize(&b)
	if err != nil {
		return nil, errors.Wrap(err, "psbt.Serialize")
	}
	return b.Bytes(), nil
}

func bytesToPsbtPacket(psbtBytes []byte) (*psbt.Packet, error) {
	return psbt.NewFromRawBytes(bytes.NewReader(psbtBytes), false)
}

func _outpointToTransactionAndIndex(outpoint string) (transaction string, index string, err error) {
	result := strings.Split(outpoint, ":")
	if len(result) < 2 {
		return "", "", errors.New("outpoint is invalid")
	}
	return result[0], result[1], nil
}
