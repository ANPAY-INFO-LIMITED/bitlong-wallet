package api

import (
	"fmt"
	"net"

	"github.com/fiatjaf/go-lnurl"
)

func Encode(url string) string {
	en, _ := lnurl.LNURLEncode(url)
	return en
}

func Decode(lnu string) string {
	de, _ := lnurl.LNURLDecode(lnu)
	return de
}

func QueryAvailablePort() uint16 {
	var startPort uint16 = 1024
	var endPort uint16 = 49151
	for port := startPort; port <= endPort; port++ {
		socket := fmt.Sprintf(":%d", port)
		listener, err := net.Listen("tcp", socket)
		if err == nil {
			_ = listener.Close()
			return port
		}
	}
	return 0
}

func QueryIsPortListening(remotePort string) bool {
	socket := fmt.Sprintf(":%s", remotePort)
	listener, err := net.Listen("tcp", socket)

	if err == nil {
		_ = listener.Close()
		return false
	}
	return false
}
