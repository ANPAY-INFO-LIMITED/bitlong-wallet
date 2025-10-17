package rpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
)

type Node string

const (
	ln  Node = "ln"
	tap Node = "tap"
	lit Node = "lit"
)

var (
	invalidNode = errors.New("invalid node type")
)

type Cfg struct {
	Host         string
	CertPath     string
	MacaroonPath string
}

type RpcCfg struct {
	Ln  Cfg
	Tap Cfg
	Lit Cfg
}

func GetConn(n Node, noMacaroon bool) (*grpc.ClientConn, error) {
	var cfg Cfg
	switch n {
	case ln:
		cfg = NodeCfg.Ln
	case tap:
		cfg = NodeCfg.Tap
	case lit:
		cfg = NodeCfg.Lit
	default:
		return nil, invalidNode
	}

	var conn *grpc.ClientConn
	var err error

	cred, err := newTlsCert(cfg.CertPath)

	if err != nil {
		return nil, errors.Wrap(err, "newTlsCert")
	}

	if noMacaroon {
		conn, err = grpc.NewClient(cfg.Host, grpc.WithTransportCredentials(cred))
		if err != nil {
			return nil, errors.Wrap(err, "grpc.NewClient noMacaroon")
		}
	} else {
		macaroon, err := getMacaroon(cfg.MacaroonPath)
		if err != nil {
			return nil, errors.Wrap(err, "getMacaroon")
		}
		conn, err = grpc.NewClient(cfg.Host, grpc.WithTransportCredentials(cred),
			grpc.WithPerRPCCredentials(NewMacaroonCredential(macaroon)), grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(10*1024*1024),
				grpc.MaxCallSendMsgSize(10*1024*1024),
			))
		if err != nil {
			return nil, errors.Wrap(err, "grpc.NewClient")
		}
	}

	return conn, nil
}

func Close(conn *grpc.ClientConn) {
	err := conn.Close()
	if err != nil {
		logrus.Errorln(errors.Wrap(err, "conn.Close"))
	}
}

func newTlsCert(tlsCertPath string) (credentials.TransportCredentials, error) {
	cert, err := os.ReadFile(tlsCertPath)
	if err != nil {
		return nil, errors.Wrap(err, "os.ReadFile")
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(cert) {
		return nil, errors.Wrap(err, "AppendCertsFromPEM")
	}
	config := &tls.Config{
		MinVersion: tls.VersionTLS12,
		RootCAs:    certPool,
	}
	cred := credentials.NewTLS(config)
	return cred, nil
}

func getMacaroon(macaroonPath string) (string, error) {
	macaroonBytes, err := os.ReadFile(macaroonPath)
	if err != nil {
		return "", errors.Wrap(err, "os.ReadFile")
	}
	macaroon := hex.EncodeToString(macaroonBytes)
	return macaroon, nil
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
