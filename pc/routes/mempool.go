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

func Mempool(r *gin.RouterGroup) *gin.RouterGroup {

	r.POST("/GetAddressTransactionsByMempool", func(c *gin.Context) {

		var req struct {
			Address string `json:"address"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[[]*api.TransactionsSimplified]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		resp, err := pcapi.GetAddressTransactionsByMempool(req.Address)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.GetAddressTransactionsByMempool"))
			c.JSON(http.StatusOK, models.RespT[[]*api.TransactionsSimplified]{
				Code: models.GetAddressTransactionsByMempoolErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespT[[]*api.TransactionsSimplified]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/GetTransactionByMempool", func(c *gin.Context) {

		var req struct {
			Txid string `json:"txid"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[*api.TransactionsSimplified]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		resp, err := pcapi.GetTransactionByMempool(req.Txid)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.GetTransactionByMempool"))
			c.JSON(http.StatusOK, models.RespT[*api.TransactionsSimplified]{
				Code: models.GetTransactionByMempoolErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespT[*api.TransactionsSimplified]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	return r

}
