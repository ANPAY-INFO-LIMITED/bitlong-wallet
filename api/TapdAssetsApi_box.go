package api

import (
	"encoding/hex"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/walletrpc"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func BtcTransferIn(listAddress *walletrpc.ListAddressesResponse, txs *lnrpc.TransactionDetails, token string) ([]*BtcTransferInInfoSimplified, error) {
	response, err := GetBtcTransferIn(listAddress, txs, token)
	if err != nil {
		return nil, errors.Wrap(err, "GetBtcTransferInInfos")
	}

	var btcTransferInInfos []*BtcTransferInInfoSimplified
	if response == nil {
		return btcTransferInInfos, nil
	}
	for _, i := range *response {
		btcTransferInInfos = append(btcTransferInInfos, &i)
	}
	return btcTransferInInfos, nil

}

func GetBtcTransferIn(listAddress *walletrpc.ListAddressesResponse, txs *lnrpc.TransactionDetails, token string) (*[]BtcTransferInInfoSimplified, error) {
	var btcTransferInInfos []BtcTransferInInfo

	var addrs []string
	for _, accountWithAddresse := range listAddress.AccountWithAddresses {
		addresses := accountWithAddresse.Addresses
		for _, address := range addresses {
			addrs = append(addrs, address.Address)
		}
	}

	transactions, err := GetThenDecodeAndQueryTransactionsWhoseLabelIsNotTapdAssetMintingIn2(txs, token)
	if err != nil {
		return nil, err
	}
	for _, transaction := range *transactions {
		for _, out := range transaction.Vout {
			voutAddress := out.ScriptPubKey.Address
			for _, address := range addrs {
				if voutAddress == address {
					btcTransferInInfos = append(btcTransferInInfos, BtcTransferInInfo{
						Address: voutAddress,
						Value:   out.Value,
						Time:    transaction.Time,
						Detail:  &transaction,
					})
				}
			}
		}
	}
	transactionsSimplified := BtcTransferOutInfoToBtcTransferInInfoSimplified(&btcTransferInInfos)
	return transactionsSimplified, nil
}

func GetThenDecodeAndQueryTransactionsWhoseLabelIsNotTapdAssetMintingIn2(txs *lnrpc.TransactionDetails, token string) (*[]PostGetRawTransactionResultSat, error) {
	getTransactions, err := GetTransactionsWhoseLabelIsNotTapdAssetMinting2(txs)
	if err != nil {
		return nil, err
	}
	var rawTransactions []string
	for _, transaction := range *getTransactions {
		if transaction.Amount >= 0 {
			rawTransactions = append(rawTransactions, transaction.RawTxHex)
		}
	}
	decodedAndQueryTransactions, err := DecodeAndQueryTransactionsWhoseLabelIsNotTapdAssetMinting(token, rawTransactions)
	if err != nil {
		return nil, err
	}
	if decodedAndQueryTransactions.Error != "" {
		return nil, errors.New(decodedAndQueryTransactions.Error)
	}
	btcResult := ProcessDecodedAndQueryTransactionsData(decodedAndQueryTransactions.Data)
	result := ProcessPostGetRawTransactionResultToUseSat(btcResult)
	return result, nil
}

func GetTransactionsWhoseLabelIsNotTapdAssetMinting2(txs *lnrpc.TransactionDetails) (*[]GetTransactionsResponse, error) {
	response, err := GetTransactionsAndGetCustomResponse2(txs)
	if err != nil {
		return nil, err
	}
	var getTransactionsResponse []GetTransactionsResponse
	for _, transaction := range *response {
		if transaction.Label != "tapd-asset-minting" {
			getTransactionsResponse = append(getTransactionsResponse, transaction)
		}
	}
	return &getTransactionsResponse, nil
}

