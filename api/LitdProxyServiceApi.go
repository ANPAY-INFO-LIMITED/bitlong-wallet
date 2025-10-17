package api

import (
	"context"
	"fmt"

	"github.com/lightninglabs/lightning-terminal/litrpc"
	"github.com/pkg/errors"
	"github.com/wallet/service/apiConnect"
)

func LitdStopDaemon() bool {
	_, err := litdStopDaemon()
	if err != nil {
		fmt.Printf("%s litrpc StopRequest err: %v\n", GetTimeNow(), err)
		return false
	}
	return true
}

func litdStopDaemon() (*litrpc.StopDaemonResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("litd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()

	client := litrpc.NewProxyClient(conn)
	request := &litrpc.StopDaemonRequest{}
	response, err := client.StopDaemon(context.Background(), request)
	return response, err
}
