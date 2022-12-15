package card

type Request struct {
	Name    string `json:"name"`
	Number  string `json:"number"`
	Holder  string `json:"holder"`
	ExpDate string `json:"exp_date"`
	CVV     string `json:"cvv"`
	Note    string `json:"note"`
}

type Response struct {
	UID     string `json:"-"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Number  string `json:"number"`
	Holder  string `json:"holder"`
	ExpDate string `json:"exp_date"`
	CVV     string `json:"cvv"`
	Note    string `json:"note"`
}
