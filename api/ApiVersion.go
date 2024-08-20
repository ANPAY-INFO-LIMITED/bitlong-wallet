package api

const (
	PrimaryVersion = "0"
	SubVersion     = "0"
	DevVersion     = "7"
)

const (
	BaseVersion = "v" + PrimaryVersion + "." + SubVersion + "." + DevVersion
	DateTime    = "2024-8-20 11:23:41"
)

func apiVersionWithMaker(maker string) string {
	if maker == "" {
		return BaseVersion + "-" + DateTime
	} else {
		return BaseVersion + "-" + maker + "-" + DateTime
	}
}

// GetApiVersion
// @Description: Get api version
func GetApiVersion() string {
	return apiVersionWithMaker("")
}
