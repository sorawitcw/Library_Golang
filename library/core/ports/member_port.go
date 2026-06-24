package ports

import "library/core/entity"

// MemberRepository — สิ่งที่ Repository ต้องทำกับ DB
type MemberRepository interface {
    Create(member entity.Member) (*entity.Member, error)
    FindById(id string) (*entity.Member, error)
    FindByUsername(username string) (*entity.Member, error)
    FindAll() ([]entity.Member, error)
}

// MemberService — สิ่งที่ Handler จะเรียกใช้
type MemberService interface {
    Register(member entity.Member) (*entity.MemberResponse, error)
    Login(username, password string) (string, error) // คืน JWT token
    GetProfile(memberId string) (*entity.MemberResponse, error)
    GetAllMembers(requesterRole string) ([]entity.MemberResponse, error)
}