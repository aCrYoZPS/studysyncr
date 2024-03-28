package notes

type Note struct {
	ID      int    `json:"id"`
	Author  string `json:"author"`
	Content string `json:"content"`
}
