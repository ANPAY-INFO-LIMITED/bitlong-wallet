package api

import (
	"context"
	"fmt"
	"github.com/lightningnetwork/lnd/lnrpc/walletrpc"
	"github.com/wallet/service/apiConnect"
)

// ListAddress
//
//	@Description: ListAddresses retrieves all the addresses along with their balance.
//	An account name filter can be provided to filter through all the wallet accounts and return the addresses of only those matching.
//	@return string
func listAddresses() (*walletrpc.ListAddressesResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
	}
	defer clearUp()
	client := walletrpc.NewWalletKitClient(conn)
	request := &walletrpc.ListAddressesRequest{}
	response, err := client.ListAddresses(context.Background(), request)
	return response, err
}

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
	response, err := listAddresses()
	if err != nil {
		fmt.Printf("%s walletrpc ListAddresses err: %v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func ListAddressesAndGetResponse() (*walletrpc.ListAddressesResponse, error) {
	return listAddresses()
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

// ListUnspent
//
//	@Description: ListUnspent returns a list of all utxos spendable by the wallet
//	with a number of confirmations between the specified minimum and maximum.
//	By default, all utxos are listed. To list only the unconfirmed utxos, set the unconfirmed_only to true.
//	@return string
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
//
//	@Description: NextAddr returns the next unused address within the wallet.
//	@return string
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
