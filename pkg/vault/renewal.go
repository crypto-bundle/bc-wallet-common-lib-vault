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

func (s *service) tokenRenew(ctx context.Context) {
	for {
		renewed, err := s.renew(ctx, s.authInfo)
		if err != nil {
			log.Fatalf("vault token renew error: %v", err)
		}

		if renewed&exitRequested != 0 {
			return
		}

		if renewed&expiringAuthToken != 0 {
			err = s.login(ctx)
			if err != nil {
				log.Fatalf("login authentication error: %v", err)
			}
		}
	}
}

func (s *service) renew(ctx context.Context, authToken *vault.Secret) (renewResult, error) {
	authTokenWatcher, err := s.client.NewLifetimeWatcher(&vault.LifetimeWatcherInput{
		Secret: authToken,
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
		case err := <-authTokenWatcher.DoneCh():
			return expiringAuthToken, err
		case info := <-authTokenWatcher.RenewCh():
			log.Printf("auth token: successfully renewed; remaining duration: %ds", info.Secret.Auth.LeaseDuration)
		}
	}
}
