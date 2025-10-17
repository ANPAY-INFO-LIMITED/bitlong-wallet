package apiConnect

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/wallet/base"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type rpccfg struct {
	grpcHost     string
	tlsCertPath  string
	macaroonPath string
}

type ConnCfg struct {
	isInit  bool
	Lndcfg  rpccfg
	Tapdcfg rpccfg
	LitdCfg rpccfg
}

var connCfg = ConnCfg{
	isInit: false,
}

func LoadConnectConfig() {
	connCfg.Lndcfg.grpcHost = base.QueryConfigByKey("lndhost")
	connCfg.Lndcfg.tlsCertPath = filepath.Join(base.Configure("lnd"), "tls.cert")
	connCfg.Lndcfg.macaroonPath = filepath.Join(base.Configure("lnd"), "data", "chain", "bitcoin", base.NetWork, "admin.macaroon")

	connCfg.Tapdcfg.grpcHost = base.QueryConfigByKey("taproothost")
	connCfg.Tapdcfg.tlsCertPath = filepath.Join(base.Configure("lit"), "tls.cert")
	connCfg.Tapdcfg.macaroonPath = filepath.Join(base.Configure("tapd"), "data", base.NetWork, "admin.macaroon")

	connCfg.LitdCfg.grpcHost = base.QueryConfigByKey("litdhost")
	connCfg.LitdCfg.tlsCertPath = filepath.Join(base.Configure("lit"), "tls.cert")
	connCfg.LitdCfg.macaroonPath = filepath.Join(base.Configure("lit"), base.NetWork, "lit.macaroon")
}

func GetConnection(grpcTarget string, isNoMacaroon bool) (*grpc.ClientConn, func(), error) {
	if !connCfg.isInit {
		LoadConnectConfig()
		connCfg.isInit = true
	}
	cfg := rpccfg{}
	switch grpcTarget {
	case "lnd":
		cfg = connCfg.Lndcfg
	case "tapd":
		cfg = connCfg.Tapdcfg
	case "litd":
		cfg = connCfg.LitdCfg
	default:
		return nil, nil, fmt.Errorf("grpcTarget not found")
	}

	var (
		conn *grpc.ClientConn
		err  error
	)
	creds, err := NewTlsCert(cfg.tlsCertPath)
	if err != nil {
		return nil, nil, errors.Wrap(err, "NewTlsCert")
	}
	if isNoMacaroon {
		conn, err = grpc.Dial(cfg.grpcHost, grpc.WithTransportCredentials(creds))
	} else {
		macaroon, err := GetMacaroon(cfg.macaroonPath)
		if err != nil {
			return nil, nil, errors.Wrap(err, "GetMacaroon")
		}
		conn, err = grpc.Dial(cfg.grpcHost, grpc.WithTransportCredentials(creds),
			grpc.WithPerRPCCredentials(NewMacaroonCredential(macaroon)), grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(10*1024*1024), // 10 MB
				grpc.MaxCallSendMsgSize(10*1024*1024), // 10 MB
			))
		if err != nil {
			return nil, nil, err
		}
	}

	cleanUp := func() {
		err := conn.Close()
		if err != nil {
			fmt.Printf("%v,%v\n", err, "conn.Close")
		}
	}
	return conn, cleanUp, err
}

type MacaroonCredential struct {
	macaroon string
}

func NewMacaroonCredential(macaroon string) *MacaroonCredential {
	return &MacaroonCredential{macaroon: macaroon}
}

func (c *MacaroonCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"macaroon": c.macaroon}, nil
}

func (c *MacaroonCredential) RequireTransportSecurity() bool {
	return true
}

func NewTlsCert(tlsCertPath string) (credentials.TransportCredentials, error) {
	cert, err := os.ReadFile(tlsCertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read tls cert: %v", err)
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(cert) {
		return nil, fmt.Errorf("failed to append cert: %v", err)
	}
	config := &tls.Config{
		MinVersion: tls.VersionTLS12,
		RootCAs:    certPool,
	}
	creds := credentials.NewTLS(config)
	return creds, nil
}

func GetMacaroon(macaroonPath string) (string, error) {
	macaroonBytes, err := os.ReadFile(macaroonPath)
	if err != nil {
		return "", errors.Wrap(err, "os.ReadFile")
	}
	macaroon := hex.EncodeToString(macaroonBytes)
	return macaroon, nil
}

func GetTimeNow() string {
	return time.Now().Format("2006/01/02 15:04:05")
}
