/*
 * Copyright © 2023 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package loader

import (
	"os"

	"github.com/photowey/nemo/pkg/collection"
	"github.com/photowey/nemo/pkg/ordered"
	"github.com/photowey/nemo/pkg/stringz"
	"gopkg.in/yaml.v2"
)

// xxx.yaml | xxx.yml

const (
	Yaml = "yaml"
	Yml  = "yml"
)

const (
	yamlStep     = 10
	yamlPriority = ordered.HighPriority + yamlStep*ordered.DefaultStep
)

var (
	ymlSupportedConfigTypes = stringz.InitStringSlice(Yaml, Yml)
)

var (
	_ ConfigLoader = (*YamlConfigLoader)(nil)
)

func init() {
	Register(NewYamlConfigLoader())
}

type YamlConfigLoader struct {
	name     string
	types    collection.StringSlice
	priority int64
}

func NewYamlConfigLoader() ConfigLoader {
	return &YamlConfigLoader{
		name:     Yaml,
		types:    ymlSupportedConfigTypes,
		priority: yamlPriority,
	}
}

func (ycl *YamlConfigLoader) Supports(strategy string) bool {
	return collection.ArrayContains(ycl.types, strategy)
}

func (ycl *YamlConfigLoader) Order() int64 {
	return ycl.priority
}

func (ycl *YamlConfigLoader) Name() string {
	return ycl.name
}

func (ycl *YamlConfigLoader) Load(path string, targetPtr any) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(bytes, targetPtr)
}
