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
	"encoding/base64"

	"github.com/hashicorp/vault/api"
)

var (
	_ Encryptor = (*encryptor)(nil)
)

const (
	plainTxt    = "plaintext"
	cipherTxt   = "ciphertext"
	encryptPath = "transit/encrypt/"
	decryptPath = "transit/decrypt/"
)

type encryptor struct {
	e errorFormatterService

	client *api.Client

	transitKey string
}

func (s *encryptor) IsHealed(_ context.Context) bool {
	status, err := s.client.Sys().Health()
	if err != nil {
		return false
	}

	return status.Standby && status.Sealed
}

// Encrypt get encrypted ciphertext bytes via vault transit secret engine.
func (s *encryptor) Encrypt(toEncrypt []byte) ([]byte, error) {
	b64Val := base64.StdEncoding.EncodeToString(toEncrypt)
	path := encryptPath + s.transitKey

	secret, err := s.client.Logical().Write(path, map[string]interface{}{
		plainTxt: b64Val,
	})
	if err != nil {
		return nil, s.e.ErrorOnly(err)
	}

	encrVal := secret.Data[cipherTxt]

	encrValStr, ok := encrVal.(string)
	if !ok {
		return nil, s.e.ErrorOnly(ErrTransitSecretFormat)
	}

	return []byte(encrValStr), nil
}

// Decrypt get decrypted value from ciphertext via vault transit secret engine.
func (s *encryptor) Decrypt(cipherBytes []byte) ([]byte, error) {
	path := decryptPath + s.transitKey

	secret, err := s.client.Logical().Write(path, map[string]interface{}{
		cipherTxt: string(cipherBytes),
	})
	if err != nil {
		return nil, s.e.ErrorOnly(err)
	}

	decrVal := secret.Data[plainTxt]

	decrValStr, ok := decrVal.(string)
	if !ok {
		return nil, s.e.ErrorOnly(ErrTransitSecretFormat)
	}

	rawData, err := base64.StdEncoding.DecodeString(decrValStr)
	if err != nil {
		return nil, s.e.ErrorOnly(err)
	}

	return rawData, nil
}

func NewEncryptService(errFmtSvc errorFormatterService,
	client clientService,
	transitKey string,
) *encryptor {
	return &encryptor{
		e: errFmtSvc,

		transitKey: transitKey,
		client:     client.GetClient(),
	}
}
