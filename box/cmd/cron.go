package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/wallet/box/crons"
	"github.com/wallet/box/loggers"
	"github.com/wallet/box/serve"
)

func StartCron() {
	var c *cron.Cron
	logger := loggers.Cron()
	if logger == nil {
		logrus.Fatalln(errors.New("cron logger is nil"))
	}
	c = cron.New(cron.WithSeconds(), cron.WithLogger(cron.VerbosePrintfLogger(logger)))

	serve.SetCron(c)

	for i, job := range crons.Jobs {
		_, err := c.AddFunc(job.Spec, job.Cmd)
		if err != nil {
			logrus.Fatalln(errors.Wrap(err, fmt.Sprintf("c.AddFunc: %d", i)))
		}
	}

	c.Start()
}
