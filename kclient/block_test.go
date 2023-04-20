package kclient

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNode_LatestBlockNumber(t *testing.T) {
	ctx := context.Background()
	node, _ := setupTestKit()
	blockNumber, err := node.LatestBlockNumber(ctx)
	assert.Nil(t, err)
	fmt.Println("Block number", blockNumber)
}
