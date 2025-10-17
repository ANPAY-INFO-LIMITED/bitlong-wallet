package exchange

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"

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
	wrpc "github.com/lightninglabs/taproot-assets/taprpc/assetwalletrpc"
	"github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	unirpc "github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	"github.com/lightninglabs/taproot-assets/tapsend"
	"github.com/lightninglabs/taproot-assets/universe"
	"github.com/lightningnetwork/lnd/keychain"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/walletrpc"
	"github.com/lightningnetwork/lnd/lntest"
	"github.com/lightningnetwork/lnd/lntest/wait"
	"github.com/wallet/service/apiConnect"
	"github.com/wallet/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func DeriveKeys() (asset.ScriptKey, keychain.KeyDescriptor, error) {
	const defaultTimeout = 100 * time.Second
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Printf("did not connect: %v\n", err)
	}
	defer clearUp()
	AssetWalletClient := assetwalletrpc.NewAssetWalletClient(conn)
	ctx := context.Background()
	ctxt, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	scriptKeyDesc, err := AssetWalletClient.NextScriptKey(
		ctxt, &wrpc.NextScriptKeyRequest{
			KeyFamily: uint32(asset.TaprootAssetsKeyFamily),
		},
	)
	if err != nil {
		return asset.ScriptKey{}, keychain.KeyDescriptor{}, err
	}
	scriptKey, err := rpcutils.UnmarshalScriptKey(scriptKeyDesc.ScriptKey)
	if err != nil {
		return asset.ScriptKey{}, keychain.KeyDescriptor{}, err
	}

	internalKeyDesc, err := AssetWalletClient.NextInternalKey(
		ctxt, &wrpc.NextInternalKeyRequest{
			KeyFamily: uint32(asset.TaprootAssetsKeyFamily),
		},
	)
	if err != nil {
		return asset.ScriptKey{}, keychain.KeyDescriptor{}, err
	}
	internalKeyLnd, err := rpcutils.UnmarshalKeyDescriptor(internalKeyDesc.InternalKey)
	if err != nil {
		return asset.ScriptKey{}, keychain.KeyDescriptor{}, err
	}

	return *scriptKey, internalKeyLnd, nil
}

func CreateVirtualPSBT(assetID asset.ID, numUnits uint64, scriptKey asset.ScriptKey, anchorInternalKey keychain.KeyDescriptor, chainParams *address.ChainParams) (*tappsbt.VPacket, error) {
	vPkt := tappsbt.ForInteractiveSend(
		assetID, numUnits-20000, scriptKey, 0, 0, 0,
		anchorInternalKey, asset.V0, chainParams,
	)
	vPkt.Outputs[0].Type = tappsbt.TypeSplitRoot
	scriptKey1, anchorInternalKey1, err := DeriveKeys()
	if err != nil {
		return nil, err
	}
	tappsbt.AddOutput(
		vPkt, 20000, scriptKey1, 1,
		anchorInternalKey1, asset.V0,
	)

	return vPkt, nil
}

func maybeFundPacket(vPkg *tappsbt.VPacket) (*assetwalletrpc.FundVirtualPsbtResponse, error) {
	const (
		maxRetries     = 5
		retryDelay     = 10 * time.Second
		defaultTimeout = 100 * time.Second
	)

	var buf bytes.Buffer
	if err := vPkg.Serialize(&buf); err != nil {
		return nil, fmt.Errorf("failed to serialize PSBT: %w", err)
	}

	for attempt := 0; attempt < maxRetries; attempt++ {
		if err := checkUTXOStatus(); err != nil {
			fmt.Printf("UTXO check failed (attempt %d): %v. Retrying in %v...\n",
				attempt+1, err, retryDelay)
			time.Sleep(retryDelay)
			continue
		}

		conn, clearUp, err := apiConnect.GetConnection("tapd", false)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to tapd: %w", err)
		}
		defer clearUp()

		assetWalletClient := assetwalletrpc.NewAssetWalletClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()

		req := &assetwalletrpc.FundVirtualPsbtRequest{
			Template: &assetwalletrpc.FundVirtualPsbtRequest_Psbt{
				Psbt: buf.Bytes(),
			},
		}

		resp, err := assetWalletClient.FundVirtualPsbt(ctx, req)
		if err == nil {
			return resp, nil
		}

		st, ok := status.FromError(err)
		if !ok {
			return nil, fmt.Errorf("non-gRPC error in FundVirtualPsbt: %w", err)
		}

		switch st.Code() {
		case codes.Unavailable, codes.DeadlineExceeded, codes.ResourceExhausted:
			fmt.Printf("Transient error in FundVirtualPsbt (attempt %d): %v. Retrying in %v...\n",
				attempt+1, err, retryDelay)
			time.Sleep(retryDelay)
		case codes.FailedPrecondition:
			fmt.Printf("Failed precondition in FundVirtualPsbt (attempt %d): %v. This might indicate insufficient funds or UTXO issues. Retrying in %v...\n",
				attempt+1, err, retryDelay)
			time.Sleep(retryDelay)
		default:
			return nil, fmt.Errorf("error in FundVirtualPsbt: %w", err)
		}
	}

	return nil, fmt.Errorf("failed to fund virtual PSBT after %d attempts", maxRetries)
}

