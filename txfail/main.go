package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	kitConfig "github.com/getbread/gokit/config"
	"github.com/getbread/gokit/db"
	//"github.com/jackc/pgx/v4/stdlib"
	"time"
)

func hang() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()


	db := setupDB(kitConfig.Postgres{
		User:     "postgres",
		Password: "example",
		Host:     "localhost",
		Port:     "7890",
		Database: "postgres",
	})

	conn, err := db.Conn(ctx)

	if err != nil {
		return err
	}

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	const q = `BROKEN SQL TO FAIL`
	_, err = tx.ExecContext(ctx, q)
	if err == nil {
		return errors.New("failed to break transaction!")
	}

	fmt.Printf("Attempting to commit broken transaction\n")

	commitCh := make(chan error, 1)
	go func() {
		commitCh <- tx.Commit()
	}()

	select {
	case err := <-commitCh:
		return err
	case <-time.After(10 * time.Second):
		fmt.Printf("Commit has hung. ctx.Err()==%v\n", ctx.Err())
	}

	return <-commitCh
}

func main() {
	err := hang()
	fmt.Printf("hang returned %v\n", err)
}

func setupDB(config kitConfig.Postgres) *sql.DB {
	sqlDB, _ := db.NewDB(
		//db.WithDriver(stdlib.GetDefaultDriver()),
		db.DatabaseName(config.Database),
		db.User(config.User),
		db.Password(config.Password.Unmask()),
		db.Host(config.Host),
		db.Port(config.Port))

	sqlDB.SetMaxIdleConns(0)
	sqlDB.SetMaxOpenConns(2)

	return sqlDB
}