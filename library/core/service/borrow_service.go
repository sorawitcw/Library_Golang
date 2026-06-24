package service

import (
    "errors"
    "library/core/entity"
    "library/core/ports"
    "time"

    "github.com/google/uuid"
)

const (
    MaxBorrowLimit  = 3           // ยืมได้สูงสุด 3 เล่ม
    BorrowDays      = 7           // กำหนดคืน 7 วัน
    FinePerDay      = 5.0         // ค่าปรับวันละ 5 บาท
)

type BorrowServiceImpl struct {
    bookRepo   ports.BookRepository
    borrowRepo ports.BorrowRepository
}

func NewBorrowService(
    bookRepo ports.BookRepository,
    borrowRepo ports.BorrowRepository,
) ports.BorrowService {
    return &BorrowServiceImpl{
        bookRepo:   bookRepo,
        borrowRepo: borrowRepo,
    }
}

func (s *BorrowServiceImpl) BorrowBook(memberId, bookId string) (*entity.BorrowResponse, error) {
    // Rule 1: ยืมได้ไม่เกิน 3 เล่มในเวลาเดียวกัน
    count, err := s.borrowRepo.CountActiveByMemberId(memberId)
    if err != nil {
        return nil, err
    }
    if count >= MaxBorrowLimit {
        return nil, errors.New("ยืมได้สูงสุด 3 เล่มเท่านั้น")
    }

    // Rule 2: หนังสือต้องมีอยู่และยังมีให้ยืม
    book, err := s.bookRepo.FindById(bookId)
    if err != nil {
        return nil, errors.New("ไม่พบหนังสือ")
    }
    if book.Available <= 0 {
        return nil, errors.New("หนังสือถูกยืมหมดแล้ว")
    }

    // ตัดจำนวน available ของหนังสือ
    book.Available--
    if _, err = s.bookRepo.Update(*book); err != nil {
        return nil, err
    }

    // สร้าง borrow record
    now := time.Now()
    record := entity.BorrowRecord{
        Id:         uuid.NewString(),
        MemberId:   memberId,
        BookId:     bookId,
        BorrowedAt: now,
        DueDate:    now.AddDate(0, 0, BorrowDays),
        Status:     "borrowed",
    }

    created, err := s.borrowRepo.Create(record)
    if err != nil {
        return nil, err
    }
    return toBorrowResponse(created), nil
}

func (s *BorrowServiceImpl) ReturnBook(recordId, memberId string) (*entity.BorrowResponse, error) {
    record, err := s.borrowRepo.FindById(recordId)
    if err != nil {
        return nil, errors.New("ไม่พบรายการยืม")
    }

    // ตรวจว่าเป็นของ member คนนี้จริง
    if record.MemberId != memberId {
        return nil, errors.New("ไม่มีสิทธิ์คืนรายการนี้")
    }

    // ตรวจว่ายังไม่ได้คืนไปแล้ว
    if record.Status == "returned" {
        return nil, errors.New("หนังสือเล่มนี้คืนไปแล้ว")
    }

    now := time.Now()

    // Rule 3: คำนวณค่าปรับถ้าคืนเกินกำหนด
    var fine float64
    if now.After(record.DueDate) {
        overdueDays := int(now.Sub(record.DueDate).Hours() / 24)
        fine = float64(overdueDays) * FinePerDay
    }

    // อัปเดต record
    record.ReturnedAt = &now
    record.Fine       = fine
    record.Status     = "returned"

    updated, err := s.borrowRepo.Update(*record)
    if err != nil {
        return nil, err
    }

    // คืน available ของหนังสือ
    book, err := s.bookRepo.FindById(record.BookId)
    if err == nil {
        book.Available++
        s.bookRepo.Update(*book)
    }

    return toBorrowResponse(updated), nil
}

func (s *BorrowServiceImpl) GetMyBorrows(memberId string) ([]entity.BorrowResponse, error) {
    records, err := s.borrowRepo.FindAllByMemberId(memberId)
    if err != nil {
        return nil, err
    }

    var result []entity.BorrowResponse
    for _, r := range records {
        result = append(result, *toBorrowResponse(&r))
    }
    return result, nil
}

func toBorrowResponse(r *entity.BorrowRecord) *entity.BorrowResponse {
    return &entity.BorrowResponse{
        Id:         r.Id,
        BookId:     r.BookId,
        BorrowedAt: r.BorrowedAt,
        DueDate:    r.DueDate,
        ReturnedAt: r.ReturnedAt,
        Fine:       r.Fine,
        Status:     r.Status,
    }
}