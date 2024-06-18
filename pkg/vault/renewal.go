package vault

import (
	"context"
	"log"

	vault "github.com/hashicorp/vault/api"
)

type renewResult uint8

const (
	renewError renewResult = 1 << iota
	exitRequested
	expiringAuthToken
)

func (s *Service) tokenRenew(ctx context.Context) {
	for {
		renewed, err := s.renew(ctx)
		if err != nil {
			log.Fatalf("vault token renew error: %v", err)
		}

		if renewed&exitRequested != 0 {
			return
		}

		if renewed&expiringAuthToken != 0 {
			vaultClient, loginErr := s.clientSvc.Login(ctx)
			if loginErr != nil {
				log.Fatalf("login authentication error: %v", err)
			}
			s.client = vaultClient
		}
	}
}

func (s *Service) renew(ctx context.Context) (renewResult, error) {
	tokenSecret, err := s.client.Auth().Token().LookupSelf()
	if err != nil {
		return renewError, NewInternalError(ErrUnableGetTokenInfo, err)
	}

	s.authInfo = tokenSecret

	authTokenWatcher, err := s.client.NewLifetimeWatcher(&vault.LifetimeWatcherInput{
		Secret: s.authInfo,
	})
	if err != nil {
		return renewError, NewInternalError(ErrInitTokenTTLWatcher, err)
	}

	go authTokenWatcher.Start()
	defer authTokenWatcher.Stop()

	for {
		select {
		case <-ctx.Done():
			return exitRequested, nil
		case doneErr := <-authTokenWatcher.DoneCh():
			return expiringAuthToken, doneErr
		case info := <-authTokenWatcher.RenewCh():
			s.authInfo = info.Secret

			log.Printf("auth token: successfully renewed; remaining duration: %ds",
				info.Secret.Auth.LeaseDuration)
		}
	}
}
