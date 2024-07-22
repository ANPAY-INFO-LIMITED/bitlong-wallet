package api

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/lightninglabs/taproot-assets/proof"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	"github.com/wallet/service/rpcclient"
	"github.com/wallet/service/universeCourier"
	"net/url"
)

var (
	InputErr      = errors.New("assetId or scriptkey or outpoint is empty")
	FetchProofErr = errors.New("fetch proof failed")
)

type Jstr struct {
	Jsonstr string `json:"jsonstr"`
}

func DeliverProof(universeUrl, assetId, groupKey, scriptKey, outpoint string) string {
	response, err := deliverProof(universeUrl, assetId, groupKey, scriptKey, outpoint)
	if err != nil {
		return MakeJsonErrorResult(deliverProofErr, err.Error(), "")
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func DeliverIssuanceProof(assetId string) string {
	err := deliverIssuanceProof(Cfg.UniverseUrl, assetId, "")
	if err != nil {
		return MakeJsonErrorResult(deliverIssuanceProofErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", nil)
}

func ReceiveProof(universeUrl, assetId, groupKey, scriptkey, outpoint string) string {
	err := receiveProof(universeUrl, assetId, groupKey, scriptkey, outpoint)
	if err != nil {
		return MakeJsonErrorResult(receiveProofErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", nil)
}

func ReadProof(assetId, groupKey, scriptkey, outpoint string) string {
	p, err := readProof(assetId, groupKey, scriptkey, outpoint)
	if err != nil {
		return MakeJsonErrorResult(readProofErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", p)
}

func QueryAssetProofs(assetId string) string {
	outPoints, err := queryAssetProofs(Cfg.UniverseUrl, assetId)
	if err != nil {
		return MakeJsonErrorResult(queryAssetProofsErr, err.Error(), nil)
	}
	if len(*outPoints) == 0 {
		return MakeJsonErrorResult(NotFoundData, "", nil)
	}
	result := struct {
		OutPoints *[]string `json:"outpoints"`
	}{
		OutPoints: outPoints,
	}
	return MakeJsonErrorResult(SUCCESS, "", result)
}

func deliverProof(universeUrl, assetId, groupKey, scriptkey, outpoint string) (string, error) {
	if assetId == "" || scriptkey == "" || outpoint == "" {
		return "", InputErr
	}
	loc := universeCourier.NewProofLoc(assetId, groupKey, scriptkey, outpoint)
	fetchProof, err := universeCourier.FetchProof(*loc)
	if err != nil {
		return "", FetchProofErr
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

func receiveProof(universeUrl, assetId, groupKey, scriptkey, outpoint string) error {
	if assetId == "" || scriptkey == "" || outpoint == "" {
		return InputErr
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

func readProof(assetId, groupKey, scriptkey, outpoint string) (*Jstr, error) {
	if assetId == "" || scriptkey == "" || outpoint == "" {
		return nil, InputErr
	}
	loc := universeCourier.NewProofLoc(assetId, groupKey, scriptkey, outpoint)
	if loc == nil {
		locRequestErr := errors.New("asset info is fail, please check the assetId, groupKey, scriptkey, outpoint")
		return nil, locRequestErr
	}
	fetchProof, err := universeCourier.FetchProof(*loc)
	if err != nil {
		return nil, FetchProofErr
	}
	p, err := rpcclient.DecodeProof(fetchProof.Blob, 0, false, false)
	if err != nil {
		return nil, errors.New("read proof failed")
	}

	json, err := taprpc.ProtoJSONMarshalOpts.Marshal(p)
	if err != nil {
		return nil, errors.New("json unmarshal failed")
	}
	str := string(json)
	js := Jstr{Jsonstr: str}
	//fmt.Println(str)
	return &js, nil
}

func queryAssetProofs(universeUrl, assetId string) (*[]string, error) {
	if assetId == "" {
		return nil, InputErr
	}
	keys, err := universeCourier.QueryAssetProof(universeUrl, assetId)
	if err != nil {
		fmt.Println("query asset proof failed, err:", err)
		return nil, FetchProofErr
	}

	var OutPoints []string
	for _, key := range keys.AssetKeys {
		if opStr, ok := key.Outpoint.(*universerpc.AssetKey_OpStr); ok {
			opStr := opStr.OpStr
			OutPoints = append(OutPoints, opStr)
		}
	}
	return &OutPoints, nil
}

func newAddrsMap() map[string]*url.URL {
	addrs := make(map[string]*url.URL)
	addrs[Cfg.UniverseUrl] = nil
	if Cfg.Network == "mainnet" {
		addrs["universerpc://mainnet.universe.lightning.finance:10029"] = nil
	}
	return addrs
}
