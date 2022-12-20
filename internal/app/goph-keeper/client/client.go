package client

type KeeperClientConfig interface {
	GetAPIAddress() string
}

type KeeperClient interface {
	Login(user, password string) error
	Logout() error
	Register(user, password string) error
}

func NewClient() (KeeperClient, error) {
	return NewHTTPClient()
}
