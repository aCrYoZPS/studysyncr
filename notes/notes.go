package notes

import (
	"encoding/json"
	"fmt"
)

type Note struct {
	ID      int     `json:"id"`
	Author  *string `json:"author"`
	Content *string `json:"content"`
}

func (note *Note) JsonifyNote() error {
	b, err := json.Marshal(note)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
