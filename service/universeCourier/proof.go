package universeCourier

import (
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/wire"
	"github.com/lightninglabs/taproot-assets/asset"
	"github.com/lightninglabs/taproot-assets/proof"
	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
	"path/filepath"
	"strings"
)

var defaultProofPath = "/home/en/tapdtest/.tapd/data/regtest/proofs"

func serverDialOpts() ([]grpc.DialOption, error) {
	var opts []grpc.DialOption

	// Skip TLS certificate verification.
	tlsConfig := tls.Config{InsecureSkipVerify: true}
	transportCredentials := credentials.NewTLS(&tlsConfig)
	opts = append(opts, grpc.WithTransportCredentials(transportCredentials))

	return opts, nil
}

// FetchProofs retrieves all proofs for a given asset ID from the proof
func FetchProofs(id asset.ID) ([]*proof.AnnotatedProof, error) {
	assetID := hex.EncodeToString(id[:])
	assetPath := filepath.Join(defaultProofPath, assetID)
	entries, err := os.ReadDir(assetPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read dir %s: %w", assetPath,
			err)
	}

	proofs := make([]*proof.AnnotatedProof, len(entries))
	for idx := range entries {
		// We'll skip any files that don't end with our suffix, this
		// will include directories as well, so we don't need to check
		// for those.
		fileName := entries[idx].Name()
		if !strings.HasSuffix(fileName, proof.TaprootAssetsFileSuffix) {
			continue
		}

		parts := strings.Split(strings.ReplaceAll(
			fileName, proof.TaprootAssetsFileSuffix, "",
		), "-")
		if len(parts) != 3 {
			return nil, fmt.Errorf("malformed proof file name "+
				"'%s', expected two parts, got %d", fileName,
				len(parts))
		}

		fullPath := filepath.Join(assetPath, fileName)
		proofFile, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read proof: %w", err)
		}

		proofs[idx] = &proof.AnnotatedProof{
			Blob: proofFile,
		}
	}

	return proofs, nil
}

// todo: add a function to fetch the proof
func FetchProof(id proof.Locator) (proof.Blob, error) {
	// All our on-disk storage is based on asset IDs, so to look up a path,
	// we just need to compute the full file path and see if it exists on
	// disk.
	proofPath, err := lookupProofFilePath(defaultProofPath, id)
	if err != nil {
		return nil, fmt.Errorf("unable to make proof file path: %w",
			err)
	}

	proofFile, err := os.ReadFile(proofPath)
	switch {
	case os.IsNotExist(err):
		return nil, proof.ErrProofNotFound
	case err != nil:
		return nil, fmt.Errorf("unable to find proof: %w", err)
	}

	return proofFile, nil
}

func NewProofLoc(assetId, groupKey, scriptKey, outpoint string) *proof.Locator {
	var _assetId asset.ID
	assetIdBytes, err := hex.DecodeString(assetId)
	if err != nil {
		return nil
	}
	copy(_assetId[:], assetIdBytes)

	scriptKeyBytes, err := hex.DecodeString(scriptKey)
	if err != nil {
		return nil
	}
	_scriptKey, err := btcec.ParsePubKey(scriptKeyBytes)
	if err != nil {
		return nil
	}

	_outpoint, err := wire.NewOutPointFromString(outpoint)
	if err != nil {
		return nil
	}

	if groupKey != "" {
		groupKeyBytes, err := hex.DecodeString(groupKey)
		if err != nil {
			return nil
		}
		_groupKey, err := btcec.ParsePubKey(groupKeyBytes)
		if err != nil {
			return nil
		}

		return &proof.Locator{
			AssetID:   &_assetId,
			ScriptKey: *_scriptKey,
			OutPoint:  _outpoint,
			GroupKey:  _groupKey,
		}

	}

	return &proof.Locator{
		AssetID:   &_assetId,
		ScriptKey: *_scriptKey,
		OutPoint:  _outpoint,
	}
}

func ImportProofs(proofPath string, replace bool,
	proofs ...*proof.AnnotatedProof) error {
	for _, p := range proofs {
		proofPath, err := genProofFileStoragePath(
			proofPath, p.Locator,
		)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(proofPath), 0750); err != nil {
			return err
		}

		// Can't replace a file that doesn't exist yet.
		if replace && !lnrpc.FileExists(proofPath) {
			return fmt.Errorf("cannot replace p because file "+
				"%s does not exist", proofPath)
		}

		err = os.WriteFile(proofPath, p.Blob, 0666)
		if err != nil {
			return fmt.Errorf("unable to store p: %v", err)
		}
	}
	return nil
}

func lookupProofFilePath(rootPath string, loc proof.Locator) (string, error) {
	// If an outpoint is specified, we want to look up a very specific file
	// on disk.
	if loc.OutPoint != nil {
		fullName, err := genProofFileStoragePath(rootPath, loc)
		if err != nil {
			return "", err
		}

		// If the file doesn't exist under the full name, we know there
		// just isn't a proof file for that asset yet.
		if !lnrpc.FileExists(fullName) {
			return "", fmt.Errorf("proof file %s does not "+
				"exist: %w", fullName, proof.ErrProofNotFound)
		}

		return fullName, nil
	}

	// If the user didn't specify an outpoint, we look up all proof files
	// that start with the script key given. If there is exactly one, we
	// return it.
	var emptyKey btcec.PublicKey
	switch {
	case loc.AssetID == nil:
		return "", proof.ErrInvalidLocatorID

	case loc.ScriptKey.IsEqual(&emptyKey):
		return "", proof.ErrInvalidLocatorKey
	}
	assetID := hex.EncodeToString(loc.AssetID[:])
	scriptKey := hex.EncodeToString(loc.ScriptKey.SerializeCompressed())

	searchPattern := filepath.Join(rootPath, assetID, scriptKey+"*")
	matches, err := filepath.Glob(searchPattern)
	if err != nil {
		return "", fmt.Errorf("error listing proof files: %w", err)
	}

	switch {
	// We have no proof for this script key.
	case len(matches) == 0:
		return "", proof.ErrProofNotFound

	// Exactly one proof for this script key, we'll return it.
	case len(matches) == 1:
		return matches[0], nil

	// User needs to specify the outpoint as well, since we have multiple
	// proofs for this script key.
	default:
		return "", proof.ErrMultipleProofs
	}
}
func genProofFileStoragePath(rootPath string, loc proof.Locator) (string, error) {
	var emptyKey btcec.PublicKey
	switch {
	case loc.AssetID == nil:
		return "", proof.ErrInvalidLocatorID

	case loc.ScriptKey.IsEqual(&emptyKey):
		return "", proof.ErrInvalidLocatorKey

	case loc.OutPoint == nil:

		return "", proof.ErrOutPointMissing
	}

	assetID := hex.EncodeToString(loc.AssetID[:])

	truncatedHash := loc.OutPoint.Hash.String()[:32]
	fileName := fmt.Sprintf("%x-%s-%d.%s",
		loc.ScriptKey.SerializeCompressed(), truncatedHash,
		loc.OutPoint.Index, proof.TaprootAssetsFileEnding)

	return filepath.Join(rootPath, assetID, fileName), nil
}
