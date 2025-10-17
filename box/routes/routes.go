package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wallet/box/middleware"
	"strings"
	"time"
)

const (
	Prefix = "/api/v1"
)

func SetupGinRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			if strings.HasPrefix(origin, "https://localhost") {
				return true
			}
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	apiV1 := r.Group(Prefix)

	apiV1.Use(middleware.IpWhitelist())

	Serve(apiV1)
	Btc(apiV1)
	Assets(apiV1)
	Channel(apiV1)
	Scripts(apiV1)
	Conf(apiV1)

	return r
}
