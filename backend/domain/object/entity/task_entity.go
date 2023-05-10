package entity

import (
	"errors"
	"time"
)

type Task struct {
	ID          string
	UserID      string
	Name        string
	IsCompleted bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewTask(id string, userID string, name string, now time.Time) (*Task, error) {
	task := &Task{
		ID:          id,
		UserID:      userID,
		Name:        name,
		IsCompleted: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := task.Validate(); err != nil {
		return nil, err
	}
	return task, nil
}

// 自身のフィールドをバリデーションする
func (t *Task) Validate() error {
	if t.ID == "" || t.UserID == "" || t.Name == "" {
		return errors.New("error: validation failed")
	}
	return nil
}
