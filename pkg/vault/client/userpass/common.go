package userpass

type configService interface {
	GetAddress() string
	GetHost() string
	GetPort() uint32
	IsUseHTTPS() bool
	GetAuthToken() string

	GetUserName() string
	GetUserPassword() string
}
