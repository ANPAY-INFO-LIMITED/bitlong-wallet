package api

import (
	"errors"
	"github.com/wallet/base"
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
	return nil
}

func GetPath() string {
	return base.GetFilePath()
}

// 检查目录的完整性
func CheckDir() error {
	baseDir := GetPath()
	//检查快照文件位置是否存在
	neutrinoPath := filepath.Join(baseDir, defaultlndpath, defaultbitcoinpath, base.NetWork)
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

func FileTestConfig() bool {
	return base.FileConfig(GetPath())
}

func ReadConfigFile() {
	base.ReadConfig(GetPath())
}

func ReadConfigFile1() {
	base.ReadConfig1(GetPath())
}

func ReadConfigFile2() {
	base.ReadConfig2(GetPath())
}

func CreateDir() {
	base.CreateDir(GetPath())
}

func CreateDir2() {
	base.CreateDir2(GetPath())
}

func Visit() {
	base.VisitAll()
}
