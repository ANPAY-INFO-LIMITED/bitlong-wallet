package api

import (
	"errors"
	"fmt"
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

	defaultbitcoinpath = "data/chain/bitcoin"
)

func SetPath(path string, network string) error {
	err := base.SetFilePath(path)
	if err != nil {
		return errors.New("path not exist")
	}
	if network != "mainnet" && network != "testnet" && network != "regtest" {
		return errors.New("network not exist")
	}
	base.SetNetwork(network)
	err = LoadServiceConfig()
	if err != nil {
		return errors.New("load config error ")
	}
	return nil
}

func LoadServiceConfig() error {
	apiConnect.LoadConnectConfig()
	universeCourier.LoadUniverseCourierConfig()
	return nil
}

func GetPath() string {
	return base.GetFilePath()
}

// CheckDir Check the integrity of the directory
func CheckDir(dir string) error {
	baseDir := dir
	//Check whether the snapshot file location exists
	neutrinoPath := filepath.Join(baseDir, defaultlndpath, defaultbitcoinpath, base.NetWork)
	fmt.Println(neutrinoPath)
	if !fileExists(neutrinoPath) {
		if err := os.MkdirAll(neutrinoPath, 0700); err != nil {
			return err
		}
	}
	return nil
}

// fileExists reports whether the named file or directory exists.
func fileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
