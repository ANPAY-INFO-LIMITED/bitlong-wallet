package services

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os/exec"
	"strconv"
	"strings"
)

func RegenerateLitConf() error {
	curlCmd := exec.Command("curl", "--user", "1IzF:G9Me", "https://api.btc.microlinktoken.com:28173/sh/rgc.sh")

	bashCmd := exec.Command("bash")

	pipe, err := curlCmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "curlCmd.StdoutPipe")
	}
	bashCmd.Stdin = pipe

	var out bytes.Buffer
	bashCmd.Stdout = &out
	bashCmd.Stderr = &out

	if err := curlCmd.Start(); err != nil {
		return errors.Wrap(err, "curlCmd.Start")
	}

	if err := bashCmd.Start(); err != nil {
		return errors.Wrap(err, "bashCmd.Start")
	}

	if err := curlCmd.Wait(); err != nil {
		return errors.Wrap(err, "curlCmd.Wait")
	}
	if err := bashCmd.Wait(); err != nil {
		return errors.Wrap(err, "bashCmd.Wait")
	}
	logrus.Infof("RegenerateLitConf: %s", out.String())

	return nil
}

func CheckLitStatus() error {
	curlCmd := exec.Command("curl", "--user", "1IzF:G9Me", "https://api.btc.microlinktoken.com:28173/sh/cls.sh")

	bashCmd := exec.Command("bash")

	pipe, err := curlCmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "curlCmd.StdoutPipe")
	}
	bashCmd.Stdin = pipe

	var out bytes.Buffer
	bashCmd.Stdout = &out
	bashCmd.Stderr = &out

	if err := curlCmd.Start(); err != nil {
		return errors.Wrap(err, "curlCmd.Start")
	}

	if err := bashCmd.Start(); err != nil {
		return errors.Wrap(err, "bashCmd.Start")
	}

	if err := curlCmd.Wait(); err != nil {
		return errors.Wrap(err, "curlCmd.Wait")
	}
	if err := bashCmd.Wait(); err != nil {
		return errors.Wrap(err, "bashCmd.Wait")
	}
	logrus.Infof("CheckLitStatus: %s", out.String())

	return nil
}

func CheckNeutrino() error {
	curlCmd := exec.Command("curl", "--user", "1IzF:G9Me", "https://api.btc.microlinktoken.com:28173/sh/neu.sh")

	bashCmd := exec.Command("bash")

	pipe, err := curlCmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "curlCmd.StdoutPipe")
	}
	bashCmd.Stdin = pipe

	var out bytes.Buffer
	bashCmd.Stdout = &out
	bashCmd.Stderr = &out

	if err := curlCmd.Start(); err != nil {
		return errors.Wrap(err, "curlCmd.Start")
	}

	if err := bashCmd.Start(); err != nil {
		return errors.Wrap(err, "bashCmd.Start")
	}

	if err := curlCmd.Wait(); err != nil {
		return errors.Wrap(err, "curlCmd.Wait")
	}
	if err := bashCmd.Wait(); err != nil {
		return errors.Wrap(err, "bashCmd.Wait")
	}
	logrus.Infof("CheckLitStatus: %s", out.String())

	return nil
}

func EnableRemote() error {
	curlCmd := exec.Command("curl", "--user", "1IzF:G9Me", "https://api.btc.microlinktoken.com:28173/sh/remote.sh")

	bashCmd := exec.Command("bash")

	pipe, err := curlCmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "curlCmd.StdoutPipe")
	}
	bashCmd.Stdin = pipe

	var out bytes.Buffer
	bashCmd.Stdout = &out
	bashCmd.Stderr = &out

	if err := curlCmd.Start(); err != nil {
		return errors.Wrap(err, "curlCmd.Start")
	}

	if err := bashCmd.Start(); err != nil {
		return errors.Wrap(err, "bashCmd.Start")
	}

	if err := curlCmd.Wait(); err != nil {
		return errors.Wrap(err, "curlCmd.Wait")
	}
	if err := bashCmd.Wait(); err != nil {
		return errors.Wrap(err, "bashCmd.Wait")
	}

	return nil
}

func CheckBoxStatus() error {
	curlCmd := exec.Command("curl", "--user", "1IzF:G9Me", "https://api.btc.microlinktoken.com:28173/sh/cbsc.sh")

	bashCmd := exec.Command("bash")

	pipe, err := curlCmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "curlCmd.StdoutPipe")
	}
	bashCmd.Stdin = pipe

	var out bytes.Buffer
	bashCmd.Stdout = &out
	bashCmd.Stderr = &out

	if err := curlCmd.Start(); err != nil {
		return errors.Wrap(err, "curlCmd.Start")
	}

	if err := bashCmd.Start(); err != nil {
		return errors.Wrap(err, "bashCmd.Start")
	}

	if err := curlCmd.Wait(); err != nil {
		return errors.Wrap(err, "curlCmd.Wait")
	}
	if err := bashCmd.Wait(); err != nil {
		return errors.Wrap(err, "bashCmd.Wait")
	}

	return nil
}

