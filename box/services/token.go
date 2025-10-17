package services

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/box/loggers"
	"github.com/wallet/box/models"
	"github.com/wallet/box/sc"
	"github.com/wallet/box/st"
)

var (
	invalidRespType = errors.New("invalid response type")
	invalidToken    = errors.New("invalid token")
	invalidNPub     = errors.New("invalid nPub")
	invalidBlkUUID  = errors.New("invalid BlkUUID")
)

func Token() {
	if err := updateToken(); err != nil {
		logrus.Errorln(errors.Wrap(err, "updateToken"))
		loggers.Token().Println(errors.Wrap(err, "updateToken"))
	}
}

func updateToken() error {
	token, err := login()
	if err != nil {
		return errors.Wrap(err, "login")
	}

	if token == "" {
		return errors.Wrap(invalidToken, "empty token")
	}

	st.Set(token)

	loggers.Token().Println(token)

	return nil
}
func UpdateToken() error {
	return updateToken()
}

func login() (string, error) {

	client := resty.New()

	nPub, err := GetNPub()
	if err != nil {
		return "", errors.Wrap(err, "GetNPub")
	}
	if nPub == nullSeedNPub {
		return "", errors.Wrap(invalidNPub, "null seed nPub")
	}

	if nPub == "" {
		return "", errors.Wrap(invalidNPub, "empty nPub")
	}

	info, err := GetInfo()
	if err != nil {
		return "", errors.Wrap(err, "GetInfo")
	}

	if info.BlkUUID == "" {
		return "", errors.Wrap(invalidBlkUUID, "empty BlkUUID")
	}

	body := map[string]string{
		"userName": nPub,
		"password": info.BlkUUID,
	}

	a := "/login"
	url := fmt.Sprintf("%s%s", sc.BaseUrl, a)

	var r models.LoginResp
	var e models.ErrResp

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&r).
		SetError(&e).
		Post(url)

	if err != nil {
		return "", errors.Wrap(err, "client.R()")
	}

	if e.Error != "" {
		return "", errors.New(fmt.Sprintf("error: %s", e.Error))
	}

	if _, ok := resp.Result().(*models.LoginResp); !ok {
		return "", invalidRespType
	}

	return r.Token, nil
}
