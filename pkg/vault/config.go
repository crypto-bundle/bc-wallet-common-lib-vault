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

import "fmt"

type BaseConfig struct {
	// config dependencies
	baseAppCfgSrv baseApplicationConfigService
	e             errorFormatterService
	// config fields
	Host       string `envconfig:"VAULT_SERVICE_HOST" default:"vault"`
	DataPath   string `envconfig:"VAULT_APP_DATA_PATH" default:""`
	AuthMethod string `envconfig:"VAULT_AUTH_METHOD" default:"token"`
	// TokenRenewTTL - value in seconds of increment token expiration time.
	// If you want to disable renewal flow - you can set zero value
	TokenRenewTTL int `envconfig:"VAULT_AUTH_TOKEN_RENEW_TTL" default:"240"`

	Port uint32 `envconfig:"VAULT_SERVICE_PORT" default:"8200"`

	UseHTTPS bool `envconfig:"VAULT_USE_HTTPS" default:"true"`
}

func (c *BaseConfig) GetAddress() string {
	protocol := "http"
	if c.UseHTTPS {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s:%d", protocol, c.Host, c.Port)
}

func (c *BaseConfig) GetHost() string {
	return c.Host
}

func (c *BaseConfig) GetPort() uint32 {
	return c.Port
}

func (c *BaseConfig) IsUseHTTPS() bool {
	return c.UseHTTPS
}

func (c *BaseConfig) GetDataPath() string {
	return c.DataPath
}

func (c *BaseConfig) GetAuthMethod() string {
	return c.AuthMethod
}

func (c *BaseConfig) GetTokenRenewTTL() int {
	return c.TokenRenewTTL
}

func (c *BaseConfig) GetApplicationStageName() string {
	return c.baseAppCfgSrv.GetStageName()
}

func (c *BaseConfig) GetApplicationEnvironment() string {
	return c.baseAppCfgSrv.GetEnvironmentName()
}

func (c *BaseConfig) Prepare() error {
	return nil
}

func (c *BaseConfig) PrepareWith(dependentCfgList ...interface{}) error {
	for _, cfgSrv := range dependentCfgList {
		switch castedDep := cfgSrv.(type) {
		case baseApplicationConfigService:
			c.baseAppCfgSrv = castedDep
		case errorFormatterService:
			c.e = castedDep
		default:
			continue
		}
	}

	return nil
}