func checkUTXOStatus() error {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return fmt.Errorf("failed to connect to LND: %v", err)
	}
	defer clearUp()

	client := lnrpc.NewLightningClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.ListUnspent(ctx, &lnrpc.ListUnspentRequest{
		MinConfs: 1, // 只查看已确认的 UTXO
		MaxConfs: 9999999,
	})
	if err != nil {
		return fmt.Errorf("failed to list unspent: %v", err)
	}

	if len(resp.Utxos) == 0 {
		return fmt.Errorf("no confirmed UTXOs available")
	}

	var totalBalance int64
	for _, utxo := range resp.Utxos {
		totalBalance += utxo.AmountSat
	}

	fmt.Printf("Found %d confirmed UTXOs with total balance: %d satoshis\n", len(resp.Utxos), totalBalance)
	return nil
}

type AliceData struct {
	BtcPsbt          []byte
	B                []byte
	Resp             *assetwalletrpc.CommitVirtualPsbtsResponse
	SignedVpsbtBytes []byte
	AssetIDStr       string
	NumUnits         uint64
}

func FundVirtualPSBT(vPkt *tappsbt.VPacket) (*tappsbt.VPacket, error) {
	fundResp, err := maybeFundPacket(vPkt)
	if err != nil {
		return nil, err
	}
	vPkt, err = tappsbt.Decode(fundResp.FundedPsbt)
	if err != nil {
		return nil, err
	}
	return vPkt, nil
}

func SignVirtualPSBT(vPkt *tappsbt.VPacket) (*assetwalletrpc.SignVirtualPsbtResponse, error) {
	vPktBytes, err := tappsbt.Encode(vPkt)
	if err != nil {
		return nil, err
	}

	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Printf("did not connect: %v\n", err)
	}
	defer clearUp()
	AssetWalletClient := assetwalletrpc.NewAssetWalletClient(conn)

	signedResp, err := AssetWalletClient.SignVirtualPsbt(context.Background(), &assetwalletrpc.SignVirtualPsbtRequest{
		FundedPsbt: vPktBytes,
	})
	if err != nil {
		return nil, err
	}
	return signedResp, nil
}

func PrepareBitcoinPSBT(vPkt *tappsbt.VPacket) (*psbt.Packet, error) {
	btcpsbt, err := tapsend.PrepareAnchoringTemplate([]*tappsbt.VPacket{vPkt})
	if err != nil {
		return nil, err
	}

	return btcpsbt, nil
}

func GeneratePaymentScript() ([]byte, string, error) {
	const defaultTimeout = 100 * time.Second
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		fmt.Printf("did not connect: %v\n", err)
		return nil, "", err
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	ctxb := context.Background()
	ctxt, cancel := context.WithTimeout(ctxb, defaultTimeout)
	defer cancel()
	req := &lnrpc.NewAddressRequest{
		Type: lnrpc.AddressType_TAPROOT_PUBKEY,
	}
	addrResp, _ := client.NewAddress(ctxt, req)
	decodeAddress, err := btcutil.DecodeAddress(addrResp.Address, &chaincfg.RegressionNetParams)
	if err != nil {
		return nil, "", err
	}
	pkScript, err := txscript.PayToAddrScript(decodeAddress)
	if err != nil {
		return nil, "", err
	}
	return pkScript, addrResp.Address, nil
}

type MacaroonCredential struct {
	macaroon string
}

func NewMacaroonCredential(macaroon string) *MacaroonCredential {
	return &MacaroonCredential{macaroon: macaroon}
}
func (c *MacaroonCredential) RequireTransportSecurity() bool {
	return true
}

func (c *MacaroonCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"macaroon": c.macaroon}, nil
}

func GetMacaroon(macaroonPath string) string {
	macaroonBytes, err := os.ReadFile(macaroonPath)
	if err != nil {
		panic(err)
	}
	macaroon := hex.EncodeToString(macaroonBytes)
	return macaroon
}

func serverDialOpts(macaroon string) ([]grpc.DialOption, error) {
	var opts []grpc.DialOption

	tlsConfig := tls.Config{InsecureSkipVerify: true}
	transportCredentials := credentials.NewTLS(&tlsConfig)
	opts = append(opts, grpc.WithTransportCredentials(transportCredentials), grpc.WithPerRPCCredentials(NewMacaroonCredential(macaroon)))

	return opts, nil
}

func CommitVirtualPsbts(vPktBytes, btcPsbtBytes []byte) (*assetwalletrpc.CommitVirtualPsbtsResponse, error) {
	if len(vPktBytes) == 0 {
		return nil, fmt.Errorf("vPktBytes is nil or empty")
	}
	if len(btcPsbtBytes) == 0 {
		return nil, fmt.Errorf("btcPsbtBytes is nil or empty")
	}

	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Printf("did not connect: %v\n", err)
		return nil, err
	}
	defer clearUp()

	AssetWalletClient := assetwalletrpc.NewAssetWalletClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	resp, err := AssetWalletClient.CommitVirtualPsbts(ctx, &assetwalletrpc.CommitVirtualPsbtsRequest{
		VirtualPsbts: [][]byte{vPktBytes},
		AnchorPsbt:   btcPsbtBytes,
		AnchorChangeOutput: &assetwalletrpc.CommitVirtualPsbtsRequest_Add{
			Add: true,
		},
		Fees: &assetwalletrpc.CommitVirtualPsbtsRequest_TargetConf{
			TargetConf: 12,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("CommitVirtualPsbts RPC call failed: %v", err)
	}
	return resp, nil
}

func SignBitcoinPSBT(btcpsbt []byte) (*psbt.Packet, error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		fmt.Printf("did not connect: %v\n", err)
	}
	defer clearUp()
	client := walletrpc.NewWalletKitClient(conn)
	req := &walletrpc.SignPsbtRequest{
		FundedPsbt: btcpsbt,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	signResp, err := client.SignPsbt(ctx, req)
	if err != nil {
		return nil, err
	}
	return psbt.NewFromRawBytes(bytes.NewReader(signResp.SignedPsbt), false)
}

func VerifyAssetTransfer(expectedAmount uint64) (bool, error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Printf("did not connect: %v\n", err)
	}
	defer clearUp()
	assetsClient := taprpc.NewTaprootAssetsClient(conn)

	assets, err := assetsClient.ListAssets(context.Background(), &taprpc.ListAssetRequest{})
	if err != nil {
		return false, err
	}
	if len(assets.Assets) > 0 && assets.Assets[0].Amount == expectedAmount {
		return true, nil
	}
	return false, nil
}

