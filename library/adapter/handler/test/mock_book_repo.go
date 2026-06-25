package test

import (
	"fmt"
	"library/core/entity"
)

// MockBookRepo implement ports.BookRepository ทั้งหมด
// แต่ไม่ได้คุยกับ DB จริงเลย
type MockBookRepo struct {
    Books []entity.Book // เก็บข้อมูลใน memory
}

func (m *MockBookRepo) Create(book entity.Book) (*entity.Book, error) {
    m.Books = append(m.Books, book)
    return &book, nil
}

func (m *MockBookRepo) FindAll() ([]entity.Book, error) {
    return m.Books, nil
}

func (m *MockBookRepo) FindById(id string) (*entity.Book, error) {
    for _, b := range m.Books {
        if b.Id == id {
            return &b, nil
        }
    }
    return nil, fmt.Errorf("ไม่พบหนังสือ")
}

func (m *MockBookRepo) FindByISBN(isbn string) (*entity.Book, error) {
    for _, b := range m.Books {
        if b.ISBN == isbn {
            return &b, nil
        }
    }
    return nil, fmt.Errorf("ไม่พบ")
}

func (m *MockBookRepo) Update(book entity.Book) (*entity.Book, error) {
    for i, b := range m.Books {
        if b.Id == book.Id {
            m.Books[i] = book
            return &book, nil
        }
    }
    return nil, fmt.Errorf("ไม่พบ")
}

func (m *MockBookRepo) Delete(id string) error {
    for i, b := range m.Books {
        if b.Id == id {
            m.Books = append(m.Books[:i], m.Books[i+1:]...)
            return nil
        }
    }
    return fmt.Errorf("ไม่พบ")
}