package storage

import "notes"

type notesStorage interface {
	Add(note notes.Note) error
	Get(ID int, Author string) (notes.Note, error)
	GetList(Author string) (map[string]notes.Note, error)
	Update(ID int, Author string, note notes.Note) error
	Delete(ID int, Author string) error
}
