package api

import (
	"github.com/lightninglabs/taproot-assets/tappsbt"
	"github.com/pkg/errors"
	"github.com/wallet/base"
)

func PcGetListEligibleCoins(assetId string) ([]*EligibleCoin, error) {
	resp, err := getListEligibleCoins(assetId)
	if err != nil {
		return nil, errors.Wrap(err, "getListEligibleCoins")
	}

	return btlrpcCoinsToGetListEligibleCoinsResults(resp), nil
}

func PcCreateSellOrderSign(assetId string, assetNum uint64, price int64, anchorPoint string, internalKey string) (string, error) {
	resp, err := psbtTrustlessSwapCreateSellOrderSignWithOneFilter(assetId, assetNum, price, 0, 1, false, anchorPoint, internalKey)
	if err != nil {
		return "", errors.Wrap(err, "psbtTrustlessSwapCreateSellOrderSignWithOneFilter")
	}
	return resp.String(), nil
}

func PcBuySOrderSign(signedSellOrder string, feeRate uint64) (string, error) {
	var orderResp psbtTrustlessSwapCreateSellOrderResponse
	err := orderResp.FromStr(signedSellOrder)
	if err != nil {
		return "", errors.Wrap(err, "orderResp.FromStr")
	}
	signedVPsbtBytes, signedPsbtBytes := orderResp.SignedVPsbtBytes, orderResp.SignedPsbtBytes
	deliveryAddr := base.QueryConfigByKey("universeHost")

	resp, err := psbtTrustlessSwapBuySOrderSign(signedVPsbtBytes, signedPsbtBytes, feeRate, deliveryAddr, false)
	if err != nil {
		return "", errors.Wrap(err, "psbtTrustlessSwapBuySOrderSign")
	}

	return resp.String(), nil
}

func PcPublishSOrderTx(signedBoughtSOrder string) (string, error) {
	var boughtResp psbtTrustlessSwapBuySOrderResponse
	err := boughtResp.FromStr(signedBoughtSOrder)
	if err != nil {
		return "", errors.Wrap(err, "boughtResp.FromStr")
	}
	signedPktBytes, bobVPsbtBytes, cvpResp := boughtResp.SignedPktBytes, boughtResp.BobVPsbtBytes, boughtResp.CvpResp

	signedPkt, err := bytesToPsbtPacket(signedPktBytes)
	if err != nil {
		return "", errors.Wrap(err, "bytesToPsbtPacket")
	}
	bobVPsbt, err := tappsbt.Decode(bobVPsbtBytes)
	if err != nil {
		return "", errors.Wrap(err, "tappsbt.Decode")
	}

	resp, err := psbtTrustlessSwapPublishTx(signedPkt, bobVPsbt, cvpResp)
	if err != nil {
		return "", errors.Wrap(err, "psbtTrustlessSwapPublishTx")
	}
	var txHash []byte
	var txid string
	if resp != nil && resp.Transfer != nil {
		txHash = resp.Transfer.AnchorTxHash
		txid = TxHashEncodeToString(txHash)
	}

	return txid, nil
}

func PcAllowFederationSyncInsertAndExport() error {
	_, err := setFederationSyncConfig(true, true)
	if err != nil {
		return errors.Wrap(err, "setFederationSyncConfig")
	}
	return nil
}

func PcInsertProofAndRegisterTransfer(assetId string, signedBoughtSOrder string, lastProofStr string) (string, error) {
	lastProof, err := LastProofFromStr(lastProofStr)
	if err != nil {
		return "", errors.Wrap(err, "LastProofFromStr")
	}

	var boughtResp psbtTrustlessSwapBuySOrderResponse
	err = boughtResp.FromStr(signedBoughtSOrder)
	if err != nil {
		return "", errors.Wrap(err, "boughtResp.FromStr")
	}
	scriptKeyBytes, assetOutpoint := boughtResp.ScriptKeyBytes, boughtResp.AssetOutpoint

	err = psbtTrustlessSwapBuySOrderProof(scriptKeyBytes, assetOutpoint, assetId, lastProof)
	if err != nil {
		return "", errors.Wrap(err, "psbtTrustlessSwapBuySOrderProof")
	}
	return assetOutpoint, nil
}
