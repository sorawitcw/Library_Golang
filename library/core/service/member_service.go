package service

import (
    "errors"
    "library/core/entity"
    "library/core/ports"
    "time"

    "github.com/dgrijalva/jwt-go"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("library-secret-key")

type MemberServiceImpl struct {
    repo ports.MemberRepository
}

func NewMemberService(repo ports.MemberRepository) ports.MemberService {
    return &MemberServiceImpl{repo: repo}
}

func (s *MemberServiceImpl) Register(member entity.Member) (*entity.MemberResponse, error) {
    // ตรวจ username ซ้ำ
    existing, _ := s.repo.FindByUsername(member.Username)
    if existing != nil {
        return nil, errors.New("username นี้ถูกใช้ไปแล้ว")
    }

    // hash password ก่อนเก็บ — ไม่เก็บ plain text เด็ดขาด
    hashed, err := bcrypt.GenerateFromPassword([]byte(member.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    member.Id       = uuid.NewString()
    member.Password = string(hashed)
    member.Role     = "member" // default role

    created, err := s.repo.Create(member)
    if err != nil {
        return nil, err
    }
    return toMemberResponse(created), nil
}

func (s *MemberServiceImpl) Login(username, password string) (string, error) {
    member, err := s.repo.FindByUsername(username)
    if err != nil {
        return "", errors.New("username หรือ password ไม่ถูกต้อง")
    }

    // เปรียบ password กับ hash ที่เก็บไว้
    err = bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(password))
    if err != nil {
        return "", errors.New("username หรือ password ไม่ถูกต้อง")
    }

    // สร้าง JWT token พร้อม claim
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "memberId": member.Id,
        "role":     member.Role,
        "exp":      time.Now().Add(24 * time.Hour).Unix(),
    })

    return token.SignedString(jwtSecret)
}

func (s *MemberServiceImpl) GetProfile(memberId string) (*entity.MemberResponse, error) {
    member, err := s.repo.FindById(memberId)
    if err != nil {
        return nil, errors.New("ไม่พบ member")
    }
    return toMemberResponse(member), nil
}

func (s *MemberServiceImpl) GetAllMembers(requesterRole string) ([]entity.MemberResponse, error) {
    // RBAC — เฉพาะ admin เท่านั้น
    if requesterRole != "admin" {
        return nil, errors.New("ไม่มีสิทธิ์เข้าถึง")
    }

    members, err := s.repo.FindAll()
    if err != nil {
        return nil, err
    }

    var result []entity.MemberResponse
    for _, m := range members {
        result = append(result, *toMemberResponse(&m))
    }
    return result, nil
}

func toMemberResponse(m *entity.Member) *entity.MemberResponse {
    return &entity.MemberResponse{
        Id:       m.Id,
        Username: m.Username,
        Email:    m.Email,
        Role:     m.Role,
    }
}