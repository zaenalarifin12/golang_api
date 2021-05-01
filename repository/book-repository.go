package repository

import (
	"book/entity"
	"gorm.io/gorm"
)

type BookRepository interface {
	InsertBook(book entity.Book) entity.Book
	UpdateBook(book entity.Book) entity.Book
	DeleteBook(book entity.Book)
	AllBook() []entity.Book
	FindBookByID(bookID uint64) entity.Book
}

type bookConnection struct {
	connection *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookConnection{
		connection: db,
	}
}

func (b *bookConnection) InsertBook(book entity.Book) entity.Book {
	b.connection.Save(&book)
	b.connection.Preload("User").Find(&b)
	return book
}

func (b *bookConnection) UpdateBook(book entity.Book) entity.Book {
	b.connection.Save(&book)
	b.connection.Preload("User").Find(&book)
	return book
}

func (b *bookConnection) DeleteBook(book entity.Book) {
	b.connection.Delete(&book)
}

func (b *bookConnection) AllBook() []entity.Book {
	var books []entity.Book
	b.connection.Preload("User").Find(&books)
	return books
}

func (b *bookConnection) FindBookByID(bookID uint64) entity.Book {
	var book entity.Book
	b.connection.Preload("User").Find(&book, bookID)
	return book
}
