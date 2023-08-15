package entity

import (
	"time"

	"github.com/7oh2020/connect-tasklist/backend/domain"
	"github.com/7oh2020/connect-tasklist/backend/domain/object/value"
)

type Task struct {
	ID          *value.ID
	UserID      *value.ID
	Name        string
	IsCompleted bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// フィールドの妥当性を検証する
func (t *Task) Validate() error {
	if err := t.ID.Validate(); err != nil {
		return err
	}
	if err := t.UserID.Validate(); err != nil {
		return err
	}
	if t.Name == "" {
		return &domain.ErrValidationFailed{Msg: "name is empty"}
	}
	return nil
}