func getTimeNow() string {
	return time.Now().Format("2006/01/02 15:04:05")
}
func assetLeaves(isGroup bool, id string, proofType universerpc.ProofType) (*universerpc.AssetLeafResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", getTimeNow(), err)
	}
	defer clearUp()
	request := &universerpc.ID{
		ProofType: proofType,
	}

	if isGroup {
		groupKey := &universerpc.ID_GroupKeyStr{
			GroupKeyStr: id,
		}
		request.Id = groupKey
	} else {
		AssetId := &universerpc.ID_AssetIdStr{
			AssetIdStr: id,
		}
		request.Id = AssetId
	}

	client := universerpc.NewUniverseClient(conn)
	response, err := client.AssetLeaves(context.Background(), request)
	return response, err
}

func GetAssetInfo(id string) *taprpc.Asset {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	client := universerpc.NewUniverseClient(conn)

	defer clearUp()

	in := &universerpc.AssetRootQuery{}
	in.Id = &universerpc.ID{
		Id: &universerpc.ID_AssetIdStr{
			AssetIdStr: id,
		},
	}
	roots, err := client.QueryAssetRoots(context.Background(), in)
	if err != nil {
		return nil
	}

	if roots == nil || roots.IssuanceRoot.Id == nil {
		return nil
	}
	queryId := id
	isGroup := false
	if groupKey, ok := roots.IssuanceRoot.Id.Id.(*universerpc.ID_GroupKey); ok {
		isGroup = true
		queryId = hex.EncodeToString(groupKey.GroupKey)
	}
	response, err := assetLeaves(isGroup, queryId, universerpc.ProofType_PROOF_TYPE_ISSUANCE)
	if err != nil {
		return nil
	}
	if response.Leaves == nil {
		return nil
	}
	for _, leaf := range response.Leaves {
		if hex.EncodeToString(leaf.Asset.AssetGenesis.GetAssetId()) == id {
			return leaf.Asset
		}
	}
	return nil
}

func ListAssets(id string) *taprpc.Asset {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	client := universerpc.NewUniverseClient(conn)

	defer clearUp()

	in := &universerpc.AssetRootQuery{}
	in.Id = &universerpc.ID{
		Id: &universerpc.ID_AssetIdStr{
			AssetIdStr: id,
		},
	}
	roots, err := client.QueryAssetRoots(context.Background(), in)
	if err != nil {
		return nil
	}

	if roots == nil || roots.IssuanceRoot.Id == nil {
		return nil
	}
	queryId := id
	isGroup := false
	if groupKey, ok := roots.IssuanceRoot.Id.Id.(*universerpc.ID_GroupKey); ok {
		isGroup = true
		queryId = hex.EncodeToString(groupKey.GroupKey)
	}
	response, err := assetLeaves(isGroup, queryId, universerpc.ProofType_PROOF_TYPE_ISSUANCE)
	if err != nil {
		return nil
	}
	if response.Leaves == nil {
		return nil
	}
	for _, leaf := range response.Leaves {
		if hex.EncodeToString(leaf.Asset.AssetGenesis.GetAssetId()) == id {
			return leaf.Asset
		}
	}
	return nil
}

