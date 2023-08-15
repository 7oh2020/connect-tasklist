package dto

import "github.com/7oh2020/connect-tasklist/backend/app"

type CreateTaskParams struct {
	userID IDParam
	name   string
}

func NewCreateTaskParams(userID string, name string) *CreateTaskParams {
	return &CreateTaskParams{
		userID: *NewIDParam(userID),
		name:   name,
	}
}

func (f *CreateTaskParams) UserID() string {
	return f.userID.Value()
}

func (f *CreateTaskParams) Name() string {
	return f.name
}

func (f *CreateTaskParams) Validate() error {
	if err := f.userID.Validate(); err != nil {
		return err
	}
	if len([]rune(f.name)) > 100 {
		return &app.ErrInputValidationFailed{Msg: "name must be 100 characters or less"}
	}
	return nil
}
