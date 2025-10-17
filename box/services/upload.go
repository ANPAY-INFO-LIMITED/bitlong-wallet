package services

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/wallet/api"
	"github.com/wallet/box/db"
	"github.com/wallet/box/loggers"
	"github.com/wallet/box/models"
	"github.com/wallet/box/rpc"
	"github.com/wallet/box/sc"
	"gorm.io/gorm"
)

func UploadWalletBalance(t func() string, deviceId string) error {

	token := t()

	host := sc.BaseUrl
	targetUrl := fmt.Sprintf("%s/btc_balance/set", host)

	client := resty.New()

	balance, err := rpc.GetWalletBalance()
	if err != nil {
		return errors.Wrap(err, "PcGetWalletBalance")
	}

	body := map[string]any{
		"total_balance":       balance.TotalBalance,
		"confirmed_balance":   balance.ConfirmedBalance,
		"unconfirmed_balance": balance.UnconfirmedBalance,
		"locked_balance":      balance.LockedBalance,
		"device_id":           deviceId,
	}

	var r models.JResult
	var e models.ErrResp

	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(token).
		SetBody(body).
		SetResult(&r).
		SetError(&e).
		Post(targetUrl)

	if err != nil {
		return errors.Wrap(err, "client.R.Post")
	}

	if e.Error != "" {
		return errors.New(fmt.Sprintf("error: %s", e.Error))
	}

	if r.Error != "" {
		return errors.New(fmt.Sprintf("error: %s", r.Error))
	}

	return nil
}

func UploadAssetManagedUtxos(t func() string, deviceId string) error {
	token := t()
	managedUtxos, err := rpc.GetListUtxos(token)
	if err != nil {
		return errors.Wrap(err, "GetListUtxos")
	}
	assetManagedUtxoSetRequests := api.ManagedUtxosToAssetManagedUtxoSetRequests(deviceId, managedUtxos)

	host := sc.BaseUrl
	targetUrl := fmt.Sprintf("%s/asset_managed_utxo/set", host)

	client := resty.New()

	var r models.JResult
	var e models.ErrResp

	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(token).
		SetBody(assetManagedUtxoSetRequests).
		SetResult(&r).
		SetError(&e).
		Post(targetUrl)

	if err != nil {
		return errors.Wrap(err, "client.R.Post")
	}

	if e.Error != "" {
		return errors.New(fmt.Sprintf("error: %s", e.Error))
	}

	if r.Error != "" {
		return errors.New(fmt.Sprintf("error: %s", r.Error))
	}

	return nil

}

func UploadBtcListUnspent(t func() string) error {
	token := t()
	utxos, err := rpc.GetListUnspent()
	if err != nil {
		return errors.Wrap(err, "GetListUnspent")
	}

	host := sc.BaseUrl
	targetUrl := fmt.Sprintf("%s/btc_utxo/set", host)

	client := resty.New()

	var r models.JResult2
	var e models.ErrResp

	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(token).
		SetBody(utxos).
		SetResult(&r).
		SetError(&e).
		Post(targetUrl)

	if err != nil {
		return errors.Wrap(err, "client.R.Post")
	}

	if e.Error != "" {
		return errors.New(fmt.Sprintf("error: %s", e.Error))
	}

	if r.ErrMsg != "" {
		return errors.New(fmt.Sprintf("error: %s", r.ErrMsg))
	}

	return nil
}

func UploadBoxChanInfo(t func() string) error {
	token := t()
	chanInfos, err := GetChanBoxInfos()
	if err != nil {
		return errors.Wrap(err, "GetChanBoxInfos")
	}

	host := sc.BaseUrl
	targetUrl := fmt.Sprintf("%s/box_device/set_channels_info", host)

	client := resty.New()

	var r models.JResult2
	var e models.ErrResp

	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(token).
		SetBody(chanInfos).
		SetResult(&r).
		SetError(&e).
		Post(targetUrl)

	if err != nil {
		return errors.Wrap(err, "client.R.Post")
	}

	if e.Error != "" {
		return errors.New(fmt.Sprintf("error: %s", e.Error))
	}

	if r.ErrMsg != "" {
		return errors.New(fmt.Sprintf("error: %s", r.ErrMsg))
	}

	return nil
}

type BackAssetsRecordRequest struct {
	NpubKey                 string `json:"npub_key"`
	MachineCoding           string `json:"machine_coding"`
	IdentityPubkey          string `json:"identity_pubkey"`
	ChanId                  int64  `json:"chan_id"`
	AssetId                 string `json:"asset_id"`
	AssetAmount             int64  `json:"asset_amount"`
	AssetAddr               string `json:"asset_addr"`
	ServerIdentityPubkey    string `json:"server_identity_pubkey"`
	IsReceiveAsset          bool   `json:"is_receive_asset"`
	KeySendAssetAmount      int64  `json:"key_send_asset_amount"`
	IsSendToActiveAssetChan bool   `json:"is_send_to_active_asset_chan"`
}

