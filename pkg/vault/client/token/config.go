package token

type AuthConfig struct {
	AuthToken string `envconfig:"VAULT_AUTH_TOKEN" default:""`
}

func (c *AuthConfig) GetAuthToken() string {
	return c.AuthToken
}

func (c *AuthConfig) Prepare() error {
	return nil
}

func (c *AuthConfig) PrepareWith(dependentCfgList ...interface{}) error {
	return nil
}
