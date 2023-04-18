package kclient

import (
	"context"
	kai "github.com/kardiachain/go-kardia"
	"github.com/kardiachain/go-kardia/lib/common"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestNode_FilterLogs(t *testing.T) {
	ctx := context.Background()
	node, lgr := setupTestKit()
	currentBlockNumber := uint64(0)
	latestBlock, err := node.LatestBlockNumber(ctx)
	assert.Nil(t, err)
	for {
		if (latestBlock - currentBlockNumber) > 2000 {
			filter := kai.FilterQuery{
				FromBlock: currentBlockNumber,
				ToBlock:   currentBlockNumber + 2000,
				Addresses: []common.Address{common.HexToAddress(StakingContractAddr)},
			}
			results, err := node.FilterLogs(context.Background(), filter)
			assert.Nil(t, err)
			lgr.Info("Query from ", zap.Uint64("From", currentBlockNumber), zap.Uint64("To", currentBlockNumber+2000))
			for _, r := range results {
				lgr.Info("Parse log", zap.Any("Log", r))
			}
			currentBlockNumber = currentBlockNumber + 2000
		}
	}

}
