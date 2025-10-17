package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/pc/models"
	"github.com/wallet/pc/pcapi"
	"net/http"
)

func SyncUpload(r *gin.RouterGroup) *gin.RouterGroup {

	r.POST("/UploadWalletBalance", func(c *gin.Context) {

		var req struct {
			Token    string `json:"token"`
			DeviceId string `json:"device_id"`
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

		err := pcapi.UploadWalletBalance(req.Token, req.DeviceId)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.UploadWalletBalance"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.UploadWalletBalanceErr,
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

	r.POST("/UploadAssetManagedUtxos", func(c *gin.Context) {

		var req struct {
			Token    string `json:"token"`
			DeviceId string `json:"device_id"`
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

		err := pcapi.UploadAssetManagedUtxos(req.Token, req.DeviceId)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.UploadAssetManagedUtxos"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.UploadAssetManagedUtxosErr,
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

	r.POST("/UploadAssetLocalMintHistory", func(c *gin.Context) {

		var req struct {
			Token    string `json:"token"`
			DeviceId string `json:"device_id"`
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

		err := pcapi.UploadAssetLocalMintHistory(req.Token, req.DeviceId)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.UploadAssetLocalMintHistory"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.UploadAssetLocalMintHistoryErr,
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

	r.POST("/UploadAssetListInfo", func(c *gin.Context) {

		var req struct {
			Token    string `json:"token"`
			DeviceId string `json:"device_id"`
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

		err := pcapi.UploadAssetListInfo(req.Token, req.DeviceId)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.UploadAssetListInfo"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.UploadAssetListInfoErr,
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

	r.POST("/UploadAddrReceives", func(c *gin.Context) {

		var req struct {
			Token    string `json:"token"`
			DeviceId string `json:"device_id"`
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

		err := pcapi.UploadAddrReceives(req.Token, req.DeviceId)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.UploadAddrReceives"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.UploadAddrReceivesErr,
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

	r.POST("/UploadAssetTransfer", func(c *gin.Context) {

		var req struct {
			Token    string `json:"token"`
			DeviceId string `json:"device_id"`
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

		err := pcapi.UploadAssetTransfer(req.Token, req.DeviceId)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.UploadAssetTransfer"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.UploadAssetTransferErr,
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

	r.POST("/UploadAssetBalanceInfo", func(c *gin.Context) {

		var req struct {
			Token    string `json:"token"`
			DeviceId string `json:"device_id"`
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

		err := pcapi.UploadAssetBalanceInfo(req.Token, req.DeviceId)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.UploadAssetBalanceInfo"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.UploadAssetBalanceInfoErr,
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

	r.POST("/UploadAssetBalanceHistories", func(c *gin.Context) {

		var req struct {
			Token string `json:"token"`
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

		err := pcapi.UploadAssetBalanceHistories(req.Token)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.UploadAssetBalanceHistories"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.UploadAssetBalanceHistoriesErr,
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

	r.POST("/UploadBtcListUnspentUtxos", func(c *gin.Context) {

		var req struct {
			Token string `json:"token"`
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
		err := pcapi.UploadBtcListUnspentUtxos(req.Token)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.UploadBtcListUnspentUtxos"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.UploadBtcListUnspentUtxosErr,
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

	r.POST("/AutoMintReserved", func(c *gin.Context) {

		var req struct {
			Token    string `json:"token"`
			DeviceId string `json:"device_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[[]string]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		ops, err := pcapi.AutoMintReserved(req.Token, req.DeviceId)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.AutoMintReserved"))
			c.JSON(http.StatusOK, models.RespT[[]string]{
				Code: models.AutoMintReservedErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespT[[]string]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: ops,
		})
		return

	})

	return r

}
