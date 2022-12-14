package client

type KeeperClientConfig interface {
	GetAPIAddress() string
}

type KeeperClient interface {
	Login(user, password string)
}

func NewClient() (KeeperClient, error) {
	return NewHTTPClient()
}
