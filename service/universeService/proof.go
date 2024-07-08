package universeService

import (
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/lightninglabs/taproot-assets/asset"
	"github.com/lightninglabs/taproot-assets/proof"
	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
	"path/filepath"
	"strings"
)

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
	proofPath := "/home/en/tapdtest/.tapd/data/regtest/proofs"
	assetID := hex.EncodeToString(id[:])
	assetPath := filepath.Join(proofPath, assetID)
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