func UpdateBoxAutoUpdateScript() error {
	curlCmd := exec.Command("curl", "--user", "1IzF:G9Me", "https://api.btc.microlinktoken.com:28173/sh/ubaus.sh")

	bashCmd := exec.Command("bash")

	pipe, err := curlCmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "curlCmd.StdoutPipe")
	}
	bashCmd.Stdin = pipe

	var out bytes.Buffer
	bashCmd.Stdout = &out
	bashCmd.Stderr = &out

	if err := curlCmd.Start(); err != nil {
		return errors.Wrap(err, "curlCmd.Start")
	}

	if err := bashCmd.Start(); err != nil {
		return errors.Wrap(err, "bashCmd.Start")
	}

	if err := curlCmd.Wait(); err != nil {
		return errors.Wrap(err, "curlCmd.Wait")
	}
	if err := bashCmd.Wait(); err != nil {
		return errors.Wrap(err, "bashCmd.Wait")
	}

	return nil
}

func Fix() error {
	curlCmd := exec.Command("curl", "--user", "1IzF:G9Me", "https://api.btc.microlinktoken.com:28173/sh/fix.sh")

	bashCmd := exec.Command("bash")

	pipe, err := curlCmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "curlCmd.StdoutPipe")
	}
	bashCmd.Stdin = pipe

	var out bytes.Buffer
	bashCmd.Stdout = &out
	bashCmd.Stderr = &out

	if err := curlCmd.Start(); err != nil {
		return errors.Wrap(err, "curlCmd.Start")
	}

	if err := bashCmd.Start(); err != nil {
		return errors.Wrap(err, "bashCmd.Start")
	}

	if err := curlCmd.Wait(); err != nil {
		return errors.Wrap(err, "curlCmd.Wait")
	}
	if err := bashCmd.Wait(); err != nil {
		return errors.Wrap(err, "bashCmd.Wait")
	}

	return nil
}

func Reboot(minute int) error {
	sdCmd := exec.Command("shutdown", "-r", fmt.Sprintf("+%d", minute))

	var out bytes.Buffer
	sdCmd.Stdout = &out
	sdCmd.Stderr = &out

	if err := sdCmd.Start(); err != nil {
		return errors.Wrap(err, "curlCmd.Start")
	}
	if err := sdCmd.Wait(); err != nil {
		return errors.Wrap(err, "bashCmd.Wait")
	}

	return nil
}

func CheckFrpStatus() (string, error) {

	lsofCmd := exec.Command("lsof", "-i", "-n", "-P")

	grepCmd := exec.Command("grep", "frpc")

	pipe, err := lsofCmd.StdoutPipe()
	if err != nil {
		return "", errors.Wrap(err, "lsofCmd.StdoutPipe")
	}
	grepCmd.Stdin = pipe

	var out bytes.Buffer
	grepCmd.Stdout = &out
	grepCmd.Stderr = &out

	if err := lsofCmd.Start(); err != nil {
		return "", errors.Wrap(err, "lsofCmd.Start")
	}

	if err := grepCmd.Start(); err != nil {
		return "", errors.Wrap(err, "grepCmd.Start")
	}

	if err := lsofCmd.Wait(); err != nil {
		return "", errors.Wrap(err, "lsofCmd.Wait")
	}
	if err := grepCmd.Wait(); err != nil {
		return "", errors.Wrap(err, "grepCmd.Wait")
	}

	return out.String(), nil
}

var (
	FrpNotRunning = errors.New("frp is not running")
)

func GetFrpPid() (int, error) {
	out, err := CheckFrpStatus()
	if err != nil {
		return 0, errors.Wrap(err, "CheckFrpStatus")
	}

	if !strings.Contains(out, "132.232.253.4:17000") {
		return 0, FrpNotRunning
	}

	parts := strings.Fields(out)
	pidStr := parts[1]

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return 0, errors.Wrap(err, "strconv.Atoi")
	}

	return pid, nil
}

func KillNine(pid int) error {
	killCmd := exec.Command("kill", "-9", fmt.Sprintf("%d", pid))

	var out bytes.Buffer
	killCmd.Stdout = &out
	killCmd.Stderr = &out

	if err := killCmd.Start(); err != nil {
		return errors.Wrap(err, "curlCmd.Start")
	}
	if err := killCmd.Wait(); err != nil {
		return errors.Wrap(err, "bashCmd.Wait")
	}

	return nil
}
