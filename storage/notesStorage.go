package storage

import "notes"

type notesStorage interface {
	Add(note *notes.Note) error
	Get(ID int) (notes.Note, error)
	GetList(Author string) ([]notes.Note, error)
	Update(ID int, note *notes.Note) error
	Delete(ID int) error
}
