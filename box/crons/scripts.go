package crons

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/box/config"
	"github.com/wallet/box/services"
)

func scripts() []Job {
	return []Job{
		{
			minute(9),
			func() {
				if !config.Conf().DisableCheck {
					if err := services.CheckLitStatus(); err != nil {
						logrus.Errorln(errors.Wrap(err, "services.CheckLitStatus"))
					}
				}
			},
		},
		{
			minute(180),
			func() {
			},
		},
	}
}
