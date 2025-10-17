package api

import (
	"github.com/pkg/errors"
)

func PcAddSession() (*Session, error) {
	resp, err := addSession("admin")
	if err != nil {
		return nil, errors.Wrap(err, "addSession")
	}
	return tranSession(resp.Session), nil
}

func PcListSessions(filter SessionFilter) ([]*Session, error) {
	resp, err := listSessions(filter)
	if err != nil {
		return nil, errors.Wrap(err, "listSessions")
	}
	return tranSessions(resp.Sessions), nil
}

func PcRevokeSession(localPubkey string) error {
	_, err := revokeSession(localPubkey)
	if err != nil {
		return errors.Wrap(err, "revokeSession")
	}
	return nil
}

func PcRevokeAllSessions() error {
	return revokeAllSessions()
}
