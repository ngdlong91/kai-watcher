package kclient

import (
	"context"
	"fmt"
	"github.com/ngdlong91/kai-watcher/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidator_APIValidatorInfo(t *testing.T) {
	ctx := context.Background()
	node, _ := setupTestKit()

	validatorAddresses, err := node.ValidatorSMCAddresses(ctx, 0)
	assert.Nil(t, err)
	fmt.Println("Length", len(validatorAddresses))
	for _, addr := range validatorAddresses {
		validator, err := node.APIValidatorInfo(ctx, addr.String())
		assert.Nil(t, err)
		fmt.Printf("Validator: %+v \n", validator)
	}

}

func TestValidator_GetDelegations(t *testing.T) {
	const (
		SampleValidatorAddr = "0xf151515fa44527e203cb457086cda630da80c4b8"
	)
	var delegatorMap map[string]types.DelegatorInfo
	node, _ := setupTestKit()
	ctx := context.Background()

	validatorInfo, err := node.RPCValidator(ctx, SampleValidatorAddr)
	assert.Nil(t, err)
	fmt.Println("-----------------")
	fmt.Printf("%+v", validatorInfo)
	for _, d := range validatorInfo.Delegators {
		fmt.Printf("Delegator: %+v \n", d)
		info := delegatorMap[d.Address]
		if info.Address == "" {
			info.Address = d.Address
		}
		record := types.ValidatorStakeRecord{
			ValidatorName:       validatorInfo.Name,
			ValidatorAddress:    validatorInfo.Address,
			ValidatorSMCAddress: validatorInfo.SmcAddress,
			Amount:              d.StakedAmount,
		}
		info.ValidatorRecords = append(info.ValidatorRecords, record)
		fmt.Printf("DelegatorInfo: %+v \n", info)
	}
}
