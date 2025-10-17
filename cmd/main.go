package main

import (
	"github.com/wallet/api"
)

const PATH = "D:\\share\\project\\wallet\\config"
const PATH2 = "/home/en/test3"
const PATH5 = "/home/en/lit/test"

func main() {
	api.StartLitd()
}

func init() {
	api.SetPath(PATH5, "regtest")
}
