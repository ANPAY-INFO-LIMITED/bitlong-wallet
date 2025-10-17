package services

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/wallet/box/config"

	"github.com/go-resty/resty/v2"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lntypes"
	"github.com/lightningnetwork/lnd/record"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/box/db"
	"github.com/wallet/box/loggers"
	"github.com/wallet/box/models"
	"github.com/wallet/box/rpc"
	"github.com/wallet/box/sc"
	"github.com/wallet/box/st"
	"github.com/wallet/box/ver"
	"gorm.io/gorm"
)

var (
	noLanIpFound = errors.New("no valid LAN IP address found")
)

func UpdateInfo() error {

	uuid, err := getUUID()
	if err != nil {
		return errors.Wrap(err, "getUUID")
	}
	mid, err := getMachineID()
	if err != nil {
		return errors.Wrap(err, "getMachineID")
	}

	mc, err := getMc()
	if err != nil {
		logrus.Errorln(errors.Wrap(err, "getMc"))
		loggers.BdInfo().Println(errors.Wrap(err, "getMc"))
	}

	ipk, err := getIpk()
	if err != nil {
		logrus.Errorln(errors.Wrap(err, "getIpk"))
		loggers.BdInfo().Println(errors.Wrap(err, "getIpk"))
	}

	tx := db.Sqlite().Begin()

	var i models.Info
	err = tx.Model(&models.Info{}).First(&i).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = tx.Model(&models.Info{}).
				Create(&models.Info{
					BlkUUID:        uuid,
					MachineID:      mid,
					MachineCoding:  mc,
					IdentityPubkey: ipk,
				}).Error
			if err != nil {
				tx.Rollback()
				return errors.Wrap(err, "tx.Model(&models.Info{}).Create")
			}
		} else {
			tx.Rollback()
			return errors.Wrap(err, "tx.Model(&models.Info{}).First")
		}
	} else {
		if uuid != i.BlkUUID || mid != i.MachineID || mc != i.MachineCoding || ipk != i.IdentityPubkey {
			logrus.Infof("ipk: %s, mc: %s, uuid: %s, mid: %s", ipk, mc, uuid, mid)

			if uuid == "" {
				uuid = i.BlkUUID
			}
			if mid == "" {
				mid = i.MachineID
			}
			if mc == "" {
				mc = i.MachineCoding
			}
			if ipk == "" {
				ipk = i.IdentityPubkey
			}

			err = tx.Model(&models.Info{}).
				Where("id = ?", i.ID).
				Updates(map[string]any{
					"blk_uuid":        uuid,
					"machine_id":      mid,
					"machine_coding":  mc,
					"identity_pubkey": ipk,
				}).Error
			if err != nil {
				tx.Rollback()
				return errors.Wrap(err, "tx.Model(&models.Info{}).Updates")
			}
		} else {
			return tx.Rollback().Error
		}
	}

	return tx.Commit().Error
}

func BoxDev() {
	if err := UploadBoxDevInfo(); err != nil {
		logrus.Errorln(errors.Wrap(err, "UploadBoxDevInfo"))
		loggers.BdInfo().Println(errors.Wrap(err, "UploadBoxDevInfo"))
	}
}

func LanIp() {
	if err := UploadLanIp(); err != nil {
		logrus.Errorln(errors.Wrap(err, "UploadLanIp"))
		loggers.BdInfo().Println(errors.Wrap(err, "UploadLanIp"))
	}
}

