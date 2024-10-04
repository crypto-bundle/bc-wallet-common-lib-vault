package github

import (
	"context"
	"errors"

	vaultApi "github.com/hashicorp/vault/api"
)

var (
	ErrEmptySecret = errors.New("unable to get secret")
)

type service struct {
	e   errorFormatterService
	cfg configService

	vaultConfig *vaultApi.Config
	client      *vaultApi.Client
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
		return nil, s.e.ErrorOnly(err)
	}

	if secret == nil {
		return nil, s.e.ErrorOnly(ErrEmptySecret)
	}

	s.client.SetToken(secret.Auth.ClientToken)

	return s.client, nil
}

// NewClient initialize vault client with github_token authorization.
func NewClient(_ context.Context,
	errFmtSvc errorFormatterService,
	cfg configService,
) (*service, error) {
	clientOpts := vaultApi.DefaultConfig()
	clientOpts.Address = cfg.GetAddress()

	client, err := vaultApi.NewClient(clientOpts)
	if err != nil {
		return nil, errFmtSvc.ErrorOnly(err)
	}

	return &service{
		e:           errFmtSvc,
		vaultConfig: nil,
		client:      client,
		cfg:         cfg,
	}, nil
}
