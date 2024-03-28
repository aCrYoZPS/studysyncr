package storage

import (
	"fmt"

	"notes"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConnected struct {
	DB *gorm.DB
}

func (dbc *DBConnected) Init(dbadress string) error {
	var err error
	dbc.DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbadress,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return err
	}
	sqlDB, err := dbc.DB.DB()
	if err != nil {
		return err
	}
	if err = sqlDB.Ping(); err != nil {
		return err
	}
	fmt.Println("Connected")
	return nil
}

func (dbc *DBConnected) Add(note notes.Note) error {
	err := dbc.DB.Create(&note).Error
	if err != nil {
		return err
	}
	return nil
}

func (dbc *DBConnected) Get(ID int, Author string) {
}