func AliceBiz(assetID asset.ID, numUnits uint64, ctxb context.Context) (*psbt.Packet, bytes.Buffer, *wrpc.CommitVirtualPsbtsResponse, []byte) {
	fmt.Printf("Starting AliceBiz function with assetID: %v, numUnits: %d\n", assetID, numUnits)

	scriptKey, anchorInternalKey, err := DeriveKeys()
	if err != nil {
		fmt.Printf("Key derivation failed: %v\n", err)
		return nil, bytes.Buffer{}, nil, nil
	}
	fmt.Printf("Keys derived successfully for Alice\n")

	vPkt, err := CreateVirtualPSBT(assetID, numUnits, scriptKey, anchorInternalKey, &address.RegressionNetTap)
	if err != nil {
		fmt.Printf("Creating virtual PSBT failed: %v\n", err)
		return nil, bytes.Buffer{}, nil, nil
	}
	fmt.Printf("Virtual PSBT created successfully\n")

	vPkt, err = FundVirtualPSBT(vPkt)
	if err != nil {
		fmt.Printf("Funding virtual PSBT failed: %v\n", err)
		return nil, bytes.Buffer{}, nil, nil
	}
	fmt.Printf("Virtual PSBT funded successfully\n")

	vPkt.Inputs[0].SighashType = txscript.SigHashSingle
	fmt.Printf("Set SighashType to SigHashNone and Output Type to TypeSimple\n")
	err = tapsend.PrepareOutputAssets(ctxb, vPkt)
	if err != nil {
		fmt.Printf("Preparing output assets failed: %v\n", err)
		return nil, bytes.Buffer{}, nil, nil
	}
	if vPkt.Outputs[1].Asset.SplitCommitmentRoot == nil {
		fmt.Printf("SplitCommitmentRoot nil/n")
	}
	if vPkt.Outputs[1].Asset.PrevWitnesses == nil {
		fmt.Printf("PrevWitnesses nil/n")
	}
	if vPkt.Outputs[1].Asset.PrevWitnesses[0].SplitCommitment == nil {
		fmt.Printf("SplitCommitment nil/n")
	}
	signedResp, err := SignVirtualPSBT(vPkt)
	if err != nil {
		fmt.Printf("Signing virtual PSBT failed: %v\n", err)
		return nil, bytes.Buffer{}, nil, nil
	}
	fmt.Printf("Virtual PSBT signed successfully\n")

	for _, input := range signedResp.SignedInputs {
		if input != uint32(0) {
			fmt.Printf("Unexpected SignedInputs: expected 0, got %d\n", input)
			return nil, bytes.Buffer{}, nil, nil
		}
	}

	vPkt, err = tappsbt.Decode(signedResp.SignedPsbt)
	if err != nil {
		fmt.Printf("Error decoding signed PSBT: %v\n", err)
		return nil, bytes.Buffer{}, nil, nil
	}
	fmt.Printf("Signed PSBT decoded successfully\n")

	btcPsbt, err := PrepareBitcoinPSBT(vPkt)
	if err != nil {
		fmt.Printf("Preparing Bitcoin PSBT failed: %v\n", err)
		return nil, bytes.Buffer{}, nil, nil
	}
	fmt.Printf("Bitcoin PSBT prepared successfully\n")

	pkScript, address1, err := GeneratePaymentScript()
	if err != nil {
		fmt.Printf("Generating payment script failed: %v\n", err)
		return nil, bytes.Buffer{}, nil, nil
	}
	fmt.Printf("Payment script generated successfully for address: %s\n", address1)

	btcPsbt.UnsignedTx.TxOut[0].PkScript = pkScript
	btcPsbt.UnsignedTx.TxOut[0].Value = int64(numUnits)
	fmt.Printf("Set PkScript and Value for first output in Bitcoin PSBT\n")

	derivation, trDerivation, err := getAddressBip32Derivation(address1)
	if err != nil {
		fmt.Printf("Error getting address BIP32 derivation: %v\n", err)
		return nil, bytes.Buffer{}, nil, nil
	}
	fmt.Printf("Got BIP32 derivation for address\n")

	btcPsbt.Outputs[0].Bip32Derivation = []*psbt.Bip32Derivation{derivation}
	btcPsbt.Outputs[0].TaprootBip32Derivation = []*psbt.TaprootBip32Derivation{trDerivation}
	btcPsbt.Outputs[0].TaprootInternalKey = trDerivation.XOnlyPubKey
	btcPsbt.Outputs[1].Bip32Derivation = []*psbt.Bip32Derivation{derivation}
	btcPsbt.Outputs[1].TaprootBip32Derivation = []*psbt.TaprootBip32Derivation{trDerivation}
	btcPsbt.Outputs[1].TaprootInternalKey = trDerivation.XOnlyPubKey
	fmt.Printf("Set BIP32 derivation and Taproot internal key for first output in Bitcoin PSBT\n")

	var b bytes.Buffer
	err = btcPsbt.Serialize(&b)
	if err != nil {
		fmt.Printf("Error serializing Bitcoin PSBT: %v\n", err)
		return nil, bytes.Buffer{}, nil, nil
	}
	fmt.Printf("Bitcoin PSBT serialized successfully\n")

	resp, err := CommitVirtualPsbts(signedResp.SignedPsbt, b.Bytes())
	if err != nil {
		fmt.Printf("Committing PSBTs failed: %v\n", err)
		return nil, bytes.Buffer{}, nil, nil
	}
	fmt.Printf("PSBTs committed successfully\n")

	btcPsbt, err = psbt.NewFromRawBytes(bytes.NewReader(resp.AnchorPsbt), false)
	if err != nil {
		fmt.Printf("Error creating PSBT from raw bytes: %v\n", err)
		return nil, bytes.Buffer{}, nil, nil
	}
	fmt.Printf("Created new Bitcoin PSBT from raw bytes\n")

	btcPsbt.Inputs[0].SighashType = txscript.SigHashSingle | txscript.SigHashAnyOneCanPay
	btcPsbt.Inputs = append(btcPsbt.Inputs[:1], btcPsbt.Inputs[2:]...)
	btcPsbt.UnsignedTx.TxIn = append(btcPsbt.UnsignedTx.TxIn[:1], btcPsbt.UnsignedTx.TxIn[2:]...)
	fmt.Printf("Modified Bitcoin PSBT inputs\n")

	btcPsbt.Outputs = btcPsbt.Outputs[:2]
	btcPsbt.UnsignedTx.TxOut = btcPsbt.UnsignedTx.TxOut[:2]
	fmt.Printf("Removed change output from Bitcoin PSBT\n")

	b.Reset()
	err = btcPsbt.Serialize(&b)
	if err != nil {
		fmt.Printf("Error serializing modified Bitcoin PSBT: %v\n", err)
		return nil, bytes.Buffer{}, nil, nil
	}
	fmt.Printf("Modified Bitcoin PSBT serialized successfully\n")

	btcPsbt, err = SignBitcoinPSBT(b.Bytes())
	if err != nil {
		fmt.Printf("Signing Bitcoin PSBT failed: %v\n", err)
		return nil, bytes.Buffer{}, nil, nil
	}
	fmt.Printf("Bitcoin PSBT signed successfully\n")

	signedVpsbtBytes, err := tappsbt.Encode(vPkt)
	if err != nil {
		fmt.Printf("Error encoding signed virtual PSBT: %v\n", err)
		return nil, bytes.Buffer{}, nil, nil
	}
	fmt.Printf("Signed virtual PSBT encoded successfully\n")

	fmt.Printf("AliceBiz function completed successfully\n")
	return btcPsbt, b, resp, signedVpsbtBytes
}

