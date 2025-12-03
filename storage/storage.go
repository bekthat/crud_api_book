package storage

import (
	"crud-api/models"
	"database/sql"
)

type Storage struct {
	DB *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{DB: db}
}

func (s *Storage) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		author TEXT
	);`
	_, err := s.DB.Exec(query)
	return err
}

func (s *Storage) Create(book *models.Book) error {
	res, err := s.DB.Exec("INSERT INTO books (title, author) VALUES (?, ?)", book.Title, book.Author)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	book.ID = int(id)
	return nil
}

func (s *Storage) GetBooks() ([]models.Book, error) {
	rows, err := s.DB.Query("SELECT id, title, author FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

func (s *Storage) GetBook(id int) (models.Book, error) {
	var b models.Book
	row := s.DB.QueryRow("SELECT id, title, author FROM books WHERE id = ?", id)
	err := row.Scan(&b.ID, &b.Title, &b.Author)
	return b, err
}

func (s *Storage) Update(id int, book models.Book) error {
	_, err := s.DB.Exec("UPDATE books SET title = ?, author = ? WHERE id = ?", book.Title, book.Author, id)
	return err
}

func (s *Storage) Delete(id int) error {
	_, err := s.DB.Exec("DELETE FROM books WHERE id = ?", id)
	return err
}
