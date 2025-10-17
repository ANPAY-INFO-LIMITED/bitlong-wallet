package tapdlock

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcec/v2/schnorr/musig2"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/lightninglabs/lndclient"
	"github.com/lightninglabs/taproot-assets/asset"
	"github.com/lightninglabs/taproot-assets/commitment"
	"github.com/lightninglabs/taproot-assets/fn"
	"github.com/lightninglabs/taproot-assets/rpcutils"
	"github.com/lightninglabs/taproot-assets/tappsbt"
	"github.com/lightninglabs/taproot-assets/taprpc"
	wrpc "github.com/lightninglabs/taproot-assets/taprpc/assetwalletrpc"
	"github.com/lightninglabs/taproot-assets/taprpc/tapdevrpc"
	unirpc "github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	"github.com/lightninglabs/taproot-assets/tapsend"
	"github.com/lightningnetwork/lnd/input"
	"github.com/lightningnetwork/lnd/keychain"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/chainrpc"
	"github.com/lightningnetwork/lnd/lnrpc/signrpc"
	"github.com/lightningnetwork/lnd/lnrpc/walletrpc"
	"github.com/lightningnetwork/lnd/lntest/wait"
	"github.com/lightningnetwork/lnd/lnwallet/chainfee"
	"github.com/wallet/service/apiConnect"
	"github.com/wallet/types"
)

type TapdClient struct {
	TaprootAssetsClient taprpc.TaprootAssetsClient
	AssetWalletClient   wrpc.AssetWalletClient
	ChainKitClient      chainrpc.ChainKitClient
	UniverseClient      unirpc.UniverseClient
}

type LndClient struct {
	LN        lnrpc.LightningClient
	Signer    signrpc.SignerClient
	WalletKit walletrpc.WalletKitClient
}

func BobGenerateAddr(ctx context.Context, req *types.BobGenerateAddrReq) (resp *types.BobGenerateAddrResp, err error) {

	conn2, clearUp2, err := GetTaprootAssetsClient()
	if err != nil {
		return nil, err
	}
	defer clearUp2()

	bobScriptKey, bobInternalKey, err := DeriveKeys(ctx)
	if err != nil {
		return nil, err
	}

	withdrawAddr, err := conn2.NewAddr(ctx, &taprpc.NewAddrRequest{
		AssetId: req.AssetId,
		Amt:     req.Amount,
	})
	if err != nil {
		return nil, err
	}

	return &types.BobGenerateAddrResp{
		WithdrawAddr:   withdrawAddr,
		BobScriptKey:   bobScriptKey.ScriptKey,
		BobInternalKey: bobInternalKey.InternalKey,
	}, nil
}

