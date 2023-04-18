package kclient

import (
	"context"

	"github.com/kardiachain/go-kardia/lib/common"
)

func (n *Node) Call(ctx context.Context, args SMCCallArgs) ([]byte, error) {
	var result common.Bytes
	err := n.client.CallContext(ctx, &result, "kai_kardiaCall", args, "latest")
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (n *Node) CallAt(ctx context.Context, args SMCCallArgs, blockNumber uint64) ([]byte, error) {
	var result common.Bytes
	err := n.client.CallContext(ctx, &result, "kai_kardiaCall", args, blockNumber)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (n *Node) GetBalance(ctx context.Context, account string) (string, error) {
	var (
		result string
		err    error
	)
	err = n.client.CallContext(ctx, &result, "account_balance", common.HexToAddress(account), "latest")
	return result, err
}
