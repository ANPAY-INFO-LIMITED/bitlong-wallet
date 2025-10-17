package middleware

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/wallet/box/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	ipNotAllowed = errors.New("Forbidden: IP address not allowed")
)

func IpWhitelist() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if strings.Contains(clientIP, ":") {
			clientIP = strings.Split(clientIP, ":")[0]
		}

		if !isAllowedIp(clientIP) {
			c.JSON(http.StatusForbidden, models.RespStr{
				Code: models.IpNotAllowed,
				Msg:  ipNotAllowed.Error(),
				Data: models.NullStr,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func isAllowedIp(ip string) bool {
	allowedPrefixes := []string{
		"127.0.0.1",
		"192.168.",
		"172.",
	}

	for _, prefix := range allowedPrefixes {
		if strings.HasPrefix(ip, prefix) {
			if strings.HasPrefix(prefix, "172.") {
				parts := strings.Split(ip, ".")
				if len(parts) >= 2 {
					secondOctet := parts[1]
					var second int
					_, err := fmt.Sscanf(secondOctet, "%d", &second)
					if err != nil {
						return false
					}
					if second >= 16 && second <= 31 {
						return true
					}
					return false
				}
			}
			return true
		}
	}
	return false
}
