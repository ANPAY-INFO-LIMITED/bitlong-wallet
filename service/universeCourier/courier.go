package universeCourier

import (
	"bytes"
	"context"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/lightninglabs/taproot-assets/fn"
	"github.com/lightninglabs/taproot-assets/proof"
	"github.com/lightninglabs/taproot-assets/taprpc"
	unirpc "github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	"google.golang.org/grpc"
	"net/url"
)

type courier struct {
	proof.Courier

	// client is the RPC client that the courier will use to interact with
	// the universe RPC server.
	client unirpc.UniverseClient

	// rawConn is the raw connection that the courier will use to interact
	// with the remote gRPC service.
	rawConn *grpc.ClientConn
}

func (c *courier) DeliverProof(ctx context.Context,
	annotatedProof *proof.AnnotatedProof) error {
	// Decode annotated proof into proof file.
	proofFile := &proof.File{}
	err := proofFile.Decode(bytes.NewReader(annotatedProof.Blob))
	if err != nil {
		return err
	}

	// Iterate over each proof in the proof file and submit to the courier
	// service.
	for i := 0; i < proofFile.NumProofs(); i++ {
		transitionProof, err := proofFile.ProofAt(uint32(i))
		if err != nil {
			return err
		}
		proofAsset := transitionProof.Asset

		// Construct asset leaf.
		rpcAsset, err := taprpc.MarshalAsset(
			ctx, &proofAsset, true, true, nil, fn.None[uint32](),
		)
		if err != nil {
			return err
		}

		var proofBuf bytes.Buffer
		if err := transitionProof.Encode(&proofBuf); err != nil {
			return fmt.Errorf("error encoding proof file: %w", err)
		}

		assetLeaf := unirpc.AssetLeaf{
			Asset: rpcAsset,
			Proof: proofBuf.Bytes(),
		}

		// Construct universe key.
		outPoint := transitionProof.OutPoint()
		assetKey := unirpc.MarshalAssetKey(
			outPoint, proofAsset.ScriptKey.PubKey,
		)
		assetID := proofAsset.ID()

		var (
			groupPubKey      *btcec.PublicKey
			groupPubKeyBytes []byte
		)
		if proofAsset.GroupKey != nil {
			groupPubKey = &proofAsset.GroupKey.GroupPubKey
			groupPubKeyBytes = groupPubKey.SerializeCompressed()
		}

		universeID := unirpc.MarshalUniverseID(
			assetID[:], groupPubKeyBytes,
		)
		universeKey := unirpc.UniverseKey{
			Id:      universeID,
			LeafKey: assetKey,
		}
		// Submit proof to courier.
		_, err = c.client.InsertProof(ctx, &unirpc.AssetProof{
			Key:       &universeKey,
			AssetLeaf: &assetLeaf,
		})
		if err != nil {
			return fmt.Errorf("error inserting proof "+
				"into universe courier service: %w",
				err)
		}

	}
	return err
}
func (c *courier) ReceiveProof(ctx context.Context,
	originLocator proof.Locator) (*proof.AnnotatedProof, error) {

	fetchProof := func(ctx context.Context, loc proof.Locator) (proof.Blob, error) {
		var groupKeyBytes []byte
		if loc.GroupKey != nil {
			groupKeyBytes = loc.GroupKey.SerializeCompressed()
		}

		if loc.OutPoint == nil {
			return nil, fmt.Errorf("proof locator for asset %x "+
				"is missing outpoint", loc.AssetID[:])
		}

		universeKey := unirpc.UniverseKey{
			Id: unirpc.MarshalUniverseID(
				loc.AssetID[:], groupKeyBytes,
			),
			LeafKey: unirpc.MarshalAssetKey(
				*loc.OutPoint, &loc.ScriptKey,
			),
		}

		// Setup proof receive/query routine and start backoff
		// procedure.
		var proofBlob []byte
		receiveFunc := func() error {
			// Retrieve proof from courier.
			resp, err := c.client.QueryProof(ctx, &universeKey)
			if err != nil {
				return err
			}
			if err != nil {
				return fmt.Errorf("error retreving proof "+
					"from universe courier service: %w",
					err)
			}

			proofBlob = resp.AssetLeaf.Proof

			return nil
		}
		err := receiveFunc()
		if err != nil {
			return nil, err
		}
		return proofBlob, nil
	}

	proofFile, err := proof.FetchProofProvenance(
		ctx, nil, originLocator, fetchProof,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching proof provenance: %w",
			err)
	}

	// Encode the full proof file.
	var buf bytes.Buffer
	if err := proofFile.Encode(&buf); err != nil {
		return nil, fmt.Errorf("error encoding proof file: %w", err)
	}
	proofFileBlob := buf.Bytes()

	return &proof.AnnotatedProof{
		Locator: originLocator,
		Blob:    proofFileBlob,
	}, nil
}

func (c *courier) Close() error {
	err := c.rawConn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *courier) QueryAssetKey(assetId string) (*unirpc.AssetLeafKeyResponse, error) {
	i := unirpc.ID{
		Id: &unirpc.ID_AssetIdStr{
			AssetIdStr: assetId,
		},
		ProofType: unirpc.ProofType_PROOF_TYPE_TRANSFER,
	}
	keys, err := c.client.AssetLeafKeys(context.Background(), &unirpc.AssetLeafKeysRequest{
		Id: &i,
	})
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func newCourier(addr *url.URL) (*courier, error) {
	switch addr.Scheme {
	case proof.HashmailCourierType:
	case proof.UniverseRpcCourierType:

		// Connect to the universe RPC server.
		dialOpts, err := serverDialOpts()
		if err != nil {
			return nil, err
		}

		serverAddr := fmt.Sprintf("%s:%s", addr.Hostname(), addr.Port())
		conn, err := grpc.Dial(serverAddr, dialOpts...)
		if err != nil {
			return nil, err
		}

		client := unirpc.NewUniverseClient(conn)

		return &courier{
			rawConn: conn,
			client:  client,
		}, nil

	default:
		return nil, fmt.Errorf("unknown courier address protocol "+
			"(consider updating tapd): %v", addr.Scheme)
	}
	return nil, nil
}
