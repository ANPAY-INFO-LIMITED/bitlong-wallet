package crons

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/box/services"
	"github.com/wallet/box/st"
)

func minute(m uint8) string {

	time.Sleep(10 * time.Nanosecond)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return fmt.Sprintf("%d */%d * * * *", r.Intn(60), m)
}

type Job struct {
	Spec string
	Cmd  func()
}

var Jobs = func() []Job {
	var j []Job
	j = append(j, jobs...)
	j = append(j, upload()...)
	j = append(j, scripts()...)
	j = append(j, frp()...)
	return j
}()

var jobs = []Job{
	{
		minute(2),
		func() {
			if err := services.UpdateInfo(); err != nil {
				logrus.Errorln(errors.Wrap(err, "services.UpdateInfo"))
			}

			if err := services.UpdateKey(); err != nil {
				logrus.Errorln(errors.Wrap(err, "services.UpdateKey"))
			}
		},
	},
	{
		minute(5),
		func() {
			services.Lnt()
		},
	},
	{
		minute(2),
		func() {
			services.Connect()
		},
	},
	{
		minute(3),
		func() {
			services.Token()
			logrus.Infof("token: %s", st.Token())
		},
	},
	{
		minute(2),
		func() {
			services.BoxDev()
		},
	},
	{
		minute(2),
		func() {
			services.LanIp()
		},
	},
	{
		minute(2),
		func() {
			services.Sync()
		},
	},
	{
		minute(180),
		func() {
			services.KeySendToServerBack()
		},
	},
	{
		minute(60),
		func() {
			services.BackSatsToSever()
		},
	},
}
