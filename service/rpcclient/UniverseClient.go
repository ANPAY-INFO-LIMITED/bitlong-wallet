package rpcclient

import (
	"context"
	"fmt"
	"github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	"github.com/wallet/service/apiConnect"
)

func getUniverseClient() (universerpc.UniverseClient, func(), error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	client := universerpc.NewUniverseClient(conn)
	return client, clearUp, nil
}

func QueryAssetRoots(Id string) *universerpc.QueryRootResponse {
	client, clearUp, err := getUniverseClient()
	if err != nil {
		fmt.Println(err)
		return nil
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
		return nil
	}
	return roots
}

func AddFederationServer(server string) error {
	client, clearUp, err := getUniverseClient()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer clearUp()

	quest := &universerpc.AddFederationServerRequest{
		Servers: []*universerpc.UniverseFederationServer{
			{
				Host: server,
			},
		},
	}
	_, err = client.AddFederationServer(context.Background(), quest)
	if err != nil {
		return err
	}
	return nil
}

func ListFederationServers() *universerpc.ListFederationServersResponse {
	client, clearUp, err := getUniverseClient()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer clearUp()
	quest := &universerpc.ListFederationServersRequest{}
	servers, err := client.ListFederationServers(context.Background(), quest)
	if err != nil {
		return nil
	}
	return servers
}
