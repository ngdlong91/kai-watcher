// Package kardia
package kardia

import (
	"encoding/hex"
	"fmt"
	"go.uber.org/zap"

	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/kardiachain/go-kardia/lib/abi"
	"github.com/kardiachain/go-kardia/lib/common"
)

func CheckAddress(address string) string {
	return common.HexToAddress(address).String()
}
func (n *node) UnpackLog(log *Log, smcABI *abi.ABI) (*Log, error) {
	if strings.HasPrefix(log.Data, "0x") {
		log.Data = log.Data[2:]
	}
	log.Address = CheckAddress(log.Address)
	event, err := smcABI.EventByID(common.HexToHash(log.Topics[0]))
	if err != nil {
		return nil, err
	}
	argumentsValue := make(map[string]interface{})
	err = unpackLogIntoMap(smcABI, argumentsValue, event.RawName, log)
	if err != nil {
		return nil, err
	}
	//convert address, bytes and string arguments into their hex representations
	for i, arg := range argumentsValue {
		argumentsValue[i] = parseBytesArrayIntoString(arg)
	}
	// append unpacked data
	log.Arguments = argumentsValue
	log.MethodName = event.RawName
	order := int64(1)
	for _, arg := range event.Inputs {
		if arg.Indexed {
			log.ArgumentsName += "index_topic_" + strconv.FormatInt(order, 10) + " "
			order++
		}
		log.ArgumentsName += arg.Type.String() + " " + arg.Name + ", "
	}
	log.ArgumentsName = strings.TrimRight(log.ArgumentsName, ", ")
	return log, nil
}

// UnpackLogIntoMap unpacks a retrieved log into the provided map.
func unpackLogIntoMap(a *abi.ABI, out map[string]interface{}, eventName string, log *Log) error {
	lgr, _ := zap.NewDevelopment()
	data, err := hex.DecodeString(log.Data)
	if err != nil {
		return err
	}

	if len(data) > 0 {
		if err := a.UnpackIntoMap(out, eventName, data); err != nil {
			return err
		}
	}
	lgr.Info("Event Name", zap.String("Event", eventName), zap.Any("Inputs", a.Events[eventName].Inputs))
	// unpacking indexed arguments
	var indexed abi.Arguments
	for _, arg := range a.Events[eventName].Inputs {
		lgr.Info("Args", zap.Any("Arg", arg))
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	lgr.Info("Indexed", zap.Any("IndexedSize", len(indexed)), zap.Any("Indexed", indexed))

	topicSize := len(log.Topics)
	if topicSize <= 1 {
		return nil
	}
	topics := make([]common.Hash, len(log.Topics)-1)
	for i, topic := range log.Topics[1:] { // exclude the eventID (log.Topic[0])
		topics[i] = common.HexToHash(topic)
	}
	lgr.Info("Topics", zap.Any("TopicSize", len(topics)), zap.Any("Topics", topics))
	return abi.ParseTopicsIntoMap(out, indexed, topics)
}

// DecodeInputData returns decoded transaction input data if it match any function in staking and validator contract.
func (n *node) DecodeInputData(to string, input string) (*FunctionCall, error) {
	// return nil if input data is too short
	if len(input) <= 2 {
		return nil, nil
	}
	data, err := hex.DecodeString(strings.TrimLeft(input, "0x"))
	if err != nil {
		return nil, err
	}
	sig := data[0:4] // get the function signature (first 4 bytes of input data)
	var (
		a      *abi.ABI
		method *abi.Method
	)
	// check if the to address is staking contract, then we search for staking method in staking contract ABI
	if n.stakingSMC.ContractAddress.Equal(common.HexToAddress(to)) {
		fmt.Println("--------using stakingABI")
		a = n.stakingSMC.Abi
		method, err = n.stakingSMC.Abi.MethodById(sig)
		if err != nil {
			return nil, err
		}
	} else { // otherwise, search for a validator method
		fmt.Println("--------using validatorABI")
		a = n.validatorSMC.Abi
		method, err = n.validatorSMC.Abi.MethodById(sig)
		if err != nil {
			return nil, err
		}
	}
	// exclude the function signature, only decode and unpack the arguments
	var body []byte
	if len(data) <= 4 {
		body = []byte{}
	} else {
		body = data[4:]
	}
	args, err := n.getInputArguments(a, method.Name, body)
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

// getInputArguments get input arguments of a contract call
func (n *node) getInputArguments(a *abi.ABI, name string, data []byte) (abi.Arguments, error) {
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
