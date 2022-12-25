package models

type TextRequest struct {
	Name string `json:"name"`
	Data string `json:"data"`
	Note string `json:"note"`
}

type TextResponse struct {
	UID  string `json:"-"`
	ID   string `json:"id"`
	Name string `json:"name"`
	Data string `json:"data"`
	Note string `json:"note"`
}

func (t TextResponse) TableRow() []string {
	return []string{t.ID, t.Name, t.Data, t.Note}
}
