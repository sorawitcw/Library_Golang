package test

import (
    "library/core/entity"
    "library/core/service"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func TestBorrowBook_Success(t *testing.T) {
    bookRepo := &MockBookRepo{
        Books: []entity.Book{
            {Id: "book-1", Title: "Clean Code", Available: 2},
        },
    }
    borrowRepo := &MockBorrowRepo{}
    svc := service.NewBorrowService(bookRepo, borrowRepo)

    result, err := svc.BorrowBook("member-1", "book-1")

    assert.NoError(t, err)
    assert.Equal(t, "borrowed", result.Status)
    assert.Equal(t, "book-1", result.BookId)

    // ตรวจว่า available ถูกตัดแล้ว
    assert.Equal(t, 1, bookRepo.Books[0].Available)
}

func TestBorrowBook_ExceedLimit(t *testing.T) {
    bookRepo := &MockBookRepo{
        Books: []entity.Book{
            {Id: "book-4", Available: 1},
        },
    }
    // member นี้ยืมครบ 3 เล่มแล้ว
    borrowRepo := &MockBorrowRepo{
        Records: []entity.BorrowRecord{
            {Id: "r1", MemberId: "member-1", Status: "borrowed"},
            {Id: "r2", MemberId: "member-1", Status: "borrowed"},
            {Id: "r3", MemberId: "member-1", Status: "borrowed"},
        },
    }
    svc := service.NewBorrowService(bookRepo, borrowRepo)

    result, err := svc.BorrowBook("member-1", "book-4")

    assert.Error(t, err)
    assert.Nil(t, result)
    assert.Equal(t, "ยืมได้สูงสุด 3 เล่มเท่านั้น", err.Error())
}

func TestBorrowBook_NoAvailable(t *testing.T) {
    bookRepo := &MockBookRepo{
        Books: []entity.Book{
            {Id: "book-1", Available: 0}, // หมดแล้ว
        },
    }
    borrowRepo := &MockBorrowRepo{}
    svc := service.NewBorrowService(bookRepo, borrowRepo)

    result, err := svc.BorrowBook("member-1", "book-1")

    assert.Error(t, err)
    assert.Nil(t, result)
    assert.Equal(t, "หนังสือถูกยืมหมดแล้ว", err.Error())
}

func TestReturnBook_WithFine(t *testing.T) {
    bookRepo := &MockBookRepo{
        Books: []entity.Book{
            {Id: "book-1", Available: 0},
        },
    }
    // due date เป็นอดีต 3 วัน → ค่าปรับ 15 บาท
    overdue := time.Now().AddDate(0, 0, -3)
    borrowRepo := &MockBorrowRepo{
        Records: []entity.BorrowRecord{
            {
                Id:       "record-1",
                MemberId: "member-1",
                BookId:   "book-1",
                DueDate:  overdue,
                Status:   "borrowed",
            },
        },
    }
    svc := service.NewBorrowService(bookRepo, borrowRepo)

    result, err := svc.ReturnBook("record-1", "member-1")

    assert.NoError(t, err)
    assert.Equal(t, "returned", result.Status)
    assert.Equal(t, 15.0, result.Fine) // 3 วัน × 5 บาท
}