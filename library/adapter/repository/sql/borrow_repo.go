package sql

import (
    "library/core/entity"
    "library/core/ports"
    "time"

    "gorm.io/gorm"
)

type BorrowModel struct {
    Id         string     `gorm:"primaryKey"`
    MemberId   string     `gorm:"not null;index"`
    BookId     string     `gorm:"not null"`
    BorrowedAt time.Time
    DueDate    time.Time
    ReturnedAt *time.Time // pointer = nullable
    Fine       float64    `gorm:"default:0"`
    Status     string     `gorm:"not null;default:borrowed"`
}

type borrowRepository struct {
    db *gorm.DB
}

func NewBorrowRepository(db *gorm.DB) ports.BorrowRepository {
    return &borrowRepository{db: db}
}

func (r *borrowRepository) Create(record entity.BorrowRecord) (*entity.BorrowRecord, error) {
    model := toBorrowModel(record)
    result := r.db.Create(&model)
    if result.Error != nil {
        return nil, result.Error
    }
    return toBorrowEntity(model), nil
}

func (r *borrowRepository) FindById(id string) (*entity.BorrowRecord, error) {
    var model BorrowModel
    result := r.db.First(&model, "id = ?", id)
    if result.Error != nil {
        return nil, result.Error
    }
    return toBorrowEntity(model), nil
}

func (r *borrowRepository) FindActiveByMemberId(memberId string) ([]entity.BorrowRecord, error) {
    var models []BorrowModel
    result := r.db.Where("member_id = ? AND status = ?", memberId, "borrowed").Find(&models)
    if result.Error != nil {
        return nil, result.Error
    }

    var records []entity.BorrowRecord
    for _, m := range models {
        records = append(records, *toBorrowEntity(m))
    }
    return records, nil
}

func (r *borrowRepository) FindAllByMemberId(memberId string) ([]entity.BorrowRecord, error) {
    var models []BorrowModel
    // เรียงจากล่าสุดก่อน
    result := r.db.Where("member_id = ?", memberId).
        Order("borrowed_at DESC").
        Find(&models)
    if result.Error != nil {
        return nil, result.Error
    }

    var records []entity.BorrowRecord
    for _, m := range models {
        records = append(records, *toBorrowEntity(m))
    }
    return records, nil
}

func (r *borrowRepository) Update(record entity.BorrowRecord) (*entity.BorrowRecord, error) {
    model := toBorrowModel(record)
    // ใช้ Save เพื่อ update ทุก field รวมถึง zero value
    result := r.db.Save(&model)
    if result.Error != nil {
        return nil, result.Error
    }
    return toBorrowEntity(model), nil
}

func (r *borrowRepository) CountActiveByMemberId(memberId string) (int, error) {
    var count int64
    result := r.db.Model(&BorrowModel{}).
        Where("member_id = ? AND status = ?", memberId, "borrowed").
        Count(&count)
    return int(count), result.Error
}

func toBorrowModel(e entity.BorrowRecord) BorrowModel {
    return BorrowModel{
        Id:         e.Id,
        MemberId:   e.MemberId,
        BookId:     e.BookId,
        BorrowedAt: e.BorrowedAt,
        DueDate:    e.DueDate,
        ReturnedAt: e.ReturnedAt,
        Fine:       e.Fine,
        Status:     e.Status,
    }
}

func toBorrowEntity(m BorrowModel) *entity.BorrowRecord {
    return &entity.BorrowRecord{
        Id:         m.Id,
        MemberId:   m.MemberId,
        BookId:     m.BookId,
        BorrowedAt: m.BorrowedAt,
        DueDate:    m.DueDate,
        ReturnedAt: m.ReturnedAt,
        Fine:       m.Fine,
        Status:     m.Status,
    }
}