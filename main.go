package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"crud-api/handlers"
	"crud-api/storage"

	_ "modernc.org/sqlite"
)

func main() {

	db, err := sql.Open("sqlite", "./library.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := storage.New(db)
	handler := handlers.New(store)

	if err := store.CreateTable(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/books", handler.BooksHandler)
	http.HandleFunc("/books/", handler.BookHandler)

	fmt.Println("Сервер запущен на :8080 (Чистая архитектура)...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
