package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wallet/box/models"
	"github.com/wallet/box/ver"
	"net/http"
)

func Serve(r *gin.RouterGroup) *gin.RouterGroup {

	r.POST("/Version", func(c *gin.Context) {
		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: ver.Version(),
		})
	})

	return r

}
