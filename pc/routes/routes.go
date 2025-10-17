package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

const (
	Prefix = "/api/rest/v1"
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

	apiRV1 := r.Group(Prefix)

	Serve(apiRV1)
	LitdLaunch(apiRV1)
	NpubKey(apiRV1)
	Login(apiRV1)
	SyncUpload(apiRV1)
	BtcChain(apiRV1)
	AssetsChain(apiRV1)
	BtcAddress(apiRV1)
	Mempool(apiRV1)
	Sign(apiRV1)
	PsbtTlSwap(apiRV1)
	Lnurl(apiRV1)
	Session(apiRV1)

	return r
}
