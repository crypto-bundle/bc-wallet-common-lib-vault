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

type errorFormatterService interface {
	ErrorWithCode(err error, code int) error
	ErrWithCode(err error, code int) error
	ErrorGetCode(err error) int
	ErrGetCode(err error) int
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
