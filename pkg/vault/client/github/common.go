package github

import (
	"context"

	vaultApi "github.com/hashicorp/vault/api"
)

type selfService interface {
	GetAuthMethod() string
	GetClient() *vaultApi.Client
	Login(ctx context.Context) (*vaultApi.Client, error)
}

type configService interface {
	GetAddress() string
	GetHost() string
	GetPort() uint32
	IsUseHTTPS() bool

	GetGithubAuthPath() string
	GetGithubAuthToken() string
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
