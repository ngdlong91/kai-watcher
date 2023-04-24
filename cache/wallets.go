package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type Wallet struct {
	Client *redis.Client
	Logger *zap.Logger
}

const (
	KeyWallet = "wallet#%s"
)

func (c *Wallet) SetWallets(ctx context.Context, wallet string, name string) error {
	key := fmt.Sprintf(KeyWallet, wallet)
	if err := c.Client.Set(ctx, key, name, 0).Err(); err != nil {
		return err
	}
	return nil

}

func (c *Wallet) WalletName(ctx context.Context, wallet string) (string, error) {
	return c.Client.Get(ctx, fmt.Sprintf(KeyWallet, wallet)).Result()

}
