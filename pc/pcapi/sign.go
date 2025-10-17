package pcapi

import "github.com/wallet/api"

func SignSchnorr(hexPrivateKey string, message string) (string, error) {
	return api.SignSchnorr(hexPrivateKey, message)
}
