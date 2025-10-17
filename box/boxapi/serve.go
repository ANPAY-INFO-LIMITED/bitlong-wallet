package boxapi

import (
	"github.com/wallet/pc/pcapi"
)

func GetApiVersion() string {
	return pcapi.GetApiVersion()
}
