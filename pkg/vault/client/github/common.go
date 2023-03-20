package github

type configService interface {
	GetAddress() string
	GetHost() string
	GetPort() uint32
	IsUseHTTPS() bool

	GetGithubAuthPath() string
	GetGithubAuthToken() string
}
