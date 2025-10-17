package api

import (
	"context"

	"github.com/lightninglabs/lightning-terminal/litrpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/service/apiConnect"
)

func PcGetLndState() (*lnrpc.GetStateResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", true)
	if err != nil {
		logrus.Println(errors.Wrap(err, "apiConnect.GetConnection"))
	}
	defer clearUp()
	sc := lnrpc.NewStateClient(conn)
	request := &lnrpc.GetStateRequest{}
	return sc.GetState(context.Background(), request)
}

func PcGetLitSubServersStatus() (*litrpc.SubServerStatusResp, error) {
	conn, clearUp, err := apiConnect.GetConnection("litd", true)
	if err != nil {
		logrus.Println(errors.Wrap(err, "apiConnect.GetConnection"))
	}
	defer clearUp()
	sc := litrpc.NewStatusClient(conn)
	request := &litrpc.SubServerStatusReq{}
	return sc.SubServerStatus(context.Background(), request)
}
