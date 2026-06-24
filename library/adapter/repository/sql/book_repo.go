package sql

import (
    "library/core/entity"
    "library/core/ports"
    "time"

    "gorm.io/gorm"
)

// GORM model — มี tag และ gorm.Model สำหรับ auto timestamps
type BookModel struct {
    Id          string `gorm:"primaryKey"`
    Title       string `gorm:"not null"`
    Author      string `gorm:"not null"`
    ISBN        string `gorm:"uniqueIndex;not null"`
    TotalCopies int    `gorm:"not null"`
    Available   int    `gorm:"not null"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type bookRepository struct {
    db *gorm.DB
}

func NewBookRepository(db *gorm.DB) ports.BookRepository {
    return &bookRepository{db: db}
}

func (r *bookRepository) Create(book entity.Book) (*entity.Book, error) {
    model := toBookModel(book)
    result := r.db.Create(&model)
    if result.Error != nil {
        return nil, result.Error
    }
    return toBookEntity(model), nil
}

func (r *bookRepository) FindAll() ([]entity.Book, error) {
    var models []BookModel
    result := r.db.Find(&models)
    if result.Error != nil {
        return nil, result.Error
    }

    var books []entity.Book
    for _, m := range models {
        books = append(books, *toBookEntity(m))
    }
    return books, nil
}

func (r *bookRepository) FindById(id string) (*entity.Book, error) {
    var model BookModel
    result := r.db.First(&model, "id = ?", id)
    if result.Error != nil {
        return nil, result.Error
    }
    return toBookEntity(model), nil
}

func (r *bookRepository) FindByISBN(isbn string) (*entity.Book, error) {
    var model BookModel
    result := r.db.First(&model, "isbn = ?", isbn)
    if result.Error != nil {
        return nil, result.Error
    }
    return toBookEntity(model), nil
}

func (r *bookRepository) Update(book entity.Book) (*entity.Book, error) {
    model := toBookModel(book)
    result := r.db.Save(&model)
    if result.Error != nil {
        return nil, result.Error
    }
    return toBookEntity(model), nil
}

func (r *bookRepository) Delete(id string) error {
    return r.db.Delete(&BookModel{}, "id = ?", id).Error
}

// ── แปลงไปมา ──

func toBookModel(e entity.Book) BookModel {
    return BookModel{
        Id:          e.Id,
        Title:       e.Title,
        Author:      e.Author,
        ISBN:        e.ISBN,
        TotalCopies: e.TotalCopies,
        Available:   e.Available,
    }
}

func toBookEntity(m BookModel) *entity.Book {
    return &entity.Book{
        Id:          m.Id,
        Title:       m.Title,
        Author:      m.Author,
        ISBN:        m.ISBN,
        TotalCopies: m.TotalCopies,
        Available:   m.Available,
        CreatedAt:   m.CreatedAt,
        UpdatedAt:   m.UpdatedAt,
    }
}