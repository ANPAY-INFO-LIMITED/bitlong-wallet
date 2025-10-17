package routes

import (
	"bytes"
	"crypto/tls"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/pc/models"
	"io"
	"net/http"
)

var insecureClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

func Proxy(r *gin.RouterGroup) *gin.RouterGroup {

	r.POST("/Proxy", func(c *gin.Context) {

		url := c.Query("url")

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			err = errors.Wrap(err, "proxy io.ReadAll(c.Request.Body)")
			logrus.Errorln(err)
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.ProxyIoReadAllErr,
				Msg:  err.Error(),
				Data: "",
			})
			return
		}

		req, err := http.NewRequest("POST", url, bytes.NewReader(body))
		if err != nil {
			err = errors.Wrap(err, "proxy http.NewRequest")
			logrus.Errorln(err)
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.ProxyHttpNewRequestErr,
				Msg:  err.Error(),
				Data: "",
			})
			return
		}

		req.Header.Set("Content-Type", c.GetHeader("Content-Type"))

		resp, err := insecureClient.Do(req)
		if err != nil {
			err = errors.Wrap(err, "proxy insecureClient.Do")
			logrus.Errorln(err)
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.ProxyInsecureClientDoErr,
				Msg:  err.Error(),
				Data: "",
			})
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				logrus.Errorf(errors.Wrap(err, "Body.Close").Error())
			}
		}(resp.Body)

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			err = errors.Wrap(err, "proxy io.ReadAll(resp.Body)")
			logrus.Errorln(err)
			c.JSON(http.StatusOK, models.RespStr{
				Code: models.ProxyIoReadAllErr,
				Msg:  err.Error(),
				Data: "",
			})
			return
		}

		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
	})

	return r
}
