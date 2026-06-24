package entity

import "time"

type BorrowRecord struct {
    Id         string
    MemberId   string
    BookId     string
    BorrowedAt time.Time
    DueDate    time.Time     // กำหนดคืน (7 วัน)
    ReturnedAt *time.Time    // nil = ยังไม่คืน
    Fine       float64       // ค่าปรับ (ถ้ามี)
    Status     string        // "borrowed" / "returned" / "overdue"
}

type BorrowResponse struct {
    Id         string     `json:"id"`
    BookId     string     `json:"book_id"`
    BorrowedAt time.Time  `json:"borrowed_at"`
    DueDate    time.Time  `json:"due_date"`
    ReturnedAt *time.Time `json:"returned_at"`
    Fine       float64    `json:"fine"`
    Status     string     `json:"status"`
}