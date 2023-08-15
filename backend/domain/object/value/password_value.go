package value

import "github.com/7oh2020/connect-tasklist/backend/domain"

type Password struct {
	value string
}

func NewPassword(value string) *Password {
	return &Password{value}
}

func (p *Password) Value() string {
	return p.value
}

func (p *Password) Validate() error {
	if p.value == "" {
		return &domain.ErrValidationFailed{Msg: "password is empty"}
	}
	return nil
}