func DeepCopyWitness(witnesses []asset.Witness) []asset.Witness {
	copied := make([]asset.Witness, len(witnesses))
	for i, w := range witnesses {
		copied[i] = asset.Witness{
			PrevID:          w.PrevID,
			TxWitness:       w.TxWitness,
			SplitCommitment: w.SplitCommitment,
		}
	}
	return copied
}

type BobData struct {
	Psbt              []byte
	Tappsbt           []byte
	BobScriptKeyBytes []byte
	AssetIDstr        string
	NumUnits          uint64
	Resp1             *wrpc.CommitVirtualPsbtsResponse
}

func BobBiz(ctx context.Context, req *types.BobBizReq) (resp *types.BobBizResp) {
	fmt.Printf("=== BobBiz Start ===\n")
	fmt.Printf("AssetID: %s, NumUnits: %d\n", req.AssetID.String(), req.NumUnits)

	bobVPsbt, err := tappsbt.Decode(req.SignedVpsbtBytes)
	if err != nil {
		fmt.Printf("Error decoding signed VPSBT: %v\n", err)
		return
	}
	fmt.Printf("1. Decoded VPSBT: Inputs=%d, Outputs=%d\n", len(bobVPsbt.Inputs), len(bobVPsbt.Outputs))

	bobScriptKey, bobAnchorInternalKey, err := DeriveKeys()
	if err != nil {
		fmt.Printf("Error deriving keys: %v\n", err)
		return
	}
	fmt.Printf("2. Derived Bob's keys\n")

	bobVOut := bobVPsbt.Outputs[0]
	bobVOut.ScriptKey = bobScriptKey
	bobVOut.AnchorOutputBip32Derivation = nil
	bobVOut.AnchorOutputTaprootBip32Derivation = nil
	bobVOut.SetAnchorInternalKey(bobAnchorInternalKey, 1)
	fmt.Printf("3. Updated VPSBT output\n")

	deliveryAddrStr := fmt.Sprintf("%s://%s", proof.UniverseRpcCourierType, "132.232.109.84:8443")
	deliveryAddr, err := url.Parse(deliveryAddrStr)
	if err != nil {
		fmt.Printf("Error parsing delivery address: %v\n", err)
		return
	}
	bobVPsbt.Outputs[0].ProofDeliveryAddress = deliveryAddr
	fmt.Printf("4. Set proof delivery address: %s\n", deliveryAddrStr)

	req.BtcPsbt.Outputs[1].TaprootInternalKey = schnorr.SerializePubKey(bobAnchorInternalKey.PubKey)
	req.BtcPsbt.Outputs[1].Bip32Derivation = bobVOut.AnchorOutputBip32Derivation
	req.BtcPsbt.Outputs[1].TaprootBip32Derivation = bobVOut.AnchorOutputTaprootBip32Derivation
	fmt.Printf("5. Updated BTC PSBT outputs\n")

	witnessBackup := bobVPsbt.Outputs[0].Asset.PrevWitnesses
	fmt.Printf("6. Backed up %d witnesses\n", len(witnessBackup))

	err = tapsend.PrepareOutputAssets(ctx, bobVPsbt)
	if err != nil {
		fmt.Printf("Error preparing output assets: %v\n", err)
		return
	}
	fmt.Printf("7. Prepared output assets\n")

	if bobVPsbt.Outputs[0].ScriptKey != bobVPsbt.Outputs[0].Asset.ScriptKey {
		fmt.Printf("ScriptKey mismatch: expected %v, got %v\n",
			bobVPsbt.Outputs[0].ScriptKey, bobVPsbt.Outputs[0].Asset.ScriptKey)
		return
	}
	fmt.Printf("8. ScriptKey consistency verified\n")

	bobVPsbt.Outputs[0].Asset.PrevWitnesses = witnessBackup
	fmt.Printf("9. Restored %d witnesses\n", len(bobVPsbt.Outputs[0].Asset.PrevWitnesses))

	bobVPsbtBytes, err := tappsbt.Encode(bobVPsbt)
	if err != nil {
		fmt.Printf("Error encoding bobVPsbt: %v\n", err)
		return
	}
	fmt.Printf("10. Encoded VPSBT, size: %d bytes\n", len(bobVPsbtBytes))

	b := req.B
	b.Reset()
	err = req.BtcPsbt.Serialize(&b)
	if err != nil {
		fmt.Printf("Error serializing btcPsbt: %v\n", err)
		return
	}
	fmt.Printf("11. Serialized BTC PSBT, size: %d bytes\n", b.Len())

	resp1, err := CommitVirtualPsbts(bobVPsbtBytes, b.Bytes())
	if err != nil {
		fmt.Printf("Committing PSBTs failed: %v\n", err)
		return
	}
	fmt.Printf("12. Committed virtual PSBTs\n")

	bobVPsbt, err = tappsbt.Decode(resp1.VirtualPsbts[0])
	if err != nil {
		fmt.Printf("Error decoding returned bobVPsbt: %v\n", err)
		return
	}
	fmt.Printf("13. Decoded returned VPSBT\n")

	signResp1, err := signPsbt(resp1.AnchorPsbt)
	if err != nil {
		fmt.Printf("Error signing returned psbt: %v\n", err)
		return
	}
	fmt.Printf("14. Signed anchor PSBT\n")

	finalPsbt, err := psbt.NewFromRawBytes(bytes.NewReader(signResp1.SignedPsbt), false)
	if err != nil {
		fmt.Printf("Error creating PSBT from raw bytes: %v\n", err)
		return
	}
	fmt.Printf("15. Created final PSBT: Inputs=%d, Outputs=%d\n", len(finalPsbt.Inputs), len(finalPsbt.Outputs))

	psbtBase64, err := encodePSBTtoBase64(finalPsbt)
	if err != nil {
		fmt.Printf("Error encoding PSBT to Base64: %v\n", err)
		return
	}
	fmt.Println("Encoded PSBT in Base64:", psbtBase64)

	signedPkt1, err := finalizePacket(finalPsbt)
	if err != nil {
		fmt.Printf("Error finalizing packet: %v\n", err)
		return
	}
	fmt.Printf("16. Finalized PSBT\n")

	var btcPsbtBuf1 bytes.Buffer
	err = finalPsbt.Serialize(&btcPsbtBuf1)
	if err != nil {
		fmt.Printf("Failed to serialize btcPsbtBuf1: %v", err)
		return
	}

	var btcPsbtBuf bytes.Buffer
	err = signedPkt1.Serialize(&btcPsbtBuf)
	if err != nil {
		fmt.Printf("Failed to serialize btcPsbt: %v", err)
		return
	}

	signedVpsbtBytes1, err := tappsbt.Encode(bobVPsbt)
	if err != nil {
		fmt.Printf("Error encoding bobVPsbt: %v\n", err)
		return
	}
	return &types.BobBizResp{
		Bobpsbt:           btcPsbtBuf.Bytes(),
		FinalPsbt:         btcPsbtBuf1.Bytes(),
		Vpsbt:             signedVpsbtBytes1,
		BobScriptKeyBytes: bobScriptKey.PubKey.SerializeCompressed(),
		AssetID:           req.AssetID,
		NumUnits:          req.NumUnits,
		CommitResp:        resp1,
	}
}