func UploadBoxDevInfo() error {

	token := st.Token()
	if token == "" {
		return errors.Wrap(invalidToken, "st.Token() is empty")
	}

	client := resty.New()

	info, err := GetInfo()
	if err != nil {
		return errors.Wrap(err, "GetInfo")
	}

	alias, externalIP, err := getAliasAndEip()
	if err != nil {
		logrus.Errorln(errors.Wrap(err, "getAliasAndEip"))
		loggers.BdInfo().Println(errors.Wrap(err, "getAliasAndEip"))
	}

	peerNum, err := getPeerNum()
	if err != nil {
		logrus.Errorln(errors.Wrap(err, "getPeerNum"))
		loggers.BdInfo().Println(errors.Wrap(err, "getPeerNum"))
	}
	chanNum, err := getChanNum()
	if err != nil {
		logrus.Errorln(errors.Wrap(err, "getChanNum"))
		loggers.BdInfo().Println(errors.Wrap(err, "getChanNum"))
	}
	db := db.Sqlite()
	var lnt models.Lnt
	err = db.Model(&models.Lnt{}).First(&lnt).Error

	body := map[string]any{
		"machine_coding":  info.MachineCoding,
		"identity_pubkey": info.IdentityPubkey,
		"blk_uuid":        info.BlkUUID,
		"machine_id":      info.MachineID,
		"alias":           alias,
		"external_ip":     externalIP,
		"peer_num":        peerNum,
		"chan_num":        chanNum,
		"note":            "",
		"box_version":     ver.Version(),
	}

	if lnt.State == models.LntStateInit {
		serverConf := sc.ServerConf()
		body["server_identity_pubkey"] = serverConf.IdentityPubkey
	}

	a := "/box_device/set"
	url := fmt.Sprintf("%s%s", sc.BaseUrl, a)

	var r models.JResult
	var e models.ErrResp

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(token).
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

	if _, ok := resp.Result().(*models.JResult); !ok {
		return invalidRespType
	}

	if r.Error != "" {
		return errors.New(fmt.Sprintf("error: %s", r.Error))
	}

	if resp.StatusCode() != 200 {
		return errors.New(fmt.Sprintf("error: %s", resp.Status()))
	}

	return nil
}

func UploadLanIp() error {

	client := resty.New()

	mc, err := getMc()
	if err != nil {
		return errors.Wrap(err, "getMc")
	}

	ip, err := getLocalIP()
	if err != nil {
		return errors.Wrap(err, "getLocalIP")
	}

	body := map[string]any{
		"mc":     mc,
		"lan_ip": ip,
	}

	a := "/box_ip/upload"
	url := fmt.Sprintf("%s%s", sc.BaseUrl, a)

	var r models.JResult
	var e models.ErrResp

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(sc.BoxIpBasicUser, sc.BoxIpBasicPass).
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

	if _, ok := resp.Result().(*models.JResult); !ok {
		return invalidRespType
	}

	if r.Error != "" {
		return errors.New(fmt.Sprintf("error: %s", r.Error))
	}

	if resp.StatusCode() != 200 {
		return errors.New(fmt.Sprintf("error: %s", resp.Status()))
	}

	return nil
}

func Sync() {
	if err := SyncLnt(); err != nil {
		logrus.Errorln(errors.Wrap(err, "SyncLnt"))
		loggers.BdInfo().Println(errors.Wrap(err, "SyncLnt"))
	}
}

func SyncLnt() error {
	var t rpc.Tap
	if _, err := t.SyncUniverse(sc.UniverseHost, sc.LntAssetID); err != nil {
		return errors.Wrap(err, "t.SyncUniverse")
	}
	return nil
}

func GetInfo() (*models.Info, error) {
	var i models.Info
	tx := db.Sqlite().Begin()
	err := tx.Model(&models.Info{}).First(&i).Error
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "tx.Model(&models.Info{}).First")
	}
	tx.Rollback()
	return &i, nil
}

func getUUID() (string, error) {

	if config.Conf().VirtualUUID != "" {
		return config.Conf().VirtualUUID, nil
	}

	lsblkCmd := exec.Command("lsblk", "-o", "UUID", "-n", "/dev/sda2")
	lsblkOutput, err := lsblkCmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "exec.Command lsblk -o UUID -n /dev/sda2")
	}

	return strings.TrimSpace(string(lsblkOutput)), nil
}

func getMachineID() (string, error) {
	catCmd := exec.Command("cat", "/etc/machine-id")
	catOutput, err := catCmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "exec.Command cat /etc/machine-id")
	}
	return strings.TrimSpace(string(catOutput)), nil
}

func getMc() (string, error) {
	mc, err := os.ReadFile(sc.MachineCodingPath)
	if err != nil {
		return "", errors.Wrapf(err, "os.ReadFile mc")
	}
	return strings.TrimSpace(string(mc)), nil
}

func getPassword() (string, error) {
	password, err := os.ReadFile(sc.PassPath)
	if err != nil {
		return "", errors.Wrapf(err, "os.ReadFile password")
	}
	return strings.TrimSpace(string(password)), nil
}

func CheckPassword(password string) error {
	_password, err := getPassword()
	if err != nil {
		return errors.Wrap(err, "getPassword")
	}
	if password != _password {
		return errors.New("invalid password")
	}
	return nil
}

