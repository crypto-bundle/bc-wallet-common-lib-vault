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
	l *log.Logger
	e errorFormatterService

	client        clientService
	currentSecret *vaultApi.Secret

	defaultTTL int
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

func (s *renewer) prepareRenew(_ context.Context) error {
	token := s.client.GetClient().Auth().Token()

	secret, loopErr := token.RenewSelf(s.defaultTTL)
	if loopErr != nil {
		s.l.Printf("vault token renew error: %v", loopErr)

		return s.e.ErrorOnly(loopErr, ErrUnableGetTokenInfoDetail)
	}

	s.currentSecret = secret

	return nil
}

func (s *renewer) PrepareAndStartRenew(ctx context.Context) error {
	err := s.prepareRenew(ctx)
	if err != nil {
		return s.e.ErrorNoWrap(err)
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
		return s.e.ErrorNoWrap(err)
	}

	return nil
}

func (s *renewer) startRenew(ctx context.Context, renewClb func()) error {
	for {
		renewed, loopErr := s.renew(ctx, renewClb)
		if loopErr != nil {
			s.l.Printf("vault token renew error: %v", loopErr)
		}

		if renewed&exitRequested != 0 {
			return nil
		}

		if renewed&expiringAuthToken != 0 {
			_, loginErr := s.client.Login(ctx)
			if loginErr != nil {
				s.l.Printf("login authentication error: %v", loginErr)
			}

			s.l.Printf("reconnect and renew")
		}
	}
}

func (s *renewer) renew(ctx context.Context, renewClb func()) (renewResult, error) {
	authTokenWatcher, err := s.client.GetClient().NewLifetimeWatcher(&vaultApi.LifetimeWatcherInput{
		Secret:        s.currentSecret,
		Grace:         0,
		Rand:          nil,
		RenewBuffer:   0,
		Increment:     0,
		RenewBehavior: 0,
	})
	if err != nil {
		return renewError, s.e.ErrorOnly(err, ErrInitTokenTTLWatcherDetail)
	}

	go authTokenWatcher.Start()
	defer authTokenWatcher.Stop()

	renewClb()

	for {
		select {
		case <-ctx.Done():
			return exitRequested, nil

		case doneErr, isClosed := <-authTokenWatcher.DoneCh():
			s.l.Printf("auth token: done with err: %s and chan is closed: %t",
				doneErr, isClosed)

			return expiringAuthToken, s.e.ErrorNoWrap(doneErr)

		case info := <-authTokenWatcher.RenewCh():
			s.currentSecret = info.Secret

			s.l.Printf("auth token: successfully renewed; remaining duration: %ds",
				info.Secret.Auth.LeaseDuration)
		}
	}
}

func newRenewer(logger *log.Logger,
	errFmtSvc errorFormatterService,
	clientSvc clientService,
	defaultTTL int,
) *renewer {
	return &renewer{
		l: logger,
		e: errFmtSvc,

		client:        clientSvc,
		defaultTTL:    defaultTTL,
		currentSecret: nil,
	}
}
