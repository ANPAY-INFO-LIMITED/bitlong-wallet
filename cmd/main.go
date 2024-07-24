package main

import (
	"github.com/wallet/api"
)

const PATH = "D:\\share\\project\\wallet\\config"
const PATH2 = "/home/en/test"

func main() {
	//api.StartLitd()
	println(api.FixAsset("f0de44ec39a53b88c698296a718e0ea6ff19819d7c9efe4b65afe09d09fe773c:1"))
}

func init() {
	api.SetPath(PATH2, "mainnet")

}
