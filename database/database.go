package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (pgx.Tx, error)
	Close()
}

var DBInstance DB

func InitDB() (DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := pgxpool.New(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		return nil, err
	}

	DBInstance = db

	return db, nil
}

func QuerySingleRow(sql string, param []interface{}) pgx.Row {
	return DBInstance.QueryRow(context.Background(), sql, param...)
}

func QueryRow(sql string, param []interface{}) (pgx.Rows, error) {
	return DBInstance.Query(context.Background(), sql, param...)
}

func Exec(sql string, param []interface{}) error {
	_, err := DBInstance.Exec(context.Background(), sql, param...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Execute failed: %v\n", err)
		return err
	}

	return nil
}

func GetTx() (pgx.Tx, error) {
	return DBInstance.Begin(context.Background())
}

func TxExec(tx pgx.Tx, sql string, param []interface{}) error {
	_, err := tx.Exec(context.Background(), sql, param...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Transaction execute failed: %v\n", err)
		return err
	}
	return nil
}

func TxCommit(tx pgx.Tx) error {
	err := tx.Commit(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Transaction commit failed: %v\n", err)
		return err
	}

	return nil
}

func TxRollBack(tx pgx.Tx) error {
	err := tx.Rollback(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Transaction roll back failed: %v\n", err)
		return err
	}

	return nil
}
