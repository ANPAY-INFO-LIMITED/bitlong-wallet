package api

import (
	"fmt"

	"github.com/wallet/base"
)

func GetAllConfig() string {
	dirpath := base.QueryConfigByKey("dirpath")
	lndhost := base.QueryConfigByKey("lndhost")
	taproothost := base.QueryConfigByKey("taproothost")
	litdhost := base.QueryConfigByKey("litdhost")
	LnurlServerHost := base.QueryConfigByKey("LnurlServerHost")
	serverAddr := base.QueryConfigByKey("serverAddr")
	universeHost := base.QueryConfigByKey("universeHost")
	BasicAuthUser := base.QueryConfigByKey("BasicAuthUser")
	BasicAuthPass := base.QueryConfigByKey("BasicAuthPass")
	frpcToken := base.QueryConfigByKey("frpcToken")
	return fmt.Sprintf("dirpath: %s\n"+
		"lndhost: %s\n"+
		"taproothost: %s\n"+
		"litdhost: %s\n"+
		"LnurlServerHost: %s\n"+
		"serverAddr: %s\n"+
		"universeHost: %s\n"+
		"BasicAuthUser: %s\n"+
		"BasicAuthPass: %s\n"+
		"frpcToken: %s\n", dirpath, lndhost, taproothost, litdhost, LnurlServerHost, serverAddr, universeHost, BasicAuthUser, BasicAuthPass, frpcToken)
}
