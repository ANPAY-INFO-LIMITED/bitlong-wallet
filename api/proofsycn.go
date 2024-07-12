package api

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/lightninglabs/taproot-assets/proof"
	"github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	"github.com/wallet/service/universeCourier"
	"net/url"
)

func DeliverProof(universeUrl, assetId, groupKey, scriptKey, outpoint string) string {
	response, err := deliverProof(universeUrl, assetId, groupKey, scriptKey, outpoint)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), "")
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}
func DeliverIssuanceProof(assetId string) string {
	err := deliverIssuanceProof(Cfg.UniverseUrl, assetId, "")
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", nil)
}

func ReceiveProof(universeUrl, assetId, groupKey, scriptkey, outpoint string) string {
	err := receiveProof(universeUrl, assetId, groupKey, scriptkey, outpoint)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", nil)
}

func receiveProof(universeUrl, assetId, groupKey, scriptkey, outpoint string) error {
	if assetId == "" || scriptkey == "" || outpoint == "" {
		return errors.New("assetId or scriptkey or outpoint is empty")
	}
	loc := universeCourier.NewProofLoc(assetId, groupKey, scriptkey, outpoint)
	addrs := newAddrsMap()
	if universeUrl != "" {
		addrs[universeUrl] = nil
	}
	for key := range addrs {
		err := universeCourier.ReceiveProof(key, loc)
		if err == nil {
			return nil
		}
	}
	return errors.New("receive proof failed")
}

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
	var total, complete = len(addrs), 0
	for key := range addrs {
		err := universeCourier.DeliverProof(key, fetchProof)
		if err != nil {
			fmt.Println("deliver proof failed, url:", key, "err:", err)
			continue
		}
		complete++
	}
	if complete == 0 {
		return "", errors.New("deliver proof failed")
	}
	return fmt.Sprintf("deliver proof success, total:%d, complete:%d", total, complete), nil
}

func deliverIssuanceProof(universeUrl, assetId, groupKey string) error {
	if assetId == "" {
		return errors.New("assetId is empty")
	}
	var (
		res *universerpc.AssetLeafResponse
		err error
	)
	if groupKey == "" {
		//TODO
	}
	res, err = assetLeaves(false, assetId, universerpc.ProofType_PROOF_TYPE_ISSUANCE)
	if err != nil {
		return errors.New("get asset leaves failed")
	}
	var blob = proof.Blob{}
	blob = res.Leaves[0].Proof

	file, err := blob.AsFile()
	if err != nil {
		fmt.Println("as file failed")
		return errors.New("as file failed")
	}

	var buf bytes.Buffer
	if err := file.Encode(&buf); err != nil {
		fmt.Println("encode failed")
		return errors.New("encode failed")
	}
	err = universeCourier.DeliverProof(universeUrl, &proof.AnnotatedProof{
		Blob: buf.Bytes(),
	})
	if err != nil {
		return errors.New("deliver proof failed")
	}
	return nil
}

func newAddrsMap() map[string]*url.URL {
	addrs := make(map[string]*url.URL)
	addrs[Cfg.UniverseUrl] = nil
	if Cfg.Network == "mainnet" {
		addrs["mainnet.universe.lightning.finance:10029"] = nil
	}
	return addrs
}
