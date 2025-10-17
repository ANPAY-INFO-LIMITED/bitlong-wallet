package services

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/box/db"
	"github.com/wallet/box/loggers"
	"github.com/wallet/box/models"
	"github.com/wallet/box/rpc"
	"github.com/wallet/box/sc"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

const (
	requiredSat = 3e4 + 1e3
	feeRate     = 4
	localFund   = 2e4
	pushSat     = 1e4
	chanMemo    = ""
	accDefault  = "default"
)

var (
	insufficientSat     = errors.New("insufficient Sat balance, please deposit at least 31000 Sat to your wallet")
	identityPubkeyEmpty = errors.New("identity pubkey is empty, cannot open channel")
	mcEmpty             = errors.New("machine code is empty, cannot open channel")
)

func Connect() {
	if err := ConnectPeer(); err != nil {
		logrus.Errorln(errors.Wrap(err, "ConnectPeer"))
		loggers.Lnt().Println(errors.Wrap(err, "ConnectPeer"))
	}
}

func ConnectPeer() error {
	var l rpc.Ln
	if _, err := l.ConnectPeer(sc.ServerConf().IdentityPubkey, sc.ServerConf().ServerHost); err != nil && !strings.Contains(err.Error(), "already connected to peer") {
		return errors.Wrap(err, "l.ConnectPeer")
	}
	return nil
}

func Lnt() {
	if err := LntOpenChan(); err != nil {
		logrus.Errorln(errors.Wrap(err, "LntOpenChan"))
		loggers.Lnt().Println(errors.Wrap(err, "LntOpenChan"))
	}
}

func LntOpenChan() error {

	st, err := getState()
	if err != nil {
		return errors.Wrap(err, "getState")
	}
	if st != models.LntStateInit {
		loggers.Lnt().Printf("state: %d\n", st)
		return nil
	}

	if err := checkSat(); err != nil {
		return errors.Wrap(err, "checkSat")
	}

	ipk, err := getIpk()
	if err != nil {
		return errors.Wrap(err, "getIpk")
	}
	if ipk == "" {
		return identityPubkeyEmpty
	}

	mc, err := getMc()
	if err != nil {
		return mcEmpty
	}

	if mc == "" {
		return errors.Wrap(mcEmpty, "machine code is empty, cannot open channel")
	}

	op, err := rpc.OpenChan(sc.ServerConf().IdentityPubkey, sc.ServerConf().ServerHost, localFund, feeRate, pushSat, mc)
	if err != nil {
		return errors.Wrap(err, "rpc.OpenChan")
	}

	loggers.Lnt().Printf("OpenChan: %s\n", op)
	logrus.Infof("OpenChan: %s", op)

	tx := db.Sqlite().Begin()

	var l models.Lnt
	if err = tx.Model(&models.Lnt{}).First(&l).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "tx.Model(&models.Lnt{}).First")
	}

	if err = tx.Model(&models.Lnt{}).Where("id = ?", l.ID).Update("state", models.LntStatePending).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "tx.Model(&models.Lnt{}).Update")
	}
	if err = tx.Commit().Error; err != nil {
		return errors.Wrap(err, "tx.Commit()")
	}

	return nil
}

func getState() (models.LntState, error) {

	tx := db.Sqlite().Begin()

	var l models.Lnt
	err := tx.Model(&models.Lnt{}).First(&l).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = tx.Model(&models.Lnt{}).
				Create(&models.Lnt{
					State: models.LntStateInit,
				}).Error
			if err != nil {
				tx.Rollback()
				return models.LntStateUnknown, errors.Wrap(err, "tx.Model(&models.Lnt{}).Create")
			}
			if err = tx.Commit().Error; err != nil {
				return models.LntStateUnknown, errors.Wrap(err, "tx.Commit()")
			}
			return models.LntStateInit, nil
		} else {
			tx.Rollback()
			return models.LntStateUnknown, errors.Wrap(err, "tx.Model(&models.Lnt{}).First")
		}
	}
	tx.Rollback()
	return l.State, nil

}

func checkSat() error {

	var l rpc.Ln
	wb, err := l.WalletBalance()
	if err != nil {
		return errors.Wrap(err, "l.WalletBalance")
	}
	if cf := (wb.AccountBalance[accDefault].ConfirmedBalance) - wb.ReservedBalanceAnchorChan; cf < requiredSat {
		return errors.Wrap(insufficientSat, strconv.Itoa(int(cf)))
	}

	return nil
}