func AliceSetMusigLock(ctx context.Context, req *types.AliceSetBizReq) (resp *types.AliceSetBizResp, err error) {
	conn, clearUp, err := GetTaprootAssetsClient()
	if err != nil {
		return nil, err
	}
	defer clearUp()

	conn2, clearUp2, err := getChainKitClient()
	if err != nil {
		return nil, err
	}
	defer clearUp2()

	blockInfo, err := conn2.GetBestBlock(ctx, &chainrpc.GetBestBlockRequest{})
	if err != nil {
		return nil, err
	}

	bestBlock := blockInfo.GetBlockHeight()
	lockTimeBlocks := int64(bestBlock) + req.LockTime

	aliceScriptKey, aliceInternalKey, err := DeriveKeys(ctx)
	if err != nil {
		return nil, err
	}

	scriptKey, scriptKey1, err := unmarshalScriptKeys(req.BobScriptKey, aliceScriptKey.ScriptKey)
	if err != nil {
		return nil, err
	}

	bobInternalKey, err := rpcutils.UnmarshalKeyDescriptor(req.BobInternalKey)
	if err != nil {
		return nil, err
	}

	aliceInternalKey1, err := rpcutils.UnmarshalKeyDescriptor(aliceInternalKey.InternalKey)
	if err != nil {
		return nil, err
	}

	btcTapscript, err := txscript.NewScriptBuilder().
		AddData(schnorr.SerializePubKey(aliceInternalKey1.PubKey)). // 添加 Alice 的公钥
		AddOp(txscript.OP_CHECKSIG).                                // Alice 的签名验证
		AddData(schnorr.SerializePubKey(bobInternalKey.PubKey)).    //添加 Bob 的公钥
		AddOp(txscript.OP_CHECKSIGADD).                             // Bob 的签名验证
		AddInt64(2).                                                // 最少需要 2 个签名
		AddOp(txscript.OP_EQUAL).                                   // 比较是否满足签名要求
		AddInt64(lockTimeBlocks).                                   // 锁定的区块高度或时间
		AddOp(txscript.OP_CHECKLOCKTIMEVERIFY).                     // 时间锁验证
		AddOp(txscript.OP_DROP).                                    // 丢弃栈顶的时间锁值
		Script()

	if err != nil {
		return nil, err
	}

	btcTapLeaf := txscript.TapLeaf{
		LeafVersion: txscript.BaseLeafVersion,
		Script:      btcTapscript,
	}

	btcInternalKey := asset.NUMSPubKey
	btcControlBlock := &txscript.ControlBlock{
		LeafVersion: txscript.BaseLeafVersion,
		InternalKey: btcInternalKey,
	}

	siblingPreimage, err := commitment.NewPreimageFromLeaf(btcTapLeaf)
	if err != nil {
		return nil, err
	}

	siblingPreimageBytes, _, err := commitment.MaybeEncodeTapscriptPreimage(
		siblingPreimage,
	)
	if err != nil {
		return nil, err
	}

	aliceNonces, bobNonces, err := generateNonces(scriptKey1.RawKey.PubKey, scriptKey.RawKey.PubKey)
	if err != nil {
		return nil, err
	}

	muSig2Key, err := input.MuSig2CombineKeys(
		input.MuSig2Version100RC2, []*btcec.PublicKey{
			scriptKey1.RawKey.PubKey,
			scriptKey.RawKey.PubKey,
		}, true, &input.MuSig2Tweaks{TaprootBIP0086Tweak: true},
	)
	if err != nil {
		return nil, err
	}

	tapScriptKey, tapLeaves, _, tapControlBlock, err := createMuSigLeaves(
		muSig2Key,
	)
	if err != nil {
		return nil, err
	}

	tapControlBlockBytes, err := tapControlBlock.ToBytes()
	if err != nil {
		return nil, err
	}

	muSig2Addr, err := conn.NewAddr(ctx, &taprpc.NewAddrRequest{
		AssetId:   req.AssetIdBytes,
		Amt:       req.Amount,
		ScriptKey: rpcutils.MarshalScriptKey(tapScriptKey),
		InternalKey: &taprpc.KeyDescriptor{
			RawKeyBytes: pubKeyBytes(btcInternalKey),
		},
		TapscriptSibling: siblingPreimageBytes,
	})
	if err != nil {
		return nil, err
	}

	sendResp, err := conn.SendAsset(ctx, &taprpc.SendAssetRequest{
		TapAddrs: []string{muSig2Addr.Encoded},
	})
	if err != nil {
		return nil, err
	}

	multiSigOutAnchor := sendResp.Transfer.Outputs[1].Anchor
	btcControlBlock.InclusionProof = multiSigOutAnchor.TaprootAssetRoot

	rootHash := btcControlBlock.RootHash(btcTapscript)
	tapKey := txscript.ComputeTaprootOutputKey(btcInternalKey, rootHash)
	if tapKey.SerializeCompressed()[0] == secp256k1.PubKeyFormatCompressedOdd {
		btcControlBlock.OutputKeyYIsOdd = true
	}

	btcControlBlockBytes, err := btcControlBlock.ToBytes()
	if err != nil {
		return nil, err
	}

	return &types.AliceSetBizResp{
		Leaves:               tapLeaves,
		TapControlBlockBytes: tapControlBlockBytes,
		BtcControlBlockBytes: btcControlBlockBytes,
		AliceScriptKey:       aliceScriptKey.ScriptKey,
		AliceInternalKey:     aliceInternalKey.InternalKey,
		BtcTapLeaf:           btcTapLeaf,
		AliceNonces:          aliceNonces,
		BobNonces:            bobNonces,
		LockTime:             lockTimeBlocks,
	}, nil
}

func FundVirtualPsbt(ctx context.Context, req *types.FundVirtualPsbtReq) (resp *types.FundVirtualPsbtResp, err error) {
	conn, clearUp, err := GetAssetWalletClient()
	if err != nil {
		return nil, err
	}
	defer clearUp()

	withdrawRecipients := map[string]uint64{
		req.WithdrawAddr.Encoded: req.WithdrawAddr.Amount,
	}

	withdrawFundResp, err := conn.FundVirtualPsbt(
		ctx, &wrpc.FundVirtualPsbtRequest{
			Template: &wrpc.FundVirtualPsbtRequest_Raw{
				Raw: &wrpc.TxTemplate{
					Recipients: withdrawRecipients,
				},
			},
		},
	)

	if err != nil {
		return nil, err
	}

	return &types.FundVirtualPsbtResp{
		FundedPsbt: withdrawFundResp.FundedPsbt,
	}, nil
}

