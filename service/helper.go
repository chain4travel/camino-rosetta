package service

import (
	"context"
	"encoding/json"
	"math/big"
	"strings"

	"github.com/chain4travel/camino-rosetta/client"
	"github.com/coinbase/rosetta-sdk-go/types"

	ethtypes "github.com/chain4travel/caminoethvm/core/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	nativeTransferGasLimit = uint64(21000)
	erc20TransferGasLimit  = uint64(250000)
	genesisTimestamp       = 946713601000 // min allowable timestamp
)

func makeGenesisBlock(hash string) *types.Block {
	return &types.Block{
		BlockIdentifier: &types.BlockIdentifier{
			Index: 0,
			Hash:  hash,
		},
		ParentBlockIdentifier: &types.BlockIdentifier{
			Index: 0,
			Hash:  hash,
		},
		Timestamp: genesisTimestamp,
	}
}

func blockHeaderFromInput(
	ctx context.Context,
	c client.Client,
	input *types.PartialBlockIdentifier,
) (*ethtypes.Header, *types.Error) {
	var (
		header *ethtypes.Header
		err    error
	)

	if input == nil {
		header, err = c.HeaderByNumber(ctx, nil)
	} else {
		if input.Hash == nil && input.Index == nil {
			return nil, errInvalidInput
		}

		if input.Index != nil {
			header, err = c.HeaderByNumber(ctx, big.NewInt(*input.Index))
		} else {
			header, err = c.HeaderByHash(ctx, ethcommon.HexToHash(*input.Hash))
		}
	}

	if err != nil {
		return nil, wrapError(errInternalError, err)
	}

	return header, nil
}

// unmarshalJSONMap converts map[string]interface{} into a interface{}.
func unmarshalJSONMap(m map[string]interface{}, i interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, i)
}

// marshalJSONMap converts an interface into a map[string]interface{}.
func marshalJSONMap(i interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return m, nil
}

// ChecksumAddress ensures an Ethereum hex address
// is in Checksum Format. If the address cannot be converted,
// it returns !ok.
func ChecksumAddress(address string) (string, bool) {
	if !strings.HasPrefix(address, "0x") {
		return "", false
	}

	addr, err := ethcommon.NewMixedcaseAddressFromString(address)
	if err != nil {
		return "", false
	}

	return addr.Address().Hex(), true
}
