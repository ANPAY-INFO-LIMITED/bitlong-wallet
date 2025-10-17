package api

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/wallet/service"
)

func PcGenerateKeys(mnemonic string) (string, error) {
	keys, err := service.GenerateKeys(mnemonic, "")
	if err != nil {
		return "", errors.Wrap(err, "service.GenerateKeys")
	}
	publicKeyHex := fmt.Sprintf("%x", keys)
	return publicKeyHex, nil
}

func PcGetNPublicKey() (string, error) {
	_, nPub, err := service.GetPublicKey()
	if err != nil {
		return "", errors.Wrap(err, "service.GetPublicKey")
	}
	return nPub, nil
}

func PcGetPublicKey() (string, error) {
	pb, err := service.GetPublicRawKey()
	if err != nil {
		return "", errors.Wrap(err, "service.GetPublicRawKey")
	}
	return pb, nil
}

func PcGetNBPublicKey() (string, error) {
	_, nPub, err := service.GetNewPublicKey()
	if err != nil {
		return "", errors.Wrap(err, "service.GetNewPublicKey")
	}
	return nPub, nil
}
