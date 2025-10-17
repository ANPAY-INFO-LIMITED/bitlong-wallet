package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/box/services"
	"github.com/wallet/pc/models"
	"net/http"
)

func Scripts(r *gin.RouterGroup) *gin.RouterGroup {

	r.POST("/Remote", func(c *gin.Context) {

		go func() {
			err := services.EnableRemote()
			if err != nil {
				logrus.Errorln(errors.Wrap(err, "services.EnableRemote"))
				c.JSON(http.StatusOK, models.RespStr{
					Code: models.EnableRemoteErr,
					Msg:  err.Error(),
					Data: models.NullStr,
				})
				return
			}
		}()
		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: models.NullStr,
		})
		return

	})

	r.POST("/Reboot", func(c *gin.Context) {
		var req struct {
			Min int `json:"min"`
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

		if err := services.Reboot(req.Min); err != nil {
			logrus.Errorln(errors.Wrap(err, "services.Reboot"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.ShouldBindJSONErr,
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

	r.POST("/CheckRemoteStatus", func(c *gin.Context) {

		resp, err := services.CheckFrpStatus()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "services.CheckFrpStatus"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.CheckFrpStatusErr,
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

	r.POST("/StopRemote", func(c *gin.Context) {

		pid, err := services.GetFrpPid()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "services.GetFrpPid"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.GetFrpPidErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}

		err = services.KillNine(pid)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "services.KillNine"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.KillNineErr,
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