func AliceLogAndPublish(ctx context.Context, req *types.AlicePublishReq) (resp *types.AlicePublishResp, err error) {
	resp1, err := logAndPublish(req.Psbt, []*tappsbt.VPacket{req.Tappsbt}, nil, req.CommitResp)
	if err != nil {
		fmt.Printf("Error logging and publishing: %v\n", err)
		return nil, err
	}
	fmt.Printf("resp1: %v\n", resp1)
	fmt.Printf("17. Logged and published transaction\n")
	fmt.Printf("Creating Bitcoin Core client...\n")

	walletName := "wlt"
	connCfg := &rpcclient.ConnConfig{
		Host:         "132.232.109.84:18443/wallet/" + walletName,
		User:         "rpcuser",
		Pass:         "rpcpassword",
		HTTPPostMode: true,
		DisableTLS:   true,
	}
	coreClient, err := rpcclient.New(connCfg, nil)
	if err != nil {
		fmt.Printf("Error creating new coreClient: %v\n", err)
		return nil, err
	}
	defer coreClient.Shutdown()
	fmt.Printf("Successfully created Bitcoin Core client\n")

	fmt.Printf("Mining new block...\n")
	_, err = MineBlocks(coreClient, 1, 1)
	if err != nil {
		fmt.Printf("Error mining blocks: %v\n", err)
		return nil, err
	}
	fmt.Printf("Successfully mined new block\n")

	fmt.Printf("Getting asset info for assetID: %s\n", req.AssetIDstr)
	assetInfo := GetAssetInfo(req.AssetIDstr)
	if assetInfo == nil {
		fmt.Printf("Asset info is nil for assetID: %s\n", req.AssetIDstr)
		return nil, err
	}
	genInfo := assetInfo.AssetGenesis
	group := assetInfo.AssetGroup
	fmt.Printf("Successfully retrieved asset info\n")

	return &types.AlicePublishResp{
		GenInfo: genInfo,
		Group:   group,
	}, nil
}

