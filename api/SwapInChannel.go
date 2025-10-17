package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lightninglabs/taproot-assets/rfqmsg"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/wallet/base"
	"github.com/wallet/service/rpcclient"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"sync"
)

func SwapExactSendInChannel(sendAssetId string, sendAmountStr string, receiveAssetId string,
	predictAmountStr string, slippageStr string, token string) string {
	sendAmount, err := strconv.ParseInt(sendAmountStr, 10, 64)
	if err != nil {
		return MakeJsonErrorResult(RequestError, "Illegal sendAmount:"+err.Error(), nil)
	}
	predictAmount, err := strconv.ParseInt(predictAmountStr, 10, 64)
	if err != nil {
		return MakeJsonErrorResult(RequestError, "Illegal predictAmount:"+err.Error(), nil)
	}
	slippage, err := strconv.ParseInt(slippageStr, 10, 16)
	if err != nil {
		return MakeJsonErrorResult(RequestError, "Illegal slippage:"+err.Error(), nil)
	}

	err = swapExactSendInChannel(sendAssetId, sendAmount, receiveAssetId, predictAmount, uint16(slippage), token)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", nil)

}
func SwapExactReceiveInChannel(sendAssetId string, sendAmountStr string, receiveAssetId string,
	predictAmountStr string, slippageStr string, token string) string {
	sendAmount, err := strconv.ParseInt(sendAmountStr, 10, 64)
	if err != nil {
		return MakeJsonErrorResult(RequestError, "Illegal sendAmount:"+err.Error(), nil)
	}
	predictAmount, err := strconv.ParseInt(predictAmountStr, 10, 64)
	if err != nil {
		return MakeJsonErrorResult(RequestError, "Illegal predictAmount:"+err.Error(), nil)
	}
	slippage, err := strconv.ParseInt(slippageStr, 10, 16)
	if err != nil {
		return MakeJsonErrorResult(RequestError, "Illegal slippage:"+err.Error(), nil)
	}
	err = swapExactReceiveInChannel(sendAssetId, sendAmount, receiveAssetId, predictAmount, uint16(slippage), token)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", nil)
}

type swapInfoRequest struct {
	sendInfo    *predictInfo
	receiveInfo *predictInfo
	routePeer   string
	slippage    uint16
	Success     bool
	RfqId       uint
}
type predictInfo struct {
	assetId       string
	predictAmount int64
	channel       *lnrpc.Channel
}
type estimateRfq struct {
	Coefficient string `json:"coefficient"`
	Scale       uint32 `json:"scale"`
	RfqID       uint   `json:"rfq_id"`
}

var swapMutex sync.Mutex

