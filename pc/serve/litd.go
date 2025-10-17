package serve

var (
	_litdStarted bool
)

func LitdStarted() bool {
	return _litdStarted
}

func SetLitdStarted(started bool) {
	_litdStarted = started
}
