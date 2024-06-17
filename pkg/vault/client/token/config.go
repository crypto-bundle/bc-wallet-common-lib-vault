package token

import (
	"errors"
	"os"
)

var (
	ErrMissedVaultTokenData = errors.New("missed vault token data")
)

type AuthConfig struct {
	AuthToken         string `envconfig:"VAULT_AUTH_TOKEN" default:""`
	AuthTokenFilePath string `envconfig:"VAULT_AUTH_TOKEN_FILE_PATH" default:"/vault/secrets/token"`
}

func (c *AuthConfig) GetAuthToken() string {
	return c.AuthToken
}

func (c *AuthConfig) Prepare() error {
	if c.AuthToken != "" {
		return nil
	}

	if c.AuthTokenFilePath == "" {
		return ErrMissedVaultTokenData
	}

	fileContent, err := os.ReadFile(c.AuthTokenFilePath)
	if err != nil {
		return err
	}

	c.AuthToken = string(fileContent)

	return nil
}

func (c *AuthConfig) PrepareWith(dependentCfgList ...interface{}) error {
	return nil
}
