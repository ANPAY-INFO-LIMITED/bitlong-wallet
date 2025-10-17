package api

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/pkg/errors"
)

func SignSchnorr(hexPrivateKey string, message string) (string, error) {
	s, err := hex.DecodeString(hexPrivateKey)
	if err != nil {
		return "", errors.Wrap(err, "hex.DecodeString")
	}
	sk, _ := btcec.PrivKeyFromBytes(s)
	h := sha256.Sum256([]byte(message))
	sig, err := schnorr.Sign(sk, h[:])
	if err != nil {
		return "", errors.Wrap(err, "schnorr.Sign")
	}
	return hex.EncodeToString(sig.Serialize()), nil
}
