package models

import "github.com/lightninglabs/lightning-terminal/litrpc"

func ToSubServerStatus(status map[string]*litrpc.SubServerStatus) SubServers {
	subServers := make(SubServers)
	for k, v := range status {
		if k == "taproot-assets" {
			k = "taproot_assets"
		}
		subServers[k] = &SubServerStatus{
			Disabled:     v.Disabled,
			Running:      v.Running,
			Error:        v.Error,
			CustomStatus: v.CustomStatus,
		}
	}
	return subServers
}
