package ports

import "library/core/entity"

// BorrowRepository — สิ่งที่ Repository ต้องทำกับ DB
type BorrowRepository interface {
    Create(record entity.BorrowRecord) (*entity.BorrowRecord, error)
    FindById(id string) (*entity.BorrowRecord, error)
    FindActiveByMemberId(memberId string) ([]entity.BorrowRecord, error)
    FindAllByMemberId(memberId string) ([]entity.BorrowRecord, error)
    Update(record entity.BorrowRecord) (*entity.BorrowRecord, error)
    CountActiveByMemberId(memberId string) (int, error)
}

// BorrowService — สิ่งที่ Handler จะเรียกใช้
type BorrowService interface {
    BorrowBook(memberId, bookId string) (*entity.BorrowResponse, error)
    ReturnBook(recordId, memberId string) (*entity.BorrowResponse, error)
    GetMyBorrows(memberId string) ([]entity.BorrowResponse, error)
}