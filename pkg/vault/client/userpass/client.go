package userpass

import (
	"context"
	"errors"

	vaultApi "github.com/hashicorp/vault/api"
	userpassAuth "github.com/hashicorp/vault/api/auth/userpass"
)

var (
	ErrEmptySecret = errors.New("unable to get secret")
)

type service struct {
	cfg configService

	auth        *userpassAuth.UserpassAuth
	vaultConfig *vaultApi.Config
	client      *vaultApi.Client
}

func (s *service) GetClient() *vaultApi.Client {
	return s.client
}

func (s *service) Login(ctx context.Context) (*vaultApi.Client, error) {
	secret, err := s.client.Auth().Login(ctx, s.auth)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, ErrEmptySecret
	}

	s.client.SetToken(secret.Auth.ClientToken)

	return s.client, nil
}

// NewClient initialize vault client with user and password
// authentication.
func NewClient(ctx context.Context, cfg configService) (*service, error) {
	clientOpts := vaultApi.DefaultConfig()
	clientOpts.Address = cfg.GetAddress()
	client, err := vaultApi.NewClient(clientOpts)
	if err != nil {
		return nil, err
	}

	auth, err := userpassAuth.NewUserpassAuth(cfg.GetUserName(), &userpassAuth.Password{
		FromString: cfg.GetUserPassword(),
	})
	if err != nil {
		return nil, err
	}

	return &service{
		client: client,
		cfg:    cfg,
		auth:   auth,
	}, nil
}
