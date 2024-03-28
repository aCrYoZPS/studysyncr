package notes

import (
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	ID      int    `json:"id"`
	Author  string `json:"author"`
	Content string `json:"content"`
}
