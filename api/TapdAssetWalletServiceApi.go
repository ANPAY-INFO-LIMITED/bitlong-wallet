package api

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/lightninglabs/taproot-assets/taprpc/assetwalletrpc"
	"github.com/wallet/service/apiConnect"
)

func AnchorVirtualPsbts(virtualPsbts []string) bool {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return false
	}
	defer clearUp()
	client := assetwalletrpc.NewAssetWalletClient(conn)
	_virtualPsbts := make([][]byte, 0)
	for _, i := range virtualPsbts {
		str, _ := hex.DecodeString(i)
		_virtualPsbts = append(_virtualPsbts, str)
	}
	request := &assetwalletrpc.AnchorVirtualPsbtsRequest{
		VirtualPsbts: _virtualPsbts,
	}
	response, err := client.AnchorVirtualPsbts(context.Background(), request)
	if err != nil {
		fmt.Printf("%s assetwalletrpc AnchorVirtualPsbts Error: %v\n", GetTimeNow(), err)
		return false
	}
	fmt.Printf("%s %v\n", GetTimeNow(), response)
	return true
}

func FundVirtualPsbt(isPsbtNotRaw bool, psbt ...string) bool {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return false
	}
	defer clearUp()
	client := assetwalletrpc.NewAssetWalletClient(conn)
	request := &assetwalletrpc.FundVirtualPsbtRequest{}
	if isPsbtNotRaw {
		_psbtByteSlice, _ := hex.DecodeString(psbt[0])
		request.Template = &assetwalletrpc.FundVirtualPsbtRequest_Psbt{Psbt: _psbtByteSlice}
	} else {
		request.Template = &assetwalletrpc.FundVirtualPsbtRequest_Raw{
			Raw: &assetwalletrpc.TxTemplate{
				Recipients: nil,
			}}
	}
	response, err := client.FundVirtualPsbt(context.Background(), request)
	if err != nil {
		fmt.Printf("%s assetwalletrpc FundVirtualPsbt Error: %v\n", GetTimeNow(), err)
		return false
	}
	fmt.Printf("%s %v\n", GetTimeNow(), response)
	return true
}

func NextInternalKey(keyFamily int) string {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return ""
	}
	defer clearUp()
	client := assetwalletrpc.NewAssetWalletClient(conn)
	request := &assetwalletrpc.NextInternalKeyRequest{
		KeyFamily: uint32(keyFamily),
	}
	response, err := client.NextInternalKey(context.Background(), request)
	if err != nil {
		fmt.Printf("%s assetwalletrpc NextInternalKey Error: %v\n", GetTimeNow(), err)
		return ""
	}
	return response.String()
}

func NextScriptKey(keyFamily int) string {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return ""
	}
	defer clearUp()
	client := assetwalletrpc.NewAssetWalletClient(conn)
	request := &assetwalletrpc.NextScriptKeyRequest{
		KeyFamily: uint32(keyFamily),
	}
	response, err := client.NextScriptKey(context.Background(), request)
	if err != nil {
		fmt.Printf("%s assetwalletrpc NextScriptKey Error: %v\n", GetTimeNow(), err)
		return ""
	}
	return response.String()
}

func ProveAssetOwnership(assetId, scriptKey string) bool {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return false
	}
	defer clearUp()
	client := assetwalletrpc.NewAssetWalletClient(conn)
	_assetIdByteSlice, _ := hex.DecodeString(assetId)
	_scriptKeyByteSlice, _ := hex.DecodeString(scriptKey)
	request := &assetwalletrpc.ProveAssetOwnershipRequest{
		AssetId:   _assetIdByteSlice,
		ScriptKey: _scriptKeyByteSlice,
		Outpoint:  nil,
	}
	response, err := client.ProveAssetOwnership(context.Background(), request)
	if err != nil {
		fmt.Printf("%s assetwalletrpc ProveAssetOwnership Error: %v\n", GetTimeNow(), err)
		return false
	}
	fmt.Printf("%s %v\n", GetTimeNow(), response)
	return true
}

func RemoveUTXOLease() bool {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return false
	}
	defer clearUp()
	client := assetwalletrpc.NewAssetWalletClient(conn)
	request := &assetwalletrpc.RemoveUTXOLeaseRequest{
		Outpoint: nil,
	}
	response, err := client.RemoveUTXOLease(context.Background(), request)
	if err != nil {
		fmt.Printf("%s assetwalletrpc RemoveUTXOLease Error: %v\n", GetTimeNow(), err)
		return false
	}
	fmt.Printf("%s %v\n", GetTimeNow(), response)
	return true
}

func SignVirtualPsbt(fundedPsbt string) bool {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return false
	}
	defer clearUp()
	client := assetwalletrpc.NewAssetWalletClient(conn)
	_fundedPsbtByteSlice, _ := hex.DecodeString(fundedPsbt)
	request := &assetwalletrpc.SignVirtualPsbtRequest{
		FundedPsbt: _fundedPsbtByteSlice,
	}
	response, err := client.SignVirtualPsbt(context.Background(), request)
	if err != nil {
		fmt.Printf("%s assetwalletrpc SignVirtualPsbt Error: %v\n", GetTimeNow(), err)
		return false
	}
	fmt.Printf("%s %v\n", GetTimeNow(), response)
	return true
}

func VerifyAssetOwnership(proofWithWitness string) bool {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		return false
	}
	defer clearUp()
	client := assetwalletrpc.NewAssetWalletClient(conn)
	_proofWithWitnessByteSlice, _ := hex.DecodeString(proofWithWitness)
	request := &assetwalletrpc.VerifyAssetOwnershipRequest{
		ProofWithWitness: _proofWithWitnessByteSlice,
	}
	response, err := client.VerifyAssetOwnership(context.Background(), request)
	if err != nil {
		fmt.Printf("%s assetwalletrpc VerifyAssetOwnership Error: %v\n", GetTimeNow(), err)
		return false
	}
	fmt.Printf("%s %v\n", GetTimeNow(), response)
	return true
}
