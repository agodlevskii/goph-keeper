package data

type StorageType int

const (
	SBinary StorageType = iota
	SCard
	SPassword
	SText
)

type SecureData struct {
	UID  string      `json:"-"`
	ID   string      `json:"id"`
	Data []byte      `json:"data"`
	Type StorageType `json:"-"`
}
