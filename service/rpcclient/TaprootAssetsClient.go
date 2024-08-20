package rpcclient

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/wire"
	"github.com/lightninglabs/taproot-assets/asset"
	"github.com/lightninglabs/taproot-assets/commitment"
	"github.com/lightninglabs/taproot-assets/fn"
	"github.com/lightninglabs/taproot-assets/proof"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/wallet/base"
	"github.com/wallet/service/apiConnect"
)

const (
	mainnetProofCourierAddr = "universerpc://132.232.109.84:8444"
	testnetProofCourierAddr = "universerpc://testnet.universe.lightning.finance:10029"
	regtestProofCourierAddr = "universerpc://132.232.109.84:8443"
)

func getTaprootAssetsClient() (taprpc.TaprootAssetsClient, func(), error) {
	conn, clearUp, err := apiConnect.GetConnection("tapd", false)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	client := taprpc.NewTaprootAssetsClient(conn)
	return client, clearUp, nil
}

func AddrReceives() (*taprpc.AddrReceivesResponse, error) {
	client, clearUp, err := getTaprootAssetsClient()
	if err != nil {
		return nil, err
	}
	defer clearUp()

	request := &taprpc.AddrReceivesRequest{}
	response, err := client.AddrReceives(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func BurnAsset(AssetIdStr string, amountToBurn uint64) (*taprpc.BurnAssetResponse, error) {
	client, clearUp, err := getTaprootAssetsClient()
	if err != nil {
		return nil, err
	}
	defer clearUp()

	request := &taprpc.BurnAssetRequest{
		Asset: &taprpc.BurnAssetRequest_AssetIdStr{
			AssetIdStr: AssetIdStr,
		},
		AmountToBurn:     amountToBurn,
		ConfirmationText: "assets will be destroyed",
	}
	response, err := client.BurnAsset(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func DecodeAddr(addr string) (*taprpc.Addr, error) {

	client, clearUp, err := getTaprootAssetsClient()
	if err != nil {
		return nil, err
	}
	defer clearUp()

	request := &taprpc.DecodeAddrRequest{
		Addr: addr,
	}
	response, err := client.DecodeAddr(context.Background(), request)
	if err != nil {
		fmt.Printf("%s taprpc DecodeAddr Error: %v\n", GetTimeNow(), err)
		return nil, err
	}
	return response, nil
}

func QueryAddr() (*taprpc.QueryAddrResponse, error) {
	client, clearUp, err := getTaprootAssetsClient()
	if err != nil {
		return nil, err
	}
	defer clearUp()
	request := &taprpc.QueryAddrRequest{}
	response, err := client.QueryAddrs(context.Background(), request)
	if err != nil {
		fmt.Printf("%s taprpc QueryAddr Error: %v\n", GetTimeNow(), err)
		return nil, err
	}
	return response, nil
}

func NewAddr(assetId string, amt int) (*taprpc.Addr, error) {
	client, clearUp, err := getTaprootAssetsClient()
	if err != nil {
		return nil, err
	}
	defer clearUp()
	var ProofCourierAddr string
	switch base.NetWork {
	case "mainnet":
		ProofCourierAddr = mainnetProofCourierAddr
	case "testnet":
		ProofCourierAddr = testnetProofCourierAddr
	case "regtest":
		ProofCourierAddr = regtestProofCourierAddr
	default:
		return nil, fmt.Errorf("invalid network: %s", base.NetWork)

	}

	_assetIdByteSlice, _ := hex.DecodeString(assetId)
	request := &taprpc.NewAddrRequest{
		AssetId:          _assetIdByteSlice,
		Amt:              uint64(amt),
		ProofCourierAddr: ProofCourierAddr,
	}
	response, err := client.NewAddr(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func ListTransfers() (*taprpc.ListTransfersResponse, error) {
	client, clearUp, err := getTaprootAssetsClient()
	if err != nil {
		return nil, err
	}
	defer clearUp()
	request := &taprpc.ListTransfersRequest{}
	response, err := client.ListTransfers(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return response, err
}

func ListGroups() (*taprpc.ListGroupsResponse, error) {
	client, clearUp, err := getTaprootAssetsClient()
	if err != nil {
		return nil, err
	}
	defer clearUp()

	request := &taprpc.ListGroupsRequest{}
	response, err := client.ListGroups(context.Background(), request)
	if err != nil {
		fmt.Printf("%s taprpc ListGroups Error: %v\n", GetTimeNow(), err)
		return nil, err
	}
	return response, nil
}

func DecodeProof(proof []byte, depth uint32, withMetaReveal bool, withPrevWitnesses bool) (*taprpc.DecodeProofResponse, error) {
	request := &taprpc.DecodeProofRequest{
		RawProof:          proof,
		ProofAtDepth:      depth,
		WithMetaReveal:    withMetaReveal,
		WithPrevWitnesses: withPrevWitnesses,
	}
	if withMetaReveal || withPrevWitnesses {
		conn, clearUp, err := apiConnect.GetConnection("tapd", false)
		if err != nil {
			fmt.Printf("%s did not connect: %v\n", GetTimeNow(), err)
		}
		defer clearUp()
		client := taprpc.NewTaprootAssetsClient(conn)
		response, err := client.DecodeProof(context.Background(), request)
		return response, err
	} else {
		o := newDecodeProofOffline()
		response, err := o.decodeProof(context.Background(), request)
		return response, err
	}
}

type decodeProofOffline struct {
	//withPrevWitnesses and withMetaReveal need an online node
	withPrevWitnesses bool
	withMetaReveal    bool
}

func newDecodeProofOffline() *decodeProofOffline {
	return &decodeProofOffline{
		withPrevWitnesses: false,
		withMetaReveal:    false,
	}
}

func (d *decodeProofOffline) decodeProof(ctx context.Context,
	req *taprpc.DecodeProofRequest) (*taprpc.DecodeProofResponse, error) {

	if req.WithPrevWitnesses || req.WithMetaReveal {
		return nil, fmt.Errorf("unable to marshal proof: cannot set WithPrevWitnesses" +
			"WithMetaReveal when decoding offline")
	}

	var (
		proofReader = bytes.NewReader(req.RawProof)
		rpcProof    *taprpc.DecodedProof
	)
	switch {
	case proof.IsSingleProof(req.RawProof):
		var p proof.Proof
		err := p.Decode(proofReader)
		if err != nil {
			return nil, fmt.Errorf("unable to decode proof: %w",
				err)
		}

		rpcProof, err = d.marshalProof(
			ctx, &p, d.withMetaReveal, d.withPrevWitnesses,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal proof: %w",
				err)
		}

		rpcProof.NumberOfProofs = 1

	case proof.IsProofFile(req.RawProof):
		if err := proof.CheckMaxFileSize(req.RawProof); err != nil {
			return nil, fmt.Errorf("invalid proof file: %w", err)
		}

		var proofFile proof.File
		if err := proofFile.Decode(proofReader); err != nil {
			return nil, fmt.Errorf("unable to decode proof file: "+
				"%w", err)
		}

		latestProofIndex := uint32(proofFile.NumProofs() - 1)
		if req.ProofAtDepth > latestProofIndex {
			return nil, fmt.Errorf("invalid depth %d is greater "+
				"than latest proof index of %d",
				req.ProofAtDepth, latestProofIndex)
		}

		// Default to latest proof.
		index := latestProofIndex - req.ProofAtDepth
		p, err := proofFile.ProofAt(index)
		if err != nil {
			return nil, err
		}

		rpcProof, err = d.marshalProof(
			ctx, p, req.WithPrevWitnesses,
			req.WithMetaReveal,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal proof: %w",
				err)
		}

		rpcProof.ProofAtDepth = req.ProofAtDepth
		rpcProof.NumberOfProofs = uint32(proofFile.NumProofs())

	default:
		return nil, fmt.Errorf("invalid raw proof, could not " +
			"identify decoding format")
	}

	return &taprpc.DecodeProofResponse{
		DecodedProof: rpcProof,
	}, nil
}

func (d *decodeProofOffline) marshalProof(ctx context.Context, p *proof.Proof,
	withPrevWitnesses, withMetaReveal bool) (*taprpc.DecodedProof, error) {

	var (
		rpcMeta        *taprpc.AssetMeta
		rpcGenesis     = p.GenesisReveal
		rpcGroupKey    = p.GroupKeyReveal
		anchorOutpoint = wire.OutPoint{
			Hash:  p.AnchorTx.TxHash(),
			Index: p.InclusionProof.OutputIndex,
		}
		txMerkleProof  = p.TxMerkleProof
		inclusionProof = p.InclusionProof
		splitRootProof = p.SplitRootProof
	)

	var txMerkleProofBuf bytes.Buffer
	if err := txMerkleProof.Encode(&txMerkleProofBuf); err != nil {
		return nil, fmt.Errorf("unable to encode serialized Bitcoin "+
			"merkle proof: %w", err)
	}

	var inclusionProofBuf bytes.Buffer
	if err := inclusionProof.Encode(&inclusionProofBuf); err != nil {
		return nil, fmt.Errorf("unable to encode inclusion proof: %w",
			err)
	}

	if inclusionProof.CommitmentProof == nil {
		return nil, fmt.Errorf("inclusion proof is missing " +
			"commitment proof")
	}
	tsSibling, tsHash, err := commitment.MaybeEncodeTapscriptPreimage(
		inclusionProof.CommitmentProof.TapSiblingPreimage,
	)
	if err != nil {
		return nil, fmt.Errorf("error encoding tapscript sibling: %w",
			err)
	}

	tapProof, err := inclusionProof.CommitmentProof.DeriveByAssetInclusion(
		&p.Asset,
	)
	if err != nil {
		return nil, fmt.Errorf("error deriving inclusion proof: %w",
			err)
	}
	merkleRoot := tapProof.TapscriptRoot(tsHash)

	var exclusionProofs [][]byte
	for _, exclusionProof := range p.ExclusionProofs {
		var exclusionProofBuf bytes.Buffer
		err := exclusionProof.Encode(&exclusionProofBuf)
		if err != nil {
			return nil, fmt.Errorf("unable to encode exclusion "+
				"proofs: %w", err)
		}
		exclusionProofs = append(
			exclusionProofs, exclusionProofBuf.Bytes(),
		)
	}

	var splitRootProofBuf bytes.Buffer
	if splitRootProof != nil {
		err := splitRootProof.Encode(&splitRootProofBuf)
		if err != nil {
			return nil, fmt.Errorf("unable to encode split root "+
				"proof: %w", err)
		}
	}

	rpcAsset, err := d.marshalChainAsset(ctx, &asset.ChainAsset{
		Asset:                  &p.Asset,
		AnchorTx:               &p.AnchorTx,
		AnchorBlockHash:        p.BlockHeader.BlockHash(),
		AnchorBlockHeight:      p.BlockHeight,
		AnchorOutpoint:         anchorOutpoint,
		AnchorInternalKey:      p.InclusionProof.InternalKey,
		AnchorMerkleRoot:       merkleRoot[:],
		AnchorTapscriptSibling: tsSibling,
	}, withPrevWitnesses)
	if err != nil {
		return nil, err
	}

	if withMetaReveal {
		//metaHash := rpcAsset.AssetGenesis.MetaHash
		//if len(metaHash) == 0 {
		//	return nil, fmt.Errorf("asset does not contain meta " +
		//		"data")
		//}
		//
		//rpcMeta, err = r.FetchAssetMeta(
		//	ctx, &taprpc.FetchAssetMetaRequest{
		//		Asset: &taprpc.FetchAssetMetaRequest_MetaHash{
		//			MetaHash: metaHash,
		//		},
		//	},
		//)
		//if err != nil {
		//	return nil, err
		//}
	}

	decodedAssetID := p.Asset.ID()
	var genesisReveal *taprpc.GenesisReveal
	if rpcGenesis != nil {
		genesisReveal = &taprpc.GenesisReveal{
			GenesisBaseReveal: &taprpc.GenesisInfo{
				GenesisPoint: rpcGenesis.FirstPrevOut.String(),
				Name:         rpcGenesis.Tag,
				MetaHash:     rpcGenesis.MetaHash[:],
				AssetId:      decodedAssetID[:],
				OutputIndex:  rpcGenesis.OutputIndex,
				AssetType:    taprpc.AssetType(p.Asset.Type),
			},
		}
	}

	var GroupKeyReveal taprpc.GroupKeyReveal
	if rpcGroupKey != nil {
		GroupKeyReveal = taprpc.GroupKeyReveal{
			RawGroupKey:   rpcGroupKey.RawKey[:],
			TapscriptRoot: rpcGroupKey.TapscriptRoot,
		}
	}

	return &taprpc.DecodedProof{
		Asset:               rpcAsset,
		MetaReveal:          rpcMeta,
		TxMerkleProof:       txMerkleProofBuf.Bytes(),
		InclusionProof:      inclusionProofBuf.Bytes(),
		ExclusionProofs:     exclusionProofs,
		SplitRootProof:      splitRootProofBuf.Bytes(),
		NumAdditionalInputs: uint32(len(p.AdditionalInputs)),
		ChallengeWitness:    p.ChallengeWitness,
		IsBurn:              p.Asset.IsBurn(),
		GenesisReveal:       genesisReveal,
		GroupKeyReveal:      &GroupKeyReveal,
	}, nil
}

func (d *decodeProofOffline) marshalChainAsset(ctx context.Context, a *asset.ChainAsset,
	withWitness bool) (*taprpc.Asset, error) {

	rpcAsset, err := taprpc.MarshalAsset(
		ctx, a.Asset, a.IsSpent, withWitness, nil, fn.None[uint32](),
	)
	if err != nil {
		return nil, err
	}

	var anchorTxBytes []byte
	if a.AnchorTx != nil {
		var anchorTxBuf bytes.Buffer
		err := a.AnchorTx.Serialize(&anchorTxBuf)
		if err != nil {
			return nil, fmt.Errorf("unable to serialize anchor "+
				"tx: %w", err)
		}
		anchorTxBytes = anchorTxBuf.Bytes()
	}

	rpcAsset.ChainAnchor = &taprpc.AnchorInfo{
		AnchorTx:         anchorTxBytes,
		AnchorBlockHash:  a.AnchorBlockHash.String(),
		AnchorOutpoint:   a.AnchorOutpoint.String(),
		InternalKey:      a.AnchorInternalKey.SerializeCompressed(),
		MerkleRoot:       a.AnchorMerkleRoot,
		TapscriptSibling: a.AnchorTapscriptSibling,
		BlockHeight:      a.AnchorBlockHeight,
	}

	if a.AnchorLeaseOwner != [32]byte{} {
		rpcAsset.LeaseOwner = a.AnchorLeaseOwner[:]
		rpcAsset.LeaseExpiry = a.AnchorLeaseExpiry.UTC().Unix()
	}

	return rpcAsset, nil
}
