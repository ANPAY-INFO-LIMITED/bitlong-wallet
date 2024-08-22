package api

const (
	PrimaryVersion = "0"
	SubVersion     = "1"
	DevVersion     = "3"
)

const (
	BaseVersion = "v" + PrimaryVersion + "." + SubVersion + "." + DevVersion
	DateTime    = "2024-8-22 09:38:20"
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
