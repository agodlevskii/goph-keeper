package auth

type Payload struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
