package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/box/config"
	"github.com/wallet/box/crt"
	"github.com/wallet/box/loggers"
	"github.com/wallet/box/routes"
	"github.com/wallet/box/rpc"
	"github.com/wallet/box/serve"
)

const (
	ginMode  = gin.DebugMode
	certFile = ".box/etc/openssl/server.crt"
	keyFile  = ".box/etc/openssl/server.key"
)

func StartServe() error {

	gin.SetMode(ginMode)
	gin.DefaultWriter = config.Writer()
	r := routes.SetupGinRouter()

	bind, port := config.Conf().Serve.Bind, config.Conf().Serve.Port
	if bind == "" {
		bind = config.DefaultServeBind
	}
	if port == 0 {
		port = config.DefaultServePort
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", bind, port),
		Handler: r,
	}
	serve.SetSrv(srv)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return errors.Wrap(err, "os.UserHomeDir")
	}
	certAbsolutePath := filepath.Join(homeDir, crt.CertPath)
	keyAbsolutePath := filepath.Join(homeDir, crt.KeyPath)

	go func() {
		if err := serve.Srv().ListenAndServeTLS(certAbsolutePath, keyAbsolutePath); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalln(errors.Wrap(err, "serve.Srv().ListenAndServeTLS"))
		}
	}()
	logrus.Infoln("Server started.")
	log.Printf("\tServer Listening on:\thttps://%s:%d/\n", bind, port)
	return nil
}

var subscribeOnce sync.Once

func SubscribeInvoiceToUpdate() {
	loggers.Chan().Println("SubscribeInvoiceToUpdate")
	subscribeOnce.Do(func() {
		loggers.Chan().Println("SubscribeInvoiceToUpdate do")
		ctx := context.Background()
		go runSubscribeInvoice(ctx)
	})
}

func runSubscribeInvoice(ctx context.Context) {
	var l rpc.Ln
	var t rpc.Lit
	baseDelay := 2 * time.Second
	maxDelay := 30 * time.Second
	curDelay := baseDelay
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		loggers.Chan().Println("SubscribeInvoiceToUpdate runSubscribeInvoice")
		resp, err := t.SubServerStatusWithCtx(ctx)
		if err != nil {
			loggers.Chan().Println(errors.Wrap(err, "SubServerStatusWithCtx"))
			time.Sleep(curDelay)
			if curDelay < maxDelay {
				curDelay *= 2
				if curDelay > maxDelay {
					curDelay = maxDelay
				}
			}
			continue
		}

		loggers.Chan().Println("SubServerStatusWithCtx success")

		s := resp.SubServers
		lndRunning := s["lnd"] != nil && s["lnd"].Running
		if lndRunning {
			loggers.Chan().Println("SubscribeInvoices lndRunning")
			l.SubscribeInvoices(ctx)
			loggers.Chan().Println("SubscribeInvoices success")
			curDelay = baseDelay
			time.Sleep(curDelay)
			continue
		}

		loggers.Chan().Println("SubscribeInvoices failed")

		time.Sleep(curDelay)
		if curDelay < maxDelay {
			curDelay *= 2
			if curDelay > maxDelay {
				curDelay = maxDelay
			}
		}
	}
}
