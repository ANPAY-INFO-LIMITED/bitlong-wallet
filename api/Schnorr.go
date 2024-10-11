package api

import (
	"encoding/hex"
	"github.com/decred/dcrd/crypto/blake256"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/decred/dcrd/dcrec/secp256k1/v4/schnorr"
	_ "golang.org/x/crypto/sha3"
)

func Sign(privKey string, msg string) (sigHex string, err error) {
	// Decode a hex-encoded private key
	pkBytes, err := hex.DecodeString(privKey)
	if err != nil {
		return "", err
	}
	privateKey := secp256k1.PrivKeyFromBytes(pkBytes)
	// Sign a msg using the private key
	messageHash := blake256.Sum256([]byte(msg))
	signature, err := schnorr.Sign(privateKey, messageHash[:])
	if err != nil {
		return "", err
	}
	sigBytes := signature.Serialize()
	sigHex = hex.EncodeToString(sigBytes)
	return sigHex, nil
}

func Verify(pubKey string, msg string, sigHex string) (verified bool, err error) {
	// Decode hex-encoded serialized public key
	pubKeyBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		return false, err
	}
	publicKey, err := schnorr.ParsePubKey(pubKeyBytes)
	if err != nil {
		return false, err
	}

	// Decode hex-encoded serialized signature
	sigBytes, err := hex.DecodeString(sigHex)
	if err != nil {
		return false, err
	}
	signature, err := schnorr.ParseSignature(sigBytes)
	if err != nil {
		return false, err
	}
	// Verify the signature for the message using the public key
	messageHash := blake256.Sum256([]byte(msg))
	verified = signature.Verify(messageHash[:], publicKey)
	return verified, nil
}

// SchnorrSign
// @Description: Schnorr sign
// @dev: privKey is 32 Bytes, string of length 64
// @dev: sig is 64 Bytes, string of length 128
func SchnorrSign(privKey string, msg string) string {
	sigHex, err := Sign(privKey, msg)
	if err != nil {
		return MakeJsonErrorResult(SchnorrSignErr, err.Error(), "")
	}
	return MakeJsonErrorResult(SUCCESS, "", sigHex)
}

// SchnorrVerify
// @Description: Schnorr verify
// @dev: pubKey is 33 Bytes, string of length 66
func SchnorrVerify(pubKey string, msg string, sigHex string) string {
	verified, err := Verify(pubKey, msg, sigHex)
	if err != nil {
		return MakeJsonErrorResult(SchnorrVerifyErr, err.Error(), false)
	}
	return MakeJsonErrorResult(SUCCESS, "", verified)
}