func GetMc() (string, error) {
	mc, err := os.ReadFile(sc.MachineCodingPath)
	if err != nil {
		return "", errors.Wrapf(err, "os.ReadFile mc")
	}
	return strings.TrimSpace(string(mc)), nil
}

func getIpk() (string, error) {
	var l rpc.Ln
	i, err := l.GetInfo()
	if err != nil {
		return "", errors.Wrap(err, "l.GetInfo")
	}
	return i.IdentityPubkey, nil
}

func getAliasAndEip() (string, string, error) {

	file, err := os.Open(sc.LitConfPath)
	if err != nil {
		return "", "", errors.Wrap(err, "os.Open")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "file.Close"))
		}
	}(file)

	var alias, externalIP string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "lnd.alias" {
			alias = value
		} else if key == "lnd.externalip" {
			externalIP = value
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", errors.Wrap(err, "scanner.Err")
	}

	return alias, externalIP, nil
}

func GetAliasAndEip() (string, string, error) {

	file, err := os.Open(sc.LitConfPath)
	if err != nil {
		return "", "", errors.Wrap(err, "os.Open")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "file.Close"))
		}
	}(file)

	var alias, externalIP string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "lnd.alias" {
			alias = value
		} else if key == "lnd.externalip" {
			externalIP = value
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", errors.Wrap(err, "scanner.Err")
	}

	return alias, externalIP, nil
}

func getPeerNum() (int, error) {
	var l rpc.Ln
	resp, err := l.ListPeers()
	if err != nil {
		return 0, errors.Wrap(err, "l.ListPeers")
	}
	return len(resp.Peers), nil
}

func getChanNum() (int, error) {
	var l rpc.Ln
	resp, err := l.ListChannels(false, false)
	if err != nil {
		return 0, errors.Wrap(err, "l.ListChannels")
	}
	return len(resp.Channels), nil
}

func getLocalIP() (string, error) {

	interfaces, err := net.Interfaces()
	if err != nil {
		return "", errors.Wrap(err, " net.Interfaces")
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip != nil && ip.To4() != nil && !ip.IsLoopback() {
				return ip.String(), nil
			}
		}
	}

	return "", noLanIpFound
}

func KeySendToServerBack() (string, error) {
	var l rpc.Ln
	resp, err := l.ListChannels(false, true)
	if err != nil {
		return "", errors.Wrap(err, "l.ListChannels")
	}
	for _, channel := range resp.Channels {
		if channel.CustomChannelData != nil && channel.LocalBalance > 2000 {
			if channel.LocalBalance-2000 < 354 {
				continue
			}
			_, err := SendPaymentV2ByKeySend(channel.RemotePubkey, channel.LocalBalance-2000, 10, channel.ChanId)
			if err != nil {
				continue
			}
		}
	}
	return "", nil
}

var (
	invalidDestNodePubkey = errors.New("dest node pubkey must be exactly 33 bytes")
	invalidRespStatusCode = errors.New("invalid response status code, expected 2 or 3")
)

func SendPaymentV2ByKeySend(dest string, amt int64, feeLimitSat int64, outgoingChanId uint64) (*lnrpc.Payment, error) {

	destNode, err := hex.DecodeString(dest)
	if err != nil {
		return nil, errors.Wrap(err, "hex.DecodeString")
	}
	if len(destNode) != 33 {
		return nil, errors.Wrap(invalidDestNodePubkey, strconv.Itoa(len(dest)))
	}
	destCustomRecords := make(map[uint64][]byte)
	outgoingChanIds := make([]uint64, 1)

	outgoingChanIds[0] = outgoingChanId

	var rHash []byte
	var preimage lntypes.Preimage

	if _, err := rand.Read(preimage[:]); err != nil {
		return nil, err
	}

	destCustomRecords[record.KeySendType] = preimage[:]
	hash := preimage.Hash()
	rHash = hash[:]
	paymentHash := rHash

	var l rpc.Ln
	resp, err := l.SendPaymentV2(30, destNode, amt, feeLimitSat, destCustomRecords, outgoingChanIds, paymentHash)
	if err != nil {
		return nil, errors.Wrap(err, "l.SendPaymentV2")
	}

	if resp.Status == 2 || resp.Status == 3 || resp.Status == 1 {
		return resp, nil
	}

	return nil, errors.Wrap(invalidRespStatusCode, strconv.Itoa(int(resp.Status)))
}
