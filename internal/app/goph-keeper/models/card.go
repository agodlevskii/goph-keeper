package models

type CardRequest struct {
	Name    string `json:"name"`
	Number  string `json:"number"`
	Holder  string `json:"holder"`
	ExpDate string `json:"exp_date"`
	CVV     string `json:"cvv"`
	Note    string `json:"note"`
}

type CardResponse struct {
	UID     string `json:"-"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Number  string `json:"number"`
	Holder  string `json:"holder"`
	ExpDate string `json:"exp_date"`
	CVV     string `json:"cvv"`
	Note    string `json:"note"`
}

func (c CardResponse) TableRow() []string {
	return []string{c.ID, c.Name, c.Number, c.Holder, c.ExpDate, c.CVV, c.Note}
}
