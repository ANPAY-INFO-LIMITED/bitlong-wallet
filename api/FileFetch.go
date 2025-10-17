package api

import (
	"fmt"
	"github.com/lightninglabs/taproot-assets/tapdbtlutil"
	"github.com/pkg/errors"
	"github.com/wallet/base"
	"github.com/wallet/service/apiConnect"
	"github.com/wallet/service/universeCourier"
	"os"
	"path/filepath"
)

const (
	defaultlndpath  = ".lnd"
	defaultlitpath  = ".lit"
	defaulttapdpath = ".tapd"
)

type Network string

const (
	Mainnet Network = "mainnet"
	Testnet Network = "testnet"
	Regtest Network = "regtest"
)

func (n Network) String() string {
	return string(n)
}

var (
	invalidNetwork = errors.New("invalid network, must be mainnet, testnet or regtest")
)

type Config struct {
	Network       string `json:"network"`
	UniverseUrl   string `json:"postServiceUrl"`
	BtlServerHost string `json:"btlServerHost"`
}

var Cfg Config

func SetPath(path string, network string) error {
	err := base.SetFilePath(path)
	if err != nil {
		return errors.New("path not exist")
	}
	if network != "mainnet" && network != "testnet" && network != "regtest" {
		return errors.New("network not exist")
	}
	base.SetNetwork(network)

	err = Cfg.loadConfig()
	if err != nil {
		return errors.New("load config error ")
	}
	tapdbtlutil.SetFeeParams(Cfg.Network)
	return nil
}

func BoxSetPath(network Network) error {
	switch network {

	case Mainnet:
		Cfg = Config{
			Network:       Mainnet.String(),
			UniverseUrl:   UniverseHostMainnet,
			BtlServerHost: BtlServerMainnet,
		}

	case Testnet:
		Cfg = Config{
			Network:       Testnet.String(),
			UniverseUrl:   UniverseHostTestnet,
			BtlServerHost: BtlServerTestNet,
		}
	case Regtest:
		Cfg = Config{
			Network:       Regtest.String(),
			UniverseUrl:   UniverseHostRegtest,
			BtlServerHost: BtlServerRegTest,
		}
	default:
		return errors.Wrap(invalidNetwork, network.String())
	}
	return nil

}

func (c *Config) loadConfig() error {
	err := c.loadServiceConfig()
	if err != nil {
		return errors.New("load config error ")
	}
	c.Network = base.NetWork
	switch {
	case base.NetWork == "mainnet":
		c.UniverseUrl = UniverseHostMainnet
		c.BtlServerHost = BtlServerMainnet
	case base.NetWork == "testnet":
		c.UniverseUrl = UniverseHostTestnet
		c.BtlServerHost = BtlServerTestNet
	case base.NetWork == "regtest":
		c.UniverseUrl = UniverseHostRegtest
		c.BtlServerHost = BtlServerRegTest
	default:
		return errors.New("network not exist")
	}
	return nil
}

func (c *Config) loadServiceConfig() error {
	apiConnect.LoadConnectConfig()
	universeCourier.LoadUniverseCourierConfig()
	return nil
}

func GetPath() string {
	return base.GetFilePath()
}

const defaultbitcoinpath = "data/chain/bitcoin"

func CheckDir(dir string) error {
	baseDir := dir
	neutrinoPath := filepath.Join(baseDir, defaultlndpath, defaultbitcoinpath, base.NetWork)
	fmt.Println(neutrinoPath)
	if !fileExists(neutrinoPath) {
		if err := os.MkdirAll(neutrinoPath, 0700); err != nil {
			return err
		}
	}
	return nil
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
