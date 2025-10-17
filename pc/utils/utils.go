package utils

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	dialTimeout = 2 * time.Second
)

func CreateFile(path, content string) error {
	dir := filepath.Dir(path)
	if dir != "." {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0644)
			if err != nil {
				return errors.Wrap(err, "os.MkdirAll")
			}
		}
	}
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return errors.Wrap(err, "os.WriteFile")
	}
	return nil
}

func PathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, errors.Wrap(err, "os.Stat")
}

func AvailablePort() (int, error) {
	listener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		return 0, errors.Wrap(err, "net.Listen")
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "listener.Close"))
		}
	}(listener)
	port := listener.Addr().(*net.TCPAddr).Port
	return port, nil
}

func RandStr(length int) string {
	const (
		chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteByte(chars[r.Intn(len(chars))])
	}
	return b.String()
}

func ToJsonStr(value any) string {
	result, err := json.MarshalIndent(value, "", "\t")
	if err != nil {
		logrus.Errorln(errors.Wrap(err, "json.MarshalIndent"))
		return ""
	}
	return string(result)
}

func GetTimeStr() string {
	return time.Now().Format("20060102150405")
}

func IsPortInUse(port string) bool {
	localhost := "127.0.0.1"
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", localhost, port))
	if err != nil {
		return true
	}
	_ = listener.Close()
	return false
}
