package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bobafetch/dashboard/utils"
	_ "github.com/denisenkom/go-mssqldb"
)

var db *sql.DB

func InitDB() error {
	var err error
	config, err := utils.LoadConfig()
	if err != nil {
		return err
	}

	conn := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", 
	config.SQL_HOST, config.SQL_USER, config.SQL_PASSWORD, config.SQL_PORT, config.SQL_DB)

	db, err = sql.Open("sqlserver", conn)
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func GetDB() *sql.DB {
	return db
}
