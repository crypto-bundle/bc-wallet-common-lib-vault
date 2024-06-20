package vault

import (
	"errors"
	"fmt"
)

var (
	ErrConfigIsNotValid = errors.New("config is not valid.some fields are missing")
)

const (
	ErrVaultAPIClientInit  InternalError = "unable to initialize vault client"
	ErrK8sAuthInit         InternalError = "unable to initialize kubernetes auth method"
	ErrUserpassInit        InternalError = "unable to initialize userpass auth method"
	ErrK8sLogin            InternalError = "unable to log in with kubernetes auth"
	ErrUserpassLogin       InternalError = "unable to log in with user and password"
	ErrNotExistingAuthInfo InternalError = "no auth info was returned after login"
	ErrReadSecret          InternalError = "unable to read secret"
	ErrEmptySecret         InternalError = "unable to get secret"
	ErrCastSecret          InternalError = "secret casting error"
	ErrNotExistingKey      InternalError = "missed key in secret"
	ErrKeyType             InternalError = "unexpected key type in secret"
	ErrTransitSecretFormat InternalError = "unexpected format for transit secret"
	ErrInitTokenTTLWatcher InternalError = "unable to initialize auth token lifetime watcher"
	ErrUnableGetTokenInfo  InternalError = "unable to get token info"
	ErrEmptyConfig         InternalError = "config is empty"
)

type InternalError string

func (e InternalError) Error() string {
	return string(e)
}

// WithMsg returns internal error with additional text.
func (e InternalError) WithMsg(msg string) error {
	return fmt.Errorf("%w: %s", e, msg)
}

// NewInternalError returns new error.
func NewInternalError(vErr InternalError, err error) error {
	return fmt.Errorf("%w: %q", vErr, err)
}
