package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/pc/config"
	"github.com/wallet/pc/routes"
	"github.com/wallet/pc/serve"
	"log"
	"net/http"
)

func StartProxyServe() error {

	gin.SetMode(ginMode)
	gin.DefaultWriter = config.Writer()
	r := routes.SetupGinDefaultRouter()
	apiRv1 := r.Group(routes.Prefix)
	routes.Proxy(apiRv1)

	bind, port := config.Conf().ProxyServe.Bind, config.Conf().ProxyServe.Port
	if bind == "" {
		bind = config.DefaultProxyServeBind
	}
	if port == 0 {
		port = config.DefaultProxyServePort
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", bind, port),
		Handler: r,
	}
	serve.SetSrvProxy(srv)

	go func() {
		if err := serve.SrvProxy().ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalln(errors.Wrap(err, "serve.SrvProxy().ListenAndServe"))
		}
	}()
	logrus.Infoln("ProxyServer started.")
	log.Printf("\tProxyServer Listening on:\thttp://%s:%d/\n", bind, port)

	return nil
}
