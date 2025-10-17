package api

import (
	"github.com/pkg/errors"
)

func PcBtcUtxos(token string) (*[]ListUnspentUtxo, error) {
	response, err := ListUnspentAndProcess(token)
	if err != nil {
		return nil, errors.Wrap(err, "ListUnspentAndProcess")
	}
	return response, nil
}

func PcUploadBtcListUnspentUtxos(token string) error {
	err := uploadBtcListUnspentUtxos(token)
	if err != nil {
		return errors.Wrap(err, "uploadBtcListUnspentUtxos")
	}
	return nil
}
