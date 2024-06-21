package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	commonVault "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault"
	commonVaultTokenClient "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault/client/token"
)

func main() {
	type VaultWrappedConfig struct {
		*commonVault.BaseConfig
		*commonVaultTokenClient.AuthConfig
	}

	vaultCfg := &VaultWrappedConfig{
		BaseConfig: &commonVault.BaseConfig{
			Host:          "127.0.0.1",
			Port:          8200,
			UseHTTPS:      false,
			AuthMethod:    "token",
			TokenRenewTTL: 960,
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

	ctx, cancelCtxFunc := context.WithCancel(context.Background())

	vaultClientSrv, err := commonVaultTokenClient.NewClient(ctx, vaultCfg)
	if err != nil {
		log.Fatal(err)
	}

	vaultSrv, err := commonVault.NewService(log.Default(), vaultCfg, vaultClientSrv)
	if err != nil {
		log.Fatal(err)
	}

	_, err = vaultSrv.Login(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//err = vaultSrv.LoadSecrets(ctx)
	//if err != nil {
	//	log.Fatal(err)
	//}

	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	cancelCtxFunc()

	log.Printf("stopped")
}
