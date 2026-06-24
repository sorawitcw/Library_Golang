package entity

import "time"

type Member struct {
    Id        string
    Username  string
    Password  string
    Email     string
    Role      string    // "admin" หรือ "member"
    CreatedAt time.Time
    UpdatedAt time.Time
}

type MemberResponse struct {
    Id       string `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Role     string `json:"role"`
}