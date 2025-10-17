package itest

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/txscript"
	"github.com/davecgh/go-spew/spew"
	"github.com/lightninglabs/taproot-assets/address"
	"github.com/lightninglabs/taproot-assets/asset"
	"github.com/lightninglabs/taproot-assets/proof"
	"github.com/lightninglabs/taproot-assets/tappsbt"
	"github.com/lightninglabs/taproot-assets/taprpc"
	wrpc "github.com/lightninglabs/taproot-assets/taprpc/assetwalletrpc"
	"github.com/lightninglabs/taproot-assets/taprpc/mintrpc"
	"github.com/lightninglabs/taproot-assets/tapsend"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/walletrpc"
	"github.com/lightningnetwork/lnd/lntest"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
	"time"
)

var testCases = []*testCase{
	{
		name: "get info",
		test: testGetInfo2,
	},
	{
		name: "psbt trustless swap",
		test: testPsbtTrustlessSwap2,
	},
}

var optionalTestCases = []*testCase{}

var optionalTests = flag.Bool("optional", false, "if true, the optional test list will be used")

func TestTaprootAssetsDaemon2(t *testing.T) {
	testList := testCases
	if *optionalTests {
		testList = optionalTestCases
	}

	if len(testList) == 0 {
		t.Skip("integration tests not selected with flag 'itest'")
	}

	ht := &harnessTest{t: t}
	ht.setupLogging()

	feeService := lntest.NewFeeService(t)
	lndHarness := lntest.SetupHarness(t, "./lnd-itest.exe", "bbolt", true, feeService)
	defer func() {
		time.Sleep(100 * time.Millisecond)
		lndHarness.Stop()
	}()

	lndHarness.SetupStandbyNodes()

	t.Log("Starting universe server LND node")
	uniServerLndHarness := lndHarness.NewNode("uni-server-lnd", nil)

	lndHarness.WaitForBlockchainSync(uniServerLndHarness)
	t.Logf("Running %v integration tests", len(testList))
	for _, _testCase := range testList {
		logLine := fmt.Sprintf("STARTING ============ %v ============\n", _testCase.name)

		success := t.Run(_testCase.name, func(t1 *testing.T) {
			_tapdHarness, uniHarness, proofCourier := setupHarnesses(t1, ht, lndHarness, uniServerLndHarness, _testCase.proofCourierType)
			lndHarness.EnsureConnected(lndHarness.Alice, lndHarness.Bob)
			lndHarness.EnsureConnected(lndHarness.Alice, uniServerLndHarness)
			lndHarness.EnsureConnected(lndHarness.Bob, uniServerLndHarness)
			lndHarness.Alice.AddToLogf(logLine)
			lndHarness.Bob.AddToLogf(logLine)

			_ht := ht.newHarnessTest(t1, lndHarness, uniHarness, _tapdHarness, proofCourier)

			_ht.RunTestCase(_testCase)

			err := _ht.shutdown(t1)
			require.NoError(t1, err)
		})

		if !success {
			return
		}
	}
}

func testGetInfo2(t *harnessTest) {
	ctxb := context.Background()
	ctxt, cancel := context.WithTimeout(ctxb, defaultWaitTimeout)
	defer cancel()

	resp, err := t.tapd.GetInfo(ctxt, &taprpc.GetInfoRequest{})
	require.NoError(t.t, err)

	expectedNetwork := t.tapd.cfg.NetParams.Name
	require.Equal(t.t, expectedNetwork, resp.Network)

	respGeneric, err := ExecTapCLI(ctxt, t.tapd, "getinfo")
	require.NoError(t.t, err)

	respCli := respGeneric.(*taprpc.GetInfoResponse)

	require.Equal(t.t, resp, respCli)
}

