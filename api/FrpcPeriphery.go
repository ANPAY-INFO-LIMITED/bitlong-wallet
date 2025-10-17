package api

import (
	"fmt"
	"os"
	"strconv"

	_ "github.com/fatedier/frp/assets/frpc"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/wallet/base"
)

func WriteConfigFrpcRunTest() {
	id := uuid.New().String()
	port := strconv.Itoa(RequestServerGetPortAvailable(base.QueryConfigByKey("LnurlServerHost")))
	FrpcConfig(id, port)
	FrpcRun()
}

func InitPhoneDBTest() {
	err := InitPhoneDB()
	if err != nil {
		fmt.Println("init phone db error,", err)
	}
}

func WriteConfig(filePath string, serverAddr string, serverPort int, proxyName string, proxyType string, localIP string, localPort int, remotePort int) bool {
	content := fmt.Sprintf("serverAddr = \"%s\"\nserverPort = %d\n\n[[proxies]]\nname = \"%s\"\ntype = \"%s\"\nlocalIP = \"%s\"\nlocalPort = %d\nremotePort = %d",
		serverAddr, serverPort, proxyName, proxyType, localIP, localPort, remotePort)
	contentByte := []byte(content)
	err := os.WriteFile(filePath, contentByte, 0644)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func WriteConf(filePath string, serverAddr string, serverPort int, proxyName string, proxyType string, localIP string, localPort int, remotePort int, token string) error {
	content := fmt.Sprintf("serverAddr = \"%s\"\nserverPort = %d\nauth.token = \"%s\"\n\n[[proxies]]\nname = \"%s\"\ntype = \"%s\"\nlocalIP = \"%s\"\nlocalPort = %d\nremotePort = %d",
		serverAddr, serverPort, token, proxyName, proxyType, localIP, localPort, remotePort)
	err := WriteToFile(filePath, content)
	if err != nil {
		return errors.Wrap(err, "WriteToFile")
	}
	return nil
}
