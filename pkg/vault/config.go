package vault

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Address         string `envconfig:"VAULT_ADDRESS" default:""`
	AppRole         string `envconfig:"VAULT_APP_ROLE" default:""`
	AuthPath        string `envconfig:"VAULT_AUTH_PATH" default:""`
	KubeSATokenPath string `envconfig:"VAULT_KUBE_SA_TOKEN_PATH" default:""`
	Username        string `envconfig:"VAULT_USERNAME" default:""`
	Password        string `envconfig:"VAULT_PASSWORD" default:""`
	DataPath        string `envconfig:"VAULT_DATA_PATH" default:""`
	TransitKey      string `envconfig:"VAULT_TRANSIT_KEY" default:""`
}

func (c *Config) GetAddress() string {
	return c.Address
}

func (c *Config) GetAppRole() string {
	return c.AppRole
}

func (c *Config) GetAuthPath() string {
	return c.AuthPath
}

func (c *Config) GetKubernatesSATokenPath() string {
	return c.KubeSATokenPath
}

func (c *Config) GetUsername() string {
	return c.Username
}

func (c *Config) GePassword() string {
	return c.Password
}

func (c *Config) GeDataPath() string {
	return c.TransitKey
}

func (c *Config) GeTransitKey() string {
	return c.TransitKey
}

// IsEmpty returns true if some prop doesn't initialized except for transit key
func (c *Config) IsEmpty() bool {
	return len(c.KubeSATokenPath) == 0 ||
		len(c.AuthPath) == 0 ||
		len(c.AppRole) == 0 ||
		len(c.DataPath) == 0 ||
		len(c.Address) == 0
}

// Prepare variables to static configuration
func (c *Config) Prepare() error {
	err := envconfig.Process("", c)
	if err != nil {
		return err
	}

	var configIsNotValid = len(c.KubeSATokenPath) == 0 ||
		len(c.AuthPath) == 0 ||
		len(c.AppRole) == 0 ||
		len(c.DataPath) == 0 ||
		len(c.Address) == 0

	if configIsNotValid {
		return ErrConfigIsNotValid
	}

	return nil
}