func testPsbtTrustlessSwap2(t *harnessTest) {

	rpcAssets := MintAssetsConfirmBatch(t.t, t.lndHarness.Miner().Client, t.tapd, []*mintrpc.MintAssetRequest{issuableAssets[0]})

	mintedAsset := rpcAssets[0]
	genInfo := mintedAsset.AssetGenesis
	ctxb := context.Background()

	var (
		aliceTapd    = t.tapd
		numUnits     = mintedAsset.Amount
		_chainParams = &address.RegressionNetTap
		assetID      asset.ID
	)
	copy(assetID[:], genInfo.AssetId)

	aliceDummyScriptKey, aliceAnchorInternalKey := DeriveKeys(t.t, aliceTapd)
	vPkt := tappsbt.ForInteractiveSend(
		assetID, numUnits, aliceDummyScriptKey, 0, 0, 1,
		aliceAnchorInternalKey, asset.V0, _chainParams,
	)

	fundResp := fundPacket(t, aliceTapd, vPkt)

	var err error
	vPkt, err = tappsbt.Decode(fundResp.FundedPsbt)
	require.NoError(t.t, err)

	require.Len(t.t, vPkt.Inputs, 1)
	require.Len(t.t, vPkt.Outputs, 1)

	vPkt.Inputs[0].SighashType = txscript.SigHashNone

	require.Equal(t.t, vPkt.Outputs[0].Type, tappsbt.TypeSimple)
	require.NoError(t.t, tapsend.PrepareOutputAssets(ctxb, vPkt))
	require.Nil(t.t, vPkt.Outputs[0].Asset.SplitCommitmentRoot)
	require.Len(t.t, vPkt.Outputs[0].Asset.PrevWitnesses, 1)
	require.Nil(t.t, vPkt.Outputs[0].Asset.PrevWitnesses[0].SplitCommitment)

	fundedPsbtBytes, err := tappsbt.Encode(vPkt)
	require.NoError(t.t, err)

	signedResp, err := aliceTapd.SignVirtualPsbt(
		ctxb, &wrpc.SignVirtualPsbtRequest{
			FundedPsbt: fundedPsbtBytes,
		},
	)
	require.NoError(t.t, err)
	require.Contains(t.t, signedResp.SignedInputs, uint32(0))

	vPkt, err = tappsbt.Decode(signedResp.SignedPsbt)
	require.NoError(t.t, err)

	btcPacket, err := tapsend.PrepareAnchoringTemplate([]*tappsbt.VPacket{
		vPkt,
	})
	require.NoError(t.t, err)

	require.Len(t.t, btcPacket.Inputs, 1)
	require.Len(t.t, btcPacket.Outputs, 2)

	addrResp := t.lndHarness.Alice.RPC.NewAddress(&lnrpc.NewAddressRequest{
		Type: lnrpc.AddressType_TAPROOT_PUBKEY,
	})

	aliceP2TR, err := btcutil.DecodeAddress(
		addrResp.Address, harnessNetParams,
	)
	require.NoError(t.t, err)

	alicePkScript, err := txscript.PayToAddrScript(aliceP2TR)
	require.NoError(t.t, err)

	btcPacket.UnsignedTx.TxOut[0].PkScript = alicePkScript
	btcPacket.UnsignedTx.TxOut[0].Value = 69420
	derivation, trDerivation := getAddressBip32Derivation(
		t.t, addrResp.Address, t.lndHarness.Alice,
	)

	btcPacket.Outputs[0].Bip32Derivation = []*psbt.Bip32Derivation{
		derivation,
	}
	btcPacket.Outputs[0].TaprootBip32Derivation =
		[]*psbt.TaprootBip32Derivation{trDerivation}
	btcPacket.Outputs[0].TaprootInternalKey = trDerivation.XOnlyPubKey

	var b bytes.Buffer
	err = btcPacket.Serialize(&b)
	require.NoError(t.t, err)

	resp, err := aliceTapd.CommitVirtualPsbts(
		ctxb, &wrpc.CommitVirtualPsbtsRequest{
			VirtualPsbts: [][]byte{signedResp.SignedPsbt},
			AnchorPsbt:   b.Bytes(),
			AnchorChangeOutput: &wrpc.CommitVirtualPsbtsRequest_Add{
				Add: true,
			},
			Fees: &wrpc.CommitVirtualPsbtsRequest_TargetConf{
				TargetConf: 12,
			},
		},
	)
	require.NoError(t.t, err)

	btcPacket, err = psbt.NewFromRawBytes(
		bytes.NewReader(resp.AnchorPsbt), false,
	)
	require.NoError(t.t, err)

	btcPacket.Inputs[0].SighashType = txscript.SigHashSingle |
		txscript.SigHashAnyOneCanPay

	btcPacket.Inputs = append(
		btcPacket.Inputs[:1], btcPacket.Inputs[2:]...,
	)
	btcPacket.UnsignedTx.TxIn = append(
		btcPacket.UnsignedTx.TxIn[:1], btcPacket.UnsignedTx.TxIn[2:]...,
	)

	btcPacket.Outputs = btcPacket.Outputs[:2]
	btcPacket.UnsignedTx.TxOut = btcPacket.UnsignedTx.TxOut[:2]

	t.Logf("Alice BTC PSBT: %v", spew.Sdump(btcPacket))

	b.Reset()
	err = btcPacket.Serialize(&b)
	require.NoError(t.t, err)

	signPsbtResp := t.lndHarness.Alice.RPC.SignPsbt(
		&walletrpc.SignPsbtRequest{
			FundedPsbt: b.Bytes(),
		},
	)

	require.Len(t.t, signPsbtResp.SignedInputs, 1)
	require.Equal(t.t, uint32(0), signPsbtResp.SignedInputs[0])

	btcPacket, err = psbt.NewFromRawBytes(
		bytes.NewReader(signPsbtResp.SignedPsbt), false,
	)
	require.NoError(t.t, err)

	require.Len(t.t, btcPacket.Inputs, 1)
	require.Len(t.t, btcPacket.Outputs, 2)

	signedVpsbtBytes, err := tappsbt.Encode(vPkt)
	require.NoError(t.t, err)

	secondTapd := setupTapdHarness(
		t.t, t, t.lndHarness.Bob, t.universeServer,
	)
	defer func() {
		require.NoError(t.t, secondTapd.stop(!*noDelete))
	}()

	var bob = secondTapd

	bobVPsbt, err := tappsbt.Decode(signedVpsbtBytes)
	require.NoError(t.t, err)

	require.Len(t.t, bobVPsbt.Outputs, 1)

	bobScriptKey, bobAnchorInternalKey := DeriveKeys(t.t, bob)

	bobVOut := bobVPsbt.Outputs[0]
	bobVOut.ScriptKey = bobScriptKey
	bobVOut.AnchorOutputBip32Derivation = nil
	bobVOut.AnchorOutputTaprootBip32Derivation = nil
	bobVOut.SetAnchorInternalKey(
		bobAnchorInternalKey, harnessNetParams.HDCoinType,
	)
	deliveryAddrStr := fmt.Sprintf(
		"%s://%s", proof.UniverseRpcCourierType,
		t.universeServer.ListenAddr,
	)
	deliveryAddr, err := url.Parse(deliveryAddrStr)
	require.NoError(t.t, err)
	bobVPsbt.Outputs[0].ProofDeliveryAddress = deliveryAddr

	btcPacket.Outputs[1].TaprootInternalKey = schnorr.SerializePubKey(
		bobAnchorInternalKey.PubKey,
	)
	btcPacket.Outputs[1].Bip32Derivation =
		bobVOut.AnchorOutputBip32Derivation
	btcPacket.Outputs[1].TaprootBip32Derivation =
		bobVOut.AnchorOutputTaprootBip32Derivation

	witnessBackup := bobVPsbt.Outputs[0].Asset.PrevWitnesses

	err = tapsend.PrepareOutputAssets(ctxb, bobVPsbt)
	require.NoError(t.t, err)

	require.Len(t.t, bobVPsbt.Outputs, 1)
	require.Equal(
		t.t, bobVPsbt.Outputs[0].ScriptKey,
		bobVPsbt.Outputs[0].Asset.ScriptKey,
	)

	bobVPsbt.Outputs[0].Asset.PrevWitnesses = witnessBackup

	bobVPsbtBytes, err := tappsbt.Encode(bobVPsbt)
	require.NoError(t.t, err)

	b.Reset()
	err = btcPacket.Serialize(&b)
	require.NoError(t.t, err)

	resp, err = bob.CommitVirtualPsbts(
		ctxb, &wrpc.CommitVirtualPsbtsRequest{
			VirtualPsbts: [][]byte{bobVPsbtBytes},
			AnchorPsbt:   b.Bytes(),
			AnchorChangeOutput: &wrpc.CommitVirtualPsbtsRequest_Add{
				Add: true,
			},
			Fees: &wrpc.CommitVirtualPsbtsRequest_TargetConf{
				TargetConf: 12,
			},
		},
	)
	require.NoError(t.t, err)

	bobVPsbt, err = tappsbt.Decode(resp.VirtualPsbts[0])
	require.NoError(t.t, err)

	signResp := t.lndHarness.Bob.RPC.SignPsbt(
		&walletrpc.SignPsbtRequest{
			FundedPsbt: resp.AnchorPsbt,
		},
	)
	require.NoError(t.t, err)

	finalPsbt, err := psbt.NewFromRawBytes(
		bytes.NewReader(signResp.SignedPsbt), false,
	)
	require.NoError(t.t, err)

	require.Len(t.t, finalPsbt.Inputs, 2)

	bobInputIdx := uint32(1)

	require.Len(t.t, signResp.SignedInputs, 1)

	require.Equal(t.t, bobInputIdx, signResp.SignedInputs[0])
	require.NoError(t.t, finalPsbt.SanityCheck())

	signedPkt := finalizePacket(t.t, t.lndHarness.Bob, finalPsbt)
	require.True(t.t, signedPkt.IsComplete())

	logResp := logAndPublish(
		t.t, aliceTapd, signedPkt, []*tappsbt.VPacket{bobVPsbt}, nil, resp,
	)
	t.Logf("Logged transaction: %v", toJSON(t.t, logResp))

	MineBlocks(t.t, t.lndHarness.Miner().Client, 1, 1)

	bobScriptKeyBytes := bobScriptKey.PubKey.SerializeCompressed()
	bobOutputIndex := uint32(1)
	transferTXID := finalPsbt.UnsignedTx.TxHash()
	bobAssetOutpoint := fmt.Sprintf("%s:%d", transferTXID.String(),
		bobOutputIndex)
	transferProofUniRPC(
		t, t.universeServer.service, bob, bobScriptKeyBytes, genInfo,
		mintedAsset.AssetGroup, bobAssetOutpoint,
	)

	registerResp, err := bob.RegisterTransfer(
		ctxb, &taprpc.RegisterTransferRequest{
			AssetId:   assetID[:],
			GroupKey:  mintedAsset.AssetGroup.TweakedGroupKey,
			ScriptKey: bobScriptKeyBytes,
			Outpoint: &taprpc.OutPoint{
				Txid:        transferTXID[:],
				OutputIndex: bobOutputIndex,
			},
		},
	)
	require.NoError(t.t, err)
	require.Equal(
		t.t, bobScriptKeyBytes, registerResp.RegisteredAsset.ScriptKey,
	)

	bobAssets, err := bob.ListAssets(ctxb, &taprpc.ListAssetRequest{})
	require.NoError(t.t, err)

	require.Len(t.t, bobAssets.Assets, 1)
	require.Equal(t.t, bobAssets.Assets[0].Amount, numUnits)

	require.Equal(t.t, bobScriptKeyBytes, bobAssets.Assets[0].ScriptKey)
}
