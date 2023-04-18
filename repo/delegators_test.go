package repo

import (
	"context"
	"github.com/joho/godotenv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ngdlong91/kai-watcher/entities"
)

func init() {
	if err := godotenv.Load("dev.env"); err != nil {
		panic(err)
	}
}

func TestDelegator_Insert(t *testing.T) {
	ctx := context.Background()
	conn, lgr := SetupDBTestKit(ctx)
	d := &Delegator{
		Conn:   conn,
		Logger: lgr,
	}
	entity := &entities.Delegator{
		Address:      "0x1",
		StakedAmount: "1",
		UpdatedAt:    time.Now().Unix(),
	}
	err := d.Insert(ctx, entity)
	assert.Nil(t, err)
}

func TestDelegator_BulkInsert(t *testing.T) {
	ctx := context.Background()
	conn, lgr := SetupDBTestKit(ctx)
	d := &Delegator{
		Conn:   conn,
		Logger: lgr,
	}
	assert.Nil(t, d.DeleteAll(ctx))
	delegators := []entities.Delegator{
		{
			Address:      "0x1",
			StakedAmount: "1",
			UpdatedAt:    time.Now().Unix(),
		},
		{
			Address:      "0x2",
			StakedAmount: "2",
			UpdatedAt:    time.Now().Unix(),
		},
	}
	var records [][]interface{}
	for _, d := range delegators {
		var data []interface{}
		data = append(data, d.Address)
		data = append(data, d.StakedAmount)
		data = append(data, d.UpdatedAt)
		records = append(records, data)
	}

	assert.Nil(t, d.BulkInsert(ctx, records))
}
