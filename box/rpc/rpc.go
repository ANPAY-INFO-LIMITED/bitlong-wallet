package rpc

import (
	"encoding/hex"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/wallet/box/models"
)

func ToDecodeAddrResp(r *taprpc.Addr) *models.DecodeAddrResp {
	return &models.DecodeAddrResp{
		Encoded:          r.Encoded,
		AssetID:          hex.EncodeToString(r.AssetId),
		AssetType:        r.AssetType.String(),
		Amount:           int64(r.Amount),
		GroupKey:         hex.EncodeToString(r.GroupKey),
		ScriptKey:        hex.EncodeToString(r.ScriptKey),
		InternalKey:      hex.EncodeToString(r.InternalKey),
		TapscriptSibling: hex.EncodeToString(r.TapscriptSibling),
		TaprootOutputKey: hex.EncodeToString(r.TaprootOutputKey),
		ProofCourierAddr: r.ProofCourierAddr,
		AssetVersion:     r.AssetVersion.String(),
		AddressVersion:   r.AddressVersion.String(),
	}
}
