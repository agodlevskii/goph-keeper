package models

import "fmt"

type BinaryRequest struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
	Note string `json:"note"`
}

type BinaryResponse struct {
	UID  string `json:"-"`
	ID   string `json:"id"`
	Name string `json:"name"`
	Data []byte `json:"data"`
	Note string `json:"note"`
}

func (b BinaryResponse) TableRow() []string {
	var data string
	if len(b.Data) > 0 {
		data = fmt.Sprintf("%b", b.Data)
	}
	return []string{b.ID, b.Name, data, b.Note}
}
