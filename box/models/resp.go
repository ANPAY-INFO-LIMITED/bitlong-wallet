package models

type Resp struct {
	Code Code   `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type RespT[T any] struct {
	Code Code   `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

type (
	RespStr RespT[string]
	RespInt RespT[int]
)

type RespLnc[T any] struct {
	Code Code    `json:"code"`
	Msg  string  `json:"msg"`
	Data LncT[T] `json:"data"`
}

const (
	NullStr = ""
)

type Lnc struct {
	List  any   `json:"list"`
	Count int64 `json:"count"`
}

type LncT[T any] struct {
	List  T     `json:"list"`
	Count int64 `json:"count"`
}

type Code int

const (
	Success Code = 0
)

const (
	ShouldBindJSONErr = iota + 1000

	IpNotAllowed

	UserHomeDirErr
	CreateFileErr
	ReadFileErr

	ProxyIoReadAllErr
	ProxyHttpNewRequestErr
	ProxyInsecureClientDoErr

	ChannelBalanceErr
	BoxListPaymentsRecordsErr
	BoxListInvoicesRecordsErr
	BoxBtcAddInvoicesErr
	BoxBtcPayInvoiceErr
	BoxAddAssetInvoiceErr
	BoxAssetChannelSendPaymentErr
	BoxListChannelsErr
	BoxBtcDecodePayReqErr
	BoxAssetDecodePayReqErr
	BoxGetInfoErr

	HexDecodeStringErr
	NewAddrErr
	SendAssetErr
	DecodeAddrErr

	BtcTransferInErr
	BtcTransferOutErr
	BtcUtxoErr
	AssetTransferInErr
	AssetTransferOutErr
	AssetUtxoErr
	GetWalletBalanceResponseErr

	CheckPasswordErr
	UpdateTokenErr
)

type ErrResp struct {
	Error string `json:"error"`
}

type JResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Data    any    `json:"data"`
}

type JResult2 struct {
	Errnos int         `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}
