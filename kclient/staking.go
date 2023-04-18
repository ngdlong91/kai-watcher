package kclient

import (
	"context"
	"math/big"

	"github.com/kardiachain/go-kardia/lib/common"
	"go.uber.org/zap"
)

func (n *Node) TotalStakedAmount(ctx context.Context, blockNumber uint64) (*big.Int, error) {
	payload, err := n.stakingSMC.Abi.Pack("totalBonded")
	if err != nil {
		n.lgr.Error("Error packing UDB entry payload: ", zap.Error(err))
		return nil, err
	}
	res, err := n.CallAt(ctx, ConstructCallArgs(n.stakingSMC.ContractAddress.Hex(), payload), blockNumber)
	if err != nil {
		n.lgr.Error("Get totalBonded KardiaCall error: ", zap.Error(err))
		return nil, err
	}

	var result struct {
		TotalBonded *big.Int
	}
	// unpack result
	err = n.stakingSMC.Abi.UnpackIntoInterface(&result, "totalBonded", res)
	if err != nil {
		n.lgr.Error("Error unpacking UDB entry: ", zap.Error(err))
		return nil, err
	}
	return result.TotalBonded, nil
}

func (n *Node) ValidatorSMCAddresses(ctx context.Context, blockNumber uint64) ([]common.Address, error) {
	payload, err := n.stakingSMC.Abi.Pack("getAllValidator")
	if err != nil {
		return nil, err
	}
	var res []byte
	if blockNumber != 0 {
		res, err = n.CallAt(ctx, ConstructCallArgs(n.stakingSMC.ContractAddress.Hex(), payload), blockNumber)
		if err != nil {
			return nil, err
		}
	} else {
		res, err = n.Call(ctx, ConstructCallArgs(n.stakingSMC.ContractAddress.Hex(), payload))
		if err != nil {
			return nil, err
		}
	}

	if len(res) == 0 {
		return nil, ErrEmptyList
	}

	var validatorSMCAddresses []common.Address
	// unpack result
	err = n.stakingSMC.Abi.UnpackIntoInterface(&validatorSMCAddresses, "getAllValidator", res)
	if err != nil {
		return nil, err
	}

	return validatorSMCAddresses, nil
}
