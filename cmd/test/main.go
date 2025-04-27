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

package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	commonVault "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault"
	commonVaultTokenClient "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault/client/token"
)

var ErrMockFormatter = errors.New("mock_err_formatter")

type loggBuilder struct{}

func (f *loggBuilder) NewSlogLoggerEntry(fields ...any) *slog.Logger {
	return slog.Default()
}

func (f *loggBuilder) NewSlogNamedLoggerEntry(named string, fields ...any) *slog.Logger {
	return slog.Default()
}

func (f *loggBuilder) NewSlogLoggerEntryWithFields(fields ...slog.Attr) *slog.Logger {
	return slog.Default()
}

type errFmt struct {
}

func (f *errFmt) NewErrorWithCode(text string, code int) error {
	return ErrMockFormatter
}

func (f *errFmt) ErrorCodeIsOneOf(err error, codes ...int) (int, bool) {
	return -1, false
}

func (f *errFmt) ErrCodeIsOneOf(err error, codes ...int) (int, bool) {
	return -1, false
}

func (f *errFmt) ErrorWithCode(_ error, _ int) error {
	return ErrMockFormatter
}

func (f *errFmt) ErrWithCode(_ error, _ int) error {
	return ErrMockFormatter
}

func (f *errFmt) ErrorGetCode(_ error) int {
	return -1
}

func (f *errFmt) ErrGetCode(_ error) int {
	return -1
}

func (f *errFmt) ErrorNoWrap(_ error) error {
	return ErrMockFormatter
}

func (f *errFmt) ErrNoWrap(_ error) error {
	return ErrMockFormatter
}

func (f *errFmt) ErrorOnly(_ error, _ ...string) error {
	return ErrMockFormatter
}

func (f *errFmt) Error(_ error, _ ...string) error {
	return ErrMockFormatter
}

func (f *errFmt) Errorf(_ error, _ string, _ ...interface{}) error {
	return ErrMockFormatter
}

func (f *errFmt) NewError(_ ...string) error {
	return ErrMockFormatter
}

func (f *errFmt) NewErrorf(_ string, _ ...interface{}) error {
	return ErrMockFormatter
}

func main() {
	const (
		DefaultTokenRenewTTL        = 960
		DefaultVaultConnectionPOrt  = 8200
		DefaultOsSignalsChannelSize = 2
	)

	type VaultWrappedConfig struct {
		*commonVault.BaseConfig
		*commonVaultTokenClient.AuthConfig
	}

	vaultCfg := &VaultWrappedConfig{
		BaseConfig: &commonVault.BaseConfig{
			Host:          "127.0.0.1",
			Port:          DefaultVaultConnectionPOrt,
			UseHTTPS:      false,
			AuthMethod:    "token",
			TokenRenewTTL: DefaultTokenRenewTTL,
			DataPath:      "kv/data/crypto-bundle/bc-wallet-common/transit,kv/data/crypto-bundle/bc-wallet-ethereum-hdwallet/common",
		},
		AuthConfig: &commonVaultTokenClient.AuthConfig{
			AuthToken:         "",
			AuthTokenFilePath: "./vault.token",
		},
	}

	err := vaultCfg.AuthConfig.Prepare()
	if err != nil {
		log.Fatal(err)
	}

	mockErrFmtSvc := &errFmt{}
	mockLoggBuilder := &loggBuilder{}
	ctx, cancelCtxFunc := context.WithCancel(context.Background())

	vaultClientSrv, err := commonVaultTokenClient.NewClient(ctx, mockErrFmtSvc, vaultCfg)
	if err != nil {
		log.Fatal(err)
	}

	vaultSrv, err := commonVault.NewService(mockLoggBuilder, mockErrFmtSvc, vaultCfg, vaultClientSrv)
	if err != nil {
		log.Fatal(err)
	}

	_, err = vaultSrv.Login(ctx)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, DefaultOsSignalsChannelSize)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	cancelCtxFunc()

	log.Printf("stopped")
}
