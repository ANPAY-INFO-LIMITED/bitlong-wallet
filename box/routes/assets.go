package routes

import (
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/api"
	"github.com/wallet/box/models"
	"github.com/wallet/box/rpc"
	"github.com/wallet/box/services"
	"github.com/wallet/box/utils"
)

func Assets(r *gin.RouterGroup) *gin.RouterGroup {

	r.POST("/NewAddr", func(c *gin.Context) {

		var req struct {
			AssetId string `json:"asset_id"`
			Amt     int    `json:"amt"`
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

		var t rpc.Tap
		assetID, err := hex.DecodeString(req.AssetId)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "hex.DecodeString"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.HexDecodeStringErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}
		resp, err := t.NewAddr(assetID, uint64(req.Amt))
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "t.NewAddr"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.NewAddrErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp.Encoded,
		})
		return

	})

	r.POST("/SendAssets", func(c *gin.Context) {

		var req struct {
			TapAddrs []string `json:"tap_addrs"`
			FeeRate  uint32   `json:"fee_rate"`
			Password string   `json:"password"`
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

		if err := services.CheckPassword(req.Password); err != nil {
			logrus.Errorln(errors.Wrap(err, "services.CheckPassword"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.CheckPasswordErr,
				Msg:  err.Error(),
				Data: "",
			})
			return
		}

		var t rpc.Tap

		resp, err := t.SendAsset(req.TapAddrs, req.FeeRate)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "t.SendAsset"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.SendAssetErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: utils.TxHashEncodeToString(resp.Transfer.GetAnchorTxHash()),
		})
		return

	})

	r.POST("/DecodeAddr", func(c *gin.Context) {

		var req struct {
			Addr string `json:"addr"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[*models.DecodeAddrResp]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		var t rpc.Tap
		resp, err := t.DecodeAddr(req.Addr)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "t.DecodeAddr"))
			c.JSON(http.StatusOK, models.RespT[*models.DecodeAddrResp]{
				Code: models.DecodeAddrErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespT[*models.DecodeAddrResp]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: rpc.ToDecodeAddrResp(resp),
		})
		return

	})

	r.POST("/AssetTransferIn", func(c *gin.Context) {

		var req struct {
			AssetID string `json:"asset_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[[]*api.AddrEvent]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		resp, err := services.AssetTransferIn(req.AssetID)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "services.AssetTransferIn"))
			c.JSON(http.StatusOK, models.RespT[[]*api.AddrEvent]{
				Code: models.AssetTransferInErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespT[[]*api.AddrEvent]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/AssetTransferOut", func(c *gin.Context) {

		var req struct {
			AssetID string `json:"asset_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[[]*api.AssetTransferSimplified]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		resp, err := services.AssetTransferOut(req.AssetID)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "services.AssetTransferOut"))
			c.JSON(http.StatusOK, models.RespT[[]*api.AssetTransferSimplified]{
				Code: models.AssetTransferOutErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespT[[]*api.AssetTransferSimplified]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/AssetUtxo", func(c *gin.Context) {

		var req struct {
			AssetID string `json:"asset_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[[]*api.ManagedUtxo]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		resp, err := services.AssetUtxo(req.AssetID)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "services.AssetUtxo"))
			c.JSON(http.StatusOK, models.RespT[[]*api.ManagedUtxo]{
				Code: models.AssetUtxoErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespT[[]*api.ManagedUtxo]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/ListBalances", func(c *gin.Context) {

		var t rpc.Tap

		resp, err := t.ListBalances()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "t.ListBalances"))
			c.JSON(http.StatusOK, models.RespT[[]*services.ListAssetBalanceInfo]{
				Code: models.NewAddrErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		processed := services.ProcessListBalancesResponse(resp)
		filtered := services.ExcludeListBalancesResponseCollectible(processed)

		c.JSON(http.StatusOK, models.RespT[[]*services.ListAssetBalanceInfo]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: filtered,
		})
		return

	})

	return r

}
