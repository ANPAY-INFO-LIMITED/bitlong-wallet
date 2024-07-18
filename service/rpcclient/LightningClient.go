package rpcclient

import (
	"context"
	"fmt"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/walletrpc"
	"github.com/wallet/service/apiConnect"
)

func getLightningClient() (lnrpc.LightningClient, func(), error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	client := lnrpc.NewLightningClient(conn)
	return client, clearUp, nil
}

func DecodePayReq(payReq string) (*lnrpc.PayReq, error) {
	client, clearUp, err := getLightningClient()
	if err != nil {
		return nil, err
	}
	defer clearUp()

	request := &lnrpc.PayReqString{
		PayReq: payReq,
	}
	response, err := client.DecodePayReq(context.Background(), request)
	if err != nil {
		fmt.Printf("%s client.DecodePayReq :%v\n", GetTimeNow(), err)
		return nil, err
	}
	return response, nil
}

func getWalletKitClient() (walletrpc.WalletKitClient, func(), error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	client := walletrpc.NewWalletKitClient(conn)
	return client, clearUp, nil
}
func ListAddresses() (*walletrpc.ListAddressesResponse, error) {
	client, clearUp, err := getWalletKitClient()
	if err != nil {
		return nil, err
	}
	defer clearUp()

	request := &walletrpc.ListAddressesRequest{}
	response, err := client.ListAddresses(context.Background(), request)
	return response, err
}

func BumpFee(txId string, fee int) (*walletrpc.BumpFeeResponse, error) {
	client, clearUp, err := getWalletKitClient()
	if err != nil {
		return nil, err
	}
	defer clearUp()

	//获取用户的未确认的utxo
	unspend, err := client.ListUnspent(context.Background(), &walletrpc.ListUnspentRequest{
		UnconfirmedOnly: true,
	})
	if err != nil || unspend == nil {
		return nil, fmt.Errorf("could not get unspent: %v", err)
	}
	//获取用户默认账户的地址
	addResponce, err := client.ListAddresses(context.Background(), &walletrpc.ListAddressesRequest{
		AccountName: "default",
	})
	if err != nil || addResponce == nil {
		return nil, fmt.Errorf("could not get user addresses: %v", err)
	}
	var addrs []string
	for _, account := range addResponce.AccountWithAddresses {
		for _, addr := range account.Addresses {
			addrs = append(addrs, addr.Address)
		}
	}
	//创建请求体
	request := &walletrpc.BumpFeeRequest{
		SatPerVbyte: uint64(fee),
	}
	//找出到可使用的未确认utxo
	for _, u := range unspend.Utxos {
		if u.Outpoint.TxidStr == txId {
			for _, addr := range addrs {
				if u.Address == addr {
					request.Outpoint = u.Outpoint
					break
				}
			}
		}
	}
	if request.Outpoint == nil {
		return nil, fmt.Errorf("could not find usable outpoint for this txid %v", txId)
	}
	//发送请求
	response, err := client.BumpFee(context.Background(), request)
	if err != nil {
		fmt.Printf("%s watchtowerrpc BumpFee err: %v\n", GetTimeNow(), err)
		return nil, err
	}
	return response, nil
}
