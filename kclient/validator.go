package kclient

import (
	"context"
	"math/big"
	"strings"

	"github.com/kardiachain/go-kardia/lib/common"
	"go.uber.org/zap"

	"github.com/ngdlong91/kai-watcher/types"
)

func (n *Node) Validators(ctx context.Context) ([]*types.Validator, error) {
	var (
		validators []*types.Validator
	)

	validatorSMCAddresses, err := n.ValidatorSMCAddresses(ctx, 0)
	if err != nil {
		return nil, err
	}

	for _, smcAddr := range validatorSMCAddresses {
		v, err := n.Validator(ctx, smcAddr.Hex())
		if err != nil {
			return nil, err
		}
		validators = append(validators, v)
	}
	return validators, nil
}

func (n *Node) APIValidatorInfo(ctx context.Context, validatorSMCAddress string) (*types.Validator, error) {
	lgr := n.lgr.With(zap.String("method", "Validator"))
	payload, err := n.validatorSMC.Abi.Pack("inforValidator")
	if err != nil {
		lgr.Error("Error packing validator info payload: ", zap.Error(err))
		return nil, err
	}
	res, err := n.Call(ctx, ConstructCallArgs(validatorSMCAddress, payload))
	if err != nil {
		lgr.Error("GetValidatorInfo KardiaCall error: ", zap.Error(err))
		return nil, err
	}
	var valInfo RPCValidator
	// unpack result
	err = n.validatorSMC.Abi.UnpackIntoInterface(&valInfo, "inforValidator", res)
	if err != nil {
		lgr.Error("Error unpacking validator info: ", zap.Error(err))
		return nil, err
	}
	validator := types.Validator{
		Name:                  validatorName(valInfo.Name),
		Signer:                strings.ToLower(valInfo.Signer.String()),
		SMCAddress:            strings.ToLower(validatorSMCAddress),
		Tokens:                valInfo.Tokens,
		Jailed:                valInfo.Jailed,
		DelegationShares:      valInfo.DelegationShares,
		AccumulatedCommission: valInfo.AccumulatedCommission,
		UbdEntryCount:         valInfo.UbdEntryCount,
		UpdateTime:            valInfo.UpdateTime,
		MinSelfDelegation:     valInfo.MinSelfDelegation,
		Status:                valInfo.Status,
		UnbondingTime:         valInfo.UnbondingTime,
		UnbondingHeight:       valInfo.UnbondingHeight,
	}

	return &validator, nil
}

func (n *Node) ValidatorInfo(ctx context.Context, validatorSMCAddress string, blockNumber uint64) (*types.Validator, error) {
	lgr := n.lgr.With(zap.String("method", "Validator"))
	payload, err := n.validatorSMC.Abi.Pack("inforValidator")
	if err != nil {
		lgr.Error("Error packing validator info payload: ", zap.Error(err))
		return nil, err
	}
	res, err := n.CallAt(ctx, ConstructCallArgs(validatorSMCAddress, payload), blockNumber)
	if err != nil {
		lgr.Error("GetValidatorInfo KardiaCall error: ", zap.Error(err))
		return nil, err
	}
	var valInfo RPCValidator
	// unpack result
	err = n.validatorSMC.Abi.UnpackIntoInterface(&valInfo, "inforValidator", res)
	if err != nil {
		lgr.Error("Error unpacking validator info: ", zap.Error(err))
		return nil, err
	}
	commission, err := n.Commission(ctx, validatorSMCAddress, blockNumber)
	if err != nil {
		return nil, err
	}

	signingInfo, err := n.SigningInfo(ctx, validatorSMCAddress, blockNumber)
	if err != nil {
		return nil, err
	}

	validator := types.Validator{
		Name:                  validatorName(valInfo.Name),
		Signer:                strings.ToLower(valInfo.Signer.String()),
		SMCAddress:            strings.ToLower(validatorSMCAddress),
		Tokens:                valInfo.Tokens,
		Jailed:                valInfo.Jailed,
		DelegationShares:      valInfo.DelegationShares,
		AccumulatedCommission: valInfo.AccumulatedCommission,
		UbdEntryCount:         valInfo.UbdEntryCount,
		UpdateTime:            valInfo.UpdateTime,
		MinSelfDelegation:     valInfo.MinSelfDelegation,
		Status:                valInfo.Status,
		UnbondingTime:         valInfo.UnbondingTime,
		UnbondingHeight:       valInfo.UnbondingHeight,
		CommissionRate:        commission.Rate,
		MaxChangeRate:         commission.MaxChangeRate,
		MaxRate:               commission.MaxRate,
		SigningInfo: &types.SigningInfo{
			StartHeight:        signingInfo.StartHeight.Uint64(),
			IndexOffset:        signingInfo.IndexOffset.Uint64(),
			Tombstoned:         signingInfo.Tombstoned,
			MissedBlockCounter: signingInfo.MissedBlockCounter.Uint64(),
			IndicatorRate:      signingInfo.IndicatorRate,
			JailedUntil:        signingInfo.JailedUntil.Uint64(),
		},
	}
	return &validator, nil
}

