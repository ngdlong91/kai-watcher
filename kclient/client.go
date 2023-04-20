package kclient

import (
	"fmt"
	"github.com/kardiachain/go-kardia/lib/abi"
	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/rpc"
	"go.uber.org/zap"
)

const (
	StakingContractAddr = "0x0000000000000000000000000000000000001337"
)

type Node struct {
	client *rpc.Client
	isLive bool
	url    string

	lgr *zap.Logger

	// SMC
	stakingSMC   *Contract
	validatorSMC *Contract
	paramsSMC    *Contract
	krc20SMC     *Contract
	krc721SMC    *Contract
}

func NewNode(url string, lgr *zap.Logger) (*Node, error) {
	rpcClient, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	node := &Node{
		client: rpcClient,
		url:    url,
		lgr:    lgr,
	}
	if err := node.setupSMC(); err != nil {
		return nil, err
	}
	return node, nil
}

func (n *Node) StakingABI() *abi.ABI {
	return n.stakingSMC.Abi
}

func (n *Node) ValidatorABI() *abi.ABI {
	return n.validatorSMC.Abi
}

func (n *Node) setupSMC() error {
	stakingABI, err := readABIFromFile("/abi/staking.json")
	if err != nil {
		fmt.Println("---staking", err)
		return err
	}
	stakingUtil := &Contract{
		Abi:             &stakingABI,
		ContractAddress: common.HexToAddress(StakingContractAddr),
	}
	n.stakingSMC = stakingUtil
	validatorSmcAbi, err := readABIFromFile("/abi/validator.json")
	if err != nil {
		fmt.Println("---validator", err)
		return err
	}
	validatorUtil := &Contract{
		Abi: &validatorSmcAbi,
	}
	n.validatorSMC = validatorUtil
	paramsSmcAddr, err := getParamsSMCAddress(stakingUtil, n.client)
	if err != nil {
		return err
	}
	paramsSmcAbi, err := readABIFromFile("/abi/params.json")
	if err != nil {
		fmt.Println("---params", err)
		return err
	}
	paramsUtil := &Contract{
		Abi:             &paramsSmcAbi,
		ContractAddress: paramsSmcAddr,
	}
	n.paramsSMC = paramsUtil

	krc20SmcABI, err := readABIFromFile("/abi/krc20.json")
	if err != nil {
		fmt.Println("---krc20", err)
		return err
	}
	krc20Util := &Contract{
		Abi: &krc20SmcABI,
	}
	n.krc20SMC = krc20Util

	return nil
}

///////////////////////////
//	TESTING
///////////////////////////

func setupTestKit() (*Node, *zap.Logger) {
	lgr, _ := zap.NewDevelopment()
	lgr = lgr.With(zap.Namespace("TestNode"))
	node, err := NewNode("https://dev.kardiachain.io", lgr)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	return node, lgr
}
