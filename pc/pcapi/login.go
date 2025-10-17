package pcapi

import "github.com/wallet/api"

func Login(username, password string) (string, error) {
	return api.Login(username, password)
}
