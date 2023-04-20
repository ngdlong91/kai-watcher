package kclient

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path"
	"strings"

	"github.com/kardiachain/go-kardia/lib/abi"
	"github.com/kardiachain/go-kardia/lib/abi/bind"
	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/lib/crypto"
	"github.com/kardiachain/go-kardia/rpc"
	"github.com/kardiachain/go-kardia/types"
	"github.com/shopspring/decimal"
)

var (
	ErrNotAValidatorAddress    = errors.New("address is not a validator")
	ErrMethodNotFound          = errors.New("abi: could not locate named method or event")
	ErrEmptyList               = errors.New("empty list")
	ErrParsingBigIntFromString = errors.New("cannot parse big.Int from string")
)
var (
	tenPoweredBy5  = new(big.Int).Exp(big.NewInt(10), big.NewInt(5), nil)
	tenPoweredBy4  = new(big.Int).Exp(big.NewInt(10), big.NewInt(4), nil)
	tenPoweredBy18 = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
)

func validatorName(data [32]byte) string {
	var name []byte
	for _, b := range data {
		if b != 0 {
			name = append(name, b)
		}
	}
	return string(name)
}

func convertBigIntToPercentage(raw string) (string, error) {
	input, ok := new(big.Int).SetString(raw, 10)
	if !ok {
		return "", ErrParsingBigIntFromString
	}
	tmp := new(big.Int).Mul(input, tenPoweredBy18)
	result := new(big.Int).Div(tmp, tenPoweredBy18).String()
	result = fmt.Sprintf("%020s", result)
	result = strings.TrimLeft(strings.TrimRight(strings.TrimRight(result[:len(result)-16]+"."+result[len(result)-16:], "0"), "."), "0")
	if strings.HasPrefix(result, ".") {
		result = "0" + result
	}
	return result, nil
}

func CalculateVotingPower(raw, total *big.Int) (string, error) {

	//valStakedAmount, ok := new(big.Int).SetString(raw, 10)
	//if !ok {
	//	return "", ErrParsingBigIntFromString
	//}
	valStakedAmount := new(big.Int).Set(raw)
	tmp := new(big.Int).Mul(valStakedAmount, tenPoweredBy5)
	result := new(big.Int).Div(tmp, total).String()
	result = fmt.Sprintf("%020s", result)
	result = strings.TrimLeft(strings.TrimRight(strings.TrimRight(result[:len(result)-3]+"."+result[len(result)-3:], "0"), "."), "0")
	if strings.HasPrefix(result, ".") {
		result = "0" + result
	}
	if strings.EqualFold(result, "") {
		result = "0"
	}
	return result, nil
}

func getParamsSMCAddress(stakingSMC *Contract, client *rpc.Client) (common.Address, error) {
	payload, err := stakingSMC.Abi.Pack("params")
	if err != nil {
		return common.Address{}, err
	}

	var (
		res common.Bytes
		ctx = context.Background()
	)
	err = client.CallContext(ctx, &res, "kai_kardiaCall", ConstructCallArgs(stakingSMC.ContractAddress.Hex(), payload), "latest")
	if err != nil {
		return common.Address{}, err
	}

	var result struct {
		ParamsSmcAddr common.Address
	}
	err = stakingSMC.Abi.UnpackIntoInterface(&result, "params", res)
	if err != nil {
		return common.Address{}, err
	}

	return result.ParamsSmcAddr, nil
}

func ConstructCallArgs(address string, payload []byte) SMCCallArgs {
	return SMCCallArgs{
		From:     address,
		To:       &address,
		Gas:      100000000,
		GasPrice: big.NewInt(0),
		Value:    big.NewInt(0),
		Data:     common.Bytes(payload).String(),
	}
}

// NewKeyedTransactor is a utility method to easily create a transaction signer
// from a single private key.
func NewKeyedTransactor(key *ecdsa.PrivateKey) *bind.TransactOpts {
	keyAddr := crypto.PubkeyToAddress(key.PublicKey)
	return &bind.TransactOpts{
		From: keyAddr,
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != keyAddr {
				return nil, errors.New("not authorized to sign this account")
			}
			signature, err := crypto.Sign(signer.Hash(tx).Bytes(), key)
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
	}
}

func DecodeWithABI(input string, a *abi.ABI) (*FunctionCall, error) {
	// return nil if input data is too short
	if len(input) <= 2 {
		return nil, nil
	}
	data, err := hex.DecodeString(strings.TrimLeft(input, "0x"))
	if err != nil {
		return nil, err
	}
	sig := data[0:4] // get the function signature (first 4 bytes of input data)
	method, err := a.MethodById(sig)
	if err != nil {
		return nil, err
	}

	fmt.Println("Method", method)
	// exclude the function signature, only decode and unpack the arguments
	var body []byte
	if len(data) <= 4 {
		body = []byte{}
	} else {
		body = data[4:]
	}
	args, err := getInputArguments(a, method.Name, body)
	if err != nil {
		return nil, err
	}
	arguments := make(map[string]interface{})
	err = args.UnpackIntoMap(arguments, body)
	if err != nil {
		return nil, err
	}
	// convert address, bytes and string arguments into their hex representations
	for i, arg := range arguments {
		arguments[i] = parseBytesArrayIntoString(arg)
	}
	return &FunctionCall{
		Function:   method.String(),
		MethodID:   "0x" + hex.EncodeToString(sig),
		MethodName: method.Name,
		Arguments:  arguments,
	}, nil
}

func getInputArguments(a *abi.ABI, name string, data []byte) (abi.Arguments, error) {
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

func CheckAddress(address string) string {
	return common.HexToAddress(address).String()
}

func GenerateWallet() (common.Address, ecdsa.PrivateKey, error) {
	privKey, err := crypto.GenerateKey()
	if err != nil {
		return common.Address{}, ecdsa.PrivateKey{}, err
	}
	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, ecdsa.PrivateKey{}, fmt.Errorf("error casting public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return address, *privKey, nil
}

func FloatToBigInt(amount float64, decimals int32) *big.Int {
	return decimal.NewFromFloat(amount).Mul(decimal.New(1, decimals)).BigInt()
}

func constructCallArgs(address string, payload []byte) SMCCallArgs {
	return SMCCallArgs{
		From:     address,
		To:       &address,
		Gas:      100000000,
		GasPrice: big.NewInt(0),
		Value:    big.NewInt(0),
		Data:     common.Bytes(payload).String(),
	}
}

func readABIFromFile(filePath string) (abi.ABI, error) {
	wd, _ := os.Getwd()
	abiFileContent, err := os.ReadFile(path.Join(wd, filePath))
	if err != nil {
		return abi.ABI{}, fmt.Errorf("cannot read ABI file %s", filePath)
	}
	//var smcABIStr string
	//if err := json.Unmarshal(abiFileContent, &smcABIStr); err != nil {
	//	return abi.ABI{}, err
	//}
	smcABI, err := abi.JSON(strings.NewReader(string(abiFileContent)))
	if err != nil {
		return abi.ABI{}, err
	}
	return smcABI, nil
}

func (n *Node) KardiaCall(ctx context.Context, args SMCCallArgs) ([]byte, error) {
	var result common.Bytes
	err := n.client.CallContext(ctx, &result, "kai_kardiaCall", args, "latest")
	if err != nil {
		return nil, err
	}
	return result, nil
}
