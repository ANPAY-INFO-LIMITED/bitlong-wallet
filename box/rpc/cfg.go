package rpc

var NodeCfg = RpcCfg{
	Ln: Cfg{
		Host:         "127.0.0.1:10009",
		CertPath:     "/root/.lnd/tls.cert",
		MacaroonPath: "/root/.lnd/data/chain/bitcoin/mainnet/admin.macaroon",
	},
	Tap: Cfg{
		Host:         "127.0.0.1:8443",
		CertPath:     "/root/.lit/tls.cert",
		MacaroonPath: "/root/.tapd/data/mainnet/admin.macaroon",
	},
	Lit: Cfg{
		Host:         "127.0.0.1:8443",
		CertPath:     "/root/.lit/tls.cert",
		MacaroonPath: "/root/.lit/mainnet/lit.macaroon",
	},
}

const (
	proofCourierAddr = "universerpc://132.232.109.84:8444"
)
