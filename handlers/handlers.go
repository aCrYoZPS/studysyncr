package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"notes"
	"storage"

	"github.com/gin-gonic/gin"
)

func GetNote(dbc *storage.DBConnected) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid ID")
		}
		author := c.Param("user")
		fmt.Println(author)
		note, err := dbc.Get(id, author)
		if err != nil {
			c.String(http.StatusNotFound, "ID not found")
			return
		} else {
			c.JSON(http.StatusOK, note)
		}
	}
	return gin.HandlerFunc(fn)
}

func GetAllNotes(dbc *storage.DBConnected) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		author := c.Param("user")
		notes, err := dbc.GetList(author)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid author")
			return
		}
		if len(notes) == 0 {
			c.String(http.StatusBadRequest, "Author not found")
			return
		}
		c.JSON(http.StatusOK, notes)
	}
	return gin.HandlerFunc(fn)
}

func DeleteNote(dbc *storage.DBConnected) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		author := c.Param("user")
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}
		err = dbc.Delete(id, author)
		if err != nil {
			c.String(http.StatusNotFound, "ID not found")
			return
		}
	}
	return gin.HandlerFunc(fn)
}

func PostNote(dbc *storage.DBConnected) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var note notes.Note
		err := c.ShouldBindJSON(&note)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid payload")
			return
		}
		err = dbc.Add(&note)
		if err != nil {
			c.String(http.StatusBadRequest, "Couldnt add a record")
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		}
	}
	return gin.HandlerFunc(fn)
}

func PatchNote(dbc *storage.DBConnected) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		author := c.Param("user")
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}
		var note notes.Note
		err = c.ShouldBindJSON(&note)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid payload")
			return
		}
		err = dbc.Update(id, author, &note)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid data")
			return
		}
		c.String(http.StatusOK, "Updated successfully")
	}
	return gin.HandlerFunc(fn)
}
