package entity

import "time"

type Book struct {
    Id          string    `json:"id"`
    Title       string    `json:"title"`
    Author      string    `json:"author"`
    ISBN        string    `json:"isbn"`
    TotalCopies int       `json:"total_copies"`
    Available   int       `json:"available"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type BookResponse struct {
    Id        string `json:"id"`
    Title     string `json:"title"`
    Author    string `json:"author"`
    ISBN      string `json:"isbn"`
    Available int    `json:"available"`
}