package clock

import "time"

// 時刻を操作する
type IClockManager interface {
	// 現在時刻を取得する
	GetNow() time.Time
}

type ClockManager struct{}

func NewClockManager() *ClockManager {
	return &ClockManager{}
}

func (m *ClockManager) GetNow() time.Time {
	return time.Now().UTC()
}
