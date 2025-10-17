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

func BtcAddress(r *gin.RouterGroup) *gin.RouterGroup {

	r.POST("/GetNewAddressP2tr", func(c *gin.Context) {

		resp, err := pcapi.GetNewAddressP2tr()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.GetNewAddressP2tr"))
			c.JSON(http.StatusOK, models.RespT[*api.Addr]{
				Code: models.GetNewAddressP2trErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		c.JSON(http.StatusOK, models.RespT[*api.Addr]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/GetNewAddressP2wkh", func(c *gin.Context) {

		resp, err := pcapi.GetNewAddressP2wkh()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.GetNewAddressP2wkh"))
			c.JSON(http.StatusOK, models.RespT[*api.Addr]{
				Code: models.GetNewAddressP2wkhErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		c.JSON(http.StatusOK, models.RespT[*api.Addr]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/GetNewAddressNp2wkh", func(c *gin.Context) {

		resp, err := pcapi.GetNewAddressNp2wkh()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.GetNewAddressNp2wkh"))
			c.JSON(http.StatusOK, models.RespT[*api.Addr]{
				Code: models.GetNewAddressNp2wkhErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		c.JSON(http.StatusOK, models.RespT[*api.Addr]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/GetNewAddressP2trExample", func(c *gin.Context) {

		resp, err := pcapi.GetNewAddressP2trExample()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.GetNewAddressP2trExample"))
			c.JSON(http.StatusOK, models.RespT[*api.Addr]{
				Code: models.GetNewAddressP2trExampleErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		c.JSON(http.StatusOK, models.RespT[*api.Addr]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/GetNewAddressP2wkhExample", func(c *gin.Context) {

		resp, err := pcapi.GetNewAddressP2wkhExample()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.GetNewAddressP2wkhExample"))
			c.JSON(http.StatusOK, models.RespT[*api.Addr]{
				Code: models.GetNewAddressP2wkhExampleErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		c.JSON(http.StatusOK, models.RespT[*api.Addr]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/GetNewAddressNp2wkhExample", func(c *gin.Context) {

		resp, err := pcapi.GetNewAddressNp2wkhExample()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.GetNewAddressNp2wkhExample"))
			c.JSON(http.StatusOK, models.RespT[*api.Addr]{
				Code: models.GetNewAddressNp2wkhExampleErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		c.JSON(http.StatusOK, models.RespT[*api.Addr]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/StoreAddr", func(c *gin.Context) {

		var req struct {
			Name           string `json:"name"`
			Address        string `json:"address"`
			Balance        int    `json:"balance"`
			AddressType    string `json:"address_type"`
			DerivationPath string `json:"derivation_path"`
			IsInternal     bool   `json:"is_internal"`
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

		err := pcapi.StoreAddr(req.Name, req.Address, req.Balance, req.Address, req.DerivationPath, req.IsInternal)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.StoreAddr"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.StoreAddrErr,
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

	r.POST("/RemoveAddr", func(c *gin.Context) {

		var req struct {
			Address string `json:"address"`
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

		err := pcapi.RemoveAddr(req.Address)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.RemoveAddr"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.RemoveAddrErr,
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

	r.POST("/QueryAllAddr", func(c *gin.Context) {

		resp, err := pcapi.QueryAllAddr()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.QueryAllAddr"))
			c.JSON(http.StatusOK, models.RespT[[]*api.Addr]{
				Code: models.QueryAllAddrErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespT[[]*api.Addr]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/UpdateAllAddressesByGnzba", func(c *gin.Context) {

		err := pcapi.UpdateAllAddressesByGnzba()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.UpdateAllAddressesByGnzba"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.UpdateAllAddressesByGnzbaErr,
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
