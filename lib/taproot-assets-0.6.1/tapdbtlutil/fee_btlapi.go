package tapdbtlutil

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/lightningnetwork/lnd/lnwallet/chainfee"
)

var (
	// AddrCharge 和 AddrChargeTr 都必须使用p2wkh格式地址来收手续费，为了兼容官方的
	// AddrCharge 为铸币手续费地址
	// AddrChargeTr 为多签手续费地址
	//这里手续费参数
	AddrCharge   = "bc1q7mnlw0nsxpxzgjw79mcjhekdun8h6hngwvlkcg"
	AddrChargeTr = "bc1q8srgdudydhpjv3892qffc24xmqu2j7cy0mryhl"
	//regtest
	//AddrCharge                     = "bcrt1q8amqylya2ahv8ftheq9z95unzveelrxe8ex0rf"
	//AddrChargeTr                   = "bcrt1qn93wjzmc77tav5h99azx27zg2mv0dl275jxq2a"
	TwoKw                          = float64(1.19)
	ThanOneKw                      = float64(0.172)
	MinFeee                        = int64(1500)
	Percentage                     = 0.2
	Network                        = "mainnet" // "regtest/testnet/ default:mainnet"
	MIntFinalizeChargeAmount int64 = 3000
)

func SetFeeParams(network string) {
	switch network {
	case "mainnet":
		AddrCharge = "bc1q7mnlw0nsxpxzgjw79mcjhekdun8h6hngwvlkcg"
		AddrChargeTr = "bc1q8srgdudydhpjv3892qffc24xmqu2j7cy0mryhl"
		Network = "mainnet"
	case "regtest":
		AddrCharge = "bcrt1q2ptgfr0xkvfrfc9vd6nny5lg5numyf3j2qku6d"
		AddrChargeTr = "bcrt1qqsk3c7chzv2x3pz4njnj7gcplvq2xttzys0g0k"
		Network = "regtest"
	}
}

// 解析taproot地址
func DecodeTaprootAddress(strAddr string, cfg *chaincfg.Params) ([]byte,
	error) {
	taprootAddr, err := btcutil.DecodeAddress(strAddr, cfg)
	if err != nil {
		return nil, err
	}

	byteAddr, err := txscript.PayToAddrScript(taprootAddr)
	if err != nil {
		return nil, err
	}
	return byteAddr, nil
}

// 手续费计算
func FeeWe(outPutLenth int64, feeRate chainfee.SatPerKWeight) int64 {
	//实际的交易个数
	lenthTx := outPutLenth - 1
	//计算
	thanTwo := lenthTx - 2
	var txKw float64
	if thanTwo == 0 {
		txKw = TwoKw
	} else {
		txKw = float64(thanTwo)*ThanOneKw + TwoKw
	}
	fee := float64(feeRate) * txKw
	weFee := int64(float64(fee) * Percentage)
	if weFee < MinFeee {
		weFee = MinFeee
	}
	return weFee
}

func GetNetWorkParams(network string) *chaincfg.Params {
	var networkParams *chaincfg.Params
	if Network == "regtest" {
		networkParams = &chaincfg.RegressionNetParams
	} else if Network == "testnet" {
		networkParams = &chaincfg.TestNet3Params
	} else {
		networkParams = &chaincfg.MainNetParams
	}
	return networkParams
}
