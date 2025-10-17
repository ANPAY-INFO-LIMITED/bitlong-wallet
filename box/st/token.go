package st

var (
	token string
)

func Token() string {
	return token
}

func Set(t string) {
	token = t
}
