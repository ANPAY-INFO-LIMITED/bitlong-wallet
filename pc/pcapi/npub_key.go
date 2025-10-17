package pcapi

import (
	"github.com/wallet/api"
	"strings"
)

func GenerateKeys(mnemonic string) (string, error) {
	spaceMnemonic := strings.Replace(mnemonic, ",", " ", -1)
	return api.PcGenerateKeys(spaceMnemonic)
}

func GetPrivateKey() (string, error) {
	return api.GetPrivateKey()
}

func GetNPublicKey() (string, error) {
	return api.PcGetNPublicKey()R
}

func GetPublicKey() (string, error) {
	return api.PcGetPublicKey()
}

func GetNBPublicKey() (string, error) {
	return api.PcGetNBPublicKey()
}
