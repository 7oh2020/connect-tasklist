package value

import (
	"github.com/7oh2020/connect-tasklist/backend/domain"
)

type ID struct {
	value string
}

func NewID(value string) *ID {
	return &ID{value}
}

func (i *ID) Value() string {
	return i.value
}

func (i *ID) Validate() error {
	if i.value == "" {
		return &domain.ErrValidationFailed{Msg: "id is empty"}
	}
	return nil
}

func (i *ID) Equal(value string) bool {
	return i.value == value
}
