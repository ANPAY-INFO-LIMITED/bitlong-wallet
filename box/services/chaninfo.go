package services

import (
	"encoding/json"
	"fmt"

	"github.com/lightninglabs/taproot-assets/rfqmsg"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/pkg/errors"
	"github.com/wallet/box/rpc"
)

type ChanBoxInfo struct {
	BoxStatus          string `json:"box_status"`
	Active             bool   `json:"active"`
	IdentityPubkey     string `json:"identity_pubkey"`
	RemotePubkey       string `json:"remote_pubkey"`
	ChannelPoint       string `json:"channel_point"`
	ChanId             int64  `json:"chan_id"`
	Capacity           int64  `json:"capacity"`
	LocalBalance       int64  `json:"local_balance"`
	RemoteBalance      int64  `json:"remote_balance"`
	Private            bool   `json:"private"`
	Initiator          bool   `json:"initiator"`
	ChanStatusFlags    string `json:"chan_status_flags"`
	CommitmentType     string `json:"commitment_type"`
	PushAmountSat      int64  `json:"push_amount_sat"`
	PeerAlias          string `json:"peer_alias"`
	Memo               string `json:"memo"`
	AssetCapacity      int64  `json:"asset_capacity"`
	AssetId            string `json:"asset_id"`
	AssetLocalBalance  int64  `json:"asset_local_balance"`
	AssetRemoteBalance int64  `json:"asset_remote_balance"`
}

func GetChanBoxInfos() ([]*ChanBoxInfo, error) {
	var l rpc.Ln
	chans, err := l.ListChannels(false, false)
	if err != nil {
		return nil, errors.Wrap(err, "l.ListChannels")
	}
	boxPubkey, err := GetBoxPubkey()
	if err != nil {
		return nil, errors.Wrap(err, "GetBoxPubkey")
	}
	boxStatus, err := GetBoxStatus()
	if err != nil {
		return nil, errors.Wrap(err, "GetBoxStatus")
	}

	infos := make([]*ChanBoxInfo, 0)

	for _, c := range chans.Channels {
		if len(c.CustomChannelData) == 0 {
			infos = append(infos, &ChanBoxInfo{
				BoxStatus:          boxStatus,
				Active:             c.Active,
				IdentityPubkey:     boxPubkey,
				RemotePubkey:       c.RemotePubkey,
				ChannelPoint:       c.ChannelPoint,
				ChanId:             int64(c.ChanId),
				Capacity:           c.Capacity,
				LocalBalance:       c.LocalBalance,
				RemoteBalance:      c.RemoteBalance,
				Private:            c.Private,
				Initiator:          c.Initiator,
				ChanStatusFlags:    c.ChanStatusFlags,
				CommitmentType:     c.CommitmentType.String(),
				PushAmountSat:      int64(c.PushAmountSat),
				PeerAlias:          c.PeerAlias,
				Memo:               c.Memo,
				AssetCapacity:      0,
				AssetId:            "00",
				AssetLocalBalance:  0,
				AssetRemoteBalance: 0,
			})
		} else {
			var result rfqmsg.JsonAssetChannel
			err := json.Unmarshal(c.CustomChannelData, &result)
			if err != nil {
				return nil, errors.Wrap(err, "json.Unmarshal")
			}
			infos = append(infos, &ChanBoxInfo{
				BoxStatus:          boxStatus,
				Active:             c.Active,
				IdentityPubkey:     boxPubkey,
				RemotePubkey:       c.RemotePubkey,
				ChannelPoint:       c.ChannelPoint,
				ChanId:             int64(c.ChanId),
				Capacity:           c.Capacity,
				LocalBalance:       c.LocalBalance,
				RemoteBalance:      c.RemoteBalance,
				Private:            c.Private,
				Initiator:          c.Initiator,
				ChanStatusFlags:    c.ChanStatusFlags,
				CommitmentType:     c.CommitmentType.String(),
				PushAmountSat:      int64(c.PushAmountSat),
				PeerAlias:          c.PeerAlias,
				Memo:               c.Memo,
				AssetCapacity:      int64(result.Capacity),
				AssetId:            result.FundingAssets[0].AssetGenesis.AssetID,
				AssetLocalBalance:  int64(result.LocalBalance),
				AssetRemoteBalance: int64(result.RemoteBalance),
			})
		}
	}
	return infos, nil
}

func GetBoxPubkey() (string, error) {
	var l rpc.Ln
	status, err := l.GetInfo()
	if err != nil {
		return "", errors.Wrap(err, "l.GetInfo")
	}
	return status.IdentityPubkey, nil
}

func GetBoxStatus() (string, error) {
	var l rpc.Ln
	status, err := l.GetState()
	if err != nil {
		return "", errors.Wrap(err, "l.GetState")
	}
	return status.State.String(), nil
}

