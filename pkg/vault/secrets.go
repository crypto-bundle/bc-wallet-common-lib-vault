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

	for _, vaultBucketPath := range paths {
		path := strings.TrimSpace(vaultBucketPath)
		if path == "" {
			continue
		}

		rawBytes, err := s.GetCredentialsBytesByPath(path)
		if err != nil {
			return s.e.ErrorNoWrap(err)
		}

		vars := make(map[string]string)

		err = json.Unmarshal(rawBytes, &vars)
		if err != nil {
			return s.e.ErrorOnly(err)
		}

		currentVars := make(map[string]string)

		for key, vaultValue := range vars {
			fk := strings.TrimPrefix(key, vaultPrefix+vaultVarPrefixDelimiter)

			_, existsInFinalVars := finalVars[key]
			_, existsInFinalVarsWithoutPrefix := finalVars[fk]

			if !existsInFinalVars && !existsInFinalVarsWithoutPrefix {
				currentVars[key] = vaultValue
			}
		}

		for k, vv := range currentVars {
			finalVars[k] = vv
		}
	}

	s.loadedSecrets = finalVars

	return nil
}
