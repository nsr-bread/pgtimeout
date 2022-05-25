package main

import (
	"context"
	"database/sql"
	"fmt"
	kitConfig "github.com/getbread/gokit/config"
	"github.com/getbread/gokit/db"
	"github.com/jackc/pgx/v4/stdlib"
	"log"
	"strconv"
	"time"
)

func main() {
	db := setupDB(kitConfig.Postgres{
		User:     "postgres",
		Password: "example",
		Host:     "localhost",
		Port:     "7890",
		Database: "postgres",
	})


	for i := 0 ; i < 100; i++ {
		go func() {

			i := 100
			for {
				ctx := context.Background()

				tx, err := db.BeginTx(ctx, nil)
				if err != nil {
					log.Println("1" + err.Error())
				}

				iS := strconv.Itoa(i)

				query := fmt.Sprintf("insert into test (id)  values (%s)", iS)
				stmt, _ := tx.PrepareContext(ctx, query)
				stmt.ExecContext(ctx)

				err = tx.Commit()

				//time.Sleep(5 *time.Second)

				if err != nil {
					log.Println("1 commit error " + err.Error())
					err = tx.Rollback()
					if err != nil {
						log.Println("1 rollback error " + err.Error())
					}
				}

				time.Sleep(1 *time.Second)
			}


		}()
	}



	time.Sleep(5 *time.Second)

}

func setupDB(config kitConfig.Postgres) *sql.DB {
	sqlDB, _ := db.NewDB(
		db.WithDriver(stdlib.GetDefaultDriver()),
		db.DatabaseName(config.Database),
		db.User(config.User),
		db.Password(config.Password.Unmask()),
		db.Host(config.Host),
		db.Port(config.Port))

	return sqlDB
}