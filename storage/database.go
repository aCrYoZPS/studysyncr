package storage

import (
	"database/sql"
)

type DBConnected struct {
	db *sql.DB
}

func (dbc *DBConnected) init(dbadress string) error {
	var err error
	dbc.db, err = sql.Open("postgres", dbadress)
	if err != nil {
		return err
	}
	return nil
}

func (dbc *DBConnected) Add(note Note) {
}