func swapExactSendInChannel(sendAssetId string, sendAmount int64, receiveAssetId string, predictAmount int64, slippage uint16, token string) error {
	if !swapMutex.TryLock() {
		return fmt.Errorf("当前任务执行中，请勿重复操作")
	}
	defer swapMutex.Unlock()
	req := swapInfoRequest{
		sendInfo: &predictInfo{
			assetId:       sendAssetId,
			predictAmount: sendAmount,
		},
		receiveInfo: &predictInfo{
			assetId:       receiveAssetId,
			predictAmount: predictAmount,
		},
		slippage: slippage,
		Success:  false,
		RfqId:    0,
	}
	pubkey, err := getRoutePeerPubkey(token)
	if err != nil {
		return fmt.Errorf("get route peer pubkey error:%v", err)
	}
	req.routePeer = pubkey
	_ = getChannelInfo(&req)
	err = CheckChannelBalance(&req)
	if err != nil {
		return fmt.Errorf("check channel  error:%v", err)
	}
	defer func() {
		if sendAssetId == "00" {
			updateSwapChannelTrans(token, receiveAssetId, &req)
		} else {
			updateSwapChannelTrans(token, sendAssetId, &req)
		}
	}()
	e, err := swapExactSendInChannelHttp(token, &req)
	if err != nil {
		return fmt.Errorf("get estimate rfq error:%v", err)
	}
	req.RfqId = e.RfqID
	err = paySwapInChannel(&req, e, 0)
	if err != nil {
		return fmt.Errorf("the swap operation failed, error:%v", err)
	}
	req.Success = true
	return nil
}
func swapExactReceiveInChannel(sendAssetId string, sendAmount int64, receiveAssetId string, predictAmount int64, slippage uint16, token string) error {
	if !swapMutex.TryLock() {
		return fmt.Errorf("当前任务执行中，请勿重复操作")
	}
	defer swapMutex.Unlock()
	req := swapInfoRequest{
		sendInfo: &predictInfo{
			assetId:       sendAssetId,
			predictAmount: sendAmount,
		},
		receiveInfo: &predictInfo{
			assetId:       receiveAssetId,
			predictAmount: predictAmount,
		},
		slippage: slippage,
		Success:  false,
		RfqId:    0,
	}
	pubkey, err := getRoutePeerPubkey(token)
	if err != nil {
		return fmt.Errorf("get route peer pubkey error:%v", err)
	}
	req.routePeer = pubkey
	_ = getChannelInfo(&req)
	err = CheckChannelBalance(&req)
	if err != nil {
		return fmt.Errorf("check channel  error:%v", err)
	}

	defer func() {
		if sendAssetId == "00" {
			updateSwapChannelTrans(token, receiveAssetId, &req)
		} else {
			updateSwapChannelTrans(token, sendAssetId, &req)
		}
	}()
	e, err := swapExactReceiveInChannelHttp(token, &req)
	if err != nil {
		return fmt.Errorf("get estimate rfq error:%v", err)
	}
	req.RfqId = e.RfqID
	err = paySwapInChannel(&req, e, 1)
	if err != nil {
		return fmt.Errorf("the swap operation failed, error:%v", err)
	}
	return nil
}

func getChannelInfo(req *swapInfoRequest) error {
	list, err := rpcclient.ListChannels(false, false, req.routePeer)
	if err != nil {
		return err
	}
	chooseChannel := func(assetId string, isSend bool) *lnrpc.Channel {
		var c *lnrpc.Channel
		if assetId == "00" {
			for _, channel := range list.GetChannels() {
				if channel.Private == false {
					if isSend {
						if c == nil || channel.LocalBalance > c.LocalBalance {
							c = channel
						}
					} else {
						if c == nil || channel.RemoteBalance > c.RemoteBalance {
							c = channel
						}
					}
				}
			}
		} else {
			for _, channel := range list.GetChannels() {
				if channel.Private && channel.CustomChannelData != nil {
					var customData rfqmsg.JsonAssetChannel
					if err := json.Unmarshal(channel.CustomChannelData, &customData); err != nil {
						continue
					}
					if isSend {
						if len(customData.LocalAssets) > 0 && customData.LocalAssets[0].AssetID == assetId {
							c = channel
							break
						}
					} else {
						if len(customData.RemoteAssets) > 0 && customData.RemoteAssets[0].AssetID == assetId {
							c = channel
							break
						}
					}
				}
			}
		}
		return c
	}
	req.sendInfo.channel = chooseChannel(req.sendInfo.assetId, true)
	req.receiveInfo.channel = chooseChannel(req.receiveInfo.assetId, false)
	return nil
}