func UploadBackAssetsToServer(t func() string) error {
	req := BackAssetsRecordRequest{}
	token := t()
	Balance, err := AssetListBalance()
	if err != nil {
		return errors.Wrap(err, "ListBalance")
	}

	chanId, err := AssetChanId()
	if err != nil {
		return errors.Wrap(err, "AssetChanId")
	}
	assetId := "97b98f3c45f926057d430ef71f20a6d3e25d7a00fbd1d7b72b306a49d48c9d8c"

	for _, asset := range Balance.AssetBalances {
		if string(asset.AssetGenesis.AssetId) == assetId {
			req.AssetAmount = int64(asset.Balance)
		}
	}
	if req.AssetAmount == 0 {
		return errors.New("asset amount is 0, did not send.")
	}

	req.NpubKey = token
	req.AssetId = assetId
	req.ChanId = int64(chanId)

	host := sc.BaseUrl
	targetUrl := fmt.Sprintf("%s/box_device/get_tap_addrs", host)

	client := resty.New()

	var r models.JResult2
	var e models.ErrResp

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(token).
		SetBody(req).
		SetResult(&r).
		SetError(&e).
		Post(targetUrl)

	if err != nil {
		return errors.Wrap(err, "client.R.Post")
	}

	if e.Error != "" {
		return errors.New(fmt.Sprintf("error: %s", e.Error))
	}

	if r.ErrMsg != "" {
		return errors.New(fmt.Sprintf("error: %s", r.ErrMsg))
	}

	if resp.StatusCode() != 200 {
		return errors.New(fmt.Sprintf("error: %s", resp.Status()))
	}

	err = PushBackAssetsToServer(r.Data.(string))
	if err != nil {
		return errors.Wrap(err, "PushBackAssetsToServer")
	}

	return nil
}

type UpdateTotalAssetPushReq struct {
	NpubKey        string `json:"npub_key"`
	IdentityPubkey string `json:"identity_pubkey"`
	NewTotal       int64  `json:"new_total"`
}

func UpdateTotalAssetPush(t func() string) error {
	token := t()
	cpaState, err := getCpaState()
	if err != nil {
		return errors.Wrap(err, "getCpaState")
	}
	loggers.Chan().Printf("cpaState: %d\n", cpaState)
	if cpaState == models.CpaStateExecuted {
		loggers.Chan().Printf("cpaState is executed, return nil\n")
		return nil
	}
	loggers.Chan().Println("cpaState is not executed")
	var req UpdateTotalAssetPushReq
	req.NpubKey = token
	identityPubkey, err := GetBoxPubkey()
	if err != nil {
		return errors.Wrap(err, "GetBoxPubkey")
	}
	req.IdentityPubkey = identityPubkey
	totalReceivedAssetAmount, err := TotalReceivedAssetAmount()
	if err != nil {
		loggers.Chan().Println("UpdateTotalAssetPush TotalReceivedAssetAmount error: ", err)
		return errors.Wrap(err, "TotalReceivedAssetAmount")
	}
	req.NewTotal = totalReceivedAssetAmount
	loggers.Chan().Println("UpdateTotalAssetPush TotalReceivedAssetAmount: ", totalReceivedAssetAmount)

	host := sc.BaseUrl
	targetUrl := fmt.Sprintf("%s/box_device/update_total_asset_push", host)

	client := resty.New()

	var r models.JResult2
	var e models.ErrResp

	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(token).
		SetBody(req).
		SetResult(&r).
		SetError(&e).
		Post(targetUrl)

	if err != nil {
		loggers.Chan().Println("UpdateTotalAssetPush client.R.Post error: ", err)
		return errors.Wrap(err, "client.R.Post")
	}

	if e.Error != "" {
		loggers.Chan().Println("UpdateTotalAssetPush e.Error: ", e.Error)
		return errors.New(fmt.Sprintf("error: %s", e.Error))
	}

	if r.ErrMsg != "" {
		loggers.Chan().Println("UpdateTotalAssetPush r.ErrMsg: ", r.ErrMsg)
		return errors.New(fmt.Sprintf("error: %s", r.ErrMsg))
	}

	return nil
}

func getCpaState() (models.CpaState, error) {
	tx := db.Sqlite().Begin()

	var l models.Cpa
	err := tx.Model(&models.Cpa{}).First(&l).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = tx.Model(&models.Cpa{}).
				Create(&models.Cpa{
					State: models.CpaStateInit,
				}).Error
			if err != nil {
				tx.Rollback()
				return models.CpaStateExecuted, errors.Wrap(err, "tx.Model(&models.Cpa{}).Create")
			}
			if err = tx.Commit().Error; err != nil {
				return models.CpaStateExecuted, errors.Wrap(err, "tx.Commit()")
			}
			return models.CpaStateInit, nil
		} else {
			tx.Rollback()
			return models.CpaStateExecuted, errors.Wrap(err, "tx.Model(&models.Cpa{}).First")
		}
	}

	if l.State == models.CpaStateExecuted {
		if err := tx.Delete(&l).Error; err != nil {
			tx.Rollback()
			return models.CpaStateExecuted, errors.Wrap(err, "tx.Delete(&l)")
		}
		if err = tx.Commit().Error; err != nil {
			return models.CpaStateExecuted, errors.Wrap(err, "tx.Commit()")
		}
		return models.CpaStateInit, nil
	}

	tx.Rollback()
	return l.State, nil
}
