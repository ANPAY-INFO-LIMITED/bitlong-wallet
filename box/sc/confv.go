package sc

import "github.com/wallet/box/config"

type ServerLit struct {
	IdentityPubkey string `json:"identity_pubkey"`
	ServerHost     string `json:"server_host"`
}

func ServerConf() ServerLit {

	var identityPubkey string
	var serverHost string

	if config.Conf().ServerLit.IdentityPubkey != "" {
		identityPubkey = config.Conf().ServerLit.IdentityPubkey
	} else {
		identityPubkey = ThailandI
	}

	if config.Conf().ServerLit.ServerHost != "" {
		serverHost = config.Conf().ServerLit.ServerHost
	} else {
		serverHost = ThailandS
	}

	return ServerLit{
		IdentityPubkey: identityPubkey,
		ServerHost:     serverHost,
	}
}