func encodePSBTtoBase64(finalPsbt *psbt.Packet) (string, error) {
	var buf bytes.Buffer
	err := finalPsbt.Serialize(&buf)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func waitForNTxsInMempool(miner *rpcclient.Client, n int, timeout time.Duration) ([]*chainhash.Hash, error) {
	breakTimeout := time.After(timeout)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var err error
	var mempool []*chainhash.Hash
	for {
		select {
		case <-breakTimeout:
			return nil, fmt.Errorf("wanted %v, found %v txs in mempool: %v", n, len(mempool), mempool)
		case <-ticker.C:
			mempool, err = miner.GetRawMempool()
			if err != nil {
				return nil, err
			}

			if len(mempool) == n {
				return mempool, nil
			}
		}
	}
}

func MineBlocks(client *rpcclient.Client, num uint32, numTxs int) ([]*wire.MsgBlock, error) {
	const minerMempoolTimeout = 50 * time.Second
	var txids []*chainhash.Hash
	var err error
	if numTxs > 0 {
		txids, err = waitForNTxsInMempool(client, numTxs, minerMempoolTimeout)
		if err != nil {
			return nil, fmt.Errorf("unable to find txns in mempool: %v", err)
		}
	}

	blocks := make([]*wire.MsgBlock, num)

	backend, err := client.BackendVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get backend version: %v", err)
	}

	var blockHashes []*chainhash.Hash

	regtestMiningAddr, err := client.GetNewAddress("")
	if err != nil {
		log.Fatalf("Error getting new address: %v", err)
	}

	switch backend.(type) {
	case *rpcclient.BitcoindVersion:
		addr, err := btcutil.DecodeAddress(regtestMiningAddr.String(), &chaincfg.RegressionNetParams)
		if err != nil {
			return nil, fmt.Errorf("failed to decode address: %v", err)
		}

		blockHashes, err = client.GenerateToAddress(int64(num), addr, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to generate blocks to address: %v", err)
		}

	case rpcclient.BtcdVersion:
		blockHashes, err = client.Generate(num)
		if err != nil {
			return nil, fmt.Errorf("failed to generate blocks: %v", err)
		}

	default:
		return nil, fmt.Errorf("unknown chain backend: %v", backend)
	}

	for i, blockHash := range blockHashes {
		block, err := client.GetBlock(blockHash)
		if err != nil {
			return nil, fmt.Errorf("unable to get block: %v", err)
		}
		blocks[i] = block
	}

	for _, txid := range txids {
		err := AssertTxInBlock(blocks[0], txid)
		if err == nil {
			return nil, fmt.Errorf("transaction not found in block: %v", err)
		}
	}

	return blocks, nil
}

func AssertTxInBlock(block *wire.MsgBlock,
	txid *chainhash.Hash) *wire.MsgTx {

	for _, tx := range block.Transactions {
		sha := tx.TxHash()
		if bytes.Equal(txid[:], sha[:]) {
			return tx
		}
	}

	return nil
}

func logAndPublish(btcPkt *psbt.Packet, activeAssets []*tappsbt.VPacket, passiveAssets []*tappsbt.VPacket,
	commitResp *wrpc.CommitVirtualPsbtsResponse) (*taprpc.SendAssetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var buf bytes.Buffer
	err := btcPkt.Serialize(&buf)
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
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Printf("did not connect: %v\n", err)
		return nil, err
	}
	defer clearUp()
	tapedClient := assetwalletrpc.NewAssetWalletClient(conn)
	resp, err := tapedClient.PublishAndLogTransfer(ctx, request)
	if err != nil {
		log.Fatalf("Publish And LogTransfer failed: %v", err)
		return nil, err
	}
	return resp, nil
}

func finalizePacket(pkt *psbt.Packet) (*psbt.Packet, error) {
	var buf bytes.Buffer
	err := pkt.Serialize(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize PSBT: %v", err)
	}
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to lnd: %v", err)
	}
	defer clearUp()
	client := walletrpc.NewWalletKitClient(conn)
	req := &walletrpc.FinalizePsbtRequest{
		FundedPsbt: buf.Bytes(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	finalizeResp, err := client.FinalizePsbt(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("FinalizePsbt failed: %v", err)
	}
	signedPacket, err := psbt.NewFromRawBytes(
		bytes.NewReader(finalizeResp.SignedPsbt), false,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse finalized PSBT: %v", err)
	}
	return signedPacket, nil
}

func signPsbt(anchorPsbt []byte) (*walletrpc.SignPsbtResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		fmt.Printf("did not connect: %v\n", err)
		return nil, err
	}
	defer clearUp()
	client := walletrpc.NewWalletKitClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	req1 := &walletrpc.SignPsbtRequest{
		FundedPsbt: anchorPsbt}
	signPsbt, err := client.SignPsbt(ctx, req1)
	if err != nil {
		return nil, err
	}
	return signPsbt, nil
}

func getAddressBip32Derivation(addr string) (*psbt.Bip32Derivation, *psbt.TaprootBip32Derivation, error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		fmt.Printf("did not connect: %v\n", err)
		return nil, nil, fmt.Errorf("failed to list addresses: %v", err)
	}
	defer clearUp()
	client := walletrpc.NewWalletKitClient(conn)
	ctxb := context.Background()

	req := &walletrpc.ListAddressesRequest{}
	addresses, err := client.ListAddresses(ctxb, req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list addresses: %v", err)
	}

	var (
		path        []uint32
		pubKeyBytes []byte
	)

	for _, account := range addresses.AccountWithAddresses {
		for _, address := range account.Addresses {
			if address.Address == addr {
				path, err = lntest.ParseDerivationPath(address.DerivationPath)
				if err != nil {
					return nil, nil, fmt.Errorf("failed to parse derivation path for address %s: %v", addr, err)
				}
				pubKeyBytes = address.PublicKey
			}
		}
	}

	if len(path) != 5 || len(pubKeyBytes) == 0 {
		return nil, nil, fmt.Errorf("derivation path for address %s not found or invalid", addr)
	}

	path[0] += hdkeychain.HardenedKeyStart
	path[1] += hdkeychain.HardenedKeyStart
	path[2] += hdkeychain.HardenedKeyStart

	return &psbt.Bip32Derivation{
			PubKey:    pubKeyBytes,
			Bip32Path: path,
		}, &psbt.TaprootBip32Derivation{
			XOnlyPubKey: pubKeyBytes[1:], // 去掉第一个字节，用于 Taproot
			Bip32Path:   path,
		}, nil
}

