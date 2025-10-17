package api

import (
	"github.com/lightninglabs/taproot-assets/tappsbt"
	"github.com/lightninglabs/taproot-assets/taprpc/btlrpc"
	"github.com/wallet/base"
)

func GetListEligibleCoins(assetId string) string {
	coins, err := getListEligibleCoins(assetId)
	if err != nil {
		return MakeJsonErrorResult2(getListEligibleCoinsErr, err.Error(), []*btlrpc.Coin{})
	}
	return MakeJsonErrorResult2(SUCCESS_2, SUCCESS_2.Error(), btlrpcCoinsToGetListEligibleCoinsResults(coins))
}

func CreateSellOrderSign(assetId string, assetNum int64, price int64, anchorPoint string, internalKey string) string {
	resp, err := psbtTrustlessSwapCreateSellOrderSignWithOneFilter(assetId, uint64(assetNum), price, 0, 1, false, anchorPoint, internalKey)
	if err != nil {
		return MakeJsonErrorResult2(psbtTrustlessSwapCreateSellOrderSignWithOneFilterErr, err.Error(), "")
	}

	return MakeJsonErrorResult2(SUCCESS_2, SUCCESS_2.Error(), resp)
}

func BuySOrderSign(signedSellOrder string, feeRate int64) string {
	var orderResp psbtTrustlessSwapCreateSellOrderResponse
	err := orderResp.FromStr(signedSellOrder)
	if err != nil {
		return MakeJsonErrorResult2(FromStrErr, err.Error(), "")
	}
	signedVPsbtBytes, signedPsbtBytes := orderResp.SignedVPsbtBytes, orderResp.SignedPsbtBytes
	deliveryAddr := base.QueryConfigByKey("universeHost")

	resp, err := psbtTrustlessSwapBuySOrderSign(signedVPsbtBytes, signedPsbtBytes, uint64(feeRate), deliveryAddr, false)
	if err != nil {
		return MakeJsonErrorResult2(psbtTrustlessSwapBuySOrderSignErr, err.Error(), "")
	}

	return MakeJsonErrorResult2(SUCCESS_2, SUCCESS_2.Error(), resp)
}

func PublishSOrderTx(signedBoughtSOrder string) string {
	var boughtResp psbtTrustlessSwapBuySOrderResponse
	err := boughtResp.FromStr(signedBoughtSOrder)
	if err != nil {
		return MakeJsonErrorResult2(FromStrErr, err.Error(), "")
	}
	signedPktBytes, bobVPsbtBytes, cvpResp := boughtResp.SignedPktBytes, boughtResp.BobVPsbtBytes, boughtResp.CvpResp

	signedPkt, err := bytesToPsbtPacket(signedPktBytes)
	if err != nil {
		return MakeJsonErrorResult2(bytesToPsbtPacketErr, err.Error(), "")
	}
	bobVPsbt, err := tappsbt.Decode(bobVPsbtBytes)
	if err != nil {
		return MakeJsonErrorResult2(tappsbtDecodeErr, err.Error(), "")
	}

	resp, err := psbtTrustlessSwapPublishTx(signedPkt, bobVPsbt, cvpResp)
	if err != nil {
		return MakeJsonErrorResult2(psbtTrustlessSwapPublishTxErr, err.Error(), "")
	}
	var txHash []byte
	var txid string
	if resp != nil && resp.Transfer != nil {
		txHash = resp.Transfer.AnchorTxHash
		txid = TxHashEncodeToString(txHash)
	}

	return MakeJsonErrorResult2(SUCCESS_2, SUCCESS_2.Error(), txid)
}

func GetLastProof(token string, scriptKey string, outpoint string, assetId string) string {
	lastProofStr, err := requestToGetLastProof(token, scriptKey, outpoint, assetId)
	if err != nil {
		return MakeJsonErrorResult2(requestToGetLastProofErr, err.Error(), "")
	}
	return MakeJsonErrorResult2(SUCCESS_2, SUCCESS_2.Error(), lastProofStr)
}

func AllowFederationSyncInsertAndExport() string {
	return SetFederationSyncConfig(true, true)
}

func InsertProofAndRegisterTransfer(assetId string, signedBoughtSOrder string, lastProofStr string) string {
	lastProof, err := LastProofFromStr(lastProofStr)
	if err != nil {
		return MakeJsonErrorResult2(LastProofFromStrErr, err.Error(), "")
	}

	var boughtResp psbtTrustlessSwapBuySOrderResponse
	err = boughtResp.FromStr(signedBoughtSOrder)
	if err != nil {
		return MakeJsonErrorResult2(FromStrErr, err.Error(), "")
	}
	scriptKeyBytes, assetOutpoint := boughtResp.ScriptKeyBytes, boughtResp.AssetOutpoint

	err = psbtTrustlessSwapBuySOrderProof(scriptKeyBytes, assetOutpoint, assetId, lastProof)
	if err != nil {
		return MakeJsonErrorResult2(psbtTrustlessSwapBuySOrderProofErr, err.Error(), "")
	}
	return MakeJsonErrorResult2(SUCCESS_2, SUCCESS_2.Error(), assetOutpoint)
}
