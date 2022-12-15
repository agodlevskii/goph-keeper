package password

type Request struct {
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Note     string `json:"note"`
}

type Response struct {
	UID      string `json:"-"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Note     string `json:"note"`
}