func GetTransactionsAndGetCustomResponse2(txs *lnrpc.TransactionDetails) (*[]GetTransactionsResponse, error) {
	var getTransactionsResponse []GetTransactionsResponse
	for _, transaction := range txs.Transactions {
		var outputDetails []GetTransactionsOutputDetails
		for _, output := range transaction.OutputDetails {
			outputDetails = append(outputDetails, GetTransactionsOutputDetails{
				OutputType:   output.OutputType.String(),
				Address:      output.Address,
				PkScript:     output.PkScript,
				OutputIndex:  int(output.OutputIndex),
				Amount:       int(output.Amount),
				IsOurAddress: output.IsOurAddress,
			})
		}
		var previousOutpoints []GetTransactionsPreviousOutpoints
		for _, previousOutpoint := range transaction.PreviousOutpoints {
			previousOutpoints = append(previousOutpoints, GetTransactionsPreviousOutpoints{
				Outpoint:    previousOutpoint.Outpoint,
				IsOurOutput: previousOutpoint.IsOurOutput,
			})
		}
		getTransactionsResponse = append(getTransactionsResponse, GetTransactionsResponse{
			TxHash:            transaction.TxHash,
			Amount:            int(transaction.Amount),
			NumConfirmations:  int(transaction.NumConfirmations),
			BlockHash:         transaction.BlockHash,
			BlockHeight:       int(transaction.BlockHeight),
			TimeStamp:         int(transaction.TimeStamp),
			TotalFees:         int(transaction.TotalFees),
			DestAddresses:     transaction.DestAddresses,
			OutputDetails:     outputDetails,
			RawTxHex:          transaction.RawTxHex,
			Label:             transaction.Label,
			PreviousOutpoints: previousOutpoints,
		})
	}
	return &getTransactionsResponse, nil
}

func BtcTransferOut(listAddress *walletrpc.ListAddressesResponse, txs *lnrpc.TransactionDetails, token string) ([]*BtcTransferOutInfoSimplified, error) {
	response, err := GetBtcTransferOut(listAddress, txs, token)
	if err != nil {
		return nil, errors.Wrap(err, "GetBtcTransferOutInfos")
	}

	var btcTransferOutInfos []*BtcTransferOutInfoSimplified
	if response == nil {
		return btcTransferOutInfos, nil
	}
	for _, o := range *response {
		btcTransferOutInfos = append(btcTransferOutInfos, &o)
	}
	return btcTransferOutInfos, nil
}

func GetBtcTransferOut(listAddress *walletrpc.ListAddressesResponse, txs *lnrpc.TransactionDetails, token string) (*[]BtcTransferOutInfoSimplified, error) {
	var btcTransferOutInfos []BtcTransferOutInfo

	var addrs []string
	for _, accountWithAddresse := range listAddress.AccountWithAddresses {
		addresses := accountWithAddresse.Addresses
		for _, address := range addresses {
			addrs = append(addrs, address.Address)
		}
	}

	transactions, err := GetThenDecodeAndQueryTransactionsWhoseLabelIsNotTapdAssetMintingIn2(txs, token)
	if err != nil {
		return nil, err
	}
	for _, transaction := range *transactions {
		for _, vin := range transaction.Vin {
			vinAddress := vin.Prevout.ScriptPubKey.Address
			for _, address := range addrs {
				if vinAddress == address {
					btcTransferOutInfos = append(btcTransferOutInfos, BtcTransferOutInfo{
						Address: vinAddress,
						Value:   vin.Prevout.Value,
						Time:    transaction.Time,
						Detail:  &transaction,
					})
				}
			}
		}
	}
	transactionsSimplified := BtcTransferOutInfoToBtcTransferOutInfoSimplified(&btcTransferOutInfos)
	return transactionsSimplified, nil
}

func BtcUtxo(unspent *walletrpc.ListUnspentResponse, token string) ([]*ListUnspentUtxo, error) {
	resp, err := GetBtcUtxo(unspent, token)
	if err != nil {
		return nil, err
	}
	var btcUtxos []*ListUnspentUtxo
	if resp == nil {
		return btcUtxos, nil
	}
	for _, u := range *resp {
		btcUtxos = append(btcUtxos, &u)
	}
	return btcUtxos, nil
}

func GetBtcUtxo(unspent *walletrpc.ListUnspentResponse, token string) (*[]ListUnspentUtxo, error) {
	response, err := ListUnspentAndProcess2(unspent, token)
	if err != nil {
		return nil, errors.Wrap(err, "ListUnspentAndProcess")
	}
	return response, nil
}

func ListUnspentAndProcess2(unspent *walletrpc.ListUnspentResponse, token string) (*[]ListUnspentUtxo, error) {

	btcUtxos := ListUnspentResponseToListUnspentUtxos(unspent)
	btcUtxos = ListUnspentUtxoFilterByDefaultAddress(btcUtxos)
	btcUtxos, err := GetTimeForListUnspentUtxoByBitcoind(token, btcUtxos)
	if err != nil {
		return nil, err
	}
	return btcUtxos, nil
}

func AssetTransferIn(receive *taprpc.AddrReceivesResponse, assetId string) ([]*AddrEvent, error) {
	resp, err := GetAssetTransferIn(receive, assetId)
	if err != nil {
		return nil, err
	}
	var addrReceives []*AddrEvent
	if resp == nil {
		return addrReceives, nil
	}
	for _, r := range *resp {
		addrReceives = append(addrReceives, &r)
	}
	return addrReceives, nil
}

