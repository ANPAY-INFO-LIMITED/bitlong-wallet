package pcapi

import (
	"github.com/wallet/api"
)

func GetNewAddressP2tr() (*api.Addr, error) {
	return api.PcGetNewAddress_P2TR()
}

func GetNewAddressP2wkh() (*api.Addr, error) {
	return api.PcGetNewAddress_P2WKH()
}

func GetNewAddressNp2wkh() (*api.Addr, error) {
	return api.PcGetNewAddress_NP2WKH()
}

func GetNewAddressP2trExample() (*api.Addr, error) {
	return api.PcGetNewAddress_P2TR_Example()
}

func GetNewAddressP2wkhExample() (*api.Addr, error) {
	return api.PcGetNewAddress_P2WKH_Example()
}

func GetNewAddressNp2wkhExample() (*api.Addr, error) {
	return api.PcGetNewAddress_NP2WKH_Example()
}

func StoreAddr(name string, address string, balance int, addressType string, derivationPath string, isInternal bool) error {
	return api.PcStoreAddr(name, address, balance, addressType, derivationPath, isInternal)
}

func RemoveAddr(address string) error {
	return api.PcRemoveAddr(address)
}

func QueryAllAddr() ([]*api.Addr, error) {
	return api.PcQueryAllAddr()
}

func UpdateAllAddressesByGnzba() error {
	return api.PcUpdateAllAddressesByGNZBA()
}
