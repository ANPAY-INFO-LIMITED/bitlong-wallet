package universeCourier

import (
	"context"
	"fmt"
	"github.com/lightninglabs/taproot-assets/proof"
	unirpc "github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	"github.com/wallet/base"
	"path/filepath"
)

func DeliverProof(universeHost string, proofFile *proof.AnnotatedProof) error {
	addr, err := proof.ParseCourierAddress(universeHost)
	if err != nil {
		return err
	}
	c, err := newCourier(addr)
	if err != nil {
		return err
	}
	defer func(c *courier) {
		err := c.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(c)
	err = c.DeliverProof(context.Background(), proofFile)
	if err != nil {
		return err
	}
	return nil
}
func ReceiveProof(universeHost string, loc *proof.Locator) error {
	addr, err := proof.ParseCourierAddress(universeHost)
	if err != nil {
		return err
	}
	c, err := newCourier(addr)
	if err != nil {
		return err
	}
	defer func(c *courier) {
		err := c.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(c)
	p, err := c.ReceiveProof(context.Background(), *loc)
	if err != nil {
		return err
	}
	err = ImportProofs(false, p)
	if err != nil {
		return err
	}
	return nil
}

func QueryAssetProof(universeHost string, asset string) (*unirpc.AssetLeafKeyResponse, error) {
	addr, err := proof.ParseCourierAddress(universeHost)
	if err != nil {
		return nil, err
	}
	c, err := newCourier(addr)
	if err != nil {
		return nil, err
	}
	defer func(c *courier) {
		err := c.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(c)
	keys, err := c.QueryAssetKey(asset)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func LoadUniverseCourierConfig() {
	defaultProofPath = filepath.Join(base.Configure("tapd"), "data", base.NetWork, "proofs")
}
