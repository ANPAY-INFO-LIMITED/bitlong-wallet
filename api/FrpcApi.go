package api

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/fatedier/frp/cmd/frpc/sub"
	"github.com/pkg/errors"
	"github.com/wallet/base"
)

func FrpcConfig(id, remotePortStr string) {
	remotePort, _ := strconv.Atoi(remotePortStr)
	_ = WriteConfig(filepath.Join(base.QueryConfigByKey("dirpath"), "frpc.ini"), base.QueryConfigByKey("serverAddr"), 7000, id, "tcp", "127.0.0.1", 9090, remotePort)
}

func FrpcRun() {
	fmt.Println("=========== LNURL =========== Before system.EnableCompatibilityMode")
	fmt.Println("=========== LNURL =========== Before sub.Execute")
	sub.Execute()
}

func FrpcConf(id, remotePortStr string) error {
	remotePort, err := strconv.Atoi(remotePortStr)
	if err != nil {
		return errors.Wrap(err, "strconv.Atoi")
	}
	return WriteConf(filepath.Join(base.QueryConfigByKey("dirpath"), "frpc.ini"), base.QueryConfigByKey("serverAddr"), 7000, id, "tcp", "127.0.0.1", 9090, remotePort, base.QueryConfigByKey("frpcToken"))
}
