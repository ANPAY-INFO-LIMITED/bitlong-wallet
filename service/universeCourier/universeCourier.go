package universeCourier

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/wire"
	"github.com/lightninglabs/taproot-assets/asset"
	"github.com/lightninglabs/taproot-assets/proof"
	"github.com/wallet/base"
	"path/filepath"
)

const (
	UniverseHostMainnet = "universerpc://132.232.109.84:8444"
	UniverseHostTestnet = "universerpc://127.0.0.1:1235"
	UniverseHostRegtest = "universerpc://132.232.109.84:8443"
)

//todo:send a specified proof to universe

// AutoDeliverProof It's a test function to auto deliver proofs to the courier service.
func AutoDeliverProof() {
	addr, err := proof.ParseCourierAddress(UniverseHostMainnet)
	if err != nil {
		fmt.Println(err)
	}
	courier, _ := NewCourier(addr)
	if courier == nil {
		return
	}
	defer func(courier proof.Courier) {
		err := courier.Close()
		if err != nil {
			return
		}
	}(courier)

	id := "1f8c52ffd0c88e6f9584d50206496769acf6aa1ba9e12a0abd661ac4a949c57b"
	b, err := hex.DecodeString(id)
	if err != nil {
		fmt.Println("id erro")
		return
	}

	var assetId asset.ID
	copy(assetId[:], b)
	p, err := FetchProofs(assetId)
	if err != nil {
		fmt.Println("fetch proofs error")
		return
	}
	for index, useproof := range p {
		err := courier.DeliverProof(context.Background(), useproof)

		if err != nil {
			fmt.Println(index)
			return
		}
	}

}

// AutoReceiveProof It's a test function to auto deliver proofs to the courier service.
func AutoReceiveProof(assetId, GroupKey, ScriptKey, op string) {
	// Attempt to receive proof via proof courier service.

	// Parse locator from arguments.
	var (
		_assetId   asset.ID
		_scriptKey btcec.PublicKey
		_op        *wire.OutPoint
		err        error
	)
	if assetId == "" || ScriptKey == "" {
		return
	}
	a, _ := hex.DecodeString(assetId)
	copy(_assetId[:], a)
	b, err := hex.DecodeString(ScriptKey)
	p, err := btcec.ParsePubKey(b)
	if err != nil {
		fmt.Println(err)
		return
	}
	_scriptKey = *p

	_op, err = wire.NewOutPointFromString(op)
	if err != nil {
		fmt.Println(err)
		return
	}
	loc := proof.Locator{
		AssetID:   &_assetId,
		ScriptKey: _scriptKey,
		OutPoint:  _op,
	}

	if GroupKey != "" {
		_groupKey, err := btcec.ParsePubKey([]byte(GroupKey))
		if err != nil {
			fmt.Println(err)
			return
		}
		loc.GroupKey = _groupKey
	}

	// Parse courier address.
	addr, err := proof.ParseCourierAddress(UniverseHostRegtest)
	if err != nil {
		fmt.Println(err)
	}
	// Create a new courier instance.
	courier, _ := NewCourier(addr)
	if courier == nil {
		return
	}
	defer func(courier proof.Courier) {
		err := courier.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(courier)
	// Retrieve proof from courier.
	addrProof, err := courier.ReceiveProof(context.Background(), loc)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Import proofs into the proof directory.
	err = ImportProofs(false, addrProof)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func LoadUniverseCourierConfig() {
	defaultProofPath = filepath.Join(base.Configure("tapd"), "data", base.NetWork, "proofs")
}
