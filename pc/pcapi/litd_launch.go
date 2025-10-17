package pcapi

import (
	"github.com/lightninglabs/lightning-terminal/litrpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/wallet/api"
)

func SetPath(path string, network string) error {
	return api.SetPath(path, network)
}

func StartLitd() {
	api.PcStartLitd()
}

func CreateWallet(password string) (string, error) {
	return api.PcCreateWallet(password)
}

func RestoreWallet(mnemonic string, password string) (string, error) {
	return api.PcRestoreWallet(mnemonic, password)
}

func UnlockWallet(password string) error {
	return api.PcUnlockWallet(password)
}

func GetState() (*lnrpc.GetStateResponse, error) {
	return api.PcGetLndState()
}

func SubServersStatus() (*litrpc.SubServerStatusResp, error) {
	return api.PcGetLitSubServersStatus()
}

func LndGetInfo() (*api.GetInfoResp, error) {
	return api.PcLndGetInfo()
}

func LitdStopDaemon() error {
	return api.PcLitdStopDaemon()
}

func LndStopDaemon() error {
	return api.PcLndStopDaemon()
}
