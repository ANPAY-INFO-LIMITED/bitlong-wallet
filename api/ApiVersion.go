package api

const (
	PrimaryVersion = "0"
	SubVersion     = "1"
	DevVersion     = "11"
)

const (
	BaseVersion = "v" + PrimaryVersion + "." + SubVersion + "." + DevVersion
	DateTime    = "2024-09-13-17-11-18-158373"
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
