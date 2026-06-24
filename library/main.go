package main

import (
    "library/adapter/handler"
    repoSQL "library/adapter/repository/sql"
    "library/core/middleware"
    "library/core/service"
    "log"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
    // โหลดไฟล์ .env
    if err := godotenv.Load(); err != nil {
        log.Println("Warning: .env file not found, using system environment variables")
    }

    // 1. เชื่อม DB
    db, err := repoSQL.NewPostgresDB()
    if err != nil {
        log.Fatalf("DB error: %v", err)
    }

    // 2. สร้าง Repository
    bookRepo   := repoSQL.NewBookRepository(db)
    memberRepo := repoSQL.NewMemberRepository(db)
    borrowRepo := repoSQL.NewBorrowRepository(db)

    // 3. สร้าง Service (ส่ง repo เข้าไป)
    bookSvc   := service.NewBookService(bookRepo)
    memberSvc := service.NewMemberService(memberRepo)
    borrowSvc := service.NewBorrowService(bookRepo, borrowRepo)

    // 4. สร้าง Handler (ส่ง service เข้าไป)
    bookH   := handler.NewBookHandler(bookSvc)
    memberH := handler.NewMemberHandler(memberSvc)
    borrowH := handler.NewBorrowHandler(borrowSvc)

    // 5. ลงทะเบียน Route
    r := gin.Default()

    // Public routes — ไม่ต้อง token
    r.POST("/register", memberH.Register)
    r.POST("/login",    memberH.Login)

    // Protected routes — ต้องมี JWT token
    auth := r.Group("/", middleware.JWTMiddleware())
    {
        // Books
        auth.GET("/books",          bookH.GetAllBooks)
        auth.GET("/books/:id",      bookH.GetBookById)
        auth.POST("/books",         bookH.AddBook)
        auth.PUT("/books/:id",      bookH.UpdateBook)
        auth.DELETE("/books/:id",   bookH.DeleteBook)

        // Members
        auth.GET("/profile",        memberH.GetProfile)
        auth.GET("/members",        memberH.GetAllMembers) // admin only (ตรวจใน service)

        // Borrow
        auth.POST("/borrow",        borrowH.BorrowBook)
        auth.PUT("/borrow/:id/return", borrowH.ReturnBook)
        auth.GET("/my-borrows",     borrowH.GetMyBorrows)
    }

    r.Run(":8080")
}