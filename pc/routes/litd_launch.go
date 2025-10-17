package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/api"
	"github.com/wallet/pc/models"
	"github.com/wallet/pc/pcapi"
	"github.com/wallet/pc/serve"
	"net/http"
	"os"
	"path/filepath"
)

func LitdLaunch(r *gin.RouterGroup) *gin.RouterGroup {

	r.POST("/SetPath", func(c *gin.Context) {

		var req struct {
			Network string `json:"network"`
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
		switch req.Network {
		case "mainnet":
			logrus.Infoln("Set mainnet.")
		case "regtest":
			logrus.Infoln("Set regtest.")
		default:
			logrus.Errorln(invalidNetwork)
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.InvalidNetwork,
				Msg:  invalidNetwork.Error(),
				Data: models.NullStr,
			})
			return
		}

		var configFolder string
		configFolder = fmt.Sprintf("bitlong_pc/etc/%s", req.Network)

		err = pcapi.SetPath(filepath.Join(homeDir, configFolder), req.Network)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.SetPath"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.SetPathErr,
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

	r.POST("/StartLitd", func(c *gin.Context) {

		go func() {
			pcapi.StartLitd()
		}()

		serve.SetLitdStarted(true)

		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: models.NullStr,
		})
		return
	})

	r.POST("/CreateWallet", func(c *gin.Context) {

		var req struct {
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

		mnemonic, err := pcapi.CreateWallet(req.Password)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.CreateWallet"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.CreateWalletErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: mnemonic,
		})
		return
	})

	r.POST("/RestoreWallet", func(c *gin.Context) {

		var req struct {
			Mnemonic string `json:"mnemonic"`
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

		mnemonic, err := pcapi.RestoreWallet(req.Mnemonic, req.Password)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.RestoreWallet"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.RestoreWalletErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: mnemonic,
		})
		return
	})

	r.POST("/UnlockWallet", func(c *gin.Context) {

		var req struct {
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

		err := pcapi.UnlockWallet(req.Password)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.UnlockWallet"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.UnlockWalletErr,
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

	r.POST("/GetState", func(c *gin.Context) {

		resp, err := pcapi.GetState()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.GetState"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.GetStateErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}

		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp.State.String(),
		})
		return
	})

	r.POST("/SubServersStatus", func(c *gin.Context) {

		resp, err := pcapi.SubServersStatus()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.SubServersStatus"))
			c.JSON(http.StatusOK, models.RespT[models.SubServers]{
				Code: models.SubServersStatusErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		c.JSON(http.StatusOK, models.RespT[models.SubServers]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: models.ToSubServerStatus(resp.GetSubServers()),
		})
		return
	})

	r.POST("/LndGetInfo", func(c *gin.Context) {

		resp, err := pcapi.LndGetInfo()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.LndGetInfo"))
			c.JSON(http.StatusOK, models.RespT[*api.GetInfoResp]{
				Code: models.LndGetInfoErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespT[*api.GetInfoResp]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return
	})

	r.POST("/LitdStopDaemon", func(c *gin.Context) {

		err := pcapi.LitdStopDaemon()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.LitdStopDaemon"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.LitdStopDaemonErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}

		serve.SetLitdStarted(false)

		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: models.NullStr,
		})
		return
	})

	r.POST("/LndStopDaemon", func(c *gin.Context) {

		err := pcapi.LndStopDaemon()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.LndStopDaemon"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.LndStopDaemonErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}

		serve.SetLitdStarted(false)

		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: models.NullStr,
		})
		return
	})

	return r

}
