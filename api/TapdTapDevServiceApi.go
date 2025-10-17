package api

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/lightninglabs/taproot-assets/taprpc/tapdevrpc"
	"github.com/wallet/service/apiConnect"
)

func ImportProof(proofFile, genesisPoint string) bool {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return false
	}
	defer clearUp()
	client := tapdevrpc.NewTapDevClient(conn)
	_proofFileByteSlice, _ := hex.DecodeString(proofFile)
	request := &tapdevrpc.ImportProofRequest{
		ProofFile:    _proofFileByteSlice,
		GenesisPoint: genesisPoint,
	}
	response, err := client.ImportProof(context.Background(), request)
	if err != nil {
		fmt.Printf("%s tapdevrpc QueryRfqAcceptedQuotes Error: %v\n", GetTimeNow(), err)
		return false
	}
	fmt.Printf("%s %v\n", GetTimeNow(), response)
	return true
}
