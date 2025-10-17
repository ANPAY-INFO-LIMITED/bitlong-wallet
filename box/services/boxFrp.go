package services

import (
	"fmt"
	"github.com/fatedier/frp/cmd/frpc/sub"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/wallet/box/config"
	"github.com/wallet/box/models"
	"github.com/wallet/box/sc"
	"github.com/wallet/box/st"
	"github.com/wallet/box/utils"
	"os"
	"path"
	"strings"
)

const (
	sshPort     = 22
	boxPort     = config.DefaultProxyServePort
	frpConfPath = ".box/etc/frp/frpc.ini"
)

var (
	invalidPort = errors.New("invalid port")
)

func FrpBeforeRun() error {
	mc, err := getMc()
	if err != nil {
		return errors.Wrap(err, "getMc")
	}

	name, remotePort, err := GenFrpConf()
	if err != nil {
		return errors.Wrap(err, "GenFrpConf")
	}

	ipk, err := getIpk()
	if err != nil {
		return errors.Wrap(err, "getIpk")
	}

	if err = UploadBoxFrp(mc, name, remotePort, ipk); err != nil {
		return errors.Wrap(err, "UploadBoxFrp")
	}

	return nil
}

func Frp() error {
	mc, err := getMc()
	if err != nil {
		return errors.Wrap(err, "getMc")
	}

	name, remotePort, err := GenFrpConf()
	if err != nil {
		return errors.Wrap(err, "GenFrpConf")
	}

	ipk, err := getIpk()
	if err != nil {
		return errors.Wrap(err, "getIpk")
	}

	if err = UploadBoxFrp(mc, name, remotePort, ipk); err != nil {
		return errors.Wrap(err, "UploadBoxFrp")
	}

	if err = RunFrp(); err != nil {
		return errors.Wrap(err, "RunFrp")
	}

	return nil
}

func GenFrpConf() (string, int, error) {
	c, name, remotePort, err := frpConf()
	if err != nil {
		return "", 0, errors.Wrap(err, "frpConf")
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", 0, errors.Wrap(err, "os.UserHomeDir")
	}
	p := path.Join(homeDir, frpConfPath)
	err = utils.CreateFile(p, c)
	if err != nil {
		return "", 0, errors.Wrap(err, "utils.CreateFile")
	}
	sub.SetCfgFile(p)
	return name, remotePort, nil
}

func UploadBoxFrp(mc string, name string, remotePort int, pubKey string) error {

	client := resty.New()

	body := map[string]any{
		"mc":          mc,
		"name":        name,
		"remote_port": remotePort,
		"pub_key":     pubKey,
	}

	a := "/box_frp/set"
	url := fmt.Sprintf("%s%s", sc.BaseUrl, a)

	var r models.Resp
	var e models.ErrResp

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(st.Token()).
		SetBody(body).
		SetResult(&r).
		SetError(&e).
		Post(url)

	if err != nil {
		return errors.Wrap(err, "client.R()")
	}

	if e.Error != "" {
		return errors.New(fmt.Sprintf("error: %s", e.Error))
	}

	if _, ok := resp.Result().(*models.Resp); !ok {
		return invalidRespType
	}

	if r.Msg != "" {
		return errors.New(fmt.Sprintf("error: %s", r.Msg))
	}

	if resp.StatusCode() != 200 {
		return errors.New(fmt.Sprintf("error: %s", resp.Status()))
	}

	return nil
}

func RunFrp() error {
	return sub.ExecuteNoExit()
}

func StopFrp() {
	sub.StopSrv()
}

func frpConf() (string, string, int, error) {
	name, err := getName()
	if err != nil {
		return "", "", 0, errors.Wrap(err, "getName")
	}
	remotePort, err := getRemotePort()
	if err != nil {
		return "", "", 0, errors.Wrap(err, "getRemotePort")
	}

	return conf(name, remotePort), name, remotePort, nil
}

func conf(name string, remotePort int) string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("serverAddr = \"%s\"\n", sc.ServerAddrOne))
	s.WriteString(fmt.Sprintf("serverPort = %d\n", sc.ServerPort))
	s.WriteString(fmt.Sprintf("auth.token = \"%s\"\n", sc.AuthToken))
	s.WriteString("\n")
	s.WriteString("[[proxies]]\n")
	s.WriteString(fmt.Sprintf("name = \"%s\"\n", name))
	s.WriteString(fmt.Sprintf("type = \"%s\"\n", "tcp"))
	s.WriteString(fmt.Sprintf("localIP = \"%s\"\n", "127.0.0.1"))
	s.WriteString(fmt.Sprintf("localPort = %d\n", boxPort))
	s.WriteString(fmt.Sprintf("remotePort = %d\n", remotePort))
	return s.String()

}

func getRemotePort() (int, error) {

	client := resty.New()

	a := "/port/available"
	url := fmt.Sprintf("%s%s", sc.BaseUrl, a)

	var r models.Resp

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(sc.BoxIpBasicUser, sc.BoxIpBasicPass).
		SetResult(&r).
		SetError(&r).
		Get(url)

	if err != nil {
		return 0, errors.Wrap(err, "client.R()")
	}

	if r.Msg != "" {
		return 0, errors.New(r.Msg)
	}

	if _, ok := resp.Result().(*models.Resp); !ok {
		return 0, invalidRespType
	}

	if resp.StatusCode() != 200 {
		return 0, errors.New(fmt.Sprintf("error: %s", resp.Status()))
	}

	if r.Data == 0 {
		return 0, invalidPort
	}

	return int(r.Data.(float64)), nil
}

func getName() (string, error) {
	mc, err := getMc()
	if err != nil {
		return "", errors.Wrap(err, "getMc")
	}
	utils.GetTimeStr()
	ip, err := getLocalIP()
	if err != nil {
		return "", errors.Wrap(err, "getLocalIP")
	}

	return fmt.Sprintf("%s-%s-%s", mc, utils.GetTimeStr(), ip), nil
}
