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

func AssetsChain(r *gin.RouterGroup) *gin.RouterGroup {

	r.POST("/ListNormalBalances", func(c *gin.Context) {

		resp, err := pcapi.ListNormalBalances()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.ListNormalBalances"))
			c.JSON(http.StatusOK, models.RespT[[]*api.ListAssetBalanceInfo2]{
				Code: models.ListNormalBalancesErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		c.JSON(http.StatusOK, models.RespT[[]*api.ListAssetBalanceInfo2]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/CheckAssetIssuanceIsLocal", func(c *gin.Context) {

		var req struct {
			AssetId string `json:"asset_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[*api.IsLocalResult]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		resp, err := pcapi.CheckAssetIssuanceIsLocal(req.AssetId)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.CheckAssetIssuanceIsLocal"))
			c.JSON(http.StatusOK, models.RespT[*api.IsLocalResult]{
				Code: models.CheckAssetIssuanceIsLocalErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespT[*api.IsLocalResult]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/AddrReceives", func(c *gin.Context) {

		var req struct {
			AssetId string `json:"asset_id"`
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

		resp, err := pcapi.AddrReceives(req.AssetId)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.AddrReceives"))
			c.JSON(http.StatusOK, models.RespT[[]*api.AddrEvent]{
				Code: models.AddrReceivesErr,
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

	r.POST("/QueryAssetTransfers", func(c *gin.Context) {

		var req struct {
			AssetId string `json:"asset_id"`
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

		resp, err := pcapi.QueryAssetTransfers(req.AssetId)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.QueryAssetTransfers"))
			c.JSON(http.StatusOK, models.RespT[[]*api.AssetTransferSimplified]{
				Code: models.QueryAssetTransfersErr,
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

	r.POST("/AssetUtxos", func(c *gin.Context) {

		var req struct {
			Token   string `json:"token"`
			AssetId string `json:"asset_id"`
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

		resp, err := pcapi.AssetUtxos(req.Token, req.AssetId)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.AssetUtxos"))
			c.JSON(http.StatusOK, models.RespT[[]*api.ManagedUtxo]{
				Code: models.AssetUtxosErr,
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

	r.POST("/NewAddr", func(c *gin.Context) {

		var req struct {
			AssetId  string `json:"asset_id"`
			Amt      int    `json:"amt"`
			Token    string `json:"token"`
			DeviceId string `json:"device_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[*api.QueriedAddr]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		resp, err := pcapi.NewAddr(req.AssetId, req.Amt, req.Token, req.DeviceId)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.NewAddr"))
			c.JSON(http.StatusOK, models.RespT[*api.QueriedAddr]{
				Code: models.NewAddrErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespT[*api.QueriedAddr]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/QueryAddrs", func(c *gin.Context) {

		var req struct {
			AssetId string `json:"asset_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[[]*api.QueriedAddr]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		resp, err := pcapi.QueryAddrs(req.AssetId)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.QueryAddrs"))
			c.JSON(http.StatusOK, models.RespT[[]*api.QueriedAddr]{
				Code: models.QueryAddrsErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespT[[]*api.QueriedAddr]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/SendAssets", func(c *gin.Context) {

		var req struct {
			JsonAddrs string `json:"json_addrs"`
			FeeRate   int64  `json:"fee_rate"`
			Token     string `json:"token"`
			DeviceId  string `json:"device_id"`
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

		resp, err := pcapi.SendAssets(req.JsonAddrs, req.FeeRate, req.Token, req.DeviceId)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.SendAssets"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.SendAssetsErr,
				Msg:  err.Error(),
				Data: models.NullStr,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespStr{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/ListNftGroups", func(c *gin.Context) {

		resp, err := pcapi.ListNftGroups()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.ListNftGroups"))
			c.JSON(http.StatusOK, models.RespT[[]*api.NftGroup]{
				Code: models.ListNftGroupsErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespT[[]*api.NftGroup]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/ListNonGroupNftAssets", func(c *gin.Context) {

		resp, err := pcapi.ListNonGroupNftAssets()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.ListNonGroupNftAssets"))
			c.JSON(http.StatusOK, models.RespT[[]*api.ListAssetsResponse]{
				Code: models.ListNonGroupNftAssetsErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespT[[]*api.ListAssetsResponse]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/GetSpentNftAssets", func(c *gin.Context) {

		resp, err := pcapi.GetSpentNftAssets()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.GetSpentNftAssets"))
			c.JSON(http.StatusOK, models.RespT[[]*api.ListAssetsSimplifiedResponse]{
				Code: models.GetSpentNftAssetsErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespT[[]*api.ListAssetsSimplifiedResponse]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/MintAsset", func(c *gin.Context) {

		var req struct {
			Name                   string `json:"name"`
			AssetTypeIsCollectible bool   `json:"asset_type_is_collectible"`
			Description            string `json:"description"`
			ImagePath              string `json:"image_path"`
			GroupName              string `json:"group_name"`
			Amount                 int    `json:"amount"`
			DecimalDisplay         int    `json:"decimal_display"`
			NewGroupedAsset        bool   `json:"new_grouped_asset"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[*api.PendingBatch]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		resp, err := pcapi.MintAsset(req.Name, req.AssetTypeIsCollectible, req.Description, req.ImagePath, req.GroupName, req.Amount, req.DecimalDisplay, req.NewGroupedAsset)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.MintAsset"))
			c.JSON(http.StatusOK, models.RespT[*api.PendingBatch]{
				Code: models.MintAssetErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespT[*api.PendingBatch]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/AddGroupAsset", func(c *gin.Context) {

		var req struct {
			Name                   string `json:"name"`
			AssetTypeIsCollectible bool   `json:"asset_type_is_collectible"`
			Description            string `json:"description"`
			ImagePath              string `json:"image_path"`
			GroupName              string `json:"group_name"`
			Amount                 int    `json:"amount"`
			GroupKey               string `json:"group_key"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[*api.PendingBatch]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		resp, err := pcapi.AddGroupAsset(req.Name, req.AssetTypeIsCollectible, req.Description, req.ImagePath, req.GroupName, req.Amount, req.GroupKey)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.AddGroupAsset"))
			c.JSON(http.StatusOK, models.RespT[*api.PendingBatch]{
				Code: models.AddGroupAssetErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespT[*api.PendingBatch]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/CancelBatch", func(c *gin.Context) {

		err := pcapi.CancelBatch()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.CancelBatch"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.CancelBatchErr,
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

	r.POST("/FinalizeBatch", func(c *gin.Context) {

		var req struct {
			FeeRate  int    `json:"fee_rate"`
			Token    string `json:"token"`
			DeviceId string `json:"device_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[*api.PendingBatch]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		resp, err := pcapi.FinalizeBatch(req.FeeRate, req.Token, req.DeviceId)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.FinalizeBatch"))
			c.JSON(http.StatusOK, models.RespT[*api.PendingBatch]{
				Code: models.FinalizeBatchErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespT[*api.PendingBatch]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/GetIssuanceTransactionFee", func(c *gin.Context) {

		var req struct {
			Token   string `json:"token"`
			FeeRate int    `json:"fee_rate"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespInt{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: 0,
			})
			return
		}

		resp, err := pcapi.GetIssuanceTransactionFee(req.Token, req.FeeRate)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.GetIssuanceTransactionFee"))
			c.JSON(http.StatusOK, models.RespInt{
				Code: models.GetIssuanceTransactionFeeErr,
				Msg:  err.Error(),
				Data: 0,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespInt{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/GetAssetInfo", func(c *gin.Context) {

		var req struct {
			AssetId string `json:"asset_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[*api.AssetInfo]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		resp, err := pcapi.GetAssetInfo(req.AssetId)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.GetAssetInfo"))
			c.JSON(http.StatusOK, models.RespT[*api.AssetInfo]{
				Code: models.GetAssetInfoErr,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespT[*api.AssetInfo]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/GetWalletBalanceTotalValue", func(c *gin.Context) {

		var req struct {
			Token string `json:"token"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorln(errors.Wrap(err, "c.ShouldBindJSON"))
			c.JSON(http.StatusOK, models.RespT[float64]{
				Code: models.ShouldBindJSONErr,
				Msg:  err.Error(),
				Data: 0,
			})
			return
		}

		resp, err := pcapi.GetWalletBalanceTotalValue(req.Token)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.GetWalletBalanceTotalValue"))
			c.JSON(http.StatusOK, models.RespT[float64]{
				Code: models.GetWalletBalanceTotalValueErr,
				Msg:  err.Error(),
				Data: 0,
			})
			return
		}
		c.JSON(http.StatusOK, models.RespT[float64]{
			Code: models.Success,
			Msg:  models.NullStr,
			Data: resp,
		})
		return

	})

	r.POST("/SyncUniverse", func(c *gin.Context) {

		var req struct {
			UniverseHost string `json:"universe_host"`
			AssetId      string `json:"asset_id"`
			IsTransfer   bool   `json:"is_transfer"`
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

		_, err := pcapi.SyncUniverse(req.UniverseHost, req.AssetId, req.IsTransfer)
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.SyncUniverse"))
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.SyncUniverseErr,
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
