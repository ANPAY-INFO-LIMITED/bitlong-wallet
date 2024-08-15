package api

import (
	"context"
	"fmt"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/wallet/base"
	"github.com/wallet/service/apiConnect"
	"golang.org/x/exp/rand"
	"os"
	"path/filepath"
	"strings"
)

// GenSeed
//
//	@Description: GenSeed is the first method that should be used to instantiate a new lnd instance.
//	This method allows a caller to generate a new aezeed cipher seed given an optional passphrase.
//	If provided, the passphrase will be necessary to decrypt the cipherseed to expose the internal wallet seed.
//	Once the cipherseed is obtained and verified by the user, the InitWallet method should be used to commit the newly generated seed, and create the wallet.
//	@return string
func GenSeed() string {
	return genSeed()
}

// InitWallet
//
//	@Description:InitWallet is used when lnd is starting up for the first time to fully initialize the daemon and its internal wallet. At the very least a wallet password must be provided.
//	This will be used to encrypt sensitive material on disk.
//	In the case of a recovery scenario, the user can also specify their aezeed mnemonic and passphrase.
//	If set, then the daemon will use this prior state to initialize its internal wallet.
//	Alternatively, this can be used along with the GenSeed RPC to obtain a seed, then present it to the user.
//	Once it has been verified by the user, the seed can be fed into this RPC in order to commit the new wallet.
//	@return bool
func InitWallet(seed, password string) bool {
	return initWallet(seed, password)
}

// UnlockWallet
//
//	@Description: UnlockWallet is used at startup of lnd to provide a password to unlock the wallet database.
//	@return bool
func UnlockWallet(password string) bool {
	return unlockWallet(password)
}

// ChangePassword
//
//	@Description:ChangePassword changes the password of the encrypted wallet.
//	This will automatically unlock the wallet database if successful.
//	@return bool
func ChangePassword(currentPassword, newPassword string) bool {
	return changePassword(currentPassword, newPassword)
}

func RecoverWallet(password, mnemonic string) string {
	err := recoverWallet(password, mnemonic, "")
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), "")
	}
	return MakeJsonErrorResult(SUCCESS, "", "")
}

func genSeed() string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", true)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
	}
	defer clearUp()
	client := lnrpc.NewWalletUnlockerClient(conn)
	//passphrase := ""
	//var aezeedPassphrase = []byte(passphrase)
	seedEntropy := make([]byte, 16)
	_, err = rand.Read(seedEntropy)
	if err != nil {
		fmt.Printf("%s could not generate seed entropy: %v\n", GetTimeNow(), err)
	}
	request := &lnrpc.GenSeedRequest{
		//AezeedPassphrase: aezeedPassphrase,
		//SeedEntropy:      seedEntropy,
	}
	response, err := client.GenSeed(context.Background(), request)
	if err != nil {
		fmt.Printf("%s Error calling GenSeed: %v\n", GetTimeNow(), err)
	}
	return strings.Join(response.CipherSeedMnemonic, ",")
}

func initWallet(seed, password string) bool {
	conn, clearUp, err := apiConnect.GetConnection("lnd", true)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
	}
	defer clearUp()

	var (
		cipherSeedMnemonic      []string
		aezeedPass              []byte
		extendedRootKey         string
		extendedRootKeyBirthday uint64
		recoveryWindow          int32
	)

	client := lnrpc.NewWalletUnlockerClient(conn)
	//seedrequest := &lnrpc.GenSeedRequest{}
	//seedresponse, err := client.GenSeed(context.Background(), seedrequest)
	//cipherSeedMnemonic = seedresponse.CipherSeedMnemonic
	//
	//recoveryWindow = 2500
	cipherSeedMnemonic = strings.Split(seed, ",")
	request := &lnrpc.InitWalletRequest{
		WalletPassword:                     []byte(password),
		CipherSeedMnemonic:                 cipherSeedMnemonic,
		AezeedPassphrase:                   aezeedPass,
		RecoveryWindow:                     recoveryWindow,
		ChannelBackups:                     nil,
		StatelessInit:                      false,
		ExtendedMasterKey:                  extendedRootKey,
		ExtendedMasterKeyBirthdayTimestamp: extendedRootKeyBirthday,
	}
	response, err := client.InitWallet(context.Background(), request)
	if err != nil {
		fmt.Printf("%s Error calling InitWallet: %v\n", GetTimeNow(), err)
	}
	return writeMacaroon(response.AdminMacaroon)
}

