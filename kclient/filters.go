package kclient

import (
	"context"
	"fmt"

	kai "github.com/kardiachain/go-kardia"

	"github.com/ngdlong91/kai-watcher/types"
)

func (n *Node) FilterLogs(ctx context.Context, query kai.FilterQuery) ([]*types.Log, error) {
	var results []*types.Log
	args, err := toFilterArg(query)
	if err != nil {
		return nil, err
	}
	if err := n.client.CallContext(ctx, &results, "kai_getLogs", args); err != nil {
		return nil, err
	}
	return results, nil
}

func toFilterArg(q kai.FilterQuery) (interface{}, error) {
	arg := map[string]interface{}{
		"address": q.Addresses,
		"topics":  q.Topics,
	}
	if q.BlockHash != nil {
		arg["blockHash"] = *q.BlockHash
		if q.FromBlock != 0 || q.ToBlock != 0 {
			return nil, fmt.Errorf("cannot specify both BlockHash and FromBlock/ToBlock")
		}
	} else {
		if q.FromBlock == 0 {
			arg["fromBlock"] = uint64(1)
		} else {
			arg["fromBlock"] = q.FromBlock
		}
		arg["toBlock"] = q.ToBlock
		if q.ToBlock == 0 {
			arg["toBlock"] = "latest"
		}
	}
	return arg, nil
}
