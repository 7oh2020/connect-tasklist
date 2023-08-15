package dto

type UserInfo struct {
	id    string
	email string
	token string
}

func NewUserInfo(id string, email string, token string) *UserInfo {
	return &UserInfo{id, email, token}
}

func (i *UserInfo) ID() string {
	return i.id
}

func (i *UserInfo) Email() string {
	return i.email
}

func (i *UserInfo) Token() string {
	return i.token
}
