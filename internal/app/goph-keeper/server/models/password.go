package models

type PasswordRequest struct {
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Note     string `json:"note"`
}

type PasswordResponse struct {
	UID      string `json:"-"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Note     string `json:"note"`
}
