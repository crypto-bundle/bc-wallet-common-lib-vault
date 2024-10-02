package vault

import (
	"errors"
)

const (
	ErrVaultAPIClientInitDetail  = "unable to initialize vault client"
	ErrK8sAuthInitDetail         = "unable to initialize kubernetes auth method"
	ErrUserpassInitDetail        = "unable to initialize userpass auth method"
	ErrK8sLoginDetail            = "unable to log in with kubernetes auth"
	ErrUserpassLoginDetail       = "unable to log in with user and password"
	ErrNotExistingAuthInfoDetail = "no auth info was returned after login"
	ErrReadSecretDetail          = "unable to read secret"
	ErrEmptySecretDetail         = "unable to get secret"
	ErrCastSecretDetail          = "secret casting error"
	ErrNotExistingKeyDetail      = "missed key in secret"
	ErrKeyTypeDetail             = "unexpected key type in secret"
	ErrTransitSecretFormatDetail = "unexpected format for transit secret"
	ErrInitTokenTTLWatcherDetail = "unable to initialize auth token lifetime watcher"
	ErrUnableGetTokenInfoDetail  = "unable to get token info"
	ErrEmptyConfigDetail         = "config is empty"
)

var (
	ErrConfigIsNotValid    = errors.New("config is not valid.some fields are missing")
	ErrReadSecret          = errors.New(ErrReadSecretDetail)
	ErrEmptySecret         = errors.New(ErrEmptySecretDetail)
	ErrCastSecret          = errors.New(ErrCastSecretDetail)
	ErrNotExistingKey      = errors.New(ErrNotExistingKeyDetail)
	ErrKeyType             = errors.New(ErrKeyTypeDetail)
	ErrUnableGetTokenInfo  = errors.New(ErrUnableGetTokenInfoDetail)
	ErrTransitSecretFormat = errors.New(ErrTransitSecretFormatDetail)
)
