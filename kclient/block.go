package kclient

import (
	"context"

	"github.com/kardiachain/go-kardia"
	"github.com/kardiachain/go-kardia/lib/common"

	"github.com/ngdlong91/kai-watcher/types"
)

// LatestBlockNumber gets the latest block number
func (n *Node) LatestBlockNumber(ctx context.Context) (uint64, error) {
	var result uint64
	err := n.client.CallContext(ctx, &result, "kai_blockNumber")
	return result, err
}

// BlockByHeight returns a block from the current canonical chain.
// Use HeaderByNumber if you don't need all transactions or uncle headers.
func (n *Node) BlockByHeight(ctx context.Context, height uint64) (*types.Block, error) {
	return n.getBlock(ctx, "kai_getBlockByNumber", height)
}

func (n *Node) getBlock(ctx context.Context, method string, args ...interface{}) (*types.Block, error) {
	var raw types.Block
	err := n.client.CallContext(ctx, &raw, method, args...)
	if err != nil {
		return nil, err
	}
	return &raw, nil
}

// BlockHeaderByNumber returns a block header from the current canonical chain.
func (n *Node) BlockHeaderByNumber(ctx context.Context, number uint64) (*types.Header, error) {
	return n.getBlockHeader(ctx, "kai_getBlockHeaderByNumber", number)
}

func (n *Node) getBlockHeader(ctx context.Context, method string, args ...interface{}) (*types.Header, error) {
	var raw types.Header
	err := n.client.CallContext(ctx, &raw, method, args...)
	if err != nil {
		return nil, err
	}
	return &raw, nil
}

func (n *Node) GetTransactionReceipt(ctx context.Context, txHash string) (*types.Receipt, error) {
	var r *types.Receipt
	err := n.client.CallContext(ctx, &r, "tx_getTransactionReceipt", common.HexToHash(txHash))
	if err == nil {
		if r == nil {
			return nil, kardia.NotFound
		}
	}
	return r, err
}
