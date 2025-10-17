package api

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	"github.com/wallet/service/apiConnect"
	"github.com/wallet/service/rpcclient"
)

func SendPaymentV2(invoice string, amt int, feelimit int, outgoingChanId int, allowSelfPayment bool) string {
	finalResponse, err := rpcclient.SendPaymentV2(invoice, amt, feelimit, uint64(outgoingChanId), allowSelfPayment)
	if err != nil {
		fmt.Printf("%s rpcclient SendPaymentV2 :%v\n", GetTimeNow(), err)
		return MakeJsonErrorResult(SendPaymentV2Err, err.Error(), nil)
	}
	if finalResponse != nil {
		if finalResponse.Status == 2 {
			return MakeJsonErrorResult(SUCCESS, "", nil)
		} else if finalResponse.Status == 3 {
			fmt.Printf("%s %v\n", GetTimeNow(), finalResponse)
			return MakeJsonErrorResult(SendPaymentV2Err, finalResponse.FailureReason.String(), finalResponse.FailureReason)
		}
	}
	fmt.Printf("%s finalResponse is nil,but is not have error\n", GetTimeNow())
	return MakeJsonErrorResult(SendPaymentV2Err, "finalResponse is nil,but is not have error", nil)
}

func TrackPaymentV2(payhash string) string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return MakeJsonErrorResult(GetConnectionErr, err.Error(), nil)
	}
	defer clearUp()
	client := routerrpc.NewRouterClient(conn)
	_payhashByteSlice, _ := hex.DecodeString(payhash)
	request := &routerrpc.TrackPaymentRequest{
		PaymentHash: _payhashByteSlice,
	}
	stream, err := client.TrackPaymentV2(context.Background(), request)

	if err != nil {
		return MakeJsonErrorResult(TrackPaymentV2Err, err.Error(), nil)
	}
	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return MakeJsonErrorResult(streamRecvInfoErr, err.Error(), nil)
			}
			return MakeJsonErrorResult(streamRecvErr, err.Error(), nil)
		}
		status := response.Status.String()
		return MakeJsonErrorResult(SUCCESS, "", status)
	}
}

func SendToRouteV2(payhash []byte, route *lnrpc.Route) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return
	}
	defer clearUp()
	client := routerrpc.NewRouterClient(conn)
	request := &routerrpc.SendToRouteRequest{
		PaymentHash: payhash,
		Route:       route,
	}
	response, err := client.SendToRouteV2(context.Background(), request)
	if err != nil {
		fmt.Printf("%s routerrpc SendToRouteV2 :%v\n", GetTimeNow(), err)
	}
	fmt.Printf("%s %v\n", GetTimeNow(), response)
}

func EstimateRouteFee(dest string, amtsat int64) string {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return ""
	}
	defer clearUp()
	client := routerrpc.NewRouterClient(conn)

	bDest, _ := hex.DecodeString(dest)
	request := &routerrpc.RouteFeeRequest{
		Dest:   bDest,
		AmtSat: amtsat,
	}
	response, err := client.EstimateRouteFee(context.Background(), request)
	if err != nil {
		fmt.Printf("%s routerrpc EstimateRouteFee :%v\n", GetTimeNow(), err)
	}
	fmt.Printf("%s  %v\n", GetTimeNow(), response.RoutingFeeMsat)
	return response.String()
}
