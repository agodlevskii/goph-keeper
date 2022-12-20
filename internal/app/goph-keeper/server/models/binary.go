package models

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
