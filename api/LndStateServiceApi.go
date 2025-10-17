package api

import (
	"context"
	"fmt"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/pkg/errors"
	"github.com/wallet/service/apiConnect"
)

func GetStateForSubscribe() bool {
	conn, clearUp, err := apiConnect.GetConnection("lnd", true)
	if err != nil {
		return false
	}
	defer clearUp()
	client := lnrpc.NewStateClient(conn)
	request := &lnrpc.SubscribeStateRequest{}
	response, err := client.SubscribeState(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc SubscribeState err: %v\n", GetTimeNow(), err)
		return false
	}
	fmt.Printf("%s %v\n", GetTimeNow(), response)
	return true
}

func GetState() string {
	response, err := getState()
	if err != nil {
		fmt.Printf("%s watchtowerrpc GetState err: %v\n", GetTimeNow(), err)
		return "NO_START_LND"
	}
	return response.State.String()
}

func getState() (*lnrpc.GetStateResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", true)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()
	client := lnrpc.NewStateClient(conn)
	request := &lnrpc.GetStateRequest{}
	response, err := client.GetState(context.Background(), request)
	return response, err
}
