package handlers

import (
	"net/http"
	"strconv"

	"auth"
	_ "docs"
	"notes"
	"storage"

	"github.com/gin-gonic/gin"
)

// @Summary      Get a note
// @Description  Returns JSON of note
// @Tags         notes
// @Param        id path int true "note id"
// @Param        Authorization header string true "Bearer "
// @Produce      json
// @Success      200   {object}  notes.Note
// @Router       /notes/{id} [get]
func GetNote(dbc *storage.DBConnected) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		author, err := auth.ExtractTokenUsername(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"errror": "Token shenanigans"})
			return
		}
		note, err := dbc.Get(id, author)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ID not found"})
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
// @Param        Authorization header string true "Bearer "
// @Produce      json
// @Success      200   {object}  []notes.Note
// @Router       /notes/ [get]
func GetAllNotes(dbc *storage.DBConnected) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		author, err := auth.ExtractTokenUsername(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"errror": "Token shenanigans"})
			return
		}
		notes, err := dbc.GetList(author)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author"})
			return
		}
		c.JSON(http.StatusOK, notes)
	}
	return gin.HandlerFunc(fn)
}

// @Summary      Delete a note
// @Description  Deletes a note from db
// @Tags         notes
// @Param        Authorization header string true "Bearer "
// @Param        id path int true "note id"
// @Produce      json
// @Success      200   {object}  notes.Note
// @Router       /notes/{id} [delete]
func DeleteNote(dbc *storage.DBConnected) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		author, err := auth.ExtractTokenUsername(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"errror": "Token shenanigans"})
			return
		}
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		err = dbc.Delete(id, author)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ID not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
	return gin.HandlerFunc(fn)
}

// @Summary      Post a note
// @Description  Upload a new note
// @Tags         notes
// @Param        note body notes.Note true "Note JSON"
// @Param        Authorization header string true "Bearer "
// @Accept       json
// @Produce      json
// @Success      200   {object}  notes.Note
// @Router       /notes [post]
func PostNote(dbc *storage.DBConnected) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var note notes.Note
		author, err := auth.ExtractTokenUsername(c)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"errror": "Token shenanigans"})
			return
		}
		err = c.ShouldBindJSON(&note)
		if *note.Author != author {
			c.JSON(
				http.StatusForbidden,
				gin.H{"error": "You can't upload notes that are not yours"},
			)
			return
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
			return
		}
		err = dbc.Add(&note)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldnt add a record"})
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
// @Param        id path int true "note id"
// @Param        Authorization header string true "Bearer "
// @Param        note body notes.Note true "Note JSON"
// @Accept       json
// @Produce      json
// @Success      200   {object}  notes.Note
// @Router       /notes/{id} [patch]
func PatchNote(dbc *storage.DBConnected) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		author, err := auth.ExtractTokenUsername(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"errror": "Token shenanigans"})
			return
		}
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var note notes.Note
		err = c.ShouldBindJSON(&note)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
			return
		}
		if *note.Author != author {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot access note that is not yours"})
			return
		}
		err = dbc.Update(id, author, &note)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		}
	}
	return gin.HandlerFunc(fn)
}
