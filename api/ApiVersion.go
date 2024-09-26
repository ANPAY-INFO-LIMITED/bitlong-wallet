package api

const (
	PrimaryVersion = "0"
	SubVersion     = "1"
	DevVersion     = "20"
)

const (
	BaseVersion = "v" + PrimaryVersion + "." + SubVersion + "." + DevVersion
	// GetNowTimeStringWithHyphens
	DateTime = "2024-9-26 15:34:25"
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