func BobPartialSig(ctx context.Context, req *types.BobPartialSigReq) (resp *types.BobPartialSigResp, err error) {
	fundedWithdrawPkt := deserializeVPacket(
		req.FundedPsbt,
	)

	bobScriptKey, aliceScriptKey, err := unmarshalScriptKeys(req.BobScriptKey, req.AliceScriptKey)
	if err != nil {
		return nil, err
	}

	leafToSign := req.Leaves[0]
	bobPartialSig, _, err := tapCreatePartialSig(
		ctx, &chaincfg.RegressionNetParams, fundedWithdrawPkt, leafToSign,
		bobScriptKey.RawKey, req.BobNonces, aliceScriptKey.RawKey.PubKey,
		req.AliceNonces.PubNonce,
	)

	if err != nil {
		return nil, err
	}

	return &types.BobPartialSigResp{
		BobPartialSig: &bobPartialSig,
	}, nil
}

func SubmitVirtualTransaction(ctx context.Context, req *types.SubmitVtxReq) (resp *types.SubmitVtxResp, err error) {
	bobScriptKey, aliceScriptKey, err := unmarshalScriptKeys(req.BobScriptKey, req.AliceScriptKey)
	if err != nil {
		return nil, err
	}

	fundedWithdrawPkt := deserializeVPacket(
		req.FundedPsbt,
	)

	tapControlBlock, err := txscript.ParseControlBlock(req.TapControlBlockBytes)
	if err != nil {
		return nil, err
	}

	leafToSign := req.Leaves[0]
	_, aliceSessID, err := tapCreatePartialSig(
		ctx, &chaincfg.RegressionNetParams, fundedWithdrawPkt, leafToSign,
		aliceScriptKey.RawKey, req.AliceNonces, bobScriptKey.RawKey.PubKey,
		req.BobNonces.PubNonce,
	)
	if err != nil {
		return nil, err
	}

	tree := txscript.AssembleTaprootScriptTree(req.Leaves...)

	finalTapWitness, err := combineSigs(
		ctx, aliceSessID, req.BobPartialSig, leafToSign, tree,
		tapControlBlock,
	)

	if err != nil {
		return nil, err
	}

	for idx := range fundedWithdrawPkt.Outputs {
		updateWitness(
			fundedWithdrawPkt.Outputs[idx].Asset, finalTapWitness,
		)
	}

	vPackets := []*tappsbt.VPacket{fundedWithdrawPkt}
	withdrawBtcPkt, err := tapsend.PrepareAnchoringTemplate(vPackets)
	if err != nil {
		return nil, err
	}

	btcWithdrawPkt, finalizedWithdrawPackets, _, commitResp := CommitVirtualPsbts(
		ctx, withdrawBtcPkt, vPackets, nil, -1,
	)

	btcWithdrawPkt.UnsignedTx.LockTime = uint32(req.LockTime)

	btcWithdrawPktStr, err := btcWithdrawPkt.B64Encode()
	if err != nil {
		return nil, err
	}

	finalizedWithdrawPacketsBytes, err := EncodeVPackets(finalizedWithdrawPackets)
	if err != nil {
		return nil, err
	}

	return &types.SubmitVtxResp{
		BtcWithdrawPkt:           btcWithdrawPktStr,
		FinalizedWithdrawPackets: finalizedWithdrawPacketsBytes,
		CommitResp:               commitResp,
	}, nil
}

func BobBtcPartialSig(ctx context.Context, req *types.BtcPartialSigReq) (resp *types.BtcPartialSigResp, err error) {
	bobInternalKey, err := rpcutils.UnmarshalKeyDescriptor(req.BobInternalKey)
	if err != nil {
		return nil, err
	}

	btcWithdrawPkt, err := DecodeBase64ToPSBT(req.BtcWithdrawPkt)
	if err != nil {
		return nil, err
	}

	assetInputIdx := uint32(0)
	bobBtcPartialSig := partialSignWithKey(
		ctx, &chaincfg.RegressionNetParams, btcWithdrawPkt, assetInputIdx,
		bobInternalKey, req.BtcControlBlockBytes, req.BtcTapLeaf,
	)
	return &types.BtcPartialSigResp{
		BobBtcPartialSig: bobBtcPartialSig,
	}, nil
}

func SignAndFinalizeBtcTransaction(ctx context.Context, req *types.SignAndFinalizeBtcTransactionReq) error {
	aliceInternalKey, err := rpcutils.UnmarshalKeyDescriptor(req.AliceInternalKey)
	if err != nil {
		return err
	}

	btcWithdrawPkt, err := DecodeBase64ToPSBT(req.BtcWithdrawPkt)
	if err != nil {
		return err
	}

	assetInputIdx := uint32(0)
	aliceBtcPartialSig := partialSignWithKey(
		ctx, &chaincfg.RegressionNetParams, btcWithdrawPkt, assetInputIdx,
		aliceInternalKey, req.BtcControlBlockBytes, req.BtcTapLeaf,
	)

	txWitness := wire.TxWitness{
		req.BobBtcPartialSig,
		aliceBtcPartialSig,
		req.BtcTapLeaf.Script,
		req.BtcControlBlockBytes,
	}

	var buf bytes.Buffer
	err = psbt.WriteTxWitness(&buf, txWitness)
	if err != nil {
		return err
	}

	btcWithdrawPkt.Inputs[assetInputIdx].FinalScriptWitness = buf.Bytes()
	finalizedWithdrawPackets, err := DecodeVPackets(req.FinalizedWithdrawPackets)
	if err != nil {
		return err
	}

	signedPkt, err := FinalizePacket(ctx, btcWithdrawPkt)
	if err != nil {
		return err
	}

	logResp, err := LogAndPublish(
		ctx, signedPkt, finalizedWithdrawPackets, nil, req.CommitResp,
	)

	if err != nil {
		return err
	}

	fmt.Printf("交易已发布: %s\n", logResp)
	return nil
}

