package pcapi

import "github.com/wallet/api"

func AddSession() (*api.Session, error) {
	return api.PcAddSession()
}

func NewSession(sessionType string) (string, error) {
	return api.NewSession(sessionType)
}

func ListSessions(filter api.SessionFilter) ([]*api.Session, error) {
	return api.PcListSessions(filter)
}

func RevokeSession(localPubkey string) error {
	return api.PcRevokeSession(localPubkey)
}

func RevokeAllSessions() error {
	return api.PcRevokeAllSessions()
}
