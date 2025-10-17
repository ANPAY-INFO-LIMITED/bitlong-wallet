package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/pc/models"
	"github.com/wallet/pc/pcapi"
	"net/http"
)

func NpubKey(r *gin.RouterGroup) *gin.RouterGroup {

	r.POST("/GenerateKeys", func(c *gin.Context) {

		var req struct {
			Mnemonic string `json:"mnemonic"`
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

		key, err := pcapi.GenerateKeys(req.Mnemonic)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.GenerateKeys"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.GenerateKeysErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: key,
		})
		return
	})

	r.POST("/GetPrivateKey", func(c *gin.Context) {

		key, err := pcapi.GetPrivateKey()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.GetPrivateKey"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.GetPrivateKeyErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: key,
		})
		return

	})

	r.POST("/GetNPublicKey", func(c *gin.Context) {

		key, err := pcapi.GetNPublicKey()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.GetNPublicKey"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.GetNPublicKeyErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: key,
		})
		return

	})

	r.POST("/GetPublicKey", func(c *gin.Context) {

		key, err := pcapi.GetPublicKey()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.GetPublicKey"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.GetPublicKeyErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: key,
		})
		return

	})

	r.POST("/GetNBPublicKey", func(c *gin.Context) {

		key, err := pcapi.GetNBPublicKey()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.GetNBPublicKey"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.GetNBPublicKeyErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: key,
		})
		return

	})

	return r

}
