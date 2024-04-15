# Change Log

## [v0.0.1] - 08.03.2022 00:49 MSK
### Added 
* Added vault client
* Added vault client config

## [v0.0.2] - 09.03.2022 18:41 MSK
### Changed
* Moved from common dir to root of lib-vault directory

## [v0.0.3] - 07.02.2023 23:11 MSK
### Changed
* Lib-vault moved to another repository - https://github.com/crypto-bundle/bc-wallet-common-lib-vault
* Added MIT license

## [v0.0.4] - 08.02.2023 10:25 MSK
### Changed
* Added config init flow via envconfig

## [v0.0.5] - 20.02.2023 00:28 MSK
### Changed
* Added config init flow via lib-config

## [v0.0.6] - 21.03.2023 23:20 MSK
### Changed
* Encrypt/Decrypt helpers
* Secrets loading flow
* Add vault auth
  * github
  * token
  * userpass
  * kubernates

## [v0.0.7] - 09.02.2024 15:41 MSK
### Added
* Added Vault helm-chart for local development. Chart cloned from [official Vault repository](https://github.com/hashicorp/vault-helm)
### Changed
* Changed config struct for support init with dependent base-config(from lib-config repo) service-component
* Changed go-namespace
* Changed README.md file-content - LICENSE section changed

## [v0.0.8 - v0.0.10] - 10.04.2024
### Changed
* Changed vault service-component - removed data encryption flow
  * Moved encryption flow to another service-component
  * Changed ENV variables name
* Changed README.md file-content - added "Data encryption" example block

## [v0.0.11] - 13.04.2024
### Added
* Added support of healthcheck flow, which required by [lib-healthcheck](https://github.com/crypto-bundle/bc-wallet-common-lib-healthcheck)

## [v0.0.12] - 15.04.2024
### Changed
* Changed env variables name for supporting kubernetes naming standard:
  * VAULT_HOST -> VAULT_SERVICE_HOST
  * VAULT_PORT -> VAULT_SERVICE_PORT

