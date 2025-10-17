package types

import (
	"bytes"

	"github.com/btcsuite/btcd/btcec/v2/schnorr/musig2"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/txscript"
	"github.com/lightninglabs/taproot-assets/asset"
	"github.com/lightninglabs/taproot-assets/tappsbt"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/lightninglabs/taproot-assets/taprpc/assetwalletrpc"
	wrpc "github.com/lightninglabs/taproot-assets/taprpc/assetwalletrpc"
)

type SetMusigLockTimeHandlerReq struct {
	AssetIdStr     string                `json:"asset_id"`
	Amount         uint64                `json:"amount"`
	BobScriptKey   *taprpc.ScriptKey     `json:"bob_script_key"`
	BobInternalKey *taprpc.KeyDescriptor `json:"bob_internal_key"`
	LockTime       int64                 `json:"lock_time"`
}

type AliceSetBizReq struct {
	AssetIdBytes   []byte                `json:"asset_id_bytes"`
	Amount         uint64                `json:"amount"`
	BobScriptKey   *taprpc.ScriptKey     `json:"bob_script_key"`
	BobInternalKey *taprpc.KeyDescriptor `json:"bob_internal_key"`
	LockTime       int64                 `json:"lock_time"`
	WithdrawAddr   *taprpc.Addr          `json:"withdraw_addr"`
}

type AliceSetBizResp struct {
	Leaves               []txscript.TapLeaf    `json:"leaves"`
	TapControlBlockBytes []byte                `json:"tap_control_block_bytes"`
	BtcControlBlockBytes []byte                `json:"btc_control_block_bytes"`
	AliceScriptKey       *taprpc.ScriptKey     `json:"alice_script_key"`
	AliceInternalKey     *taprpc.KeyDescriptor `json:"alice_internal_key"`
	BtcTapLeaf           txscript.TapLeaf      `json:"btc_tap_leaf"`
	AliceNonces          *musig2.Nonces        `json:"alice_nonces"`
	BobNonces            *musig2.Nonces        `json:"bob_nonces"`
	LockTime             int64                 `json:"lock_time"`
}

type BobGenerateAddrReq struct {
	AssetId []byte `json:"asset_id"`
	Amount  uint64 `json:"amount"`
}

type BobGenerateAddrResp struct {
	WithdrawAddr   *taprpc.Addr          `json:"withdraw_addr"`
	BobScriptKey   *taprpc.ScriptKey     `json:"bob_script_key"`
	BobInternalKey *taprpc.KeyDescriptor `json:"bob_internal_key"`
}

type FundVirtualPsbtReq struct {
	WithdrawAddr *taprpc.Addr `json:"withdraw_addr"`
}

type FundVirtualPsbtResp struct {
	FundedPsbt []byte `json:"funded_psbt"`
}

type BobPartialSigReq struct {
	AliceScriptKey *taprpc.ScriptKey  `json:"alice_script_key"`
	BobScriptKey   *taprpc.ScriptKey  `json:"bob_script_key"`
	Leaves         []txscript.TapLeaf `json:"leaves"`
	FundedPsbt     []byte             `json:"funded_psbt"`
	AliceNonces    *musig2.Nonces     `json:"alice_nonces"`
	BobNonces      *musig2.Nonces     `json:"bob_nonces"`
}

type BobPartialSigResp struct {
	BobPartialSig *[]byte `json:"bob_partial_sig"`
}

type SubmitVtxReq struct {
	AliceScriptKey       *taprpc.ScriptKey  `json:"alice_script_key"`
	BobScriptKey         *taprpc.ScriptKey  `json:"bob_script_key"`
	Leaves               []txscript.TapLeaf `json:"leaves"`
	TapControlBlockBytes []byte             `json:"tap_control_block_bytes"`
	FundedPsbt           []byte             `json:"funded_psbt"`
	AliceNonces          *musig2.Nonces     `json:"alice_nonces"`
	BobNonces            *musig2.Nonces     `json:"bob_nonces"`
	BobPartialSig        []byte             `json:"bob_partial_sig"`
	LockTime             int64              `json:"lock_time"`
}

type SubmitVtxResp struct {
	BtcWithdrawPkt           string                           `json:"btc_withdraw_pkt"`
	FinalizedWithdrawPackets []byte                           `json:"finalized_withdraw_packets"`
	CommitResp               *wrpc.CommitVirtualPsbtsResponse `json:"commit_resp"`
}

type BtcPartialSigReq struct {
	BtcWithdrawPkt       string                `json:"btc_withdraw_pkt"`
	BobInternalKey       *taprpc.KeyDescriptor `json:"bob_internal_key"`
	BtcControlBlockBytes []byte                `json:"btc_control_block_bytes"`
	BtcTapLeaf           txscript.TapLeaf      `json:"btc_tap_leaf"`
}

type BtcPartialSigResp struct {
	BobBtcPartialSig []byte `json:"bob_btc_partial_sig"`
}

type SignAndFinalizeBtcTransactionReq struct {
	BtcWithdrawPkt           string                           `json:"btc_withdraw_pkt"`
	FinalizedWithdrawPackets []byte                           `json:"finalized_withdraw_packets"`
	BobBtcPartialSig         []byte                           `json:"bob_btc_partial_sig"`
	BtcTapLeaf               txscript.TapLeaf                 `json:"btc_tap_leaf"`
	CommitResp               *wrpc.CommitVirtualPsbtsResponse `json:"commit_resp"`
	BtcControlBlockBytes     []byte                           `json:"btc_control_block_bytes"`
	AliceInternalKey         *taprpc.KeyDescriptor            `json:"alice_internal_key"`
}

type FuncAliceBizReq struct {
	AssetID  asset.ID
	NumUnits uint64
}

type FuncAliceBizResp struct {
	BtcPsbt          *psbt.Packet
	B                bytes.Buffer
	Resp             *wrpc.CommitVirtualPsbtsResponse
	SignedVpsbtBytes []byte
}

type AlicePublishReq struct {
	Psbt              *psbt.Packet
	Tappsbt           *tappsbt.VPacket
	AssetIDstr        string
	BobScriptKeyBytes []byte
	CommitResp        *wrpc.CommitVirtualPsbtsResponse
}

type AlicePublishResp struct {
	ProofFile []byte
	Outpoint  string
	GenInfo   *taprpc.GenesisInfo
	Group     *taprpc.AssetGroup
}

type BobBizReq struct {
	SignedVpsbtBytes []byte
	BtcPsbt          *psbt.Packet
	B                bytes.Buffer
	CommitResp       *assetwalletrpc.CommitVirtualPsbtsResponse
	AssetID          asset.ID
	NumUnits         uint64
}

type BobBizResp struct {
	Bobpsbt           []byte
	FinalPsbt         []byte
	Vpsbt             []byte
	BobScriptKeyBytes []byte
	AssetID           asset.ID
	NumUnits          uint64
	CommitResp        *wrpc.CommitVirtualPsbtsResponse
}

type BobSendProofReq struct {
	AssetId           asset.ID
	NumUnits          uint64
	BobScriptKeyBytes []byte
	FinalPsbt         []byte
	GenInfo           *taprpc.GenesisInfo
	Group             *taprpc.AssetGroup
}
