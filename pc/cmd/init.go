package cmd

import (
	frpLog "github.com/fatedier/frp/pkg/util/log"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/pc/config"
	"github.com/wallet/pc/crt"
	"github.com/wallet/pc/db"
	"github.com/wallet/pc/logf"
	"github.com/wallet/pc/utils"
	"io"
	"log"
	"os"
	"path"
)

const (
	serverName = "bitlong_pc"
	confPath   = "bitlong_pc/etc/srv/config.yaml"
	logPath    = "bitlong_pc/logs/srv/bitlong_pc.log"
)

func Init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logrus.Fatalln(errors.Wrap(err, "os.UserHomeDir"))
	}
	confP := path.Join(homeDir, confPath)
	logP := path.Join(homeDir, logPath)

	exist, err := utils.PathExist(logP)
	if err != nil {
		logrus.Fatalln(errors.Wrap(err, "PathExist"))
	}
	if !exist {
		err = os.MkdirAll(path.Dir(logP), 0755)
		if err != nil {
			logrus.Fatalln(errors.Wrap(err, "os.MkdirAll"))
		}
	}
	f, err := os.OpenFile(logP, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "os.OpenFile"))
	}

	logf.Set(f.Close)

	multiWriter := io.MultiWriter(os.Stdout, f)
	config.SetWriter(multiWriter)
	frpLog.SetLogWriter(multiWriter)

	logrus.SetOutput(config.Writer())
	logrus.SetReportCaller(true)

	log.SetOutput(config.Writer())
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("======================================== %s ========================================\n", serverName)

	exist, err = utils.PathExist(confP)
	if err != nil {
		logrus.Fatalln(errors.Wrap(err, "PathExist"))
	}
	if !exist {
		err = config.CreateConfSample(confP)
		if err != nil {
			logrus.Fatalln(errors.Wrap(err, "config.CreateConfSample"))
		}
	}
	_, err = config.LoadConfig(confP)
	if err != nil {
		logrus.Fatalln(errors.Wrap(err, "config.LoadConfig"))
	}

	err = crt.CheckCertExist()
	if err != nil {
		logrus.Infoln(errors.Wrap(err, "crt.CheckCertExist"))
		_, err = crt.GenerateSelfSignedTlsCert()
		if err != nil {
			logrus.Fatalln(errors.Wrap(err, "crt.GenerateSelfSignedTlsCert"))
		} else {
			logrus.Infoln("Self-signed TLS certificate generated successfully.")
		}
	} else {
		logrus.Infoln("Certificate check passed.")
	}

	if err = db.InitSqlite(); err != nil {
		logrus.Fatalln(errors.Wrap(err, "db.InitSqlite"))
	}

}
