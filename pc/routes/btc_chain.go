package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/api"
	"github.com/wallet/pc/models"
	"github.com/wallet/pc/pcapi"
	"net/http"
)

func BtcChain(r *gin.RouterGroup) *gin.RouterGroup {

	r.POST("/GetWalletBalance", func(c *gin.Context) {

		resp, err := pcapi.GetWalletBalance()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.GetWalletBalance"))
			c.JSON(http.StatusOK, models.RespT[*api.WalletBalanceResponse]{
				Code: models.GetWalletBalanceErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		c.JSON(http.StatusOK, models.RespT[*api.WalletBalanceResponse]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/GetBtcTransferInInfosJsonResult", func(c *gin.Context) {

		var req struct {
			Token string `json:"token"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[[]*api.BtcTransferInInfoSimplified]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		resp, err := pcapi.GetBtcTransferInInfosJsonResult(req.Token)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.GetBtcTransferInInfosJsonResult"))
			c.JSON(http.StatusOK, models.RespT[[]*api.BtcTransferInInfoSimplified]{
				Code: models.GetBtcTransferInInfosJsonResultErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		c.JSON(http.StatusOK, models.RespT[[]*api.BtcTransferInInfoSimplified]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/GetBtcTransferOutInfosJsonResult", func(c *gin.Context) {

		var req struct {
			Token string `json:"token"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[[]*api.BtcTransferOutInfoSimplified]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		resp, err := pcapi.GetBtcTransferOutInfosJsonResult(req.Token)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.GetBtcTransferOutInfosJsonResult"))
			c.JSON(http.StatusOK, models.RespT[[]*api.BtcTransferOutInfoSimplified]{
				Code: models.GetBtcTransferOutInfosJsonResultErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		c.JSON(http.StatusOK, models.RespT[[]*api.BtcTransferOutInfoSimplified]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/BtcUtxos", func(c *gin.Context) {

		var req struct {
			Token string `json:"token"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[[]*api.ListUnspentUtxo]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		resp, err := pcapi.BtcUtxos(req.Token)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.BtcUtxos"))
			c.JSON(http.StatusOK, models.RespT[[]*api.ListUnspentUtxo]{
				Code: models.BtcUtxosErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		c.JSON(http.StatusOK, models.RespT[[]*api.ListUnspentUtxo]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/GetNewAddress", func(c *gin.Context) {

		resp, err := pcapi.GetNewAddress()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.GetNewAddress"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.GetNewAddressErr,
				Msg:  err.Error(),
				Data: "",
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

	r.POST("/SendCoins", func(c *gin.Context) {

		var req struct {
			Addr    string `json:"addr"`
			Amount  int64  `json:"amount"`
			FeeRate int64  `json:"fee_rate"`
			SendAll bool   `json:"send_all"`
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

		resp, err := pcapi.SendCoins(req.Addr, req.Amount, req.FeeRate, req.SendAll)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.SendCoins"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.SendCoinsErr,
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

	r.POST("/MergeUTXO", func(c *gin.Context) {

		var req struct {
			FeeRate int64 `json:"fee_rate"`
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

		resp, err := pcapi.MergeUTXO(req.FeeRate)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.MergeUTXO"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.MergeUTXOErr,
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