func (n *Node) Commission(ctx context.Context, valSmcAddr string, blockNumber uint64) (*Commission, error) {
	payload, err := n.validatorSMC.Abi.Pack("commission")
	if err != nil {
		n.lgr.Debug("cannot pack data", zap.Error(err))
		return nil, err
	}
	res, err := n.CallAt(ctx, ConstructCallArgs(valSmcAddr, payload), blockNumber)
	if err != nil {
		n.lgr.Debug("cannot make request", zap.Error(err))
		return nil, err
	}

	var commission Commission
	// unpack result
	err = n.validatorSMC.Abi.UnpackIntoInterface(&commission, "commission", res)
	if err != nil {
		n.lgr.Debug("cannot unpack data", zap.Error(err))
		return nil, err
	}
	return &commission, nil
}

func (n *Node) SigningInfo(ctx context.Context, validatorSMCAddress string, blockNumber uint64) (*SigningInfo, error) {
	lgr := n.lgr.With(zap.String("method", "getSigningInfo"))
	payload, err := n.validatorSMC.Abi.Pack("signingInfo")
	if err != nil {
		lgr.Error("Error packing get signingInfo payload: ", zap.Error(err))
		return nil, err
	}
	res, err := n.CallAt(ctx, ConstructCallArgs(validatorSMCAddress, payload), blockNumber)
	if err != nil {
		lgr.Error("GetSigningInfo KardiaCall error: ", zap.Error(err))
		return nil, err
	}
	var result SigningInfo
	// unpack result
	err = n.validatorSMC.Abi.UnpackIntoInterface(&result, "signingInfo", res)
	if err != nil {
		lgr.Error("Error unpack get signingInfo: ", zap.Error(err))
		return nil, err
	}
	return &result, nil
}

func (n *Node) DelegationRewards(ctx context.Context, validatorSMCAddr, delegatorAddress string, blockNumber uint64) (*big.Int, error) {
	payload, err := n.validatorSMC.Abi.Pack("getDelegationRewards", common.HexToAddress(delegatorAddress))
	if err != nil {
		n.lgr.Error("Error packing delegation rewards payload: ", zap.Error(err))
		return nil, err
	}
	res, err := n.CallAt(ctx, ConstructCallArgs(validatorSMCAddr, payload), blockNumber)
	if err != nil {
		n.lgr.Error("GetDelegationRewards KardiaCall error: ", zap.Error(err))
		return nil, err
	}
	var result struct {
		Rewards *big.Int
	}
	// unpack result
	err = n.validatorSMC.Abi.UnpackIntoInterface(&result, "getDelegationRewards", res)
	if err != nil {
		n.lgr.Error("Error unpacking delegation rewards: ", zap.Error(err))
		return nil, err
	}
	return result.Rewards, nil
}

func (n *Node) ValidatorsByDelegator(ctx context.Context, delegatorAddress common.Address, blockNumber uint64) ([]common.Address, error) {
	// construct input data
	payload, err := n.stakingSMC.Abi.Pack("getValidatorsByDelegator", delegatorAddress)
	if err != nil {
		return nil, err
	}

	res, err := n.CallAt(ctx, ConstructCallArgs(n.stakingSMC.ContractAddress.String(), payload), blockNumber)
	if err != nil {
		return nil, err
	}
	// get validators list of delegator
	var valAddrs struct {
		ValAddrs []common.Address
	}
	err = n.stakingSMC.Abi.UnpackIntoInterface(&valAddrs, "getValidatorsByDelegator", res)
	if err != nil {
		return nil, err
	}
	return valAddrs.ValAddrs, nil

}

func (n *Node) RPCValidator(ctx context.Context, address string) (*types.APIValidator, error) {
	var validator *types.APIValidator
	err := n.client.CallContext(ctx, &validator, "kai_validator", address, true)
	if err != nil {
		return nil, err
	}

	return validator, nil
}

func (n *Node) Validator(ctx context.Context, validatorSMCAddress string) (*types.Validator, error) {
	lgr := n.lgr.With(zap.String("method", "Validator"))
	payload, err := n.validatorSMC.Abi.Pack("inforValidator")
	if err != nil {
		lgr.Error("Error packing validator info payload: ", zap.Error(err))
		return nil, err
	}
	res, err := n.KardiaCall(ctx, constructCallArgs(validatorSMCAddress, payload))
	if err != nil {
		lgr.Error("GetValidatorInfo KardiaCall error: ", zap.Error(err))
		return nil, err
	}
	var valInfo types.RPCValidator
	// unpack result
	err = n.validatorSMC.Abi.UnpackIntoInterface(&valInfo, "inforValidator", res)
	if err != nil {
		lgr.Error("Error unpacking validator info: ", zap.Error(err))
		return nil, err
	}

	valInfo.ValStakingSmc = common.HexToAddress(validatorSMCAddress)
	validator := valInfo.ToValidator()
	return &validator, nil
}
