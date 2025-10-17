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

func PsbtTlSwap(r *gin.RouterGroup) *gin.RouterGroup {

	r.POST("/GetListEligibleCoins", func(c *gin.Context) {

		var req struct {
			AssetId string `json:"asset_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[[]*api.EligibleCoin]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		resp, err := pcapi.GetListEligibleCoins(req.AssetId)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.GetListEligibleCoins"))
			c.JSON(http.StatusOK, models.RespT[[]*api.EligibleCoin]{
				Code: models.GetListEligibleCoinsErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespT[[]*api.EligibleCoin]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/CreateSellOrderSign", func(c *gin.Context) {

		var req struct {
			AssetId     string `json:"asset_id"`
			AssetNum    uint64 `json:"asset_num"`
			Price       int64  `json:"price"`
			AnchorPoint string `json:"anchor_point"`
			InternalKey string `json:"internal_key"`
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

		resp, err := pcapi.CreateSellOrderSign(req.AssetId, req.AssetNum, req.Price, req.AnchorPoint, req.InternalKey)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.CreateSellOrderSign"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.CreateSellOrderSignErr,
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

	r.POST("/BuySOrderSign", func(c *gin.Context) {

		var req struct {
			SignedSellOrder string `json:"signed_sell_order"`
			FeeRate         uint64 `json:"fee_rate"`
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

		resp, err := pcapi.BuySOrderSign(req.SignedSellOrder, req.FeeRate)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.BuySOrderSign"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.BuySOrderSignErr,
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

	r.POST("/PublishSOrderTx", func(c *gin.Context) {

		var req struct {
			SignedBoughtSOrder string `json:"signed_bought_s_order"`
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

		resp, err := pcapi.PublishSOrderTx(req.SignedBoughtSOrder)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.PublishSOrderTx"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.PublishSOrderTxErr,
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

	r.POST("/AllowFederationSyncInsertAndExport", func(c *gin.Context) {

		err := pcapi.AllowFederationSyncInsertAndExport()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.AllowFederationSyncInsertAndExport"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.AllowFederationSyncInsertAndExportErr,
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

	r.POST("/InsertProofAndRegisterTransfer", func(c *gin.Context) {

		var req struct {
			AssetId            string `json:"asset_id"`
			SignedBoughtSOrder string `json:"signed_bought_s_order"`
			LastProofStr       string `json:"last_proof_str"`
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

		resp, err := pcapi.InsertProofAndRegisterTransfer(req.AssetId, req.SignedBoughtSOrder, req.LastProofStr)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.InsertProofAndRegisterTransfer"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.InsertProofAndRegisterTransferErr,
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
