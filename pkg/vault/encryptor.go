package vault

import (
	"context"
	"encoding/base64"

	"github.com/hashicorp/vault/api"
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
