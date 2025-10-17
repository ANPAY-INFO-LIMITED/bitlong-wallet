package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/pc/config"
	"github.com/wallet/pc/crt"
	"github.com/wallet/pc/routes"
	"github.com/wallet/pc/serve"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	ginMode = gin.DebugMode
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
