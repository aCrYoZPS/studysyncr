package storage

import (
	"fmt"

	"notes"
	"users"

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
	result := dbc.DB.Select("author", "content").Create(note)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (dbc *DBConnected) AddUser(user *users.User) error {
	var testUser users.User
	res := dbc.DB.Where("username = ?", user.Username).
		First(&testUser)
	if err := res.Error; err == nil {
		return fmt.Errorf("Username taken")
	}
	result := dbc.DB.Select("username", "password").Create(user)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (dbc *DBConnected) Get(ID int, author string) (notes.Note, error) {
	var note notes.Note = notes.Note{Author: nil}
	if err := dbc.DB.Where("author = ? AND id = ?", author, ID).First(&note).Error; err != nil {
		return note, err
	}
	return note, nil
}

func (dbc *DBConnected) GetUser(username string) (users.User, error) {
	var user users.User
	if err := dbc.DB.Where("username= ?", username).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (dbc *DBConnected) GetList(author string) ([]notes.Note, error) {
	var noteSlice []notes.Note
	result := dbc.DB.Where("author = ?", author).Find(&noteSlice)
	if err := result.Error; err != nil {
		return noteSlice, err
	}
	return noteSlice, nil
}

func (dbc *DBConnected) Update(ID int, author string, note *notes.Note) error {
	result := dbc.DB.Model(notes.Note{}).
		Where("author = ? AND id = ?", author, ID).
		Updates(notes.Note{Content: note.Content})
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (dbc *DBConnected) Delete(ID int, author string) error {
	var note notes.Note = notes.Note{Author: nil}
	result := dbc.DB.Where("author = ? AND id = ?", author, ID).Delete(&note)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}
