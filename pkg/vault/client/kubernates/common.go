package kubernates

type configService interface {
	GetAddress() string
	GetHost() string
	GetPort() uint32
	IsUseHTTPS() bool

	// GetAppRole ...
	GetKubernatesAppRole() string
	GetKubernatesAuthPath() string
	GetKubernatesSATokenPath() string
}
