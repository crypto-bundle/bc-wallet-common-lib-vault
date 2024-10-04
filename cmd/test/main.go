package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	commonVault "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault"
	commonVaultTokenClient "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault/client/token"
)

var ErrMockFormatter = errors.New("mock_err_formatter")

type errFmt struct {
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
	ctx, cancelCtxFunc := context.WithCancel(context.Background())

	vaultClientSrv, err := commonVaultTokenClient.NewClient(ctx, mockErrFmtSvc, vaultCfg)
	if err != nil {
		log.Fatal(err)
	}

	vaultSrv, err := commonVault.NewService(log.Default(), mockErrFmtSvc, vaultCfg, vaultClientSrv)
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
