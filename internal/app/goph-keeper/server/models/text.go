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
