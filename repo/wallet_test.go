package repo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ngdlong91/kai-watcher/entities"
)

func TestRepo_GetWallet(t *testing.T) {
	ctx := context.Background()
	pool, lgr := SetupDBTestKit(ctx)
	d := &Wallet{
		Pool:   pool,
		Logger: lgr,
	}
	name, err := d.Retrieve(ctx, "0xdd8b713e656919d1a643c71f1e271c55ca3fb7bc")

	assert.Nil(t, err)
	fmt.Println("name", name)
}

func TestRepo_InsertWallet(t *testing.T) {
	ctx := context.Background()
	pool, lgr := SetupDBTestKit(ctx)
	d := &Wallet{
		Pool:   pool,
		Logger: lgr,
	}
	entity := &entities.Wallet{
		Address:   "0x1",
		Name:      "Test",
		UpdatedAt: time.Now().Unix(),
	}
	err := d.Insert(ctx, entity)
	assert.Nil(t, err)
}
