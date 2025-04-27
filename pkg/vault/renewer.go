/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2025 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

package vault

import (
	"context"
	"log/slog"
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
	l *slog.Logger
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
		s.l.Error("vault token renew error", slog.Any("error", loopErr))

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
		s.startRenew(ctx, sync.OnceFunc(func() {
			wg.Done()
		}))
	}()

	wg.Wait()

	return nil
}

func (s *renewer) startRenew(ctx context.Context, renewClb func()) {
	for {
		renewed, loopErr := s.renew(ctx, renewClb)
		if loopErr != nil {
			s.l.Error("vault token renew error", slog.Any("error", loopErr))
		}

		if renewed&exitRequested != 0 {
			return
		}

		if renewed&expiringAuthToken != 0 {
			_, loginErr := s.client.Login(ctx)
			if loginErr != nil {
				s.l.Error("login authentication error", slog.Any("error", loginErr))
			}

			s.l.Info("reconnect and renew")
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
			s.l.Error("renew auth token done with err, and chan is closed",
				slog.Any("error", doneErr), slog.Bool(RenewTokenChannelStatusTag, isClosed))

			return expiringAuthToken, s.e.ErrorNoWrap(doneErr)

		case info := <-authTokenWatcher.RenewCh():
			s.currentSecret = info.Secret

			s.l.Info("successfully renewed",
				slog.Int(RenewTokenLeaseDurationTag, info.Secret.Auth.LeaseDuration))
		}
	}
}

func newRenewer(logger *slog.Logger,
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