func paySwapInChannel(req *swapInfoRequest, eRfq *estimateRfq, exact uint8) error {
	var invoice string
	if req.receiveInfo.assetId == "00" {
		var amount int64
		if exact == 0 {
			amount = assetToSat(req.sendInfo.predictAmount, *eRfq)
		} else {
			amount = req.receiveInfo.predictAmount
		}
		resp, err := rpcclient.AddInvoice(amount, "swapinvoice", req.receiveInfo.channel.ChanId)
		if err != nil {
			return fmt.Errorf("add invoice error:%v", err)
		}
		invoice = resp.PaymentRequest
	} else {
		var amount uint64
		if exact == 0 {
			amount = uint64(satToAsset(req.sendInfo.predictAmount, *eRfq))
		} else {
			amount = uint64(req.receiveInfo.predictAmount)
		}
		resp, err := rpcclient.AddAssetInvoice(amount, req.receiveInfo.assetId, req.routePeer, "swapinvoice")
		if err != nil {
			return fmt.Errorf("add asset invoice error:%v", err)
		}
		invoice = resp.InvoiceResult.PaymentRequest
		if resp.AcceptedBuyQuote.AskAssetRate.Coefficient != eRfq.Coefficient || resp.AcceptedBuyQuote.AskAssetRate.Scale != eRfq.Scale {
			return fmt.Errorf("fail price check")
		}
	}
	_ = invoice
	var result *lnrpc.Payment
	var err error
	if req.sendInfo.assetId == "00" {
		result, err = rpcclient.SendPaymentV2(invoice, 0, 50000,
			req.sendInfo.channel.ChanId, true)
		if err != nil {
			return fmt.Errorf("PayInvoice error:%v", err)
		}
	} else {
		result, err = rpcclient.PayAssetInvoice(invoice, req.sendInfo.assetId, req.routePeer,
			req.sendInfo.channel.ChanId, 50000, true)
		if err != nil {
			return fmt.Errorf("PayAssetInvoice error:%v", err)
		}
	}
	if result.Status != lnrpc.Payment_SUCCEEDED {
		return fmt.Errorf("payment failed,status:%v,reason:%s", result.Status, result.FailureReason)
	}
	return nil
}

func getRoutePeerPubkey(token string) (string, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	if base.NetWork == "regtest" || base.NetWork == "" {
		serverDomainOrSocket = "http://localhost:8081"
	}
	url := serverDomainOrSocket + "/channelSwap/v1/getRfqPeerPubkey"
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		var result struct {
			Error string `json:"error"`
		}
		_ = json.Unmarshal(body, &result)
		return "", fmt.Errorf("get Error code,status:%d,error:%s", res.StatusCode, result.Error)
	}

	var result struct {
		RfqPeerPubkey string `json:"rfqPeerPubkey"`
	}
	_ = json.Unmarshal(body, &result)
	return result.RfqPeerPubkey, nil
}

