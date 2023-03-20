package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const vaultVarPrefixDelimiter = "-" // delimiter between vault variable name and the prefix

const vaultDataPathDelimiter = "," // vault DataPath env variable can be a list of paths, separated by delimiter

const vaultTestDefaultPath = "test" // for testing purposes

const vaultPrefixEnvName = "ENVIRONMENT" // name of env variable which may contain vault variable prefix

func (s *Service) GetByName(keyName string) (data string, isExists bool) {
	secretData, isExists := s.loadedSecrets[keyName]
	if !isExists {
		return "", false
	}

	return secretData, true
}

func (s *Service) LoadSecrets(_ context.Context) error {
	finalVars := make(map[string]string)

	paths := strings.Split(s.cfg.GetDataPath(), vaultDataPathDelimiter)

	vaultPrefix, isExists := os.LookupEnv(vaultPrefixEnvName)
	if !isExists {
		vaultPrefix = "dev"
	}

	for _, v := range paths {
		path := strings.TrimSpace(v)
		if path == "" {
			continue
		}

		b, err := s.GetCredentialsBytesByPath(path)
		if err != nil {
			return fmt.Errorf("get bytes from path %s: %w", path, err)
		}

		vars := make(map[string]string)
		err = json.Unmarshal(b, &vars)
		if err != nil {
			return fmt.Errorf("unmarshal path from %s: %w", path, err)
		}

		currentVars := make(map[string]string)
		for k, vv := range vars {
			fk := strings.TrimPrefix(k, vaultPrefix+vaultVarPrefixDelimiter)

			_, existsInFinalVars := finalVars[k]
			_, existsInFinalVarsWithoutPrefix := finalVars[fk]

			if !existsInFinalVars && !existsInFinalVarsWithoutPrefix {
				currentVars[k] = vv
			}
		}

		for k, vv := range currentVars {
			finalVars[k] = vv
		}
	}

	s.loadedSecrets = finalVars

	return nil
}
