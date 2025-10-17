package services

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/wallet/box/sc"
	"os"
	"os/exec"
	"strings"
)

func getSeed() (string, error) {

	_, err := os.ReadFile(sc.SeedPath)
	if err != nil {
		return "", errors.Wrapf(err, "os.ReadFile sc.SeedPath")
	}

	password, err := os.ReadFile(sc.PassPath)
	if err != nil {
		return "", errors.Wrapf(err, "os.ReadFile password")
	}
	password = bytes.TrimSpace(password)

	osCmd := exec.Command("openssl", "enc", "-aes-256-cbc", "-d", "-pbkdf2", "-iter", "100000", "-in", sc.SeedPath, "-pass", fmt.Sprintf("pass:%s", password))

	osOutput, err := osCmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "Execute openssl enc -aes-256-cbc -d -pbkdf2 -iter 100000 -in /root/.lit/seed -pass pass:\"$(cat /root/.lit/password)\"")
	}

	return strings.TrimSpace(string(osOutput)), nil
}
