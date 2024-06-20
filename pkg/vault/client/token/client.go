package token

import (
	"context"

	vaultApi "github.com/hashicorp/vault/api"
)

type service struct {
	vaultConfig *vaultApi.Config
	client      *vaultApi.Client
	cfg         configService
}

func (s *service) GetClient() *vaultApi.Client {
	return s.client
}

func (s *service) Login(ctx context.Context) (*vaultApi.Client, error) {
	return s.client, nil
}

// NewClient initialize vault client with single token
// authentication.
func NewClient(ctx context.Context, cfg configService) (*service, error) {
	clientOpts := vaultApi.DefaultConfig()
	clientOpts.Address = cfg.GetAddress()

	client, err := vaultApi.NewClient(clientOpts)
	if err != nil {
		return nil, err
	}

	client.SetToken(cfg.GetAuthToken())

	return &service{
		vaultConfig: clientOpts,
		client:      client,
		cfg:         cfg,
	}, nil
}
