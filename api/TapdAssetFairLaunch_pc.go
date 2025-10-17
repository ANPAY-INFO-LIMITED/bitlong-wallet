package api

import (
	"github.com/pkg/errors"
)

func PcGetIssuanceTransactionFee(token string, feeRate int) (int, error) {
	result, err := RequestToGetIssuanceTransactionFee(token, feeRate)
	if err != nil {
		return 0, errors.Wrap(err, "RequestToGetIssuanceTransactionFee")
	}
	return result, nil
}
