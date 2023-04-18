package repo

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"

	"github.com/ngdlong91/kai-watcher/entities"
)

func GeneratePlaceholders(n int) string {
	if n <= 0 {
		return ""
	}

	var builder strings.Builder
	sep := ", "
	for i := 1; i <= n; i++ {
		if i == n {
			sep = ""
		}
		builder.WriteString("$" + strconv.Itoa(i) + sep)
	}

	return builder.String()
}

// InsertReturning entity using repo Executor
func InsertReturning(ctx context.Context, e entities.Entity, db *pgx.Conn,
	returnFieldName string,
	returnFieldValue interface{},
) error {
	fields := entities.GetFieldNames(e)
	fieldNames := strings.Join(fields, ",")
	placeHolders := GeneratePlaceholders(len(fields))
	stmt := "INSERT INTO " + e.TableName() +
		" (" + fieldNames + ") VALUES (" + placeHolders + ") RETURNING " + returnFieldName + ";"
	args := entities.GetScanFields(e, fields)
	return db.QueryRow(ctx, stmt, args...).Scan(returnFieldValue)
}

// Insert entity using repo Executor
func Insert(ctx context.Context, e entities.Entity, db *pgx.Conn) (pgconn.CommandTag, error) {
	fields := entities.GetFieldNames(e)
	fieldNames := strings.Join(fields, ",")
	fmt.Println("Fields", fieldNames)
	placeHolders := GeneratePlaceholders(len(fields))
	fmt.Println("Placeholders", placeHolders)
	stmt := "INSERT INTO " + e.TableName() + " (" + fieldNames + ") VALUES (" + placeHolders + ");"
	args := entities.GetScanFields(e, fields)
	return db.Exec(ctx, stmt, args...)
}

func GenerateUpdatePlaceholders(fields []string) string {
	var builder strings.Builder
	sep := ", "

	totalField := len(fields)
	for i, field := range fields {
		if i == totalField-1 {
			sep = ""
		}

		builder.WriteString(field + " = $" + strconv.Itoa(i+1) + sep)
	}

	return builder.String()
}

func transformError(err error) error {
	pqErr := err.(*pgconn.PgError)
	switch pqErr.Code {
	case pgerrcode.UniqueViolation:
		return errors.New("unique violation")
	default:
		return err
	}
}
