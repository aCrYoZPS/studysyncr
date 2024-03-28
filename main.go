package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"storage"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

// var db *sql.DB
//
// func init_db() {
// 	var err error
// 	conStr := "postgres://postgres:postgres@localhost/test_db?sslmode=disable"
// 	db, err = sql.Open("postgres", conStr)
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = db.Ping()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Connected")
// }

func main() {
	conStr := "host=localhost user=postgres password=postgres dbname=test_db port=5432 sslmode=disable"
	var db *storage.DBConnected = &storage.DBConnected{DB: new(gorm.DB)}
	db.Init(conStr)
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprint(rw, "Hello Web")
	})
	http.HandleFunc("/echo", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprint(rw, req.URL.Query().Get("message"))
	})
	http.HandleFunc("/shape", func(rw http.ResponseWriter, req *http.Request) {
		area, err := strconv.ParseFloat(req.URL.Query().Get("radius"), 64)
		area = area * area * math.Pi
		if err != nil {
			panic(err)
		}
		fmt.Fprint(rw, area)
	})
	// insertSth := "insert into Notes (id, author, content) values (2, 'Kolya', 'SNAPPS')"
	// _, err := db.Exec(insertSth)
	// if err != nil {
	// 	panic(err)
	// }
	http.ListenAndServe(":8080", nil)
}
