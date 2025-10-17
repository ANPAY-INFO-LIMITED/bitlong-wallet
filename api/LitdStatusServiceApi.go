package api

import (
	"context"

	"github.com/pkg/errors"

	"github.com/lightninglabs/lightning-terminal/litrpc"
	"github.com/wallet/service/apiConnect"
)

func SubServerStatus() string {
	response, err := subServerStatus()
	if err != nil {
		return MakeJsonErrorResult(subServerStatusErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}

func GetTapdStatus() bool {
	response, err := subServerStatus()
	if err != nil {
		return false
	}
	if len(response.SubServers) == 0 {
		return false
	}
	return response.SubServers["taproot-assets"].Running
}

func GetLitStatus() bool {
	response, err := subServerStatus()
	if err != nil {
		return false
	}
	if len(response.SubServers) == 0 {
		return false
	}
	return response.SubServers["lit"].Running
}

func subServerStatus() (*litrpc.SubServerStatusResp, error) {
	conn, clearUp, err := apiConnect.GetConnection("litd", true)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()

	client := litrpc.NewStatusClient(conn)
	request := &litrpc.SubServerStatusReq{}
	response, err := client.SubServerStatus(context.Background(), request)
	return response, err
}

type SubServerStatusInfo struct {
	Accounts struct {
		Disabled     bool   `json:"disabled"`
		Running      bool   `json:"running"`
		Error        string `json:"error"`
		CustomStatus string `json:"custom_status"`
	} `json:"accounts"`
	Faraday struct {
		Disabled     bool   `json:"disabled"`
		Running      bool   `json:"running"`
		Error        string `json:"error"`
		CustomStatus string `json:"custom_status"`
	} `json:"faraday"`
	Lit struct {
		Disabled     bool   `json:"disabled"`
		Running      bool   `json:"running"`
		Error        string `json:"error"`
		CustomStatus string `json:"custom_status"`
	} `json:"lit"`
	Lnd struct {
		Disabled     bool   `json:"disabled"`
		Running      bool   `json:"running"`
		Error        string `json:"error"`
		CustomStatus string `json:"custom_status"`
	} `json:"lnd"`
	Loop struct {
		Disabled     bool   `json:"disabled"`
		Running      bool   `json:"running"`
		Error        string `json:"error"`
		CustomStatus string `json:"custom_status"`
	} `json:"loop"`
	Pool struct {
		Disabled     bool   `json:"disabled"`
		Running      bool   `json:"running"`
		Error        string `json:"error"`
		CustomStatus string `json:"custom_status"`
	} `json:"pool"`
	TaprootAssets struct {
		Disabled     bool   `json:"disabled"`
		Running      bool   `json:"running"`
		Error        string `json:"error"`
		CustomStatus string `json:"custom_status"`
	} `json:"taproot-assets"`
}

func getSubServerStatusInfo() (subServerStatusInfo SubServerStatusInfo, err error) {
	response, err := subServerStatus()
	if err != nil {
		return subServerStatusInfo, AppendErrorInfo(err, "subServerStatus")
	}
	for k, v := range response.SubServers {
		if k == "accounts" {
			if v != nil {
				subServerStatusInfo.Accounts.Disabled = v.Disabled
				subServerStatusInfo.Accounts.Running = v.Running
				subServerStatusInfo.Accounts.Error = v.Error
				subServerStatusInfo.Accounts.CustomStatus = v.CustomStatus
			}
		}
		if k == "faraday" {
			if v != nil {
				subServerStatusInfo.Faraday.Disabled = v.Disabled
				subServerStatusInfo.Faraday.Running = v.Running
				subServerStatusInfo.Faraday.Error = v.Error
				subServerStatusInfo.Faraday.CustomStatus = v.CustomStatus
			}
		}
		if k == "lit" {
			if v != nil {
				subServerStatusInfo.Lit.Disabled = v.Disabled
				subServerStatusInfo.Lit.Running = v.Running
				subServerStatusInfo.Lit.Error = v.Error
				subServerStatusInfo.Lit.CustomStatus = v.CustomStatus
			}
		}
		if k == "lnd" {
			if v != nil {
				subServerStatusInfo.Lnd.Disabled = v.Disabled
				subServerStatusInfo.Lnd.Running = v.Running
				subServerStatusInfo.Lnd.Error = v.Error
				subServerStatusInfo.Lnd.CustomStatus = v.CustomStatus
			}
		}
		if k == "loop" {
			if v != nil {
				subServerStatusInfo.Loop.Disabled = v.Disabled
				subServerStatusInfo.Loop.Running = v.Running
				subServerStatusInfo.Loop.Error = v.Error
				subServerStatusInfo.Loop.CustomStatus = v.CustomStatus
			}
		}
		if k == "pool" {
			if v != nil {
				subServerStatusInfo.Pool.Disabled = v.Disabled
				subServerStatusInfo.Pool.Running = v.Running
				subServerStatusInfo.Pool.Error = v.Error
				subServerStatusInfo.Pool.CustomStatus = v.CustomStatus
			}
		}
		if k == "taproot-assets" {
			if v != nil {
				subServerStatusInfo.TaprootAssets.Disabled = v.Disabled
				subServerStatusInfo.TaprootAssets.Running = v.Running
				subServerStatusInfo.TaprootAssets.Error = v.Error
				subServerStatusInfo.TaprootAssets.CustomStatus = v.CustomStatus
			}
		}
	}
	return subServerStatusInfo, nil
}

func GetSubServerStatusInfo() string {
	response, err := getSubServerStatusInfo()
	if err != nil {
		return MakeJsonErrorResult(getSubServerStatusInfoErr, err.Error(), nil)
	}
	return MakeJsonErrorResult(SUCCESS, "", response)
}
