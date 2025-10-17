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

func Session(r *gin.RouterGroup) *gin.RouterGroup {

	r.POST("/AddSession", func(c *gin.Context) {

		resp, err := pcapi.AddSession()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.AddSession"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.AddSessionErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}

		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp.PairingSecretMnemonic,
		})
		return

	})

	r.POST("/NewSession", func(c *gin.Context) {

		var req struct {
			SessionType string `json:"session_type"`
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

		resp, err := pcapi.NewSession(req.SessionType)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.NewSession"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.NewSessionErr,
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

	r.POST("/ListSessions", func(c *gin.Context) {

		var req struct {
			Filter uint32 `json:"filter"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[[]*api.Session]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		resp, err := pcapi.ListSessions(api.SessionFilter(req.Filter))
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.ListSessions"))
			c.JSON(http.StatusOK, models.RespT[[]*api.Session]{
				Code: models.ListSessionsErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		c.JSON(http.StatusOK, models.RespT[[]*api.Session]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/RevokeAllSessions", func(c *gin.Context) {

		err := pcapi.RevokeAllSessions()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.RevokeAllSessions"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.RevokeAllSessionsErr,
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

	return r
}
