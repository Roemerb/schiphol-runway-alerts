package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/roemerb/schiphol-runway-alerts/config"
)

// OpenDatabase opens a new database connection
func OpenDatabase() (*sql.DB, error) {
	db, err := sql.Open("mysql", config.Load().DBCon)
	if err != nil {
		return &sql.DB{}, err
	}

	return db, nil
}
