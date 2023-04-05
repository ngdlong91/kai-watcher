/*
 *  Copyright 2018 KardiaChain
 *  This file is part of the go-kardia library.
 *
 *  The go-kardia library is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU Lesser General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  The go-kardia library is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 *  GNU Lesser General Public License for more details.
 *
 *  You should have received a copy of the GNU Lesser General Public License
 *  along with the go-kardia library. If not, see <http://www.gnu.org/licenses/>.
 */
// Package kardia
package kardia

import (
	"context"
	"github.com/kardiachain/go-kardia/lib/abi"
	"math/big"
)

type Node interface {
	IsAlive() bool
	Info(ctx context.Context) (*NodeInfo, error)

	IAddress

	LatestBlockNumber(ctx context.Context) (uint64, error)
	BlockByHash(ctx context.Context, hash string) (*Block, error)
	BlockByHeight(ctx context.Context, height uint64) (*Block, error)
	BlockHeaderByHash(ctx context.Context, hash string) (*Header, error)
	BlockHeaderByNumber(ctx context.Context, number uint64) (*Header, error)
	GetTransaction(ctx context.Context, hash string) (*Transaction, error)
	GetTransactionReceipt(ctx context.Context, txHash string) (*Receipt, error)

	GetCirculatingSupply(ctx context.Context) (*big.Int, error)

	SendRawTransaction(ctx context.Context, tx string) error

	KardiaCall(ctx context.Context, args SMCCallArgs) ([]byte, error)
	IValidator

	DecodeInputData(to string, input string) (*FunctionCall, error)
	UnpackLog(log *Log, a *abi.ABI) (*Log, error)
	ValidatorABI() abi.ABI
	StakingABI() abi.ABI
}
