package cmd

import (
	"database/sql"
	"time"

	"github.com/harrydayexe/Omni/internal/config"
)

func GetDBConnection(config config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", config.DataSourceName)
	if err != nil {
		return nil, err
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * time.Duration(config.ConnMaxLifetime))
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db, nil
}
