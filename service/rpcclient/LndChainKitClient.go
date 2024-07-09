package rpcclient

import (
	"context"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/lightningnetwork/lnd/lnrpc/chainrpc"
	"github.com/wallet/service/apiConnect"
)

func getChainKitClient() (chainrpc.ChainKitClient, func(), error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	client := chainrpc.NewChainKitClient(conn)
	return client, clearUp, nil
}

func GetBlock(blockHashStr string) (response *chainrpc.GetBlockResponse, err error) {
	client, clearUp, err := getChainKitClient()
	if err != nil {
		return nil, err
	}
	defer clearUp()

	blockHash, err := chainhash.NewHashFromStr(blockHashStr)
	if err != nil {
		return nil, err
	}
	request := &chainrpc.GetBlockRequest{
		BlockHash: blockHash.CloneBytes(),
	}
	response, err = client.GetBlock(context.Background(), request)
	return response, err
}

func GetBlockHash(height int64) ([]byte, error) {
	client, clearUp, err := getChainKitClient()
	if err != nil {
		return nil, err
	}
	defer clearUp()
	request := &chainrpc.GetBlockHashRequest{
		BlockHeight: height,
	}
	response, err := client.GetBlockHash(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return response.BlockHash, nil

}
