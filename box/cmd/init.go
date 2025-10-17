package cmd

import (
	frpLog "github.com/fatedier/frp/pkg/util/log"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/api"
	"github.com/wallet/box/config"
	"github.com/wallet/box/crt"
	"github.com/wallet/box/db"
	"github.com/wallet/box/logf"
	"github.com/wallet/box/loggers"
	"github.com/wallet/box/utils"
	"io"
	"log"
	"os"
	"path"
)

const (
	serverName    = "box"
	confPath      = ".box/etc/srv/config.yaml"
	logPath       = ".box/logs/srv/box.log"
	cronLogPath   = ".box/logs/cron/cron.log"
	lntLogPath    = ".box/logs/lnt/lnt.log"
	tokenLogPath  = ".box/logs/token/token.log"
	bdInfoLogPath = ".box/logs/bdi/bd_info.log"
	litLogPath    = ".box/logs/lit/lit.log"
	boxLogPath    = ".box/logs/box/box.log"
	frpLogPath    = ".box/logs/frp/frp.log"
	chanLogPath   = ".box/logs/chan/chan.log"
)

func Init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logrus.Fatalln(errors.Wrap(err, "os.UserHomeDir"))
	}
	confP := path.Join(homeDir, confPath)
	logP := path.Join(homeDir, logPath)
	cronLogP := path.Join(homeDir, cronLogPath)
	lntLogP := path.Join(homeDir, lntLogPath)
	tokenLogP := path.Join(homeDir, tokenLogPath)
	bdInfoLogP := path.Join(homeDir, bdInfoLogPath)
	litLogP := path.Join(homeDir, litLogPath)
	boxLogP := path.Join(homeDir, boxLogPath)
	frpLogP := path.Join(homeDir, frpLogPath)
	chanLogP := path.Join(homeDir, chanLogPath)

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

	exist, err = utils.PathExist(cronLogP)
	if err != nil {
		logrus.Fatalln(errors.Wrap(err, "PathExist"))
	}
	if !exist {
		err = os.MkdirAll(path.Dir(cronLogP), 0755)
		if err != nil {
			logrus.Fatalln(errors.Wrap(err, "os.MkdirAll"))
		}
	}

	exist, err = utils.PathExist(lntLogP)
	if err != nil {
		logrus.Fatalln(errors.Wrap(err, "PathExist"))
	}
	if !exist {
		err = os.MkdirAll(path.Dir(lntLogP), 0755)
		if err != nil {
			logrus.Fatalln(errors.Wrap(err, "os.MkdirAll"))
		}
	}

	exist, err = utils.PathExist(tokenLogP)
	if err != nil {
		logrus.Fatalln(errors.Wrap(err, "PathExist"))
	}
	if !exist {
		err = os.MkdirAll(path.Dir(tokenLogP), 0755)
		if err != nil {
			logrus.Fatalln(errors.Wrap(err, "os.MkdirAll"))
		}
	}

	exist, err = utils.PathExist(bdInfoLogP)
	if err != nil {
		logrus.Fatalln(errors.Wrap(err, "PathExist"))
	}
	if !exist {
		err = os.MkdirAll(path.Dir(bdInfoLogP), 0755)
		if err != nil {
			logrus.Fatalln(errors.Wrap(err, "os.MkdirAll"))
		}
	}

	exist, err = utils.PathExist(litLogP)
	if err != nil {
		logrus.Fatalln(errors.Wrap(err, "PathExist"))
	}
	if !exist {
		err = os.MkdirAll(path.Dir(litLogP), 0755)
		if err != nil {
			logrus.Fatalln(errors.Wrap(err, "os.MkdirAll"))
		}
	}

	exist, err = utils.PathExist(boxLogP)
	if err != nil {
		logrus.Fatalln(errors.Wrap(err, "PathExist"))
	}
	if !exist {
		err = os.MkdirAll(path.Dir(boxLogP), 0755)
		if err != nil {
			logrus.Fatalln(errors.Wrap(err, "os.MkdirAll"))
		}
	}

	exist, err = utils.PathExist(frpLogP)
	if err != nil {
		logrus.Fatalln(errors.Wrap(err, "PathExist"))
	}
	if !exist {
		err = os.MkdirAll(path.Dir(frpLogP), 0755)
		if err != nil {
			logrus.Fatalln(errors.Wrap(err, "os.MkdirAll"))
		}
	}

	exist, err = utils.PathExist(chanLogP)
	if err != nil {
		logrus.Fatalln(errors.Wrap(err, "PathExist"))
	}
	if !exist {
		err = os.MkdirAll(path.Dir(chanLogP), 0755)
		if err != nil {
			logrus.Fatalln(errors.Wrap(err, "os.MkdirAll"))
		}
	}

	f, err := os.OpenFile(logP, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "os.OpenFile logP"))
	}

	cronF, err := os.OpenFile(cronLogP, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "os.OpenFile cronLogP"))
	}

	lntF, err := os.OpenFile(lntLogP, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "os.OpenFile lntLogP"))
	}

	tokenF, err := os.OpenFile(tokenLogP, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "os.OpenFile tokenLogP"))
	}

	bdInfoF, err := os.OpenFile(bdInfoLogP, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "os.OpenFile bdInfoLogP"))
	}

	litF, err := os.OpenFile(litLogP, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "os.OpenFile litLogP"))
	}

	boxF, err := os.OpenFile(boxLogP, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "os.OpenFile boxLogP"))
	}

	frpF, err := os.OpenFile(frpLogP, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "os.OpenFile frpLogP"))
	}

	chanF, err := os.OpenFile(chanLogP, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "os.OpenFile frpLogP"))
	}

	loggers.SetCron(cronF)
	loggers.SetLnt(lntF)
	loggers.SetToken(tokenF)
	loggers.SetBdInfo(bdInfoF)
	loggers.SetLit(litF)
	loggers.SetBox(boxF)
	loggers.SetFrp(frpF)
	loggers.SetChan(chanF)

	logf.Set(f.Close)
	logf.SetCron(cronF.Close)
	logf.SetLnt(lntF.Close)
	logf.SetToken(tokenF.Close)
	logf.SetBdInfo(bdInfoF.Close)
	logf.SetLit(litF.Close)
	logf.SetBox(boxF.Close)
	logf.SetFrp(frpF.Close)
	logf.SetChan(chanF.Close)

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

	if err = db.Migrate(); err != nil {
		logrus.Fatalln(errors.Wrap(err, "db.Migrate"))
		return
	}

	if err = api.BoxSetPath(api.Mainnet); err != nil {
		logrus.Fatalln(errors.Wrap(err, "api.BoxSetPath"))
	}

}
