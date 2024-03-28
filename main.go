package main

import (
	"net/http"

	"storage"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

var (
	db     *storage.DBConnected = &storage.DBConnected{DB: new(gorm.DB)}
	conStr                      = "host=localhost user=postgres password=postgres dbname=test_db port=5432 sslmode=disable"
	router                      = gin.Default()
)

func main() {
	db.Init(conStr)
	router.GET("/", homePage)
	router.Run("localhost:8080")
}

func homePage(c *gin.Context) {
	c.String(http.StatusOK, "homePage")
}
