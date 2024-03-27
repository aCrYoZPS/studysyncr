package main

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init_db() {
	var err error
	conStr := "postgres://postgres:postgres@localhost/test_db?sslmode=disable"
	db, err = sql.Open("postgres", conStr)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected")
}

func main() {
	init_db()
	defer db.Close()
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
