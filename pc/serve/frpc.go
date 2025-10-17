package serve

var (
	_frpcStarted bool
)

func FrpcStarted() bool {
	return _frpcStarted
}

func SetFrpcStarted(started bool) {
	_frpcStarted = started
}
