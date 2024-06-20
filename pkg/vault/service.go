package vault

import (
	"context"
	"encoding/json"
	"log"
	"time"

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

type Service struct {
	logger *log.Logger

	client    *vaultApi.Client
	clientSvc clientService
	authInfo  *vaultApi.Secret
	cfg       configService

	loadedSecrets map[string]string
}

func (s *Service) IsHealed(_ context.Context) bool {
	status, err := s.client.Sys().Health()
	if err != nil {
		return false
	}

	serverOk := status.Standby && status.Sealed
	if !serverOk {
		return false
	}

	secretData, err := s.client.Auth().Token().LookupSelf()
	if err != nil {
		return false
	}

	currentTime := time.Now()
	if currentTime.Unix() > int64(secretData.LeaseDuration) {
		return false
	}

	return true
}

// GetCredentialsBytes returns all secrets bytes from default path.
func (s *Service) GetCredentialsBytes() (b []byte, err error) {
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
func (s *Service) GetCredentialsBytesByPath(path string) (b []byte, err error) {
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
func (s *Service) GetCredentialsByPathAndKey(path, key string) (string, error) {
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
func (s *Service) GetCredentialsByPathAndKeys(path string, keys ...string) (map[string]string, error) {
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
		keyVal, isExists := data[k]
		if !isExists {
			return res, ErrNotExistingKey.WithMsg(k)
		}

		keyString, isExists := keyVal.(string)
		if !isExists {
			return res, ErrKeyType.WithMsg(k)
		}

		res[k] = keyString
	}

	return res, nil
}

func (s *Service) Login(ctx context.Context) (*vaultApi.Client, error) {
	loggedInVaultClient, err := s.clientSvc.Login(ctx)
	if err != nil {
		return nil, err
	}

	s.client = loggedInVaultClient

	err = s.prepareAndStartRenew(ctx)
	if err != nil {
		return nil, err
	}

	return s.client, nil
}

func (s *Service) GetClient() *vaultApi.Client {
	return s.client
}

func NewService(
	logger *log.Logger,
	cfg configService,
	client clientService,
) (*Service, error) {
	return &Service{
		logger: logger,

		clientSvc: client,
		cfg:       cfg,

		loadedSecrets: make(map[string]string),
	}, nil
}
