package cmd

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/api"
	"github.com/wallet/box/config"
	"github.com/wallet/box/serve"
	"log"
	"net/http"
)

func StartLnuServe() {

	srv := api.PcLnurlSetServer(config.Writer(), config.Conf().LnuServe.BasicAuthUser, config.Conf().LnuServe.BasicAuthPass)

	serve.SetSrvLnu(srv)

	go func() {
		if err := serve.SrvLnu().ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalln(errors.Wrap(err, "serve.SrvLnu().ListenAndServe"))
		}
	}()

	logrus.Infoln("LnuServer started.")
	log.Printf("\tLnuServer Listening on:\thttp://%s:%d/\n", api.LnuBind, api.LnuPort)

}
