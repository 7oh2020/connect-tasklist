package dto

import "github.com/7oh2020/connect-tasklist/backend/app"

type ChangeTaskNameParams struct {
	id     IDParam
	userID IDParam
	name   string
}

func NewChangeTaskNameParams(id string, userID string, name string) *ChangeTaskNameParams {
	return &ChangeTaskNameParams{
		id:     *NewIDParam(id),
		userID: *NewIDParam(userID),
		name:   name,
	}
}

func (f *ChangeTaskNameParams) ID() string {
	return f.id.Value()
}

func (f *ChangeTaskNameParams) UserID() string {
	return f.userID.Value()
}

func (f *ChangeTaskNameParams) Name() string {
	return f.name
}

func (f *ChangeTaskNameParams) Validate() error {
	if err := f.id.Validate(); err != nil {
		return err
	}
	if err := f.userID.Validate(); err != nil {
		return err
	}
	if len([]rune(f.name)) > 100 {
		return &app.ErrInputValidationFailed{Msg: "name must be 100 characters or less"}
	}
	return nil
}
