package crons

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/box/loggers"
	"github.com/wallet/box/serve"
	"github.com/wallet/box/services"
)

const (
	frpEnable = true
)

func frp() []Job {
	return []Job{
		{
			minute(1),
			func() {
				if frpEnable {
					if !serve.FrpStarted() {

						err := services.FrpBeforeRun()
						if err != nil {
							loggers.Frp().Println(errors.Wrap(err, "services.FrpBeforeRun"))
							logrus.Errorln(errors.Wrap(err, "services.FrpBeforeRun"))
							return
						}
						serve.SetFrpStarted(true)

						if err = services.RunFrp(); err != nil {
							loggers.Frp().Println(errors.Wrap(err, "services.RunFrp"))
							logrus.Errorln(errors.Wrap(err, "services.RunFrp"))
						}

					}
				}
			},
		},
	}
}
