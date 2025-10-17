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

	SetPathErr
	CreateWalletErr
	RestoreWalletErr
	UnlockWalletErr
	GetStateErr
	SubServersStatusErr
	LndGetInfoErr
	LitdStopDaemonErr
	LndStopDaemonErr

	GenerateKeysErr
	GetPrivateKeyErr
	GetNPublicKeyErr
	GetPublicKeyErr
	GetNBPublicKeyErr

	LoginErr

	UploadWalletBalanceErr
	UploadAssetManagedUtxosErr
	UploadAssetLocalMintHistoryErr
	UploadAssetListInfoErr
	UploadAddrReceivesErr
	UploadAssetTransferErr
	UploadAssetBalanceInfoErr
	UploadAssetBalanceHistoriesErr
	UploadBtcListUnspentUtxosErr
	AutoMintReservedErr
	MergeUTXOErr
	GetWalletBalanceTotalValueErr
	SyncUniverseErr

	GetWalletBalanceErr
	GetBtcTransferInInfosJsonResultErr
	GetBtcTransferOutInfosJsonResultErr
	BtcUtxosErr
	GetNewAddressErr
	SendCoinsErr

	ListNormalBalancesErr
	CheckAssetIssuanceIsLocalErr
	AddrReceivesErr
	QueryAssetTransfersErr
	AssetUtxosErr
	NewAddrErr
	QueryAddrsErr
	SendAssetsErr
	ListNftGroupsErr
	ListNonGroupNftAssetsErr
	GetSpentNftAssetsErr
	MintAssetErr
	AddGroupAssetErr
	FinalizeBatchErr
	CancelBatchErr
	GetIssuanceTransactionFeeErr
	GetAssetInfoErr

	UserHomeDirErr
	CreateFileErr
	InvalidNetwork
	InvalidReq

	GetNewAddressP2trErr
	GetNewAddressP2wkhErr
	GetNewAddressNp2wkhErr
	GetNewAddressP2trExampleErr
	GetNewAddressP2wkhExampleErr
	GetNewAddressNp2wkhExampleErr
	StoreAddrErr
	RemoveAddrErr
	QueryAllAddrErr
	UpdateAllAddressesByGnzbaErr

	GetAddressTransactionsByMempoolErr
	GetTransactionByMempoolErr

	SignSchnorrErr

	GetListEligibleCoinsErr
	CreateSellOrderSignErr
	BuySOrderSignErr
	PublishSOrderTxErr
	GetLastProofErr
	AllowFederationSyncInsertAndExportErr
	InsertProofAndRegisterTransferErr

	LnurlGetAvailPortErr
	LnurlRunFrpcConfErr
	LnurlRunFrpcErr
	LnurlRequestErr
	LnurlRequestInvoiceErr

	AddSessionErr
	NewSessionErr
	ListSessionsErr
	RevokeAllSessionsErr

	ProxyIoReadAllErr
	ProxyHttpNewRequestErr
	ProxyInsecureClientDoErr

	EnableRemoteErr

	CheckFrpStatusErr
	GetFrpPidErr
	KillNineErr
)
