package test

import (
    "fmt"
    "library/core/entity"
)

type MockBorrowRepo struct {
    Records []entity.BorrowRecord
}

func (m *MockBorrowRepo) Create(record entity.BorrowRecord) (*entity.BorrowRecord, error) {
    m.Records = append(m.Records, record)
    return &record, nil
}

func (m *MockBorrowRepo) FindById(id string) (*entity.BorrowRecord, error) {
    for _, r := range m.Records {
        if r.Id == id {
            return &r, nil
        }
    }
    return nil, fmt.Errorf("ไม่พบ")
}

func (m *MockBorrowRepo) FindActiveByMemberId(memberId string) ([]entity.BorrowRecord, error) {
    var result []entity.BorrowRecord
    for _, r := range m.Records {
        if r.MemberId == memberId && r.Status == "borrowed" {
            result = append(result, r)
        }
    }
    return result, nil
}

func (m *MockBorrowRepo) FindAllByMemberId(memberId string) ([]entity.BorrowRecord, error) {
    var result []entity.BorrowRecord
    for _, r := range m.Records {
        if r.MemberId == memberId {
            result = append(result, r)
        }
    }
    return result, nil
}

func (m *MockBorrowRepo) Update(record entity.BorrowRecord) (*entity.BorrowRecord, error) {
    for i, r := range m.Records {
        if r.Id == record.Id {
            m.Records[i] = record
            return &record, nil
        }
    }
    return nil, fmt.Errorf("ไม่พบ")
}

func (m *MockBorrowRepo) CountActiveByMemberId(memberId string) (int, error) {
    count := 0
    for _, r := range m.Records {
        if r.MemberId == memberId && r.Status == "borrowed" {
            count++
        }
    }
    return count, nil
}