func GetAssetTransferIn(receive *taprpc.AddrReceivesResponse, assetId string) (*[]AddrEvent, error) {

	var addrEvents []AddrEvent
	for _, event := range receive.Events {
		if assetId != "" && assetId != hex.EncodeToString(event.Addr.AssetId) {
			continue
		}
		e := AddrEvent{}
		e.CreationTimeUnixSeconds = int64(event.CreationTimeUnixSeconds)
		a := QueriedAddr{}
		a.GetData(event.Addr)
		e.Addr = &a
		e.Status = event.Status.String()
		e.Outpoint = event.Outpoint
		e.Txid, _ = outpointToTransactionAndIndex(event.Outpoint)
		e.UtxoAmtSat = int64(event.UtxoAmtSat)
		e.TaprootSibling = hex.EncodeToString(event.TaprootSibling)
		e.ConfirmationHeight = int64(event.ConfirmationHeight)
		e.HasProof = event.HasProof
		addrEvents = append(addrEvents, e)
	}
	if len(addrEvents) == 0 {
		return &addrEvents, nil
	}
	result := SortAddrEvents(&addrEvents)
	return result, nil
}

func AssetTransferOut(transfer *taprpc.ListTransfersResponse, assetId string) ([]*AssetTransferSimplified, error) {
	resp, err := GetAssetTransferOut(transfer, assetId)
	if err != nil {
		return nil, err
	}
	var assetTransfers []*AssetTransferSimplified
	if resp == nil {
		return assetTransfers, nil
	}
	for _, t := range *resp {
		assetTransfers = append(assetTransfers, &t)
	}
	return assetTransfers, nil
}

func GetAssetTransferOut(transfer *taprpc.ListTransfersResponse, assetId string) (*[]AssetTransferSimplified, error) {
	token := ""
	assetTransfers, err := QueryAssetTransferSimplified2(transfer, token, assetId)
	if err != nil {
		return nil, errors.Wrap(err, "QueryAssetTransferSimplified")
	}
	assetTransfers = SortAssetTransferSimplified(assetTransfers)
	return assetTransfers, nil
}

func QueryAssetTransferSimplified2(transfer *taprpc.ListTransfersResponse, token string, assetId string) (*[]AssetTransferSimplified, error) {
	var assetTransferSimplified *[]AssetTransferSimplified
	assetTransfers, err := QueryAssetTransfersAndGetResponse2(transfer, assetId)
	if err != nil {
		return nil, err
	}
	if assetTransfers == nil {
		return nil, nil
	}
	_ = token
	assetTransferSimplified, err = ProcessAssetTransfer(assetTransfers)
	return assetTransferSimplified, nil
}

func QueryAssetTransfersAndGetResponse2(transfer *taprpc.ListTransfersResponse, assetId string) (*[]Transfer, error) {

	var transfers []Transfer
	for _, t := range transfer.Transfers {
		if assetId != "" && assetId != hex.EncodeToString(t.Inputs[0].AssetId) {
			continue
		}
		newTransfer := Transfer{}
		newTransfer.GetData(t)
		transfers = append(transfers, newTransfer)
	}
	if len(transfers) == 0 {
		return nil, nil
	}
	return &transfers, nil
}

func AssetUtxo(utxo *taprpc.ListUtxosResponse, token string, assetId string) ([]*ManagedUtxo, error) {
	resp, err := GetAssetUtxo(utxo, token, assetId)
	if err != nil {
		return nil, err
	}
	var assetUtxos []*ManagedUtxo
	if resp == nil {
		return assetUtxos, nil
	}
	for _, u := range *resp {
		assetUtxos = append(assetUtxos, &u)
	}
	return assetUtxos, nil
}

func GetAssetUtxo(utxo *taprpc.ListUtxosResponse, token string, assetId string) (*[]ManagedUtxo, error) {

	managedUtxos := ListUtxosResponseToManagedUtxos(utxo)
	managedUtxos = ManagedUtxosFilterByAssetId(managedUtxos, assetId)
	managedUtxos, err := GetTimeForManagedUtxoByBitcoind(token, managedUtxos)
	if err != nil {
		logrus.Infoln("GetTimeForManagedUtxoByBitcoind", err)
	}
	managedUtxos = SortAssetUtxos(managedUtxos)
	return managedUtxos, nil
}
