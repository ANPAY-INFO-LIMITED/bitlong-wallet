package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/pc/models"
	"github.com/wallet/pc/pcapi"
	"net/http"
)

func Login(r *gin.RouterGroup) *gin.RouterGroup {

	r.POST("/Login", func(c *gin.Context) {

		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
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

		token, err := pcapi.Login(req.Username, req.Password)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.Login"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.LoginErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: token,
		})
		return
	})

	return r

}
