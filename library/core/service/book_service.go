package service

import (
    "errors"
    "library/core/entity"
    "library/core/ports"

    "github.com/google/uuid"
)

type BookServiceImpl struct {
    repo ports.BookRepository
}

func NewBookService(repo ports.BookRepository) ports.BookService {
    return &BookServiceImpl{repo: repo}
}

func (s *BookServiceImpl) AddBook(book entity.Book) (*entity.BookResponse, error) {
    // ตรวจ ISBN ซ้ำก่อนเพิ่ม
    existing, _ := s.repo.FindByISBN(book.ISBN)
    if existing != nil {
        return nil, errors.New("ISBN นี้มีในระบบแล้ว")
    }

    book.Id = uuid.NewString()
    book.Available = book.TotalCopies // ตอนเพิ่งเพิ่ม available = total ทั้งหมด

    created, err := s.repo.Create(book)
    if err != nil {
        return nil, err
    }
    return toBookResponse(created), nil
}

func (s *BookServiceImpl) GetAllBooks() ([]entity.BookResponse, error) {
    books, err := s.repo.FindAll()
    if err != nil {
        return nil, err
    }

    var result []entity.BookResponse
    for _, b := range books {
        result = append(result, *toBookResponse(&b))
    }
    return result, nil
}

func (s *BookServiceImpl) GetBookById(id string) (*entity.BookResponse, error) {
    book, err := s.repo.FindById(id)
    if err != nil {
        return nil, errors.New("ไม่พบหนังสือ")
    }
    return toBookResponse(book), nil
}

func (s *BookServiceImpl) UpdateBook(id string, book entity.Book) (*entity.BookResponse, error) {
    existing, err := s.repo.FindById(id)
    if err != nil {
        return nil, errors.New("ไม่พบหนังสือ")
    }

    existing.Title  = book.Title
    existing.Author = book.Author

    updated, err := s.repo.Update(*existing)
    if err != nil {
        return nil, err
    }
    return toBookResponse(updated), nil
}

func (s *BookServiceImpl) DeleteBook(id string) error {
    _, err := s.repo.FindById(id)
    if err != nil {
        return errors.New("ไม่พบหนังสือ")
    }
    return s.repo.Delete(id)
}

// helper แปลง entity → response (ไม่ส่ง field ที่ไม่จำเป็นออก API)
func toBookResponse(b *entity.Book) *entity.BookResponse {
    return &entity.BookResponse{
        Id:        b.Id,
        Title:     b.Title,
        Author:    b.Author,
        ISBN:      b.ISBN,
        Available: b.Available,
    }
}