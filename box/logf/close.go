package logf

import "github.com/pkg/errors"

var (
	closeLog       func() error
	closeLogCron   func() error
	closeLogLnt    func() error
	closeLogToken  func() error
	closeLogBdInfo func() error
	closeLogLit    func() error
	closeLogBox    func() error
	closeLogFrp    func() error
	closeLogChan   func() error
)

func CloseLog() error {

	if err := closeLog(); err != nil {
		return errors.Wrap(err, "closeLog")
	}
	if err := closeLogCron(); err != nil {
		return errors.Wrap(err, "closeLogCron")
	}
	if err := closeLogLnt(); err != nil {
		return errors.Wrap(err, "closeLogLnt")
	}
	if err := closeLogToken(); err != nil {
		return errors.Wrap(err, "closeLogToken")
	}
	if err := closeLogBdInfo(); err != nil {
		return errors.Wrap(err, "closeLogBdInfo")
	}
	if err := closeLogLit(); err != nil {
		return errors.Wrap(err, "closeLogLit")
	}
	if err := closeLogBox(); err != nil {
		return errors.Wrap(err, "closeLogBox")
	}
	if err := closeLogFrp(); err != nil {
		return errors.Wrap(err, "closeLogFrp")
	}
	if err := closeLogChan(); err != nil {
		return errors.Wrap(err, "closeLogChan")
	}

	return nil
}

func Set(f func() error) {
	closeLog = f
}

func SetCron(f func() error) {
	closeLogCron = f
}

func SetLnt(f func() error) {
	closeLogLnt = f
}

func SetToken(f func() error) {
	closeLogToken = f
}

func SetBdInfo(f func() error) {
	closeLogBdInfo = f
}

func SetLit(f func() error) {
	closeLogLit = f
}

func SetBox(f func() error) {
	closeLogBox = f
}

func SetFrp(f func() error) {
	closeLogFrp = f
}

func SetChan(f func() error) {
	closeLogChan = f
}
