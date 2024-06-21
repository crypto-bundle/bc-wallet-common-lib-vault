package vault

import (
	"context"
	"log"
	"sync"

	vaultApi "github.com/hashicorp/vault/api"
)

type renewResult uint8

const (
	renewError renewResult = 1 << iota
	exitRequested
	expiringAuthToken
)

type renewer struct {
	logger *log.Logger

	client     clientService
	defaultTTL int

	currentSecret *vaultApi.Secret
}

func (s *renewer) IsHealed(_ context.Context) bool {
	isRenewable, err := s.currentSecret.TokenIsRenewable()
	if err != nil {
		return false
	}

	if !isRenewable {
		return false
	}

	ttl, err := s.currentSecret.TokenTTL()
	if err != nil {
		return false
	}

	if ttl == 0 {
		return false
	}

	return true
}

func (s *renewer) prepareRenew(ctx context.Context) error {
	token := s.client.GetClient().Auth().Token()

	secret, loopErr := token.RenewSelf(s.defaultTTL)
	if loopErr != nil {
		s.logger.Printf("vault token renew error: %v", loopErr)

		return NewInternalError(ErrUnableGetTokenInfo, loopErr)
	}

	s.currentSecret = secret

	return nil
}

func (s *renewer) PrepareAndStartRenew(ctx context.Context) error {
	err := s.prepareRenew(ctx)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err = s.startRenew(ctx, sync.OnceFunc(func() {
			wg.Done()
		}))
	}()
	wg.Wait()

	if err != nil {
		return err
	}

	return nil
}

func (s *renewer) startRenew(ctx context.Context, renewClb func()) error {
	for {
		renewed, loopErr := s.renew(ctx, renewClb)
		if loopErr != nil {
			s.logger.Printf("vault token renew error: %v", loopErr)
		}

		if renewed&exitRequested != 0 {
			return nil
		}

		if renewed&expiringAuthToken != 0 {
			_, loginErr := s.client.Login(ctx)
			if loginErr != nil {
				s.logger.Printf("login authentication error: %v", loginErr)
			}

			s.logger.Printf("reconnect and renew")
		}
	}
}

func (s *renewer) renew(ctx context.Context, renewClb func()) (renewResult, error) {
	authTokenWatcher, err := s.client.GetClient().NewLifetimeWatcher(&vaultApi.LifetimeWatcherInput{
		Secret: s.currentSecret,
	})
	if err != nil {
		return renewError, NewInternalError(ErrInitTokenTTLWatcher, err)
	}

	go authTokenWatcher.Start()
	defer authTokenWatcher.Stop()

	renewClb()

	for {
		select {
		case <-ctx.Done():
			return exitRequested, nil

		case doneErr, isClosed := <-authTokenWatcher.DoneCh():
			s.logger.Printf("auth token: done with err: %s and chan is closed: %t",
				doneErr, isClosed)

			return expiringAuthToken, doneErr

		case info := <-authTokenWatcher.RenewCh():
			s.currentSecret = info.Secret

			s.logger.Printf("auth token: successfully renewed; remaining duration: %ds",
				info.Secret.Auth.LeaseDuration)
		}
	}
}

func newRenewer(logger *log.Logger,
	clientSvc clientService,
	defaultTTL int,
) *renewer {
	return &renewer{
		logger:        logger,
		client:        clientSvc,
		defaultTTL:    defaultTTL,
		currentSecret: nil,
	}
}
