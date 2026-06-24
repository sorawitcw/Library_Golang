package sql

import (
    "library/core/entity"
    "library/core/ports"
    "time"

    "gorm.io/gorm"
)

type MemberModel struct {
    Id        string `gorm:"primaryKey"`
    Username  string `gorm:"uniqueIndex;not null"`
    Password  string `gorm:"not null"`
    Email     string `gorm:"uniqueIndex;not null"`
    Role      string `gorm:"not null;default:member"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

type memberRepository struct {
    db *gorm.DB
}

func NewMemberRepository(db *gorm.DB) ports.MemberRepository {
    return &memberRepository{db: db}
}

func (r *memberRepository) Create(member entity.Member) (*entity.Member, error) {
    model := toMemberModel(member)
    result := r.db.Create(&model)
    if result.Error != nil {
        return nil, result.Error
    }
    return toMemberEntity(model), nil
}

func (r *memberRepository) FindById(id string) (*entity.Member, error) {
    var model MemberModel
    result := r.db.First(&model, "id = ?", id)
    if result.Error != nil {
        return nil, result.Error
    }
    return toMemberEntity(model), nil
}

func (r *memberRepository) FindByUsername(username string) (*entity.Member, error) {
    var model MemberModel
    result := r.db.First(&model, "username = ?", username)
    if result.Error != nil {
        return nil, result.Error
    }
    return toMemberEntity(model), nil
}

func (r *memberRepository) FindAll() ([]entity.Member, error) {
    var models []MemberModel
    result := r.db.Find(&models)
    if result.Error != nil {
        return nil, result.Error
    }

    var members []entity.Member
    for _, m := range models {
        members = append(members, *toMemberEntity(m))
    }
    return members, nil
}

func toMemberModel(e entity.Member) MemberModel {
    return MemberModel{
        Id:       e.Id,
        Username: e.Username,
        Password: e.Password,
        Email:    e.Email,
        Role:     e.Role,
    }
}

func toMemberEntity(m MemberModel) *entity.Member {
    return &entity.Member{
        Id:        m.Id,
        Username:  m.Username,
        Password:  m.Password,
        Email:     m.Email,
        Role:      m.Role,
        CreatedAt: m.CreatedAt,
        UpdatedAt: m.UpdatedAt,
    }
}