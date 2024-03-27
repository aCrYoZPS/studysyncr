package entities

type Note struct {
	Id      int    `json:"id"`
	Author  string `json:"author"`
	Content string `json:"content"`
}
