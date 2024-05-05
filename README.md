# bc-wallet-common-lib-vault

## Description

Library for manage Hashicord Vault config, connections, and secret vault data. 
Library supports github, k8s, userpass and token authentication. 

Library contains:
* common Vault config structs for all supported auth method
* secret data loading from multiple vault_data path
* encryption/decryption

## Usage example

Examples of create connection to Vault, retrieving data and encryption usage

### Config and connection

```go
package main

import (
	"context"
	"errors"
	"log"

	commonEnvConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/config"
	commonVault "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault"
	commonVaultTokenClient "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault/client/token"
)

type VaultWrappedConfig struct {
	*commonVault.BaseConfig
	*commonVaultTokenClient.AuthConfig
}

func main() {
	ctx := context.Background()

	cfgPreparerSrv := commonEnvConfig.NewConfigManager()
	vaultCfg := &VaultWrappedConfig{
		BaseConfig: &commonVault.BaseConfig{},
		AuthConfig: &commonVaultTokenClient.AuthConfig{},
	}
	err := cfgPreparerSrv.PrepareTo(vaultCfg).With(baseCfgSrv).Do(ctx)
	if err != nil {
		panic(err)
	}

	vaultClientSvc, err := commonVaultTokenClient.NewClient(ctx, vaultCfg)
	if err != nil {
		panic(err)
	}

	// vault prepare 
	vaultSvc, err := commonVault.NewService(ctx, vaultCfg, vaultClientSvc)
	if err != nil {
		panic(err)
	}

	_, err = vaultSvc.Login(ctx)
	if err != nil {
		panic(err)
	}

	err = vaultSvc.LoadSecrets(ctx)
	if err != nil {
		panic(err)
	}

	data, isExists := vaultSvc.GetByName("key_name")
	if !isExists {
		panic(errors.New("missing value"))
	}

	log.Printf("result: %s", data)

	...
}
```

### Data encryption

```go
package main

import (
	"context"
	"log"

	commonEnvConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/config"
	commonVault "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault"
	commonVaultTokenClient "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault/client/token"
)

type VaultWrappedConfig struct {
	*commonVault.BaseConfig
	*commonVaultTokenClient.AuthConfig
}

func main() {
	ctx := context.Background()

	cfgPreparerSrv := commonEnvConfig.NewConfigManager()
	vaultCfg := &VaultWrappedConfig{
		BaseConfig: &commonVault.BaseConfig{},
		AuthConfig: &commonVaultTokenClient.AuthConfig{},
	}
	err := cfgPreparerSrv.PrepareTo(vaultCfg).With(baseCfgSrv).Do(ctx)
	if err != nil {
		panic(err)
	}

	vaultClientSvc, err := commonVaultTokenClient.NewClient(ctx, vaultCfg)
	if err != nil {
		panic(err)
	}

	// vault prepare 
	vaultSvc, err := commonVault.NewService(ctx, vaultCfg, vaultClientSvc)
	if err != nil {
		panic(err)
	}

	_, err = vaultSvc.Login(ctx)
	if err != nil {
		panic(err)
	}

	vaultCrypterSvc, err := commonVault.NewEncryptService(ctx, vaultSvc.GetClient())
	if err != nil {
		panic(err)
	}

	encryptedData, err := vaultCrypterSvc.Encrypt([]byte("Hello world"))
	if err != nil {
		panic(err)
	}
	
	log.Printf("result: %s", encryptedData)

	...
}
```

## Contributors

* Author and maintainer - [@gudron (Alex V Kotelnikov)](https://github.com/gudron)

## Licence

**bc-wallet-common-lib-vault** is licensed under the [MIT](./LICENSE) License.