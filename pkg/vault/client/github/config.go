package github

type AuthConfig struct {
	AuthPath  string `envconfig:"VAULT_GITHUB_AUTH_PATH" default:""`
	AuthToken string `envconfig:"VAULT_GITHUB_AUTH_TOKEN" default:""`
}

func (c *AuthConfig) GetGithubAuthPath() string {
	return c.AuthPath
}

func (c *AuthConfig) GetGithubAuthToken() string {
	return c.AuthToken
}

func (c *AuthConfig) Prepare() error {
	if c.AuthPath == "" {
		c.AuthPath = githubDefaultAuthPath
	}

	return nil
}

func (c *AuthConfig) PrepareWith(dependentCfgList ...interface{}) error {
	return nil
}
