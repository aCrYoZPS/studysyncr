package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	_ "docs"
	"notes"
	"storage"

	"github.com/gin-gonic/gin"
)

// @Summary      Get a note
// @Description  Returns JSON of note
// @Tags         notes
// @Param        user path string true "current user"
// @Param        id path int true "note id"
// @Produce      json
// @Success      200   {object}  notes.Note
// @Router       /notes/{user}/{id} [get]
func GetNote(dbc *storage.DBConnected) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid ID")
			return
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

// @Summary      Get all notes of user
// @Description  Returns JSON array of notes
// @Tags         notes
// @Param        user path string true "current user"
// @Produce      json
// @Success      200   {object}  []notes.Note
// @Router       /notes/{user} [get]
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

// @Summary      Delete a note
// @Description  Deletes a note from db
// @Tags         notes
// @Param        user path string true "current user"
// @Param        id path int true "note id"
// @Produce      json
// @Success      200   {object}  notes.Note
// @Router       /notes/{user}/{id} [delete]
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

// @Summary      Post a note
// @Description  Upload a new note
// @Tags         notes
// @Param        user path string true "current user"
// @Param        note body notes.Note true "Note JSON"
// @Accept       json
// @Produce      json
// @Success      200   {object}  notes.Note
// @Router       /notes/{user} [post]
func PostNote(dbc *storage.DBConnected) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var note notes.Note
		err := c.ShouldBindJSON(&note)
		author := c.Param("user")
		if *note.Author != author {
			c.String(http.StatusForbidden, "You can't upload notes that are not yours")
			return
		}
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

// @Summary      Patch a note
// @Description  Updates a note in the db
// @Tags         notes
// @Param        user path string true "current user"
// @Param        id path int true "note id"
// @Param        note body notes.Note true "Note JSON"
// @Accept       json
// @Produce      json
// @Success      200   {object}  notes.Note
// @Router       /notes/{user}/{id} [patch]
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
		if *note.Author != author {
			c.String(http.StatusForbidden, "Cannot access note that is not yours")
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