func GetTapdClient() (TapdClient, error) {
	conn1, clearUp1, err := GetTaprootAssetsClient()
	if err != nil {
		return TapdClient{}, err
	}
	defer clearUp1()

	conn2, clearUp2, err := GetAssetWalletClient()
	if err != nil {
		return TapdClient{}, err
	}
	defer clearUp2()

	client := TapdClient{
		TaprootAssetsClient: conn1,
		AssetWalletClient:   conn2,
	}
	return client, nil
}

func GetLndClient() (LndClient, error) {
	conn1, clearUp1, err := getLightningClient()
	if err != nil {
		return LndClient{}, err
	}
	defer clearUp1()

	conn2, clearUp2, err := getSignerClient()
	if err != nil {
		return LndClient{}, err
	}
	defer clearUp2()

	conn3, clearUp3, err := getWalletKitClient()
	if err != nil {
		return LndClient{}, err
	}
	defer clearUp3()

	client := LndClient{
		LN:        conn1,
		Signer:    conn2,
		WalletKit: conn3,
	}
	return client, nil
}

func getChainKitClient() (chainrpc.ChainKitClient, func(), error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	client := chainrpc.NewChainKitClient(conn)
	return client, clearUp, nil
}

func getWalletKitClient() (walletrpc.WalletKitClient, func(), error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	client := walletrpc.NewWalletKitClient(conn)
	return client, clearUp, nil
}

func getLightningClient() (lnrpc.LightningClient, func(), error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	client := lnrpc.NewLightningClient(conn)
	return client, clearUp, nil
}

func getSignerClient() (signrpc.SignerClient, func(), error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	client := signrpc.NewSignerClient(conn)
	return client, clearUp, nil
}

func GetTaprootAssetsClient() (taprpc.TaprootAssetsClient, func(), error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	client := taprpc.NewTaprootAssetsClient(conn)
	return client, clearUp, nil
}

func GetAssetWalletClient() (wrpc.AssetWalletClient, func(), error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	client := wrpc.NewAssetWalletClient(conn)
	return client, clearUp, nil
}

func getTapDevClient() (tapdevrpc.TapDevClient, func(), error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	client := tapdevrpc.NewTapDevClient(conn)
	return client, clearUp, nil
}

func unmarshalScriptKeys(bobKey, aliceKey *taprpc.ScriptKey) (bobScriptKey, aliceScriptKey *asset.ScriptKey, err error) {
	bobScriptKey, err = rpcutils.UnmarshalScriptKey(bobKey)
	if err != nil {
		return nil, nil, err
	}

	aliceScriptKey, err = rpcutils.UnmarshalScriptKey(aliceKey)
	if err != nil {
		return nil, nil, err
	}

	return bobScriptKey, aliceScriptKey, nil
}

