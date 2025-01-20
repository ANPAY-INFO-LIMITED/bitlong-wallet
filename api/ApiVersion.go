package api

const (
	PrimaryVersion = "0"
	SubVersion     = "1"
	DevVersion     = "54"
)

const (
	BaseVersion = "v" + PrimaryVersion + "." + SubVersion + "." + DevVersion
	// GetNowTimeStringWithHyphens
	DateTime = "2025-1-20 17:25:01"
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
