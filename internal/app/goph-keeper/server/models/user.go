package models

type UserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID       string `json:"-"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
