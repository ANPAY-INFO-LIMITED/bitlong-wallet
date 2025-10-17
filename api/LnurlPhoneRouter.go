package api

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lightninglabs/taproot-assets/taprpc/tapchannelrpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/pkg/errors"
	"github.com/wallet/base"
	"github.com/wallet/service/apiConnect"
	"github.com/wallet/service/rpcclient"
	"gopkg.in/resty.v1"
)

func RouterRunOnPhone() {
	router := setupRouterOnPhone()
	go func() {
		err := router.Run("0.0.0.0:9090")
		if err != nil {
			return
		}
	}()
}

func setupRouterOnPhone() *gin.Engine {
	r := gin.Default()

	username := base.QueryConfigByKey("BasicAuthUser")
	password := base.QueryConfigByKey("BasicAuthPass")
	lnurl := r.Group("/lnurl", gin.BasicAuth(gin.Accounts{
		username: password,
	}))
	lnurl.POST("/gen_invoice", GenInvoiceHandler)
	lnurl.POST("/set_token", SetTokenHandler)
	lnurl.POST("/get_token", GetTokenHandler)
	return r
}

func GenInvoiceHandler(c *gin.Context) {
	id := uuid.New().String()

	var req GenInvoiceReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: ReqShouldBindJSON,
			Msg:  err.Error(),
			Data: "",
		})
		return
	}

	invoice, err := GenInvoice(&req)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: GenInvoiceErr,
			Msg:  err.Error(),
			Data: "",
		})
		return
	}

	err = InitPhoneDB()
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: InitPhoneDBErr,
			Msg:  err.Error(),
			Data: "",
		})
		return
	}

	db, err := bolt.Open(filepath.Join(base.QueryConfigByKey("dirpath"), "phone.db"), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: boltOpenErr,
			Msg:  err.Error(),
			Data: "",
		})
		return
	}
	defer func(db *bolt.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("%s db.Close :%v\n", GetTimeNow(), err)
		}
	}(db)
	s := &PhoneStore{DB: db}

	err = s.CreateOrUpdateInvoice("gen_invoice", &Invoice{
		ID:          id,
		InvoiceType: req.InvoiceType,
		AssetID:     req.AssetID,
		Amount:      req.Amount,
		PubKey:      req.PubKey,
		Memo:        req.Memo,
		Invoice:     invoice,
	})
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: CreateOrUpdateInvoiceErr,
			Msg:  err.Error(),
			Data: "",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code: Success,
		Msg:  "",
		Data: invoice,
	})
	return
}

func SetTokenHandler(c *gin.Context) {
	token := c.PostForm("token")
	setToken(token)
	c.JSON(http.StatusOK, Response{
		Code: Success,
		Msg:  "",
		Data: "",
	})
	return
}

func GetTokenHandler(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: Success,
		Msg:  "",
		Data: getToken(),
	})
	return
}

var (
	invalidInvoiceType = errors.New("invalid invoice type")
)

func GenInvoice(req *GenInvoiceReq) (string, error) {
	switch req.InvoiceType {
	case InvoiceTypeBtcOnChain:
		return genBtcOnChain(req)
	case InvoiceTypeAssetOnChain:
		return genAssetOnChain(req)
	case InvoiceTypeBtcChannel:
		return genBtcChannel(req)
	case InvoiceTypeAssetChannel:
		return genAssetChannel(req)
	case InvoiceTypeAccountBtc:
		return genAccountBtcOnChain(req)
	case InvoiceTypeAccountAssetOnChain:
		return genAccountAssetOnChain(req)
	case InvoiceTypeAccountAssetChannel:
		return genAccountAssetChannel(req)
	default:
		return "", invalidInvoiceType
	}
}

