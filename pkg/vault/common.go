/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2025 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

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
