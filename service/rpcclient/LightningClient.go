package rpcclient

import (
	"context"
	"fmt"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/walletrpc"
	"github.com/lightningnetwork/lnd/routing/route"
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

	unspend, err := client.ListUnspent(context.Background(), &walletrpc.ListUnspentRequest{
		UnconfirmedOnly: true,
	})
	if err != nil || unspend == nil {
		return nil, fmt.Errorf("could not get unspent: %v", err)
	}
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
	request := &walletrpc.BumpFeeRequest{
		SatPerVbyte: uint64(fee),
		Immediate:   true,
	}
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
	response, err := client.BumpFee(context.Background(), request)
	if err != nil {
		fmt.Printf("%s watchtowerrpc BumpFee err: %v\n", GetTimeNow(), err)
		return nil, err
	}
	return response, nil
}

func ListSwaps() (*walletrpc.ListSweepsResponse, error) {
	client, clearUp, err := getWalletKitClient()
	if err != nil {
		return nil, err
	}
	defer clearUp()
	response, err := client.ListSweeps(context.Background(), &walletrpc.ListSweepsRequest{})
	if err != nil {
		fmt.Printf("%s watchtowerrpc ListSwaps err: %v\n", GetTimeNow(), err)
		return nil, err
	}
	return response, nil
}

func ListPendingSwaps() (*walletrpc.PendingSweepsResponse, error) {
	client, clearUp, err := getWalletKitClient()
	if err != nil {
		return nil, err
	}
	defer clearUp()
	response, err := client.PendingSweeps(context.Background(), &walletrpc.PendingSweepsRequest{})
	if err != nil {
		fmt.Printf("%s watchtowerrpc ListSwaps err: %v\n", GetTimeNow(), err)
		return nil, err
	}
	return response, nil
}

func ListInvoices(maxNum int, reversed bool) (*lnrpc.ListInvoiceResponse, error) {
	client, clearUp, err := getLightningClient()
	if err != nil {
		return nil, err
	}
	defer clearUp()

	request := &lnrpc.ListInvoiceRequest{
		Reversed: reversed,
	}
	if maxNum > 0 {
		request.NumMaxInvoices = uint64(maxNum)
	}
	response, err := client.ListInvoices(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return response, err
}

func ListChannels(activeOnly bool, private bool, peer string) (*lnrpc.ListChannelsResponse, error) {
	client, clearUp, err := getLightningClient()
	if err != nil {
		return nil, err
	}
	defer clearUp()

	request := &lnrpc.ListChannelsRequest{
		ActiveOnly:  activeOnly,
		PrivateOnly: private,
	}
	var peerKey []byte
	if len(peer) > 0 {
		pk, err := route.NewVertexFromStr(peer)
		if err != nil {
			return nil, fmt.Errorf("invalid --peer pubkey: %w", err)
		}
		peerKey = pk[:]
	}
	if len(peer) > 10 {
		request.Peer = peerKey
	}
	response, err := client.ListChannels(context.Background(), request)
	if err != nil {
		fmt.Printf("%s lnrpc ListChannels err: %v\n", GetTimeNow(), err)
		return nil, err
	}
	return response, nil
}

func AddInvoice(amount int64, memo string, chainId uint64) (*lnrpc.AddInvoiceResponse, error) {
	client, clearUp, err := getLightningClient()
	if err != nil {
		return nil, err
	}
	defer clearUp()

	request := &lnrpc.Invoice{
		Value: amount,
		Memo:  memo,
	}
	if chainId > 0 {
		chanInfo, err := client.GetChanInfo(context.Background(), &lnrpc.ChanInfoRequest{
			ChanId: chainId,
		})
		if err != nil {
			return nil, fmt.Errorf("could not get channel info: %v", err)
		}
		lndInfo, err := client.GetInfo(context.Background(), &lnrpc.GetInfoRequest{})
		if err != nil {
			return nil, fmt.Errorf("could not get lnd info: %v", err)
		}
		var hop *lnrpc.HopHint
		if chanInfo.Node1Pub == lndInfo.IdentityPubkey {
			hop = &lnrpc.HopHint{
				NodeId:                    chanInfo.Node2Pub,
				ChanId:                    chanInfo.ChannelId,
				FeeBaseMsat:               uint32(chanInfo.Node2Policy.FeeBaseMsat),
				FeeProportionalMillionths: uint32(chanInfo.Node2Policy.FeeRateMilliMsat),
				CltvExpiryDelta:           chanInfo.Node2Policy.TimeLockDelta,
			}
		} else {
			hop = &lnrpc.HopHint{
				NodeId:                    chanInfo.Node1Pub,
				ChanId:                    chanInfo.ChannelId,
				FeeBaseMsat:               uint32(chanInfo.Node1Policy.FeeBaseMsat),
				FeeProportionalMillionths: uint32(chanInfo.Node1Policy.FeeRateMilliMsat),
				CltvExpiryDelta:           chanInfo.Node1Policy.TimeLockDelta,
			}
		}
		request.RouteHints = []*lnrpc.RouteHint{
			{
				HopHints: []*lnrpc.HopHint{hop},
			},
		}
	}

	response, err := client.AddInvoice(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return response, nil
}
