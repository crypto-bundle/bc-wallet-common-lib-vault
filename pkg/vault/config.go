package vault

import "fmt"

type BaseConfig struct {
	Host     string `envconfig:"VAULT_HOST" default:"vault"`
	Port     uint32 `envconfig:"VAULT_PORT" default:"8200"`
	UseHTTPS bool   `envconfig:"VAULT_USE_HTTPS" default:"true"`

	DataPath   string `envconfig:"VAULT_DATA_PATH" default:""`
	TransitKey string `envconfig:"VAULT_TRANSIT_KEY" default:""`
}

func (c *BaseConfig) GetAddress() string {
	protocol := "http"
	if c.UseHTTPS {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s,%d", protocol, c.Host, c.Port)
}

func (c *BaseConfig) GetHost() string {
	return c.Host
}

func (c *BaseConfig) GetPort() uint32 {
	return c.Port
}

func (c *BaseConfig) IsUseHTTPS() bool {
	return c.UseHTTPS
}

func (c *BaseConfig) GeDataPath() string {
	return c.DataPath
}

func (c *BaseConfig) GeTransitKey() string {
	return c.TransitKey
}

func (c *BaseConfig) Prepare() error {
	return nil
}

func (c *BaseConfig) PrepareWith(dependentCfgList ...interface{}) error {
	return nil
}
