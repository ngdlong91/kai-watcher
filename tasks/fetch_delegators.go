package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"github.com/ngdlong91/kai-watcher/kclient"
	"github.com/ngdlong91/kai-watcher/types"
	"github.com/ngdlong91/kai-watcher/utils"
)

type FetchDelegators struct {
	Logger      *zap.Logger
	Pool        *pgxpool.Pool
	Node        *kclient.Node
	DelegatorDB interface {
		DeleteAll(ctx context.Context) error
		BulkInsert(ctx context.Context, records [][]interface{}) error
	}
}

func (t *FetchDelegators) Start(ctx context.Context) error {
	//TODO implement me
	panic("implement me")

}

func (t *FetchDelegators) IsTriggerSpecialEvent(ctx context.Context) bool {
	//TODO implement me
	panic("implement me")
}

func (t *FetchDelegators) Execute() {
	ctx := context.Background()
	lgr := utils.LoggerForMethod(t.Logger, "FetchDelegator")
	lgr.Info("Start fetch delegator")
	records, err := t.fetch()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := t.DelegatorDB.DeleteAll(ctx); err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := t.DelegatorDB.BulkInsert(ctx, records); err != nil {
		fmt.Println(err.Error())
		return
	}

}

var headers = []string{"ID", "Delegator", "Balance"}

func (t *FetchDelegators) fetch() ([][]interface{}, error) {
	var validatorSigners []string

	ctx := context.Background()
	validatorAddresses, err := t.Node.ValidatorSMCAddresses(ctx, 0)
	if err != nil {
		return nil, err
	}

	for _, addr := range validatorAddresses {
		validator, err := t.Node.APIValidatorInfo(ctx, addr.String())
		if err != nil {
			return nil, err
		}
		headers = append(headers, validator.Name)

		validatorSigners = append(validatorSigners, validator.Signer)
	}

	delegatorMap := make(map[string]types.DelegatorInfo)
	var delegatorList []string
	for _, s := range validatorSigners {
		validatorInfo, err := t.Node.RPCValidator(ctx, s)
		if err != nil {
			return nil, err
		}
		var info types.DelegatorInfo
		for _, d := range validatorInfo.Delegators {
			info = delegatorMap[d.Address]
			if info.Address == "" {
				var boilerStakeRecord []types.ValidatorStakeRecord
				for _, sClone := range validatorSigners {
					stakeRecord := types.ValidatorStakeRecord{
						ValidatorAddress: sClone,
						Amount:           "",
					}
					boilerStakeRecord = append(boilerStakeRecord, stakeRecord)
				}
				info.Address = d.Address
				info.ValidatorRecords = boilerStakeRecord
				delegatorList = append(delegatorList, info.Address)
			}
			for id, r := range info.ValidatorRecords {
				if r.ValidatorAddress == s {
					stakedAmount := utils.ToDecimal(d.StakedAmount, 18)
					info.ValidatorRecords[id].Amount = stakedAmount.String()
					currentTotalStaked := utils.ToDecimal(info.TotalStaked, 18)
					info.TotalStaked, _ = decimal.Sum(currentTotalStaked, stakedAmount).Float64()
				}
			}
			fmt.Printf("Delegator info final: %+v \n", info)

			delegatorMap[d.Address] = info
		}
	}

	fmt.Println("Total Delegator: ", len(delegatorList))
	var records [][]interface{}
	now := time.Now().Unix()
	for _, dAddr := range delegatorList {
		info := delegatorMap[dAddr]
		balance, err := t.Node.GetBalance(ctx, dAddr)
		if err != nil {
			return nil, err
		}
		info.Balance, _ = utils.ToDecimal(balance, 18).Float64()
		info.TotalAmount = info.TotalStaked + info.Balance
		fmt.Printf("Address: %s | DelegatorInfo: %+v \n", dAddr, info)

		for _, record := range info.ValidatorRecords {
			if record.Amount == "" {
				record.Amount = "0"
			}
		}

		var data []interface{}
		data = append(data, info.Address)
		data = append(data, info.TotalStaked)
		data = append(data, now)
		data = append(data, info.Balance)
		data = append(data, info.TotalAmount)
		records = append(records, data)
	}

	return records, nil
}
