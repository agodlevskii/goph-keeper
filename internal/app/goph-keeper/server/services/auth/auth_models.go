package auth

type Request struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
