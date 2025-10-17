package api

import (
	"context"
	"strings"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/pkg/errors"
	"github.com/wallet/service/apiConnect"
	"golang.org/x/exp/rand"
)

func PcUnlockWallet(password string) error {
	conn, clearUp, err := apiConnect.GetConnection("lnd", true)
	if err != nil {
		return errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()
	wuc := lnrpc.NewWalletUnlockerClient(conn)
	request := &lnrpc.UnlockWalletRequest{
		WalletPassword: []byte(password),
	}
	_, err = wuc.UnlockWallet(context.Background(), request)
	if err != nil {
		return errors.Wrap(err, "wuc.UnlockWallet")
	}
	return nil
}

func PcCreateWallet(password string) (string, error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", true)
	if err != nil {
		return "", errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()
	wuc := lnrpc.NewWalletUnlockerClient(conn)
	seedEntropy := make([]byte, 16)
	_, err = rand.Read(seedEntropy)
	if err != nil {
		return "", errors.Wrap(err, "rand.Read")
	}
	genSeedReq := &lnrpc.GenSeedRequest{}
	seedResp, err := wuc.GenSeed(context.Background(), genSeedReq)
	if err != nil {
		return "", errors.Wrap(err, "wuc.GenSeed")
	}

	initReq := &lnrpc.InitWalletRequest{
		WalletPassword:                     []byte(password),
		CipherSeedMnemonic:                 seedResp.CipherSeedMnemonic,
		AezeedPassphrase:                   nil,
		RecoveryWindow:                     0,
		ChannelBackups:                     nil,
		StatelessInit:                      false,
		ExtendedMasterKey:                  "",
		ExtendedMasterKeyBirthdayTimestamp: 0,
	}
	_, err = wuc.InitWallet(context.Background(), initReq)
	if err != nil {
		return "", errors.Wrap(err, "wuc.InitWallet")
	}
	return strings.Join(seedResp.CipherSeedMnemonic, ","), nil
}

func PcRestoreWallet(mnemonic string, password string) (string, error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", true)
	if err != nil {
		return "", errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()
	wuc := lnrpc.NewWalletUnlockerClient(conn)

	mnemonic = strings.TrimSpace(mnemonic)
	mnemonic = strings.ToLower(mnemonic)
	mnemonic = strings.ReplaceAll(mnemonic, " ", ",")
	mnemos := strings.Split(mnemonic, ",")

	initReq := &lnrpc.InitWalletRequest{
		WalletPassword:                     []byte(password),
		CipherSeedMnemonic:                 mnemos,
		AezeedPassphrase:                   nil,
		RecoveryWindow:                     0,
		ChannelBackups:                     nil,
		StatelessInit:                      false,
		ExtendedMasterKey:                  "",
		ExtendedMasterKeyBirthdayTimestamp: 0,
	}
	_, err = wuc.InitWallet(context.Background(), initReq)
	if err != nil {
		return "", errors.Wrap(err, "wuc.InitWallet")
	}
	return strings.Join(mnemos, ","), nil
}
