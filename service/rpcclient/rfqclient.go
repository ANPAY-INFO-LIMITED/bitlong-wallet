package rpcclient

import (
	"context"
	"fmt"
	"github.com/lightninglabs/taproot-assets/taprpc/rfqrpc"
	"github.com/wallet/service/apiConnect"
)

func getRfqClient() (rfqrpc.RfqClient, func(), error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	client := rfqrpc.NewRfqClient(conn)
	return client, clearUp, nil
}

func SubscribeRfqEventNtfns() {
	client, clearUp, err := getRfqClient()
	if err != nil {
		fmt.Println(err)
	}
	defer clearUp()
	request := &rfqrpc.SubscribeRfqEventNtfnsRequest{}

	stream, err := client.SubscribeRfqEventNtfns(context.Background(), request)
	if err != nil {
		fmt.Println(err)
	}
	for {
		event, err := stream.Recv()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(event)
	}
}
