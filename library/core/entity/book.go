package entity

import "time"

type Book struct {
    Id          string
    Title       string
    Author      string
    ISBN        string
    TotalCopies int
    Available   int    // จำนวนที่ยืมได้ตอนนี้
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type BookResponse struct {
    Id        string `json:"id"`
    Title     string `json:"title"`
    Author    string `json:"author"`
    ISBN      string `json:"isbn"`
    Available int    `json:"available"`
}