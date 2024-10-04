package kubernates

import (
	"context"
	"errors"

	vaultApi "github.com/hashicorp/vault/api"
	k8sAuth "github.com/hashicorp/vault/api/auth/kubernetes"
)

var (
	ErrEmptySecret = errors.New("unable to get secret")
)

type service struct {
	e   errorFormatterService
	cfg configService

	vaultConfig *vaultApi.Config
	client      *vaultApi.Client
	k8sAuth     *k8sAuth.KubernetesAuth
}

func (s *service) GetClient() *vaultApi.Client {
	return s.client
}

func (s *service) Login(ctx context.Context) (*vaultApi.Client, error) {
	authInfo, err := s.client.Auth().Login(ctx, s.k8sAuth)
	if err != nil {
		return nil, s.e.ErrorOnly(err)
	}

	if authInfo == nil {
		return nil, s.e.ErrorOnly(ErrEmptySecret)
	}

	s.client.SetToken(authInfo.Auth.ClientToken)

	return s.client, nil
}

// NewClient initialize vault client with service account token authorization.
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

	auth, err := k8sAuth.NewKubernetesAuth(
		cfg.GetKubernatesAppRole(),
		k8sAuth.WithServiceAccountTokenPath(cfg.GetKubernatesSATokenPath()),
		k8sAuth.WithMountPath(cfg.GetKubernatesAuthPath()),
	)
	if err != nil {
		return nil, errFmtSvc.ErrorOnly(err)
	}

	vaultSvc := &service{
		e:   errFmtSvc,
		cfg: cfg,

		client:      client,
		k8sAuth:     auth,
		vaultConfig: nil,
	}

	return vaultSvc, nil
}
