package crons

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/box/services"
	"github.com/wallet/box/st"
)

const (
	di = "box_device"
)

func upload() []Job {

	deviceId := di

	return []Job{
		{
			minute(10),
			func() {
				if err := services.UploadWalletBalance(st.Token, deviceId); err != nil {
					logrus.Errorln(errors.Wrap(err, "services.UploadWalletBalance"))
				}
			},
		},
		{
			minute(10),
			func() {
				if err := services.UploadAssetManagedUtxos(st.Token, deviceId); err != nil {
					logrus.Errorln(errors.Wrap(err, "services.UploadAssetManagedUtxos"))
				}
			},
		},
		{
			minute(10),
			func() {
				if err := services.UploadBtcListUnspent(st.Token); err != nil {
					logrus.Errorln(errors.Wrap(err, "services.UploadBtcListUnspent"))
				}
			},
		},
		{
			minute(2),
			func() {
				if err := services.UploadBoxChanInfo(st.Token); err != nil {
					logrus.Errorln(errors.Wrap(err, "services.UploadBtcListUnspent"))
				}
			},
		},
		{
			minute(178),
			func() {
				if err := services.UpdateTotalAssetPush(st.Token); err != nil {
					logrus.Errorln(errors.Wrap(err, "services.UpdateTotalAssetPush"))
				}
			},
		},
	}
}
