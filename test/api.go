package main

import (
	"flag"
	"fmt"

	"github.com/wallet/api"
)

func main() {
	err := api.SetPath("/home/shui/wallet01/wallet/config", "regtest")
	if err != nil {
		fmt.Println(err)
		return
	}

	var (
		assetId string
		pubkey  string
		amount  int
	)
	flag.StringVar(&assetId, "assetId", "1abd22d7538f1b1213732ee93d391ce946a98575cfa2ca1b4d28c3ebc3656e3d", "assetId")
	flag.StringVar(&pubkey, "pubkey", "02f50eb89fa03c64ba5a110cb0158e80d01d7407182217bb5c48be14ad76e51d91", "pubkey")
	flag.IntVar(&amount, "amount", 500, "amount")
	flag.Parse()

	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}
}