func BackSatsToSever() error {
	var l rpc.Ln
	txns, err := l.ListChaintxns()
	if err != nil {
		return errors.Wrap(err, "l.ListChaintxns")
	}

	closedChans, err := l.GetClosedChannels()
	if err != nil {
		return errors.Wrap(err, "l.GetClosedChannels")
	}

	targetAddr := "bc1p64h5vmktqqdntq762k2c0vcjzwhkcmm5a5qpaaknmd5m27sm7s7see95d8"
	targetAmount := int64(9600)
	sentTxCount := 0

	for _, tx := range txns.Transactions {
		for _, output := range tx.OutputDetails {
			if output.Address == targetAddr && output.Amount == targetAmount {
				sentTxCount++
			}
		}
	}

	assetChannelCount := 0
	for _, c := range closedChans.Channels {
		if len(c.CustomChannelData) == 0 || c.ClosingTxHash == "0000000000000000000000000000000000000000000000000000000000000000" {
			continue
		}
		var result rfqmsg.JsonAssetChannel
		err := json.Unmarshal(c.CustomChannelData, &result)
		if err != nil {
			return errors.Wrap(err, "json.Unmarshal")
		}
		if result.FundingAssets[0].AssetGenesis.AssetID == "97b98f3c45f926057d430ef71f20a6d3e25d7a00fbd1d7b72b306a49d48c9d8c" {
			assetChannelCount++
		}
	}

	if sentTxCount >= assetChannelCount {
		return nil
	}

	remainingCount := assetChannelCount - sentTxCount

	sentCount := 0
	for _, c := range closedChans.Channels {
		if len(c.CustomChannelData) == 0 || c.ClosingTxHash == "0000000000000000000000000000000000000000000000000000000000000000" {
			continue
		}
		var result rfqmsg.JsonAssetChannel
		err := json.Unmarshal(c.CustomChannelData, &result)
		if err != nil {
			return errors.Wrap(err, "json.Unmarshal")
		}
		if result.FundingAssets[0].AssetGenesis.AssetID == "97b98f3c45f926057d430ef71f20a6d3e25d7a00fbd1d7b72b306a49d48c9d8c" {
			if sentCount >= remainingCount {
				break
			}

			resp, err := SendBackSats(targetAddr, targetAmount)
			if err != nil {
				return errors.Wrap(err, "SendBackSats")
			}
			fmt.Printf("发送第 %d 笔交易: %v\n", sentCount+1, resp)
			sentCount++
		}
	}

	fmt.Printf("成功发送了 %d 笔交易\n", sentCount)
	return nil
}

func SendBackSats(Addr string, Amount int64) (*lnrpc.SendCoinsResponse, error) {
	var l rpc.Ln
	resp, err := l.SendCoins(Addr, Amount, 3, false)
	if err != nil {
		return nil, errors.Wrap(err, "l.SendCoins")
	}
	return resp, nil
}

func AssetListBalance() (resp *taprpc.ListBalancesResponse, err error) {
	var t rpc.Tap
	assets, err := t.AssetsListBalances()
	if err != nil {
		return nil, errors.Wrap(err, "l.AssetsListBalances")
	}

	return assets, nil
}

func AssetChanId() (uint64, error) {
	var l rpc.Ln
	channels, err := l.ListChannels(true, true)
	if err != nil {
		return 0, errors.Wrap(err, "l.AssetsListBalances")
	}
	for _, channel := range channels.Channels {
		if channel.CustomChannelData != nil {
			return channel.ChanId, nil
		}
	}
	return 0, nil
}

func PushBackAssetsToServer(addr string) error {
	var t rpc.Tap
	var addrs []string
	err := json.Unmarshal([]byte(addr), &addrs)
	if err != nil {
		return errors.Wrap(err, "json.Unmarshal")
	}

	_, err = t.SendAsset(addrs, 3)
	if err != nil {
		return errors.Wrap(err, "t.SendAsset")
	}
	return nil
}

func TotalReceivedAssetAmount() (int64, error) {
	var l rpc.Ln
	assetId := "97b98f3c45f926057d430ef71f20a6d3e25d7a00fbd1d7b72b306a49d48c9d8c"
	totalReceivedAssetAmount := int64(0)
	amount, err := l.ListAssetPaymentsAmount(assetId)
	if err != nil {
		return 0, errors.Wrap(err, "l.ListAssetPaymentsAmount")
	}
	totalReceivedAssetAmount += amount
	resp, err := l.ListChannelsAll()
	if err != nil {
		return 0, errors.Wrap(err, "l.ListChannelsAll")
	}
	for _, channel := range resp.Channels {
		if len(channel.CustomChannelData) == 0 {
			continue
		}
		var result rfqmsg.JsonAssetChannel
		err := json.Unmarshal(channel.CustomChannelData, &result)
		if err != nil {
			return 0, errors.Wrap(err, "json.Unmarshal")
		}
		if result.FundingAssets[0].AssetGenesis.AssetID == assetId {
			totalReceivedAssetAmount += int64(result.LocalBalance)
		}
	}

	pendingAmount, err := l.BoxPendingChannelsAmount(assetId)
	if err != nil {
		return 0, errors.Wrap(err, "l.BoxPendingChannelsAmount")
	}

	totalReceivedAssetAmount += pendingAmount

	var t rpc.Tap
	assetBalance, err := t.AssetsListBalances()
	if err != nil {
		return 0, errors.Wrap(err, "t.AssetsListBalances")
	}

	totalReceivedAssetAmount += int64(assetBalance.AssetBalances[assetId].Balance)

	return totalReceivedAssetAmount, nil
}
