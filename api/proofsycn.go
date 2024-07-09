package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/lightninglabs/taproot-assets/proof"
	"github.com/wallet/service/universeCourier"
	"net/url"
)

func deliverProof(universeUrl, assetId, groupKey, scriptkey, outpoint string) (string, error) {
	if assetId == "" || scriptkey == "" || outpoint == "" {
		return "", errors.New("assetId or scriptkey or outpoint is empty")
	}
	loc := universeCourier.NewProofLoc(assetId, groupKey, scriptkey, outpoint)
	fetchProof, err := universeCourier.FetchProof(*loc)
	if err != nil {
		return "", errors.New("fetch proof failed")
	}
	addrs := newAddrsMap()
	if universeUrl != "" {
		addrs[universeUrl] = nil
	}
	var total, complelte = len(addrs), 0
	for key := range addrs {
		addr, err := proof.ParseCourierAddress(key)
		if err != nil {
			fmt.Println(err)
			continue
		}
		NewCourier, err := universeCourier.NewCourier(addr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer func(NewCourier *universeCourier.Courier) {
			err := NewCourier.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(NewCourier)
		err = NewCourier.DeliverProof(context.Background(), fetchProof)
		if err != nil {
			fmt.Println(err)
			continue
		}
		complelte++
	}
	if complelte == 0 {
		return "", errors.New("deliver proof failed")
	}
	return fmt.Sprintf("deliver proof success, total:%d, complete:%d", total, complelte), nil
}

func newAddrsMap() map[string]*url.URL {
	addrs := make(map[string]*url.URL)
	addrs[Cfg.PostServiceUrl] = nil
	if Cfg.Network == "mainnet" {
		addrs["mainnet.universe.lightning.finance:10029"] = nil
	}
	return addrs
}
