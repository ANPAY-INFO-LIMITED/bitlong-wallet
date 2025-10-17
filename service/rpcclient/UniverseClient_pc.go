package rpcclient

import (
	"context"
	"github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	"github.com/pkg/errors"
)

func PcQueryAssetRoots(Id string) (*universerpc.QueryRootResponse, error) {
	client, clearUp, err := getUniverseClient()
	if err != nil {
		return nil, errors.Wrap(err, "getUniverseClient")
	}
	defer clearUp()

	in := &universerpc.AssetRootQuery{}
	in.Id = &universerpc.ID{
		Id: &universerpc.ID_AssetIdStr{
			AssetIdStr: Id,
		},
	}
	roots, err := client.QueryAssetRoots(context.Background(), in)
	if err != nil {
		return nil, errors.Wrap(err, "QueryAssetRoots")
	}
	return roots, nil
}
