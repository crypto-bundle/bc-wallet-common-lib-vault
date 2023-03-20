package vault

import (
	"context"
	"encoding/json"
	vaultApi "github.com/hashicorp/vault/api"
)

// Values for vault authentication method used by this library.
// Remember the environment values are case-sensitive.
const (
	authMethodGithub     = "github"
	authMethodKubernetes = "kubernetes"
	authMethodUserpass   = "userpass"
	authMethodToken      = "token"
	authMethodNone       = "none" // no vault is used, just get values from env variables

	defaultAuthMethod = authMethodKubernetes
)

type Vaulter interface {
	Encrypt(toEncrypt []byte) ([]byte, error)
	Decrypt(cipherBytes []byte) ([]byte, error)

	GetCredentialsBytes() (b []byte, err error)
	GetCredentialsBytesByPath(path string) (b []byte, err error)
	GetCredentialsByPathAndKey(path, field string) (string, error)
	GetCredentialsByPathAndKeys(path string, fields ...string) (map[string]string, error)
}

type service struct {
	client    *vaultApi.Client
	clientSvc clientService
	authInfo  *vaultApi.Secret
	cfg       configService
}

// GetCredentialsBytes returns all secrets bytes from default path.
func (s *service) GetCredentialsBytes() (b []byte, err error) {
	secret, err := s.client.Logical().Read(s.cfg.GetDataPath())
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

func (s *service) Login(ctx context.Context) (*vaultApi.Client, error) {
	loggedInVaultClient, err := s.clientSvc.Login(ctx)
	if err != nil {
		return nil, err
	}

	s.client = loggedInVaultClient

	return s.client, nil
}

func NewService(ctx context.Context,
	cfg configService,
	client clientService,
) (*service, error) {
	//var clientSvc clientService = nil
	//switch cfg.GetAuthMethod() {
	//case authMethodGithub:
	//	svc, err := github.NewClient(ctx, cfg)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	clientSvc = svc
	//
	//case authMethodKubernetes:
	//	svc, err := kubernates.NewClient(ctx, cfg)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	clientSvc = svc
	//
	//case authMethodUserpass:
	//	svc, err := userpass.NewClient(ctx, cfg)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	clientSvc = svc
	//default:
	//	return nil, fmt.Errorf("unknown auth method: %s", cfg.GetAuthMethod())
	//}

	//loggedInClient, err := clientSvc.Login(ctx)
	//if err != nil {
	//	return nil, err
	//}

	return &service{
		clientSvc: client,
		cfg:       cfg,
	}, nil
}
