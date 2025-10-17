package api

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/fatedier/frp/cmd/frpc/sub"
	frpLog "github.com/fatedier/frp/pkg/util/log"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/wallet/base"
)

const (
	ginMode = gin.DebugMode
	LnuBind = "0.0.0.0"
	LnuPort = 9090
)

func PcLnurlSetServerRouter(username string, password string) *gin.Engine {
	r := gin.Default()

	lnurl := r.Group("/lnurl", gin.BasicAuth(gin.Accounts{
		username: password,
	}))
	lnurl.POST("/gen_invoice", GenInvoiceHandler)
	lnurl.POST("/set_token", SetTokenHandler)
	lnurl.POST("/get_token", GetTokenHandler)
	return r
}

func PcLnurlSetServer(writer io.Writer, username string, password string) *http.Server {

	gin.SetMode(ginMode)
	gin.DefaultWriter = writer
	r := PcLnurlSetServerRouter(username, password)

	bind, port := LnuBind, LnuPort

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", bind, port),
		Handler: r,
	}
	return srv
}

func PcLnurlGetAvailPort() (int, error) {
	pstr, err := GetServerRequestAvailablePort(base.QueryConfigByKey("LnurlServerHost"))
	if err != nil {
		return 0, errors.Wrap(err, "GetServerRequestAvailablePort")
	}
	return pstr, nil
}

func PcLnurlRunFrpcConf(id, remotePort string) error {
	isListening, err := PostServerRequestIsPortListening(remotePort)
	if err != nil {
		return errors.Wrap(err, "PostServerRequestIsPortListening")
	}
	if isListening {
		return errors.New(fmt.Sprintf("isListening: %v\n", serverPortNotAvailable))
	}
	err = FrpcConf(id, remotePort)
	if err != nil {
		return errors.Wrap(err, "FrpcConf")
	}
	dirPath := base.QueryConfigByKey("dirpath")
	cfgPath := filepath.Join(dirPath, "frpc.ini")
	sub.SetCfgFile(cfgPath)
	c := sub.GetCfgFile()
	frpLog.Infof("cfgFile: %s,	dirPath: %s", c, dirPath)
	return nil
}

func PcLnurlRunFrpc() error {
	return sub.ExecuteNoExit()
}

func PcLnurlStopFrpc() {
	sub.StopSrv()
}

func PcLnurlRequest(id, name, localPort, remotePort string) (string, error) {
	lnu, err := PostServerToRequestLnurl(id, name, localPort, remotePort)
	if err != nil {
		return "", errors.Wrap(err, "PostServerToRequestLnurl")
	}
	return lnu, nil
}

func PcLnurlRequestInvoice(lnu string, invoiceType int, assetID string, amount int, pubkey string, memo string) (string, error) {
	invoice, err := PostServerToRequestInvoice(lnu, invoiceType, assetID, amount, pubkey, memo)
	if err != nil {
		return "", errors.Wrap(err, "PostServerToRequestInvoice")
	}
	return invoice, nil
}
