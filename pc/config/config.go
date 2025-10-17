package config

import (
	"github.com/pkg/errors"
	"github.com/wallet/pc/utils"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
)

const (
	DefaultServeBind = localhost
	DefaultServePort = 9876

	DefaultProxyServeBind = localhost
	DefaultProxyServePort = 9880

	localhost = "127.0.0.1"

	LnuBasicAuthUser string = "bitlonguser"
	LnuBasicAuthPass string = "bitlongpass"

	DefaultDbPath = "bitlong_pc/db/bitlong_pc.db"
)

type Config struct {
	Serve      Serve    `json:"serve" yaml:"serve"`
	ProxyServe Serve    `json:"proxy_serve" yaml:"proxy_serve"`
	LnuServe   LnuServe `json:"lnu_serve" yaml:"lnu_serve"`
	Db         Db       `json:"db" yaml:"db"`
}

type Serve struct {
	Bind string `json:"bind" yaml:"bind"`
	Port int    `json:"port" yaml:"port"`
}

type LnuServe struct {
	BasicAuthUser string `json:"basic_auth_user" yaml:"basic_auth_user"`
	BasicAuthPass string `json:"basic_auth_pass" yaml:"basic_auth_pass"`
}

type Db struct {
	Path string `json:"path" yaml:"path"`
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

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return errors.Wrap(err, "os.UserHomeDir")
	}
	dbAbsolutePath := filepath.Join(homeDir, DefaultDbPath)

	conf, _ := yaml.Marshal(&Config{
		Serve: Serve{
			Bind: DefaultServeBind,
			Port: DefaultServePort,
		},
		ProxyServe: Serve{
			Bind: DefaultProxyServeBind,
			Port: DefaultProxyServePort,
		},
		LnuServe: LnuServe{
			BasicAuthUser: LnuBasicAuthUser,
			BasicAuthPass: LnuBasicAuthPass,
		},
		Db: Db{Path: dbAbsolutePath},
	})

	err = utils.CreateFile(path, string(conf))
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
