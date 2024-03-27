package entities

type notesStorage interface {
	Add(note Note) error
	Get(Id int, Author string) (Note, error)
	GetList(Author string) (map[int]Note, error)
	Update(Id int, Author string, note Note) error
	Delete(Id int, Author string) error
}
