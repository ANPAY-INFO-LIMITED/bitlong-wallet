package services

import (
	"github.com/pkg/errors"
	"github.com/wallet/api"
	"github.com/wallet/box/rpc"
	"github.com/wallet/box/st"
)

func BtcTransferIn() ([]*api.BtcTransferInInfoSimplified, error) {
	err := updateToken()
	if err != nil {
		return nil, errors.Wrap(err, "updateToken")
	}
	token := st.Token()
	var l rpc.Ln
	listAddress, err := l.ListAddresses()
	if err != nil {
		return nil, errors.Wrap(err, "l.ListAddresses")
	}
	txs, err := l.GetTransactions()
	if err != nil {
		return nil, errors.Wrap(err, "l.GetTransactions")
	}
	return api.BtcTransferIn(listAddress, txs, token)
}

func BtcTransferOut() ([]*api.BtcTransferOutInfoSimplified, error) {
	err := updateToken()
	if err != nil {
		return nil, errors.Wrap(err, "updateToken")
	}
	token := st.Token()
	var l rpc.Ln
	listAddress, err := l.ListAddresses()
	if err != nil {
		return nil, errors.Wrap(err, "l.ListAddresses")
	}
	txs, err := l.GetTransactions()
	if err != nil {
		return nil, errors.Wrap(err, "l.GetTransactions")
	}
	return api.BtcTransferOut(listAddress, txs, token)
}

func BtcUtxo() ([]*api.ListUnspentUtxo, error) {
	err := updateToken()
	if err != nil {
		return nil, errors.Wrap(err, "updateToken")
	}
	token := st.Token()
	var l rpc.Ln
	unspent, err := l.ListUnspent()
	if err != nil {
		return nil, errors.Wrap(err, "l.ListUnspent")
	}

	return api.BtcUtxo(unspent, token)
}
