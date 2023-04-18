package kclient

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/kardiachain/go-kardia/lib/abi"
	"github.com/kardiachain/go-kardia/lib/common"
)

type Contract struct {
	Abi             *abi.ABI
	Bytecode        string
	ContractAddress common.Address
	OwnerAddress    common.Address
}

// parseBytesArrayIntoString is a utility function. It converts address, bytes and string arguments into their hex representation.
func parseBytesArrayIntoString(v interface{}) interface{} {
	vType := reflect.TypeOf(v).Kind()
	switch vType {
	case reflect.Array:
		//common.Address{}
		addr, ok := v.(common.Address)
		if ok {
			return common.Bytes(addr[:]).String()
		}

		hash, ok := v.([32]byte)
		if ok {
			return common.Bytes(hash[:]).String()
		}
		return v
	case reflect.Ptr:
		if value, ok := v.(*big.Int); ok {
			return value.String()
		}
	default:
		return v
	}
	return v
}

// getInputArguments get input arguments of a contract call
func (n *Node) getInputArguments(a *abi.ABI, name string, data []byte) (abi.Arguments, error) {
	var args abi.Arguments
	if method, ok := a.Methods[name]; ok {
		if len(data)%32 != 0 {
			return nil, fmt.Errorf("abi: improperly formatted output: %s - Bytes: [%+v]", string(data), data)
		}
		args = method.Inputs
	}
	if args == nil {
		return nil, ErrMethodNotFound
	}
	return args, nil
}

func (c *Contract) SetBytecode(bytecode string) {
	c.Bytecode = bytecode
}

func (c *Contract) SetOwnerAddress(address string) {
	c.OwnerAddress = common.HexToAddress(address)
}

func NewContract(abi *abi.ABI, addr common.Address) *Contract {
	c := &Contract{
		Abi:             abi,
		ContractAddress: addr,
	}
	return c
}

func (c *Contract) ABI() *abi.ABI {
	return c.Abi
}
