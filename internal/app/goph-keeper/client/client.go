package client

type KeeperClient interface {
	Login(user, password string)
}

func NewClient() (KeeperClient, error) {
	return NewHTTPClient()
}