func swapExactSendInChannelHttp(token string, req *swapInfoRequest) (*estimateRfq, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	if base.NetWork == "regtest" || base.NetWork == "" {
		serverDomainOrSocket = "http://localhost:8081"
	}
	url := serverDomainOrSocket + "/channelSwap/v1/swapExactSendInChannel"
	method := "POST"
	data := struct {
		SendAsset  string  `json:"send_asset"`
		SendAmount float64 `json:"send_amount"`
		RecvAsset  string  `json:"recv_asset"`
		RecvAmount float64 `json:"recv_amount"`
		Slippage   uint16  `json:"slippage"`
	}{
		SendAsset:  req.sendInfo.assetId,
		SendAmount: float64(req.sendInfo.predictAmount),
		RecvAsset:  req.receiveInfo.assetId,
		RecvAmount: float64(req.receiveInfo.predictAmount),
		Slippage:   req.slippage,
	}
	jsonData, _ := json.Marshal(data)
	payload := bytes.NewReader(jsonData)
	client := &http.Client{}
	r, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		var result struct {
			Error string `json:"error"`
		}
		_ = json.Unmarshal(body, &result)
		return nil, fmt.Errorf("get Error code,status:%d,error:%s", res.StatusCode, result.Error)
	}
	result := estimateRfq{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func swapExactReceiveInChannelHttp(token string, req *swapInfoRequest) (*estimateRfq, error) {
	serverDomainOrSocket := Cfg.BtlServerHost
	if base.NetWork == "regtest" || base.NetWork == "" {
		serverDomainOrSocket = "http://localhost:8081"
	}
	url := serverDomainOrSocket + "/channelSwap/v1/swapExactReceiveInChannel"
	method := "POST"
	data := struct {
		SendAsset  string  `json:"send_asset"`
		SendAmount float64 `json:"send_amount"`
		RecvAsset  string  `json:"recv_asset"`
		RecvAmount float64 `json:"recv_amount"`
		Slippage   uint16  `json:"slippage"`
	}{
		SendAsset:  req.sendInfo.assetId,
		SendAmount: float64(req.sendInfo.predictAmount),
		RecvAsset:  req.receiveInfo.assetId,
		RecvAmount: float64(req.receiveInfo.predictAmount),
		Slippage:   req.slippage,
	}
	jsonData, _ := json.Marshal(data)
	payload := bytes.NewReader(jsonData)
	client := &http.Client{}
	r, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		var result struct {
			Error string `json:"error"`
		}
		_ = json.Unmarshal(body, &result)
		return nil, fmt.Errorf("get Error code,status:%d,error:%s", res.StatusCode, result.Error)
	}
	result := estimateRfq{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func satToAsset(sat int64, rfq estimateRfq) int64 {
	coefficient, _ := strconv.ParseFloat(rfq.Coefficient, 64)
	price := 1e8 / (coefficient / math.Pow(10, float64(rfq.Scale)))

	return int64(float64(sat) / price)
}

func assetToSat(asset int64, rfq estimateRfq) int64 {
	coefficient, _ := strconv.ParseFloat(rfq.Coefficient, 64)
	price := 1e8 / (coefficient / math.Pow(10, float64(rfq.Scale)))
	return int64(float64(asset) * price)
}

func updateSwapChannelTrans(token string, assetId string, SwapReq *swapInfoRequest) {
	serverDomainOrSocket := Cfg.BtlServerHost
	if base.NetWork == "regtest" || base.NetWork == "" {
		serverDomainOrSocket = "http://localhost:8081"
	}
	url := serverDomainOrSocket + "/channelSwap/v1/updateSwapChannelTrans"
	method := "POST"
	data := struct {
		AssetId string `json:"asset_id"`
		RfqId   uint   `json:"rfq_id"`
		Success bool   `json:"success"`
	}{
		AssetId: assetId,
		RfqId:   SwapReq.RfqId,
		Success: SwapReq.Success,
	}
	jsonData, _ := json.Marshal(data)
	payload := bytes.NewReader(jsonData)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
}

func CheckChannelBalance(req *swapInfoRequest) error {
	if req.sendInfo.channel == nil {
		return fmt.Errorf("not found send channel,assetId:%s", req.sendInfo.assetId)
	}
	if req.receiveInfo.channel == nil {
		return fmt.Errorf("not found receive channel,assetId:%s", req.receiveInfo.assetId)
	}
	if req.sendInfo.assetId == "00" {
		if req.sendInfo.channel.LocalBalance-req.sendInfo.predictAmount <
			int64(float64(req.sendInfo.channel.Capacity)*0.1) {
			return fmt.Errorf("not enough balance,assetId:%s,balance:%d",
				req.sendInfo.assetId, req.sendInfo.channel.LocalBalance)
		}
	} else {
		var customData rfqmsg.JsonAssetChannel
		if err := json.Unmarshal(req.sendInfo.channel.CustomChannelData, &customData); err != nil {
			return fmt.Errorf("get assetChannel data error:%v", err)
		}
		if int64(customData.LocalBalance)-req.sendInfo.predictAmount < 1 {
			return fmt.Errorf("not enough balance,assetId:%s,balance:%d", req.sendInfo.assetId, customData.LocalBalance)
		}
	}

	if req.receiveInfo.assetId == "00" {
		if req.receiveInfo.channel.RemoteBalance-req.receiveInfo.predictAmount <
			int64(float64(req.receiveInfo.channel.Capacity)*0.1) {
			return fmt.Errorf("not enough recieve balance,assetId:%s,balance:%d",
				req.receiveInfo.assetId, req.receiveInfo.channel.RemoteBalance)
		}
	} else {
		var customData rfqmsg.JsonAssetChannel
		if err := json.Unmarshal(req.receiveInfo.channel.CustomChannelData, &customData); err != nil {
			return fmt.Errorf("get assetChannel data error:%v", err)
		}
		if int64(customData.RemoteBalance)-req.receiveInfo.predictAmount < int64(float64(customData.Capacity)*0.1) {
			return fmt.Errorf("not enough recieve balance,assetId:%s,balance:%d", req.receiveInfo.assetId, customData.RemoteBalance)
		}
	}
	return nil
}
