package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/pc/models"
	"github.com/wallet/pc/pcapi"
	"github.com/wallet/pc/serve"
	"net/http"
)

func Lnurl(r *gin.RouterGroup) *gin.RouterGroup {

	r.POST("/LnurlGetAvailPort", func(c *gin.Context) {

		resp, err := pcapi.LnurlGetAvailPort()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.LnurlGetAvailPort"))
			c.JSON(http.StatusOK, models.RespInt{
				Code: models.LnurlGetAvailPortErr,
				Msg:  err.Error(),
				Data: 0,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespInt{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/LnurlRunFrpcConf", func(c *gin.Context) {

		var req struct {
			Id         string `json:"id"`
			RemotePort string `json:"remote_port"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}

		err := pcapi.LnurlRunFrpcConf(req.Id, req.RemotePort)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.LnurlRunFrpcConf"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.LnurlRunFrpcConfErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: models.NullStr,
		})
		return

	})

	r.POST("/LnurlRunFrpc", func(c *gin.Context) {

		go func() {
			err := pcapi.LnurlRunFrpc()
			if err != nil {
				logrus.Errorln(errors.Wrap(err, "pcapi.LnurlRunFrpc"))
				c.JSON(http.StatusOK, models.RespStr{
					Code: models.LnurlRunFrpcErr,
					Msg:  models.NullStr,
					Data: models.NullStr,
				})
				return
			}
		}()

		serve.SetFrpcStarted(true)

		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: models.NullStr,
		})
		return
	})

	r.POST("/LnurlStopFrpc", func(c *gin.Context) {

		pcapi.LnurlStopFrpc()

		serve.SetFrpcStarted(false)

		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: models.NullStr,
		})
		return
	})

	r.POST("/LnurlRequest", func(c *gin.Context) {

		var req struct {
			Id         string `json:"id"`
			Name       string `json:"name"`
			LocalPort  string `json:"local_port"`
			RemotePort string `json:"remote_port"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}

		resp, err := pcapi.LnurlRequest(req.Id, req.Name, req.LocalPort, req.RemotePort)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.LnurlRequest"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.LnurlRequestErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/LnurlRequestInvoice", func(c *gin.Context) {

		var req struct {
			Lnu         string `json:"lnu"`
			InvoiceType int    `json:"invoice_type"`
			AssetID     string `json:"asset_id"`
			Amount      int    `json:"amount"`
			Pubkey      string `json:"pubkey"`
			Memo        string `json:"memo"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}

		resp, err := pcapi.LnurlRequestInvoice(req.Lnu, req.InvoiceType, req.AssetID, req.Amount, req.Pubkey, req.Memo)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.LnurlRequestInvoice"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.LnurlRequestInvoiceErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	return r
}
