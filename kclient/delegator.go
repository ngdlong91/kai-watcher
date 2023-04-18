package kclient

import (
	"context"

	"go.uber.org/zap"

	"github.com/ngdlong91/kai-watcher/types"
)

func (n *Node) UnbondedRecords(ctx context.Context, validatorSMCAddress, delegatorAddress string, blockNumber uint64) (*types.UnbondedRecord, error) {
	lgr := n.lgr.With(zap.String("method", "UnbondedRecords"))
	payload, err := n.validatorSMC.Abi.Pack("getUBDEntries", delegatorAddress)
	if err != nil {
		lgr.Error("Error packing UDB entry payload: ", zap.Error(err))
		return nil, err
	}
	res, err := n.CallAt(ctx, ConstructCallArgs(validatorSMCAddress, payload), blockNumber)
	if err != nil {
		lgr.Error("GetUDBEntry KardiaCall error: ", zap.Error(err))
		return nil, err
	}
	if len(res) == 0 {
		return nil, ErrEmptyList
	}

	var result *types.UnbondedRecord
	// unpack result
	err = n.validatorSMC.Abi.UnpackIntoInterface(&result, "getUBDEntries", res)
	if err != nil {
		lgr.Error("Error unpacking UDB entry: ", zap.Error(err))
		return nil, err
	}

	return result, nil
}
