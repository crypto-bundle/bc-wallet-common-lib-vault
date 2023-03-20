package vault

import (
	"context"
	vaultApi "github.com/hashicorp/vault/api"
)

type configService interface {
	GetAddress() string
	GetHost() string
	GetPort() uint32
	IsUseHTTPS() bool
	GetAuthMethod() string
	GetDataPath() string
	GetTransitKey() string
	// --------------------------------------------
	// Token auth config methods
	// --------------------------------------------

	//GetAuthToken() string
	// --------------------------------------------
	// Github token auth config methods
	// --------------------------------------------

	//GetGithubAuthPath() string
	//GetGithubAuthToken() string
	// --------------------------------------------
	// User and password auth config methods

	//GetUserName() string
	//GetUserPassword() string

	// --------------------------------------------
	// Kubernates config methods
	// --------------------------------------------

	//GetKubernatesAppRole() string
	//GetKubernatesSATokenPath() string
	//GetKubernatesAuthPath() string
}

type clientService interface {
	GetClient() *vaultApi.Client
	Login(context.Context) (*vaultApi.Client, error)
}
