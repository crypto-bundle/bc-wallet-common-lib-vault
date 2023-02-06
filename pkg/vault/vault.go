package vault

import (
	"context"
	"encoding/json"

	vaultApi "github.com/hashicorp/vault/api"
	k8sAuth "github.com/hashicorp/vault/api/auth/kubernetes"
	userpassAuth "github.com/hashicorp/vault/api/auth/userpass"
)

const githubAuthPath = "auth/github/login"

type Vaulter interface {
	Encrypt(toEncrypt []byte) ([]byte, error)
	Decrypt(cipherBytes []byte) ([]byte, error)

	GetCredentialsBytes() (b []byte, err error)
	GetCredentialsBytesByPath(path string) (b []byte, err error)
	GetCredentialsByPathAndKey(path, field string) (string, error)
	GetCredentialsByPathAndKeys(path string, fields ...string) (map[string]string, error)
}

type service struct {
	client   *vaultApi.Client
	authInfo *vaultApi.Secret
	cfg      *Config
}

// GetCredentialsBytes returns all secrets bytes from default path.
func (s *service) GetCredentialsBytes() (b []byte, err error) {
	secret, err := s.client.Logical().Read(s.cfg.DataPath)
	if err != nil {
		return nil, NewInternalError(ErrReadSecret, err)
	}
	if secret == nil {
		return nil, ErrEmptySecret
	}

	return json.Marshal(secret.Data["data"])
}

// GetCredentialsBytesByPath returns all secrets bytes from the specified path.
func (s *service) GetCredentialsBytesByPath(path string) (b []byte, err error) {
	secret, err := s.client.Logical().Read(path)
	if err != nil {
		return nil, NewInternalError(ErrReadSecret, err)
	}
	if secret == nil {
		return nil, ErrEmptySecret
	}

	return json.Marshal(secret.Data["data"])
}

// GetCredentialsByPathAndKey returns secret by path and field.
func (s *service) GetCredentialsByPathAndKey(path, key string) (string, error) {
	secret, err := s.client.Logical().Read(path)
	if err != nil {
		return "", NewInternalError(ErrReadSecret, err)
	}
	if secret == nil {
		return "", ErrEmptySecret
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return "", ErrCastSecret
	}

	keyVal, ok := data[key]
	if !ok {
		return "", ErrNotExistingKey.WithMsg(key)
	}

	keyString, ok := keyVal.(string)
	if !ok {
		return "", ErrKeyType.WithMsg(key)
	}

	return keyString, nil
}

// GetCredentialsByPathAndKeys returns fields sets from path.
func (s *service) GetCredentialsByPathAndKeys(path string, keys ...string) (map[string]string, error) {
	res := make(map[string]string, len(keys))

	secret, err := s.client.Logical().Read(path)
	if err != nil {
		return nil, NewInternalError(ErrReadSecret, err)
	}
	if secret == nil {
		return nil, ErrEmptySecret
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, ErrCastSecret
	}

	for _, k := range keys {
		keyVal, ok := data[k]
		if !ok {
			return res, ErrNotExistingKey.WithMsg(k)
		}

		keyString, ok := keyVal.(string)
		if !ok {
			return res, ErrKeyType.WithMsg(k)
		}

		res[k] = keyString
	}

	return res, nil
}

func (s *service) login(ctx context.Context) error {
	auth, err := k8sAuth.NewKubernetesAuth(
		s.cfg.AppRole,
		k8sAuth.WithServiceAccountTokenPath(s.cfg.KubeSATokenPath),
		k8sAuth.WithMountPath(s.cfg.AuthPath),
	)
	if err != nil {
		return NewInternalError(ErrK8sAuthInit, err)
	}

	authInfo, err := s.client.Auth().Login(ctx, auth)
	if err != nil {
		return NewInternalError(ErrK8sLogin, err)
	}
	if authInfo == nil {
		return ErrNotExistingAuthInfo
	}

	s.setToken(authInfo)

	return nil
}

func (s *service) setToken(auth *vaultApi.Secret) {
	s.authInfo = auth
	s.client.SetToken(auth.Auth.ClientToken)
}

// NewClientByGithubToken initialize vault client with github_token authorization.
func NewClientByGithubToken(cfg *Config, ghToken string) (Vaulter, error) {
	clientOpts := vaultApi.DefaultConfig()
	clientOpts.Address = cfg.Address

	client, err := vaultApi.NewClient(clientOpts)
	if err != nil {
		return nil, NewInternalError(ErrVaultAPIClientInit, err)
	}

	secret, err := client.Logical().Write(githubAuthPath, map[string]interface{}{"token": ghToken})
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, ErrEmptySecret
	}

	client.SetToken(secret.Auth.ClientToken)

	return &service{
		client: client,
		cfg:    cfg,
	}, nil
}

// NewClientByUserPass initialize vault client with user and password
// authentication.
func NewClientByUserPass(ctx context.Context, cfg *Config) (Vaulter, error) {
	clientOpts := vaultApi.DefaultConfig()
	clientOpts.Address = cfg.Address
	client, err := vaultApi.NewClient(clientOpts)
	if err != nil {
		return nil, NewInternalError(ErrVaultAPIClientInit, err)
	}

	auth, err := userpassAuth.NewUserpassAuth(cfg.Username, &userpassAuth.Password{FromString: cfg.Password})
	if err != nil {
		return nil, NewInternalError(ErrUserpassInit, err)
	}

	secret, err := client.Auth().Login(ctx, auth)
	if err != nil {
		return nil, NewInternalError(ErrUserpassLogin, err)
	}
	if secret == nil {
		return nil, ErrEmptySecret
	}

	client.SetToken(secret.Auth.ClientToken)

	return &service{
		client: client,
		cfg:    cfg,
	}, nil
}

// NewClient initialize vault client with service account token authorization.
func NewClient(ctx context.Context, cfg *Config) (Vaulter, error) {
	// to have the fallback way for application config initialization only from envs.
	if cfg.IsEmpty() {
		return nil, ErrEmptyConfig
	}

	clientOpts := vaultApi.DefaultConfig()
	clientOpts.Address = cfg.Address

	client, err := vaultApi.NewClient(clientOpts)
	if err != nil {
		return nil, NewInternalError(ErrVaultAPIClientInit, err)
	}

	vaultSvc := &service{
		client: client,
		cfg:    cfg,
	}

	err = vaultSvc.login(ctx)
	if err != nil {
		return nil, err
	}

	go vaultSvc.tokenRenew(ctx)

	return vaultSvc, nil
}
