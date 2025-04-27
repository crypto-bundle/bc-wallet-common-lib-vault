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
	"encoding/json"
	"log/slog"
	"time"

	vaultApi "github.com/hashicorp/vault/api"
)

var (
	_ Vaulter = (*Service)(nil)
)

type Service struct {
	l *slog.Logger
	e errorFormatterService

	client     *vaultApi.Client
	clientSvc  clientService
	renewerSvc renewerService
	authInfo   *vaultApi.Secret
	cfg        configService

	loadedSecrets map[string]string
}

func (s *Service) IsHealed(ctx context.Context) bool {
	status, err := s.client.Sys().Health()
	if err != nil {
		return false
	}

	serverOk := status.Standby && status.Sealed
	if !serverOk {
		return false
	}

	isHealed := s.renewerSvc.IsHealed(ctx)
	if !isHealed {
		return isHealed
	}

	secretData, err := s.client.Auth().Token().LookupSelf()
	if err != nil {
		return false
	}

	currentTime := time.Now()

	return currentTime.Unix() > int64(secretData.LeaseDuration)
}

func (s *Service) GetAuthMethod() string {
	return s.clientSvc.GetAuthMethod()
}

// GetCredentialsBytes returns all secrets bytes from default path.
func (s *Service) GetCredentialsBytes() ([]byte, error) {
	secret, err := s.client.Logical().Read(s.cfg.GetDataPath())
	if err != nil {
		return nil, s.e.ErrorOnly(err, ErrReadSecretDetail)
	}

	if secret == nil {
		return nil, s.e.ErrorOnly(ErrEmptySecret)
	}

	rawJSONData, err := json.Marshal(secret.Data["data"])
	if err != nil {
		return nil, s.e.ErrorOnly(err)
	}

	return rawJSONData, nil
}

// GetCredentialsBytesByPath returns all secrets bytes from the specified path.
func (s *Service) GetCredentialsBytesByPath(path string) ([]byte, error) {
	secret, err := s.client.Logical().Read(path)
	if err != nil {
		return nil, s.e.ErrorOnly(err, ErrReadSecretDetail)
	}

	if secret == nil {
		return nil, s.e.ErrorOnly(ErrEmptySecret)
	}

	rawJSONData, err := json.Marshal(secret.Data["data"])
	if err != nil {
		return nil, s.e.ErrorOnly(err)
	}

	return rawJSONData, nil
}

// GetCredentialsByPathAndKey returns secret by path and field.
func (s *Service) GetCredentialsByPathAndKey(path, key string) (string, error) {
	secret, err := s.client.Logical().Read(path)
	if err != nil {
		return "", s.e.ErrorOnly(err, ErrReadSecretDetail)
	}

	if secret == nil {
		return "", s.e.ErrorOnly(ErrEmptySecret)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return "", s.e.ErrorOnly(ErrCastSecret)
	}

	keyVal, ok := data[key]
	if !ok {
		return "", s.e.ErrorOnly(ErrNotExistingKey, key)
	}

	keyString, ok := keyVal.(string)
	if !ok {
		return "", s.e.ErrorOnly(ErrKeyType, key)
	}

	return keyString, nil
}

// GetCredentialsByPathAndKeys returns fields sets from path.
func (s *Service) GetCredentialsByPathAndKeys(path string, keys ...string) (map[string]string, error) {
	res := make(map[string]string, len(keys))

	secret, err := s.client.Logical().Read(path)
	if err != nil {
		return nil, s.e.ErrorOnly(err, ErrReadSecretDetail)
	}

	if secret == nil {
		return nil, s.e.ErrorOnly(ErrEmptySecret)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, s.e.ErrorOnly(ErrCastSecret)
	}

	for _, key := range keys {
		keyVal, isExists := data[key]
		if !isExists {
			return res, s.e.ErrorOnly(ErrNotExistingKey, key)
		}

		keyString, isExists := keyVal.(string)
		if !isExists {
			return res, s.e.ErrorOnly(ErrKeyType, key)
		}

		res[key] = keyString
	}

	return res, nil
}

func (s *Service) Login(ctx context.Context) (*vaultApi.Client, error) {
	loggedInVaultClient, err := s.clientSvc.Login(ctx)
	if err != nil {
		return nil, s.e.ErrorOnly(err)
	}

	s.client = loggedInVaultClient

	renewTTL := s.cfg.GetTokenRenewTTL()
	if renewTTL == 0 {
		return s.client, nil
	}

	renewSvc := newRenewer(s.l, s.e,
		s.clientSvc, renewTTL)

	err = renewSvc.PrepareAndStartRenew(ctx)
	if err != nil {
		return nil, s.e.ErrorOnly(err)
	}

	s.renewerSvc = renewSvc

	return s.client, nil
}

func (s *Service) GetClient() *vaultApi.Client {
	return s.client
}

func NewService(logBuilder loggerFabricService,
	errFmtSvc errorFormatterService,
	cfg configService,
	client clientService,
) (*Service, error) {
	return &Service{
		l: logBuilder.NewSlogNamedLoggerEntry("vault",
			slog.String(AuthMethodNameTag, client.GetAuthMethod())),
		e: errFmtSvc,

		renewerSvc: nil,
		authInfo:   nil,

		clientSvc: client,
		client:    nil,

		cfg: cfg,

		loadedSecrets: make(map[string]string),
	}, nil
}
