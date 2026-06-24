package sql

import (
    "fmt"
    "os"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func NewPostgresDB() (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, fmt.Errorf("เชื่อมต่อ DB ไม่ได้: %w", err)
    }

    // AutoMigrate สร้าง table อัตโนมัติตาม model
    err = db.AutoMigrate(
        &BookModel{},
        &MemberModel{},
        &BorrowModel{},
    )
    if err != nil {
        return nil, fmt.Errorf("migrate ไม่ได้: %w", err)
    }

    return db, nil
}