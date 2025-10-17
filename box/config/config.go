package config

import (
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/wallet/box/utils"
	"gopkg.in/yaml.v3"
)

const (
	DefaultServeBind = allHost
	DefaultServePort = 8027

	allHost = "0.0.0.0"

	DefaultProxyServeBind = allHost
	DefaultProxyServePort = 9880

	UserPassLength = 64

	LnuBasicAuthUser string = "bitlonguser"
	LnuBasicAuthPass string = "bitlongpass"
)

type Config struct {
	Serve        Serve     `json:"serve" yaml:"serve"`
	ProxyServe   Serve     `json:"proxy_serve" yaml:"proxy_serve"`
	LnuServe     LnuServe  `json:"lnu_serve" yaml:"lnu_serve"`
	DisableCheck bool      `json:"disable_check" yaml:"disable_check"`
	ServerLit    ServerLit `json:"server_lit" yaml:"server_lit"`
	VirtualUUID  string    `json:"virtual_uuid" yaml:"virtual_uuid"`
}

type ServerLit struct {
	IdentityPubkey string `json:"identity_pubkey" yaml:"identity_pubkey"`
	ServerHost     string `json:"server_host" yaml:"server_host"`
}

type Serve struct {
	Bind          string `json:"bind" yaml:"bind"`
	Port          int    `json:"port" yaml:"port"`
	BasicAuthUser string `json:"basic_auth_user" yaml:"basic_auth_user"`
	BasicAuthPass string `json:"basic_auth_pass" yaml:"basic_auth_pass"`
}

type LnuServe struct {
	BasicAuthUser string `json:"basic_auth_user" yaml:"basic_auth_user"`
	BasicAuthPass string `json:"basic_auth_pass" yaml:"basic_auth_pass"`
}

var (
	config Config
)

func Conf() *Config {
	return &config
}

func LoadConfig(path string) (*Config, error) {
	conf, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "os.ReadFile")
	}
	err = yaml.Unmarshal(conf, &config)
	if err != nil {
		return nil, errors.Wrap(err, "yaml.Unmarshal")
	}
	return &config, nil
}

func CreateConfSample(path string) error {
	randUser := utils.RandStr(UserPassLength)
	time.Sleep(10 * time.Nanosecond)
	randPass := utils.RandStr(UserPassLength)

	time.Sleep(10 * time.Nanosecond)
	proxyRandUser := utils.RandStr(UserPassLength)
	time.Sleep(10 * time.Nanosecond)
	proxyRandPass := utils.RandStr(UserPassLength)

	conf, _ := yaml.Marshal(&Config{
		Serve: Serve{
			Bind:          DefaultServeBind,
			Port:          DefaultServePort,
			BasicAuthUser: randUser,
			BasicAuthPass: randPass,
		},
		ProxyServe: Serve{
			Bind:          DefaultProxyServeBind,
			Port:          DefaultProxyServePort,
			BasicAuthUser: proxyRandUser,
			BasicAuthPass: proxyRandPass,
		},
		LnuServe: LnuServe{
			BasicAuthUser: LnuBasicAuthUser,
			BasicAuthPass: LnuBasicAuthPass,
		},
	})
	err := utils.CreateFile(path, string(conf))
	if err != nil {
		return errors.Wrap(err, "CreateFile")
	}
	return nil
}

var (
	logWriter io.Writer
)

func Writer() io.Writer {
	return logWriter
}

func SetWriter(w io.Writer) {
	logWriter = w
}
