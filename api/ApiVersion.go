package api

const (
	PrimaryVersion = "0"
	SubVersion     = "1"
	DevVersion     = "43"
)

const (
	BaseVersion = "v" + PrimaryVersion + "." + SubVersion + "." + DevVersion
	// GetNowTimeStringWithHyphens
	DateTime = "2024-11-11 14:15:31"
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
