package token

import (
	"errors"
	"os"
	"strings"
)

var (
	ErrMissedVaultTokenData = errors.New("missed vault token data")
)

type AuthConfig struct {
	AuthToken         string `envconfig:"VAULT_AUTH_TOKEN" default:""`
	AuthTokenFilePath string `envconfig:"VAULT_AUTH_TOKEN_FILE_PATH" default:"/vault/secrets/token"`
	// config dependencies
	e errorFormatterService
}

func (c *AuthConfig) GetAuthToken() string {
	return c.AuthToken
}

func (c *AuthConfig) Prepare() error {
	if c.AuthToken != "" {
		return nil
	}

	if c.AuthTokenFilePath == "" {
		return c.e.ErrorOnly(ErrMissedVaultTokenData)
	}

	fileContent, err := os.ReadFile(c.AuthTokenFilePath)
	if err != nil {
		return c.e.ErrorOnly(err)
	}

	c.AuthToken = strings.TrimRight(string(fileContent), "\n")

	return nil
}

func (c *AuthConfig) PrepareWith(dependentCfgList ...interface{}) error {
	for _, dependency := range dependentCfgList {
		switch castedDep := dependency.(type) {
		case errorFormatterService:
			c.e = castedDep
		default:
			continue
		}
	}

	return nil
}
