package pcapi

import (
	"github.com/wallet/api"
)

func LnurlGetAvailPort() (int, error) {
	return api.PcLnurlGetAvailPort()
}

func LnurlRunFrpcConf(id string, remotePort string) error {
	return api.PcLnurlRunFrpcConf(id, remotePort)
}

func LnurlRunFrpc() error {
	return api.PcLnurlRunFrpc()
}

func LnurlStopFrpc() {
	api.PcLnurlStopFrpc()
}

func LnurlRequest(id string, name string, localPort string, remotePort string) (string, error) {
	return api.PcLnurlRequest(id, name, localPort, remotePort)
}

func LnurlRequestInvoice(lnu string, invoiceType int, assetID string, amount int, pubkey string, memo string) (string, error) {
	return api.PcLnurlRequestInvoice(lnu, invoiceType, assetID, amount, pubkey, memo)
}
