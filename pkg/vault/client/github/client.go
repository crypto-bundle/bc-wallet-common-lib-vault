package github

import (
	"errors"
	vaultApi "github.com/hashicorp/vault/api"
	"golang.org/x/net/context"
)

var (
	ErrEmptySecret = errors.New("unable to get secret")
)

const githubAuthPath = "auth/github/login"

type service struct {
	vaultConfig *vaultApi.Config
	client      *vaultApi.Client
	cfg         configService
}

func (s *service) GetClient() *vaultApi.Client {
	return s.client
}

func (s *service) Login(_ context.Context) (*vaultApi.Client, error) {
	writeData := map[string]interface{}{
		"token": s.cfg.GetGithubAuthToken(),
	}

	secret, err := s.client.Logical().Write(s.cfg.GetGithubAuthPath(), writeData)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, ErrEmptySecret
	}

	s.client.SetToken(secret.Auth.ClientToken)

	return s.client, nil
}

// NewClient initialize vault client with github_token authorization.
func NewClient(_ context.Context, cfg configService) (*service, error) {
	clientOpts := vaultApi.DefaultConfig()
	clientOpts.Address = cfg.GetAddress()

	client, err := vaultApi.NewClient(clientOpts)
	if err != nil {
		return nil, err
	}

	return &service{
		client: client,
		cfg:    cfg,
	}, nil
}
