package vault

import "fmt"

type BaseConfig struct {
	Host       string `envconfig:"VAULT_SERVICE_HOST" default:"vault"`
	Port       uint32 `envconfig:"VAULT_SERVICE_PORT" default:"8200"`
	UseHTTPS   bool   `envconfig:"VAULT_USE_HTTPS" default:"true"`
	AuthMethod string `envconfig:"VAULT_AUTH_METHOD" default:"token"`

	TokenRenewTTL int `envconfig:"VAULT_TOKEN_RENEW_TTL" default:"240"`

	DataPath string `envconfig:"VAULT_APP_DATA_PATH" default:""`

	// dependencies
	baseAppCfgSrv baseApplicationConfigService
}

func (c *BaseConfig) GetAddress() string {
	protocol := "http"
	if c.UseHTTPS {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s:%d", protocol, c.Host, c.Port)
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

func (c *BaseConfig) GetDataPath() string {
	return c.DataPath
}

func (c *BaseConfig) GetAuthMethod() string {
	return c.AuthMethod
}

func (c *BaseConfig) GetTokenRenewTTL() int {
	return c.TokenRenewTTL
}

func (c *BaseConfig) GetApplicationStageName() string {
	return c.baseAppCfgSrv.GetStageName()
}

func (c *BaseConfig) GetApplicationEnvironment() string {
	return c.baseAppCfgSrv.GetEnvironmentName()
}

func (c *BaseConfig) Prepare() error {
	return nil
}

func (c *BaseConfig) PrepareWith(dependentCfgList ...interface{}) error {
	for _, cfgSrv := range dependentCfgList {
		switch castedCfg := cfgSrv.(type) {
		case baseApplicationConfigService:
			c.baseAppCfgSrv = castedCfg
		default:
			continue
		}
	}

	return nil
}
