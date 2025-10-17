package serve

import "net/http"

var (
	srv      *http.Server
	srvProxy *http.Server
	srvLnu   *http.Server
)

func Srv() *http.Server {
	return srv
}

func SetSrv(s *http.Server) {
	srv = s
}

func SrvLnu() *http.Server {
	return srvLnu
}

func SetSrvLnu(s *http.Server) {
	srvLnu = s
}

func SrvProxy() *http.Server {
	return srvProxy
}

func SetSrvProxy(s *http.Server) {
	srvProxy = s
}
