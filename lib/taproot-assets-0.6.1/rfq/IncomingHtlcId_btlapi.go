//go:build btlapi
// +build btlapi

package rfq

var incomingHtlcIds []uint64

func SaveIncomingHtlcIds(htlcId uint64) {
	incomingHtlcIds = append(incomingHtlcIds, htlcId)
}

func CheckIncomingHtlcId(htlcId uint64) bool {
	for _, id := range incomingHtlcIds {
		if id == htlcId {
			return true
		}
	}
	return false
}
