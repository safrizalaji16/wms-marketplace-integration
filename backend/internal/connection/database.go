package connection

import (
	"backend/internal/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func GetDatabase(conf config.Database) *bun.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
		conf.Host,
		conf.Port,
		conf.User,
		conf.Password,
		conf.Name,
		conf.Tz,
	)
	sqldb, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Fatal("failed to open connection: ", err.Error())
	}

	err = sqldb.Ping()

	if err != nil {
		log.Fatal("failed to ping database: ", err.Error())
	}

	return bun.NewDB(sqldb, pgdialect.New())
}
