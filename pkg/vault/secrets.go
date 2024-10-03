package vault

import (
	"context"
	"encoding/json"
	"strings"
)

const (
	vaultVarPrefixDelimiter = "-" // delimiter between vault variable name and the prefix
	vaultDataPathDelimiter  = "," // vault DataPath env variable can be a list of paths, separated by delimiter
)

func (s *Service) GetByName(keyName string) (string, bool) {
	secretData, isExists := s.loadedSecrets[keyName]
	if !isExists {
		return "", false
	}

	return secretData, true
}

func (s *Service) LoadSecrets(_ context.Context) error {
	finalVars := make(map[string]string)

	paths := strings.Split(s.cfg.GetDataPath(), vaultDataPathDelimiter)

	vaultPrefix := s.cfg.GetApplicationStageName()

	for _, v := range paths {
		path := strings.TrimSpace(v)
		if path == "" {
			continue
		}

		b, err := s.GetCredentialsBytesByPath(path)
		if err != nil {
			return s.e.ErrorNoWrap(err)
		}

		vars := make(map[string]string)

		err = json.Unmarshal(b, &vars)
		if err != nil {
			return s.e.ErrorOnly(err)
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
