package api

func lnurlPay(invoiceType int, amount int, invoice string, feeRate int64, token string, deviceID string, assetID string, pubkey string, outgoingChanId int, numSatoshis int, feeLimitSat int) string {
	if invoice == "" {
		return MakeJsonErrorResult(NullInvoice, "invoice is null", "")
	}
	switch InvoiceType(invoiceType) {
	case InvoiceTypeBtcOnChain:
		return lnurlPayBtcOnChain(amount, invoice, feeRate)
	case InvoiceTypeAssetOnChain:
		return lnurlPayAssetOnChain(invoice, feeRate, token, deviceID)
	case InvoiceTypeBtcChannel:
		return lnurlPayBtcChannel(invoice, outgoingChanId, numSatoshis, feeLimitSat)
	case InvoiceTypeAssetChannel:
		return lnurlPayAssetChannel(assetID, pubkey, invoice, 60, outgoingChanId, numSatoshis, feeLimitSat)
	default:
		return MakeJsonErrorResult(lnurlPayErr, "lnurl pay error", "")
	}
}

func lnurlPayBtcOnChain(amount int, invoice string, feeRate int64) string {
	return SendCoins(invoice, int64(amount), feeRate, false)
}

func lnurlPayAssetOnChain(invoice string, feeRate int64, token string, deviceID string) string {
	return SendAssets(invoice, feeRate, token, deviceID)
}

func lnurlPayBtcChannel(invoice string, outgoingChanId int, numSatoshis int, feeLimitSat int) string {
	return SendPaymentV2(invoice, numSatoshis, feeLimitSat, outgoingChanId, false)
}

func lnurlPayAssetChannel(assetId string, pubkey string, paymentReq string, timeoutSeconds int, outgoingChanId int, numSatoshis int, feeLimitSat int) string {
	return AssetChannelSendPayment(assetId, pubkey, paymentReq, timeoutSeconds, outgoingChanId, feeLimitSat, false)
}
