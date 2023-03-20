package kubernates

type AuthConfig struct {
	KubeAppRole     string `envconfig:"VAULT_KUBE_APP_ROLE" default:""`
	KubeSATokenPath string `envconfig:"VAULT_KUBE_SA_TOKEN_PATH" default:""`
}

func (c *AuthConfig) GetAppRole() string {
	return c.KubeAppRole
}

func (c *AuthConfig) GetKubernatesSATokenPath() string {
	return c.KubeSATokenPath
}

func (c *AuthConfig) Prepare() error {
	return nil
}

func (c *AuthConfig) PrepareWith(dependentCfgList ...interface{}) error {
	return nil
}
