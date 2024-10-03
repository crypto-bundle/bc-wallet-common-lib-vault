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
	l *log.Logger
	e errorFormatterService

	client     *vaultApi.Client
	clientSvc  clientService
	renewerSvc renewerService
	authInfo   *vaultApi.Secret
	cfg        configService

	loadedSecrets map[string]string
}

func (s *Service) IsHealed(ctx context.Context) bool {
	status, err := s.client.Sys().Health()
	if err != nil {
		return false
	}

	serverOk := status.Standby && status.Sealed
	if !serverOk {
		return false
	}

	isHealed := s.renewerSvc.IsHealed(ctx)
	if !isHealed {
		return isHealed
	}

	secretData, err := s.client.Auth().Token().LookupSelf()
	if err != nil {
		return false
	}

	currentTime := time.Now()

	return currentTime.Unix() > int64(secretData.LeaseDuration)
}

// GetCredentialsBytes returns all secrets bytes from default path.
func (s *Service) GetCredentialsBytes() ([]byte, error) {
	secret, err := s.client.Logical().Read(s.cfg.GetDataPath())
	if err != nil {
		return nil, s.e.ErrorOnly(err, ErrReadSecretDetail)
	}

	if secret == nil {
		return nil, s.e.ErrorOnly(ErrEmptySecret)
	}

	return json.Marshal(secret.Data["data"])
}

// GetCredentialsBytesByPath returns all secrets bytes from the specified path.
func (s *Service) GetCredentialsBytesByPath(path string) ([]byte, error) {
	secret, err := s.client.Logical().Read(path)
	if err != nil {
		return nil, s.e.ErrorOnly(err, ErrReadSecretDetail)
	}

	if secret == nil {
		return nil, s.e.ErrorOnly(ErrEmptySecret)
	}

	return json.Marshal(secret.Data["data"])
}

// GetCredentialsByPathAndKey returns secret by path and field.
func (s *Service) GetCredentialsByPathAndKey(path, key string) (string, error) {
	secret, err := s.client.Logical().Read(path)
	if err != nil {
		return "", s.e.ErrorOnly(err, ErrReadSecretDetail)
	}

	if secret == nil {
		return "", s.e.ErrorOnly(ErrEmptySecret)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return "", s.e.ErrorOnly(ErrCastSecret)
	}

	keyVal, ok := data[key]
	if !ok {
		return "", s.e.ErrorOnly(ErrNotExistingKey, key)
	}

	keyString, ok := keyVal.(string)
	if !ok {
		return "", s.e.ErrorOnly(ErrKeyType, key)
	}

	return keyString, nil
}

// GetCredentialsByPathAndKeys returns fields sets from path.
func (s *Service) GetCredentialsByPathAndKeys(path string, keys ...string) (map[string]string, error) {
	res := make(map[string]string, len(keys))

	secret, err := s.client.Logical().Read(path)
	if err != nil {
		return nil, s.e.ErrorOnly(err, ErrReadSecretDetail)
	}

	if secret == nil {
		return nil, s.e.ErrorOnly(ErrEmptySecret)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, s.e.ErrorOnly(ErrCastSecret)
	}

	for _, k := range keys {
		keyVal, isExists := data[k]
		if !isExists {
			return res, s.e.ErrorOnly(ErrNotExistingKey, k)
		}

		keyString, isExists := keyVal.(string)
		if !isExists {
			return res, s.e.ErrorOnly(ErrKeyType, k)
		}

		res[k] = keyString
	}

	return res, nil
}

func (s *Service) Login(ctx context.Context) (*vaultApi.Client, error) {
	loggedInVaultClient, err := s.clientSvc.Login(ctx)
	if err != nil {
		return nil, s.e.ErrorOnly(err)
	}

	s.client = loggedInVaultClient

	renewTTL := s.cfg.GetTokenRenewTTL()
	if renewTTL == 0 {
		return s.client, nil
	}

	renewSvc := newRenewer(s.l, s.e,
		s.clientSvc, renewTTL)

	err = renewSvc.PrepareAndStartRenew(ctx)
	if err != nil {
		return nil, s.e.ErrorOnly(err)
	}

	s.renewerSvc = renewSvc

	return s.client, nil
}

func (s *Service) GetClient() *vaultApi.Client {
	return s.client
}

func NewService(
	logger *log.Logger,
	errFmtSvc errorFormatterService,
	cfg configService,
	client clientService,
) (*Service, error) {
	return &Service{
		l: logger,
		e: errFmtSvc,

		renewerSvc: nil,
		authInfo:   nil,

		clientSvc: client,
		cfg:       cfg,

		loadedSecrets: make(map[string]string),
	}, nil
}
