package ver

import (
	"github.com/sirupsen/logrus"
)

func Version() string {
	return "v0.2.95"
}

func Print() {
	logrus.Infof("Version: %s", Version())
}
