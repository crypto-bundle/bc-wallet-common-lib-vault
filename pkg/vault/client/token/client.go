package token

import (
	"context"

	vaultApi "github.com/hashicorp/vault/api"
)

const AuthMethodName = "token"

var _ selfService = (*service)(nil)

type service struct {
	e errorFormatterService

	vaultConfig *vaultApi.Config
	client      *vaultApi.Client
	cfg         configService
}

func (s *service) GetAuthMethod() string {
	return AuthMethodName
}

func (s *service) GetClient() *vaultApi.Client {
	return s.client
}

func (s *service) Login(ctx context.Context) (*vaultApi.Client, error) {
	return s.client, nil
}

// NewClient initialize vault client with single token
// authentication.
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

	client.SetToken(cfg.GetAuthToken())

	return &service{
		e: errFmtSvc,

		vaultConfig: clientOpts,
		client:      client,
		cfg:         cfg,
	}, nil
}
