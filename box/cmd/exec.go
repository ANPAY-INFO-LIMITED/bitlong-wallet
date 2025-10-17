package cmd

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/box/loggers"
	"github.com/wallet/box/services"
	"github.com/wallet/box/ver"
)

func Exec() {

	ver.Print()

	if err := services.UpdateInfo(); err != nil {
		logrus.Errorln(errors.Wrap(err, "services.UpdateInfo"))
	}

	services.Token()

	if err := services.RegenerateLitConf(); err != nil {
		loggers.Lit().Println(errors.Wrap(err, "services.RegenerateLitConf"))
		logrus.Errorln(errors.Wrap(err, "services.RegenerateLitConf"))
	}

	if err := services.CheckBoxStatus(); err != nil {
		loggers.Box().Println(errors.Wrap(err, "services.CheckBoxStatus"))
		logrus.Errorln(errors.Wrap(err, "services.CheckBoxStatus"))
	}

	if err := services.UpdateBoxAutoUpdateScript(); err != nil {
		loggers.Box().Println(errors.Wrap(err, "services.UpdateBoxAutoUpdateScript"))
		logrus.Errorln(errors.Wrap(err, "services.UpdateBoxAutoUpdateScript"))
	}

	if err := services.Fix(); err != nil {
		loggers.Box().Println(errors.Wrap(err, "services.Fix"))
		logrus.Errorln(errors.Wrap(err, "services.Fix"))
	}

}
