package models

type SubServers map[string]*SubServerStatus

type SubServerStatus struct {
	Disabled     bool   `json:"disabled"`
	Running      bool   `json:"running"`
	Error        string `json:"error"`
	CustomStatus string `json:"custom_status"`
}
