package vault

import (
	"context"
	vault "github.com/hashicorp/vault/api"
	"sync"
	"time"
)

type renewResult uint8

const (
	renewError renewResult = 1 << iota
	exitRequested
	expiringAuthToken
)

func (s *Service) prepareRenew(ctx context.Context) error {
	token := s.client.Auth().Token()

	secret, loopErr := token.RenewSelf(s.cfg.GetTokenRenewTTL())
	if loopErr != nil {
		s.logger.Printf("vault token renew error: %v", loopErr)

		return NewInternalError(ErrUnableGetTokenInfo, loopErr)
	}

	s.authInfo = secret

	return nil
}

func (s *Service) prepareAndStartRenew(ctx context.Context) error {
	err := s.prepareRenew(ctx)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err = s.startRenew(ctx, wg)
	}()
	wg.Wait()

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) startRenew(ctx context.Context, wg *sync.WaitGroup) error {
	for {
		renewed, loopErr := s.renew(ctx, wg)
		if loopErr != nil {
			s.logger.Printf("vault token renew error: %v", loopErr)
		}

		if renewed&exitRequested != 0 {
			return nil
		}

		if renewed&expiringAuthToken != 0 {
			vaultClient, loginErr := s.clientSvc.Login(ctx)
			if loginErr != nil {
				s.logger.Printf("login authentication error: %v", loginErr)
			}

			s.client = vaultClient
			s.logger.Printf("reconnect and renew")
		}

		time.Sleep(time.Second * 2)
	}
}

func (s *Service) renew(ctx context.Context, wg *sync.WaitGroup) (renewResult, error) {
	authTokenWatcher, err := s.client.NewLifetimeWatcher(&vault.LifetimeWatcherInput{
		Secret: s.authInfo,
		//Increment: s.cfg.GetTokenRenewTTL(),
	})
	if err != nil {
		return renewError, NewInternalError(ErrInitTokenTTLWatcher, err)
	}

	go authTokenWatcher.Start()
	defer authTokenWatcher.Stop()

	wg.Done()

	for {
		select {
		case <-ctx.Done():
			return exitRequested, nil

		case doneErr, isClosed := <-authTokenWatcher.DoneCh():
			s.logger.Printf("auth token: done with err: %s and chan is closed: %t",
				doneErr, isClosed)

			return expiringAuthToken, doneErr

		case info := <-authTokenWatcher.RenewCh():
			s.authInfo = info.Secret

			s.logger.Printf("auth token: successfully renewed; remaining duration: %ds",
				info.Secret.Auth.LeaseDuration)
		}
	}
}
