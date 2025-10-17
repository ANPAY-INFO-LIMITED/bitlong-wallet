package api

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/lightninglabs/lightning-terminal/litrpc"
	"github.com/pkg/errors"
	"github.com/wallet/service/apiConnect"
)

func addSession(t string) (*litrpc.AddSessionResponse, error) {
	var sessionType litrpc.SessionType
	switch t {
	case "readonly":
		sessionType = litrpc.SessionType_TYPE_MACAROON_READONLY
	case "admin":
		sessionType = litrpc.SessionType_TYPE_MACAROON_ADMIN
	default:
		return nil, errors.New("invalid session type, must be readonly or admin")
	}
	conn, clearUp, err := apiConnect.GetConnection("litd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()

	sc := litrpc.NewSessionsClient(conn)
	req := &litrpc.AddSessionRequest{
		Label:                  fmt.Sprintf("s%s", time.Now().Format("20060102150405")),
		SessionType:            sessionType,
		ExpiryTimestampSeconds: 60 * 60 * 24 * 365 * 100,
		MailboxServerAddr:      "mailbox.terminal.lightning.today:443",
	}
	resp, err := sc.AddSession(context.Background(), req)
	if err != nil {
		return nil, errors.Wrap(err, "sc.AddSession")
	}
	return resp, err
}

func NewSession(sessionType string) (string, error) {
	var suffix string
	switch sessionType {
	case "readonly":
		suffix = "Read-Only"
	case "admin":
		suffix = "Admin"
	default:
		return "", errors.New("invalid session type")
	}

	resp, err := addSession(sessionType)
	if err != nil {
		return "", errors.Wrap(err, "addSession")
	}

	mnemonic := resp.Session.PairingSecretMnemonic
	s := fmt.Sprintf("%s||mailbox.terminal.lightning.today:443||%s", mnemonic, suffix)

	encoded := base64.StdEncoding.EncodeToString([]byte(s))

	prefix := "https://terminal.lightning.engineering/#/connect/pair/"

	result := fmt.Sprintf("%s%s", prefix, encoded)

	return result, nil
}

type SessionFilter uint32

const (
	SessionFilterAll SessionFilter = iota
	SessionFilterExpired
	SessionFilterInUse
	SessionFilterRevoked
	SessionFilterCreated
)

var sessionStateMap = map[litrpc.SessionState]SessionFilter{
	litrpc.SessionState_STATE_CREATED: SessionFilterCreated,
	litrpc.SessionState_STATE_EXPIRED: SessionFilterExpired,
	litrpc.SessionState_STATE_IN_USE:  SessionFilterInUse,
	litrpc.SessionState_STATE_REVOKED: SessionFilterRevoked,
}

func listSessions(filter SessionFilter) (*litrpc.ListSessionsResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("litd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()

	sc := litrpc.NewSessionsClient(conn)
	req := &litrpc.ListSessionsRequest{}
	resp, err := sc.ListSessions(context.Background(), req)
	if err != nil {
		return nil, errors.Wrap(err, "sc.ListSessions")
	}

	if filter == SessionFilterAll {
		return resp, nil
	}

	var sessions []*litrpc.Session
	for _, session := range resp.Sessions {
		if sessionStateMap[session.SessionState] != filter {
			continue
		}

		sessions = append(sessions, session)
	}

	return &litrpc.ListSessionsResponse{Sessions: sessions}, nil
}

func revokeSession(localPubkey string) (*litrpc.RevokeSessionResponse, error) {
	conn, clearUp, err := apiConnect.GetConnection("litd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()

	sc := litrpc.NewSessionsClient(conn)
	pubkey, err := hex.DecodeString(localPubkey)
	if err != nil {
		return nil, errors.Wrap(err, "hex.DecodeString")
	}
	req := &litrpc.RevokeSessionRequest{LocalPublicKey: pubkey}
	resp, err := sc.RevokeSession(context.Background(), req)
	if err != nil {
		return nil, errors.Wrap(err, "sc.ListSessions")
	}

	return resp, nil
}

type Session struct {
	ID                     string                      `json:"id"`
	Label                  string                      `json:"label"`
	SessionState           string                      `json:"session_state"`
	SessionType            string                      `json:"session_type"`
	ExpiryTimestampSeconds uint64                      `json:"expiry_timestamp_seconds"`
	MailboxServerAddr      string                      `json:"mailbox_server_addr"`
	DevServer              bool                        `json:"dev_server"`
	PairingSecret          string                      `json:"pairing_secret"`
	PairingSecretMnemonic  string                      `json:"pairing_secret_mnemonic"`
	LocalPublicKey         string                      `json:"local_public_key"`
	RemotePublicKey        string                      `json:"remote_public_key"`
	CreatedAt              uint64                      `json:"created_at"`
	MacaroonRecipe         *MacaroonRecipe             `json:"macaroon_recipe"`
	AccountID              string                      `json:"account_id"`
	AutopilotFeatureInfo   map[string]*litrpc.RulesMap `json:"autopilot_feature_info"`
	RevokedAt              uint64                      `json:"revoked_at"`
	GroupID                string                      `json:"group_id"`
	FeatureConfigs         map[string]string           `json:"feature_configs"`
	PrivacyFlags           uint64                      `json:"privacy_flags"`
}

type MacaroonRecipe struct {
	Permissions []*litrpc.MacaroonPermission
	Caveats     []string
}

func tranSession(s *litrpc.Session) *Session {
	var macaroonRecipe MacaroonRecipe
	if s.MacaroonRecipe != nil {
		macaroonRecipe = MacaroonRecipe{Permissions: s.MacaroonRecipe.Permissions, Caveats: s.MacaroonRecipe.Caveats}
	}
	return &Session{
		ID:                     hex.EncodeToString(s.Id),
		Label:                  s.Label,
		SessionState:           s.SessionState.String(),
		SessionType:            s.SessionType.String(),
		ExpiryTimestampSeconds: s.ExpiryTimestampSeconds,
		MailboxServerAddr:      s.MailboxServerAddr,
		DevServer:              s.DevServer,
		PairingSecret:          hex.EncodeToString(s.PairingSecret),
		PairingSecretMnemonic:  s.PairingSecretMnemonic,
		LocalPublicKey:         hex.EncodeToString(s.LocalPublicKey),
		RemotePublicKey:        hex.EncodeToString(s.RemotePublicKey),
		CreatedAt:              s.CreatedAt,
		MacaroonRecipe:         &macaroonRecipe,
		AccountID:              s.AccountId,
		AutopilotFeatureInfo:   s.AutopilotFeatureInfo,
		RevokedAt:              s.RevokedAt,
		GroupID:                hex.EncodeToString(s.GroupId),
		FeatureConfigs:         s.FeatureConfigs,
		PrivacyFlags:           s.PrivacyFlags,
	}
}

func tranSessions(ss []*litrpc.Session) []*Session {
	var sessions []*Session
	for _, s := range ss {
		sessions = append(sessions, tranSession(s))
	}
	return sessions
}

func revokeAllSessions() error {
	sessionResp, err := listSessions(SessionFilterCreated)
	if err != nil {
		return errors.Wrap(err, "listSessions(SessionFilterCreated)")
	}
	for _, session := range sessionResp.Sessions {
		if session.SessionState != litrpc.SessionState_STATE_CREATED {
			continue
		}
		_, err := revokeSession(hex.EncodeToString(session.LocalPublicKey))
		if err != nil {
			return errors.Wrapf(err, "revokeSession(%s)", hex.EncodeToString(session.LocalPublicKey))
		}
		continue
	}
	return nil
}