func LogAndPublish(ctx context.Context, btcPkt *psbt.Packet,
	activeAssets []*tappsbt.VPacket, passiveAssets []*tappsbt.VPacket,
	commitResp *wrpc.CommitVirtualPsbtsResponse) (*taprpc.SendAssetResponse, error) {
	tapClient, clearUp, err := GetAssetWalletClient()
	if err != nil {
		return nil, err
	}
	defer clearUp()

	var buf bytes.Buffer
	err = btcPkt.Serialize(&buf)
	if err != nil {
		return nil, err
	}

	request := &wrpc.PublishAndLogRequest{
		AnchorPsbt:        buf.Bytes(),
		VirtualPsbts:      make([][]byte, len(activeAssets)),
		PassiveAssetPsbts: make([][]byte, len(passiveAssets)),
		ChangeOutputIndex: commitResp.ChangeOutputIndex,
		LndLockedUtxos:    commitResp.LndLockedUtxos,
	}

	for idx := range activeAssets {
		request.VirtualPsbts[idx], err = tappsbt.Encode(
			activeAssets[idx],
		)
		if err != nil {
			return nil, err
		}
	}

	for idx := range passiveAssets {
		request.PassiveAssetPsbts[idx], err = tappsbt.Encode(
			passiveAssets[idx],
		)
		if err != nil {
			return nil, err
		}
	}

	resp, err := tapClient.PublishAndLogTransfer(ctx, request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func generateNonces(alicePubKey, bobPubKey *btcec.PublicKey) (aliceNonces, bobNonces *musig2.Nonces, err error) {
	aliceFundingNonceOpt := musig2.WithPublicKey(alicePubKey)
	aliceNonces, err = musig2.GenNonces(aliceFundingNonceOpt)
	if err != nil {
		return nil, nil, err
	}

	bobFundingNonceOpt := musig2.WithPublicKey(bobPubKey)
	bobNonces, err = musig2.GenNonces(bobFundingNonceOpt)
	if err != nil {
		return nil, nil, err
	}

	return aliceNonces, bobNonces, nil
}

func pubKeyBytes(k *btcec.PublicKey) []byte {
	return k.SerializeCompressed()
}

func updateWitness(a *asset.Asset, witness wire.TxWitness) {
	firstPrevWitness := &a.PrevWitnesses[0]
	if a.HasSplitCommitmentWitness() {
		rootAsset := firstPrevWitness.SplitCommitment.RootAsset
		firstPrevWitness = &rootAsset.PrevWitnesses[0]
	}
	firstPrevWitness.TxWitness = witness
}

func createMuSigLeaves(keys ...*musig2.AggregateKey) (asset.ScriptKey, []txscript.TapLeaf,
	*txscript.IndexedTapScriptTree, *txscript.ControlBlock, error) {

	leaves := make([]txscript.TapLeaf, len(keys))
	for i, key := range keys {
		muSigTapscript, err := txscript.NewScriptBuilder().
			AddData(schnorr.SerializePubKey(key.FinalKey)).
			AddOp(txscript.OP_CHECKSIG).
			Script()
		if err != nil {
			return asset.ScriptKey{}, nil, nil, nil, err
		}
		leaves[i] = txscript.TapLeaf{
			LeafVersion: txscript.BaseLeafVersion,
			Script:      muSigTapscript,
		}
	}

	tree := txscript.AssembleTaprootScriptTree(leaves...)
	internalKey := asset.NUMSPubKey
	controlBlock := &txscript.ControlBlock{
		LeafVersion: txscript.BaseLeafVersion,
		InternalKey: internalKey,
	}
	merkleRootHash := tree.RootNode.TapHash()

	tapKey := txscript.ComputeTaprootOutputKey(
		internalKey, merkleRootHash[:],
	)
	tapScriptKey := asset.ScriptKey{
		PubKey: tapKey,
		TweakedScriptKey: &asset.TweakedScriptKey{
			RawKey: keychain.KeyDescriptor{
				PubKey: internalKey,
			},
			Tweak: merkleRootHash[:],
		},
	}

	if tapKey.SerializeCompressed()[0] ==
		secp256k1.PubKeyFormatCompressedOdd {

		controlBlock.OutputKeyYIsOdd = true
	}

	return tapScriptKey, leaves, tree, controlBlock, nil
}

func AssertNonInteractiveRecvComplete(ctx context.Context,
	receiver taprpc.TaprootAssetsClient, totalInboundTransfers int) {

	err := wait.NoError(func() error {
		resp, err := receiver.AddrReceives(
			ctx, &taprpc.AddrReceivesRequest{},
		)
		if err != nil {
			return err
		}
		statusCompleted := taprpc.AddrEventStatus_ADDR_EVENT_STATUS_COMPLETED
		for _, event := range resp.Events {
			if event.Status != statusCompleted {
				return fmt.Errorf("got status %v, wanted %v",
					event.Status, statusCompleted)
			}

			if !event.HasProof {
				return fmt.Errorf("wanted proof, but was false")
			}
		}

		return nil
	}, time.Second*30/2)
	if err != nil {
		fmt.Println(err)
	}
}

func deserializeVPacket(packetBytes []byte) *tappsbt.VPacket {
	p, err := tappsbt.NewFromRawBytes(bytes.NewReader(packetBytes), false)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return p
}

func tapCreatePartialSig(ctx context.Context,
	params *chaincfg.Params, vPkt *tappsbt.VPacket,
	leafToSign txscript.TapLeaf, localKey keychain.KeyDescriptor,
	localNonces *musig2.Nonces, otherKey *btcec.PublicKey,
	otherNonces [musig2.PubNonceSize]byte) ([]byte, []byte, error) {

	conn1, clearUp1, err := getLightningClient()
	if err != nil {
		return nil, nil, err
	}
	defer clearUp1()

	conn2, clearUp2, err := getSignerClient()
	if err != nil {
		return nil, nil, err
	}
	defer clearUp2()

	conn3, clearUp3, err := getWalletKitClient()
	if err != nil {
		return nil, nil, err
	}
	defer clearUp3()

	client := LndClient{
		LN:        conn1,
		Signer:    conn2,
		WalletKit: conn3,
	}

	sessID := tapMuSig2Session(
		ctx, client, localKey, otherKey.SerializeCompressed(), *localNonces,
		[][]byte{otherNonces[:]},
	)

	partialSigner := &muSig2PartialSigner{
		sessID:     sessID,
		lnd:        client,
		leafToSign: leafToSign,
	}

	vIn := vPkt.Inputs[0]
	derivation, trDerivation := tappsbt.Bip32DerivationFromKeyDesc(
		keychain.KeyDescriptor{
			PubKey: localKey.PubKey,
		}, params.HDCoinType,
	)
	vIn.Bip32Derivation = []*psbt.Bip32Derivation{derivation}
	vIn.TaprootBip32Derivation = []*psbt.TaprootBip32Derivation{
		trDerivation,
	}

	err = tapsend.SignVirtualTransaction(
		vPkt, partialSigner, partialSigner,
	)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	isSplit, err := vPkt.HasSplitCommitment()
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	newAsset := vPkt.Outputs[0].Asset
	if isSplit {
		splitOut, err := vPkt.SplitRootOutput()
		if err != nil {
			fmt.Println(err)
			return nil, nil, err
		}

		newAsset = splitOut.Asset
	}

	partialSig := newAsset.PrevWitnesses[0].TxWitness[0][32:]

	return partialSig, sessID, nil
}

func tapMuSig2Session(ctx context.Context, lnd LndClient,
	localKey keychain.KeyDescriptor, otherKey []byte,
	localNonces musig2.Nonces, otherNonces [][]byte) []byte {

	version := signrpc.MuSig2Version_MUSIG2_VERSION_V100RC2
	sess, err := lnd.Signer.MuSig2CreateSession(
		ctx, &signrpc.MuSig2SessionRequest{
			KeyLoc: &signrpc.KeyLocator{
				KeyFamily: int32(localKey.Family),
				KeyIndex:  int32(localKey.Index),
			},
			AllSignerPubkeys: [][]byte{
				localKey.PubKey.SerializeCompressed(),
				otherKey,
			},
			OtherSignerPublicNonces: otherNonces,
			TaprootTweak: &signrpc.TaprootTweakDesc{
				KeySpendOnly: true,
			},
			Version:                version,
			PregeneratedLocalNonce: localNonces.SecNonce[:],
		},
	)
	if err != nil {
		return nil
	}

	return sess.SessionId
}

type muSig2PartialSigner struct {
	sessID     []byte
	lnd        LndClient
	leafToSign txscript.TapLeaf
}

func (m *muSig2PartialSigner) ValidateWitnesses(*asset.Asset,
	[]*commitment.SplitAsset, commitment.InputSet) error {

	return nil
}

func (m *muSig2PartialSigner) SignVirtualTx(_ *lndclient.SignDescriptor,
	tx *wire.MsgTx, prevOut *wire.TxOut) (*schnorr.Signature, error) {
	ctxb := context.Background()
	ctxt, cancel := context.WithTimeout(ctxb, time.Second*30)
	defer cancel()

	prevOutputFetcher := txscript.NewCannedPrevOutputFetcher(
		prevOut.PkScript, prevOut.Value,
	)
	sighashes := txscript.NewTxSigHashes(tx, prevOutputFetcher)

	sigHash, err := txscript.CalcTapscriptSignaturehash(
		sighashes, txscript.SigHashDefault, tx, 0, prevOutputFetcher,
		m.leafToSign,
	)
	if err != nil {
		return nil, err
	}

	sign, err := m.lnd.Signer.MuSig2Sign(
		ctxt, &signrpc.MuSig2SignRequest{
			SessionId:     m.sessID,
			MessageDigest: sigHash,
			Cleanup:       false,
		},
	)
	if err != nil {
		return nil, err
	}

	var sig [schnorr.SignatureSize]byte
	copy(sig[32:], sign.LocalPartialSignature)

	return schnorr.ParseSignature(sig[:])
}

func (m *muSig2PartialSigner) Execute(*asset.Asset, []*commitment.SplitAsset,
	commitment.InputSet) error {
	return nil
}

func combineSigs(ctx context.Context, sessID,
	otherPartialSig []byte, leafToSign txscript.TapLeaf,
	tree *txscript.IndexedTapScriptTree,
	controlBlock *txscript.ControlBlock) (wire.TxWitness, error) {
	conn, clearUp, err := getSignerClient()
	if err != nil {
		return nil, err
	}
	defer clearUp()

	resp, err := conn.MuSig2CombineSig(
		ctx, &signrpc.MuSig2CombineSigRequest{
			SessionId:              sessID,
			OtherPartialSignatures: [][]byte{otherPartialSig},
		},
	)
	if err != nil {
		return nil, err
	}

	for _, leaf := range tree.LeafMerkleProofs {
		if leaf.TapHash() == leafToSign.TapHash() {
			controlBlock.InclusionProof = leaf.InclusionProof
		}
	}

	controlBlockBytes, err := controlBlock.ToBytes()
	if err != nil {
		return nil, err
	}

	commitmentWitness := make(wire.TxWitness, 3)
	commitmentWitness[0] = resp.FinalSignature
	commitmentWitness[1] = leafToSign.Script
	commitmentWitness[2] = controlBlockBytes

	return commitmentWitness, nil
}

func CommitVirtualPsbts(ctx context.Context, packet *psbt.Packet,
	activePackets []*tappsbt.VPacket, passivePackets []*tappsbt.VPacket,
	changeOutputIndex int32) (*psbt.Packet, []*tappsbt.VPacket,
	[]*tappsbt.VPacket, *wrpc.CommitVirtualPsbtsResponse) {

	tapClient, clearUp, err := GetAssetWalletClient()
	if err != nil {
		return nil, nil, nil, nil
	}
	defer clearUp()

	var feeRateSatPerKVByte chainfee.SatPerKVByte = 2000

	var buf bytes.Buffer
	err = packet.Serialize(&buf)
	if err != nil {
		return nil, nil, nil, nil
	}

	request := &wrpc.CommitVirtualPsbtsRequest{
		AnchorPsbt: buf.Bytes(),
		Fees: &wrpc.CommitVirtualPsbtsRequest_SatPerVbyte{
			SatPerVbyte: uint64(feeRateSatPerKVByte / 1000),
		},
	}

	type existingIndex = wrpc.CommitVirtualPsbtsRequest_ExistingOutputIndex
	if changeOutputIndex < 0 {
		request.AnchorChangeOutput = &wrpc.CommitVirtualPsbtsRequest_Add{
			Add: true,
		}
	} else {
		request.AnchorChangeOutput = &existingIndex{
			ExistingOutputIndex: changeOutputIndex,
		}
	}

	request.VirtualPsbts = make([][]byte, len(activePackets))
	for idx := range activePackets {
		request.VirtualPsbts[idx], err = tappsbt.Encode(
			activePackets[idx],
		)
		if err != nil {
			return nil, nil, nil, nil
		}
	}
	request.PassiveAssetPsbts = make([][]byte, len(passivePackets))
	for idx := range passivePackets {
		request.PassiveAssetPsbts[idx], err = tappsbt.Encode(
			passivePackets[idx],
		)
		if err != nil {
			return nil, nil, nil, nil
		}
	}

	commitResponse, err := tapClient.CommitVirtualPsbts(ctx, request)
	if err != nil {
		return nil, nil, nil, nil
	}

	fundedPacket, err := psbt.NewFromRawBytes(
		bytes.NewReader(commitResponse.AnchorPsbt), false,
	)
	if err != nil {
		return nil, nil, nil, nil
	}

	activePackets = make(
		[]*tappsbt.VPacket, len(commitResponse.VirtualPsbts),
	)
	for idx := range commitResponse.VirtualPsbts {
		activePackets[idx], err = tappsbt.Decode(
			commitResponse.VirtualPsbts[idx],
		)
		if err != nil {
			return nil, nil, nil, nil
		}
	}

	passivePackets = make(
		[]*tappsbt.VPacket, len(commitResponse.PassiveAssetPsbts),
	)
	for idx := range commitResponse.PassiveAssetPsbts {
		passivePackets[idx], err = tappsbt.Decode(
			commitResponse.PassiveAssetPsbts[idx],
		)
		if err != nil {
			return nil, nil, nil, nil
		}
	}

	return fundedPacket, activePackets, passivePackets, commitResponse
}

func partialSignWithKey(ctx context.Context, params *chaincfg.Params,
	pkt *psbt.Packet, inputIndex uint32, key keychain.KeyDescriptor,
	controlBlockBytes []byte, tapLeaf txscript.TapLeaf) []byte {
	conn, clearUp, err := getWalletKitClient()
	if err != nil {
		return nil
	}
	defer clearUp()

	leafToSign := []*psbt.TaprootTapLeafScript{{
		ControlBlock: controlBlockBytes,
		Script:       tapLeaf.Script,
		LeafVersion:  tapLeaf.LeafVersion,
	}}

	signInput := &pkt.Inputs[inputIndex]
	derivation, trDerivation := tappsbt.Bip32DerivationFromKeyDesc(
		key, params.HDCoinType,
	)
	trDerivation.LeafHashes = [][]byte{fn.ByteSlice(tapLeaf.TapHash())}
	signInput.Bip32Derivation = []*psbt.Bip32Derivation{derivation}
	signInput.TaprootBip32Derivation = []*psbt.TaprootBip32Derivation{
		trDerivation,
	}
	signInput.TaprootLeafScript = leafToSign
	signInput.SighashType = txscript.SigHashDefault

	var buf bytes.Buffer
	err = pkt.Serialize(&buf)
	if err != nil {
		return nil
	}

	resp, err := conn.SignPsbt(
		ctx, &walletrpc.SignPsbtRequest{
			FundedPsbt: buf.Bytes(),
		},
	)
	if err != nil {
		return nil
	}

	result, err := psbt.NewFromRawBytes(
		bytes.NewReader(resp.SignedPsbt), false,
	)
	if err != nil {
		return nil
	}

	return result.Inputs[inputIndex].TaprootScriptSpendSig[0].Signature
}

func FinalizePacket(ctx context.Context,
	pkt *psbt.Packet) (*psbt.Packet, error) {

	conn, clearUp, err := getWalletKitClient()
	if err != nil {
		return nil, err
	}
	defer clearUp()
	var buf bytes.Buffer
	err = pkt.Serialize(&buf)
	if err != nil {
		return nil, err
	}

	finalizeResp, err := conn.FinalizePsbt(ctx, &walletrpc.FinalizePsbtRequest{
		FundedPsbt: buf.Bytes(),
	})
	if err != nil {
		return nil, err
	}

	signedPacket, err := psbt.NewFromRawBytes(
		bytes.NewReader(finalizeResp.SignedPsbt), false,
	)
	if err != nil {
		return nil, err
	}

	return signedPacket, nil
}

func DeriveKeys(ctx context.Context) (*wrpc.NextScriptKeyResponse,
	*wrpc.NextInternalKeyResponse, error) {

	conn, clearUp, err := GetAssetWalletClient()
	if err != nil {
		return nil, nil, err
	}
	defer clearUp()
	scriptKeyDesc, err := conn.NextScriptKey(
		ctx, &wrpc.NextScriptKeyRequest{
			KeyFamily: uint32(asset.TaprootAssetsKeyFamily),
		},
	)
	if err != nil {
		return nil, nil, err
	}

	internalKeyDesc, err := conn.NextInternalKey(
		ctx, &wrpc.NextInternalKeyRequest{
			KeyFamily: uint32(asset.TaprootAssetsKeyFamily),
		},
	)
	if err != nil {
		return nil, nil, err
	}

	return scriptKeyDesc, internalKeyDesc, nil
}

func EncodeVPackets(vPackets []*tappsbt.VPacket) ([]byte, error) {
	var buf bytes.Buffer

	for _, vPkt := range vPackets {
		encoded, err := tappsbt.Encode(vPkt)
		if err != nil {
			return nil, err
		}

		length := uint32(len(encoded))
		err = binary.Write(&buf, binary.LittleEndian, length)
		if err != nil {
			return nil, err
		}

		_, err = buf.Write(encoded)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func DecodeVPackets(encoded []byte) ([]*tappsbt.VPacket, error) {
	var vPackets []*tappsbt.VPacket
	buf := bytes.NewReader(encoded)

	for {
		var length uint32
		err := binary.Read(buf, binary.LittleEndian, &length)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		data := make([]byte, length)
		_, err = buf.Read(data)
		if err != nil {
			return nil, err
		}

		vPkt, err := tappsbt.Decode(data)
		if err != nil {
			return nil, err
		}

		vPackets = append(vPackets, vPkt)
	}

	return vPackets, nil
}

func DecodeBase64ToPSBT(encodedPsbt string) (*psbt.Packet, error) {
	psbtBytes, err := base64.StdEncoding.DecodeString(encodedPsbt)
	if err != nil {
		return nil, err
	}

	psbtPacket, err := psbt.NewFromRawBytes(bytes.NewReader(psbtBytes), false)
	if err != nil {
		return nil, err
	}

	return psbtPacket, nil
}

func AssertAddrEventCustomTimeout(ctx context.Context, client taprpc.TaprootAssetsClient, addr *taprpc.Addr,
	numEvents int, expectedStatus taprpc.AddrEventStatus, timeout time.Duration) error {

	err := wait.NoError(func() error {
		resp, err := client.AddrReceives(
			ctx, &taprpc.AddrReceivesRequest{
				FilterAddr: addr.Encoded,
			},
		)
		if err != nil {
			return err
		}

		if len(resp.Events) != numEvents {
			return fmt.Errorf("got %d events, wanted %d",
				len(resp.Events), numEvents)
		}

		if resp.Events[0].Status != expectedStatus {
			return fmt.Errorf("got status %v, wanted %v",
				resp.Events[0].Status, expectedStatus)
		}

		return nil
	}, timeout)
	if err != nil {
		return err
	}
	return nil
}