func BobSendProof(ctx context.Context, req *types.BobSendProofReq) error {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Printf("did not connect: %v\n", err)
		return nil
	}
	defer clearUp()
	tapedClient := taprpc.NewTaprootAssetsClient(conn)
	finalPsbt, err := psbt.NewFromRawBytes(bytes.NewReader(req.FinalPsbt), false)
	if err != nil {
		fmt.Printf("Error parsing final PSBT: %v\n", err)
		return err
	}
	bobOutputIndex := uint32(1)
	transferTXID := finalPsbt.UnsignedTx.TxHash()
	bobAssetOutpoint := fmt.Sprintf("%s:%d", transferTXID.String(), bobOutputIndex)
	transferProofUniRPC(ctx, req.BobScriptKeyBytes, req.GenInfo, req.Group, bobAssetOutpoint)
	registerResp, err := tapedClient.RegisterTransfer(ctx, &taprpc.RegisterTransferRequest{
		AssetId:   req.AssetId[:],
		GroupKey:  req.Group.TweakedGroupKey,
		ScriptKey: req.BobScriptKeyBytes,
		Outpoint: &taprpc.OutPoint{
			Txid:        transferTXID[:],
			OutputIndex: bobOutputIndex,
		},
	})
	if err != nil {
		fmt.Printf("Error registering transfer: %v\n", err)
		return err
	}
	fmt.Printf("Register transfer response: %v\n", registerResp)
	fmt.Printf("Verifying asset transfer...\n")
	verified, err := VerifyAssetTransfer(req.NumUnits)
	if err != nil || !verified {
		fmt.Printf("Asset transfer verification failed: %v\n", err)
		return err
	}
	fmt.Printf("Asset transfer verified successfully\n")

	fmt.Println("PSBT trustless swap completed successfully!")
	return nil
}

func transferProofUniRPC(ctx context.Context, scriptKey []byte, genInfo *taprpc.GenesisInfo, group *taprpc.AssetGroup,
	outpoint string) *unirpc.AssetProofResponse {

	proofFile := ExportProofFileFromUniverse(ctx, genInfo.AssetId, scriptKey, outpoint, group)

	lastProof, err := proofFile.RawLastProof()
	if err != nil {
		return nil
	}

	return InsertProofIntoUniverse(ctx, lastProof)
}

func ExportProofFileFromUniverse(ctx context.Context, assetIDBytes, scriptKey []byte, outpoint string,
	group *taprpc.AssetGroup) *proof.File {

	var assetID asset.ID
	copy(assetID[:], assetIDBytes)

	scriptPubKey, err := btcec.ParsePubKey(scriptKey)
	if err != nil {
		return nil
	}

	op, err := wire.NewOutPointFromString(outpoint)
	if err != nil {
		return nil
	}

	loc := proof.Locator{
		AssetID:   &assetID,
		ScriptKey: *scriptPubKey,
		OutPoint:  op,
	}

	if group != nil {
		groupKey, err := btcec.ParsePubKey(group.TweakedGroupKey)
		if err != nil {
			return nil
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
			return nil, err
		}

		op := &unirpc.Outpoint{
			HashStr: loc.OutPoint.Hash.String(),
			Index:   int32(loc.OutPoint.Index),
		}
		scriptKeyBytes := loc.ScriptKey.SerializeCompressed()

		rawURL := "http://132.232.109.84:8443"

		addr, err := url.Parse(rawURL)
		if err != nil {
			fmt.Printf("Error parsing URL: %v\n", err)
			return nil, err
		}

		macaroon := GetMacaroon("/home/shui/temp/admin.macaroon")
		dialOpts, err := serverDialOpts(macaroon)
		if err != nil {
			return nil, err
		}

		serverAddr := fmt.Sprintf("%s:%s", addr.Hostname(), addr.Port())
		conn, err := grpc.NewClient(serverAddr, dialOpts...)
		if err != nil {
			return nil, err
		}

		client := unirpc.NewUniverseClient(conn)

		uniProof, err := client.QueryProof(ctx, &unirpc.UniverseKey{
			Id: rpcUniID,
			LeafKey: &unirpc.AssetKey{
				Outpoint: &unirpc.AssetKey_Op{
					Op: op,
				},
				ScriptKey: &unirpc.AssetKey_ScriptKeyBytes{
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
			ctx, nil, loc, fetchUniProof,
		)
		return err
	}, 10*time.Minute)
	if err != nil {
		return nil
	}

	return proofFile
}

func InsertProofIntoUniverse(ctx context.Context, proofBytes proof.Blob) *unirpc.AssetProofResponse {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Printf("did not connect: %v\n", err)
		return nil
	}
	defer clearUp()
	tapedClient := taprpc.NewTaprootAssetsClient(conn)
	resp, err := tapedClient.DecodeProof(ctx, &taprpc.DecodeProofRequest{
		RawProof:          proofBytes,
		WithMetaReveal:    true,
		WithPrevWitnesses: true,
	})
	if err != nil {
		return nil
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
			return nil
		}
	}

	rpcUniID, err := taprootassets.MarshalUniID(uniID)
	if err != nil {
		return nil
	}

	dst := unirpc.NewUniverseClient(conn)
	importResp, err := dst.InsertProof(ctx, &unirpc.AssetProof{
		Key: &unirpc.UniverseKey{
			Id: rpcUniID,
			LeafKey: &unirpc.AssetKey{
				Outpoint: &unirpc.AssetKey_OpStr{
					OpStr: rpcAnchor.AnchorOutpoint,
				},
				ScriptKey: &unirpc.AssetKey_ScriptKeyBytes{
					ScriptKeyBytes: rpcAsset.ScriptKey,
				},
			},
		},
		AssetLeaf: &unirpc.AssetLeaf{
			Proof: proofBytes,
		},
	})
	if err != nil {
		return nil
	}

	return importResp
}
