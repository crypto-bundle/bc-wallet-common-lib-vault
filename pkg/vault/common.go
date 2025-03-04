package vault

import (
	"context"
	"log/slog"

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

	// GetAuthToken() string
	// --------------------------------------------
	// Github token auth config methods
	// --------------------------------------------

	// GetGithubAuthPath() string
	// GetGithubAuthToken() string
	// --------------------------------------------
	// User and password auth config methods

	// GetUserName() string
	// GetUserPassword() string

	// --------------------------------------------
	// Kubernates config methods
	// --------------------------------------------

	// GetKubernatesAppRole() string
	// GetKubernatesSATokenPath() string
	// GetKubernatesAuthPath() string

	GetApplicationStageName() string
	GetApplicationEnvironment() string
}

type clientService interface {
	GetAuthMethod() string
	GetClient() *vaultApi.Client
	Login(ctx context.Context) (*vaultApi.Client, error)
}

type renewerService interface {
	IsHealed(_ context.Context) bool
	PrepareAndStartRenew(ctx context.Context) error
}

type baseApplicationConfigService interface {
	GetEnvironmentName() string
	GetStageName() string
}

type Encryptor interface {
	Encrypt(toEncrypt []byte) ([]byte, error)
	Decrypt(cipherBytes []byte) ([]byte, error)
}

type Vaulter interface {
	GetCredentialsBytes() (b []byte, err error)
	GetCredentialsBytesByPath(path string) (b []byte, err error)
	GetCredentialsByPathAndKey(path, field string) (string, error)
	GetCredentialsByPathAndKeys(path string, fields ...string) (map[string]string, error)
}

type loggerFabricService interface {
	NewSlogLoggerEntry(fields ...any) *slog.Logger
	NewSlogNamedLoggerEntry(named string, fields ...any) *slog.Logger
	NewSlogLoggerEntryWithFields(fields ...slog.Attr) *slog.Logger
}

type errorFormatterService interface {
	ErrWithCode(err error, code int) error
	NewErrorWithCode(text string, code int) error
	ErrorGetCode(err error) int
	ErrGetCode(err error) int
	ErrorCodeIsOneOf(err error, codes ...int) (int, bool)
	ErrCodeIsOneOf(err error, codes ...int) (int, bool)
	// ErrorNoWrap function for pseudo-wrap error, must be used in case of linter warnings...
	ErrorNoWrap(err error) error
	// ErrNoWrap same with ErrorNoWrap function, just alias for ErrorNoWrap, just short function name...
	ErrNoWrap(err error) error
	ErrorOnly(err error, details ...string) error
	Error(err error, details ...string) error
	Errorf(err error, format string, args ...interface{}) error
	NewError(details ...string) error
	NewErrorf(format string, args ...interface{}) error
}
