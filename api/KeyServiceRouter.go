package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wallet/service"
)

type SignMessage struct {
	Message string `json:"message"`
}

func RouterForKeyService() {
	router := setupRouterKeyService()
	err := router.Run("0.0.0.0:9091")
	if err != nil {
		return
	}
	return
}

func setupRouterKeyService() *gin.Engine {
	router := gin.Default()
	router.GET("/getPublicKey", func(c *gin.Context) {
		publicKey, add, err := service.GetPublicKey()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"time": GetTimeNow(),
			"pk":   publicKey,
			"npk":  add,
		})
	})
	router.POST("/sign", func(c *gin.Context) {
		var signMess SignMessage
		if err := c.ShouldBindJSON(&signMess); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		sign, err := service.SignMessage(signMess.Message)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"time": GetTimeNow(),
			"sign": sign,
		})
	})
	return router
}
