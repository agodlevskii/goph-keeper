package binary

type Request struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
	Note string `json:"note"`
}

type Response struct {
	UID  string `json:"-"`
	ID   string `json:"id"`
	Name string `json:"name"`
	Data []byte `json:"data"`
	Note string `json:"note"`
}
