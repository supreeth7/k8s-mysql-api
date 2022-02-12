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

type library struct {
	dbPass  string
	dbName  string
	apiPath string
}

type Book struct {
	Id     string
	Name   string
	Author string
}

func main() {

	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		dbPass = "password123"
	}

	apiPath := os.Getenv("API_PATH")
	if apiPath == "" {
		apiPath = "/apis/v1/books"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "library"
	}

	l := library{
		dbPass:  dbPass,
		dbName:  dbName,
		apiPath: apiPath,
	}

	r := mux.NewRouter()
	r.HandleFunc(apiPath, l.getBooks).Methods(http.MethodGet)
	r.HandleFunc(apiPath, l.createBook).Methods(http.MethodPost)
	http.ListenAndServe(":8080", r)
}

func (l *library) getBooks(w http.ResponseWriter, r *http.Request) {
	// Open DB Connection
	db, err := l.newConnection()
	if err != nil {
		log.Fatal("error:", err)
	}

	// Read Books
	rows, err := db.Query("select * from books")
	if err != nil {
		log.Fatalf("querying the books table %s\n", err.Error())
	}

	var books []Book

	for rows.Next() {
		book := Book{}
		err = rows.Scan(&book.Id, &book.Name, &book.Author)
		if err != nil {
			log.Fatal("error:", err)
		}

		books = append(books, book)
	}

	// encode in json
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		log.Fatal("error:", err)
	}

	// Close connection
	defer db.Close()
}

func (l *library) newConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", "root", l.dbPass, l.dbName))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (l *library) createBook(w http.ResponseWriter, req *http.Request) {
	book := Book{}
	json.NewDecoder(req.Body).Decode(&book)

	db, err := l.newConnection()
	if err != nil {
		log.Fatal("error:", err)
	}

	query := fmt.Sprintf("INSERT INTO books(id,title,author) values ('%s', '%s', '%s');", book.Id, book.Name, book.Author)

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal("error:", err)
	}

	res, err := stmt.Exec()
	if err != nil {
		log.Fatal("error:", err)
	}

	if rows, err := res.RowsAffected(); err != nil || rows < 1 {
		log.Fatal("error:", err)
	}

	db.Close()
}
