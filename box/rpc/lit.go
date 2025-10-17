package rpc

import (
	"context"

	"github.com/lightninglabs/lightning-terminal/litrpc"
	"github.com/pkg/errors"
)

type Lit struct{}

func (l Lit) SubServerStatus() (*litrpc.SubServerStatusResp, error) {

	conn, err := GetConn(lit, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	ctx := context.Background()
	sc := litrpc.NewStatusClient(conn)
	req := &litrpc.SubServerStatusReq{}

	resp, err := sc.SubServerStatus(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "sc.SubServerStatus")
	}
	return resp, nil
}

func (l Lit) SubServerStatusWithCtx(ctx context.Context) (*litrpc.SubServerStatusResp, error) {

	conn, err := GetConn(lit, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetConn")
	}

	defer Close(conn)

	sc := litrpc.NewStatusClient(conn)
	req := &litrpc.SubServerStatusReq{}

	resp, err := sc.SubServerStatus(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "sc.SubServerStatus")
	}
	return resp, nil
}
