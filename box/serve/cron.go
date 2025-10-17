package serve

import "github.com/robfig/cron/v3"

var (
	_cron *cron.Cron
)

func Cron() *cron.Cron {
	return _cron
}

func SetCron(c *cron.Cron) {
	_cron = c
}
