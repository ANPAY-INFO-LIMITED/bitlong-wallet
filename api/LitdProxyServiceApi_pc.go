package api

import (
	"context"

	"github.com/lightninglabs/lightning-terminal/litrpc"
	"github.com/pkg/errors"
	"github.com/wallet/service/apiConnect"
)

func PcLitdStopDaemon() error {
	conn, clearUp, err := apiConnect.GetConnection("litd", false)
	if err != nil {
		return errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()

	client := litrpc.NewProxyClient(conn)
	request := &litrpc.StopDaemonRequest{}
	_, err = client.StopDaemon(context.Background(), request)
	return err
}
