package api

import (
	"bytes"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/wallet/service/rpcclient"
)

func GetBlockWrap(blockHash string) string {
	response, err := rpcclient.GetBlock(blockHash)
	if err != nil {
		return MakeJsonErrorResult(GetBlockErr, err.Error(), nil)
	}
	msgBlock := &wire.MsgBlock{}
	blockReader := bytes.NewReader(response.RawBlock)
	err = msgBlock.Deserialize(blockReader)
	return MakeJsonErrorResult(SUCCESS, "", msgBlock)
}

func GetBlockInfoByHeight(height int64) string {
	response, err := rpcclient.GetBlockHash(height)
	if err != nil {
		return MakeJsonErrorResult(GetBlockHashErr, err.Error(), nil)
	}
	var blockHash chainhash.Hash
	copy(blockHash[:], response)
	hashstr := blockHash.String()

	return GetBlockWrap(hashstr)
}
