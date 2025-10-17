package api

import (
	"context"
	"fmt"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/pkg/errors"
	"github.com/wallet/service/apiConnect"
	"gopkg.in/resty.v1"
)

func PcLndStopDaemon() error {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return errors.Wrap(err, "apiConnect.GetConnection")
	}

	defer clearUp()

	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.StopRequest{}
	_, err = client.StopDaemon(context.Background(), request)
	if err != nil {
		return errors.Wrap(err, "client.StopDaemon")
	}
	return nil
}

type GetInfoResp struct {
	Version                   string                    `json:"version"`
	CommitHash                string                    `json:"commit_hash"`
	IdentityPubkey            string                    `json:"identity_pubkey"`
	Alias                     string                    `json:"alias"`
	Color                     string                    `json:"color"`
	NumPendingChannels        uint32                    `json:"num_pending_channels"`
	NumActiveChannels         uint32                    `json:"num_active_channels"`
	NumInactiveChannels       uint32                    `json:"num_inactive_channels"`
	NumPeers                  uint32                    `json:"num_peers"`
	BlockHeight               uint32                    `json:"block_height"`
	BlockHash                 string                    `json:"block_hash"`
	BestHeaderTimestamp       int64                     `json:"best_header_timestamp"`
	SyncedToChain             bool                      `json:"synced_to_chain"`
	SyncedToGraph             bool                      `json:"synced_to_graph"`
	Testnet                   bool                      `json:"testnet"`
	Chains                    []*lnrpc.Chain            `json:"chains"`
	Uris                      []string                  `json:"uris"`
	Features                  map[uint32]*lnrpc.Feature `json:"features"`
	RequireHtlcInterceptor    bool                      `json:"require_htlc_interceptor"`
	StoreFinalHtlcResolutions bool                      `json:"store_final_htlc_resolutions"`
}

func PcLndGetInfo() (*GetInfoResp, error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()
	lc := lnrpc.NewLightningClient(conn)
	request := &lnrpc.GetInfoRequest{}
	resp, err := lc.GetInfo(context.Background(), request)
	if err != nil {
		return nil, errors.Wrap(err, "lc.GetInfo")
	}
	return &GetInfoResp{
		Version:                   resp.Version,
		CommitHash:                resp.CommitHash,
		IdentityPubkey:            resp.IdentityPubkey,
		Alias:                     resp.Alias,
		Color:                     resp.Color,
		NumPendingChannels:        resp.NumPendingChannels,
		NumActiveChannels:         resp.NumActiveChannels,
		NumInactiveChannels:       resp.NumInactiveChannels,
		NumPeers:                  resp.NumPeers,
		BlockHeight:               resp.BlockHeight,
		BlockHash:                 resp.BlockHash,
		BestHeaderTimestamp:       resp.BestHeaderTimestamp,
		SyncedToChain:             resp.SyncedToChain,
		SyncedToGraph:             resp.SyncedToGraph,
		Testnet:                   resp.Testnet,
		Chains:                    resp.Chains,
		Uris:                      resp.Uris,
		Features:                  resp.Features,
		RequireHtlcInterceptor:    resp.RequireHtlcInterceptor,
		StoreFinalHtlcResolutions: resp.StoreFinalHtlcResolutions,
	}, nil
}

func PcGetWalletBalance() (*WalletBalanceResponse, error) {

	response, err := getWalletBalance()
	if err != nil {
		return nil, errors.Wrap(err, "getWalletBalance")
	}
	// @dev: mark imported tap addresses as locked
	response, err = ProcessGetWalletBalanceResult(response)
	if err != nil {
		return nil, errors.Wrap(err, "ProcessGetWalletBalanceResult")
	}
	return &WalletBalanceResponse{
		TotalBalance:       int(response.TotalBalance),
		ConfirmedBalance:   int(response.ConfirmedBalance),
		UnconfirmedBalance: int(response.UnconfirmedBalance),
		LockedBalance:      int(response.LockedBalance),
	}, nil

}

func PcGetNewAddress() (string, error) {

	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return "", errors.Wrap(err, "apiConnect.GetConnection")
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.NewAddressRequest{
		Type: lnrpc.AddressType_TAPROOT_PUBKEY,
	}
	response, err := client.NewAddress(context.Background(), request)
	if err != nil {
		return "", errors.Wrap(err, "client.NewAddress")
	}
	return response.Address, nil

}

func PcSendCoins(addr string, amount int64, feeRate int64, sendAll bool) (string, error) {
	if feeRate > 500 {
		return "", errors.New("fee rate exceeds max(500)")
	}
	response, err := sendCoins(addr, amount, uint64(feeRate), sendAll)
	if err != nil {
		return "", errors.Wrap(err, "sendCoins")
	}
	return response.Txid, nil
}

func PcMergeUTXO(feeRate int64) (string, error) {
	if feeRate > 500 {
		return "", errors.New("fee rate exceeds max(500)")
	}
	//创建一个地址
	addrTarget := GetNewAddress()
	//将所有utxo合并到这个地址
	response, err := sendCoins(addrTarget, 0, uint64(feeRate), true)
	if err != nil {
		return "", errors.Wrap(err, "sendCoins")
	}
	return response.Txid, nil
}

func PcUploadWalletBalance(token string, deviceId string) error {

	host := Cfg.BtlServerHost
	targetUrl := fmt.Sprintf("%s/btc_balance/set", host)

	// 创建 Resty 客户端
	client := resty.New()

	balance, err := PcGetWalletBalance()
	if err != nil {
		return errors.Wrap(err, "PcGetWalletBalance")
	}

	// 设置请求体为 map，Resty 自动设置为 application/json
	body := map[string]any{
		"total_balance":       balance.TotalBalance,
		"confirmed_balance":   balance.ConfirmedBalance,
		"unconfirmed_balance": balance.UnconfirmedBalance,
		"locked_balance":      balance.LockedBalance,
		"device_id":           deviceId,
	}

	var r JsonResult

	_, err = client.R().
		SetAuthToken(token).
		SetBody(body).
		SetResult(&r).
		SetError(&r).
		Post(targetUrl)

	if err != nil {
		return errors.Wrap(err, "client.R.Post")
	}

	if r.Error != "" {
		return errors.New(r.Error)
	}
	return nil
}
