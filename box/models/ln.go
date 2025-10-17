package models

type WalletBalanceResponse struct {
	TotalBalance              int64          `json:"total_balance"`
	ConfirmedBalance          int64          `json:"confirmed_balance"`
	UnconfirmedBalance        int64          `json:"unconfirmed_balance"`
	LockedBalance             int64          `json:"locked_balance"`
	ReservedBalanceAnchorChan int64          `json:"reserved_balance_anchor_chan"`
	AccountBalance            AccountBalance `json:"account_balance"`
}

type AccountBalance struct {
	Default Default `json:"default"`
}

type Default struct {
	ConfirmedBalance   int64 `json:"confirmed_balance"`
	UnconfirmedBalance int64 `json:"unconfirmed_balance"`
}
