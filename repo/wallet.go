package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"

	"github.com/ngdlong91/kai-watcher/entities"
)

type Wallet struct {
	Logger *zap.Logger
	Pool   *pgxpool.Pool
}

func (r *Wallet) Retrieve(ctx context.Context, wallet string) (string, error) {
	e := &entities.Wallet{}
	sql := fmt.Sprintf(`SELECT name FROM %s WHERE address = $1`, e.TableName())
	rows := r.Pool.QueryRow(ctx, sql, wallet)
	var name string
	if err := rows.Scan(&name); err != nil {
		return "", err
	}

	return name, nil
}

func (r *Wallet) Insert(ctx context.Context, e *entities.Wallet) error {
	fields := []string{"address", "name", "updated_at"}
	sql := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s) ON CONFLICT (address) DO UPDATE SET name = $%d, updated_at = $%d`,
		e.TableName(),
		strings.Join(fields, ` , `),
		GeneratePlaceholders(len(fields)),
		len(fields)+1,
		len(fields)+2,
	)
	values := entities.GetScanFields(e, fields)
	values = append(values, e.Name)
	values = append(values, e.UpdatedAt)
	if _, err := r.Pool.Exec(ctx, sql, values...); err != nil {
		return err
	}

	return nil
}
