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
	"github.com/wallet/box/cmd"
	"github.com/wallet/box/logf"
	"github.com/wallet/box/serve"
	"github.com/wallet/box/services"
)

func main() {
	cmd.Init()

	if err := cmd.StartServe(); err != nil {
		logrus.Fatalln(errors.Wrap(err, "cmd.StartServe"))
	}
	cmd.SubscribeInvoiceToUpdate()
	cmd.StartLnuServe()

	if err := cmd.StartProxyServe(); err != nil {
		logrus.Fatalln(errors.Wrap(err, "cmd.StartProxyServe"))
	}

	cmd.Exec()

	cmd.StartCron()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logrus.Infoln("Server shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if serve.FrpStarted() {
		services.StopFrp()
	}

	serve.Cron().Stop()

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
