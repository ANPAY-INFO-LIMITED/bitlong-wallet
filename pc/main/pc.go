package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/pc/cmd"
	"github.com/wallet/pc/logf"
	"github.com/wallet/pc/pcapi"
	"github.com/wallet/pc/serve"
)

func main() {
	cmd.Init()

	err := cmd.StartServe()
	if err != nil {
		logrus.Fatalln(errors.Wrap(err, "cmd.StartServe"))
	}

	cmd.StartLnuServe()

	err = cmd.StartProxyServe()
	if err != nil {
		logrus.Fatalln(errors.Wrap(err, "cmd.StartProxyServe"))
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logrus.Infoln("Server shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if serve.FrpcStarted() {
		pcapi.LnurlStopFrpc()
	}

	if serve.LitdStarted() {
		err := pcapi.LitdStopDaemon()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "pcapi.LitdStopDaemon"))
			err = pcapi.LndStopDaemon()
			if err != nil {
				logrus.Errorln(errors.Wrap(err, "pcapi.LndStopDaemon"))
			}
		}
		time.Sleep(2 * time.Second)
	}

	if err := serve.Srv().Shutdown(ctx); err != nil {
		logrus.Fatalln(errors.Wrap(err, "serve.Srv().Shutdown"))
	}
	logrus.Infoln("Server exited.")

	if err := serve.SrvLnu().Shutdown(ctx); err != nil {
		logrus.Fatalln(errors.Wrap(err, "serve.SrvLnu().Shutdown"))
	}
	logrus.Infoln("LnuServer exited.")

	if err := serve.SrvProxy().Shutdown(ctx); err != nil {
		logrus.Fatalln(errors.Wrap(err, "serve.SrvProxy().Shutdown"))
	}
	logrus.Infoln("ProxyServer exited.")

	defer func() {
		err := logf.CloseLog()
		if err != nil {
			log.Fatalln(errors.Wrap(err, "cmd.CloseLog"))
		} else {
			fmt.Println("logf closed.")
		}
	}()

}