func unlockWallet(password string) bool {
	conn, clearUp, err := apiConnect.GetConnection("lnd", true)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
	}
	defer clearUp()
	client := lnrpc.NewWalletUnlockerClient(conn)
	request := &lnrpc.UnlockWalletRequest{
		WalletPassword: []byte(password),
	}
	_, err = client.UnlockWallet(context.Background(), request)
	if err != nil {
		fmt.Printf("%s did not UnlockWallet: %v\n", GetTimeNow(), err)
		return false
	}
	fmt.Printf("%s unlockSuccess\n", GetTimeNow())
	return true
}

func changePassword(currentPassword, newPassword string) bool {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
	}
	defer clearUp()

	client := lnrpc.NewWalletUnlockerClient(conn)
	request := &lnrpc.ChangePasswordRequest{
		CurrentPassword: []byte(currentPassword),
		NewPassword:     []byte(newPassword),
	}
	_, err = client.ChangePassword(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc ChangePassword err: %v\n", GetTimeNow(), err)
		return false
	}
	fmt.Printf("%s ChangePassword Successfully\n", GetTimeNow())
	return true
}

// abandon bar bicycle license embark keen crime rain suffer nation pill blade dwarf faith play motor meadow power skull cheese follow thunder load sail
func recoverWallet(password, mnemonic, passphrase string) error {
	conn, clearUp, err := apiConnect.GetConnection("lnd", true)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
	}
	defer clearUp()
	client := lnrpc.NewWalletUnlockerClient(conn)

	var (
		cipherSeedMnemonic      []string
		aezeedPass              []byte
		extendedRootKey         string
		extendedRootKeyBirthday uint64
		recoveryWindow          int32
	)
	// We'll trim off extra spaces, and ensure the mnemonic is all
	// lower case, then populate our request.
	mnemonic = strings.TrimSpace(mnemonic)
	mnemonic = strings.ToLower(mnemonic)

	cipherSeedMnemonic = strings.Split(mnemonic, " ")

	fmt.Println()

	if len(cipherSeedMnemonic) != 24 {
		return fmt.Errorf("wrong cipher seed mnemonic "+
			"length: got %v words, expecting %v words",
			len(cipherSeedMnemonic), 24)
	}

	// Additionally, the user may have a passphrase, that will also
	// need to be provided so the daemon can properly decipher the
	// cipher seed.
	aezeedPass = []byte(passphrase)

	recoveryWindow = 2500

	// With either the user's prior cipher seed, or a newly generated one,
	// we'll go ahead and initialize the wallet.
	req := &lnrpc.InitWalletRequest{
		WalletPassword:                     []byte(password),
		CipherSeedMnemonic:                 cipherSeedMnemonic,
		AezeedPassphrase:                   aezeedPass,
		ExtendedMasterKey:                  extendedRootKey,
		ExtendedMasterKeyBirthdayTimestamp: extendedRootKeyBirthday,
		RecoveryWindow:                     recoveryWindow,
	}

	response, err := client.InitWallet(context.Background(), req)
	if err != nil {
		return err
	}
	if !writeMacaroon(response.AdminMacaroon) {
		return fmt.Errorf("write macaroon file failed")
	}
	return nil
}

func writeMacaroon(macaroon []byte) bool {
	newFilePath := filepath.Join(base.Configure("lnd"), "."+"macaroonfile")
	err := os.MkdirAll(newFilePath, os.ModePerm)
	if err != nil {
		fmt.Printf("%s Error calling InitWallet: %v\n", GetTimeNow(), err)
	}
	macaroonPath := filepath.Join(newFilePath, "admin.macaroon")
	f, err := os.Create(macaroonPath)
	if err != nil {
		fmt.Printf("%s Error calling InitWallet: %v\n", GetTimeNow(), err)
		return false
	}
	_, err = f.Write(macaroon)
	if err != nil {
		err := f.Close()
		if err != nil {
			fmt.Printf("%s f Close err: %v\n", GetTimeNow(), err)
			return false
		}
		fmt.Printf("%s Error calling InitWallet: %v\n", GetTimeNow(), err)
		return false
	}
	fmt.Printf("%s successful\n", GetTimeNow())
	err = f.Close()
	if err != nil {
		fmt.Printf("%s Error calling InitWallet: %v\n", GetTimeNow(), err)
		return false
	}
	return true
}
