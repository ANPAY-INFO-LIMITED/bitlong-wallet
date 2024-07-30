package main

import (
	"github.com/wallet/api"
)

const PATH = "D:\\share\\project\\wallet\\config"
const PATH2 = "/home/en/test"

func main() {
	//api.StartLitd()
	api.StartLitd()
}

func init() {
	api.SetPath(PATH2, "regtest")
}
