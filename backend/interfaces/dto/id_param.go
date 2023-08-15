package dto

import "github.com/7oh2020/connect-tasklist/backend/app"

type IDParam struct {
	id string
}

func NewIDParam(id string) *IDParam {
	return &IDParam{id}
}

func (i *IDParam) Value() string {
	return i.id
}

func (i *IDParam) Validate() error {
	if len(i.id) > 50 {
		return &app.ErrInputValidationFailed{Msg: "id must be 50 characters or less"}
	}
	return nil
}
