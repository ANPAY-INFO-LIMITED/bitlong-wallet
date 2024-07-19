package api

import (
	"context"
	"fmt"
	"github.com/lightningnetwork/lnd/lnrpc/walletrpc"
	"github.com/wallet/service/apiConnect"
	"github.com/wallet/service/rpcclient"
	"strconv"
)

// ListAccounts
//
//	@Description: ListAddresses retrieves all the addresses along with their balance.
//	An account name filter can be provided to filter through all the wallet accounts
//	and return the addresses of only those matching.
//	@return string
func listAccounts() (*walletrpc.ListAccountsResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
	}
	defer clearUp()
	client := walletrpc.NewWalletKitClient(conn)
	request := &walletrpc.ListAccountsRequest{}
	response, err := client.ListAccounts(context.Background(), request)
	return response, err
}

func ListAddresses() string {
	response, err := rpcclient.ListAddresses()
	if err != nil {
		fmt.Printf("%s walletrpc ListAddresses err: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func GetAllDefaultAddresses() ([]string, error) {
	var result []string
	listAddress, err := ListAddressesAndGetResponse()
	if err != nil {
		return nil, err
	}
	for _, accountWithAddresse := range listAddress.AccountWithAddresses {
		if accountWithAddresse.Name == "default" {
			addresses := accountWithAddresse.Addresses
			for _, address := range addresses {
				result = append(result, address.Address)
			}
		}
	}
	return result, nil
}

func IsIncludeAddress(addresses []string, address string) bool {
	for _, _address := range addresses {
		if _address == address {
			return true
		}
	}
	return false
}

func ListAddressesAndGetResponse() (*walletrpc.ListAddressesResponse, error) {
	return rpcclient.ListAddresses()
}

func ListAccounts() string {
	response, err := listAccounts()
	if err != nil {
		fmt.Printf("%s watchtowerrpc ListAccounts err: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)

}

func FindAccount(name string) string {
	response, err := listAccounts()
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	var accounts []*walletrpc.Account
	for _, account := range response.Accounts {
		if account.Name == name {
			accounts = append(accounts, account)
		}
	}
	if len(accounts) > 0 {
		return MakeJsonErrorResult(SUCCESS, "", accounts)
	}
	return MakeJsonErrorResult(DefaultErr, "account not found", nil)
}

// ListLeases
//
//	@Description: ListLeases lists all currently locked utxos.
//	@return string
func ListLeases() string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
	}
	defer clearUp()
	client := walletrpc.NewWalletKitClient(conn)
	request := &walletrpc.ListLeasesRequest{}
	response, err := client.ListLeases(context.Background(), request)
	if err != nil {
		fmt.Printf("%s watchtowerrpc ListLeases err: %v\n", GetTimeNow(), err)
		return ""
	}
	return response.String()
}

// ListSweeps
//
//	@Description: ListSweeps returns a list of the sweep transactions our node has produced.
//	Note that these sweeps may not be confirmed yet, as we record sweeps on broadcast, not confirmation.
//	@return string
func ListSweeps() string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
	}
	defer clearUp()
	client := walletrpc.NewWalletKitClient(conn)
	request := &walletrpc.ListSweepsRequest{}
	response, err := client.ListSweeps(context.Background(), request)
	if err != nil {
		fmt.Printf("%s watchtowerrpc ListSweeps err: %v\n", GetTimeNow(), err)
		return ""
	}
	return response.String()
}

func ListUnspentAndGetResponse() (*walletrpc.ListUnspentResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
	}
	defer clearUp()
	client := walletrpc.NewWalletKitClient(conn)
	request := &walletrpc.ListUnspentRequest{}
	return client.ListUnspent(context.Background(), request)
}

type ListUnspentUtxo struct {
	AddressType   string `json:"address_type"`
	Address       string `json:"address"`
	AmountSat     int    `json:"amount_sat"`
	PkScript      string `json:"pk_script"`
	Outpoint      string `json:"outpoint"`
	Confirmations int    `json:"confirmations"`
	Time          int    `json:"time"`
}

func ListUnspentResponseToListUnspentUtxos(listUnspentResponse *walletrpc.ListUnspentResponse) *[]ListUnspentUtxo {
	var listUnspentUtxos []ListUnspentUtxo
	for _, utxo := range listUnspentResponse.Utxos {
		listUnspentUtxos = append(listUnspentUtxos, ListUnspentUtxo{
			AddressType:   utxo.AddressType.String(),
			Address:       utxo.Address,
			AmountSat:     int(utxo.AmountSat),
			PkScript:      utxo.PkScript,
			Outpoint:      utxo.Outpoint.TxidStr + ":" + strconv.Itoa(int(utxo.Outpoint.OutputIndex)),
			Confirmations: int(utxo.Confirmations),
			Time:          0,
		})
	}
	return &listUnspentUtxos
}

func ListUnspentUtxoFilterByDefaultAddress(utxos *[]ListUnspentUtxo) *[]ListUnspentUtxo {
	var listUnspentUtxos []ListUnspentUtxo
	addresses, err := GetAllDefaultAddresses()
	if err != nil {
		return utxos
	}
	for _, utxo := range *utxos {
		address := utxo.Address
		if IsIncludeAddress(addresses, address) {
			listUnspentUtxos = append(listUnspentUtxos, utxo)
		}
	}
	return &listUnspentUtxos
}

func ListUnspentAndProcess(token string) (*[]ListUnspentUtxo, error) {
	response, err := ListUnspentAndGetResponse()
	if err != nil {
		return nil, err
	}
	btcUtxos := ListUnspentResponseToListUnspentUtxos(response)
	btcUtxos = ListUnspentUtxoFilterByDefaultAddress(btcUtxos)
	btcUtxos, err = GetTimeForListUnspentUtxoByBitcoind(token, btcUtxos)
	if err != nil {
		return nil, err
	}
	return btcUtxos, nil
}

func BtcUtxos(token string) string {
	response, err := ListUnspentAndProcess(token)
	if err != nil {
		return MakeJsonErrorResult(ListUnspentAndGetResponseErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, SuccessError, response)
}

// ListUnspent
// @Description: ListUnspent returns a list of all utxos spendable by the wallet
// with a number of confirmations between the specified minimum and maximum.
// By default, all utxos are listed. To list only the unconfirmed utxos, set the unconfirmed_only to true.
func ListUnspent() string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
	}
	defer clearUp()
	client := walletrpc.NewWalletKitClient(conn)
	request := &walletrpc.ListUnspentRequest{}
	response, err := client.ListUnspent(context.Background(), request)
	if err != nil {
		fmt.Printf("%s watchtowerrpc ListUnspent err: %v\n", GetTimeNow(), err)
		return ""
	}
	return response.String()
}

// NextAddr
// @Description: NextAddr returns the next unused address within the wallet.
func NextAddr() string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
	}
	defer clearUp()
	client := walletrpc.NewWalletKitClient(conn)
	request := &walletrpc.AddrRequest{}
	response, err := client.NextAddr(context.Background(), request)
	if err != nil {
		fmt.Printf("%s watchtowerrpc NextAddr err: %v\n", GetTimeNow(), err)
		return ""
	}
	return response.String()
}

func BumpFee(txId string, fee int) string {
	if txId == "" || fee == 0 {
		return MakeJsonErrorResult(RequestError, ErrMsgMap[RequestError].Error()+":txId or fee is empty", nil)
	}
	_, err := rpcclient.BumpFee(txId, fee)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", nil)
}
