package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/wallet/service"
	"github.com/wallet/service/untils"
)

const (
	LoginUrl       = "/login"
	RefreshUrl     = "/refresh"
	GetNonceUrl    = "/getNonce"
	GetDeviceIdUrl = "/getDeviceId"
	reChangeUrl    = "/reChange"
)

func GetServerHost() string {
	return Cfg.BtlServerHost
}

func Login(username, password string) (string, error) {
	url := GetServerHost() + LoginUrl
	return login(url, username, password)
}

func ReChange(username, password string) (string, error) {
	url := GetServerHost() + reChangeUrl
	return reChange(url, username, password)
}
func Refresh(username, password string) (string, error) {
	url := GetServerHost() + RefreshUrl
	return refresh(url, username, password)
}
func Nonce(username string) (string, error) {
	url := GetServerHost() + GetNonceUrl
	return getNonce(url, username)
}
func DeviceID(username, nonce string) (string, error) {
	url := GetServerHost() + GetDeviceIdUrl
	return getDeviceID(url, nonce, username)
}

func getNonce(url string, username string) (string, error) {
	nonce_Info := struct {
		Username string `json:"userName"`
		Nonce    string `json:"nonce"`
	}{
		Username: username,
		Nonce:    "",
	}
	requestBody, _ := json.Marshal(nonce_Info)
	a, err := SendPostRequest(url, "", requestBody)
	if err != nil {
		return "", err
	}
	result := struct {
		Error string `json:"error"`
		Nonce string `json:"nonce"`
	}{}
	err = json.Unmarshal(a, &result)
	if err != nil {
		fmt.Println("An error occurred while unmarshalling the response body:", err)
	}
	if result.Error != "" {
		return "", fmt.Errorf("%v", result.Error)
	}
	return result.Nonce, err
}
func getDeviceID(url string, nonce, username string) (string, error) {
	nonce_Info := struct {
		Username string `json:"userName"`
		Nonce    string `json:"nonce"`
	}{
		Username: username,
		Nonce:    nonce,
	}
	requestBody, _ := json.Marshal(nonce_Info)
	a, err := SendPostRequest(url, "", requestBody)
	if err != nil {
		return "", err
	}
	result := struct {
		Error           string `json:"error"`
		EncryptDeviceID string `json:"encryptDeviceID"`
		EncodedSalt     string `json:"encodedSalt"`
	}{}
	err = json.Unmarshal(a, &result)
	if err != nil {
		fmt.Println("An error occurred while unmarshalling the response body:", err)
	}
	if result.Error != "" {
		return "", fmt.Errorf("%v", result.Error)
	}
	fmt.Println("get encryptDeviceID:", result.EncryptDeviceID)
	fmt.Println("get EncodedSalt:", result.EncodedSalt)
	deviceID, err := service.BuildDecrypt(result.EncodedSalt, result.EncryptDeviceID)
	if err != nil {
		fmt.Println("get deviceID err:", err)
	}
	return deviceID, err
}

func login(url string, username string, password string) (string, error) {
	user := struct {
		Username string `gorm:"unique;column:user_name" json:"userName"` // 正确地将unique和column选项放在同一个gorm标签内
		Password string `gorm:"column:password" json:"password"`
	}{
		Username: username,
		Password: password,
	}
	requestBody, _ := json.Marshal(user)
	a, err := SendPostRequest(url, "", requestBody)
	if err != nil {
		return "", err
	}
	result := struct {
		Error string `json:"error"`
		Token string `json:"token"`
	}{}
	err = json.Unmarshal(a, &result)
	if err != nil {
		fmt.Println("An error occurred while unmarshalling the response body:", err)
	}
	if result.Error != "" {
		return "", fmt.Errorf("%v", result.Error)
	}
	return result.Token, err
}

func refresh(url string, username string, password string) (string, error) {
	user := struct {
		Username string `gorm:"unique;column:user_name" json:"userName"` // 正确地将unique和column选项放在同一个gorm标签内
		Password string `gorm:"column:password" json:"password"`
	}{
		Username: username,
		Password: untils.GenerateExtMD5WithSalt(password),
	}
	requestBody, _ := json.Marshal(user)
	a, err := SendPostRequest(url, "", requestBody)
	if err != nil {
		return "", err
	}
	result := struct {
		Error string `json:"error"`
		Token string `json:"token"`
	}{}
	err = json.Unmarshal(a, &result)
	if err != nil {
		fmt.Println("An error occurred while unmarshalling the response body:", err)
	}
	if result.Error != "" {
		return "", fmt.Errorf("%v", result.Error)
	}
	return result.Token, err
}

func upLoadLog(url string, token string, requestBody []byte) (string, error) {
	return "", nil
}

func SendPostRequest(url string, token string, requestBody []byte) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("An error occurred while creating an HTTP request:", err)
		return nil, err
	}
	req = req.WithContext(ctx)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("An error occurred while sending a POST request:", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("An error occurred while closing the HTTP response body:", err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil

}
func reChange(url string, username string, password string) (string, error) {
	user := struct {
		Username string `gorm:"unique;column:user_name" json:"userName"` // 正确地将unique和column选项放在同一个gorm标签内
		Password string `gorm:"column:password" json:"password"`
	}{
		Username: username,
		Password: password,
	}
	requestBody, _ := json.Marshal(user)
	a, err := SendPostRequest(url, "", requestBody)
	if err != nil {
		return "", err
	}
	result := struct {
		Error string `json:"error"`
		Token string `json:"token"`
	}{}
	err = json.Unmarshal(a, &result)
	if err != nil {
		fmt.Println("An error occurred while unmarshalling the response body:", err)
	}
	if result.Error != "" {
		return "", fmt.Errorf("%v", result.Error)
	}
	return result.Token, err
}
