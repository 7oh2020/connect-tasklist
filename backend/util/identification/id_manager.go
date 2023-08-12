package identification

import "github.com/google/uuid"

// IDの操作
type IIDManager interface {
	// IDの生成
	GenerateID() string
}

type UUIDManager struct{}

func NewUUIDManager() *UUIDManager {
	return &UUIDManager{}
}

func (m *UUIDManager) GenerateID() string {
	return uuid.NewString()
}
