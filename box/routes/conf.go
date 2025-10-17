package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/box/models"
	"github.com/wallet/box/services"
)

func Conf(r *gin.RouterGroup) *gin.RouterGroup {

	r.POST("/UpdateLitConfField", func(c *gin.Context) {

		var req struct {
			NewValue string `json:"new_value"`
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

		err := services.UpdateLitConfField("/root/.lit/lit.conf", "lnd.alias", req.NewValue)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "t.AssetChannelSendPayment"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.DecodeAddrErr,
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

	r.POST("/GetMcAndAlias", func(c *gin.Context) {
		mc, err := services.GetMc()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "services.GetMc"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.ReadFileErr,
				Msg:  "读取机器码文件失败: " + err.Error(),
				Data: models.NullStr,
			})
			return
		}
		alias, _, err := services.GetAliasAndEip()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "services.GetAliasAndEip"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.ReadFileErr,
				Msg:  "读取别名文件失败: " + err.Error(),
				Data: models.NullStr,
			})
			return
		}
		c.JSON(http.StatusOK, models.Resp{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: struct {
				Mc    string `json:"mc"`
				Alias string `json:"alias"`
			}{
				Mc:    mc,
				Alias: alias,
			},
		})
		return

	})

	return r

}
