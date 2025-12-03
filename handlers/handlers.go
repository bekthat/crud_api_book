package handlers

import (
	"crud-api/models"
	"crud-api/storage"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Handler struct {
	storage *storage.Storage
}

func New(storage *storage.Storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) BooksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		h.createBook(w, r)
	case "GET":
		h.getBooks(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) BookHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/books/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		h.getBook(w, r, id)
	case "PUT":
		h.updateBook(w, r, id)
	case "DELETE":
		h.deleteBook(w, r, id)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) createBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.storage.Create(&book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *Handler) getBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.storage.GetBooks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (h *Handler) getBook(w http.ResponseWriter, r *http.Request, id int) {
	book, err := h.storage.GetBook(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Книга не найдена", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *Handler) updateBook(w http.ResponseWriter, r *http.Request, id int) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.storage.Update(id, book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Книга с id %d обновлена", id)
}

func (h *Handler) deleteBook(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.storage.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Книга с id %d удалена", id)
}
