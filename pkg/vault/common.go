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
	GetTokenRenewTTL() int
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

	GetApplicationStageName() string
	GetApplicationEnvironment() string
}

type clientService interface {
	GetClient() *vaultApi.Client
	Login(context.Context) (*vaultApi.Client, error)
}

type baseApplicationConfigService interface {
	GetEnvironmentName() string
	GetStageName() string
}

type Vaulter interface {
	Encrypt(toEncrypt []byte) ([]byte, error)
	Decrypt(cipherBytes []byte) ([]byte, error)

	GetCredentialsBytes() (b []byte, err error)
	GetCredentialsBytesByPath(path string) (b []byte, err error)
	GetCredentialsByPathAndKey(path, field string) (string, error)
	GetCredentialsByPathAndKeys(path string, fields ...string) (map[string]string, error)
}
