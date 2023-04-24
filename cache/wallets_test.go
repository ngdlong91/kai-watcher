package cache

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCache_SetWallet(t *testing.T) {
	client, lgr := SetupTestCache()
	wallet := &Wallet{
		Client: client,
		Logger: lgr,
	}
	ctx := context.Background()
	assert.Nil(t, wallet.SetWallets(ctx, "123", "Test"))
}

func TestCache_GetWallet(t *testing.T) {
	client, lgr := SetupTestCache()
	wallet := &Wallet{
		Client: client,
		Logger: lgr,
	}
	ctx := context.Background()
	name, err := wallet.WalletName(ctx, "123")
	assert.Nil(t, err)
	assert.Equal(t, "Test", name)
}
