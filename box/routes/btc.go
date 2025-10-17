package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/api"
	"github.com/wallet/box/models"
	"github.com/wallet/box/rpc"
	"github.com/wallet/box/services"
)

func Btc(r *gin.RouterGroup) *gin.RouterGroup {

	r.POST("/BtcTransferIn", func(c *gin.Context) {

		resp, err := services.BtcTransferIn()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "services.BtcTransferIn"))
			c.JSON(http.StatusOK, models.RespT[[]*api.BtcTransferInInfoSimplified]{
				Code: models.BtcTransferInErr,
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

	r.POST("/BtcTransferOut", func(c *gin.Context) {

		resp, err := services.BtcTransferOut()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "services.BtcTransferOut"))
			c.JSON(http.StatusOK, models.RespT[[]*api.BtcTransferOutInfoSimplified]{
				Code: models.BtcTransferOutErr,
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

	r.POST("/BtcUtxo", func(c *gin.Context) {

		resp, err := services.BtcUtxo()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "services.BtcUtxo"))
			c.JSON(http.StatusOK, models.RespT[[]*api.ListUnspentUtxo]{
				Code: models.BtcUtxoErr,
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

	r.POST("/GetWalletBalanceResponse", func(c *gin.Context) {

		resp, err := rpc.GetWalletBalanceResponse()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "rpc.GetWalletBalanceResponse"))
			c.JSON(http.StatusOK, models.RespT[*models.WalletBalanceResponse]{
				Code: models.GetWalletBalanceResponseErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespT[*models.WalletBalanceResponse]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("SendCoins", func(c *gin.Context) {

		var req struct {
			Addr     string `json:"addr"`
			Amount   int64  `json:"amount"`
			FeeRate  uint64 `json:"fee_rate"`
			SendAll  bool   `json:"send_all"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: "",
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

		resp, err := l.SendCoins(req.Addr, req.Amount, req.FeeRate, req.SendAll)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "l.SendCoins"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.BtcTransferInErr,
				Msg:  err.Error(),
				Data: "",
			})
			return
		}
		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp.Txid,
		})
		return

	})

	r.POST("/NewAddress", func(c *gin.Context) {

		var l rpc.Ln

		resp, err := l.NewAddress()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "l.NewAddress"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.BtcTransferOutErr,
				Msg:  err.Error(),
				Data: "",
			})
			return
		}

		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp.Address,
		})
		return

	})

	return r

}
