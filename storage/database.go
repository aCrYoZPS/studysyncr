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

func (dbc *DBConnected) Add(note *notes.Note) error {
	result := dbc.DB.Create(note)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (dbc *DBConnected) Get(ID int) (notes.Note, error) {
	var note notes.Note = notes.Note{}
	if err := dbc.DB.First(&note, ID).Error; err != nil {
		return note, err
	}
	return note, nil
}

func (dbc *DBConnected) GetList(Author string) ([]notes.Note, error) {
	var noteSlice []notes.Note
	result := dbc.DB.Where("author = ?", Author).Find(&noteSlice)
	if err := result.Error; err != nil {
		return noteSlice, err
	}
	return noteSlice, nil
}

func (dbc *DBConnected) Update(ID int, note *notes.Note) error {
	result := dbc.DB.Save(note)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (dbc *DBConnected) Delete(ID int) error {
	result := dbc.DB.Delete(&notes.Note{}, ID)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}
