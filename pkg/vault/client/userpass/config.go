package userpass

type AuthConfig struct {
	Username string `envconfig:"VAULT_USERNAME" default:""`
	Password string `envconfig:"VAULT_PASSWORD" default:""`
}

func (c *AuthConfig) GetUsername() string {
	return c.Username
}

func (c *AuthConfig) GePassword() string {
	return c.Password
}

func (c *AuthConfig) Prepare() error {
	return nil
}

func (c *AuthConfig) PrepareWith(dependentCfgList ...interface{}) error {
	return nil
}
