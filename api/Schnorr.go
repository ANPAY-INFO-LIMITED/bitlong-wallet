package api

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
)

func SignSchnorr(hexPrivateKey string, message string) (string, error) {
	s, err := hex.DecodeString(hexPrivateKey)
	if err != nil {
		return "", fmt.Errorf("Sign called with invalid private key '%s': %w", hexPrivateKey, err)
	}
	sk, _ := btcec.PrivKeyFromBytes(s)
	h := sha256.Sum256([]byte(message))
	sig, err := schnorr.Sign(sk, h[:])
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(sig.Serialize()), nil
}
