package routes

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/pc/models"
	"github.com/wallet/pc/pcapi"
	"github.com/wallet/pc/serve"
	"github.com/wallet/pc/utils"
)

var (
	invalidNetwork = errors.New("invalid network")
	invalidReq     = errors.New("invalid request")
)

func Serve(r *gin.RouterGroup) *gin.RouterGroup {

	r.POST("/GetApiVersion", func(c *gin.Context) {
		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: pcapi.GetApiVersion(),
		})
	})

	r.POST("/Shutdown", func(c *gin.Context) {
		var err error
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if serve.FrpcStarted() {
			pcapi.LnurlStopFrpc()
		}

		if serve.LitdStarted() {
			err = pcapi.LitdStopDaemon()
			if err != nil {
				logrus.Errorln(errors.Wrap(err, "pcapi.LitdStopDaemon"))
				err = pcapi.LndStopDaemon()
				if err != nil {
					logrus.Errorln(errors.Wrap(err, "pcapi.LndStopDaemon"))
				}
			}
		}

		defer func() {
			if err = serve.Srv().Shutdown(ctx); err != nil {
				logrus.Fatalln(errors.Wrap(err, "serve.Srv().Shutdown"))
			}
			logrus.Infoln("Server exited.")

			if err = serve.SrvLnu().Shutdown(ctx); err != nil {
				logrus.Fatalln(errors.Wrap(err, "serve.SrvLnu().Shutdown"))
			}
			logrus.Infoln("LnuServer exited.")

			if err := serve.SrvProxy().Shutdown(ctx); err != nil {
				logrus.Fatalln(errors.Wrap(err, "serve.SrvProxy().Shutdown"))
			}
			logrus.Infoln("ProxyServer exited.")

		}()
	})

	r.POST("/CreateFile", func(c *gin.Context) {

		var req struct {
			Path    string `json:"path"`
			Content string `json:"content"`
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

		homeDir, err := os.UserHomeDir()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "os.UserHomeDir"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.UserHomeDirErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}

		err = utils.CreateFile(filepath.Join(homeDir, req.Path), req.Content)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "utils.CreateFile"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.CreateFileErr,
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
	})

	r.POST("/CreateConfig", func(c *gin.Context) {

		var req struct {
			Network string `json:"network"`
			Config  string `json:"config"`
			LitConf string `json:"lit_conf"`
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

		if req.Network == "" || req.Config == "" || req.LitConf == "" {
			logrus.Errorln(invalidReq)
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.InvalidReq,
				Msg:  invalidReq.Error(),
				Data: models.NullStr,
			})
			return
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "os.UserHomeDir"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.UserHomeDirErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}

		var configPath string
		var configFirstLine string
		var configContent string

		var litConfPath string

		switch req.Network {
		case "mainnet":
			configPath = "bitlong_pc/etc/mainnet/config.txt"
			litConfPath = "bitlong_pc/nodes/mainnet/.lit/lit.conf"
			configFirstLine = fmt.Sprintf("dirpath=%s\r\n", filepath.Join(homeDir, "bitlong_pc/nodes/mainnet"))
		case "regtest":
			configPath = "bitlong_pc/etc/regtest/config.txt"
			litConfPath = "bitlong_pc/nodes/regtest/.lit/lit.conf"
			configFirstLine = fmt.Sprintf("dirpath=%s\r\n", filepath.Join(homeDir, "bitlong_pc/nodes/regtest"))
		default:
			logrus.Errorln(invalidNetwork)
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.InvalidNetwork,
				Msg:  invalidNetwork.Error(),
				Data: models.NullStr,
			})
			return
		}
		configContent = configFirstLine + req.Config
		err = utils.CreateFile(filepath.Join(homeDir, configPath), configContent)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "utils.CreateFile; config.txt"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.CreateFileErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}

		err = utils.CreateFile(filepath.Join(homeDir, litConfPath), req.LitConf)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "utils.CreateFile; lit.conf"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.CreateFileErr,
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
	})

	return r

}
