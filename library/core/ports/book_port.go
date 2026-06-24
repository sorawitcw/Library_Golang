package ports

import "library/core/entity"

// BookRepository — สิ่งที่ Repository ต้องทำกับ DB
type BookRepository interface {
    Create(book entity.Book) (*entity.Book, error)
    FindAll() ([]entity.Book, error)
    FindById(id string) (*entity.Book, error)
    FindByISBN(isbn string) (*entity.Book, error)
    Update(book entity.Book) (*entity.Book, error)
    Delete(id string) error
}

// BookService — สิ่งที่ Handler จะเรียกใช้
type BookService interface {
    AddBook(book entity.Book) (*entity.BookResponse, error)
    GetAllBooks() ([]entity.BookResponse, error)
    GetBookById(id string) (*entity.BookResponse, error)
    UpdateBook(id string, book entity.Book) (*entity.BookResponse, error)
    DeleteBook(id string) error
}