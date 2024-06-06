package rpcclient

import (
	"context"
	"fmt"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/wallet/api/connect"
)

func getLightningClient() (lnrpc.LightningClient, func(), error) {
	conn, clearUp, err := connect.GetConnection("lnd", false)
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
