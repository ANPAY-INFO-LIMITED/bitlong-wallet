package logf

var (
	closeLog func() error
)

func CloseLog() error {
	return closeLog()
}

func Set(f func() error) {
	closeLog = f
}
