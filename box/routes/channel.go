package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lightninglabs/taproot-assets/taprpc/tapchannelrpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/box/models"
	"github.com/wallet/box/rpc"
	"github.com/wallet/box/services"
	"github.com/wallet/box/st"
)

func Channel(r *gin.RouterGroup) *gin.RouterGroup {

	r.POST("/ChannelBalance", func(c *gin.Context) {

		var l rpc.Ln
		resp, err := l.ChannelBalance()

		if err != nil {
			logrus.Errorln(errors.Wrap(err, "l.ChannelBalance"))
			c.JSON(http.StatusOK, models.RespT[*lnrpc.ChannelBalanceResponse]{
				Code: models.ChannelBalanceErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		c.JSON(http.StatusOK, models.RespT[*lnrpc.ChannelBalanceResponse]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/BoxListPaymentsRecords", func(c *gin.Context) {
		var req struct {
			AssetId           string `json:"asset_id"`
			MaxPayments       uint64 `json:"max_payments"`
			CreationDateStart uint64 `json:"creation_date_start"`
			CreationDateEnd   uint64 `json:"creation_date_end"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[*rpc.BoxListPaymentsRecords]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		var l rpc.Ln
		resp, err := l.BoxListPaymentsRecords(req.AssetId, req.MaxPayments, req.CreationDateStart, req.CreationDateEnd)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "l.BoxListPaymentsRecords"))
			c.JSON(http.StatusOK, models.RespT[*rpc.BoxListPaymentsRecords]{
				Code: models.BoxListPaymentsRecordsErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespT[*rpc.BoxListPaymentsRecords]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: &resp,
		})
		return

	})

	r.POST("/BoxListInvoicesRecords", func(c *gin.Context) {
		var req struct {
			AssetId           string `json:"asset_id"`
			PendingOnly       bool   `json:"pending_only"`
			NumMaxInvoices    uint64 `json:"num_max_invoices"`
			CreationDateStart uint64 `json:"creation_date_start"`
			CreationDateEnd   uint64 `json:"creation_date_end"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[*rpc.BoxListInvoicesRecords]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		var l rpc.Ln
		resp, err := l.BoxListInvoicesRecords(req.AssetId, req.PendingOnly, req.NumMaxInvoices, req.CreationDateStart, req.CreationDateEnd)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "l.BoxListInvoicesRecords"))
			c.JSON(http.StatusOK, models.RespT[*rpc.BoxListInvoicesRecords]{
				Code: models.BoxListInvoicesRecordsErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespT[*rpc.BoxListInvoicesRecords]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: &resp,
		})
		return

	})

	r.POST("/BoxBtcAddInvoices", func(c *gin.Context) {
		var req struct {
			Value int64  `json:"value"`
			Memo  string `json:"memo"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[string]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: "",
			})
			return
		}

		var l rpc.Ln
		resp, err := l.BoxBtcAddInvoices(req.Value, req.Memo)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "l.BoxBtcAddInvoices"))
			c.JSON(http.StatusOK, models.RespT[string]{
				Code: models.BoxBtcAddInvoicesErr,
				Msg:  err.Error(),
				Data: "",
			})
			return
		}
		c.JSON(http.StatusOK, models.RespT[string]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return
	})

	r.POST("/BoxBtcPayInvoice", func(c *gin.Context) {
		var req struct {
			Invoice          string `json:"invoice"`
			Amt              int    `json:"amt"`
			Feelimit         int    `json:"feelimit"`
			OutgoingChanId   string `json:"outgoing_chan_id"`
			AllowSelfPayment bool   `json:"allow_self_payment"`
			Password         string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[*lnrpc.Payment]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		if err := services.CheckPassword(req.Password); err != nil {
			logrus.Errorln(errors.Wrap(err, "services.CheckPassword"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.CheckPasswordErr,
				Msg:  err.Error(),
				Data: "",
			})
			return
		}

		var l rpc.Ln
		resp, err := l.BoxBtcPayInvoice(req.Invoice, req.Amt, req.Feelimit, req.OutgoingChanId, req.AllowSelfPayment)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "l.BoxBtcPayInvoice"))
			c.JSON(http.StatusOK, models.RespT[*lnrpc.Payment]{
				Code: models.BoxBtcPayInvoiceErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		if resp != nil {
			if resp.Status == 2 {
				c.JSON(http.StatusOK, models.RespT[*lnrpc.Payment]{
					Code: models.Success,
					Msg:  models.NullStr,
					Data: resp,
				})
				return
			} else if resp.Status == 3 {
				c.JSON(http.StatusOK, models.RespT[*lnrpc.Payment]{
					Code: models.BoxBtcPayInvoiceErr,
					Msg:  resp.FailureReason.String(),
					Data: resp,
				})
				return
			}
		}
		c.JSON(http.StatusOK, models.RespT[*lnrpc.Payment]{
			Code: models.BoxBtcPayInvoiceErr,
			Msg:  "unknown error",
			Data: nil,
		})
		return
	})

	r.POST("/BoxAddAssetInvoice", func(c *gin.Context) {
		var req struct {
			AssetId     string `json:"asset_id"`
			AssetAmount uint64 `json:"asset_amount"`
			Memo        string `json:"memo"`
			PeerPubkey  string `json:"peer_pubkey"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[*tapchannelrpc.AddInvoiceResponse]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		if st.Token() == "" {
			err := services.UpdateToken()
			if err != nil {
				logrus.Errorln(errors.Wrap(err, "services.UpdateToken"))
				c.JSON(http.StatusOK, models.RespStr{
					Code: models.UpdateTokenErr,
					Msg:  err.Error(),
					Data: "",
				})
				return
			}
		}

		var t rpc.Tap
		LocalInvoice, MappingInvoice, err := t.FwdBoxAddAssetInvoice(req.AssetId, req.AssetAmount, req.Memo, req.PeerPubkey)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "t.BoxAddAssetInvoice"))
			c.JSON(http.StatusOK, models.RespT[*tapchannelrpc.AddInvoiceResponse]{
				Code: models.BoxAddAssetInvoiceErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		type mappingInvoice struct {
			LocalInvoice   *tapchannelrpc.AddInvoiceResponse
			MappingInvoice string
		}
		resp := &mappingInvoice{
			LocalInvoice:   LocalInvoice,
			MappingInvoice: MappingInvoice,
		}
		c.JSON(http.StatusOK, models.RespT[*mappingInvoice]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/BoxAssetChannelSendPayment", func(c *gin.Context) {
		var req struct {
			AssetId          string `json:"asset_id"`
			Pubkey           string `json:"pubkey"`
			PaymentReq       string `json:"payment_req"`
			OutgoingChanId   string `json:"outgoing_chan_id"`
			FeeLimitSat      int    `json:"fee_limit_sat"`
			AllowSelfPayment bool   `json:"allow_self_payment"`
			Password         string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[*lnrpc.Payment]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		if err := services.CheckPassword(req.Password); err != nil {
			logrus.Errorln(errors.Wrap(err, "services.CheckPassword"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.CheckPasswordErr,
				Msg:  err.Error(),
				Data: "",
			})
			return
		}

		if st.Token() == "" {
			err := services.UpdateToken()
			if err != nil {
				logrus.Errorln(errors.Wrap(err, "services.UpdateToken"))
				c.JSON(http.StatusOK, models.RespStr{
					Code: models.UpdateTokenErr,
					Msg:  err.Error(),
					Data: "",
				})
				return
			}
		}

		var t rpc.Tap
		resp, err := t.FwdBoxAssetChannelSendPayment(req.AssetId, req.Pubkey, req.PaymentReq, req.OutgoingChanId, req.FeeLimitSat, req.AllowSelfPayment)
		if err != nil {
			c.JSON(http.StatusOK, models.RespT[*lnrpc.Payment]{
				Code: models.BoxAssetChannelSendPaymentErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		if resp != nil {
			if resp.Status == 2 {
				c.JSON(http.StatusOK, models.RespT[*lnrpc.Payment]{
					Code: models.Success,
					Msg:  models.NullStr,
					Data: resp,
				})
				return
			} else if resp.Status == 3 {
				c.JSON(http.StatusOK, models.RespT[*lnrpc.Payment]{
					Code: models.BoxAssetChannelSendPaymentErr,
					Msg:  resp.FailureReason.String(),
					Data: resp,
				})
				return
			}
		}
		c.JSON(http.StatusOK, models.RespT[*lnrpc.Payment]{
			Code: models.BoxAssetChannelSendPaymentErr,
			Msg:  "unknown error",
			Data: nil,
		})
		return
	})

	r.POST("/BoxListChannels", func(c *gin.Context) {
		var l rpc.Ln
		resp, err := l.BoxListChannels()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "l.BoxListChannels"))
			c.JSON(http.StatusOK, models.RespT[*rpc.BoxListChannelsResp]{
				Code: models.BoxListChannelsErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		c.JSON(http.StatusOK, models.RespT[*rpc.BoxListChannelsResp]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/BoxBtcDecodePayReq", func(c *gin.Context) {
		var req struct {
			Invoice string `json:"invoice"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[*lnrpc.PayReq]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		var l rpc.Ln
		resp, err := l.BoxBtcDecodePayReq(req.Invoice)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "l.BoxBtcDecodePayReq"))
			c.JSON(http.StatusOK, models.RespT[*lnrpc.PayReq]{
				Code: models.BoxBtcDecodePayReqErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		c.JSON(http.StatusOK, models.RespT[*lnrpc.PayReq]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/BoxAssetDecodePayReq", func(c *gin.Context) {
		var req struct {
			AssetId string `json:"asset_id"`
			Invoice string `json:"invoice"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[*tapchannelrpc.AssetPayReqResponse]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		var l rpc.Tap
		resp, err := l.BoxAssetDecodePayReq(req.AssetId, req.Invoice)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "l.BoxAssetDecodePayReq"))
			c.JSON(http.StatusOK, models.RespT[*tapchannelrpc.AssetPayReqResponse]{
				Code: models.BoxAssetDecodePayReqErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		c.JSON(http.StatusOK, models.RespT[*tapchannelrpc.AssetPayReqResponse]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/BoxGetInfo", func(c *gin.Context) {
		var l rpc.Ln
		resp, err := l.BoxGetInfo()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "l.BoxGetInfo"))
			c.JSON(http.StatusOK, models.RespT[*lnrpc.GetInfoResponse]{
				Code: models.BoxGetInfoErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		c.JSON(http.StatusOK, models.RespT[*lnrpc.GetInfoResponse]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	return r
}
