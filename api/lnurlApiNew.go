package api

import (
	"fmt"
	"path/filepath"

	"github.com/fatedier/frp/cmd/frpc/sub"
	frpLog "github.com/fatedier/frp/pkg/util/log"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/wallet/base"
)

var (
	serverPortNotAvailable = errors.New("server port not available")
)

const (
	LnurlRouterPort = "9090"
)

func LnurlRunRouter() {
	errChan := make(chan error)

	router := setupRouterOnPhone()
	go func() {
		err := router.Run("0.0.0.0:" + LnurlRouterPort)
		if err != nil {
			errChan <- err
		}
	}()

	err := <-errChan
	if err != nil {
		fmt.Printf("router.Run: %v\n", err)
	}
}

func LnurlGetAvailPort() string {
	pstr, err := GetServerRequestAvailablePort(base.QueryConfigByKey("LnurlServerHost"))
	if err != nil {
		return MakeJsonErrorResult2(GetServerRequestAvailablePortErr, err.Error(), 0)
	}
	return MakeJsonErrorResult2(SUCCESS_2, SUCCESS_2.Error(), pstr)
}

func LnurlGetNewUUID() string {
	return MakeJsonErrorResult2(SUCCESS_2, SUCCESS_2.Error(), uuid.New().String())
}

func LnurlRunFrpcConf(id, remotePort string) string {
	isListening, err := PostServerRequestIsPortListening(remotePort)
	if err != nil {
		return MakeJsonErrorResult2(PostServerRequestIsPortListeningErr, err.Error(), "")
	}
	if isListening {
		return MakeJsonErrorResult2(ServerRequestPortIsListening, fmt.Sprintf("isListening: %v\n", serverPortNotAvailable), "")
	}
	err = FrpcConf(id, remotePort)
	if err != nil {
		return MakeJsonErrorResult2(FrpcConfErr, err.Error(), "")
	}
	dirPath := base.QueryConfigByKey("dirpath")
	cfgPath := filepath.Join(dirPath, "frpc.ini")
	sub.SetCfgFile(cfgPath)
	c := sub.GetCfgFile()
	frpLog.Infof("cfgFile: %s,	dirPath: %s", c, dirPath)
	return MakeJsonErrorResult2(SUCCESS_2, SUCCESS_2.Error(), "")
}

func LnurlRunFrpc() {
	fmt.Println("=========== LNURL =========== Before FrpcRun")
	FrpcRun()
}

func LnurlRequest(id string, name string, localPort string, remotePort string) string {
	lnu, err := PostServerToRequestLnurl(id, name, localPort, remotePort)
	if err != nil {
		return MakeJsonErrorResult2(PostServerToRequestLnurlErr, err.Error(), "")
	}
	return MakeJsonErrorResult2(SUCCESS_2, SUCCESS_2.Error(), lnu)
}

func LnurlRequestInvoice(lnu string, invoiceType int, assetID string, amount int, pubkey string, memo string) string {
	invoice, err := PostServerToRequestInvoice(lnu, invoiceType, assetID, amount, pubkey, memo)
	if err != nil {
		return MakeJsonErrorResult2(PostServerToRequestInvoiceErr, err.Error(), "")
	}
	return MakeJsonErrorResult2(SUCCESS_2, SUCCESS_2.Error(), invoice)
}
