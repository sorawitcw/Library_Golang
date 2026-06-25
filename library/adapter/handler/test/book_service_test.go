package test

import (
    "library/core/entity"
    "library/core/service"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestAddBook_Success(t *testing.T) {
    // Arrange — เตรียม mock และ service
    repo := &MockBookRepo{}
    svc  := service.NewBookService(repo)

    book := entity.Book{
        Title:       "Clean Code",
        Author:      "Robert Martin",
        ISBN:        "978-0132350884",
        TotalCopies: 3,
    }

    // Act — เรียก method ที่ต้องการทดสอบ
    result, err := svc.AddBook(book)

    // Assert — ตรวจผลลัพธ์
    assert.NoError(t, err)
    assert.Equal(t, "Clean Code", result.Title)
    assert.Equal(t, 3, result.Available)
}

func TestAddBook_DuplicateISBN(t *testing.T) {
    repo := &MockBookRepo{
        // ใส่หนังสือที่มี ISBN นี้อยู่แล้ว
        Books: []entity.Book{
            {Id: "1", ISBN: "978-0132350884"},
        },
    }
    svc := service.NewBookService(repo)

    book := entity.Book{
        Title:       "หนังสืออื่น",
        ISBN:        "978-0132350884", // ISBN ซ้ำ
        TotalCopies: 1,
    }

    result, err := svc.AddBook(book)

    // ต้อง error และ result ต้อง nil
    assert.Error(t, err)
    assert.Nil(t, result)
    assert.Equal(t, "ISBN นี้มีในระบบแล้ว", err.Error())
}