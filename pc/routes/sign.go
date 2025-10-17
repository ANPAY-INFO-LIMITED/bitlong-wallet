package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/pc/models"
	"github.com/wallet/pc/pcapi"
	"net/http"
)

func Sign(r *gin.RouterGroup) *gin.RouterGroup {

	r.POST("/SignSchnorr", func(c *gin.Context) {

		var req struct {
			HexPrivateKey string `json:"hex_private_key"`
			Message       string `json:"message"`
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

		resp, err := pcapi.SignSchnorr(req.HexPrivateKey, req.Message)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.SignSchnorr"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.SignSchnorrErr,
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