type Response struct {
	Code Code   `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type RespStr struct {
	Code Code   `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

type Code int

const (
	Success           Code = 0
	ReqShouldBindJSON      = iota + 1000
	GenInvoiceErr
	InitPhoneDBErr
	boltOpenErr
	CreateOrUpdateInvoiceErr
)

type InvoiceType uint8

const (
	InvoiceTypeBtcOnChain InvoiceType = iota
	InvoiceTypeAssetOnChain
	InvoiceTypeBtcChannel
	InvoiceTypeAssetChannel
	InvoiceTypeAccountBtc
	InvoiceTypeAccountAssetOnChain
	InvoiceTypeAccountAssetChannel
	InvoiceTypeNostr
)

type GenInvoiceReq struct {
	InvoiceType InvoiceType `json:"invoice_type"`
	AssetID     string      `json:"asset_id"`
	Amount      uint64      `json:"amount"`
	PubKey      string      `json:"pub_key"`
	Memo        string      `json:"memo"`
	RfqPeerKey  string      `json:"rfq_peer_key"`
}

type AccountBtcOnChainResp struct {
	Invoice string `json:"invoice"`
	Error   string `json:"error"`
}

type AccountBtcOnChainErr struct {
	Invoice string `json:"invoice"`
	Error   string `json:"error"`
}

type AccountAssetOnChainResp struct {
	Success bool                    `json:"success"`
	Error   string                  `json:"error"`
	Code    ErrCode                 `json:"code"`
	Data    AccountAssetOnChainAddr `json:"data"`
}
type AccountAssetOnChainAddr struct {
	Addr string `json:"addr"`
}

type AccountAssetChannelResp struct {
	Success bool                    `json:"success"`
	Error   string                  `json:"error"`
	Code    ErrCode                 `json:"code"`
	Data    AccountAssetChannelAddr `json:"data"`
}
type AccountAssetChannelAddr struct {
	Invoice string `json:"invoice"`
}

func genBtcOnChain(req *GenInvoiceReq) (string, error) {
	return getP2trAddress()
}

func getP2trAddress() (string, error) {
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

func genAssetOnChain(req *GenInvoiceReq) (string, error) {
	addr, err := rpcclient.NewAddr(req.AssetID, int(req.Amount))
	if err != nil {
		return "", errors.Wrap(err, "rpcclient.NewAddr")
	}
	return addr.Encoded, nil

}

func genBtcChannel(req *GenInvoiceReq) (string, error) {
	var private bool = true
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return "", errors.Wrap(err, "apiConnect.GetConnection")
	}

	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.Invoice{
		Value:   int64(req.Amount),
		Memo:    req.Memo,
		Private: private,
	}
	response, err := client.AddInvoice(context.Background(), request)
	if err != nil {
		return "", errors.Wrap(err, "client.AddInvoice")
	}
	return response.PaymentRequest, nil
}

func genAssetChannel(req *GenInvoiceReq) (string, error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return "", errors.Wrap(err, "apiConnect.GetConnection")
	}

	defer clearUp()
	tacc := tapchannelrpc.NewTaprootAssetChannelsClient(conn)

	assetID, err := hex.DecodeString(req.AssetID)
	if err != nil {
		return "", errors.Wrap(err, "hex.DecodeString")
	}

	peerPubkey, err := hex.DecodeString(req.PubKey)
	if err != nil {
		return "", errors.Wrap(err, "hex.DecodeString")
	}

	resp, err := tacc.AddInvoice(context.Background(), &tapchannelrpc.AddInvoiceRequest{
		AssetId:     assetID,
		AssetAmount: req.Amount,
		PeerPubkey:  peerPubkey,
		InvoiceRequest: &lnrpc.Invoice{
			Memo: req.Memo,
		},
	})
	if err != nil {
		return "", errors.Wrap(err, "tacc.AddInvoice")
	}
	return resp.InvoiceResult.PaymentRequest, nil
}

func genAccountBtcOnChain(req *GenInvoiceReq) (string, error) {

	targetUrl := fmt.Sprintf("http://%s/custodyAccount/invoice/apply", Cfg.BtlServerHost)
	client := resty.New()

	body := map[string]any{
		"amount": req.Amount,
		"memo":   req.Memo,
	}

	var r AccountBtcOnChainResp
	var e AccountBtcOnChainErr

	resp, err := client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", getToken())).
		SetBody(body).
		SetResult(&r).
		SetError(&e).
		Post(targetUrl)

	if r.Error != "" {
		return "", errors.New(r.Error)
	}

	if err != nil {
		return "", errors.Wrap(err, "client.R.Post")
	}
	if r.Error != "" {
		return "", errors.New(r.Error)
	}
	return resp.Result().(*AccountBtcOnChainResp).Invoice, nil
}

func genAccountAssetOnChain(req *GenInvoiceReq) (string, error) {

	targetUrl := fmt.Sprintf("http://%s/custodyAccount/Asset/apply", Cfg.BtlServerHost)
	client := resty.New()

	body := map[string]any{
		"asset_id": req.AssetID,
		"amount":   req.Amount,
	}

	var r AccountAssetOnChainResp
	var e JsonResult

	resp, err := client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", getToken())).
		SetBody(body).
		SetResult(&r).
		SetError(&e).
		Post(targetUrl)

	fmt.Println(r)
	fmt.Println(e)

	if err != nil {
		return "", errors.Wrap(err, "client.R.Post")
	}

	if e.Error != "" {
		return "", errors.New(e.Error)
	}

	if r.Error != "" {
		return "", errors.New(r.Error)
	}

	return resp.Result().(*AccountAssetOnChainResp).Data.Addr, nil
}

func genAccountAssetChannel(req *GenInvoiceReq) (string, error) {

	targetUrl := fmt.Sprintf("http://%s/custodyAccount/Asset/applyInvoice", Cfg.BtlServerHost)
	client := resty.New()

	body := map[string]any{
		"asset_id":     req.AssetID,
		"amount":       req.Amount,
		"rfq_peer_key": req.RfqPeerKey,
	}

	var r AccountAssetChannelResp
	var e JsonResult

	resp, err := client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", getToken())).
		SetBody(body).
		SetResult(&r).
		SetError(&e).
		Post(targetUrl)

	fmt.Println(r)
	fmt.Println(e)

	if err != nil {
		return "", errors.Wrap(err, "client.R.Post")
	}

	if e.Error != "" {
		return "", errors.New(e.Error)
	}

	if r.Error != "" {
		return "", errors.New(r.Error)
	}

	return resp.Result().(*AccountAssetChannelResp).Data.Invoice, nil
}
