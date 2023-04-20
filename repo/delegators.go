package repo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"strings"

	"github.com/ngdlong91/kai-watcher/entities"
)

type Delegator struct {
	Conn   *pgxpool.Pool
	Logger *zap.Logger
}

func (r *Delegator) DeleteAll(ctx context.Context) error {
	e := entities.Delegator{}
	sql := fmt.Sprintf("DELETE FROM %s", e.TableName())
	if _, err := r.Conn.Exec(ctx, sql); err != nil {
		return err
	}
	return nil
}

func (r *Delegator) BulkInsert(ctx context.Context, records [][]interface{}) error {
	e := entities.Delegator{}
	fields := []string{"address", "staked_amount", "updated_at", "balance", "total_amount"}
	if _, err := r.Conn.CopyFrom(
		ctx, pgx.Identifier{e.TableName()}, fields, pgx.CopyFromRows(records),
	); err != nil {
		return err
	}
	return nil
}

func (r *Delegator) Insert(ctx context.Context, e *entities.Delegator) error {
	fields := []string{"address", "staked_amount", "updated_at", "balance", "total_amount"}
	sql := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s) ON CONFLICT (address) DO UPDATE SET staked_amount = $%d, updated_at = $%d, balance = $%d`,
		e.TableName(),
		strings.Join(fields, ` , `),
		GeneratePlaceholders(len(fields)),
		len(fields)+1,
		len(fields)+2,
		len(fields)+3,
	)
	values := entities.GetScanFields(e, fields)
	values = append(values, e.Address)
	values = append(values, e.StakedAmount)
	values = append(values, e.UpdatedAt)
	values = append(values, e.Balance)
	values = append(values, e.Balance)
	if _, err := r.Conn.Exec(ctx, sql, values...); err != nil {
		return err
	}

	return nil
}
