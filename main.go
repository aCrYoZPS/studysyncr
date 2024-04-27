package main

import (
	"fmt"
	"net/http"
	"os"

	"auth"
	_ "docs"
	"handlers"
	"storage"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // swagger embed files
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"gorm.io/gorm"
)

var (
	db        *storage.DBConnected = &storage.DBConnected{DB: new(gorm.DB)}
	userDB    *storage.DBConnected = &storage.DBConnected{DB: new(gorm.DB)}
	router    *gin.Engine          = gin.Default()
	secretKey []byte
)

// @title           Studysyncr API
// @version         1.0
// @description     API for Studysyncr practice project
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("file error")
	}
	DB_HOST := os.Getenv("DB_HOST")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	NOTE_DB_NAME := os.Getenv("NOTE_DB_NAME")
	USERS_DB_NAME := os.Getenv("USERS_DB_NAME")
	PORT := os.Getenv("PORT")
	SSL_MODE := os.Getenv("SSL_MODE")
	conStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		DB_HOST,
		DB_USER,
		DB_PASSWORD,
		NOTE_DB_NAME,
		PORT,
		SSL_MODE,
	)
	usersConStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		DB_HOST,
		DB_USER,
		DB_PASSWORD,
		USERS_DB_NAME,
		PORT,
		SSL_MODE,
	)
	fmt.Println(conStr)
	fmt.Println(usersConStr)
	fmt.Print(os.Getenv("AUTH_KEY"))
	db.Init(conStr)
	userDB.Init(usersConStr)
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "")
	})
	router.POST("/register", auth.Register(userDB))
	router.POST("/authorise", auth.Authorise(userDB))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("notes/:id", handlers.GetNote(db))
	router.GET("notes/", handlers.GetAllNotes(db))
	router.POST("notes/", handlers.PostNote(db))
	router.DELETE("notes/:id", handlers.DeleteNote(db))
	router.PATCH("notes/:id", handlers.PatchNote(db))
	router.Run("localhost:8080")
}
