package serve

var (
	_frpStarted bool
)

func FrpStarted() bool {
	return _frpStarted
}

func SetFrpStarted(started bool) {
	_frpStarted = started
}

var (
	_remoteStarted bool
)

func RemoteStarted() bool {
	return _remoteStarted
}

func SetRemoteStarted(started bool) {
	_remoteStarted = started
}
