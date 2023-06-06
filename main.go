package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	API_PATH = "/api/v1/books"
)

type Book struct {
	Id, Name, Isbn string
}

type library struct {
	dbHost, dbPass, dbName string
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost:3306"
	}
	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		dbPass = "unacademy"
	}
	apiPath := os.Getenv("API_PATH")
	if apiPath == "" {
		apiPath = API_PATH
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "library"
	}
	l := library{
		dbHost: dbHost,
		dbPass: dbPass,
		dbName: dbName,
	}
	r := mux.NewRouter()
	r.HandleFunc(apiPath, l.getBooks).Methods("GET")
	http.ListenAndServe(":8080", r)
}

func (l library) getBooks(w http.ResponseWriter, r *http.Request) {
	// open connection
	db := l.openConnection()

	// read all books
	rows, err := db.Query("select * from books")
	if err != nil {
		log.Fatalf("querying the book table %s", err.Error())
	}

	books := []Book{}
	for rows.Next() {
		var id, name, isbn string
		err := rows.Scan(&id, &name, &isbn)
		if err != nil {
			log.Fatalf("error while scaning row %s", err.Error())
		}
		aBook := Book{
			Id:   id,
			Name: name,
			Isbn: isbn,
		}
		books = append(books, aBook)
	}

	json.NewEncoder(w).Encode(books)
	// close connection
	l.closeConnection(db)
	log.Println("getbooks was called")
}

func (l library) openConnection() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s", "root", l.dbPass, l.dbHost, l.dbName))
	if err != nil {
		log.Fatalf("opening the connection to db %s\n", err.Error())
	}
	return db
}

func (l library) closeConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatalf("error while closing conection %s", err.Error())
	}
}
