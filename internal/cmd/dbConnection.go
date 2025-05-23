package cmd

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/harrydayexe/Omni/internal/config"
)

func GetDBConnection(config config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", config.DataSourceName)
	if err != nil {
		return nil, err
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * time.Duration(config.ConnMaxLifetime))
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	return db, nil
}
