package sqlc

import (
	"testing"

	"github.com/7oh2020/connect-tasklist/backend/domain/repository"
)

func TestTaskRepository_NewTaskRepository(tt *testing.T) {
	tt.Run("異常系: structがinterfaceを実装しているか", func(t *testing.T) {
		var _ repository.ITaskRepository = (*SQLCTaskRepository)(nil)
	})
}
